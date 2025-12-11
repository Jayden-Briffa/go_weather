package api

import (
	"encoding/json"
	"fmt"
	Model "go_weather/model"
	Presentation "go_weather/presentation"
	"net/http"
	"strings"
	"sync"
	"time"
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

func (c *Client) GetCityCoords(city string, country_code string) (float64, float64, error) {
	url := fmt.Sprintf("%s/geo/1.0/direct?q=%s,%s&limit=1&appid=%s", c.domain, city, country_code, c.apiKey)
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

	weather := data[0]

	if !strings.EqualFold(weather.Name, city) || !strings.EqualFold(weather.Country, country_code) {
		panic("Could not find the city which was requested")
	}

	return weather.Lat, weather.Lon, err
}

func (c *Client) GetWeather(city Model.City) (Model.Weather, error) {
	lat, lon, err := c.GetCityCoords(city.Name, city.Country_code)

	if err != nil {
		panic("Error getting city coordinates")
	}

	url := fmt.Sprintf("%s/data/2.5/weather?lat=%f&lon=%f&appid=%s", c.domain, lat, lon, c.apiKey)
	res, err := http.Get(url)

	if err != nil {
		return Model.Weather{}, err
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

	time.Sleep(600 * time.Millisecond)

	return Model.Weather{
		Description: data.Weather[0].Description,
		WindSpeed:   data.Wind.Speed,
		Temp:        data.Main.Temp,
		City:        city,
	}, err
}

func (client *Client) GetWeatherAsync(city Model.City, ch chan Model.Weather){
	weather, err := client.GetWeather(city)

	if err != nil {
		fmt.Printf("Error getting weather for %v, %v: %v\n", city.Name, city.Country_code, err)
	}

	ch<-weather
}

func (client *Client) StreamWeather(cities []Model.City) {
	i := 0
	for {
		i++

		var wg sync.WaitGroup
		ch := make(chan Model.Weather)

		for _, city := range cities {
			wg.Go(func() {client.GetWeatherAsync(city, ch)})
		}

		go func() {
			wg.Wait()
			close(ch)
		}()

		var citiesWeather string
		for weather := range ch {
			citiesWeather += Presentation.FormatWeather(weather)
			citiesWeather += " // "
		}

		fmt.Printf("%v: %v\r", i, citiesWeather)
		time.Sleep(1 * time.Second)
	}
}
