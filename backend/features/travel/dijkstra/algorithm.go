package dijkstra

import (
	"errors"
)

type Node interface {
	ID() string
	Weighted
}

type Expand[N Node] func(n N) []N

type Path[N Node] struct {
	Prev *Path[N]
	Node N
}

func (p *Path[N]) ID() string {
	return p.Node.ID()
}

func (p *Path[N]) Weight() int {
	return p.Node.Weight()
}

type Config[N Node] struct {
	Destination string
	Initial     N
	Expand      Expand[N]
}

func Algorithm[N Node](config *Config[N]) (*Path[N], error) {
	seen := Set{}
	pq := NewPriorityQueue[*Path[N]]()

	// push initial path
	pq.Push(&Path[N]{
		Prev: nil,
		Node: config.Initial,
	})

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
