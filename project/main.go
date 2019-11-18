package main

import (
	"fmt"
	"math"

	"github.com/dyedgreen/comp-phys/pkg/quad"
)

// Function supplied in problem
func wave_fn_2(z float64) float64 {
	// This is |Psi(z)|^2
	return math.Exp(-z*z) / math.SqrtPi
}

func main() {
	// Quadrature integration
	eps := 1e-6

	trap := quad.NewTrapezoidalIntegral(8)
	trap.Accuracy(&eps)
	simp := quad.NewSimpsonIntegral(8)
	simp.Accuracy(&eps)

	P1, err := quad.Integrate(wave_fn_2, 0, 2, trap)
	if err != nil {
		panic(err)
	}

	P2, err := quad.Integrate(wave_fn_2, 0, 2, simp)
	if err != nil {
		panic(err)
	}

	fmt.Printf("We find P = %v (Trapezoidal)\n        P = %v (Simpson)\n", P1, P2)
	fmt.Println("Statistics (Trapezoidal):")
	fmt.Println(trap.Stats())
	fmt.Println("Statistics (Simpson):")
	fmt.Println(simp.Stats())

	// Monte Carlo Integration
	accs := []float64{1e-3, 1e-4, 1e-5, 1e-6}
	
}
