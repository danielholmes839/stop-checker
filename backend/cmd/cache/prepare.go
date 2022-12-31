package main

import (
	"stop-checker.com/db"
	"stop-checker.com/features/osrm"
)

/* prepare cache data */
func main() {
	// read the dataset and create indexes
	database, dataset := db.NewDBFromFilesystem("./data")

	// osrm client
	client := osrm.NewClient("http://localhost:5000")

	// prepare osrm cache data
	data := osrm.PrepareCacheData(client, 300, dataset.Stops, database.StopLocationIndex)

	// save osrm cache data
	err := osrm.SaveCacheData("./data/300m-directions.json", data)

	if err != nil {
		panic(err)
	}
}
