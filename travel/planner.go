package travel

import (
	"time"

	"stop-checker.com/travel/dijkstra"
)

type Planner struct {
}

func NewPlanner() *Planner {
	return &Planner{}
}

func (p *Planner) Depart(at time.Time, origin, destination string) (Plan, error) {
	config := &dijkstra.Config[*node]{
		Destination: destination,
		Initial: &node{
			stopId:   origin,
			arrival:  at,
			duration: time.Duration(0),
		},
		Expand: p.expand,
	}

	solution, err := dijkstra.Algorithm[*node](config)
	if err != nil {
		return nil, err
	}

	return p.plan(solution), nil
}

func (p *Planner) expand(n *node) []*dijkstra.Path[*node] {

}

func (p *Planner) plan(solution *dijkstra.Path[*node]) Plan {
	return nil
}

func stopTimeDiffDuration(from, to time.Time) time.Duration {
	f := from.Hour()*60 + from.Minute()
	t := to.Hour()*60 + to.Minute()

	if t < f {
		t += 60 * 24
	}

	return time.Duration(t-f) * time.Minute
}
