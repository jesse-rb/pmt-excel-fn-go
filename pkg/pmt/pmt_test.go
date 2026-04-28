package pmt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalcPMT(t *testing.T) {
	s := NewService()
	t.Run("happy", func(t *testing.T) {
		cases := []struct {
			loanAmountCents int64
			interestRate    float64
			numPayments     int32
			expected        int64
		}{
			{loanAmountCents: 1000, interestRate: 0.1, numPayments: 1, expected: 1100},
			{loanAmountCents: 1000, interestRate: 0.25, numPayments: 2, expected: 694},
			{loanAmountCents: 123456, interestRate: 0.05, numPayments: 1, expected: 129629},
			{loanAmountCents: 130025034, interestRate: 0.3293, numPayments: 24, expected: 42863496},
			// Zero interest rate
			{loanAmountCents: 1000, interestRate: 0, numPayments: 1, expected: 1000},
			{loanAmountCents: 100000, interestRate: 0, numPayments: 10, expected: 10000},
			{loanAmountCents: 10, interestRate: 0, numPayments: 3, expected: 3}, // TODO: double check this (1 cent has gone unaccounted for)
			{loanAmountCents: 1000, interestRate: 0, numPayments: 3, expected: 333},
			{loanAmountCents: 999, interestRate: 0, numPayments: 3, expected: 333},
		}

		for _, c := range cases {
			actual, err := s.CalcPMT(c.loanAmountCents, c.interestRate, c.numPayments)
			require.NoError(t, err)
			assert.Equal(t, c.expected, actual)
		}
	})

	t.Run("unhappy - zero num payments returns error", func(t *testing.T) {
		c := struct {
			loanAmountCents int64
			interestRate    float64
			numPayments     int32
		}{loanAmountCents: 1000, interestRate: 0.1, numPayments: 0}

		_, err := s.CalcPMT(c.loanAmountCents, c.interestRate, c.numPayments)
		require.ErrorIs(t, err, ErrZeroNumPayments)
	})
}
