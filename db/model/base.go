package model

import (
	"strconv"
	"time"

	"stop-checker.com/db/gtfs"
)

type Base struct {
	Agency            Agency
	Routes            []Route
	Stops             []Stop
	StopTimes         []StopTime
	Trips             []Trip
	Services          []Service
	ServiceExceptions []ServiceException
}

func NewBaseFromGTFS(data *gtfs.Dataset, parser *BaseParser) *Base {
	agency := data.Agencies[0]

	base := &Base{
		Agency: Agency{
			Name:     agency.Name,
			URL:      agency.URL,
			Timezone: agency.Timezone,
		},
		Routes:            make([]Route, len(data.Routes)),
		Stops:             make([]Stop, len(data.Stops)),
		StopTimes:         make([]StopTime, len(data.StopTimes)),
		Trips:             make([]Trip, len(data.Trips)),
		Services:          make([]Service, len(data.Calendars)),
		ServiceExceptions: make([]ServiceException, len(data.CalendarDates)),
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
		base.Services[i] = parser.NewCalendarFromGTFS(calendar)
	}

	for i, calendarDate := range data.CalendarDates {
		base.ServiceExceptions[i] = parser.NewCalendarDateFromGTFS(calendarDate)
	}

	return base
}

type BaseParser struct {
	TimeLayout string
	DateLayout string
}

func (b *BaseParser) NewCalendarFromGTFS(data gtfs.Calendar) Service {
	start, _ := time.Parse(b.DateLayout, data.Start)
	end, _ := time.Parse(b.DateLayout, data.End)

	return Service{
		Id:        data.ServiceID,
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

func (b *BaseParser) NewCalendarDateFromGTFS(data gtfs.CalendarDate) ServiceException {
	date, _ := time.Parse(b.DateLayout, data.Date)

	return ServiceException{
		ServiceId: data.ServiceID,
		Date:      date,
		Added:     data.ExceptionType == 1,
	}
}

func (b *BaseParser) NewRouteFromGTFS(data gtfs.Route) Route {
	return Route{
		Id:        data.ID,
		Name:      data.ShortName,
		Type:      data.Type,
		Color:     data.Color,
		TextColor: data.TextColor,
	}
}

func (b *BaseParser) NewStopTimeFromGTFS(data gtfs.StopTime) StopTime {
	seq, _ := strconv.Atoi(data.StopSeq)
	arrival, _ := time.Parse(b.TimeLayout, data.Departure)

	return StopTime{
		StopId:  data.StopID,
		StopSeq: seq,
		TripId:  data.TripID,
		Time:    arrival,
	}
}

func (b *BaseParser) NewStopFromGTFS(data gtfs.Stop) Stop {
	return Stop{
		Id:   data.ID,
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
		Id:          data.ID,
		RouteId:     data.RouteID,
		ServiceId:   data.ServiceID,
		ShapeId:     data.ShapeID,
		DirectionId: data.DirectionID,
		Headsign:    data.Headsign,
	}
}
