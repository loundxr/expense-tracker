package money

import "math"

func ToCents(amount float64) int64 {
	return int64(math.Round(amount * 100))
}

func ToUnits(cents int64) float64 {
	return float64(cents) / 100
}
