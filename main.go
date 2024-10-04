package main

import (
	"os"
	"weather/internal/api"
	"weather/internal/location"

	"github.com/fatih/color"
)

func main() {
	var lat, lon float64
	var city string
	var err error

	if len(os.Args) > 1 {
		city = os.Args[1]
		color.Green("Using city: %s", city)
		err = api.GetWeatherByCity(city)
	} else {
		lat, lon, err = location.GetLocation()
		if err != nil {
			color.Red("Error getting location: %v", err)
			return
		}
		color.Green("Using coordinates: %f, %f", lat, lon)
		err = api.GetWeatherByCoords(lat, lon)
	}

	if err != nil {
		color.Red("Error: %v", err)
	}
}
