package application

import (
	"context"
	"fmt"

	"github.com/dmmitrenko/weather-app/internal/domain"
	"github.com/dmmitrenko/weather-app/internal/utils"
)

type SubscriptionProcessor struct {
	Repo            domain.SubscriptionRepository
	Sender          domain.EmailSender
	WeatherProvider domain.WeatherProvider
}

func (p *SubscriptionProcessor) Process(ctx context.Context, freq domain.Frequency) error {
	subs, err := p.Repo.GetActiveSubscriptions(ctx, freq)
	if err != nil {
		return fmt.Errorf("fetch subscriptions: %w", err)
	}
	for _, s := range subs {
		weather, err := p.WeatherProvider.GetCurrentWeather(ctx, s.City)
		if err != nil {
			fmt.Printf("notify %s failed: %v\n", s.Email, err)
			continue
		}

		subject, body, err := utils.BuildWeatherUpdateMessage(s.Email, s.City, weather)
		if err != nil {
			fmt.Printf("building message for %s failed: %v\n", s.Email, err)
			continue
		}

		if err := p.Sender.Send(ctx, s.Email, subject, body); err != nil {
			fmt.Printf("notify %s failed: %v\n", s.Email, err)
		}
	}
	return nil
}
