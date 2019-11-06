package interpolate

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func approx(a, b float64) bool {
	return math.Abs(a-b) < 1e-10
}

// Test spline interpolation
func TestSpline(t *testing.T) {
	f := func(x float64) float64 {
		return x*x*x - x*x + x + 2
	}
	ff := func(x float64) float64 {
		return 6*x - 2
	}

	// Generate random splines; they should match the cubic we choose
	rand.Seed(42)
	for i := 0; i < 100; i++ {
		start := rand.Float64() * 10
		if rand.Float64() > 0.5 {
			start *= -1
		}
		end := start + rand.Float64()*10
		spline, err := NewSplineRangeCopy([]float64{start, end}, []float64{f(start), f(end)})
		spline.yy = []float64{ff(start), ff(end)}
		if err != nil {
			t.Fatal(err)
		}

		for j := 1; j < 100; j++ {
			x := start + (end-start)/100*float64(j)
			if y, err := spline.Eval(x); err != nil || !approx(f(x), y) {
				t.Error(fmt.Sprintf("%v is not approximately %v (error: %v)", y, f(x), err))
			}
		}
	}
}
