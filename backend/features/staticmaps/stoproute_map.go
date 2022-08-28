package staticmaps

import (
	"errors"
	"fmt"

	"stop-checker.com/db/model"
)

func NewStopRouteMap(width, height int, stop model.Location, buses []model.Bus) (*Map, error) {
	m := NewMap(width, height)
	m.AddMarker(Marker{
		Location: stop,
		Color:    BLUE,
		Label:    "S",
		Size:     MEDIUM,
	})

	busCount := 0
	for _, bus := range buses {
		if bus.Location == nil {
			continue
		}
		busCount++
		m.AddMarker(Marker{
			Location: *bus.Location,
			Color:    RED,
			Label:    fmt.Sprint(busCount),
			Size:     MEDIUM,
		})
	}

	if busCount == 0 {
		return nil, errors.New("zero bus locations")
	}

	return m, nil
}
