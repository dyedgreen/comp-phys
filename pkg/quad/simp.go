package quad

import (
	"math"
)

// Implements Integral
type simpsonIntegral trapezoidalIntegral

// Create a new Integral, based on Simpson's rule.
// This will evaluate the integral concurrently.
// The argument specifies how many workers will be used
// to evaluate the function. Passing workers < 1 is
// the same as passing workers = 1.
// If more than one worker is used, integrand functions
// must be thread safe.
func NewSimpsonIntegral(workers int) Integral {
	if workers < 1 {
		workers = 1
	}
	return &simpsonIntegral{
		accuracy: defaultAccuracy,
		steps:    defaultMaxStep,
		workers:  workers,
	}
}

// Accuracy implements Integral
func (simp *simpsonIntegral) Accuracy(acc *float64) float64 {
	return (*trapezoidalIntegral)(simp).Accuracy(acc)
}

// Steps implements Integral. Note that at least 3 steps
// are always evaluated, no matter what is set here.
func (simp *simpsonIntegral) Steps(stp *int) int {
	return (*trapezoidalIntegral)(simp).Steps(stp)
}

// Function implements Integral
func (simp *simpsonIntegral) Function(fn func(float64) float64) error {
	return (*trapezoidalIntegral)(simp).Function(fn)
}

// Integrate implements Integral
func (simp *simpsonIntegral) Integrate(a, b float64) (float64, error) {
	// We don't want people messing with the function / accuracy while we are
	// working hard for them!
	simp.lock.RLock()
	defer simp.lock.RUnlock()

	out := make(chan float64)
	next := make(chan bool)
	defer close(next)

	go trap_stepper(simp.workers, simp.function, a, b, out, next)

	steps := 3

	prevTrap := <-out
	next <- true
	trap := <-out

	integral := trap*4/3 - prevTrap/3
	var prevInt float64

	for n := 2; steps+n < simp.steps || simp.steps < 0; n *= 2 {
		steps += n

		prevInt = integral
		prevTrap = trap
		next <- true // Request next step
		trap = <-out
		integral = trap*4/3 - prevTrap/3

		// Check for convergence, after the first trapezoidal 5 steps
		if n > 1<<5 && math.Abs(integral-prevInt) < simp.accuracy {
			break
		}
	}

	// Record statistics
	simp.stats = &Stats{steps, math.Abs(integral - prevInt), nil}
	if simp.stats.Accuracy > simp.accuracy {
		simp.stats.Error = ErrorConverge
	}

	return integral, simp.stats.Error
}

func (simp *simpsonIntegral) Stats() *Stats {
	return simp.stats
}
