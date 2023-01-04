package resolvers

import (
	"context"

	"stop-checker.com/db/model"
)

type LocationResolvers struct {
}

func (r *LocationResolvers) Distance(ctx context.Context, obj *model.Location, location model.Location) (float64, error) {
	return obj.Distance(location), nil
}
