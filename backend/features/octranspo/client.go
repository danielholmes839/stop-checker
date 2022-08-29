package octranspo

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"stop-checker.com/db/model"
)

type Client struct {
	Endpoint          string
	OCTRANSPO_APP_ID  string
	OCTRANSPO_API_KEY string
}

func (c *Client) Request(stopCode string) (map[string][]model.Bus, error) {
	q := url.Values{}
	q.Add("appID", c.OCTRANSPO_APP_ID)
	q.Add("apiKey", c.OCTRANSPO_API_KEY)
	q.Add("stopNo", stopCode)
	q.Add("format", "XML")

	url := fmt.Sprintf("%s?%s", c.Endpoint, q.Encode())

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	parsed := &soapEnvelope{}
	err = xml.Unmarshal(data, parsed)
	if err != nil {
		return nil, err
	}

	return parseResults(parsed.Body.Response.Results), err
}

type soapEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    soapBody `xml:"Body"`
}

type soapBody struct {
	Response responseGetRouteSummaryForStopResponse `xml:"GetRouteSummaryForStopResponse"`
}

type responseGetRouteSummaryForStopResponse struct {
	Results []responseGetRouteSummaryForStopResult `xml:"GetRouteSummaryForStopResult"`
}

type responseGetRouteSummaryForStopResult struct {
	StopNo    string         `xml:"StopNo"`
	StopLabel string         `xml:"StopLabel"`
	Error     string         `xml:"Error"`
	Routes    responseRoutes `xml:"Routes"`
}

type responseRoutes struct {
	Routes []responseRoute `xml:"Route"`
}
type responseRoute struct {
	RouteNo    string        `xml:"RouteNo"`    // route name
	RouteLabel string        `xml:"RouteLabel"` // trip headsign
	Direction  string        `xml:"Direction"`
	Error      string        `xml:"Error"`
	Trips      responseTrips `xml:"Trips"`
}

type responseTrips struct {
	Trips []responseTrip `xml:"Trip"`
}

type responseTrip struct {
	Destination          string `xml:"TripDestination"`
	StartTime            string `xml:"TripStartTime"`
	AdjustedScheduleTime string `xml:"AdjustedScheduleTime"`
	AdjustmentAge        string `xml:"AdjustmentAge"`
	LastTripOfSchedule   bool   `xml:"LastTripOfSchedule"`
	BusType              string `xml:"BusType"`
	Longitude            string `xml:"Longitude"`
	Latitude             string `xml:"Latitude"`
}

func parseResults(results []responseGetRouteSummaryForStopResult) map[string][]model.Bus {
	data := map[string][]model.Bus{}

	for _, result := range results {
		for _, route := range result.Routes.Routes {
			buses := []model.Bus{}
			for _, trip := range route.Trips.Trips {
				buses = append(buses, parseTrip(trip))
			}
			data[route.RouteNo] = buses
		}
	}

	return data
}

func parseTime(str string) time.Time {
	hours, _ := strconv.Atoi(str[0:2])
	hours %= 24

	minutes, _ := strconv.Atoi(str[3:5])

	return time.Date(0, 0, 0, hours, minutes, 0, 0, time.UTC)
}

func parseTripLocation(trip responseTrip) *model.Location {
	lat, err := strconv.ParseFloat(trip.Latitude, 64)
	if err != nil {
		return nil
	}

	lon, err := strconv.ParseFloat(trip.Longitude, 64)
	if err != nil {
		return nil
	}

	return &model.Location{
		Latitude:  lat,
		Longitude: lon,
	}
}

// return the trip arrival and last updated time
func parseTripArrival(trip responseTrip) (time.Time, time.Time) {
	// get age
	f, _ := strconv.ParseFloat(trip.AdjustmentAge, 64)
	age := int(math.Round(f))
	if age < 0 {
		age = 0
	}

	// get arrival time
	duration, _ := strconv.Atoi(trip.AdjustedScheduleTime)

	now := time.Now().In(time.Local)
	arrival := now.Add(time.Minute * time.Duration(duration))
	lastUpdated := now.Add(-time.Minute * time.Duration(age))

	return arrival, lastUpdated
}

func parseTrip(trip responseTrip) model.Bus {
	loc := parseTripLocation(trip)
	arrival, lastUpdated := parseTripArrival(trip)

	return model.Bus{
		Headsign:    trip.Destination,
		Arrival:     model.NewTimeFromDateTime(arrival),
		LastUpdated: model.NewTimeFromDateTime(lastUpdated),
		Location:    loc,
	}
}
