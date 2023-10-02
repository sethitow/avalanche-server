package aaa_api

import (
	"encoding/json"
	"errors"
	"net/http"
)

const (
	AvalancheOrgBaseURI = "https://api.avalanche.org/v2/public/products/map-layer/"
)

var (
	ErrUpstreamAPIError = errors.New("could not get forecast from avalanche.org")
)

//go:generate mockery --name=Requester
type Requester interface {
	GetForecastByCenter(string) (Root, error)
}

type APIRequester struct{}

func (APIRequester) GetForecastByCenter(centerId string) (Root, error) {
	resp, err := http.Get(AvalancheOrgBaseURI + centerId)
	if err != nil {
		return Root{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return Root{}, ErrUpstreamAPIError
	}

	decoder := json.NewDecoder(resp.Body)
	root := Root{}
	decoder.Decode(&root)
	return root, nil
}
