package main

import (
	"fmt"
	"time"

	"github.com/dyedgreen/comp-phys/pkg/casino"
	"github.com/dyedgreen/comp-phys/pkg/quad"
)

// Calculate long-running IS data point
func genHeavyData() {
	fmt.Println("Running heavy calculation. This may take a while ...")

	var eps float64 = 1e-6
	var steps int = 2e10

	dist, err := casino.NewLinearDist(A, B, -0.48, 0.98)
	if err != nil {
		panic(err)
	}
	mont := quad.NewMonteCarloIntegral(dist, 128, 128, casino.Noise(64))
	mont.Accuracy(&eps)
	mont.Steps(&steps)

	fmt.Printf("INFO: %v steps max at %v target accuracy\n", steps, eps)

	start := time.Now()
	P, err := quad.Integrate(wave_fn_2, A, B, mont)
	elapsed := time.Now().Sub(start)

	fmt.Printf("Result: P = %v\n", P)
	fmt.Println("Statistics:")
	fmt.Println(mont.Stats())
	fmt.Printf("Time elapsed: %v (%v nanosecond/sample)\n",
		elapsed,
		float64(elapsed.Nanoseconds())/float64(mont.Stats().Steps))
}
