package model

type Base struct {
	Routes            []Route
	Stops             []Stop
	StopTimes         []StopTime
	Trips             []Trip
	Services          []Service
	ServiceExceptions []ServiceException
	Shapes            []Shape
}
