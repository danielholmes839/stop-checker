package model

import "fmt"

type StopRoute struct {
	RouteId     string
	StopId      string
	DirectionId string
	Headsign    string
}

func (r StopRoute) DirectedID() string {
	return fmt.Sprintf("%s:%s", r.RouteId, r.DirectionId)
}
