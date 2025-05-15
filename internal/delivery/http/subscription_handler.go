package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/dmmitrenko/weather-app/internal/application"
	"github.com/dmmitrenko/weather-app/internal/domain"
)

type SubscriptionHandler struct {
	svc *application.SubscriptionService
}

func NewSubscriptionHandler(r *mux.Router, svc *application.SubscriptionService) {
	h := &SubscriptionHandler{
		svc: svc,
	}

	r.Handle("/subscribe", WithErrorHandling(h.Subscribe)).Methods(http.MethodPost)
	r.Handle("/confirm/{token}", WithErrorHandling(h.Confirm)).Methods(http.MethodGet)
	r.Handle("/unsubscribe/{token}", WithErrorHandling(h.Unsubscribe)).Methods(http.MethodGet)
}

func (h SubscriptionHandler) Subscribe(w http.ResponseWriter, r *http.Request) error {
	var req struct {
		Email     string `json:"email"`
		City      string `json:"city"`
		Frequency string `json:"frequency"`
	}

	ct := r.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return domain.ErrInvalidInput
		}
	} else {
		if err := r.ParseForm(); err != nil {
			return domain.ErrInvalidInput
		}
		req.Email = r.FormValue("email")
		req.City = r.FormValue("city")
		req.Frequency = r.FormValue("frequency")
	}

	if req.Email == "" || req.City == "" || req.Frequency == "" {
		return domain.ErrInvalidInput
	}

	freq, err := domain.ParseFrequency(req.Frequency)
	if err != nil {
		return domain.ErrInvalidInput
	}

	return h.svc.Subscribe(r.Context(), req.Email, freq, req.City)
}

func (h SubscriptionHandler) Confirm(w http.ResponseWriter, r *http.Request) error {
	token := mux.Vars(r)["token"]
	if token == "" {
		return domain.ErrInvalidToken
	}
	return h.svc.ConfirmSubscription(r.Context(), token)
}

func (h SubscriptionHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) error {
	token := mux.Vars(r)["token"]
	if token == "" {
		return domain.ErrInvalidToken
	}
	return h.svc.Unsubscribe(r.Context(), token)
}
