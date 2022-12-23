package osrm

type osrmResponse struct {
	Code   string              `json:"code"`
	Routes []osrmResponseRoute `json:"routes"`
}

type osrmResponseRoute struct {
	Legs []osrmResponseLeg `json:"legs"`
}

type osrmResponseLeg struct {
	Steps []osrmResponseStep `json:"steps"`
}

type osrmResponseStep struct {
	Intersections []osrmResponseIntersection `json:"intersections"`
}

type osrmResponseIntersection struct {
	Location [2]float64 `json:"location"`
}
