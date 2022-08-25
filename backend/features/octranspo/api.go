package octranspo

import (
	"errors"
	"sync"
	"time"

	"stop-checker.com/model"
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

func (api *API) StopData(stopCode string) (map[string][]model.Bus, error) {
	entry := api.getEntry(stopCode)
	entry.Ready.Wait()

	if entry.Error != nil {
		return nil, entry.Error
	}
	return entry.Routes, nil
}

// Request by stop code and route name
func (api *API) StopRouteData(stopCode string, routeName string) ([]model.Bus, error) {
	routes, err := api.StopData(stopCode)
	if err != nil {
		return nil, err
	}

	buses, ok := routes[routeName]
	if !ok {
		return nil, errors.New("route not found")
	}

	return buses, nil
}

func (api *API) getEntry(stopCode string) *entry {
	api.lock.Lock()
	defer api.lock.Unlock()

	entry, err := api.getCurrentEntry(stopCode)
	if err != nil {
		return api.update(stopCode)
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

func (api *API) update(stopCode string) *entry {
	entry := &entry{
		Routes:  nil,
		Error:   nil,
		Created: time.Now(),
		Ready:   sync.WaitGroup{},
	}

	entry.Ready.Add(1)

	api.data[stopCode] = entry

	go func() {
		routes, err := api.client.Request(stopCode)
		entry.Routes = routes
		entry.Error = err
		entry.Ready.Done()
	}()

	return entry
}
