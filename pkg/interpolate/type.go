package interpolate

// Range represents a dataset that can
// be interpolated over.
type Range interface {
	// Return the min and max values x can take
	Bounds() (min, max float64)
	// Test if a given x falls within the bounds
	InBounds(x float64) bool
	// Evaluate the interpolation at point x, if an error is returned
	// the value of y is undefined.
	Eval(x float64) (y float64, err error)
}
