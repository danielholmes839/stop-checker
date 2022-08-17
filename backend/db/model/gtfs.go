package model

import (
	"fmt"
	"time"
)

type Agency struct {
	Name     string
	URL      string
	Timezone *time.Location
}

type Service struct {
	Id    string
	On    map[time.Weekday]bool
	Start time.Time
	End   time.Time
}

func (s Service) ID() string {
	return s.Id
}

type ServiceException struct {
	ServiceId string
	Date      time.Time
	Added     bool // false when service is cancelled
}

func (s ServiceException) ID() string {
	return fmt.Sprintf("exception:%s:%s", s.ServiceId, s.Date.Format("2006-01-02"))
}

type Route struct {
	Id        string
	Name      string
	Type      int
	Color     string
	TextColor string
}

func (r Route) ID() string {
	return r.Id
}

type StopTime struct {
	TripId   string
	StopId   string
	StopSeq  int
	Time     time.Time
	Overflow bool // stop times can be past 24 hours.
}

func (st StopTime) ID() string {
	return fmt.Sprintf("stoptime:%s:%s:%d", st.StopId, st.TripId, st.StopSeq)
}

type Stop struct {
	Location
	Id   string
	Code string
	Name string
	Type string
}

func (s Stop) ID() string {
	return s.Id
}

type Trip struct {
	Id          string
	RouteId     string
	ServiceId   string
	ShapeId     string
	DirectionId string
	Headsign    string
}

func (t Trip) ID() string {
	return t.Id
}
