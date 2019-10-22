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

// Tests for automated equation solving
func TestLUSolve(t *testing.T) {
	rand.Seed(42)
	n := 20
	data, x := make([]float64, n*n, n*n), make([]float64, n, n)
	for i := 0; i < 100; i++ {
		// Generate random matrix
		for j := range data {
			data[j] = rand.Float64() * 10
			if rand.Float64() > 0.5 {
				data[j] *= -1
			}
		}
		m := mat.NewDense(n, n, data)

		// Generate random x
		for j := range x {
			x[j] = rand.Float64() * 10
		}

		xVec := mat.NewVecDense(n, x)
		yVec := mat.NewVecDense(n, nil)
		yVec.MulVec(m, xVec)

		// Compute LU decomposition
		lu, err := NewLU(m)
		if err != nil {
			t.Error(err.Error())
			continue
		}

		xSol := lu.Solve(yVec)
		if !mat.EqualApprox(xVec, xSol, 1e-10) {
			t.Error("the solution is wrong")
		}
	}
}

// Tests for matrix inversion
func TestLUInvert(t *testing.T) {
	rand.Seed(42)
	n := 30
	data := make([]float64, n*n, n*n)
	ones := make([]float64, n)
	for i := range ones {
		ones[i] = 1
	}
	unity := mat.NewDiagDense(n, ones)
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

		// Compute inverse
		u := lu.Invert()
		u.(*mat.Dense).Mul(m, u)

		if !mat.EqualApprox(u, unity, 1e-10) {
			t.Error("the inverse is wrong")
		}
	}
}
