package main

import (
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	"github.com/dyedgreen/comp-phys/pkg/interpolate"
	"github.com/dyedgreen/comp-phys/pkg/util"
)

const (
	N     = 30
	start = 0.0
	end   = 20.0
)

func generateData() plotter.XYs {
	data := make(plotter.XYs, N)
	x := start
	for i := 0; i < len(data); i++ {
		data[i] = plotter.XY{x, math.J0(x)}
		x += (end - start) / N
	}
	return data
}

func plotLines(title, file string, lines ...interface{}) error {
	p, err := plot.New()
	if err != nil {
		return err
	}
	p.Title.Text = title
	if err := plotutil.AddScatters(p, "data", lines[0]); err != nil {
		return err
	}
	rng, _ := interpolate.NewSplineRange(util.XYToSlice(lines[0].(plotter.XYs)))
	if err := plotutil.AddLines(p, "lagrange", lines[1], "spline", plotter.NewFunction(func(x float64) (y float64) {
		y, _ = rng.Eval(x)
		return
	})); err != nil {
		return err
	}
	if err := p.Save(8*vg.Inch, 5*vg.Inch, file); err != nil {
		return err
	}
	return nil
}

func shitLagrange(data plotter.XYs, x float64) float64 {
	// a O(N^2) Lagrange implementation
	var y float64
	for i := range data {
		var top float64 = 1
		var bot float64 = 1
		for j := range data {
			if j == i {
				continue
			}
			top *= (x - data[j].X)
			bot *= (data[i].X - data[j].X)
		}
		y += top * data[i].Y / bot
	}
	return y
}

func cubicSpline(data plotter.XYs, x float64) float64 {
	return 0 // TODO
}

func main() {
	// Plot the data
	data := generateData()
	interpolation := make(plotter.XYs, 100)
	x := start
	step := (end - start) / float64(len(interpolation))
	for i := range interpolation {
		interpolation[i] = plotter.XY{x, shitLagrange(data, x)}
		x += step
	}
	plotLines("What the heck is a Bessel function???", "bessel.png", data, interpolation)
}
