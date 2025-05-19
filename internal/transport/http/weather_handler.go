package http

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/dmmitrenko/weather-app/internal/domain"
	"github.com/gorilla/mux"
)

var cityRegex = regexp.MustCompile(`^[A-Za-z]+$`)

type WeatherResponse struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Description string  `json:"description"`
}

type WeatherHandler struct {
	weatherProvider domain.WeatherProvider
}

func NewWeatherHandler(r *mux.Router, w domain.WeatherProvider) {
	h := &WeatherHandler{
		weatherProvider: w,
	}

	r.HandleFunc("/api/weather", WithErrorHandling(h.GetCurrentWeather)).Methods("GET")
}

func (h *WeatherHandler) GetCurrentWeather(w http.ResponseWriter, r *http.Request) error {
	city := r.URL.Query().Get("city")
	if city == "" || !cityRegex.MatchString(city) {
		return domain.ErrInvalidInput
	}

	weather, err := h.weatherProvider.GetCurrentWeather(r.Context(), city)
	if err != nil {
		return err
	}

	resp := WeatherResponse{
		Temperature: weather.Temperature,
		Humidity:    float64(weather.Humidity),
		Description: weather.Description,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)

	return nil
}
