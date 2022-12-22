package gtfs

import (
	"io"
	"os"
	fp "path/filepath"
)

type Input struct {
	Calendars     io.ReadCloser // service
	CalendarDates io.ReadCloser // service exceptions
	Routes        io.ReadCloser
	Stoptimes     io.ReadCloser
	Stops         io.ReadCloser
	Trips         io.ReadCloser
	Shapes        io.ReadCloser
}

func FileInput(path string) (*Input, error) {
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

	return &Input{
		Calendars:     calendars,
		CalendarDates: calendarDates,
		Routes:        routes,
		Stoptimes:     stoptimes,
		Stops:         stops,
		Trips:         trips,
		Shapes:        shapes,
	}, nil
}
