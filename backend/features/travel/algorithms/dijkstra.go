package algorithms

import (
	"errors"
)

type Node interface {
	ID() string
}

type Expand[N Node] func(n N) []N

type Path[N Node] struct {
	Prev *Path[N]
	Node N
}

func (p *Path[N]) ID() string {
	return p.Node.ID()
}

type DijkstraConfig[N Node] struct {
	Destination string
	Initial     []N
	Expand      Expand[N]
	Compare     func(a, b N) bool
}

func Dijkstra[N Node](config *DijkstraConfig[N]) (*Path[N], error) {
	seen := Set{}
	pq := NewPriorityQueue(func(a, b *Path[N]) bool {
		return config.Compare(a.Node, b.Node)
	})

	// push initial path
	for _, node := range config.Initial {
		pq.Push(&Path[N]{
			Prev: nil,
			Node: node,
		})
	}

	for !pq.Empty() {
		path := pq.Pop()
		node := path.Node

		// destination
		if node.ID() == config.Destination {
			return path, nil
		}

		// already expanded
		if seen.Contains(node.ID()) {
			continue
		}
		seen.Add(node.ID())

		// expand
		for _, neighbor := range config.Expand(node) {
			pq.Push(&Path[N]{
				Prev: path,
				Node: neighbor,
			})
		}
	}

	return nil, errors.New("dijkstra's algorithm: no solution")
}
