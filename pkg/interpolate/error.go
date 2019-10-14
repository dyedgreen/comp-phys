package interpolate

type InterploationError int

const (
	ErrorOutOfBounds InterploationError = iota
	ErrorDimMissmatch
)

func (err InterploationError) Error() string {
	switch err {
	case ErrorOutOfBounds:
		return "out of bounds error"
	case ErrorDimMissmatch:
		return "dimensional miss match"
	default:
		return "unknown error"
	}
}
