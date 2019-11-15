package quad

import (
	"fmt"
	"math"
	"testing"
)

type fnSlice []func(float64) float64

func TestTrap(t *testing.T) {
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
			scheme := NewTrapezoidalIntegral(16)
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
