package geoip

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sentinel-official/health-check/libs/geoip/types"
)

func Location(transport *http.Transport) (*types.GeoIPLocation, error) {
	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
	}

	resp, err := client.Get("http://ip-api.com/json")
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	var m map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	return &types.GeoIPLocation{
		City:      m["city"].(string),
		Country:   m["country"].(string),
		IP:        m["query"].(string),
		Latitude:  m["lat"].(float64),
		Longitude: m["lon"].(float64),
	}, nil
}
