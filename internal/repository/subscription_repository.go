package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dmmitrenko/weather-app/internal/domain"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(
	ctx context.Context,
	sub *domain.Subscription,
) error {
	const query = `
        INSERT INTO subscriptions
            (email, city, frequency, token, confirmed)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
	return r.db.
		QueryRowContext(
			ctx,
			query,
			sub.Email,
			sub.City,
			string(sub.Frequency),
			sub.Token,
			sub.Confirmed,
		).
		Scan(&sub.Id)
}

func (r *SubscriptionRepository) GetByToken(
	ctx context.Context,
	token string,
) (*domain.Subscription, error) {
	const query = `
        SELECT id, email, city, frequency, token_hash, confirmed
        FROM subscriptions
        WHERE token_hash = $1
    `
	row := r.db.QueryRowContext(ctx, query, token)

	var s domain.Subscription
	var freq string
	if err := row.Scan(
		&s.Id,
		&s.Email,
		&s.City,
		&freq,
		&s.Token,
		&s.Confirmed,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSubscriptionNotFound
		}
		return nil, err
	}
	s.Frequency = domain.Frequency(freq)
	return &s, nil
}

func (r *SubscriptionRepository) ConfirmByToken(
	ctx context.Context,
	token string,
) error {
	const query = `
        UPDATE subscriptions
        SET confirmed = TRUE
        WHERE token = $1
    `
	res, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}
	if cnt, err := res.RowsAffected(); err != nil {
		return err
	} else if cnt == 0 {
		return domain.ErrSubscriptionNotFound
	}
	return nil
}

func (r *SubscriptionRepository) DeleteByToken(ctx context.Context, token string) error {
	const query = `
        DELETE FROM subscriptions
        WHERE token = $1
    `

	res, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}

	if cnt, err := res.RowsAffected(); err != nil {
		return err
	} else if cnt == 0 {
		return domain.ErrSubscriptionNotFound
	}

	return nil
}

func (r *SubscriptionRepository) GetActiveSubscriptions(
	ctx context.Context,
	freq domain.Frequency,
) ([]domain.Subscription, error) {
	const query = `
        SELECT id, email, city, token, confirmed
        FROM subscriptions
        WHERE frequency = $1 AND confirmed = TRUE
    `
	rows, err := r.db.QueryContext(ctx, query, string(freq))
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var subs []domain.Subscription

	for rows.Next() {
		var s domain.Subscription

		if err := rows.Scan(
			&s.Id,
			&s.Email,
			&s.City,
			&s.Token,
			&s.Confirmed,
		); err != nil {
			return nil, err
		}

		s.Frequency = freq
		subs = append(subs, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subs, nil
}
