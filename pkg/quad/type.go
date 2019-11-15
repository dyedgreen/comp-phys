package quad

// To be used by all implementations in this package
const defaultAccuracy = 1e-5
const defaultMaxStep = 1e6

// Integral represents an integration
// procedure that can be evaluated
// at different points.
type Integral interface {
	// Accuracy sets the accuracy parameter of the
	// underlying integration procedure and returns
	// the stored accuracy. Passing in nil leaves
	// the accuracy unchanged. This expects
	// the passed accuracy to be bigger than 0.
	Accuracy(*float64) float64

	// Steps determines how many times the function
	// will be evaluated at most, before the integration
	// is aborted. Passing in < 0 means no limit.
	Steps(*int) int

	// Function sets the function to be integrated.
	// Note that this function may be required to be
	// thread safe by some of the implementors.
	Function(func(x float64) float64) error

	// Evaluate the integral between a and b.
	// This fails if a > b. The underlying routine
	// may or may not support +/-infinity as arguments.
	Integrate(a, b float64) (float64, error)

	// Return statistics of last run.
	Stats() *Stats
}

// Stats represent Statistics about the performance
// of the last integration performed.
type Stats struct {
	Steps    int
	Accuracy float64
	Error    error
}

// Integrate fn between a, b using the supplied scheme. If no scheme is
func Integrate(fn func(float64) float64, a, b float64, scheme Integral) (float64, error) {
	if scheme == nil {
		// Simpson is the default scheme
		// Use 1 worker, so that fn does not have to
		// be thread safe.
		scheme = NewSimpsonIntegral(1)
	}
	if err := scheme.Function(fn); err != nil {
		return 0, err
	}
	return scheme.Integrate(a, b)
}
