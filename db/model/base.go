package model

import (
	"strconv"
	"time"

	"stop-checker.com/db/gtfs"
)

type Base struct {
	Agency        Agency
	Routes        []Route
	Stops         []Stop
	StopTimes     []StopTime
	Trips         []Trip
	Calendars     []Calendar
	CalendarDates []CalendarDate
}

func NewBaseFromGTFS(data *gtfs.Dataset, parser *BaseParser) *Base {
	agency := data.Agencies[0]

	base := &Base{
		Agency: Agency{
			Name:     agency.Name,
			URL:      agency.URL,
			Timezone: agency.Timezone,
		},
		Routes:        make([]Route, len(data.Routes)),
		Stops:         make([]Stop, len(data.Stops)),
		StopTimes:     make([]StopTime, len(data.StopTimes)),
		Trips:         make([]Trip, len(data.Trips)),
		Calendars:     make([]Calendar, len(data.Calendars)),
		CalendarDates: make([]CalendarDate, len(data.CalendarDates)),
	}

	for i, route := range data.Routes {
		base.Routes[i] = parser.NewRouteFromGTFS(route)
	}

	for i, stop := range data.Stops {
		base.Stops[i] = parser.NewStopFromGTFS(stop)
	}

	for i, stopTime := range data.StopTimes {
		base.StopTimes[i] = parser.NewStopTimeFromGTFS(stopTime)
	}

	for i, trips := range data.Trips {
		base.Trips[i] = parser.NewTripFromGTFS(trips)
	}

	for i, calendar := range data.Calendars {
		base.Calendars[i] = parser.NewCalendarFromGTFS(calendar)
	}

	for i, calendarDate := range data.CalendarDates {
		base.CalendarDates[i] = parser.NewCalendarDateFromGTFS(calendarDate)
	}

	return base
}

type BaseParser struct {
	TimeLayout string
	DateLayout string
}

func (b *BaseParser) NewCalendarFromGTFS(data gtfs.Calendar) Calendar {
	start, _ := time.Parse(b.DateLayout, data.Start)
	end, _ := time.Parse(b.DateLayout, data.End)

	return Calendar{
		ServiceID: data.ServiceID,
		Monday:    data.Monday == 1,
		Tuesday:   data.Tuesday == 1,
		Wednesday: data.Wednesday == 1,
		Thursday:  data.Thursday == 1,
		Friday:    data.Friday == 1,
		Saturday:  data.Saturday == 1,
		Sunday:    data.Sunday == 1,
		Start:     start,
		End:       end,
	}
}

func (b *BaseParser) NewCalendarDateFromGTFS(data gtfs.CalendarDate) CalendarDate {
	date, _ := time.Parse(b.DateLayout, data.Date)

	return CalendarDate{
		ServiceID:     data.ServiceID,
		Date:          date,
		ExceptionType: data.ExceptionType,
	}
}

func (b *BaseParser) NewRouteFromGTFS(data gtfs.Route) Route {
	return Route{
		ID:        data.ID,
		Name:      data.ShortName,
		Type:      data.Type,
		Color:     data.Color,
		TextColor: data.TextColor,
	}
}

func (b *BaseParser) NewStopTimeFromGTFS(data gtfs.StopTime) StopTime {
	seq, _ := strconv.Atoi(data.StopSeq)
	arrival, _ := time.Parse(b.TimeLayout, data.Departure)
	departure, _ := time.Parse(b.TimeLayout, data.Departure)

	return StopTime{
		StopID:    data.StopID,
		StopSeq:   seq,
		TripID:    data.TripID,
		Arrival:   arrival,
		Departure: departure,
	}
}

func (b *BaseParser) NewStopFromGTFS(data gtfs.Stop) Stop {
	return Stop{
		ID:   data.ID,
		Code: data.Code,
		Name: data.Name,
		Type: data.Type,
		Location: Location{
			Latitude:  data.Latitude,
			Longitude: data.Longitude,
		},
	}
}

func (b *BaseParser) NewTripFromGTFS(data gtfs.Trip) Trip {
	return Trip{
		ID:          data.ID,
		RouteID:     data.RouteID,
		ServiceID:   data.ServiceID,
		ShapeID:     data.ShapeID,
		DirectionID: data.DirectionID,
		Headsign:    data.Headsign,
	}
}
