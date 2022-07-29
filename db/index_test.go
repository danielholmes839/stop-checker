package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type record struct {
	id    string
	value string
}

func (r record) ID() string {
	return r.id
}

func TestIndex(t *testing.T) {
	index := NewIndex("test_index", []record{
		{id: "a", value: "1"},
	})

	r, err := index.Get("a")
	assert.NoError(t, err)
	assert.Equal(t, 1, r.value)
}

func TestInvertedIndex(t *testing.T) {
	index := NewInvertedIndex([]record{
		{id: "a", value: "1"},
		{id: "b", value: "1"},
		{id: "c", value: "1"},
	}, func(r record) (key string) {
		return r.value
	})

	results, err := index.Get("1")
	assert.Equal(t, 3, len(results))
	assert.Equal(t, "a", results[0].id)
	assert.Equal(t, "b", results[1].id)
	assert.Equal(t, "c", results[2].id)
	
}