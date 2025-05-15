package domain

import "github.com/google/uuid"

type Subscription struct {
	Id        uuid.UUID `json:"id"`
	Frequency Frequency `json:"frequency"`
	Email     string    `json:"email"`
	Confirmed bool      `json:"confirmed"`
	Token     string    `json:"token"`
	City      string    `json:"city"`
}

type SubscriptionRepository interface {
	Create(sub *Subscription) error
	GetByToken(token string) (*Subscription, error)
	ConfirmByToken(token string) error
	DeleteByToken(token string) error
	GetActiveSubscriptions(frequency Frequency) error
}
