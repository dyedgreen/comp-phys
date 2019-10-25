package comply

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/dyedgreen/comp-phys/pkg/interpolate"
	"github.com/dyedgreen/comp-phys/pkg/util"
)

// Simple test for spline-range by checking against
// the implementation in pkg/interpolate
func TestSplineRange(t *testing.T) {
	rand.Seed(42)
	n := 100
	var step float64 = 1
	xs := make([]float64, n, n)
	ys := make([]float64, n, n)
	for i := range xs {
		xs[i] = float64(i) * step
		ys[i] = rand.Float64() * 10
	}

	complyRange, err := NewSplineRange(xs, ys)
	if err != nil {
		t.Fatal(err)
	}
	splineRange, err := interpolate.NewSplineRange(xs, ys)
	if err != nil {
		t.Fatal(err)
	}

	minComply, maxComply := complyRange.Bounds()
	if min, max := splineRange.Bounds(); min != minComply || max != maxComply {
		t.Error("bounds are faulty")
	}

	for i := 1; i < (n-1)*100; i++ {
		x := float64(i) * step / 100
		yComply, err := complyRange.Eval(x)
		if err != nil {
			t.Error(err)
			continue
		}
		ySpline, err := splineRange.Eval(x)
		if err != nil {
			t.Error(err)
			continue
		}
		if !util.Approx(ySpline, yComply) {
			t.Error(fmt.Sprintf("%v is not approximately %v", yComply, ySpline))
		}
	}
}
