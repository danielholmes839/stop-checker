package db

import (
	"time"

	"stop-checker.com/db/model"
)

type ServiceExceptionIndex struct {
	index *Index[model.ServiceException]
}

func NewServiceExceptionIndex(data []model.ServiceException) *ServiceExceptionIndex {
	return &ServiceExceptionIndex{
		index: NewIndex("service exception", data),
	}
}

func (s *ServiceExceptionIndex) Get(serviceId string, t time.Time) (model.ServiceException, error) {
	id := model.ServiceException{ServiceId: serviceId, Date: t}.ID()
	return s.index.Get(id)
}
