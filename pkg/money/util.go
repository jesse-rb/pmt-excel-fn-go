package money

func ToCents(v float64) int64 {
	return int64(v * 100)
}

func FromCents(v int64) float64 {
	return float64(v) / 100
}
