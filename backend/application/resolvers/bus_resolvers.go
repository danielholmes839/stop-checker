package resolvers

import (
	"context"
	"fmt"
	"time"

	"stop-checker.com/db/model"
)

type BusResolvers struct {
}

func (r *BusResolvers) LastUpdatedMinutes(ctx context.Context, obj *model.Bus) (int, error) {
	now := model.NewTimeFromDateTime(time.Now().Local())
	diff := int(model.TimeDiff(obj.LastUpdated, now).Minutes())
	return diff, nil
}

func (r *BusResolvers) LastUpdatedMessage(ctx context.Context, obj *model.Bus) (string, error) {
	now := model.NewTimeFromDateTime(time.Now().Local())
	diff := int(model.TimeDiff(obj.LastUpdated, now).Minutes())
	if diff <= 1 {
		return "Updated just now", nil
	}
	return fmt.Sprintf("Updated %d minutes ago", diff), nil
}

func (r *BusResolvers) Distance(ctx context.Context, obj *model.Bus) (*float64, error) {
	if obj.Location == nil {
		return nil, nil
	}
	distance := obj.Location.Distance(obj.Destination.Location)
	return &distance, nil
}
