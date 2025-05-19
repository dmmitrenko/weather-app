package http

import (
	"errors"
	"net/http"

	"github.com/dmmitrenko/weather-app/internal/domain"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) error

func WithErrorHandling(fn AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err == nil {
			return
		}

		code, msg := http.StatusInternalServerError, "Internal server error"

		switch {
		case errors.Is(err, domain.ErrInvalidInput):
			code, msg = http.StatusBadRequest, "Invalid request"
		case errors.Is(err, domain.ErrAlreadySubscribed):
			code, msg = http.StatusConflict, "Email already subscribed"
		case errors.Is(err, domain.ErrSubscriptionNotFound):
			code, msg = http.StatusNotFound, "Token not found"
		case errors.Is(err, domain.ErrInvalidToken):
			code, msg = http.StatusBadRequest, "Invalid token"
		case errors.Is(err, domain.ErrCityNotFound):
			code, msg = http.StatusNotFound, "City not found"
		}

		http.Error(w, msg, code)
	}
}
