package domain

import "math"

// DailyInsight holds three dimension scores derived from a Checkin.
// Each field is 0–100 when available, or -1 when unavailable.
type DailyInsight struct {
	Intimacy   int `json:"intimacy"`
	Commitment int `json:"commitment"`
	Passion    int `json:"passion"`
}

// CalculateDimensionScore computes a normalized 0–100 score from two proxy metrics (each 0–5).
// Both 0 → -1 (unavailable). One 0 → normalize the non-zero value. Both set → normalize the average.
func CalculateDimensionScore(a, b int) int {
	if a == 0 && b == 0 {
		return -1
	}
	var avg float64
	if a == 0 {
		avg = float64(b)
	} else if b == 0 {
		avg = float64(a)
	} else {
		avg = float64(a+b) / 2.0
	}
	return int(math.Round((avg - 1) / 4 * 100))
}

// NewDailyInsight computes a DailyInsight from a Checkin.
func NewDailyInsight(c Checkin) DailyInsight {
	return DailyInsight{
		Intimacy:   CalculateDimensionScore(c.FeltUnderstood, c.MeaningfulSharing),
		Commitment: CalculateDimensionScore(c.CouldCountOnThem, c.EffortForUs),
		Passion:    CalculateDimensionScore(c.Desire, c.Spark),
	}
}
