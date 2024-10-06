package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
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

	fmt.Printf(
		"\nWeather for %s, %s: %.1f°C, %s\n\n",
		weather.Location.Name,
		weather.Location.Country,
		weather.Current.TempC,
		weather.Current.Condition.Text,
	)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Time", "Temp (°C)", "Chance of Rain (%)", "Condition"})

	// Loop through hourly forecast and fill the table
	for _, hour := range weather.Forecast.ForecastDay[0].Hour {
		hourTime := time.Unix(hour.TimeEpoch, 0)
		if time.Now().After(hourTime) {
			continue
		}

		// Add a new row for each hour of the forecast
		row := []string{
			hourTime.Format("15:04"),
			fmt.Sprintf("%.1f", hour.TempC),
			fmt.Sprintf("%.0f", hour.ChanceOfRain),
			hour.Condition.Text,
		}

		// Highlight rows in red if chance of rain is more than 40%
		if hour.ChanceOfRain > 40 {
			color.Set(color.FgRed)
			table.Rich(row, []tablewriter.Colors{
				{}, {}, {tablewriter.FgRedColor}, {},
			})
			color.Unset()
		} else {
			table.Append(row)
		}
	}

	// Render and display the table in the terminal
	table.Render()

	// Print a tip to the user regarding the meaning of red rows
	fmt.Println("\nTip: Red rows indicate a high chance of rain. Be prepared!")
	return nil
}
