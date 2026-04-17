package domain_test

import (
	"testing"

	"triangleoflove/backend/internal/domain"
)

func TestCalculateDimensionScore_GivenBothMetricsSet_ThenReturnsNormalizedAverage(t *testing.T) {
	// metrics 4 and 4: avg=4, normalized=(4-1)/4*100=75
	score := domain.CalculateDimensionScore(4, 4)
	if score != 75 {
		t.Errorf("expected 75, got %d", score)
	}
}

func TestCalculateDimensionScore_GivenOneMetricZero_WhenCalculated_ThenScoreUsesNonZeroOnly(t *testing.T) {
	// metric 0 and 3: avg=3, normalized=(3-1)/4*100=50
	score := domain.CalculateDimensionScore(0, 3)
	if score != 50 {
		t.Errorf("expected 50, got %d", score)
	}

	// metric 5 and 0: avg=5, normalized=(5-1)/4*100=100
	score = domain.CalculateDimensionScore(5, 0)
	if score != 100 {
		t.Errorf("expected 100, got %d", score)
	}
}

func TestCalculateDimensionScore_GivenBothMetricsZero_WhenCalculated_ThenDimensionUnavailable(t *testing.T) {
	// metrics 0 and 0: should return -1 (unavailable)
	score := domain.CalculateDimensionScore(0, 0)
	if score != -1 {
		t.Errorf("expected -1, got %d", score)
	}
}

func TestNewDailyInsight_GivenFullCheckin_ThenReturnsThreeScores(t *testing.T) {
	c := domain.Checkin{
		FeltUnderstood:    4,
		MeaningfulSharing: 4,
		CouldCountOnThem:  3,
		EffortForUs:       5,
		Desire:            2,
		Spark:             4,
	}
	insight := domain.NewDailyInsight(c)

	if insight.Intimacy != 75 {
		t.Errorf("expected intimacy=75, got %d", insight.Intimacy)
	}
	if insight.Commitment != 75 {
		t.Errorf("expected commitment=75, got %d", insight.Commitment)
	}
	if insight.Passion != 50 {
		t.Errorf("expected passion=50, got %d", insight.Passion)
	}
}
