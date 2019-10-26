package main

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/dyedgreen/comp-phys/assignment/comply"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// Generate plot for question a
func q_a(uni distuv.Uniform) {
	fmt.Println("Generating plot for question (a)")

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
	p.Add(h)

	// Overlay PDF
	pdf := plotter.NewFunction(func(x float64) float64 {
		// We scale by factor, since the y-axis is not normalized
		return 1000
	})
	pdf.Color = color.RGBA{255, 0, 0, 255}
	p.Add(pdf)

	if err := p.Save(8*vg.Inch, 5*vg.Inch, "plot-a.pdf"); err != nil {
		panic(err)
	}
}

// Generate plot for question b
func q_b(uni distuv.Uniform) {
	fmt.Println("Generating plot for question (b)")

	y := func(x float64) float64 {
		return 2 * math.Asin(x)
	}

	numbers := make(plotter.Values, 1e5, 1e5)
	for i := range numbers {
		numbers[i] = y(uni.Rand())
	}

	// Plot the generated numbers in a histogram
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	// Plot settings
	p.Title.Text = "1e5 Random Numbers from p(x) = 0.5 * cos(0.5 * x)"

	// Draw histogram
	h, err := plotter.NewHist(numbers, 100)
	if err != nil {
		panic(err)
	}
	p.Add(h)

	// Overlay PDF
	pdf := plotter.NewFunction(func(x float64) float64 {
		// We scale by factor, since the y-axis is not normalized
		return 3.1e3 * 0.5 * math.Cos(0.5*x)
	})
	pdf.Color = color.RGBA{255, 0, 0, 255}
	p.Add(pdf)

	if err := p.Save(8*vg.Inch, 5*vg.Inch, "plot-b.pdf"); err != nil {
		panic(err)
	}
}

// Probability distribution types needed for question c
type logProb func(x float64) float64

func (p logProb) LogProb(x float64) float64 {
	return math.Log(p(x))
}

type randLogProb struct {
	// source of randomness
	src rand.Source
	// Probability distribution
	p func(x float64) float64
	// Map from univariate x to y ~ p
	m func(x float64) float64
}

func (r *randLogProb) LogProb(x float64) float64 {
	return math.Log(r.p(x))
}
func (r *randLogProb) Rand() float64 {
	rng := rand.New(r.src)
	return r.m(rng.Float64())
}

// Generate plot for question c
func q_c(src rand.Source) {
	fmt.Println("Generating plot for question (c)")

	target := logProb(func(x float64) float64 {
		return 2 / math.Pi * math.Cos(0.5*x) * math.Cos(0.5*x)
	})
	proposal := &randLogProb{
		src,
		func(x float64) float64 {
			return 0.5 * math.Cos(0.5*x)
		},
		func(x float64) float64 {
			return 2 * math.Asin(x)
		},
	}
	// This number is carefully calibrated to minimize the difference between
	// target and proposal.
	c := 1.3
	sampler := comply.Rejection{c, target, proposal, src}

	numbers := make(plotter.Values, 1e5, 1e5)
	sampler.Sample(numbers)

	// Plot the generated numbers in a histogram
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	// Plot settings
	p.Title.Text = "1e5 Random Numbers from p(x) = 2 / pi * cos^2(0.5 * x)"

	// Draw histogram
	h, err := plotter.NewHist(numbers, 100)
	if err != nil {
		panic(err)
	}
	p.Add(h)

	// Overlay PDF
	pdf := plotter.NewFunction(func(x float64) float64 {
		// We scale by factor, since the y-axis is not normalized
		return 3.1e3 * target(x)
	})
	pdf.Color = color.RGBA{255, 0, 0, 255}
	p.Add(pdf)

	if err := p.Save(8*vg.Inch, 5*vg.Inch, "plot-c.pdf"); err != nil {
		panic(err)
	}
}

func main() {
	// Initialize a random number generator
	// NOTE: This generator is ultimately sourced by
	//       PCG XSL RR 128/64, see: https://github.com/golang/exp/blob/master/rand/rng.go
	//       and http://www.pcg-random.org/pdf/toms-oneill-pcg-family-v1.02.pdf
	src := rand.NewSource(42)
	uni := distuv.Uniform{0, 1, src}

	q_a(uni)

	// Time these functions
	a := time.Now()
	q_b(uni)
	b := time.Now()
	q_c(src)
	c := time.Now()
	fmt.Println("The relative time taken is:", float64(c.Sub(b))/float64(b.Sub(a)), "(the second routine is slower)")
}
