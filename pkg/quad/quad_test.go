package quad

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"testing"
)

type fnSlice []func(float64) float64

// Unit test helpers

func helperTestResults(scheme Integral, t *testing.T) {
	cases := fnSlice{
		// Contains function -> analytic integral
		func(x float64) float64 { return x * x },
		func(x float64) float64 { return x * x * x / 3 },
		math.Exp,
		math.Exp,
		func(x float64) float64 { return 4*x*x - 2*x + 4 },
		func(x float64) float64 { return 4*x*x*x/3 - x*x + 4*x },
	}

	borders := []float64{
		// Contains a, b pairs sequentially
		0, 10,
		-10, 4,
		-4, 3.5,
	}

	for i := 0; i < len(cases); i += 2 {
		for j := 0; j < len(borders); j += 2 {
			num, err := Integrate(cases[i], borders[j], borders[j+1], scheme)
			ana := cases[i+1](borders[j+1]) - cases[i+1](borders[j])
			if err != nil {
				t.Error(fmt.Sprintf("error: \"%v\" (%v, analytic: %v, stats: %v)", err, num, ana, scheme.Stats()), i/2, j/2)
			} else if math.Abs(ana-num) > scheme.Accuracy(nil) {
				t.Error(fmt.Sprintf("result %v is not approximately %v (stats: %v)", num, ana, scheme.Stats()), i/2, j/2)
			}
		}
	}
}

// Ensure step limit and statistic function as
// advertised.
func helperTestLimits(scheme Integral, min int, t *testing.T) {
	// Function that counts how often it is called
	called := 0
	mut := sync.Mutex{}
	rand.Seed(42)
	fn := func(_ float64) float64 {
		mut.Lock()
		called++
		mut.Unlock()
		return rand.Float64() * 100
	}

	for i := 0; i < 100; i++ {
		mut.Lock()
		called = 0
		mut.Unlock()
		N := min + (rand.Int() % 1000)
		scheme.Steps(&N)
		Integrate(fn, 0, 1, scheme)
		if called > N {
			t.Error(fmt.Sprintf("was called %v, should have been called %v", called, N))
		}
		if called != scheme.Stats().Steps {
			t.Error(fmt.Sprintf("was called %v, reported to have been called %v", called, scheme.Stats().Steps))
		}
	}

	// Test min steps works
	N := min - 1
	if N < 0 {
		N = 0
	}
	scheme.Steps(&N)
	_, err := scheme.Integrate(0, 1)
	if err != ErrorMinSteps {
		t.Error("should error if min steps is too low")
	}
}
