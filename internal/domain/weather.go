package domain

import (
	"context"
	"time"
)

type Weather struct {
	City        string    `json:"city"`
	Region      string    `json:"region"`
	Country     string    `json:"country"`
	LocalTime   time.Time `json:"local_time"`
	LastUpdated time.Time `json:"last_updated"`
	Temperature float64   `json:"temperature"`
	Humidity    int       `json:"humidity"`
	Description string    `json:"description"`
}

type WeatherProvider interface {
	GetCurrentWeather(ctx context.Context, city string) (Weather, error)
}

type EmailSender interface {
	Send(ctx context.Context, to, subject, body string) error
}
