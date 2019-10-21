package comply

import (
	"math/rand"
	"testing"

	"gonum.org/v1/gonum/mat"
)

// Tests the LU decomposer
func TestLU(t *testing.T) {
	rand.Seed(42)
	n := 20
	data := make([]float64, n*n)
	for i := 0; i < 100; i++ {
		// Generate random matrix
		for j := range data {
			data[j] = rand.Float64() * 10
			if rand.Float64() > 0.5 {
				data[j] *= -1
			}
		}
		m := mat.NewDense(n, n, data)

		// Compute LU decomposition
		lu, err := NewLU(m)
		if err != nil {
			t.Error(err.Error())
			continue
		}

		// Reconstruct original matrix
		recon := &mat.Dense{}
		recon.Mul(lu.L(), lu.U())
		if !mat.EqualApprox(m, recon, 1e-10) {
			t.Error("LU != m after decomposition")
		}
	}
}
