package model

import (
	"time"
)

type Agency struct {
	Name     string
	URL      string
	Timezone string
}

type Calendar struct {
	ServiceID string
	Monday    bool
	Tuesday   bool
	Wednesday bool
	Thursday  bool
	Friday    bool
	Saturday  bool
	Sunday    bool

	// start and end date
	Start time.Time
	End   time.Time
}

type CalendarDate struct {
	ServiceID     string
	Date          time.Time
	ExceptionType int
}

type Route struct {
	ID        string
	Name      string
	Type      int
	Color     string
	TextColor string
}

type StopTime struct {
	StopID    string
	StopSeq   int
	TripID    string
	Arrival   time.Time
	Departure time.Time
}

type Stop struct {
	ID       string
	Code     string
	Name     string
	Type     string
	Location Location
}

type Trip struct {
	ID          string
	RouteID     string
	ServiceID   string
	ShapeID     string
	DirectionID string
	Headsign    string
}

type Location struct {
	Latitude  float64
	Longitude float64
}
