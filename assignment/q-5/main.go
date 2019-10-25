package main

import (
	"fmt"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	fmt.Println("Question a)")

	// Initialize a random number generator
	// NOTE: This generator is ultimately sourced by
	//       PCG XSL RR 128/64, see: https://github.com/golang/exp/blob/master/rand/rng.go
	//       and http://www.pcg-random.org/pdf/toms-oneill-pcg-family-v1.02.pdf
	src := rand.NewSource(42)
	uni := distuv.Uniform{0, 1, src}

	numbers := make(plotter.Values, 1e5, 1e5)
	for i := range numbers {
		numbers[i] = uni.Rand()
	}

	// Plot the generated numbers in a histogram
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	// Plot settings
	p.Title.Text = "1e5 Uniform Random Numbers from [0,1]"

	// Draw histogram
	h, err := plotter.NewHist(numbers, 100)
	if err != nil {
		panic(err)
	}
	h.Normalize(1)
	p.Add(h)
	if err := p.Save(8*vg.Inch, 5*vg.Inch, "plot.pdf"); err != nil {
		panic(err)
	}
}
