package gtfs

import (
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"stop-checker.com/db/model"
)

type CSVParser struct {
	ParserFilter
	TZ         *time.Location
	TimeLayout string
	DateLayout string
}

func (p *CSVParser) ParseDataset(dataset *CSVDataset) *model.Dataset {
	t0 := time.Now()

	// create the services
	services := []model.Service{}
	for _, calendar := range dataset.Calendars {
		service := p.parseService(calendar)
		if !p.FilterService(service) {
			services = append(services, service)
		}
	}

	// create service exceptions
	serviceExceptions := []model.ServiceException{}
	for _, calendarDate := range dataset.CalendarDates {
		serviceException := p.parseServiceException(calendarDate)
		if !p.FilterServiceException(serviceException) {
			serviceExceptions = append(serviceExceptions, serviceException)
		}
	}

	// create trips
	trips := []model.Trip{}
	for _, tripRecord := range dataset.Trips {
		trip := p.parseTrip(tripRecord)
		if !p.FilterTrip(trip) {
			trips = append(trips, trip)
		}
	}

	// create routes
	routes := []model.Route{}
	for _, routeRecord := range dataset.Routes {
		route := p.parseRoute(routeRecord)
		if !p.FilterRoute(route) {
			routes = append(routes, route)
		}
	}

	// create stop times
	stoptimes := []model.StopTime{}
	for _, stoptimeRecord := range dataset.StopTimes {
		stoptime := p.parseStopTime(stoptimeRecord)
		if !p.FilterStopTime(stoptime) {
			stoptimes = append(stoptimes, stoptime)
		}
	}

	// create stops
	stops := []model.Stop{}
	for _, stopRecord := range dataset.Stops {
		stop := p.parseStop(stopRecord)
		if !p.FilterStop(stop) {
			stops = append(stops, stop)
		}
	}

	// create shapes
	shapes := []model.Shape{}
	for _, shapeRecord := range dataset.Shapes {
		shape := p.parseShape(shapeRecord)
		if !p.FilterShape(shape) {
			shapes = append(shapes, shape)
		}
	}

	log.Info().
		Dur("duration", time.Since(t0)).
		Int("routes", len(routes)).
		Int("stops", len(stops)).
		Int("stoptimes", len(stoptimes)).
		Int("trips", len(trips)).
		Int("services", len(services)).
		Int("service-exceptions", len(serviceExceptions)).
		Int("shapes", len(shapes)).
		Msg("parsed CSV dataset")

	return &model.Dataset{
		Routes:            routes,
		Stops:             stops,
		StopTimes:         stoptimes,
		Trips:             trips,
		Services:          services,
		ServiceExceptions: serviceExceptions,
		Shapes:            shapes,
	}
}

func (p *CSVParser) parseService(data Calendar) model.Service {
	start, _ := time.ParseInLocation(p.DateLayout, data.Start, p.TZ)
	end, _ := time.ParseInLocation(p.DateLayout, data.End, p.TZ)

	return model.Service{
		Id: data.ServiceID,
		On: [7]bool{
			data.Sunday == 1,
			data.Monday == 1,
			data.Tuesday == 1,
			data.Wednesday == 1,
			data.Thursday == 1,
			data.Friday == 1,
			data.Saturday == 1,
		},
		Start: start,
		End:   end,
	}
}

func (p *CSVParser) parseServiceException(data CalendarDate) model.ServiceException {
	date, _ := time.ParseInLocation(p.DateLayout, data.Date, p.TZ)

	return model.ServiceException{
		ServiceId: data.ServiceID,
		Date:      date,
		Added:     data.ExceptionType == 1,
	}
}

func (p *CSVParser) parseRoute(data Route) model.Route {
	return model.Route{
		Id:              data.ID,
		Name:            data.ShortName,
		BackgroundColor: "#" + data.Color,
		TextColor:       "#" + data.TextColor,
	}
}

func (p *CSVParser) parseStopTime(data StopTime) model.StopTime {
	seq, _ := strconv.Atoi(data.StopSeq)

	hours, _ := strconv.Atoi(data.Departure[0:2])
	overflow := hours >= 24
	hours %= 24

	minutes, _ := strconv.Atoi(data.Departure[3:5])

	t := model.NewTime(hours, minutes)

	return model.StopTime{
		StopId:   data.StopID,
		Sequence: seq,
		TripId:   data.TripID,
		Time:     t,
		Overflow: overflow,
	}
}

func (p *CSVParser) parseStop(data Stop) model.Stop {
	return model.Stop{
		Id:   data.ID,
		Code: data.Code,
		Name: strings.Title(strings.ToLower(data.Name)),
		Type: data.Type,
		Location: model.Location{
			Latitude:  data.Latitude,
			Longitude: data.Longitude,
		},
	}
}

func (p *CSVParser) parseTrip(data Trip) model.Trip {
	return model.Trip{
		Id:          data.ID,
		RouteId:     data.RouteID,
		ServiceId:   data.ServiceID,
		ShapeId:     data.ShapeID,
		DirectionId: data.DirectionID,
		Headsign:    data.Headsign,
	}
}

func (p *CSVParser) parseShape(data Shape) model.Shape {
	return model.Shape{
		Id: data.ID,
		Location: model.Location{
			Latitude:  data.Latitude,
			Longitude: data.Longitude,
		},
		Seq: data.Seq,
	}
}
