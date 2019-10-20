package util

import (
	"fmt"
	"testing"
)

func TestApprox(t *testing.T) {
	a := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 235}
	b := []float64{1.01, -2, 3.6, 4.0004, 5.023, 6.0052, -7.35, -8.44, 9.463456, 235.00000346}
	delta := 1e-42
	for i := range a {
		if Approx(a[i], b[i]) {
			t.Error(fmt.Sprintf("%v is not approximately %v", a[i], b[i]))
		}
		if !Approx(a[i], a[i]+delta) {
			t.Error(fmt.Sprintf("%v is approximately %v", a[i], a[i]+delta))
		}
	}
}
