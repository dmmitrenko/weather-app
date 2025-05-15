package domain

import "errors"

var (
	ErrInvalidInput         = errors.New("invalid input")
	ErrAlreadySubscribed    = errors.New("email already subscribed")
	ErrSubscriptionNotFound = errors.New("subscription not found")
	ErrInvalidToken         = errors.New("invalid token")
	ErrCityNotFound         = errors.New("city not found")
)
