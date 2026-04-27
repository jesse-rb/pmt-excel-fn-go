package pmt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcPMT(t *testing.T) {
	cases := []struct {
		loanAmountCents int64
		interestRate    float64
		numPayments     int32
		expected        int64
	}{
		{loanAmountCents: 1000, interestRate: 0.1, numPayments: 1, expected: 1100},
		{loanAmountCents: 1000, interestRate: 0.25, numPayments: 2, expected: 694},
		{loanAmountCents: 123456, interestRate: 0.05, numPayments: 1, expected: 129629},
		// Zero interest rate
		{loanAmountCents: 1000, interestRate: 0, numPayments: 1, expected: 1000},
		{loanAmountCents: 100000, interestRate: 0, numPayments: 10, expected: 10000},
		{loanAmountCents: 999, interestRate: 0, numPayments: 3, expected: 333},
		{loanAmountCents: 10, interestRate: 0, numPayments: 3, expected: 3}, // TODO: double check this
	}

	for _, c := range cases {
		actual := CalcPMT(c.loanAmountCents, c.interestRate, c.numPayments)
		assert.Equal(t, c.expected, actual)
	}
}
