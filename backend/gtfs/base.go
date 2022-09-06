package gtfs

import (
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"stop-checker.com/db/model"
)

type BaseOptions struct {
	Time       time.Time
	TimeZone   *time.Location
	TimeLayout string
	DateLayout string
}

func NewBase(r *raw, opts BaseOptions) (*model.Base, error) {
	dataset := newDataset(r)

	t0 := time.Now()

	validRoutes := map[string]struct{}{}
	validTrips := map[string]struct{}{}
	validServices := map[string]struct{}{}

	// create the services
	services := []model.Service{}
	for _, calendar := range dataset.Calendars {
		service := NewService(calendar, opts)

		if service.End.After(opts.Time) {
			validServices[service.Id] = struct{}{}
			services = append(services, service)
		}
	}

	// create service exceptions
	serviceExceptions := []model.ServiceException{}
	for _, calendarDate := range dataset.CalendarDates {
		serviceException := NewServiceException(calendarDate, opts)

		if _, ok := validServices[serviceException.ServiceId]; ok {
			serviceExceptions = append(serviceExceptions, serviceException)
		}
	}

	// create trips
	trips := []model.Trip{}
	for _, tripRecord := range dataset.Trips {
		trip := NewTrip(tripRecord)
		if _, ok := validServices[trip.ServiceId]; ok {
			trips = append(trips, trip)
			validTrips[trip.Id] = struct{}{}
			validRoutes[trip.RouteId] = struct{}{}
		}
	}

	// create routes
	routes := []model.Route{}
	for _, routeRecord := range dataset.Routes {
		route := NewRoute(routeRecord)
		if _, ok := validRoutes[route.Id]; ok {
			routes = append(routes, route)
		}
	}

	// create stop times
	stoptimes := []model.StopTime{}
	for _, stoptimeRecord := range dataset.StopTimes {
		stoptime := NewStopTime(stoptimeRecord)
		if _, ok := validTrips[stoptime.TripId]; ok {
			stoptimes = append(stoptimes, stoptime)
		}
	}

	// create stops
	stops := []model.Stop{}
	for _, stopRecord := range dataset.Stops {
		stop := NewStop(stopRecord)
		stops = append(stops, stop)
	}

	log.Info().
		Dur("duration", time.Since(t0)).
		Int("routes", len(routes)).
		Int("stops", len(stops)).
		Int("stoptimes", len(stoptimes)).
		Int("trips", len(trips)).
		Int("services", len(services)).
		Int("service-exceptions", len(serviceExceptions)).
		Msg("filtered dataset")

	return &model.Base{
		Routes:            routes,
		Stops:             stops,
		StopTimes:         stoptimes,
		Trips:             trips,
		Services:          services,
		ServiceExceptions: serviceExceptions,
	}, nil
}

type DatasetParser struct {
	TimeZone   *time.Location
	TimeLayout string
	DateLayout string

	ValidServices map[string]struct{} // all services with an end date before the startup time
	ValidRoutes   map[string]struct{} //
}

func NewService(data Calendar, opts BaseOptions) model.Service {
	start, _ := time.ParseInLocation(opts.DateLayout, data.Start, opts.TimeZone)
	end, _ := time.ParseInLocation(opts.DateLayout, data.End, opts.TimeZone)

	return model.Service{
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

func NewServiceException(data CalendarDate, opts BaseOptions) model.ServiceException {
	date, _ := time.ParseInLocation(opts.DateLayout, data.Date, opts.TimeZone)

	return model.ServiceException{
		ServiceId: data.ServiceID,
		Date:      date,
		Added:     data.ExceptionType == 1,
	}
}

func NewRoute(data Route) model.Route {
	return model.Route{
		Id:        data.ID,
		Name:      data.ShortName,
		Type:      data.Type,
		Color:     "#" + data.Color,
		TextColor: "#" + data.TextColor,
	}
}

func NewStopTime(data StopTime) model.StopTime {
	seq, _ := strconv.Atoi(data.StopSeq)

	hours, _ := strconv.Atoi(data.Departure[0:2])
	overflow := hours >= 24

	hours %= 24

	minutes, _ := strconv.Atoi(data.Departure[3:5])

	t := model.NewTime(hours, minutes)

	return model.StopTime{
		StopId:   data.StopID,
		StopSeq:  seq,
		TripId:   data.TripID,
		Time:     t,
		Overflow: overflow,
	}
}

func NewStop(data Stop) model.Stop {
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

func NewTrip(data Trip) model.Trip {
	return model.Trip{
		Id:          data.ID,
		RouteId:     data.RouteID,
		ServiceId:   data.ServiceID,
		ShapeId:     data.ShapeID,
		DirectionId: data.DirectionID,
		Headsign:    data.Headsign,
	}
}
