package main

import (
	"log"
	"net/http"

	"github.com/dmmitrenko/weather-app/configs"
	"github.com/dmmitrenko/weather-app/internal/application"
	delivery "github.com/dmmitrenko/weather-app/internal/delivery/http"
	weatherapi "github.com/dmmitrenko/weather-app/internal/infrastructure/weather-api"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := configs.Load("configs/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	client := weatherapi.NewClient(cfg.WeatherAPI.Key)

	weatherService := application.NewWeatherService(client)

	r := mux.NewRouter()
	delivery.NewWeatherHandler(r, weatherService)

	delivery.RegisterStatic(r)

	log.Printf("Listening on %s", cfg.Server.Address)
	log.Fatal(http.ListenAndServe(cfg.Server.Address, r))
}
