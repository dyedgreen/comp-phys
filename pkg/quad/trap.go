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

// Integrate implements Integral
func (trap *trapezoidalIntegral) Integrate(a, b float64) (float64, error) {
	// We don't want people messing with the function / accuracy while we are
	// working hard for them!
	trap.lock.RLock()
	defer trap.lock.RUnlock()

	// Channels used to gather results
	results := make(chan float64, trap.workers)
	work := make(chan float64, 2) // TODO: What is good for capacity?
	defer close(work)

	for i := 0; i < trap.workers; i++ {
		go func() {
			for x := range work {
				results <- trap.function(x)
			}
		}()
	}

	// Gather new points along the function, until the
	// relative accuracy no longer improves.
	h := b - a
	work <- a
	work <- b
	integral := 0.5 * h * (<-results + <-results)
	prevInt := integral
	steps := 2

	for n := 1; steps+n < trap.steps || trap.steps < 0; n *= 2 {
		steps += n
		prevInt = integral
		// Half step distance used
		h *= 0.5
		integral *= 0.5 // Adjust `previously used' h to be half
		// Fill in the missing evaluations
		stp := (b - a) / float64(n)
		x, s, r := a+0.5*stp, 0, 0 // x, sent, received
		for r < n {
			if s < n {
				select {
				case work <- x:
					s++
					x += stp
				case y := <-results:
					r++
					integral += h * y // inner ones are factor 1
				}
			} else {
				integral += h * <-results
				r++
			}
		}

		// Check for convergence, after the first 5 steps
		if n > 1<<5 && math.Abs(integral-prevInt) < trap.accuracy {
			break
		}
	}

	// Record statistics
	trap.stats = &Stats{steps, math.Abs(integral - prevInt), nil}
	if trap.stats.Accuracy > trap.accuracy {
		trap.stats.Error = ErrorConverge
	}

	return integral, trap.stats.Error
}

func (trap *trapezoidalIntegral) Stats() *Stats {
	return trap.stats
}
