package interpolate

type SplineRange struct {
	xs []float64
	ys []float64
	yy []float64
}

// NewSplineRange creates a cubic spline range. The data passed is
// not copied, but referenced directly. The first two values
// passed after xs and ys are used as the initial guess (if not
// given, 0 is used).
func NewSplineRange(xs, ys []float64, boundary ...float64) (*SplineRange, error) {
	if len(ys) < 1 || len(ys) != len(xs) {
		return nil, ErrorDimMissmatch
	}
	yy := make([]float64, len(ys), len(ys))
	u := make([]float64, len(ys)-1, len(ys)-1)
	// Set the left bound
	if len(boundary) > 0 {
		yy[0] = -0.5
		u[0] = (3 / (xs[1] - xs[0])) * ((ys[1] - ys[0]) / (xs[1] - xs[0]) - boundary[0])
	} else {
		// Use natural left bound
		yy[0], u[0] = 0, 0
	}
	// Decompose tridiagonal matrix (see Numerical Recipes)
	for i := 1; i < len(yy)-1; i ++ {
		s := (xs[i] - xs[i-1]) / (xs[i+1] - xs[i-1])
		p := s * yy[i-1] + 2
		yy[i] = (s-1) / p
		u[i] = (ys[i+1] -ys[i]) / (xs[i+1] - xs[i]) - (ys[i] - ys[i-1]) / (xs[i] - xs[i-1])
		u[i] = (6 * u[i] / (xs[i+1] - xs[i-1]) - s * u[i-1]) / p
	}
	// Set the right bound
	var un, yn float64
	if len(boundary) > 1 {
		yn = -0.5
		un = (3 / (xs[len(yy)-1] - xs[len(yy)-2])) * (boundary[1] - (ys[len(yy)-1] - ys[len(yy)-2]) / (xs[len(yy)-1] - xs[len(yy)-2]))
	} else {
		// Use natural left bound
		un, un = 0, 0
	}
	yy[len(yy)-1] = (un - yn * u[len(yy) - 2]) / (yn * yy[len(yy)-2] + 1)
	// Back-substitute tridiagonal matrix
	for i := len(yy)-2; i >= 0; i -- {
		yy[i] = yy[i] * yy[i+1] + u[i]
	}
	// Assemble spline range
	return &SplineRange{xs, ys, yy}, nil
}

// NewSplineRangeCopy is like NewSplineRange, except the passed
// data is copied.
func NewSplineRangeCopy(xs, ys []float64, yy0, yy1 float64) (*SplineRange, error) {
	xsCopy := make([]float64, len(xs), len(xs))
	ysCopy := make([]float64, len(ys), len(ys))
	copy(xsCopy, xs)
	copy(ysCopy, ys)
	return NewSplineRange(xsCopy, ysCopy, yy0, yy1)
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
		return 0, ErrorOutOfBounds
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
		return 0, ErrorBadSpline
	}
	a := (r.xs[top] - x) / h
	b := (x - r.xs[bot]) / h
	return a * r.ys[bot] + b * r.ys[top] + ((a*a*a - a) * r.yy[bot] + (b*b*b - b) * r.yy[top]) * (h*h) / 6, nil
}
