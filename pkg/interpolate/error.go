package interpolate

type InterploationError int

const (
	ErrorOutOfBounds InterploationError = iota
	ErrorDimMissmatch
	ErrorBadSpline
)

func (err InterploationError) Error() string {
	switch err {
	case ErrorOutOfBounds:
		return "out of bounds error"
	case ErrorDimMissmatch:
		return "dimensional miss match error"
	case ErrorBadSpline:
		return "spline error"
	default:
		return "error"
	}
}
