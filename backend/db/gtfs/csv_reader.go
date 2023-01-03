package gtfs

import (
	"io"
	"runtime"
	"time"

	csvtag "github.com/artonge/go-csv-tag/v2"
	"github.com/rs/zerolog/log"
)

// Read CSVs
func read[T any](input io.ReadCloser) []T {
	data := []T{}
	if err := csvtag.LoadFromReader(input, &data, csvtag.CsvOptions{Separator: ','}); err != nil {
		panic(err)
	}

	input.Close()
	runtime.GC() // todo: use parquet if this doesn't work
	return data
}

type CSVDataset struct {
	Calendars     []Calendar
	CalendarDates []CalendarDate
	Routes        []Route
	StopTimes     []StopTime
	Stops         []Stop
	Trips         []Trip
	Shapes        []Shape
}

type CSVReader struct {
}

func (r *CSVReader) ReadDataset(input *Input) (*CSVDataset, error) {
	t0 := time.Now()
	calendars := read[Calendar](input.Calendars)
	calendarDates := read[CalendarDate](input.CalendarDates)
	routes := read[Route](input.Routes)
	stoptimes := read[StopTime](input.Stoptimes)
	stops := read[Stop](input.Stops)
	trips := read[Trip](input.Trips)
	shapes := read[Shape](input.Shapes)

	log.Info().
		Dur("duration", time.Since(t0)).
		Int("routes", len(routes)).
		Int("stops", len(stops)).
		Int("stoptimes", len(stoptimes)).
		Int("trips", len(trips)).
		Int("services", len(calendars)).
		Int("service-exceptions", len(calendarDates)).
		Int("shapes", len(shapes)).
		Msg("read CSV dataset")

	return &CSVDataset{
		Calendars:     calendars,
		CalendarDates: calendarDates,
		Routes:        routes,
		StopTimes:     stoptimes,
		Stops:         stops,
		Trips:         trips,
		Shapes:        shapes,
	}, nil
}
