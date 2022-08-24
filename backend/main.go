package main

import (
	"fmt"
	"time"

	"stop-checker.com/features/octranspo"
	"stop-checker.com/features/travel"
)

func printLegs(legs []*travel.Leg) {
	for _, leg := range legs {
		fmt.Println(leg.String())
	}
}

func main() {
	api := octranspo.NewAPI(time.Second*30, &octranspo.Client{
		Endpoint:          "https://api.octranspo1.com/v2.0/GetNextTripsForStopAllRoutes",
		OCTRANSPO_APP_ID:  "13d12d72",
		OCTRANSPO_API_KEY: "508a0741b6945609192422d77f3a1da4",
	})

	go func() {
		api.StopData("8810")
	}()

	res, err := api.StopData("8810")

	fmt.Println(res)
	fmt.Println(err)
}
