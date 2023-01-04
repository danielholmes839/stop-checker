package resolvers

import (
	"context"

	"stop-checker.com/db/model"
	"stop-checker.com/db/repository"
)

type TravelScheduleResolvers struct {
}

func (r *TravelScheduleResolvers) Duration(ctx context.Context, obj *model.TravelSchedule) (int, error) {
	return int(obj.Duration().Minutes()), nil
}

type TravelScheduleLegResolvers struct {
}

func (r *TravelScheduleLegResolvers) Duration(ctx context.Context, obj *model.TravelScheduleLeg) (int, error) {
	dur := obj.Destination.Arrival.Sub(obj.Origin.Arrival)
	return int(dur.Minutes()), nil
}

type TravelScheduleNodeResolvers struct {
	repository.Stops
}

func (r *TravelScheduleNodeResolvers) Stop(ctx context.Context, obj *model.TravelScheduleNode) (*model.Stop, error) {
	return nullable(r.Stops.Get(obj.Id)), nil
}
