package db

import (
	"fmt"
	"sort"
	"time"

	"stop-checker.com/db/model"
)

type hashStopInfo struct {
	// used to quickly lookup stop times
	index int

	// used to sort reachable stops compared with other hashes
	sequence int
}

type ReachableResults map[string]map[string]hashStopInfo // { reachable stop id: { trip hash: hash stop info }}

type ReachableScheduleResult struct {
	Schedule *ScheduleResults
	model.Stop
}

type ReachIndex struct {
	trips                     *Index[model.Trip]
	stops                     *Index[model.Stop]
	stopTimesByTrip           *InvertedIndex[model.StopTime]
	indexesRequiredBySchedule *indexesRequiredBySchedule

	// {hash: {tripId: nil}}
	tripsByHash map[string]map[string]struct{}

	// {hash: {stopId: trip stop time index}}
	stopsByHash map[string]map[string]hashStopInfo

	// {stopId-routeId: {hash: trip stop time index}}
	hashesByStopRoute map[string]map[string]hashStopInfo
}

func NewReachIndex(indexes *BaseIndex, base *model.Base, stopTimesByTrip *InvertedIndex[model.StopTime], indexesRequiredBySchedule *indexesRequiredBySchedule) *ReachIndex {
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
		stops:                     indexes.Stops,
		stopTimesByTrip:           stopTimesByTrip,
		indexesRequiredBySchedule: indexesRequiredBySchedule,
		tripsByHash:               tripsByHash,
		stopsByHash:               stopsByHash,
		hashesByStopRoute:         hashesByStopRoute,
	}
}

/*
Returns the reachable stops given the origin and route.
If reverse is true then the function returns the stops that can reach the origin by the route
*/
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

func (r *ReachIndex) ReachableWithSchedule(originId, routeId string, at time.Time) []ReachableScheduleResult {
	reachable := r.reachable(originId, routeId, false)
	results := []ReachableScheduleResult{}

	for destinationId, hashInfo := range reachable {
		destination, _ := r.stops.Get(destinationId)
		scheduleResults := r.reachableScheduleResult(hashInfo)

		results = append(results, ReachableScheduleResult{
			Schedule: scheduleResults,
			Stop:     destination,
		})
	}

	return results
}

func (r *ReachIndex) ReachableStopWithSchedule(originId, destinationId, routeId string, at time.Time) ReachableScheduleResult {
	reachable := r.reachable(originId, routeId, false)
	hashInfo := reachable[destinationId]
	destination, _ := r.stops.Get(destinationId)

	scheduleResults := r.reachableScheduleResult(hashInfo)

	return ReachableScheduleResult{
		Schedule: scheduleResults,
		Stop:     destination,
	}
}

func (r *ReachIndex) reachable(originId string, routeId string, reverse bool) ReachableResults {
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

func (r *ReachIndex) reachableScheduleResult(hashInfo map[string]hashStopInfo) *ScheduleResults {
	destinationStopTimes := []model.StopTime{}

	for hash, info := range hashInfo {
		for tripId := range r.tripsByHash[hash] {
			stoptimes, _ := r.stopTimesByTrip.Get(tripId)
			destinationStopTime := stoptimes[info.index]
			destinationStopTimes = append(destinationStopTimes, destinationStopTime)
		}
	}

	// sort the stop times - TODO it should be possible to make scheduleResults.Next(after) linear and not need a sort
	sort.Slice(destinationStopTimes, func(i, j int) bool {
		return destinationStopTimes[i].Time < destinationStopTimes[j].Time
	})

	return &ScheduleResults{
		indexesRequiredBySchedule: r.indexesRequiredBySchedule,
		results:                   destinationStopTimes,
	}
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
