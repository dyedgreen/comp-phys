package quad

import "testing"

func TestTrap(t *testing.T) {
	scheme := NewTrapezoidalIntegral(16)
	helperTestResults(scheme, t)
}

// Ensure step limit and statistic function as
// advertised.
func TestTrapLimit(t *testing.T) {
	scheme := NewTrapezoidalIntegral(16)
	helperTestLimits(scheme, 2, t)
}
