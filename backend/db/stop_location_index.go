package db

import (
	"fmt"

	"github.com/uber/h3-go"
	"stop-checker.com/db/model"
)

type Resolution struct {
	Level      int
	EdgeLength float64 // in meters
}

type StopLocationIndex struct {
	resolution Resolution
	index      *InvertedIndex[model.Stop]
}

func NewStopLocationIndex(stops []model.Stop, resolution Resolution) *StopLocationIndex {
	return &StopLocationIndex{
		resolution: resolution,
		index: NewInvertedIndex("location", stops, func(stop model.Stop) string {
			key := h3.FromGeo(h3.GeoCoord(stop.Location), resolution.Level)
			return fmt.Sprintf("%#x", key)
		}),
	}
}

func (s *StopLocationIndex) Query(origin model.Location, searchRadius float64) []model.StopWithDistance {
	// determine the hexs that could contain stops
	// https://observablehq.com/@nrabinowitz/h3-radius-lookup
	rings := int(searchRadius/(s.resolution.EdgeLength*2)) + 1
	originHex := h3.FromGeo(h3.GeoCoord(origin), s.resolution.Level)
	ringHexs := h3.KRing(originHex, rings)

	// store stops within the hexs
	stops := []model.Stop{}

	// stops within the rings (origin included: https://h3geo.org/docs/api/traversal/#kring)
	for _, hex := range ringHexs {
		results, _ := s.index.Get(fmt.Sprintf("%#x", hex))
		stops = append(stops, results...)
	}

	filtered := []model.StopWithDistance{}

	for _, stop := range stops {
		// check the distance
		distance := origin.Distance(stop.Location)
		if searchRadius < distance {
			continue
		}

		// add the stop
		filtered = append(filtered, model.StopWithDistance{
			Stop:     stop,
			Distance: distance,
		})
	}

	return filtered
}
