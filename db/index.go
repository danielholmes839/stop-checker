package db

type WithID interface {
	ID() string
}

type KeyFunc[R any] func(record R) (key string)

type Index[R any] struct {
	key  KeyFunc[R]
	data map[string]R
}

func NewIndex[R WithID](data []R) *Index[R] {
	index := &Index[R]{
		data: map[string]R{},
	}

	for _, record := range data {
		index.data[record.ID()] = record
	}

	return index
}

func (index *Index[R]) Get(key string) (R, bool) {
	record, ok := index.data[key]
	return record, ok
}

type InvertedIndex[R any] struct {
	data map[string][]R
}

func NewInvertedIndex[R any](data []R, key KeyFunc[R]) *InvertedIndex[R] {
	index := &InvertedIndex[R]{
		data: map[string][]R{},
	}

	for _, value := range data {
		// add values
		k := key(value)
		index.data[k] = append(index.data[k], value)
	}

	return index
}

func (index *InvertedIndex[R]) Get(key string) ([]R, bool) {
	records, ok := index.data[key]
	return records, ok
}
