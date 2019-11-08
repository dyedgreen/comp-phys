package main

import (
	"fmt"
	"math"

	"github.com/dyedgreen/comp-phys/pkg/signal"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// Range over which we plot
const start = -10
const stop = 10

// Functions given in assignment
func g(t float64) float64 {
	return math.Exp(-t*t/2) / math.Sqrt(2*math.Pi)
}

func h(t float64) float64 {
	if 3 <= t && t <= 5 {
		return 4
	}
	return 0
}

// Theoretical result of integral
func theory(t float64) float64 {
	return -math.Erf(math.Sqrt(2)*(3-t)/2)*2 + math.Erf(math.Sqrt(2)*(5-t)/2)*2
}

func main() {
	// Compute convolution
	fmt.Println("Computing convolution using Fourier transforms ...")

	// See assignment PDF for reason these constants
	// are chosen (/pdf/assignment.tex)
	const res = 1 << 10
	const min float64 = -10
	const max float64 = +10
	const stp float64 = (max - min) / res

	gs := make([]float64, res)
	hs := make([]float64, res)

	for i := 0; i < res; i++ {
		x := min + stp*float64(i)
		gs[i] = g(x)
		hs[i] = h(x)
	}

	gh := signal.FFTConvolve(gs, hs)

	// Normalize (we approximate an integral, so multiply by approx dx = stp)
	for i := range gh {
		gh[i] *= stp
	}

	// Plot everything
	fmt.Println("Plotting g, h, and (g * h) ...")

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	ghPlotter := make(plotter.XYs, len(gh))
	for i := range ghPlotter {
		ghPlotter[i].X = min + stp*float64(i)
		ghPlotter[i].Y = gh[i]
	}

	gPlotter := plotter.NewFunction(g)
	gPlotter.Samples = 500
	hPlotter := plotter.NewFunction(h)
	hPlotter.Samples = 500
	tPlotter := plotter.NewFunction(theory)
	tPlotter.Samples = 500

	plotutil.AddLines(p,
		"(g * h)(t)", ghPlotter,
		"g(t)", gPlotter,
		"h(t)", hPlotter,
		"exact result", tPlotter)

	// Plot settings
	p.Title.Text = "Plots of g, h, and (g * h)"
	p.X.Min = start
	p.X.Max = stop
	p.Y.Min = 0
	p.Y.Max = 5
	p.Legend.Top, p.Legend.Left = true, true

	if err := p.Save(8*vg.Inch, 5*vg.Inch, "plot.pdf"); err != nil {
		panic(err)
	}
}
