package casino

import (
	"sync"

	"golang.org/x/exp/rand"
)

// Expectation can be used to find the
// expectation of a function under a
// given distribution.
type Expectation struct {
	// The distribution used to
	// sample from
	Distribution
	// Function to be averaged
	Function func(float64) float64
	// Seed determines how the
	// experiments random number
	// generators are seeded.
	// This is done by initializing
	// a PCG generator with seed and
	// then using it do obtain seeds
	// for each experiment.
	//
	// This has the effect that for the
	// same seed, each subsequent
	// experiment has the same result
	// if rerun.
	Seed uint64
	rng  rand.Source

	// Expectations of value and value^2
	value1, value2 float64
	// Trials and experiments run to
	// achieve these averages
	trials, experiments int

	lock sync.RWMutex
}

// init helper
// This is used in functions
// which require initialization so that
// the type can be used in it's null
// value.
func (exp *Expectation) init() {
	if exp.rng == nil {
		exp.rng = rand.NewSource(exp.Seed)
	}
}

// Refine will update the expectation
// and variance estimate using
func (exp *Expectation) Refine(trials, experiments int) Result {
	// No touching until we are done!
	exp.lock.Lock()
	defer exp.lock.Unlock()
	exp.init()

	var value1 float64
	var value2 float64

	values := make(chan struct {
		v1 float64
		v2 float64
	})
	wait := sync.WaitGroup{}

	// Run the experiments concurrently
	for i := 0; i < experiments; i++ {
		seed := exp.rng.Uint64()
		go func() {
			sampler := NewSampler(exp, seed)
			defer wait.Done()

			var sum1, sum2 float64
			for i := 0; i < trials; i++ {
				// The key bit, sample a function value
				v := exp.Function(sampler.Sample())
				sum1 += v
				sum2 += v * v
			}
			values <- struct {
				v1 float64
				v2 float64
			}{sum1, sum2}
		}()
	}

	// Close the values channel once
	// all the experiments are concluded
	wait.Add(experiments)
	go func() {
		wait.Wait()
		close(values)
	}()

	// Sum up all trials
	for v := range values {
		value1 += v.v1
		value2 += v.v2
	}

	// Add previous results if present
	// TODO: Should I keep the prev results
	// un-multiplied for efficiency/ less rounding?
	// FIXME: https://en.wikipedia.org/wiki/Algorithms_for_calculating_variance
	//        consider these algorithms for improved numeric stability ...

	if exp.trials > 0 {
		value1 += float64(exp.trials) * exp.value1
		value2 += float64(exp.trials) * exp.value2
	}

	// Compute final results
	exp.trials += trials * experiments
	exp.experiments += experiments
	exp.value1 = value1 / float64(exp.trials)
	exp.value2 = value2 / float64(exp.trials)

	return exp.result()
}

func (exp *Expectation) result() Result {
	return Result{
		Value: exp.value1,
		// This may suffer from numerical instability, consider
		// other algorithms! (check https://en.wikipedia.org/wiki/Algorithms_for_calculating_variance)
		Variance: exp.value2 - exp.value1*exp.value1,
		Stats: Stats{
			Trials:      exp.trials,
			Experiments: exp.experiments,
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
