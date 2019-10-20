package util

import "math"

const (
	Epsilon = 1e-40
)

// Approx returns if a and b are
// approximately equal.
func Approx(a, b float64, args ...float64) bool {
	// Determine accuracy to use
	var eps float64 = Epsilon
	if len(args) > 0 {
		if len(args) != 1 {
			panic("invalid number of optional arguments")
		}
		eps = args[0]
	}
	// Determine if either number is zero
	if a == 0 && b == 0 {
		return true
	} else if a == 0 {
		return math.Abs(b) < eps
	} else if b == 0 {
		return math.Abs(a) < eps
	}
	// Check sign is equal if |a| > eps
	if math.Abs(a) > eps && math.Signbit(a) != math.Signbit(b) {
		return false
	}
	// Otherwise check ||a/b|-1| < eps
	return math.Abs(math.Abs(a/b)-1) < eps
}
