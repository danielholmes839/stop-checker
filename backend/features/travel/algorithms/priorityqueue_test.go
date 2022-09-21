package algorithms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPriorityQueue(t *testing.T) {
	queue := NewPriorityQueue(func(a, b path) bool {
		return a.weight < b.weight
	})
	queue.Push(path{weight: 1})
	queue.Push(path{weight: 4})
	queue.Push(path{weight: 2})
	queue.Push(path{weight: 3})
	queue.Push(path{weight: 5})
	queue.Push(path{weight: 0})

	for i := 0; i <= 5; i++ {
		n := queue.Pop()
		assert.Equal(t, i, n.weight)
	}
}
