package travel

import (
	"fmt"
	"time"
)

type transit struct {
	TripId                string
}

type Leg struct {
	Transit     *transit      // transit info
	Origin      string        // origin stop id
	Destination string        // destination stop id
	Walk        bool          // if we walk between the two stops
	Departure   time.Time     // when do we depart from the origin
	Duration    time.Duration // duration between arriving at the destination and leaving the origin
}

func (leg *Leg) String() string {
	if leg.Walk {
		return fmt.Sprintf("{origin:%s, destination:%s, walk:%t, departure:%s, duration:%s}",
			leg.Origin, leg.Destination, leg.Walk, leg.Departure, leg.Duration)
	}
	return fmt.Sprintf("{origin:%s, destination:%s, walk:%t, departure:%s, duration:%s, trip:%s}",
		leg.Origin, leg.Destination, leg.Walk, leg.Departure, leg.Duration, leg.Transit.TripId)
}

/* FixedLeg
- leg of a travel plan without assigned times
*/
type FixedLeg struct {
	Origin      string // stop id
	Destination string // stop id
	RouteId     string // stop id
	Walk        bool   // if true then: route id is empty
}

func (fl *FixedLeg) String() string {
	if fl.Walk {
		return fmt.Sprintf("walk{origin:%s, destination:%s}", fl.Origin, fl.Destination)
	}
	return fmt.Sprintf("transit{origin:%s, destination:%s, route:%s}", fl.Origin, fl.Destination, fl.RouteId)
}
