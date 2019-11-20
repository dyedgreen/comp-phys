package quad

import "testing"

func TestMont(t *testing.T) {
	// We use a uniform integral, as this
	// allows us to easily test different ranges
	scheme := NewUniformMonteCarloIntegral(10, 1000, 42)
	acc := 0.5 // Monte Carlo takes a while to converge
	scheme.Accuracy(&acc)
	helperTestResults(scheme, t)
}

// Ensure step limit and statistic function as
// advertised.
func TestMontLimit(t *testing.T) {
	scheme := NewUniformMonteCarloIntegral(10, 1000, 42)
	helperTestLimits(scheme, 0, t)
}
