package main

import (
	"fmt"

	"github.com/dyedgreen/comp-phys/assignment/comply"
	"github.com/dyedgreen/comp-phys/pkg/util"

	"gonum.org/v1/gonum/mat"
)

func main() {
	// Matrix given in problem
	A := mat.NewDense(5, 5, []float64{
		3.0, 1.0, 0.0, 0.0, 0.0,
		3.0, 9.0, 4.0, 0.0, 0.0,
		0.0, 9.0, +20, +10, 0.0,
		0.0, 0.0, -22, +31, -25,
		0.0, 0.0, 0.0, -55, +61,
	})

	// Compute decomposition
	luA, err := comply.NewLU(A)
	if err != nil {
		panic(err)
	}

	// Compute determinant det(A) = det(L)*det(U) = 1 * det(U)
	var det float64 = 1
	for i := 0; i < 5; i++ {
		det *= luA.U().At(i, i)
	}

	fmt.Printf("The determinant of A is %v.\n", det)
	fmt.Printf("(Compare this to %v, computed by gonum's `mat.Det()`)\n", mat.Det(A))
	fmt.Println("Here the decompositions:")
	fmt.Println("L:", util.MatrixToLaTeX(luA.L(), "%v"))
	fmt.Println("U:", util.MatrixToLaTeX(luA.U(), "%v"))
}
