package main

import (
	"flag"
	"fmt"
	"image/color"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	"github.com/dyedgreen/comp-phys/pkg/casino"
	"github.com/dyedgreen/comp-phys/pkg/quad"
)

var fileType *string = flag.String("format", "png", "set type for saving plotted images")

// make accuracy over time type plots
func plotAcc(title, file string, n, vals, accs []float64, theory func(float64) float64) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	// Plot settings
	p.Title.Text = title
	p.Legend.Top = true
	p.Y.Label.Text = "Integral Value"
	p.X.Label.Text = "Function Evaluations"

	// Plot theoretical bounds
	ployPoints := make(plotter.XYs, 200, 200)
	for i := 0; i < 100; i++ {
		x := n[0] + (n[len(n)-1]-n[0])*float64(i)/99
		ployPoints[i] = plotter.XY{x, vals[len(vals)-1] + theory(x)}
		ployPoints[200-1-i] = plotter.XY{x, vals[len(vals)-1] - theory(x)}
	}
	poly, err := plotter.NewPolygon(ployPoints)
	if err != nil {
		panic(err)
	}
	poly.Color = color.RGBA{255, 240, 190, 255}
	poly.LineStyle.Color = color.RGBA{0, 0, 0, 0}
	p.Add(poly)

	// Plot function with errors + final value
	final := plotter.XYs{plotter.XY{n[0], vals[len(vals)-1]}, plotter.XY{n[len(n)-1], vals[len(vals)-1]}}
	progress := make(plotter.XYs, len(vals))
	errors := make(plotter.YErrors, len(vals))
	for i := range vals {
		progress[i] = plotter.XY{n[i], vals[i]}
		errors[i] = struct{ Low, High float64 }{accs[i], accs[i]}
	}
	points := plotter.YErrorBars{
		XYs:     progress,
		YErrors: errors,
	}

	plotutil.AddLines(p, plotter.XYs{}, "Final Value", final)
	plotutil.AddLinePoints(p, "Integral Value", points)
	plotutil.AddErrorBars(p, points)

	if err := p.Save(8*vg.Inch, 5*vg.Inch, fmt.Sprintf("%v.%v", file, *fileType)); err != nil {
		panic(err)
	}
}

func plotForScheme(title, file string, scheme quad.Integral, theory func(float64) float64) {
	const points = 25
	maxSteps := scheme.Steps(nil)

	n, vals, accs := make([]float64, 0), make([]float64, 0), make([]float64, 0)
	for i := 0; i < points; i++ {
		steps := maxSteps * (i + 1) / points

		scheme.Steps(&steps)
		val, err := quad.Integrate(wave_fn_2, A, B, scheme)
		if err != nil && err != quad.ErrorConverge {
			continue // need more steps
		}
		stats := scheme.Stats()
		if len(n) > 0 && float64(stats.Steps) == n[len(n)-1] {
			continue // already did this many steps
		}
		n = append(n, float64(stats.Steps))
		vals = append(vals, val)
		accs = append(accs, stats.Accuracy)
	}

	plotAcc(title, file, n, vals, accs, theory)
}

// Generate the nice graphs for the project report, this is slow ...
func genGraphs() {
	fmt.Printf("Generating Graphs (format .%v) ...\n", *fileType)

	// Deterministic methods
	fmt.Println("Plots for quadrature methods ...")

	eps := 1e-6
	trapMax := 300
	simpMax := 150

	trap := quad.NewTrapezoidalIntegral(8)
	trap.Accuracy(&eps)
	trap.Steps(&trapMax)
	simp := quad.NewSimpsonIntegral(8)
	simp.Accuracy(&eps)
	simp.Steps(&simpMax)

	plotForScheme("Accuracy Trapezoidal Method", "trap-accuracy", trap, func(n float64) float64 {
		h := (B - A) / n
		return 1e-2 * h * h // O(h^2)
	})
	plotForScheme("Accuracy Simpson's Method", "simp-accuracy", simp, func(n float64) float64 {
		h := (B - A) / n
		return 4e-2 * h * h * h * h // O(h^4)
	})

	// Monte Carlo Methods (this will take a few seconds to run ...)
	fmt.Println("Plots for Monte-Carlo methods ...")

	montFlat := quad.NewUniformMonteCarloIntegral(128, 128, casino.Noise(64))
	dist, err := casino.NewLinearDist(A, 2, -0.48, 0.98)
	if err != nil {
		panic(err)
	}
	montSlanted := quad.NewMonteCarloIntegral(dist, 128, 128, casino.Noise(64))
	montFlat.Accuracy(&eps)
	montSlanted.Accuracy(&eps)

	plotForScheme("Accuracy Uniform Importance Sampling", "mont-flat-accuracy", montFlat, func(n float64) float64 {
		return 0.8 / math.Sqrt(n) // sqrt(var / n)
	})
	plotForScheme("Accuracy Slanted Importance Sampling", "mont-slanted-accuracy", montSlanted, func(n float64) float64 {
		return 0.3 / math.Sqrt(n) // sqrt(var / n)
	})
}
