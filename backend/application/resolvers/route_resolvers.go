package resolvers

import (
	"context"

	"stop-checker.com/db/model"
)

type RouteResolvers struct {
}

func (r *RouteResolvers) ID(ctx context.Context, obj *model.Route) (string, error) {
	return obj.Id, nil
}

func (r *RouteResolvers) Text(ctx context.Context, obj *model.Route) (string, error) {
	return obj.TextColor, nil
}

func (r *RouteResolvers) Background(ctx context.Context, obj *model.Route) (string, error) {
	return obj.BackgroundColor, nil
}
