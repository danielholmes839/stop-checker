package osrm

import (
	"errors"
	"fmt"

	"stop-checker.com/db/model"
)

type Cache struct {
	data cacheData
}

func NewCache(data cacheData) *Cache {
	return &Cache{data: data}
}

func (c *Cache) GetDirections(originId, destinationId string) (model.Path, error) {
	directions, ok := c.data[fmt.Sprintf("%s:%s", originId, destinationId)]
	if ok {
		return directions, nil
	}

	directions, ok = c.data[fmt.Sprintf("%s:%s", destinationId, originId)]
	if ok {
		return directions, nil
	}

	return model.Path{}, errors.New("not found")
}
