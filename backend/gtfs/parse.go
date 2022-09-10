package gtfs

import (
	"io"
	"os"
	fp "path/filepath"
	"time"

	csvtag "github.com/artonge/go-csv-tag/v2"
	"github.com/rs/zerolog/log"
)

type raw struct {
	calendars     io.ReadCloser // service
	calendarDates io.ReadCloser // service exceptions
	routes        io.ReadCloser
	stoptimes     io.ReadCloser
	stops         io.ReadCloser
	trips         io.ReadCloser
	shapes        io.ReadCloser
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

	shapes, err := os.Open(fp.Join(path, "shapes.txt"))
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
		shapes:        shapes,
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
	Shapes        []Shape
}

func newDataset(r *raw) *dataset {
	t0 := time.Now()

	calendars := parseCSV[Calendar](r.calendars)
	calendarDates := parseCSV[CalendarDate](r.calendarDates)
	routes := parseCSV[Route](r.routes)
	stoptimes := parseCSV[StopTime](r.stoptimes)
	stops := parseCSV[Stop](r.stops)
	trips := parseCSV[Trip](r.trips)
	shapes := parseCSV[Shape](r.shapes)

	dataset := &dataset{
		Calendars:     calendars,
		CalendarDates: calendarDates,
		Routes:        routes,
		StopTimes:     stoptimes,
		Stops:         stops,
		Trips:         trips,
		Shapes:        shapes,
	}

	log.Info().
		Dur("duration", time.Since(t0)).
		Int("routes", len(routes)).
		Int("stops", len(stops)).
		Int("stoptimes", len(stoptimes)).
		Int("trips", len(trips)).
		Int("services", len(calendars)).
		Int("service-exceptions", len(calendarDates)).
		Int("shapes", len(shapes)).
		Msg("created dataset")

	return dataset
}
