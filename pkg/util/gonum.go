package util

import (
	"gonum.org/v1/plot/plotter"
)

func XYToSlice(xys plotter.XYs) (xs, ys []float64) {
	xs = make([]float64, len(xys), len(xys))
	ys = make([]float64, len(xys), len(xys))
	for i, xy := range xys {
		xs[i], ys[i] = xy.X, xy.Y
	}
	return
}
