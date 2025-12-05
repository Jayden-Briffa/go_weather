package presentation

import (
	"fmt"
	Model "go_weather/model"
)

func FormatWeather(weather Model.Weather) string {
	weatherStr := ""

	weatherStr = fmt.Sprintf("%v | %vmph | %vdeg", weather.Description, weather.WindSpeed, weather.Temp)

	return weatherStr
}
