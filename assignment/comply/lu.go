package comply

import (
	"errors"

	"gonum.org/v1/gonum/mat"
)

// We define out own LU decomposition here and not as
// a reusable package in pkg, since there is already
// very good LU decomposition functionality in gonum
// (based on LAPACK and BLAS), which is preferable
// to this home-cooked solution.
type LU struct {
	decomp *mat.Dense
}

// This type is used to return a mat.Matrix which represents
// the lower/upper (L/U) part of the decomposition.
type luDecompMatrix struct {
	lu    *LU
	lower bool
}

// NewLU decomposes a given matrix into lower and
// upper triangular matrix.
func NewLU(m mat.Matrix) (*LU, error) {
	if c, r := m.Dims(); c != r {
		return nil, errors.New("invalid dimensions")
	}
	N, _ := m.Dims()
	decomp := mat.NewDense(N, N, nil)
	lu := &LU{decomp}
	L, U := lu.L(), lu.U()

	// Perform a simple application of
	// Crout’s algorithm, without pivoting
	for j := 0; j < N; j++ {
		for i := 0; i <= j; i++ {
			var sum float64
			for k := 0; k < i; k++ {
				sum += L.At(i, k) * U.At(k, j)
			}
			decomp.Set(i, j, m.At(i, j)-sum)
		}
		for i := j + 1; i < N; i++ {
			var sum float64
			for k := 0; k < j; k++ {
				sum += L.At(i, k) * U.At(k, j)
			}
			if U.At(j, j) == 0 {
				return nil, errors.New("matrix is singular")
			}
			decomp.Set(i, j, (m.At(i, j)-sum)/U.At(j, j))
		}
	}

	return lu, nil
}

func (lu *LU) Solve(y mat.Vector) mat.Vector {
	n, _ := lu.decomp.Dims()
	if n != y.Len() {
		panic(errors.New("invalid dimensions for y"))
	}

	// We will solve LUx=y <=> L(Ux) = y, where Ux = x' by
	// first solving for x', and then for x.
	L := lu.L()
	U := lu.U()

	// Solve for x' using back-substitution on L
	x := make([]float64, n, n)
	for i := range x {
		var sum float64
		for j := 0; j < i; j++ {
			sum += L.At(i, j) * x[j]
		}
		x[i] = (y.AtVec(i) - sum) / L.At(i, i)
	}

	// Solve for x using back-substitution on U
	for i := n - 1; i >= 0; i-- {
		var sum float64
		for j := n - 1; j > i; j-- {
			sum += U.At(i, j) * x[j]
		}
		x[i] = (x[i] - sum) / U.At(i, i)
	}

	return mat.NewVecDense(n, x)
}

// L returns the lower triangular decomposition
// matrix.
func (lu *LU) L() mat.Matrix {
	return &luDecompMatrix{lu, true}
}

// U returns the upper triangular decomposition
// matrix.
func (lu *LU) U() mat.Matrix {
	return &luDecompMatrix{lu, false}
}

// At implements mat.Matrix
func (lu *luDecompMatrix) At(i, j int) float64 {
	// Mask underlying matrix
	if lu.lower && i == j {
		// Lower decomposition has L_ii = 1
		return 1
	} else if (lu.lower && i < j) || (!lu.lower && j < i) {
		return 0
	}
	return lu.lu.decomp.At(i, j)
}

// Dims implements mat.Matrix
func (lu *luDecompMatrix) Dims() (int, int) {
	return lu.lu.decomp.Dims()
}

// T implements mat.Matrix
func (lu *luDecompMatrix) T() mat.Matrix {
	return mat.Transpose{lu}
}
