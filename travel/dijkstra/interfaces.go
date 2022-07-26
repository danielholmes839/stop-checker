package dijkstra

import "errors"

type Node interface {
	Previous() Node
	ID() string
	Weight() int
}

type Expand[N Node] func(n N) []N

func Dijkstra[N Node](initial N, destination string, expand Expand[N]) (N, error) {
	pq := NewPriorityQueue[N]()
	pq.Push(initial)

	seen := Set{}

	for !pq.Empty() {
		node := pq.Pop()

		// destination reached!
		if node.ID() == destination {
			return node, nil
		}

		// do not revisit nodes
		if seen.Contains(node.ID()) {
			continue
		}
		seen.Add(node.ID())

		// expand the node
		for _, node := range expand(node) {
			pq.Push(node)
		}
	}

	return initial, errors.New("dijkstra's algorithm: no solution")
}
