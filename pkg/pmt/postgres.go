package pmt

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, loanAmountCents int64, interestRate float64, numPayments int32, pmtCents int64) (*PMTHistory, error) {
	new := &PMTHistory{}
	args := []interface{}{loanAmountCents, interestRate, numPayments, pmtCents}
	err := pgxscan.Get(ctx, r.db, new, "INSERT INTO pmt_histories(loan_amount_cents, interest_rate, num_payments, pmt_cents), VALUES($1, $2, $3, $4);", args)
	if err != nil {
		return nil, fmt.Errorf("error creating pmt history: %w", err)
	}
	return new, nil
}
