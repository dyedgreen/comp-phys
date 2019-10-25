package comply

import (
	"github.com/dyedgreen/comp-phys/pkg/interpolate"

	"gonum.org/v1/gonum/mat"
)

// We have a separate spline range
// here to comply with using the
// home-cooked LU decomposition
// instead of the much more stable
// and faster inline tridiagonal band
// solver...
type SplineRange struct {
	xs []float64
	ys []float64
	yy []float64
}

// NewSplineRange is a compliant version of interpolate.NewSplineRange
func NewSplineRange(xs, ys []float64) (*SplineRange, error) {
	if len(ys) < 1 || len(ys) != len(xs) {
		return nil, interpolate.ErrorDimMissmatch
	}
	yy := make([]float64, len(ys), len(ys))
	// The compliant version only supports natural
	// boundary conditions
	t := make([]float64, len(yy), len(yy))
	m := mat.NewDense(len(yy), len(yy), nil)
	for i := 1; i < len(t)-1; i++ {
		t[i] = (ys[i+1]-ys[i])/(xs[i+1]-xs[i]) - (ys[i]-ys[i-1])/(xs[i]-xs[i-1])
		m.Set(i, i-1, (xs[i]-xs[i-1])/6)
		m.Set(i, i, (xs[i+1]-xs[i-1])/3)
		m.Set(i, i+1, (xs[i+1]-xs[i])/6)
	}
	// Natural BCs
	m.Set(0, 0, 1)
	m.Set(len(yy)-1, len(yy)-1, 1)
	decomp, err := NewLU(m)
	if err != nil {
		return nil, err
	}
	tVec := mat.NewVecDense(len(t), t)
	yyVec := decomp.Solve(tVec)
	for i := range yy {
		yy[i] = yyVec.AtVec(i)
	}
	// Assemble spline range
	return &SplineRange{xs, ys, yy}, nil
}

// Bounds implements a Range.
func (r *SplineRange) Bounds() (float64, float64) {
	return r.xs[0], r.xs[len(r.xs)-1]
}

// InBounds implements a Range.
func (r *SplineRange) InBounds(x float64) bool {
	min, max := r.Bounds()
	return min <= x && max >= x
}

// Eval implements a Range.
func (r *SplineRange) Eval(x float64) (y float64, err error) {
	if !r.InBounds(x) {
		return 0, interpolate.ErrorOutOfBounds
	}
	// Perform binary search to find points for x
	bot, top := 0, len(r.xs)-1
	mid := (bot + top) / 2
	for bot+1 < top {
		if r.xs[mid] > x {
			top = mid
		} else {
			bot = mid
		}
		mid = (bot + top) / 2
	}
	// Compute spline from spline formula
	h := r.xs[top] - r.xs[bot]
	if h == 0 {
		return 0, interpolate.ErrorBadSpline
	}
	a := (r.xs[top] - x) / h
	b := (x - r.xs[bot]) / h
	return a*r.ys[bot] + b*r.ys[top] + ((a*a*a-a)*r.yy[bot]+(b*b*b-b)*r.yy[top])*(h*h)/6, nil
}
