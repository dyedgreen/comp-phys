package main

import (
	"fmt"
	"image/color"
	"sync"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	"github.com/dyedgreen/comp-phys/pkg/casino"
	"github.com/dyedgreen/comp-phys/pkg/quad"
)

// Generate graph for stability
func genStabilityGraph() {
	fmt.Println("Plots for stability ...")

	const points = 100
	const maxSteps = 5000000

	mont := quad.NewUniformMonteCarloIntegral(8, 8, casino.Noise(64))

	// this will set to machine precision
	eps := 0.0
	mont.Accuracy(&eps)

	mut := sync.Mutex{}
	n := 0
	f := func(_ float64) float64 {
		mut.Lock()
		defer mut.Unlock()
		if n == 0 {
			n = 1
			return 0
		} else {
			n = 0
			return 1
		}
	}

	vals := make(plotter.XYs, 0)
	for i := 0; i < points; i++ {
		steps := maxSteps * (i + 1) / points

		mont.Steps(&steps)
		val, err := quad.Integrate(f, 0, 1, mont)
		if err != nil && err != quad.ErrorConverge {
			fmt.Println(err.Error())
			continue // need more steps
		}
		vals = append(vals, plotter.XY{float64(mont.Stats().Steps), 0.5 - val})
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	// Plot settings
	p.Title.Text = "Numeric Stability of Monte Carlo Implementation"
	p.Legend.Top = true
	p.Y.Label.Text = "Deviation from Exact Arithmetic"
	p.X.Label.Text = "Function Evaluations"

	// Plot theoretically best precision bounds
	ployPoints := plotter.XYs{plotter.XY{0, 1e-16}, plotter.XY{vals[len(vals)-1].X, 1e-16},
		plotter.XY{vals[len(vals)-1].X, -1e-16}, plotter.XY{0, -1e-16}}
	poly, err := plotter.NewPolygon(ployPoints)
	if err != nil {
		panic(err)
	}
	poly.Color = color.RGBA{255, 240, 190, 255}
	poly.LineStyle.Color = color.RGBA{0, 0, 0, 0}
	p.Add(poly)

	plotutil.AddLines(p, vals)

	if err := p.Save(8*vg.Inch, 5*vg.Inch, fmt.Sprintf("mont-stab.%v", *fileType)); err != nil {
		panic(err)
	}
}
