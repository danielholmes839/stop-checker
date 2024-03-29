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
		index: NewIndex("service-exceptions", data, func(exception model.ServiceException) string {
			return exception.ID()
		}),
	}
}

func (s *ServiceExceptionIndex) Get(serviceId string, t time.Time) (model.ServiceException, error) {
	id := model.ServiceException{ServiceId: serviceId, Date: t}.ID()
	return s.index.Get(id)
}
