package gtfs

import (
	"io"
	"os"
	fp "path/filepath"

	csvtag "github.com/artonge/go-csv-tag/v2"
)

func parseFile[T any](file string) ([]T, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return parse[T](f)
}

func parse[T any](input io.Reader) ([]T, error) {
	data := []T{}
	err := csvtag.LoadFromReader(input, &data, csvtag.CsvOptions{Separator: ','})
	return data, err
}

type Dataset struct {
	Routes        []Route
	Stops         []Stop
	StopTimes     []StopTime
	Trips         []Trip
	Calendars     []Calendar
	CalendarDates []CalendarDate
}

func NewDatasetFromFilesystem(path string) (*Dataset, error) {
	var err error
	var gtfs *Dataset = &Dataset{}

	gtfs.CalendarDates, err = parseFile[CalendarDate](fp.Join(path, "calendar_dates.txt"))
	if err != nil {
		return nil, err
	}

	gtfs.Calendars, err = parseFile[Calendar](fp.Join(path, "calendar.txt"))
	if err != nil {
		return nil, err
	}

	gtfs.Routes, err = parseFile[Route](fp.Join(path, "routes.txt"))
	if err != nil {
		return nil, err
	}

	gtfs.StopTimes, err = parseFile[StopTime](fp.Join(path, "stop_times.txt"))
	if err != nil {
		return nil, err
	}

	gtfs.Stops, err = parseFile[Stop](fp.Join(path, "stops.txt"))
	if err != nil {
		return nil, err
	}

	gtfs.Trips, err = parseFile[Trip](fp.Join(path, "trips.txt"))
	if err != nil {
		return nil, err
	}

	return gtfs, nil
}
