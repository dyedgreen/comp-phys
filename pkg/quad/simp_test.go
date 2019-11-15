package quad

import "testing"

func TestSimp(t *testing.T) {
	scheme := NewSimpsonIntegral(16)
	helperTestResults(scheme, t)
}

// Ensure step limit and statistic function as
// advertised.
func TestSimpLimit(t *testing.T) {
	scheme := NewSimpsonIntegral(16)
	helperTestLimits(scheme, 3, t)
}
