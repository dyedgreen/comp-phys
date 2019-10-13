package main

import (
	"fmt"
	"strings"
)

// We use our own matrix type
// later would typically use mat.Dense
// from gonum (note: this is not a
// very efficient way to store matrices,
// unless you specifically need fast row-
// reordering)
type Matrix [][]float64

func (A Matrix) String() string {
	out := strings.Builder{}
	for i := range A {
		for j := range A[i] {
			if j > 0 {
				out.WriteString(" ")
			}
			out.WriteString(fmt.Sprintf("%4.2e", A[i][j]))
		}
		out.WriteString("\n")
	}
	return out.String()
}

func (A Matrix) Mult(B Matrix) Matrix {
	// Allocate resulting matrix
	na, ma := len(A), len(A[0])
	nb, mb := len(B), len(B[0])
	if ma != nb {
		panic("matrix dimensions miss-matched")
	}
	C := make(Matrix, na)
	for i := range C {
		C[i] = make([]float64, mb)
	}

	// Do the contraction
	for i := range C {
		for j := range C[i] {
			for k := 0; k < ma; k++ {
				C[i][j] += A[i][k] * B[k][j]
			}
		}
	}
	return C
}

func (A Matrix) Solve(b Matrix) Matrix {
	n := len(A)
	if n != len(A[0]) || n != len(b) || len(b[0]) != 1 {
		panic("matrix dimensions miss-matched")
	}
	cp := make(Matrix, n)
	for i := range cp {
		cp[i] = make([]float64, n+1)
		copy(cp[i], A[i])
		cp[i][n] = b[i][0]
	}
	// TODO: rearrange to have no zeros on diagonal
	// Gauss elimination
	for r := 0; r < n-1; r++ {
		frr := cp[r][r]
		for rr := r + 1; rr < n; rr++ {
			fr := cp[rr][r]
			for c := r; c < n+1; c++ {
				cp[rr][c] *= frr
				cp[rr][c] -= cp[r][c] * fr
			}
		}
	}
	// Generate result
	res := make(Matrix, n)
	for r := n - 1; r >= 0; r-- {
		if cp[r][r] == 0 {
			if cp[r][n] != 0 {
				panic("no solution")
			}
			res[r] = []float64{1}
			continue
		}
		res[r] = []float64{cp[r][n]}
		for c := n - 1; c > r; c-- {
			res[r][0] -= cp[r][c] * res[c][0]
		}
		res[r][0] /= cp[r][r]
	}
	return res
}

func main() {
	fmt.Println("Matrix Multiplication:")

	// Create Matrices
	matA := Matrix{
		[]float64{1, 2, 3},
		[]float64{2, 1, 1},
		[]float64{7, 8, 5},
	}
	matB := Matrix{
		[]float64{1, 0, 1},
		[]float64{0, 1, 0},
		[]float64{1, 0, 1},
	}

	fmt.Println(matA.String())
	fmt.Println(matB.String())

	fmt.Println(matA.Mult(matB).String())

	fmt.Println("Solving Equations:")

	x := Matrix{[]float64{6}, []float64{3}, []float64{1}}
	b := matA.Mult(x)

	fmt.Println(matA.String())
	fmt.Println(x.String())
	fmt.Println(b.String())
	fmt.Println(matA.Solve(b).String())
}
