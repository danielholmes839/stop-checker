package octranspo

import (
	"time"

	"stop-checker.com/db/model"
)

type Bus struct {
	Arrival     time.Time
	LastUpdated time.Time
	Location    *model.Location
}
