package pmt

import "math"

func CalcPMT(loanAmountCents int64, annualInterestRate float64, numPayments int32) int64 {
	if annualInterestRate == 0 {
		return loanAmountCents / int64(numPayments)
	}

	r := annualInterestRate
	n := float64(numPayments)
	pv := float64(loanAmountCents)

	pmt := r * pv / (1 - math.Pow(1+r, -n))
	return int64(math.Round(pmt))
}
