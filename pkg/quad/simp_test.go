package quad

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"testing"
)

// TODO: These is almost identical to the trap test,
// should re-factor into test helper function ...

func TestSimp(t *testing.T) {
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
			scheme := NewSimpsonIntegral(16)
			num, err := Integrate(cases[i], borders[j], borders[j+1], scheme)
			ana := cases[i+1](borders[j+1]) - cases[i+1](borders[j])
			if err != nil {
				t.Error(err, i/2, j/2)
			} else if math.Abs(ana-num) > scheme.Accuracy(nil) {
				t.Error(fmt.Sprintf("result %v is not approximately %v", num, ana), i/2, j/2)
			}
		}
	}
}

// Ensure step limit and statistic function as
// advertised.
func TestSimpLimit(t *testing.T) {
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

	scheme := NewSimpsonIntegral(16)

	for i := 0; i < 100; i++ {
		mut.Lock()
		called = 0
		mut.Unlock()
		N := 2 + (rand.Int() % 1000)
		scheme.Steps(&N)
		Integrate(fn, 0, 1, scheme)
		if called > N {
			t.Error(fmt.Sprintf("was called %v, should have been called %v", called, N))
		}
		if called != scheme.Stats().Steps {
			t.Error(fmt.Sprintf("was called %v, reported to have been called %v", called, scheme.Stats().Steps))
		}
	}
}