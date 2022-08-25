package staticmaps

import (
	"math"

	"stop-checker.com/model"
)

/*
	Based on:
	https://groups.google.com/g/google-maps-js-api-v3/c/hDRO4oHVSeM/m/osOYQYXg2oUJ?pli=1

	pixels: the number of pixels from the center of the map to the edge. For example if the height is 400 this should be 200
	distance: the largest distance between the center latitude/longitude and marker latitude/longitude
	centerLatitude: the latitude of the center of the map. Required since the number of meters per pixel changes based on the latitude
*/
func calculateZoomOneWay(pixels int, distance float64, centerLatitude float64) int {
	c := 156543.03392 * math.Cos((centerLatitude*math.Pi)/180.0)
	zoom := math.Log2(float64(pixels) * c / distance)
	return int(math.Floor(zoom))
}

// calculate the required zoom to fit all markers on the map
func calculateZoom(center model.Location, markers []Marker, width, height int) int {
	maxLatitudeDiff := 0.0
	maxLongitudeDiff := 0.0

	for _, marker := range markers {
		latitudeDiff := math.Abs(marker.Latitude - center.Latitude)
		longitudeDiff := math.Abs(marker.Longitude - center.Longitude)

		if maxLatitudeDiff < latitudeDiff {
			maxLatitudeDiff = latitudeDiff
		}

		if maxLongitudeDiff < longitudeDiff {
			maxLongitudeDiff = longitudeDiff
		}
	}

	maxDiffLatitudeMeters := maxLatitudeDiff*111139.44 + 500 // 500m buffer
	maxDiffLongitudeMeters := maxLongitudeDiff*111139.44 + 500

	zoomHeight := calculateZoomOneWay(height/2, maxDiffLatitudeMeters, center.Latitude)
	zoomWidth := calculateZoomOneWay(width/2, maxDiffLongitudeMeters, center.Latitude)

	// return the minimum (most zoomed out) zoom
	if zoomHeight < zoomWidth {
		return zoomHeight
	}
	return zoomWidth
}

func calculateCenter(markers []Marker) model.Location {
	totalLat := 0.0
	totalLon := 0.0

	for _, marker := range markers {
		totalLat += marker.Latitude
		totalLon += marker.Longitude
	}

	count := float64(len(markers))

	return model.Location{
		Latitude:  totalLat / count,
		Longitude: totalLon / count,
	}
}
