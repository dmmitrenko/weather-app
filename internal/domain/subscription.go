package domain

import (
	"context"

	"github.com/google/uuid"
)

type Subscription struct {
	Id        uuid.UUID `json:"id"`
	Frequency Frequency `json:"frequency"`
	Email     string    `json:"email"`
	Confirmed bool      `json:"confirmed"`
	Token     string    `json:"token"`
	City      string    `json:"city"`
}

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *Subscription) error
	GetByToken(ctx context.Context, token string) (*Subscription, error)
	ConfirmByToken(ctx context.Context, token string) error
	DeleteByToken(ctx context.Context, token string) error
	GetActiveSubscriptions(ctx context.Context, freq Frequency) ([]Subscription, error)
}
