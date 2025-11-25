package main

import (
	"fmt"
	WeatherClient "go_weather/api"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Env struct {
	API_KEY string
}

// func get_coords(location string)

func load_env() (*Env, error) {
	cfg := &Env{
		API_KEY: os.Getenv("API_KEY"),
	}

	return cfg, nil
}

func main() {
	env, err := load_env()

	if err != nil {
		panic("Error loading environment")
	}

	Client := WeatherClient.NewClient(env.API_KEY, "http://api.openweathermap.org")
	description, windSpeed, temp, err := Client.GetWeather("London", "gb")

	fmt.Print(description, windSpeed, temp)
}
