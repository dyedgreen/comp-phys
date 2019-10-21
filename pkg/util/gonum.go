package util

import (
	"fmt"
	"strings"

	"github.com/dyedgreen/comp-phys/pkg/interpolate"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

// XYToSlice converts between gonums plotter XYs and float64 slices
func XYToSlice(xys plotter.XYs) (xs, ys []float64) {
	xs = make([]float64, len(xys), len(xys))
	ys = make([]float64, len(xys), len(xys))
	for i, xy := range xys {
		xs[i], ys[i] = xy.X, xy.Y
	}
	return
}

// RangeToPlotter wraps an interpolation range in a gonum plot.Plotter type
func RangeToPlotter(r interpolate.Range) plot.Plotter {
	min, max := r.Bounds()
	return &plotter.Function{
		F: func(x float64) float64 {
			// We can safely ignore the error, since we
			// set the min/max x range in the plotter.Function
			y, _ := r.Eval(x)
			return y
		},
		XMin:    min,
		XMax:    max,
		Samples: 1000,
	}
}

// MatrixToLaTeX takes a matrix and an element-wise format string
// and returns LaTeX code for displaying the matrix as a string.
func MatrixToLaTeX(m mat.Matrix, format string) string {
	str := strings.Builder{}
	str.WriteString("\\begin{bmatrix}")

	// Default format string
	if format == "" {
		format = "%v"
	}

	r, c := m.Dims()
	for i := 0; i < r; i++ {
		if i > 0 {
			str.WriteString(" \\\\ ")
		}
		for j := 0; j < c; j++ {
			if j > 0 {
				str.WriteString(" & ")
			}
			str.WriteString(fmt.Sprintf(format, m.At(i, j)))
		}
	}

	str.WriteString("\\end{bmatrix}")
	return str.String()
}
