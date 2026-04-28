package pmt

import (
	"errors"
	"math"
)

var ErrZeroNumPayments = errors.New("number of payments must not be 0")

func CalcPMT(loanAmountCents int64, annualInterestRate float64, numPayments int32) (int64, error) {
	if numPayments == 0 {
		return 0, ErrZeroNumPayments
	}
	if annualInterestRate == 0 {
		return loanAmountCents / int64(numPayments), nil
	}

	r := annualInterestRate
	n := float64(numPayments)
	pv := float64(loanAmountCents)

	pmt := r * pv / (1 - math.Pow(1+r, -n))
	return int64(math.Round(pmt)), nil
}
