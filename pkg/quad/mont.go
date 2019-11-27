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
	seeds          []uint64
	stats          *Stats

	lock sync.RWMutex
}

// Pad seed array
func padSeeds(seeds []uint64, n int) []uint64 {
	for i := 0; len(seeds) < n; i++ {
		seeds = append(seeds, seeds[i]+1)
	}
	return seeds
}

// Returns an integral that is evaluated using a
// Monte-Carlo computation of the form:
//
//     int(f) = exp(f/p) ~ p
//
// i.e. this scheme implements an importance
// sampling algorithm.
//
// workers specifies how many experiments to
// run concurrently, batch specifies how many
// trials to run per experiment/ worker.
// This means, we always run workers*batch
// steps at a time to refine.
//
// You should always provide as many seeds as
// workers. If insufficient seeds are provided,
// the seed slice is padded with extra seeds
// by incrementing the previous seeds. Note
// that this strategy may lead to duplicate
// seeds.
func NewMonteCarloIntegral(dist casino.Distribution, workers, batch int, seeds []uint64) Integral {
	return &monteCaroloIntegral{
		Distribution: dist,
		accuracy:     defaultMonteCarloAccuracy,
		steps:        defaultMonteCarloStep,
		workers:      workers,
		batch:        batch,
		seeds:        padSeeds(seeds, workers),
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
		Seeds: mont.seeds,
	}

	steps := 0
	for steps+mont.batch*mont.workers < mont.steps {
		steps += mont.batch * mont.workers
		res := exp.Refine(mont.batch, mont.workers)
		// sigma on expectation estimate is ~ sqrt(variance_estimate / n)
		// (from central limit theorem)
		// -> we want to be within 2 sigma
		if mont.accuracy >= 2*math.Sqrt(res.Variance/float64(steps)) {
			// We are happy with the results
			break
		}
	}

	// Return final result
	res := exp.Result()
	mont.stats = &Stats{steps, 2 * math.Sqrt(res.Variance/float64(steps)), nil}

	// If we couldn't take any steps, then we have no
	// estimate for anything ...
	if steps == 0 {
		mont.stats.Error = ErrorMinSteps
	} else if mont.stats.Accuracy > mont.accuracy {
		mont.stats.Error = ErrorConverge
	}

	return res.Value, mont.stats.Error
}
