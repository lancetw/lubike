package ubike

import (
	"testing"

	"github.com/lancetw/lubike/geo"
)

var (
	center  = geo.LatLng{Lat: 1.0803, Lng: 53.9583}
	collect = []Station{
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
)

func TestIsNaN(t *testing.T) {
	isNaN(3.14)
}

func TestUpdateDistance(t *testing.T) {
	collect = UpdateDistance(center, collect)

	if len(collect) != 2 {
		t.Fail()
	}

	if collect[0].Distance != 0 {
		t.Fail()
	}
}

func TestUpdateDistanceEmptyCollect(t *testing.T) {
	emptyCollect := []Station{}
	emptyCenter := geo.LatLng{Lat: 0.0, Lng: 0.0}
	emptyCollect = UpdateDistance(emptyCenter, emptyCollect)

	if len(emptyCollect) != 0 {
		t.Fail()
	}
}

func TestUpdateDistanceByRouteMatrix(t *testing.T) {
	center = geo.LatLng{Lat: 25.034153, Lng: 121.568509}
	collect = []Station{
		{
			Name: "B",
			Lat:  25.0347361111,
			Lng:  121.565658333,
		},
		{
			Name: "A",
			Lat:  25.0365638889,
			Lng:  121.5686639,
		},
	}
	limit := 2
	ret := UpdateDistanceByRouteMatrix(center, collect, limit)

	if ret[0].Distance > ret[1].Distance {
		t.Fail()
	}
}

func TestUpdateDistanceByRouteMatrixCorrupted(t *testing.T) {
	center = geo.LatLng{Lat: 181, Lng: 91}
	collect = []Station{
		{
			Name: "B",
			Lat:  0,
			Lng:  0,
		},
		{
			Name: "A",
			Lat:  1,
			Lng:  1,
		},
	}
	limit := 2
	ret := UpdateDistanceByRouteMatrix(center, collect, limit)

	if ret[0].Distance > ret[1].Distance {
		t.Fail()
	}
}

func TestUpdateDistanceByRouteMatrixLimit(t *testing.T) {
	center = geo.LatLng{Lat: 181, Lng: 91}
	collect = []Station{}
	limit := -1
	ret := UpdateDistanceByRouteMatrix(center, collect, limit)

	if len(ret) > 0 {
		t.Fail()
	}

	if len(ret) > 0 && ret[0].Distance > ret[1].Distance {
		t.Fail()
	}
}
