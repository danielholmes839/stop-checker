package db

import (
	"fmt"
	"time"

	"stop-checker.com/db/model"
)

type Reach struct {
	model.Stop
	Trips map[string]struct{}
}

type ReachableStop struct {
	model.Stop
	Trips map[string]struct{}
}

type IReachIndex interface {
	Reachable(stopId, routeId string) []model.Stop

	ReachableScheduledAfter(stopId, routeId string, after time.Time)
	ReachableScheduledBefore(stopId, routeId string, before time.Time)
}

type tripHashData struct {
}

type ReachIndex struct {
	stopTimesByTrip *InvertedIndex[model.StopTime]
	// {hash: {tripId: nil}}
	tripsByHash map[string]map[string]struct{}

	// {hash: {stopId: trip stop time index}}
	stopsByHash map[string]map[string]int

	// {stopId-routeId: {hash: trip stop time index}}
	hashesByStopRoute map[string]map[string]int
}

func NewReachIndex(indexes *BaseIndex, base *model.Base, stopTimesByTrip *InvertedIndex[model.StopTime]) *ReachIndex {
	t0 := time.Now()

	// {hash: {tripId: nil}}
	tripsByHash := map[string]map[string]struct{}{} // {}

	// {hash: [stopId (in order)]}
	stopsByHash := map[string]map[string]int{}

	// {stopId: {hash: stop sequence}}
	hashesByStopRoute := map[string]map[string]int{}
	// for _, stop := range base.Stops {
	// 	hashesByStopRoute[stop.Id] = map[string]int{}
	// }

	for _, trip := range base.Trips {
		// get trip hash
		stopTimes, _ := stopTimesByTrip.Get(trip.Id)
		hash, err := hashTrip(trip, stopTimes)
		if err != nil {
			continue
		}

		// add the trip to tripsByHash
		_, seen := tripsByHash[hash]
		if !seen {
			tripsByHash[hash] = map[string]struct{}{}
			stopsByHash[hash] = map[string]int{}

			// add the hash to stops
			for i, stoptime := range stopTimes {
				srId := stopRouteId(stoptime.StopId, trip.RouteId)
				if _, ok := hashesByStopRoute[srId]; !ok {
					hashesByStopRoute[srId] = map[string]int{}
				}
				hashesByStopRoute[srId][hash] = i
				stopsByHash[hash][stoptime.StopId] = i
			}
		}

		tripsByHash[hash][trip.ID()] = struct{}{}
	}

	fmt.Println("created reach index in", time.Since(t0))
	return &ReachIndex{
		stopTimesByTrip:   stopTimesByTrip,
		tripsByHash:       tripsByHash,
		stopsByHash:       stopsByHash,
		hashesByStopRoute: hashesByStopRoute,
	}
}

func (r *ReachIndex) Reachable(originId string, routeId string) map[string]map[string]struct{} {
	originHashes := r.hashesByStopRoute[stopRouteId(originId, routeId)]

	reachable := map[string]struct{}{}

	for originHash, originSequence := range originHashes {
		for destination, destinationSequence := range r.stopsByHash[originHash] {
			if originSequence >= destinationSequence || destination == originId {
				continue
			}
			reachable[destination] = struct{}{}
		}
	}

	return reachable
}

func stopRouteId(stopId, routeId string) string {
	return fmt.Sprintf("%s:%s", stopId, routeId)
}

func hashTrip(trip model.Trip, stoptimes []model.StopTime) (string, error) {
	if len(stoptimes) == 0 {
		// trip has zero stop times.
		return "", nil
	}
	return fmt.Sprintf("%s:%s:%s:%s:%d", trip.RouteId, trip.DirectionId, stoptimes[0].StopId, stoptimes[len(stoptimes)-1].StopId, len(stoptimes)), nil
}
