package db

import (
	"fmt"
	"sort"

	"github.com/uber/h3-go"
	"stop-checker.com/db/model"
)

type ResolutionConfig struct {
	Level      int
	EdgeLength float64 // in meters
}

type StopLocationResult struct {
	model.Stop
	Distance float64
}

type StopLocationIndex struct {
	resolution ResolutionConfig
	index      *InvertedIndex[model.Stop]
}

func NewStopLocationIndex(indexes *BaseIndex, base *model.Base, resolution ResolutionConfig) *StopLocationIndex {
	return &StopLocationIndex{
		resolution: resolution,
		index: NewInvertedIndex("location", base.Stops, func(stop model.Stop) string {
			key := h3.FromGeo(h3.GeoCoord(stop.Location), resolution.Level)
			return fmt.Sprintf("%#x", key)
		}),
	}
}

func (s *StopLocationIndex) Query(origin model.Location, searchRadius float64) []StopLocationResult {
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

	filtered := []StopLocationResult{}

	for _, stop := range stops {
		// check the distance
		distance := origin.Distance(stop.Location)
		if searchRadius < distance {
			continue
		}

		// add the stop
		filtered = append(filtered, StopLocationResult{
			Stop:     stop,
			Distance: distance,
		})
	}

	// sort by ascending distance
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Distance < filtered[j].Distance
	})

	return filtered
}
