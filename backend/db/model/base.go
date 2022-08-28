package model

import (
	"strconv"
	"time"

	"stop-checker.com/db/model/gtfs"
)

type Base struct {
	Routes            []Route
	Stops             []Stop
	StopTimes         []StopTime
	Trips             []Trip
	Services          []Service
	ServiceExceptions []ServiceException
}

func NewDatasetFromGTFS(data *gtfs.Dataset, parser *DatasetParser) *Base {
	base := &Base{
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

type DatasetParser struct {
	TimeZone   *time.Location
	TimeLayout string
	DateLayout string
}

func (p *DatasetParser) NewCalendarFromGTFS(data gtfs.Calendar) Service {
	start, _ := time.ParseInLocation(p.DateLayout, data.Start, p.TimeZone)
	end, _ := time.ParseInLocation(p.DateLayout, data.End, p.TimeZone)

	return Service{
		Id: data.ServiceID,
		On: map[time.Weekday]bool{
			time.Sunday:    data.Sunday == 1,
			time.Monday:    data.Monday == 1,
			time.Tuesday:   data.Tuesday == 1,
			time.Wednesday: data.Wednesday == 1,
			time.Thursday:  data.Thursday == 1,
			time.Friday:    data.Friday == 1,
			time.Saturday:  data.Saturday == 1,
		},
		Start: start,
		End:   end.Add(time.Hour*24 - time.Minute),
	}
}

func (p *DatasetParser) NewCalendarDateFromGTFS(data gtfs.CalendarDate) ServiceException {
	date, _ := time.ParseInLocation(p.DateLayout, data.Date, p.TimeZone)

	return ServiceException{
		ServiceId: data.ServiceID,
		Date:      date,
		Added:     data.ExceptionType == 1,
	}
}

func (p *DatasetParser) NewRouteFromGTFS(data gtfs.Route) Route {
	return Route{
		Id:        data.ID,
		Name:      data.ShortName,
		Type:      data.Type,
		Color:     "#" + data.Color,
		TextColor: "#" + data.TextColor,
	}
}

func (p *DatasetParser) NewStopTimeFromGTFS(data gtfs.StopTime) StopTime {
	seq, _ := strconv.Atoi(data.StopSeq)

	hours, _ := strconv.Atoi(data.Departure[0:2])
	overflow := hours >= 24

	hours %= 24

	minutes, _ := strconv.Atoi(data.Departure[3:5])

	t := NewTime(hours, minutes)

	return StopTime{
		StopId:   data.StopID,
		StopSeq:  seq,
		TripId:   data.TripID,
		Time:     t,
		Overflow: overflow,
	}
}

func (p *DatasetParser) NewStopFromGTFS(data gtfs.Stop) Stop {
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

func (p *DatasetParser) NewTripFromGTFS(data gtfs.Trip) Trip {
	return Trip{
		Id:          data.ID,
		RouteId:     data.RouteID,
		ServiceId:   data.ServiceID,
		ShapeId:     data.ShapeID,
		DirectionId: data.DirectionID,
		Headsign:    data.Headsign,
	}
}
