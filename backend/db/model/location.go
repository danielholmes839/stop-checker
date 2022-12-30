package model

import "math"

type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

const earthRaidusKm float64 = 6371 // radius of the earth in kilometers.

func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

func (l Location) Distance(other Location) float64 {
	// harversine distance in meters
	lat1 := degreesToRadians(l.Latitude)
	lon1 := degreesToRadians(l.Longitude)
	lat2 := degreesToRadians(other.Latitude)
	lon2 := degreesToRadians(other.Longitude)

	diffLat := lat2 - lat1
	diffLon := lon2 - lon1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*
		math.Pow(math.Sin(diffLon/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return c * 6_371_000
}

type Path struct {
	Distance float64    `json:"distance"`
	Path     []Location `json:"path"`
}

func Distance(path ...Location) float64 {
	distance := 0.0
	n := len(path) - 1
	for i := 0; i < n; i++ {
		distance += path[i].Distance(path[i+1])
	}
	return distance
}
