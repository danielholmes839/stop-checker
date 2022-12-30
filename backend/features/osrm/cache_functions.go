package osrm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/rs/zerolog/log"
	"stop-checker.com/db/model"
	"stop-checker.com/db/repository"
)

type cacheData map[string]model.Path

func PrepareCacheData(
	client *Client,
	radius float64,
	stops []model.Stop,
	locationIndex repository.StopLocationSearch,
) cacheData {
	t0 := time.Now()
	paths := cacheData{}

	for _, origin := range stops {
		neighbors := locationIndex.Query(origin.Location, radius)

		for _, destination := range neighbors {
			// check both origin:destination or destination:origin
			key1 := fmt.Sprintf("%s:%s", origin.Id, destination.Id)
			_, ok1 := paths[key1]

			key2 := fmt.Sprintf("%s:%s", destination.Id, origin.Id)
			_, ok2 := paths[key2]

			// ignore this pair of stops
			if ok1 || ok2 || origin.Id == destination.Id {
				continue
			}

			path, err := client.GetDirections(origin.Location, destination.Location)
			if err != nil {
				panic(err)
			}

			paths[key1] = path
		}
	}

	log.Info().
		Dur("duration", time.Since(t0)).
		Float64("radius", radius).
		Int("paths", len(paths)).Msg("osrm-prepared-cache")

	return paths
}

func SaveCacheData(filename string, data cacheData) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, bytes, 0644)
	if err != nil {
		return err
	}

	log.Info().
		Int("bytes", len(bytes)).
		Msg("osrm-saved-cache")

	return nil
}

func ReadCacheData(filename string) (cacheData, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	data := &cacheData{}
	if err := json.Unmarshal(bytes, data); err != nil {
		return nil, err
	}

	log.Info().Msg("osrm-read-cache")

	return *data, nil
}
