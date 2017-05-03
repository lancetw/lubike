package ubike

import (
	"fmt"
	"sort"

	"github.com/lancetw/lubike/geo"
)

// Station struct
type Station struct {
	ID       string  `json:"-"`
	Name     string  `json:"station"`
	Lat      float64 `json:"-"`
	Lng      float64 `json:"-"`
	Distance float64 `json:"-"`
	NumUbike int     `json:"num_ubike"`
	IsFull   bool    `json:"-"`
	//UpdateTime time.Time `json:"-"`
}

// ByDistance : implements sort.Interface for []Station
type ByDistance []Station

func (p ByDistance) Len() int      { return len(p) }
func (p ByDistance) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p ByDistance) Less(i, j int) bool {
	return p[i].Distance < p[j].Distance || isNaN(p[i].Distance) && !isNaN(p[j].Distance)
}
func isNaN(f float64) bool {
	return f != f
}

// UpdateDistance : update distance field of slice
func UpdateDistance(center geo.LatLng, ubikeInfo []Station) []Station {
	if len(ubikeInfo) == 0 {
		return ubikeInfo
	}
	for index, item := range ubikeInfo {
		dist := geo.Haversine(center.Lat, center.Lng, item.Lat, item.Lng)
		ubikeInfo[index].Distance = dist
	}

	sort.Sort(ByDistance(ubikeInfo))

	return ubikeInfo
}

// UpdateDistanceByRouteMatrix : update distance field of slice
func UpdateDistanceByRouteMatrix(center geo.LatLng, ubikeInfo []Station, limit int) []Station {
	points := []string{}
	points = append(points, fmt.Sprintf("%f,%f", center.Lat, center.Lng))

	if limit >= 0 && len(ubikeInfo) >= limit {
		ubikeInfo = ubikeInfo[:limit]
	}

	for _, item := range ubikeInfo {
		point := fmt.Sprintf("%f,%f", item.Lat, item.Lng)
		points = append(points, point)
	}
	distance := geo.GoogleMapMatrixProvider(points).RouteMatrix()
	if len(distance) == 0 {
		return ubikeInfo
	}

	for index, dist := range distance {
		ubikeInfo[index].Distance = dist.(float64)
	}

	if len(ubikeInfo) >= limit {
		ubikeInfo = ubikeInfo[:limit]
	}

	sort.Sort(ByDistance(ubikeInfo))

	return ubikeInfo
}
