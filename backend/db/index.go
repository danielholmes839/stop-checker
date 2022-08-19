package db

import (
	"fmt"
)

type WithID interface {
	ID() string
}

type KeyFunc[R any] func(record R) (key string)

type Index[R any] struct {
	name string
	data map[string]R
}

func NewIndex[R WithID](name string, data []R) *Index[R] {
	index := &Index[R]{
		name: name,
		data: map[string]R{},
	}

	for _, record := range data {
		index.data[record.ID()] = record
	}

	return index
}

func (index *Index[R]) Get(key string) (R, error) {
	record, ok := index.data[key]
	if !ok {
		return record, fmt.Errorf("%s not found in %s index", key, index.name)
	}
	return record, nil
}

type InvertedIndex[R any] struct {
	name string
	data map[string][]R
}

func NewInvertedIndex[R any](name string, data []R, key KeyFunc[R]) *InvertedIndex[R] {
	index := &InvertedIndex[R]{
		name: name,
		data: map[string][]R{},
	}

	for _, value := range data {
		// add values
		k := key(value)
		index.data[k] = append(index.data[k], value)
	}

	return index
}

func (index *InvertedIndex[R]) Get(key string) ([]R, error) {
	records, ok := index.data[key]
	if !ok {
		return nil, fmt.Errorf("%s not found in %s inverted index", key, index.name)
	}
	return records, nil
}

func (index *InvertedIndex[R]) Data() map[string][]R {
	return index.data
}
