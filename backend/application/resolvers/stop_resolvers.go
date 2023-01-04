package resolvers

import (
	"context"

	"stop-checker.com/db/model"
	"stop-checker.com/db/repository"
)

type StopResolvers struct {
	repository.StopRoutes
}

func (r *StopResolvers) ID(ctx context.Context, obj *model.Stop) (string, error) {
	return obj.ID(), nil
}

func (r *StopResolvers) Routes(ctx context.Context, obj *model.Stop) ([]model.StopRoute, error) {
	return r.StopRoutes.Get(obj.ID()), nil
}
