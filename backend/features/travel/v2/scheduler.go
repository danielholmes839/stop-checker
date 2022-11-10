package v2

import (
	"time"

	"stop-checker.com/db"
	"stop-checker.com/db/model"
)

type Scheduler struct {
	stopIndex       *db.Index[model.Stop]
	reachIndex      *db.ReachIndex
	stopTimesByTrip *db.InvertedIndex[model.StopTime]
}

func (s *Scheduler) Depart(route Route, at time.Time) Schedule {
	schedule := []*Leg{}

	first := route[0] // always walking

	schedule = append(schedule, &Leg{
		OriginId:           first.OriginId,
		Origin:             first.Origin,
		OriginArrival:      at,
		DestinationId:      first.DestinationId,
		Destination:        first.Destination,
		DestinationArrival: at.Add(walkingDuration(first.Origin.Distance(first.Destination))),
	})

	for i, leg := range route[1:] {
		previous := schedule[i]

		if leg.RouteId == "" {
			schedule = append(schedule, &Leg{
				OriginId:           leg.OriginId,
				Origin:             leg.Origin,
				OriginArrival:      previous.DestinationArrival,
				DestinationId:      leg.DestinationId,
				Destination:        leg.Destination,
				DestinationArrival: previous.DestinationArrival.Add(walkingDuration(leg.Origin.Distance(leg.Destination))),
			})
			continue
		}
	}

	return Schedule{}
}
