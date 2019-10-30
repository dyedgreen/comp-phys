package main

import (
	"fmt"

	"github.com/dyedgreen/comp-phys/pkg/interpolate"
	"github.com/dyedgreen/comp-phys/pkg/util"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func main() {
	fmt.Println("This program will create plots of linear and natural cubic splines for the data from Q 3.")
	x := []float64{-2.1, -1.45, -1.3, -0.2, 0.1, 0.15, 0.9, 1.1, 1.5, 2.8, 3.8}
	y := []float64{0.012155, 0.122151, 0.184520, 0.960789, 0.990050, 0.977751, 0.422383, 0.298197, 0.105399, 3.936690e-4, 5.355348e-7}

	linRange, err := interpolate.NewLinearRange(x, y)
	if err != nil {
		panic(err)
	}
	// NOTE: We may equally use /assignment/comply.NewSplineRange; but this is faster
	//       and we have demonstrated in unit tests for the problem compliant spline that
	//       outputs are identical.
	//       If still unhappy with this choice, feel free to use the compliant implementation
	//       by changing this file and recompiling the report after running `make pdf`.
	splRange, err := interpolate.NewSplineRange(x, y)
	if err != nil {
		panic(err)
	}

	// Generate plot with data and ranges
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	// Plot settings
	p.Title.Text = "Linear and Spline Interpolation"
	p.Legend.Top = true
	p.Y.Min = -0.03
	p.Y.Max = 1

	// Draw interpolation
	plotutil.AddScatters(p, "Data", util.SliceToXY(x, y))
	plotutil.AddLines(p, "Linear", util.RangeToPlotter(linRange), "Natural Cubic Splines", util.RangeToPlotter(splRange))
	if err := p.Save(8*vg.Inch, 5*vg.Inch, "plot.pdf"); err != nil {
		panic(err)
	}
}
