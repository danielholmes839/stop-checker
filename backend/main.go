package main

import (
	"fmt"
	"time"

	"stop-checker.com/features/octranspo"
)

func main() {
	api := octranspo.NewAPI(time.Second*30, &octranspo.Client{
		Endpoint:          "https://api.octranspo1.com/v2.0/GetNextTripsForStopAllRoutes",
		OCTRANSPO_APP_ID:  "13d12d72",
		OCTRANSPO_API_KEY: "508a0741b6945609192422d77f3a1da4",
	})

	t0 := time.Now()

	res, err := api.StopData("8810")

	fmt.Println(time.Since(t0))

	t1 := time.Now()

	for i := 0; i < 10; i++ {
		api.StopData("8810")
	}

	fmt.Println(time.Since(t1))

	fmt.Println(res)
	fmt.Println(err)
}
