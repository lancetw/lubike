package geo

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/lancetw/lubike/setting"
	"github.com/lancetw/lubike/utils/commonutil"
)

const earthRadius = float64(6371)

// LatLng struct
type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Provider struct
type Provider struct {
	name    string
	url     string
	values  map[string]interface{}
	timeout int
}

func round(f float64) float64 {
	if math.Abs(f) < 0.5 {
		return 0
	}
	return float64(int((f + math.Copysign(0.5, f))))
}

// Haversine : Use haversine to get the diatance
// brrow from https://play.golang.org/p/MZVh5bRWqN
func Haversine(lonFrom float64, latFrom float64, lonTo float64, latTo float64) (distance float64) {
	var deltaLat = (latTo - latFrom) * (math.Pi / 180)
	var deltaLon = (lonTo - lonFrom) * (math.Pi / 180)

	var a = math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(latFrom*(math.Pi/180))*math.Cos(latTo*(math.Pi/180))*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance = earthRadius * c

	return round(distance * 1000)
}

// RouteMatrix : use Directions API to get distances
func (p Provider) RouteMatrix() []interface{} {
	timeout := time.Duration(time.Duration(p.timeout) * time.Second)
	var data []interface{}
	switch p.name {
	case "mapquest":
		jsonData := commonutil.PostAPIData(p.url, timeout, p.values).(*simplejson.Json)
		rdata := jsonData.Get("distance").MustArray()[1:]
		for _, item := range rdata {
			value := string(item.(json.Number))
			dist, err := strconv.ParseFloat(value, 64)
			dist = dist * 1000
			if err != nil {
				return nil
			}
			data = append(data, dist)
		}
	case "googlemap":
		jsonData := commonutil.FetchAPIData(p.url, timeout).(*simplejson.Json)
		rows := jsonData.Get("rows").MustArray()
		if len(rows) == 0 {
			return nil
		}
		rdata := rows[0].(map[string]interface{})["elements"].([]interface{})
		for _, ele := range rdata {
			item := ele.(map[string]interface{})
			if item["status"] == "ZERO_RESULTS" || item["status"] == "NOT_FOUND" || item["status"] == "INVALID_REQUEST" {
				return nil
			}
			value := string((item["distance"]).(map[string]interface{})["value"].(json.Number))
			dist, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil
			}
			data = append(data, dist)
		}
	}
	return data
}

// MapQuestMatrixProvider : use MapQuest Directions API to get distances
// Notice: numbers of locations can not larger than 100
func MapQuestMatrixProvider(locations []string) Provider {
	var config = setting.InitConfig()
	options := map[string]interface{}{"allToAll": false, "manyToOne": false, "unit": "m"}
	url := config.MapquestRouteMatrixEndpoint + "?key=" + config.MapquestAPIKey
	values := map[string]interface{}{"locations": locations, "options": options}
	timeout := config.MapquestRouteMatrixEndpointTimeout
	return Provider{name: "mapquest", url: url, values: values, timeout: timeout}
}

// GoogleMapMatrixProvider : use MapQuest Directions API to get distances
// Notice: numbers of locations can not be larger than 100
func GoogleMapMatrixProvider(locations []string) Provider {
	var config = setting.InitConfig()
	origins := locations[0]
	destinations := strings.Join(locations[1:], "|")
	url := config.GoogleMapMatrixEndpoint + "?key=" + config.GoogleMapMatrixAPIKey + "&origins=" + origins + "&destinations=" + destinations + "&mode=walking"
	timeout := config.GoogleMapMatrixEndpointTimeout
	return Provider{name: "googlemap", url: url, timeout: timeout}
}
