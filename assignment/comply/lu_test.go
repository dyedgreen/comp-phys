package comply

import (
	"fmt"
	"math/rand"
	"testing"

	"gonum.org/v1/gonum/mat"
)

// Tests the LU decomposer by comparing the
// result to LAPACK output
func TestLU(t *testing.T) {
	rand.Seed(42)
	n := 3
	data := make([]float64, n*n)
	for i := 0; i < 1; i++ {
		// Generate random matrix
		for j := range data {
			data[j] = rand.Float64() * 10
		}
		m := mat.NewDense(n, n, data)
		lu, err := NewLU(m)
		if err != nil {
			t.Error(err.Error())
		}
		luLAPACK := mat.LU{}
		luLAPACK.Factorize(m)

		// fix this ...
		//fmt.Println(lu.decomp)
		//fmt.Println(lu.L().At(2, 0), luLAPACK.LTo(mat.NewTriDense(n, mat.Lower, nil)).At(2, 0))
		res := mat.NewDense(n, n, nil)
		res.Mul(lu.L(), lu.U())
		fmt.Println(res)
		fmt.Println(m)

		if !mat.EqualApprox(lu.L(), luLAPACK.LTo(mat.NewTriDense(n, mat.Lower, nil)), 1e-10) {
			t.Error("L decomposition not equal")
		}
		if !mat.EqualApprox(lu.U(), luLAPACK.UTo(mat.NewTriDense(n, mat.Upper, nil)), 1e-10) {
			t.Error("U decomposition not equal")
		}
	}
}
