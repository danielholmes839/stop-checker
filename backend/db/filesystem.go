package db

import (
	"time"

	"stop-checker.com/model"
	"stop-checker.com/model/gtfs"
)

func NewDatabaseFromFilesystem(path string) (*Database, *model.Base) {
	datasetRaw, err := gtfs.NewDatasetFromFilesystem(path)
	if err != nil {
		panic(err)
	}

	tz, err := time.LoadLocation("America/Toronto")
	if err != nil {
		panic(err)
	}

	octranspo := &model.DatasetParser{
		TimeZone:   tz,
		TimeLayout: "15:04:05",
		DateLayout: "20060102",
	}

	dataset := model.NewDatasetFromGTFS(datasetRaw, octranspo)
	database := NewDatabase(dataset, tz)

	return database, dataset
}
