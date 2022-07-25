package travel

import "container/heap"

// An IntHeap is a min-heap of ints.
type nodeHeap []*Node

func (h nodeHeap) Len() int           { return len(h) }
func (h nodeHeap) Less(i, j int) bool { return h[i].Arrival.Before(h[j].Arrival) }
func (h nodeHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *nodeHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(*Node))
}

func (h *nodeHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type NodePriorityQueue struct {
	*nodeHeap
}

func NewNodePriorityQueue() *NodePriorityQueue {
	return &NodePriorityQueue{
		nodeHeap: &nodeHeap{},
	}
}

func (pq *NodePriorityQueue) Push(node *Node) {
	heap.Push(pq.nodeHeap, node)
}

func (pq *NodePriorityQueue) Pop() *Node {
	return heap.Pop(pq.nodeHeap).(*Node)
}

func (pq *NodePriorityQueue) Empty() bool {
	return len(*pq.nodeHeap) == 0
}
