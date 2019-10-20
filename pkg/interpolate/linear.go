package interpolate

type LinearRange struct {
	xs []float64
	ys []float64
}

// NewLinearRange creates a linearly interpolating range.
// The data passed is not copied, but referenced directly.
func NewLinearRange(xs, ys []float64) (*LinearRange, error) {
	if len(ys) < 1 || len(ys) != len(xs) {
		return nil, ErrorDimMissmatch
	}
	return &LinearRange{xs, ys}, nil
}

// NewLinearRangeCopy is like NewLinearRange, except the passed
// data is copied.
func NewLinearRangeCopy(xs, ys []float64) (*LinearRange, error) {
	xsCopy := make([]float64, len(xs), len(xs))
	ysCopy := make([]float64, len(ys), len(ys))
	copy(xsCopy, xs)
	copy(ysCopy, ys)
	return NewLinearRange(xsCopy, ysCopy)
}

// Bounds implements a Range.
func (r *LinearRange) Bounds() (float64, float64) {
	return r.xs[0], r.xs[len(r.xs)-1]
}

// InBounds implements a Range.
func (r *LinearRange) InBounds(x float64) bool {
	min, max := r.Bounds()
	return min <= x && max >= x
}

// Eval implements a Range.
func (r *LinearRange) Eval(x float64) (y float64, err error) {
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
	return linear(r.xs[bot], r.xs[top], r.ys[bot], r.ys[top], x), nil
}

func linear(x0, x1, y0, y1, x float64) float64 {
	A := (x1 - x) / (x1 - x0)
	return A*y0 + (1-A)*y1
}

// Linear returns the linear interpolation at x between
// two points.
func Linear(x0, x1, y0, y1, x float64) (float64, error) {
	var err error
	if x0 > x || x1 < x {
		err = ErrorOutOfBounds
	}
	return linear(x0, x1, y0, y1, x), err
}
