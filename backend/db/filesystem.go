package db

import (
	"stop-checker.com/db/gtfs"
	"stop-checker.com/db/model"
)

func NewDatabaseFromFilesystem(path string) (*Database, *model.Base) {
	dataset, err := gtfs.NewDatasetFromFilesystem(path)
	if err != nil {
		panic(err)
	}

	octranspo := &model.BaseParser{
		TimeZone:   dataset.TimeZone,
		TimeLayout: "15:04:05",
		DateLayout: "20060102",
	}

	base := model.NewBaseFromGTFS(dataset, octranspo)
	database := NewDatabase(base)

	return database, base
}
