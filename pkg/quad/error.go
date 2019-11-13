package quad

type Error int

const (
	ErrorConverge Error = iota
)

func (err Error) Error() string {
	switch err {
	case ErrorConverge:
		return "integral did not converge"
	default:
		return "unknown error"
	}
}
