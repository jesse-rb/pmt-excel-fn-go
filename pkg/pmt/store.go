package pmt

import (
	"context"
	"time"
)

type PMTHistory struct {
	LoanAmountCents int64
	InterestRate    float64
	NumPayments     int32
	PMTCents        int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Store interface {
	Create(ctx context.Context, loanAmountCents int64, interestRate float64, numPayments int32, pmtCents int64) (*PMTHistory, error)
}
