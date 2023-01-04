package resolvers

import (
	"context"
	"errors"
	"time"

	"stop-checker.com/db/model"
)

type ServiceResolvers struct {
}

func (r *ServiceResolvers) Sunday(ctx context.Context, obj *model.Service) (bool, error) {
	return obj.On[time.Sunday], nil
}

func (r *ServiceResolvers) Monday(ctx context.Context, obj *model.Service) (bool, error) {
	return obj.On[time.Monday], nil
}

func (r *ServiceResolvers) Tuesday(ctx context.Context, obj *model.Service) (bool, error) {
	return obj.On[time.Tuesday], nil
}

func (r *ServiceResolvers) Wednesday(ctx context.Context, obj *model.Service) (bool, error) {
	return obj.On[time.Wednesday], nil
}

func (r *ServiceResolvers) Thursday(ctx context.Context, obj *model.Service) (bool, error) {
	return obj.On[time.Thursday], nil
}

func (r *ServiceResolvers) Friday(ctx context.Context, obj *model.Service) (bool, error) {
	return obj.On[time.Friday], nil
}

func (r *ServiceResolvers) Saturday(ctx context.Context, obj *model.Service) (bool, error) {
	return obj.On[time.Saturday], nil
}

func (r *ServiceResolvers) Exceptions(ctx context.Context, obj *model.Service) ([]model.ServiceException, error) {
	return []model.ServiceException{}, errors.New("not implemented")
}
