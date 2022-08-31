package db

import (
	"fmt"
	"sort"

	"stop-checker.com/db/model"
)

type hashStopInfo struct {
	// used to quickly lookup stop times 
	index    int

	// used to sort reachable stops compared with other hashes
	sequence int
}

type ReachIndex struct {
	stops           *Index[model.Stop]
	stopTimesByTrip *InvertedIndex[model.StopTime]

	// {hash: {tripId: nil}}
	tripsByHash map[string]map[string]struct{}

	// {hash: {stopId: trip stop time index}}
	stopsByHash map[string]map[string]hashStopInfo

	// {stopId-routeId: {hash: trip stop time index}}
	hashesByStopRoute map[string]map[string]hashStopInfo
}

func NewReachIndex(indexes *BaseIndex, base *model.Base, stopTimesByTrip *InvertedIndex[model.StopTime]) *ReachIndex {
	tripsByHash := map[string]map[string]struct{}{}
	stopsByHash := map[string]map[string]hashStopInfo{}
	hashesByStopRoute := map[string]map[string]hashStopInfo{}

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
			stopsByHash[hash] = map[string]hashStopInfo{}

			// add the hash to stops
			for i, stoptime := range stopTimes {
				srId := stopRouteId(stoptime.StopId, trip.RouteId)
				if _, ok := hashesByStopRoute[srId]; !ok {
					hashesByStopRoute[srId] = map[string]hashStopInfo{}
				}
				info := hashStopInfo{
					index:    i,
					sequence: stoptime.StopSeq,
				}
				hashesByStopRoute[srId][hash] = info
				stopsByHash[hash][stoptime.StopId] = info
			}
		}

		tripsByHash[hash][trip.ID()] = struct{}{}
	}

	return &ReachIndex{
		stops:             indexes.Stops,
		stopTimesByTrip:   stopTimesByTrip,
		tripsByHash:       tripsByHash,
		stopsByHash:       stopsByHash,
		hashesByStopRoute: hashesByStopRoute,
	}
}

func (r *ReachIndex) Reachable(originId string, routeId string, reverse bool) []model.Stop {
	reachable := r.reachable(originId, routeId, reverse)
	order := map[string]int{}
	stops := make([]model.Stop, len(reachable))

	counter := 0
	for stopId, options := range reachable {
		stop, _ := r.stops.Get(stopId)
		stops[counter] = stop
		counter++

		for _, info := range options {
			order[stopId] = info.sequence
		}
	}

	sort.Slice(stops, func(i, j int) bool {
		si := order[stops[i].Id]
		sj := order[stops[j].Id]
		return si < sj
	})

	return stops
}

func (r *ReachIndex) reachable(originId string, routeId string, reverse bool) map[string]map[string]hashStopInfo {
	originHashes := r.hashesByStopRoute[stopRouteId(originId, routeId)]

	reachable := map[string]map[string]hashStopInfo{}

	for originHash, originInfo := range originHashes {
		for destination, destinationInfo := range r.stopsByHash[originHash] {
			if destination == originId {
				continue
			}

			if !reverse && originInfo.sequence >= destinationInfo.sequence {
				continue
			}

			if reverse && originInfo.sequence <= destinationInfo.sequence {
				continue
			}

			if _, ok := reachable[destination]; !ok {
				reachable[destination] = map[string]hashStopInfo{}
			}

			reachable[destination][originHash] = destinationInfo
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
