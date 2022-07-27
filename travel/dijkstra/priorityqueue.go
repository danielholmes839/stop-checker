package dijkstra

import "container/heap"

// An IntHeap is a min-heap of ints.
type nHeap []Node

func (h nHeap) Len() int { return len(h) }
func (h nHeap) Less(i, j int) bool {
	return h[i].Weight() < h[j].Weight()
}

func (h nHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *nHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(Node))
}

func (h *nHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1] // removing the top node
	return x
}

type Weighted interface {
	Weight() int
}

type PriorityQueue[W Weighted] struct {
	heap *nHeap
}

func NewPriorityQueue[W Weighted]() *PriorityQueue[W] {
	return &PriorityQueue[W]{
		heap: &nHeap{},
	}
}

func (pq *PriorityQueue[W]) Push(w W) {
	heap.Push(pq.heap, w)
}

func (pq *PriorityQueue[W]) Pop() W {
	return heap.Pop(pq.heap).(W)
}

func (pq *PriorityQueue[W]) Empty() bool {
	return len(*pq.heap) == 0
}
