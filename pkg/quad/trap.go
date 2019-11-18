package quad

import (
	"math"
	"sync"
)

// Implements Integral
type trapezoidalIntegral struct {
	function func(float64) float64
	accuracy float64
	steps    int
	workers  int
	stats    *Stats

	lock sync.RWMutex
}

// Create a new Integral, based on the trapezoidal
// rule. This will evaluate the integral concurrently.
// The argument specifies how many workers will be used
// to evaluate the function. Passing workers < 1 is
// the same as passing workers = 1.
// If more than one worker is used, integrand functions
// must be thread safe.
func NewTrapezoidalIntegral(workers int) Integral {
	if workers < 1 {
		workers = 1
	}
	return &trapezoidalIntegral{
		accuracy: defaultAccuracy,
		steps:    defaultMaxStep,
		workers:  workers,
	}
}

// Accuracy implements Integral
func (trap *trapezoidalIntegral) Accuracy(acc *float64) float64 {
	if acc != nil {
		trap.lock.Lock()
		defer trap.lock.Unlock()
		// Cap at machine epsilon
		trap.accuracy = math.Max(*acc, 1e-16)
	} else {
		// We only need a read lock
		trap.lock.RLock()
		defer trap.lock.RUnlock()
	}
	return trap.accuracy
}

// Steps implements Integral
func (trap *trapezoidalIntegral) Steps(stp *int) int {
	if stp != nil {
		trap.lock.Lock()
		defer trap.lock.Unlock()
		trap.steps = *stp
	} else {
		trap.lock.RLock()
		defer trap.lock.RUnlock()
	}
	return trap.steps
}

// Function implements Integral
func (trap *trapezoidalIntegral) Function(fn func(float64) float64) error {
	trap.lock.Lock()
	defer trap.lock.Unlock()
	trap.function = fn
	return nil
}

func (trap *trapezoidalIntegral) Stats() *Stats {
	return trap.stats
}

// Integrate implements Integral
func (trap *trapezoidalIntegral) Integrate(a, b float64) (float64, error) {
	// We don't want people messing with the function / accuracy while we are
	// working hard for them!
	trap.lock.RLock()
	defer trap.lock.RUnlock()

	if trap.steps < 2 && trap.steps >= 0 {
		trap.stats = &Stats{0, 0, ErrorMinSteps}
		return 0, ErrorMinSteps
	}

	out := make(chan float64)
	next := make(chan bool)
	defer close(next)

	go trap_stepper(trap.workers, trap.function, a, b, out, next)

	steps := 2

	integral := <-out
	var prevInt float64

	var n int
	for n = 1; steps+n < trap.steps || trap.steps < 0; n *= 2 {
		steps += n

		prevInt = integral
		next <- true // Request next integral
		integral = <-out

		// Check for convergence, after the first 5 steps
		if n > 1<<5 && math.Abs(integral-prevInt) < trap.accuracy {
			break
		}
	}

	// Record statistics
	trap.stats = &Stats{steps, math.Abs(integral - prevInt), nil}

	if n <= 1<<5 {
		// We are not confident in the result, unless we take 5 refining steps
		// This number is based on experience and comes from Numerical Recipes
		trap.stats.Error = ErrorInsufficientSteps
	} else if trap.stats.Accuracy > trap.accuracy {
		trap.stats.Error = ErrorConverge
	}

	return integral, trap.stats.Error
}
