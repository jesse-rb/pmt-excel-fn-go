package pmt

import (
	"errors"
	"math"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

var ErrZeroNumPayments = errors.New("number of payments must not be 0")
var ErrInternal = errors.New("something went wrong")

// CalcPMT takes a loan amount, interest rate, and num payments as input, calcualtes and returens the result of the pmt (excel) function.
// Monetary values that require precision are passed in as the smallest unit of currency (cents) so that we can transfer them as integers
// and perform calculations on whole numbers, this is to reduce the amount of opportunities for floating point rounding errors to occur.
func (s *Service) CalcPMT(loanAmountCents int64, interestRate float64, numPayments int32) (int64, error) {
	if numPayments == 0 {
		return 0, ErrZeroNumPayments
	}
	if interestRate == 0 {
		return loanAmountCents / int64(numPayments), nil
	}

	r := interestRate
	n := float64(numPayments)
	pv := float64(loanAmountCents)

	pmt := r * pv / (1 - math.Pow(1+r, -n))
	return int64(math.Round(pmt)), nil
}
