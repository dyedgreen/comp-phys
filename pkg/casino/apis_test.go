package casino

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

// Basic tests to make sure APIS is
// implemented correctly.
func TestAPIS(t *testing.T) {
	mus := make([]float64, 10)
	sigs := make([]float64, 10)
	for i := range mus {
		mus[i] = rand.Float64() * 10
		if rand.Float64() > 0.5 {
			mus[i] *= -1
		}
		sigs[i] = rand.Float64() * 10
	}

	// Can a Gaussian estimate a Gaussian?
	norm := NormalDist{0, 1}
	apis := APIS{
		func(x float64) float64 { return x * x }, norm.Prob,
		64, 32,
		mus, sigs,
		Noise(10),
	}

	I, Z := apis.Estimate()
	if math.Abs(I-1) > 1e-2 {
		t.Error(fmt.Sprintf("I = %v should be 1", I))
	}
	if math.Abs(Z-1) > 1e-2 {
		t.Error(fmt.Sprintf("Z = %v should be 1", Z))
	}

	// A slightly more interesting problem
	apis.Function = func(x float64) float64 {
		if x > 0 {
			return 1
		}
		return 0
	}

	I, Z = apis.Estimate()
	if math.Abs(I-0.5) > 1e-2 {
		t.Error(fmt.Sprintf("I = %v should be 0.5", I))
	}
	if math.Abs(Z-1) > 1e-2 {
		t.Error(fmt.Sprintf("Z = %v should be 1", Z))
	}

	// An even more interesting problem
	lambda := 2.345
	apis.Function = func(x float64) float64 {
		return x
	}
	// Un-normalized exponential distribution
	apis.Pi = func(x float64) float64 {
		if x < 0 {
			return 0
		}
		return math.Exp(-lambda * x)
	}

	I, Z = apis.Estimate()
	if math.Abs(I-1/lambda) > 1e-1 {
		t.Error(fmt.Sprintf("I = %v should be %v", I, 1/lambda))
	}
	if math.Abs(Z-1/lambda) > 1e-1 {
		t.Error(fmt.Sprintf("Z = %v should be %v", I, 1/lambda))
	}
}
