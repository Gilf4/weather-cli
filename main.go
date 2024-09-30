package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net/http"
	"os"
	"time"
)

type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`

	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`

	Forecast struct {
		ForecastDay []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	var city string

	if len(os.Args) > 1 {
		city = os.Args[1]
	} else {
		city = "Novgorod"
	}

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		color.Red("API_KEY is not set in env variables")
		return
	}

	res, err := http.Get(fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%s&q=%s&days=1&aqi=no&alerts=no", apiKey, city))
	if err != nil {
		color.Red("Error when requesting API: %v", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		color.Red("API недоступен: статус %d", res.StatusCode)
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		color.Red("Error reading response: %v", err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		color.Red("Error parsing JSON: %v", err)
		return
	}

	location := weather.Location
	current := weather.Current
	hours := weather.Forecast.ForecastDay[0].Hour

	fmt.Printf(
		"%s, %s: %.0fC, %s\n",
		location.Name, location.Country,
		current.TempC,
		current.Condition.Text)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		if time.Now().After(date) {
			continue
		}

		msg := fmt.Sprintf(
			"%s, %.0f°C, %.0f%%, %s\n",
			date.Format("15:04"),
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text)
		if hour.ChanceOfRain < 40 {
			fmt.Print(msg)
		} else {
			color.Red(msg)
		}
	}
}
