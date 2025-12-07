package presentation

import (
	"fmt"
	Model "go_weather/model"
	"strings"
)

func FormatWeather(weather Model.Weather) string {
	weatherStr := ""

	weatherStr = fmt.Sprintf("%v, %v: %v | %vmph | %vdeg", 
		strings.ToUpper(weather.City.Name), 
		strings.ToUpper(weather.City.Country_code), 
		weather.Description, 
		weather.WindSpeed, 
		weather.Temp,
	)

	return weatherStr
}
