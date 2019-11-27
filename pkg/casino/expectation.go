package casino

import (
	"sync"
)

// Expectation can be used to find the
// expectation of a function under a
// given distribution by repeatedly
// sampling from a provided distribution.
type Expectation struct {
	// The distribution used to
	// sample from
	Distribution
	// Function to be averaged
	Function func(float64) float64
	// Seed determines how the
	// workers random number
	// generators are seeded.
	//
	// If the number of seeds
	// exceeds the number of
	// desired workers, then
	// refine will panic.
	Seeds    []uint64
	samplers []Sampler

	// Expectation value (x_bar) and variance (var(x) = m2 / n)
	x_bar, m2 float64
	// total number of trials
	trials int

	lock sync.RWMutex
}

// init helper
// This is used in functions
// which require initialization so that
// the type can be used in it's null
// value.
func (exp *Expectation) init() {
	if exp.samplers == nil {
		exp.samplers = make([]Sampler, len(exp.Seeds), len(exp.Seeds))
		for i := range exp.samplers {
			exp.samplers[i] = NewSampler(exp, exp.Seeds[i])
		}
	}
}

// Refine will update the expectation
// and variance estimate. This will update
// the estimate of expectation and variance
// by evaluating the Function trials times in
// every worker.
// This calls Function trials * workers times.
func (exp *Expectation) Refine(trials, workers int) Result {
	// No touching until we are done!
	exp.lock.Lock()
	defer exp.lock.Unlock()
	exp.init()

	// Raise panic if insufficient seeds provided
	if workers > len(exp.samplers) {
		panic("insufficient seeds provided to run workers")
	}

	type valStruct struct {
		x_bar, m2 float64
	}

	values := make(chan valStruct)
	wait := sync.WaitGroup{}

	// Close the values channel once
	// all the experiments are concluded
	wait.Add(workers)
	go func() {
		wait.Wait()
		close(values)
	}()

	// Update the expectation concurrently
	for i := 0; i < workers; i++ {
		// each worker uses a separate sampler
		go func(sampler int) {
			defer wait.Done()

			var x_bar_prev float64
			var x_bar float64
			var m2 float64

			for n := 1; n <= trials; n++ {
				// Get a function value, potentially expensive
				x := exp.Function(exp.samplers[sampler].Sample())
				// Online expectation and variance
				// updates based on:
				//
				//     Welford, B. P. (1962). "Note on a method for calculating corrected sums of
				//     squares and products". Technometrics. 4 (3): 419–420.
				//
				x_bar_prev = x_bar
				x_bar = x_bar + (x-x_bar)/float64(n)
				m2 = m2 + (x-x_bar_prev)*(x-x_bar)
			}

			values <- valStruct{x_bar, m2}
		}(i)
	}

	// Combine calculated expectations based on:
	//
	//     Chan, Tony F.; Golub, Gene H.; LeVeque, Randall J. (1979), "Updating Formulae
	//     and a Pairwise Algorithm for Computing Sample Variances.", Technical Report
	//     STAN-CS-79-773, Department of Computer Science, Stanford University.
	//
	for v := range values {
		delta := v.x_bar - exp.x_bar
		exp.x_bar += delta * float64(trials) / float64(exp.trials+trials)
		exp.m2 += v.m2 + delta*delta*float64(exp.trials)*float64(trials)/float64(exp.trials+trials)
		exp.trials += trials
	}

	return exp.result()
}

func (exp *Expectation) result() Result {
	return Result{
		Value: exp.x_bar,
		// Use unbiased variance estimator
		Variance: exp.m2 / float64(exp.trials-1),
		Stats: Stats{
			Trials: exp.trials,
		},
	}
}

// Result returns the current result
// of the computation.
func (exp *Expectation) Result() Result {
	exp.lock.RLock()
	defer exp.lock.RUnlock()
	return exp.result()
}
