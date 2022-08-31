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

type reachableResults map[string]map[string]hashStopInfo // { reachable stop id: { trip hash: hash stop info }}

type ReachableSchedule struct {
	Departure   time.Time  // departure time from the origin
	Arrival     time.Time  // arrival time at the destination
	Origin      model.Stop // origin stop
	Destination model.Stop // destination stop
	Trip        model.Trip
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
		trips:                     indexes.Trips,
		stops:                     indexes.Stops,
		stopTimesByTrip:           stopTimesByTrip,
		indexesRequiredBySchedule: indexesRequiredBySchedule,
		tripsByHash:               tripsByHash,
		stopsByHash:               stopsByHash,
		hashesByStopRoute:         hashesByStopRoute,
	}
}

func (r *ReachIndex) Reachable(originId string, routeId string, reverse bool) []model.Stop {
	/*
		Returns the reachable stops given the origin and route sorted by stop sequence (always ascending)
		If reverse is true then the function returns the incident stops using this route
	*/
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

/* ReachableBetweenWithSchedule
returns the *ScheduleResults for the all stop times that run between the origin and destination
(origin, destination)
- used by the travel.Scheduler
- used to provide alternative stop times for travel plans
*/
func (r *ReachIndex) ReachableBetweenWithSchedule(originId, destinationId, routeId string) (*ScheduleResults, *ScheduleResults) {
	// hashes that visit the origin and destination
	originHashes := r.hashesByStopRoute[stopRouteId(originId, routeId)]
	destinationHashes := r.hashesByStopRoute[stopRouteId(destinationId, routeId)]

	// origin and destination stop times that are part of a trip that reaches both stops
	originStopTimes := []model.StopTime{}
	destinationStopTimes := []model.StopTime{}

	for hash, originInfo := range originHashes {
		destinationInfo, shared := destinationHashes[hash]
		if !shared {
			continue
		}

		// add stop times from each trip
		for tripId := range r.tripsByHash[hash] {
			stopTimes, _ := r.stopTimesByTrip.Get(tripId)
			originStopTimes = append(originStopTimes, stopTimes[originInfo.index])
			destinationStopTimes = append(destinationStopTimes, stopTimes[destinationInfo.index])
		}
	}

	// sort stop times and create schedule results
	model.StopTimeSort(originStopTimes)
	model.StopTimeSort(destinationStopTimes)

	originResults := &ScheduleResults{
		indexesRequiredBySchedule: r.indexesRequiredBySchedule,
		results:                   originStopTimes,
	}

	destinationResults := &ScheduleResults{
		indexesRequiredBySchedule: r.indexesRequiredBySchedule,
		results:                   destinationStopTimes,
	}

	return originResults, destinationResults
}

func (r *ReachIndex) ReachableForwardWithNext(originId, routeId string, after time.Time) []ReachableSchedule {
	/*
		1. get all stop times (as a *ScheduleResults object) for each hash
		2. get next stop time for each *ScheduleResults for each hash
		3. calculate the closest for each
	*/

	origin, _ := r.stops.Get(originId)
	originScheduleResultsByHash := r.stopTimesByHash(originId, routeId)
	originNextByHash := map[string]ScheduleResult{}

	for hash, originScheduleResults := range originScheduleResultsByHash {
		next, err := originScheduleResults.Next(after)
		if err != nil {
			continue
		}
		originNextByHash[hash] = next
	}

	reachableForward := r.reachable(originId, routeId, false)
	results := []ReachableSchedule{}

	for destinationId, destinationHashInfo := range reachableForward {
		destination, _ := r.stops.Get(destinationId)

		set := false
		result := ReachableSchedule{
			Origin:      origin,
			Destination: destination,
		}

		// find the next stop time for this destination
		for hash, destinationInfo := range destinationHashInfo {

			// if there's no next stop time for this hash
			if _, ok := originNextByHash[hash]; !ok {
				continue
			}

			originNext := originNextByHash[hash]
			departure := originNext.Time
			departureStopTime := originNext.StopTime

			if !set || departure.Before(result.Departure) {
				// find the arrival time
				stopTimes, _ := r.stopTimesByTrip.Get(departureStopTime.TripId)
				arrivalStopTime := stopTimes[destinationInfo.index]
				arrival := departure.Add(model.TimeDiff(departureStopTime.Time, arrivalStopTime.Time))

				// update result fields
				result.Departure = departure
				result.Arrival = arrival
				result.Trip, _ = r.trips.Get(departureStopTime.TripId)

				set = true
			}
		}

		if !set {
			continue
		}

		results = append(results, result)
	}

	return results
}

func (r *ReachIndex) ReachableBackwardWithPrevious(destinationId, routeId string, before time.Time) []ReachableSchedule {
	destination, _ := r.stops.Get(destinationId)
	destinationScheduleResultsByHash := r.stopTimesByHash(destinationId, routeId)
	destinationPreviousByHash := map[string]ScheduleResult{}

	for hash, destinationScheduleResults := range destinationScheduleResultsByHash {
		next, err := destinationScheduleResults.Previous(before)
		if err != nil {
			continue
		}
		destinationPreviousByHash[hash] = next
	}

	reachableBackward := r.reachable(destinationId, routeId, true)
	results := []ReachableSchedule{}

	for originId, originHashInfo := range reachableBackward {
		origin, _ := r.stops.Get(originId)

		set := false
		result := ReachableSchedule{
			Origin:      origin,
			Destination: destination,
		}

		// find the next stop time for this destination
		for hash, originInfo := range originHashInfo {

			// if there's no next stop time for this hash
			if _, ok := destinationPreviousByHash[hash]; !ok {
				continue
			}

			destinationPrevious := destinationPreviousByHash[hash]
			arrival := destinationPrevious.Time
			arrivalStopTime := destinationPrevious.StopTime

			if !set || arrival.After(result.Arrival) {
				// find the arrival time
				stopTimes, _ := r.stopTimesByTrip.Get(arrivalStopTime.TripId)
				departureStopTime := stopTimes[originInfo.index]
				departure := arrival.Add(-model.TimeDiff(departureStopTime.Time, arrivalStopTime.Time))

				// update result fields
				result.Departure = departure
				result.Arrival = arrival
				result.Trip, _ = r.trips.Get(departureStopTime.TripId)

				set = true
			}
		}

		if !set {
			continue
		}

		results = append(results, result)
	}

	return results
}

/* reachable
returns what trip hashes can be used to reach each stop
*/
func (r *ReachIndex) reachable(originId string, routeId string, reverse bool) reachableResults {
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

// {hash: schedule results}
func (r *ReachIndex) stopTimesByHash(stopId, routeId string) map[string]*ScheduleResults {
	stopTimesByHash := map[string]*ScheduleResults{}
	stopRoute := stopRouteId(stopId, routeId)

	for hash, info := range r.hashesByStopRoute[stopRoute] {
		stopTimes := []model.StopTime{}

		for tripId := range r.tripsByHash[hash] {
			tripStopTimes, _ := r.stopTimesByTrip.Get(tripId)
			stoptime := tripStopTimes[info.index]
			stopTimes = append(stopTimes, stoptime)
		}

		sort.Slice(stopTimes, func(i, j int) bool {
			return stopTimes[i].Time < stopTimes[j].Time
		})

		stopTimesByHash[hash] = &ScheduleResults{
			indexesRequiredBySchedule: r.indexesRequiredBySchedule,
			results:                   stopTimes,
		}
	}

	return stopTimesByHash
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
