package ubike

import (
	"fmt"
	"testing"

	"github.com/lancetw/lubike/geo"
)

func TestIsNaN(t *testing.T) {
	isNaN(3.14)
}

func TestUpdateDistance(t *testing.T) {
	center := geo.LatLng{Lat: 1.0803, Lng: 53.9583}
	collect := []Station{
		{
			Name: "York",
			Lat:  1.0803,
			Lng:  53.9583,
		},
		{
			Name: "Bristol",
			Lat:  2.5833,
			Lng:  51.4500,
		},
	}
	collect = UpdateDistance(center, collect)

	if len(collect) != 2 {
		t.Fail()
	}

	if collect[0].Distance != 0 {
		t.Fail()
	}

	if fmt.Sprintf("%.2f", collect[1].Distance) != "296.71" {
		t.Error(collect[1].Distance)
	}
}

func TestUpdateDistanceEmptyCollect(t *testing.T) {
	var collect []Station
	center := geo.LatLng{Lat: 0.0, Lng: 0.0}
	collect = UpdateDistance(center, collect)

	if len(collect) != 0 {
		t.Fail()
	}
}
