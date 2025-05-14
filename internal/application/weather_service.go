package application

import (
	"context"

	"github.com/dmmitrenko/weather-app/internal/domain"
)

type WeatherService struct {
	provider domain.WeatherProvider
}

func NewWeatherService(p domain.WeatherProvider) *WeatherService {
	return &WeatherService{
		provider: p,
	}
}

func (s *WeatherService) GetCurrentWeather(ctx context.Context, city string) (domain.Weather, error) {
	return s.provider.GetCurrentWeather(ctx, city)
}
