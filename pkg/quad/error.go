package quad

type Error int

const (
	ErrorConverge Error = iota
	ErrorMinSteps
	ErrorInsufficientSteps
)

func (err Error) Error() string {
	switch err {
	case ErrorConverge:
		return "integral did not converge"
	case ErrorMinSteps:
		return "the number of steps allowed is lower than the min needed"
	case ErrorInsufficientSteps:
		return "not enough steps taken to determine convergence"
	default:
		return "unknown error"
	}
}
