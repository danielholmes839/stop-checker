// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package types

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"stop-checker.com/travel"
)

type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type PageInput struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}

type TravelLegInput struct {
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Route       *string `json:"route"`
}

type TravelRoutePayload struct {
	Route  travel.Route `json:"route"`
	Errors []*Error     `json:"errors"`
}

type TravelRoutePlannerInput struct {
	Origin      string     `json:"origin"`
	Destination string     `json:"destination"`
	Departure   *time.Time `json:"departure"`
}

type TravelSchedulePayload struct {
	Schedule travel.Schedule `json:"schedule"`
	Errors   []*Error        `json:"errors"`
}

type TravelSchedulePlannerInput struct {
	Legs      []*TravelLegInput `json:"legs"`
	Departure *time.Time        `json:"departure"`
	Arrival   *time.Time        `json:"arrival"`
}

type RouteType string

const (
	RouteTypeBus   RouteType = "BUS"
	RouteTypeTrain RouteType = "TRAIN"
)

var AllRouteType = []RouteType{
	RouteTypeBus,
	RouteTypeTrain,
}

func (e RouteType) IsValid() bool {
	switch e {
	case RouteTypeBus, RouteTypeTrain:
		return true
	}
	return false
}

func (e RouteType) String() string {
	return string(e)
}

func (e *RouteType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RouteType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RouteType", str)
	}
	return nil
}

func (e RouteType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
