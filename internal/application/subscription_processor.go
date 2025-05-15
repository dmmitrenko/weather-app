package application

import (
	"context"
	"fmt"

	"github.com/dmmitrenko/weather-app/internal/domain"
)

type SubscriptionProcessor struct {
	Repo            domain.SubscriptionRepository
	Sender          domain.EmailSender
	WeatherProvider domain.WeatherProvider
}

func (p *SubscriptionProcessor) Process(ctx context.Context, freq domain.Frequency) error {
	subs, err := p.Repo.GetActiveSubscriptions(ctx, freq)
	if err != nil {
		return fmt.Errorf("fetch subs: %w", err)
	}
	for _, s := range subs {
		weather, err := p.WeatherProvider.GetCurrentWeather(ctx, s.City)
		if err != nil {
			fmt.Printf("notify %s failed: %v\n", s.Email, err)
		}

		if err := p.Sender.Send(ctx, s.Email, "Subscription update", weather.Country); err != nil {
			fmt.Printf("notify %s failed: %v\n", s.Email, err)
		}
	}
	return nil
}
