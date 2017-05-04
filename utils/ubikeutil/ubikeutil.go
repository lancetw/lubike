package ubikeutil

import (
	"log"
	"strconv"
	"sync"
	"time"

	valid "github.com/asaskevich/govalidator"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/jasonwinn/geocoder"
	"github.com/lancetw/lubike/geo"
	"github.com/lancetw/lubike/setting"
	"github.com/lancetw/lubike/ubike"
	"github.com/lancetw/lubike/utils/commonutil"
)

// GeoCoordinateValidator : check lat and lng
func GeoCoordinateValidator(slat string, slng string) bool {
	status := false
	if valid.IsLatitude(slat) && valid.IsLongitude(slng) {
		status = true
	}
	return status
}

// IsAllStationsFull : check all stations fulled
func IsAllStationsFull(stations []ubike.Station) bool {
	status := true
	for _, item := range stations {
		if !item.IsFull {
			status = false
			return status
		}
	}
	return status
}

var reverseGeocode = geocoder.ReverseGeocode

// IsInCity : target coordinate is in the city
func IsInCity(latlng geo.LatLng, city string) bool {
	status := false
	config := setting.InitConfig()
	geocoder.SetAPIKey(config.MapquestAPIKey)
	address, err := reverseGeocode(latlng.Lat, latlng.Lng)
	if err != nil {
		return status
	}
	if address.State == city || address.City == city {
		status = true
	}
	return status
}

// LoadGeoCoordinate : covert strings to geo point
func LoadGeoCoordinate(slat string, slng string) (geo.LatLng, error) {
	lat, err := strconv.ParseFloat(slat, 64)
	if err != nil {
		return geo.LatLng{}, err
	}
	lng, err := strconv.ParseFloat(slng, 64)
	if err != nil {
		return geo.LatLng{}, err
	}
	return geo.LatLng{Lat: lat, Lng: lng}, err
}

// LoadNearbyUbikes : load nearby ubike stations
func LoadNearbyUbikes(slat string, slng string, num int) ([]ubike.Station, ubike.StationErrno) {
	var errno ubike.StationErrno
	ubikeInfo := []ubike.Station{}

	// check: invalid latitude or longitude
	if !GeoCoordinateValidator(slat, slng) {
		errno = ubike.InvalidLatLng
		return ubikeInfo, errno
	}

	// check: invalid latitude or longitude
	latlng, _ := LoadGeoCoordinate(slat, slng)

	// check: given location not in Taipei City
	if !IsInCity(latlng, "Taipei City") {
		errno = ubike.GivenLocationNotInTaipeiCity
		return ubikeInfo, errno
	}

	config := setting.InitConfig()
	ubikeInfo, errno = LoadUbikeInfo(config.UbikeEndpoint, config.UbikeEndpointTimeout)
	// check: all ubike stations are full
	if errno == ubike.SystemError {
		return ubikeInfo, errno
	}

	if IsAllStationsFull(ubikeInfo) {
		errno = ubike.Full
		return ubikeInfo, errno
	}
	ubikeInfo = ubike.UpdateDistance(latlng, ubikeInfo)
	limit := 7
	ubikeInfo = ubike.UpdateDistanceByRouteMatrix(latlng, ubikeInfo, limit)

	// take number of elements
	if len(ubikeInfo) >= num {
		ubikeInfo = append(ubikeInfo[:num])
	}
	return ubikeInfo, errno
}

var fetchAPIData = commonutil.FetchAPIData
var logFatalf = log.Fatalf

// LoadUbikeInfo : load remote ubike api data
func LoadUbikeInfo(endpoint string, timeout int) ([]ubike.Station, ubike.StationErrno) {
	collect := []ubike.Station{}
	errno := ubike.Ok
	endpointTimeout := time.Duration(time.Duration(timeout) * time.Second)
	json := fetchAPIData(endpoint, endpointTimeout).(*simplejson.Json)

	status := json.Get("retCode").MustInt()
	if status != 1 {
		errno = ubike.SystemError
		return collect, errno
	}

	dataset := json.Get("retVal").MustMap()

	wgStation := make(chan ubike.Station)
	wgErr := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(len(dataset))

	for _, item := range dataset {
		go func(item interface{}) {
			defer wg.Done()
			ele := item.(map[string]interface{})
			if ele["act"] == "1" && ele["sbi"] != "0" {
				sbi, err := strconv.Atoi(ele["sbi"].(string))
				if err != nil {
					wgErr <- err
					return
				}
				lat, err := strconv.ParseFloat(ele["lat"].(string), 64)
				if err != nil {
					wgErr <- err
					return
				}
				lng, err := strconv.ParseFloat(ele["lng"].(string), 64)
				if err != nil {
					wgErr <- err
					return
				}
				/*
						mdaylayout := "20060102150405"
					  updated, err := time.Parse(mdaylayout, ele["mday"].(string))
						if err != nil {
							wgErr <- err
							return
						}*/
				full := false
				if ele["bemp"] == "0" {
					full = true
				}
				station := ubike.Station{
					ID:       ele["sno"].(string),
					Name:     ele["sna"].(string),
					Lat:      lat,
					Lng:      lng,
					NumUbike: sbi,
					IsFull:   full,
					//UpdateTime: updated,
				}

				wgStation <- ubike.Station(station)
			}
		}(item)
	}

	go func() {
		wg.Wait()
		close(wgStation)
	}()

	select {
	case <-wgStation:
		for station := range wgStation {
			collect = append(collect, station)
		}
	case err := <-wgErr:
		if err != nil {
			errno = ubike.SystemError
		}
	}

	return collect, errno
}
