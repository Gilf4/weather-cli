package location

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UserLocation struct {
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	City    string  `json:"city"`
	Country string  `json:"country"`
}

func GetLocation() (float64, float64, error) {
	resp, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return 0, 0, fmt.Errorf("Error making request to location API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return 0, 0, fmt.Errorf("Failed to get location: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("Error reading the response body: %v", err)
	}

	var usrLocation UserLocation
	err = json.Unmarshal(body, &usrLocation)
	if err != nil {
		return 0, 0, fmt.Errorf("Error parsing the JSON response: %v", err)
	}

	return usrLocation.Lat, usrLocation.Lon, nil
}
