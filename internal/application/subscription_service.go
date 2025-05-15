package application

import (
	"context"

	"github.com/dmmitrenko/weather-app/internal/domain"
)

type SubscriptionService struct {
	repository  domain.SubscriptionRepository
	emailSender domain.EmailSender
}

func NewSubscriptionService(r domain.SubscriptionRepository, emailSender domain.EmailSender) *SubscriptionService {
	return &SubscriptionService{
		repository:  r,
		emailSender: emailSender,
	}
}

func (s *SubscriptionService) ConfirmSubscription(ctx context.Context, token string) error {
	return s.repository.ConfirmByToken(ctx, token)
}

func (s *SubscriptionService) Subscribe(ctx context.Context, email string, freq domain.Frequency, city string) error {
	token, err := domain.GenerateToken()
	if err != nil {
		return err
	}

	s.emailSender.Send(ctx, email, "Subscription confirmation", token)

	sub := domain.Subscription{
		City:      city,
		Frequency: freq,
		Email:     email,
		Token:     token,
	}

	s.repository.Create(ctx, &sub)
	return nil
}

func (s *SubscriptionService) Unsubscribe(ctx context.Context, token string) error {
	return s.repository.DeleteByToken(ctx, token)
}
