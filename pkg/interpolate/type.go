package interpolate

import "math"

// Range represents a dataset that can
// be interpolated over.
type Range interface {
	// Return the min and max values x can take
	Bounds() (min, max float64)
	// Test if a given x falls within the bounds
	InBounds(x float64) bool
	// Evaluate the interpolation at point x, if an error is returned
	// the value of y is undefined
	Eval(x float64) (y float64, err error)
}

// PeriodicRange wraps a range by making the
// resulting function periodic.
type PeriodicRange struct {
	Range Range
}

// Bounds implements Range.
func (p PeriodicRange) Bounds() (float64, float64) {
	return math.Inf(-1), math.Inf(1)
}

// InBounds implements Range.
func (p PeriodicRange) InBounds(x float64) bool {
	return true
}

// Eval implements Range.
func (p PeriodicRange) Eval(x float64) (float64, error) {
	min, max := p.Range.Bounds()
	x = min + math.Mod(max-min+math.Mod(x, max-min), max-min)
	return p.Range.Eval(x)
}
