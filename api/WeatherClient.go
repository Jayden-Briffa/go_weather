package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	apiKey string
	domain string
}

func NewClient(apiKey string, domain string) *Client {
	return &Client{
		apiKey: apiKey,
		domain: domain,
	}
}

func (c *Client) GetCityCoords(location string, country_code string) (float64, float64, error) {
	url := fmt.Sprintf("%s/geo/1.0/direct?q=%s,%s&limit=1&appid=%s", c.domain, location, country_code, c.apiKey)
	res, err := http.Get(url)

	if err != nil {
		panic("Error getting city coordinates")
	}

	defer res.Body.Close()

	// Any returned response should be a slice of data structures conforming to these fields
	var data []struct {
		Lat     float64 `json:"lat"`
		Lon     float64 `json:"lon"`
		Name    string  `json:"name"`
		Country string  `json:"country"`
	}

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return 0, 0, err
	}

	city := data[0]

	if !strings.EqualFold(city.Name, location) || !strings.EqualFold(city.Country, country_code) {
		panic("Could not find the city which was requested")
	}

	return city.Lat, city.Lon, err
}

func (c *Client) GetWeather(location string, country_code string) (string, float64, float64, error) {
	lat, lon, err := c.GetCityCoords(location, country_code)

	if err != nil {
		panic("Error getting city coordinates")
	}

	url := fmt.Sprintf("%s/data/2.5/weather?lat=%f&lon=%f&appid=%s", c.domain, lat, lon, c.apiKey)
	res, err := http.Get(url)

	if err != nil {
		panic("Error getting weather response")
	}

	defer res.Body.Close()

	var data struct {
		Weather []struct {
			Description string `json:"description"`
		} `json:"weather"`

		Wind struct {
			Speed float64 `json:"speed"`
		} `json:"wind"`

		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
	}

	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		panic("Error decoding weather response body")
	}

	description := data.Weather[0].Description
	windSpeed := data.Wind.Speed
	temp := data.Main.Temp

	return description,
		windSpeed,
		temp,
		err

}
