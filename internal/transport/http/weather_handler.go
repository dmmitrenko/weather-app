package http

import (
	"encoding/json"
	"net/http"

	"github.com/dmmitrenko/weather-app/internal/domain"
	"github.com/gorilla/mux"
)

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

	r.HandleFunc("/api/weather", h.GetCurrentWeather).Methods("GET")
}

func (h *WeatherHandler) GetCurrentWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "city required", http.StatusBadRequest)
		return
	}

	weather, err := h.weatherProvider.GetCurrentWeather(r.Context(), city)
	if err != nil {
		if err == domain.ErrCityNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	resp := WeatherResponse{
		Temperature: weather.Temperature,
		Humidity:    float64(weather.Humidity),
		Description: weather.Description,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
