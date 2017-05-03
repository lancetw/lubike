package ubikeutil

import (
	"io/ioutil"
	"testing"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/lancetw/lubike/geo"
	"github.com/lancetw/lubike/setting"
	"github.com/lancetw/lubike/ubike"
)

const (
	lat = "25.034153"
	lng = "121.568509"
	num = 2
)

func TestIsInCity(t *testing.T) {
	latlng := geo.LatLng{Lat: 25.034153, Lng: 121.568509}
	if !IsInCity(latlng, "Taipei City") {
		t.Fail()
	}
}

func TestLoadGeoCoordinate(t *testing.T) {
	_, err := LoadGeoCoordinate(lat, lng)
	if err != nil {
		t.Error(err)
	}
	_, err = LoadGeoCoordinate("", lng)
	if err == nil {
		t.Fail()
	}
	_, err = LoadGeoCoordinate(lat, "")
	if err == nil {
		t.Fail()
	}
	_, err = LoadGeoCoordinate("", "")
	if err == nil {
		t.Fail()
	}
}

func TestGeoCoordinateValidator(t *testing.T) {
	var status bool
	status = GeoCoordinateValidator(lat, lng)
	if !status {
		t.Fail()
	}
	status = GeoCoordinateValidator("", lng)
	if status {
		t.Fail()
	}
	status = GeoCoordinateValidator(lat, "")
	if status {
		t.Fail()
	}
	status = GeoCoordinateValidator("", "")
	if status {
		t.Fail()
	}
}

func TestLoadNearbyUbikes(t *testing.T) {
	ubikeInfo, errno := LoadNearbyUbikes(lat, lng, num)

	if errno != ubike.Ok {
		t.Error(errno)
	}

	if len(ubikeInfo) != num {
		t.Error(len(ubikeInfo), num)
	}
}

func TestLoadNearbyUbikesInvalidLatLng(t *testing.T) {
	invaildLat := "lat"
	invaildLng := "lng"
	_, errno := LoadNearbyUbikes(invaildLat, invaildLng, num)

	if errno != ubike.InvalidLatLng {
		t.Error(errno)
	}
}

func TestLoadNearbyUbikesFull(t *testing.T) {
	var mockFetchAPIData = func(s string, t time.Duration) interface{} {
		body, _ := ioutil.ReadFile("./stataions_full.json")
		data, _ := simplejson.NewJson([]byte(body))
		return data
	}

	var origFetchAPIData = fetchAPIData
	fetchAPIData = mockFetchAPIData

	defer func() { fetchAPIData = origFetchAPIData }()

	_, errno := LoadNearbyUbikes(lat, lng, num)
	if errno != ubike.Full {
		t.Error(errno)
	}
}

func TestLoadNearbyUbikesGivenLocationNotInTaipeiCity(t *testing.T) {
	newTaipeiLat := "24.986779"
	newTaipeiLng := "121.3645554"
	_, errno := LoadNearbyUbikes(newTaipeiLat, newTaipeiLng, num)

	if errno != ubike.GivenLocationNotInTaipeiCity {
		t.Error(errno)
	}
}

func TestLoadNearbyUbikesSystemError(t *testing.T) {
	var mockFetchAPIData = func(s string, t time.Duration) interface{} {
		body, _ := ioutil.ReadFile("./stataions_noret.json")
		data, _ := simplejson.NewJson([]byte(body))
		return data
	}

	var origFetchAPIData = fetchAPIData
	fetchAPIData = mockFetchAPIData

	defer func() { fetchAPIData = origFetchAPIData }()

	_, errno := LoadNearbyUbikes(lat, lng, num)

	if errno != ubike.SystemError {
		t.Error(errno)
	}
}

func TestLoadUbikeInfo(t *testing.T) {
	config := setting.InitConfig()
	collect, errno := LoadUbikeInfo(config.UbikeEndpoint, config.UbikeEndpointTimeout)

	if errno != ubike.Ok {
		t.Error(errno)
	}

	if len(collect) == 0 {
		t.Fail()
	}
}
