package resolvers

import (
	"context"
	"time"

	"stop-checker.com/db/model"
	"stop-checker.com/db/repository"
)

type ScheduleResolvers struct {
}

func (r *ScheduleResolvers) Next(ctx context.Context, obj repository.Schedule, limit int, after *time.Time) ([]model.ScheduleResult, error) {
	if after == nil {
		now := time.Now()
		after = &now
	}
	return obj.After(*after, limit), nil
}

func (r *ScheduleResolvers) On(ctx context.Context, obj repository.Schedule, date time.Time) ([]model.ScheduleResult, error) {
	return obj.Day(date), nil
}

type ScheduleResultResolvers struct {
}

func (r *ScheduleResultResolvers) Datetime(ctx context.Context, obj *model.ScheduleResult) (time.Time, error) {
	return obj.Time, nil
}
