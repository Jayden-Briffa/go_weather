package model

type Weather struct {
	Description string
	WindSpeed float64
	Temp float64

	City City
}