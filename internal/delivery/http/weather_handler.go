package http

import (
	"encoding/json"
	"net/http"

	"github.com/dmmitrenko/weather-app/internal/application"
	"github.com/gorilla/mux"
)

type WeatherHandler struct {
	weatherService *application.WeatherService
}

func NewWeatherHandler(r *mux.Router, w *application.WeatherService) {
	h := &WeatherHandler{
		weatherService: w,
	}

	r.HandleFunc("/weather", h.GetCurrentWeather).Methods("GET")
}

func (h *WeatherHandler) GetCurrentWeather(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "city required", http.StatusBadRequest)
		return
	}

	weather, err := h.weatherService.GetCurrentWeather(ctx, city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather)
}
