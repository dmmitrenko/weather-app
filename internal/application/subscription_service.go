package application

import (
	"context"
	"fmt"

	"github.com/dmmitrenko/weather-app/internal/domain"
	"github.com/dmmitrenko/weather-app/internal/utils"
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
	alreadyExists, err := s.repository.IsExists(ctx, email)
	if err != nil {
		return fmt.Errorf("checking existing subscription: %w", err)
	}
	if alreadyExists {
		return domain.ErrAlreadySubscribed
	}

	token, err := domain.GenerateToken()
	if err != nil {
		return err
	}

	subject, body, err := utils.BuildSubscriptionConfirmationMessage(email, city, token)
	if err != nil {
		fmt.Printf("building message for %s failed: %v\n", email, err)
	}
	s.emailSender.Send(ctx, email, subject, body)

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
