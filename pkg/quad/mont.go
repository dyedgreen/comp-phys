package quad

import (
	"errors"
	"math"
	"sync"

	"github.com/dyedgreen/comp-phys/pkg/casino"
)

// Concept:
// - init integral type w/ distribution
// - when integrating, check distribution support matches bounds

// This is a bit higher to approx match defaultAccuracy
const defaultMonteCarloAccuracy = 1e-3
const defaultMonteCarloStep = 1e8

type monteCaroloIntegral struct {
	casino.Distribution
	accuracy       float64
	steps          int
	function       func(float64) float64
	workers, batch int
	seed           uint64
	stats          *Stats

	lock sync.RWMutex
}

// Returns an integral that is evaluated using a
// Monte-Carlo computation of the form:
//
//     int(f) = exp(f/p)
//
// workers specifies how many experiments to
// run concurrently, batch specifies how many
// trials to run per experiment/ worker.
// This means, we always run workers*batch
// steps at a time to refine.
func NewMonteCarloIntegral(dist casino.Distribution, workers, batch int, seed uint64) Integral {
	return &monteCaroloIntegral{
		Distribution: dist,
		accuracy:     defaultMonteCarloAccuracy,
		steps:        defaultMonteCarloStep,
		workers:      workers,
		batch:        batch,
		seed:         seed,
	}
}

// Accuracy implements Integral
func (mont *monteCaroloIntegral) Accuracy(acc *float64) float64 {
	if acc != nil {
		mont.lock.Lock()
		defer mont.lock.Unlock()
		// Cap at machine epsilon
		mont.accuracy = math.Max(*acc, 1e-16)
	} else {
		// We only need a read lock
		mont.lock.RLock()
		defer mont.lock.RUnlock()
	}
	return mont.accuracy
}

// Steps implements Integral
func (mont *monteCaroloIntegral) Steps(stp *int) int {
	if stp != nil {
		mont.lock.Lock()
		defer mont.lock.Unlock()
		mont.steps = *stp
	} else {
		mont.lock.RLock()
		defer mont.lock.RUnlock()
	}
	return mont.steps
}

// Function implements Integral
func (mont *monteCaroloIntegral) Function(fn func(float64) float64) error {
	mont.lock.Lock()
	defer mont.lock.Unlock()
	mont.function = fn
	return nil
}

func (mont *monteCaroloIntegral) Stats() *Stats {
	return mont.stats
}

func (mont *monteCaroloIntegral) Integrate(a, b float64) (float64, error) {
	mont.lock.Lock()
	defer mont.lock.Unlock()

	if min, max := mont.Support(); min != a || max != b {
		return 0, errors.New("support of distribution must match bounds")
	}

	exp := casino.Expectation{
		Distribution: mont,
		Function: func(x float64) float64 {
			return mont.function(x) / mont.Prob(x)
		},
		Seed: mont.seed,
	}

	steps := 0
	var prevInt float64
	var integral float64
	for steps+mont.batch*mont.workers < mont.steps {
		steps += mont.batch * mont.workers
		prevInt = integral
		res := exp.Refine(mont.batch, mont.workers)
		integral = res.Value
		if steps >= mont.batch*mont.workers*5 && math.Abs(integral-prevInt) < mont.accuracy*0.1 {
			// We are happy with the results
			break
		}
	}

	// Return final result
	mont.stats = &Stats{steps, math.Abs(integral - prevInt), nil}

	// If we couldn't take any steps, then we have no
	// estimate for anything ...
	if steps == 0 {
		mont.stats.Error = ErrorMinSteps
	} else if steps < mont.batch*mont.workers*5 {
		mont.stats.Error = ErrorInsufficientSteps
	} else if mont.stats.Accuracy > mont.accuracy {
		mont.stats.Error = ErrorConverge
	}

	return integral, mont.stats.Error
}
