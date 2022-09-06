package octranspo

import (
	"errors"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"stop-checker.com/db/model"
)

type entry struct {
	Routes  map[string][]model.Bus
	Error   error
	Ready   sync.WaitGroup
	Created time.Time
}

type API struct {
	lock   sync.Mutex
	client *Client
	data   map[string]*entry
	ttl    time.Duration
}

func NewAPI(ttl time.Duration, client *Client) *API {
	return &API{
		lock:   sync.Mutex{},
		client: client,
		data:   map[string]*entry{},
		ttl:    ttl,
	}
}

func (api *API) StopData(stop model.Stop) (map[string][]model.Bus, error) {
	entry := api.getEntry(stop)
	entry.Ready.Wait()

	if entry.Error != nil {
		return nil, entry.Error
	}
	return entry.Routes, nil
}

// Request by stop code and route name
func (api *API) StopRouteData(stop model.Stop, routeName string) ([]model.Bus, error) {
	routes, err := api.StopData(stop)
	if err != nil {
		return nil, err
	}

	buses, ok := routes[routeName]
	if !ok {
		return nil, errors.New("route not found")
	}

	return buses, nil
}

func (api *API) getEntry(stop model.Stop) *entry {
	api.lock.Lock()
	defer api.lock.Unlock()

	entry, err := api.getCurrentEntry(stop.Code)
	if err != nil {
		return api.update(stop)
	}
	return entry
}

func (api *API) getCurrentEntry(stopCode string) (*entry, error) {
	entry, ok := api.data[stopCode]
	if !ok {
		return nil, errors.New("not found")
	}

	if time.Since(entry.Created) > api.ttl {
		return nil, errors.New("TTL exceeded")
	}

	return entry, nil
}

func (api *API) update(stop model.Stop) *entry {
	entry := &entry{
		Routes:  nil,
		Error:   nil,
		Created: time.Now(),
		Ready:   sync.WaitGroup{},
	}

	entry.Ready.Add(1)

	api.data[stop.Code] = entry

	go func() {
		t0 := time.Now()
		routes, err := api.client.Request(stop)
		entry.Routes = routes
		entry.Error = err
		entry.Ready.Done()

		if err != nil {
			log.Error().Err(err).
				Dur("request-duration", time.Since(t0)).
				Str("request-stop", stop.Code).
				Msg("failed to request live bus data from OC Transpo")
			return
		}

		log.Info().
			Dur("request-duration", time.Since(t0)).
			Str("request-stop", stop.Code).
			Msg("requested live bus data from OC Transpo")
	}()

	return entry
}
