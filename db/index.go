package db

import (
	"container/list"
	"fmt"
	"time"
)

type KeyFunc[T any] func(t T) (key string)

type LookupFunc[T any] func(key string) (results []T, ok bool)

type InvertedIndex[T any] struct {
	key  KeyFunc[T]
	data map[string]*list.List
}

func NewInvertedIndex[T any](data []T, key KeyFunc[T]) *InvertedIndex[T] {
	now := time.Now()
	index := &InvertedIndex[T]{
		key:  key,
		data: map[string]*list.List{},
	}

	for _, value := range data {
		index.add(value)
	}

	fmt.Println("--- inverted index created", time.Since(now))
	return index
}

func (index *InvertedIndex[T]) add(data T) {
	key := index.key(data)

	if _, ok := index.data[key]; !ok {
		linkedList := list.New()
		linkedList.PushFront(data)
		index.data[key] = linkedList
		return
	}

	index.data[key].PushFront(data)
}

func (index *InvertedIndex[T]) Get(key string) ([]T, bool) {
	_, ok := index.data[key]
	return []T{}, ok
}
