package model

import "time"

/* database and index related models  */

type StopWithDistance struct {
	Stop
	Distance float64
}

type ReachableSchedule struct {
	Departure   time.Time // departure time from the origin
	Arrival     time.Time // arrival time at the destination
	Origin      Stop      // origin stop
	Destination Stop      // destination stop
	Trip        Trip
}
