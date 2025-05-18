package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type WeatherService struct {
	Address string
	Key     string
}

type WeatherResponse struct {
	Current `json:"current"`
}

type Current struct {
	Temperature float64 `json:"temp_c"`
	Humidity    float64 `json:"humidity"`
	Condition   `json:"condition"`
}

type Condition struct {
	Text string `json:"text"`
}

type Weather struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Description string  `json:"description"`
}

func (w WeatherService) GetWeather(city string) (Weather, error) {
	var weather Weather
	var weatherResponse WeatherResponse
	url := fmt.Sprintf(w.Address, w.Key, city)

	resp, err := http.Get(url)

	if resp.StatusCode == http.StatusBadRequest {
		return weather, fmt.Errorf("city `%s` not exists", city)
	}

	if err != nil {
		return weather, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return weather, err
	}

	if err := json.Unmarshal(body, &weatherResponse); err != nil {
		return weather, err
	}

	weather.Description = weatherResponse.Text
	weather.Humidity = weatherResponse.Humidity
	weather.Temperature = weatherResponse.Temperature

	return weather, nil
}
