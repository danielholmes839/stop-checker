package v3

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"stop-checker.com/db/model"
	"stop-checker.com/db/repository"
)

type edgeFactory struct {
	directions      walkingDirections
	directionsCache walkingDirectionsCache
	stops           repository.Stops
	reach           scheduleReach
}

func (f *edgeFactory) Edges(plan *model.TravelPlan) ([]scheduleEdge, error) {
	edges := []scheduleEdge{}

	current := &scheduleNode{
		Id:       "#ORIGIN",
		Location: plan.Origin,
	}

	for _, leg := range plan.Legs {
		originNode, destinationNode, err := f.getNodes(leg)
		if err != nil {
			return nil, err
		}

		if current.Id != originNode.Id {
			// add a walking edge from current to origin node
			edges = append(edges, f.getWalkingDirectionsEdge(current, originNode))
		}

		edges = append(edges, &scheduleTransitEdge{
			edge:    &edge{origin: originNode, destination: destinationNode},
			routeId: leg.RouteId,
			reach:   f.reach,
		})

		current = destinationNode
	}

	edges = append(edges, f.getWalkingDirectionsEdge(current, &scheduleNode{
		Id:       "#DESTINATION",
		Location: plan.Destination,
	}))

	return edges, nil
}

func (f *edgeFactory) getWalkingDirectionsEdge(origin, destination *scheduleNode) *scheduleWalkEdge {
	path := f.getDirections(origin, destination)

	return &scheduleWalkEdge{
		edge:     &edge{origin: origin, destination: destination},
		path:     &path,
		duration: walkingDuration(path.Distance),
	}
}

func (f *edgeFactory) getDirections(origin, destination *scheduleNode) model.Path {
	// first check the directions cache. expected to error when the final origin/destination are used
	if directions, err := f.directionsCache.GetDirections(origin.Id, destination.Id); err == nil {
		return directions
	}

	// get directions. not expected to error
	if directions, err := f.directions.GetDirections(origin.Location, destination.Location); err == nil {
		return directions
	}

	log.Warn().
		Str("origin", fmt.Sprintf("%+v", origin)).
		Str("destination", fmt.Sprintf("%+v", destination)).
		Msg("scheduler edge factory failed to get walking directions")

	return model.Path{
		Distance: model.Distance(origin.Location, destination.Location),
		Path:     []model.Location{origin.Location, destination.Location},
	}
}

// return the origin and destination nodes for a model.TravelPlanLeg
func (f *edgeFactory) getNodes(leg model.TravelPlanLeg) (*scheduleNode, *scheduleNode, error) {
	// get origin stop
	origin, err := f.stops.Get(leg.OriginId)
	if err != nil {
		return nil, nil, err
	}

	originNode := &scheduleNode{
		Id:       leg.OriginId,
		Location: origin.Location,
	}

	// get destination stop
	destination, err := f.stops.Get(leg.DestinationId)
	if err != nil {
		return nil, nil, err
	}

	destinationNode := &scheduleNode{
		Id:       leg.DestinationId,
		Location: destination.Location,
	}

	return originNode, destinationNode, nil
}
