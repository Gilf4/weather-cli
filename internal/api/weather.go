package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
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

func GetWeatherByCoords(lat, lon float64) error {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("WEATHER_API_KEY is not set in environment variables")
	}
	url := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%s&q=%f,%f&days=1&aqi=no&alerts=no", apiKey, lat, lon)
	return getWeather(url)
}

func GetWeatherByCity(city string) error {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("WEATHER_API_KEY is not set in environment variables")
	}
	url := fmt.Sprintf("http://api.weatherapi.com/v1/forecast.json?key=%s&q=%s&days=1&aqi=no&alerts=no", apiKey, city)
	return getWeather(url)
}

func getWeather(url string) error {
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error making the request to the API: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("API unavailable: status code %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Error reading the response body: %v", err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return fmt.Errorf("Error parsing the JSON response: %v", err)
	}

	location := weather.Location
	current := weather.Current
	hours := weather.Forecast.ForecastDay[0].Hour

	fmt.Printf(
		"%s, %s: %.0f°C, %s\n",
		location.Name, location.Country,
		current.TempC,
		current.Condition.Text)

	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0)

		if time.Now().After(date) {
			continue
		}

		msg := fmt.Sprintf(
			"%s, %.0f°C, %.0f%% chance of rain, %s\n",
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

	return nil
}
