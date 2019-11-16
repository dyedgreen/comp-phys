package casino

import (
	"fmt"
	"math"
	"testing"
)

const eps_exp = 1e-2
const eps_var = 1 // TODO: Get more stable algorithm and bring this down!
const trials = 1000
const experiments = 1000

// Test expectation and variance of known distributions
func TestExpect(t *testing.T) {
	unity := func(x float64) float64 {
		return x
	}

	linDist, err := NewLinearDist(0, 1, 2, 0)
	if err != nil {
		t.Error(err)
	}

	dists := []Distribution{
		UniDistAB{0, 10},
		NormalDist{3, 2},
		linDist,
	}

	exps := []float64{
		5,
		3,
		2.0 / 3.0,
	}
	vars := []float64{
		100.0 / 12.0,
		4,
		0.5 - 4.0/9.0,
	}

	fmt.Println(exps, vars)

	for i := range dists {
		e := Expectation{Distribution: dists[i], Function: unity, Seed: 42}
		stats := e.Refine(trials, experiments)
		if stats.Trials != trials*experiments {
			t.Error("wrong number of trials conducted")
		}
		if math.Abs(stats.Value-exps[i]) > eps_exp {
			t.Error(fmt.Sprintf("%v is not approximately %v", stats.Value, exps[i]))
		}
		if math.Abs(stats.Variance-vars[i]) > eps_var {
			t.Error(fmt.Sprintf("%v is not approximately %v", stats.Variance, vars[i]))
		}
	}
}
