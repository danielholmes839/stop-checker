package dijkstra

import (
	"errors"
)

type Node interface {
	ID() string
	Weighted
}

type Expand[N Node] func(n *Path[N]) []*Path[N]

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

type Input[N Node] struct {
	Destination string
	Initial     N
	Expand      Expand[N]
}

func Algorithm[N Node](input *Input[N]) (*Path[N], error) {
	seen := Set{}
	pq := NewPriorityQueue[*Path[N]]()

	// push initial path
	pq.Push(&Path[N]{
		Prev: nil,
		Node: input.Initial,
	})

	for !pq.Empty() {
		path := pq.Pop()
		node := path.Node

		// destination
		if node.ID() == input.Destination {
			return path, nil
		}

		// already expanded
		if seen.Contains(node.ID()) {
			continue
		}
		seen.Add(node.ID())

		// expand
		for _, neighbor := range input.Expand(path) {
			pq.Push(neighbor)
		}
	}

	return nil, errors.New("dijkstra's algorithm: no solution")
}
