package quad

import "github.com/dyedgreen/comp-phys/pkg/casino"

type uniformMonteCarloIntegral monteCaroloIntegral

// NewUniformIntegral is a helper for creating a Monte-Carlo
// integral with uniform sampling function.
func NewUniformMonteCarloIntegral(workers, batch int, seeds []uint64) Integral {
	return &uniformMonteCarloIntegral{
		accuracy: defaultMonteCarloAccuracy,
		steps:    defaultMonteCarloStep,
		workers:  workers,
		batch:    batch,
		seeds:    padSeeds(seeds, workers),
	}
}

// Accuracy implements Integral
func (mont *uniformMonteCarloIntegral) Accuracy(acc *float64) float64 {
	return (*monteCaroloIntegral)(mont).Accuracy(acc)
}

// Steps implements Integral. Note that at least 3 steps
// are always evaluated, no matter what is set here.
func (mont *uniformMonteCarloIntegral) Steps(stp *int) int {
	return (*monteCaroloIntegral)(mont).Steps(stp)
}

// Function implements Integral
func (mont *uniformMonteCarloIntegral) Function(fn func(float64) float64) error {
	return (*monteCaroloIntegral)(mont).Function(fn)
}

func (mont *uniformMonteCarloIntegral) Stats() *Stats {
	return mont.stats
}

// Integrate implements Integral
func (mont *uniformMonteCarloIntegral) Integrate(a, b float64) (float64, error) {
	mont.Distribution = casino.UniDistAB{a, b}
	return (*monteCaroloIntegral)(mont).Integrate(a, b)
}
