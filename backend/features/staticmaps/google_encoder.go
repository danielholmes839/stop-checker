package staticmaps

import (
	"fmt"
	"net/url"
)

type GoogleMapEncoder struct {
	Key string
}

func (e *GoogleMapEncoder) Encode(m *Map) string {
	center := calculateCenter(m.markers)
	zoom := calculateZoom(center, m.markers, m.width, m.height)

	centerParamValue := fmt.Sprintf("%f,%f", center.Latitude, center.Longitude)
	sizeParamValue := fmt.Sprintf("%dx%d", m.width, m.height)

	params := url.Values{}
	params.Add("key", e.Key)
	params.Add("center", centerParamValue)
	params.Add("size", sizeParamValue)
	params.Add("zoom", fmt.Sprint(zoom))

	for _, marker := range m.markers {
		params.Add("markers", fmt.Sprintf("color:%s|size:%s|label:%s|%f,%f",
			marker.Color,
			marker.Size,
			marker.Label,
			marker.Latitude,
			marker.Longitude,
		))
	}

	return fmt.Sprintf("https://maps.googleapis.com/maps/api/staticmap?%s", params.Encode())
}
