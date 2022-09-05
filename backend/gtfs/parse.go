package gtfs

import (
	"io"
	"os"
	fp "path/filepath"

	csvtag "github.com/artonge/go-csv-tag/v2"
)

type raw struct {
	calendars     io.ReadCloser // service
	calendarDates io.ReadCloser // service exceptions
	routes        io.ReadCloser
	stoptimes     io.ReadCloser
	stops         io.ReadCloser
	trips         io.ReadCloser
}

func RawFilesystem(path string) (*raw, error) {
	calendars, err := os.Open(fp.Join(path, "calendar.txt"))
	if err != nil {
		return nil, err
	}

	calendarDates, err := os.Open(fp.Join(path, "calendar_dates.txt"))
	if err != nil {
		return nil, err
	}

	routes, err := os.Open(fp.Join(path, "routes.txt"))
	if err != nil {
		return nil, err
	}

	stoptimes, err := os.Open(fp.Join(path, "stop_times.txt"))
	if err != nil {
		return nil, err
	}

	stops, err := os.Open(fp.Join(path, "stops.txt"))
	if err != nil {
		return nil, err
	}

	trips, err := os.Open(fp.Join(path, "trips.txt"))
	if err != nil {
		return nil, err
	}

	return &raw{
		calendars:     calendars,
		calendarDates: calendarDates,
		routes:        routes,
		stoptimes:     stoptimes,
		stops:         stops,
		trips:         trips,
	}, nil
}

func RawS3() *raw {
	return nil
}

func parseCSV[T any](input io.ReadCloser) []T {
	data := []T{}
	_ = csvtag.LoadFromReader(input, &data, csvtag.CsvOptions{Separator: ','})
	input.Close()
	return data
}

type dataset struct {
	Calendars     []Calendar
	CalendarDates []CalendarDate
	Routes        []Route
	StopTimes     []StopTime
	Stops         []Stop
	Trips         []Trip
}

func newDataset(r *raw) *dataset {
	return &dataset{
		Calendars:     parseCSV[Calendar](r.calendars),
		CalendarDates: parseCSV[CalendarDate](r.calendarDates),
		Routes:        parseCSV[Route](r.routes),
		StopTimes:     parseCSV[StopTime](r.stoptimes),
		Stops:         parseCSV[Stop](r.stops),
		Trips:         parseCSV[Trip](r.trips),
	}
}
