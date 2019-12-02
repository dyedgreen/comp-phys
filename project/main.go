package main

import (
	"flag"
	"math"
)

// Function supplied in problem
func wave_fn_2(z float64) float64 {
	// This is |Psi(z)|^2
	return math.Exp(-z*z) / math.SqrtPi
}

func main() {
	graph := flag.Bool("graph", false, "generate graphs")
	data := flag.Bool("data", false, "print data")
	flag.Parse()

	if !*graph && !*data {
		flag.Usage()
	}
	if *data {
		genData()
	}
	if *graph {
		genGraphs()
	}
}
