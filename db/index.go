package db

type KeyFunc[T any] func(t T) (key string)

type Index[T any] struct {
	key  KeyFunc[T]
	data map[string]T
}

func NewIndex[T any](data []T, key KeyFunc[T]) *Index[T] {
	index := &Index[T]{
		data: map[string]T{},
	}

	for _, value := range data {
		index.data[key(value)] = value
	}

	return index
}

func (index *Index[T]) Get(key string) (T, bool) {
	result, ok := index.data[key]
	return result, ok
}

type InvertedIndex[T any] struct {
	data map[string][]T
}

func NewInvertedIndex[T any](data []T, key KeyFunc[T]) *InvertedIndex[T] {
	index := &InvertedIndex[T]{
		data: map[string][]T{},
	}

	for _, value := range data {
		// add values
		k := key(value)
		index.data[k] = append(index.data[k], value)
	}

	return index
}

func (index *InvertedIndex[T]) Get(key string) ([]T, bool) {
	results, ok := index.data[key]
	return results, ok
}
