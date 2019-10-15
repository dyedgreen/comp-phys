package interpolate

import (
	"fmt"
	"testing"
)

func TestLinear(t *testing.T) {
	xs := []float64{0, 1, 0, 1, 0, 1, 1, 2}
	ys := []float64{0, 1, 1, 1, 1, 0, 2, 4}
	xx := []float64{0.5, 0.8, 0.5, 1.5}
	yy := []float64{0.5, 1.0, 0.5, 3.0}

	for i := range xx {
		y, err := Linear(xs[i*2], xs[i*2+1], ys[i*2], ys[i*2+1], xx[i])
		if err != nil {
			t.Fatal(err)
		}
		if y != yy[i] {
			t.Error(fmt.Sprintf("%v is no %v", y, yy[i]))
		}
	}
}
