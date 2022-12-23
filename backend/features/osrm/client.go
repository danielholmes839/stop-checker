package osrm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"stop-checker.com/db/model"
)

/*
Open Source Routing Machine client
- http://project-osrm.org
*/
type Client struct {
	host string
}

func NewClient(host string) *Client {
	return &Client{host: host}
}

func (c *Client) requestWalkingDirections(origin, destination model.Location) (*osrmResponse, error) {
	url := fmt.Sprintf("%s/route/v1/driving/%f,%f;%f,%f?alternatives=false&annotations=true&overview=false&steps=true",
		c.host, origin.Longitude, origin.Latitude, destination.Longitude, destination.Latitude)

	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("osm request error: %w", err)
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	directions := &osrmResponse{}
	if err := json.Unmarshal(data, directions); err != nil {
		return nil, err
	}

	return directions, nil
}

func (c *Client) WalkingDirections(origin, destination model.Location) ([]model.Location, error) {
	directions, err := c.requestWalkingDirections(origin, destination)
	if err != nil {
		return nil, fmt.Errorf("osm request error: %w", err)
	}

	if directions.Code != "Ok" {
		return nil, fmt.Errorf("osm error: %s", directions.Code)
	}

	// all points in the directions
	points := []model.Location{}

	for _, route := range directions.Routes {
		for _, leg := range route.Legs {
			for _, step := range leg.Steps {
				for _, intersection := range step.Intersections {
					// add point
					points = append(points, model.Location{
						Longitude: intersection.Location[0],
						Latitude:  intersection.Location[1],
					})
				}
			}
		}
	}

	return points, nil
}
