package travel

import "fmt"

type PlannerMetrics interface {
	RecordExploreNode(n *node)
}

type PlannerMetricsEmpty struct {
}

func (m *PlannerMetricsEmpty) RecordExploreNode(n *node) {
}

// https://mobisoftinfotech.com/tools/plot-multiple-points-on-map/
type PlannerMetricsMiMapTools struct {
	Shape string // circle, marker, square
	Color string // #000000 hex color
}

func (m *PlannerMetricsMiMapTools) RecordExploreNode(n *node) {
	fmt.Printf("%f,%f,%s,%s\n", n.Latitude, n.Longitude, m.Color, m.Shape)
}
