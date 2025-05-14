package domain

import "github.com/google/uuid"

type Subscription struct {
	Id        uuid.UUID `json:"id"`
	Frequency Frequency `json:"frequency"`
	Email     string    `json:"email"`
	Confirmed bool      `json:"confirmed"`
	Token     string    `json:"token"`
}
