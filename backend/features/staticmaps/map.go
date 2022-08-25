package staticmaps

import (
	"fmt"
	"net/url"

	"stop-checker.com/model"
)

const (
	RED  = "red"
	BLUE = "blue"

	SMALL  = "small"
	MEDIUM = "mid"
)

type Marker struct {
	model.Location
	Color string
	Label string
	Size  string
}

type Map struct {
	width   int
	height  int
	markers []Marker
}

func NewMap(width, height int) *Map {
	m := &Map{
		width:   width,
		height:  height,
		markers: []Marker{},
	}

	return m
}

func (m *Map) AddMarker(marker Marker) {
	m.markers = append(m.markers, marker)
}

func (m *Map) Encode(key string) string {
	center := calculateCenter(m.markers)
	zoom := calculateZoom(center, m.markers, m.width, m.height)

	centerParamValue := fmt.Sprintf("%f,%f", center.Latitude, center.Longitude)
	sizeParamValue := fmt.Sprintf("%dx%d", m.width, m.height)

	params := url.Values{}
	params.Add("key", key)
	params.Add("center", centerParamValue)
	params.Add("size", sizeParamValue)
	params.Add("zoom", fmt.Sprint(zoom))

	for _, marker := range m.markers {
		params.Add("markers", fmt.Sprintf("color:%s|size:%s|label:%s|%f,%f",
			marker.Color,
			marker.Size,
			marker.Label,
			marker.Latitude,
			marker.Longitude,
		))
	}

	return fmt.Sprintf("https://maps.googleapis.com/maps/api/staticmap?%s", params.Encode())
}
