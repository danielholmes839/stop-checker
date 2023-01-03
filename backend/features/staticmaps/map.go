package staticmaps

import (
	"stop-checker.com/db/model"
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
