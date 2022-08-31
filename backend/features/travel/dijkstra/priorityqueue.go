package dijkstra

import "container/heap"

// An IntHeap is a min-heap of ints.
type nHeap[N any] struct {
	Nodes   []N
	Compare func(a, b N) bool
}

func (h nHeap[N]) Len() int { return len(h.Nodes) }
func (h nHeap[N]) Less(i, j int) bool {
	return h.Compare(h.Nodes[i], h.Nodes[j])
}

func (h nHeap[N]) Swap(i, j int) { h.Nodes[i], h.Nodes[j] = h.Nodes[j], h.Nodes[i] }

func (h *nHeap[N]) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	h.Nodes = append(h.Nodes, x.(N))
}

func (h *nHeap[N]) Pop() any {
	old := h.Nodes
	n := len(old)
	x := old[n-1]
	h.Nodes = old[0 : n-1] // removing the top node
	return x
}

type PriorityQueue[N Node] struct {
	heap *nHeap[N]
}

func NewPriorityQueue[N Node](compare func(a, b N) bool) *PriorityQueue[N] {
	return &PriorityQueue[N]{
		heap: &nHeap[N]{
			Nodes:   []N{},
			Compare: compare,
		},
	}
}

func (pq *PriorityQueue[N]) Push(n N) {
	heap.Push(pq.heap, n)
}

func (pq *PriorityQueue[N]) Pop() N {
	return heap.Pop(pq.heap).(N)
}

func (pq *PriorityQueue[N]) Empty() bool {
	return len(pq.heap.Nodes) == 0
}
