package db

import (
	"runtime"
	"time"

	"stop-checker.com/db/model"
	"stop-checker.com/gtfs"
)

func NewDatabaseFromFilesystem(path string, t time.Time) (*Database, *model.Dataset) {
	// input
	input, err := gtfs.FileInput(path)
	if err != nil {
		panic(err)
	}

	// read the dataset
	reader := &gtfs.CSVReader{}
	raw, err := reader.ReadDataset(input)
	if err != nil {
		panic(err)
	}

	// parse the dataset
	parser := &gtfs.CSVParser{
		ParserFilter: gtfs.NewCutoffFilter(time.Now()),
		TZ:           time.Local,
		TimeLayout:   "15:04:05",
		DateLayout:   "20060102",
	}

	dataset := parser.ParseDataset(raw)

	// indexes
	database := NewDatabase(dataset, time.Local)
	runtime.GC()

	return database, dataset
}
