package main

import (
	Api "go_weather/api"
	Data "go_weather/data"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Env struct {
	API_KEY string
}

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

	client := Api.NewClient(env.API_KEY, "http://api.openweathermap.org")

	cities := Data.Cities

	client.StreamWeather(cities)
}
