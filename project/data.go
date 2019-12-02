package main

import (
	"fmt"
	"time"

	"github.com/dyedgreen/comp-phys/pkg/casino"
	"github.com/dyedgreen/comp-phys/pkg/quad"
)

// Make data and print output
func genData() {
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

	fmt.Println("\n-- Results (Quadrature Methods) --\n")
	fmt.Printf("We find P = %v (Trapezoidal)\n        P = %v (Simpson)\n", P1, P2)
	fmt.Println("Statistics:")
	fmt.Println(trap.Stats())
	fmt.Println(simp.Stats())

	// Monte Carlo Integration
	accs := []float64{1e-3, 1e-4, 1e-5, 1e-6}

	montFlat := quad.NewUniformMonteCarloIntegral(128, 128, casino.Noise(64))
	dist, err := casino.NewLinearDist(0, 2, -0.48, 0.98)
	if err != nil {
		panic(err)
	}
	montSlanted := quad.NewMonteCarloIntegral(dist, 128, 128, casino.Noise(64))

	fmt.Println("\n-- Monte Carlo Results (Flat) --\n")
	for _, acc := range accs {
		montFlat.Accuracy(&acc)

		start := time.Now()
		P, err := quad.Integrate(wave_fn_2, 0, 2, montFlat)
		elapsed := time.Now().Sub(start)

		fmt.Printf("For accuracy %v:\n        P = %v\n", acc, P)

		fmt.Println("Statistics:")
		fmt.Println(montFlat.Stats())
		fmt.Printf("Time elapsed: %v (%v nanosecond/sample)\n",
			elapsed,
			float64(elapsed.Nanoseconds())/float64(montFlat.Stats().Steps))

		if err != nil {
			break
		}
	}

	fmt.Println("\n-- Monte Carlo Results (Slanted) --\n")
	for _, acc := range accs {
		montSlanted.Accuracy(&acc)

		start := time.Now()
		P, err := quad.Integrate(wave_fn_2, 0, 2, montSlanted)
		elapsed := time.Now().Sub(start)

		fmt.Printf("For accuracy %v:\n        P = %v\n", acc, P)

		fmt.Println("Statistics:")
		fmt.Println(montSlanted.Stats())
		fmt.Printf("Time elapsed: %v (%v nanosecond/sample)\n",
			elapsed,
			float64(elapsed.Nanoseconds())/float64(montSlanted.Stats().Steps))

		if err != nil {
			break
		}
	}

	// APIS (still Monte-Carlo)
	mus, sigmas := casino.APISFamily(casino.NewSampler(casino.UniDistAB{-10, 10}, casino.Seed()), 32)
	apis := casino.APIS{
		Function: func(x float64) float64 {
			// This is the finite support we're integrating over
			if x < 0 || x > 2 {
				return 0
			}
			return 1
		},
		Pi:     wave_fn_2,
		Epochs: 128, Iterations: 32,
		Mus: mus, Sigmas: sigmas,
		Seeds: casino.Noise(32),
	}

	fmt.Println("\n-- Monte Carlo Results (APIS) --\n")
	P, Z := apis.Estimate()
	fmt.Printf("        P = %v\n        Z = %v\n", P, Z)
}