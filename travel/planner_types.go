package travel

import "time"

type node struct {
	stopId   string
	arrival  time.Time
	duration time.Duration
}

func (n *node) ID() string {
	return n.stopId
}

func (n *node) Weight() int {
	return int(n.duration)
}
