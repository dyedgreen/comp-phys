package comply

import (
	"math"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

// RejectionSample implements rejection
// sampling. Note that this is implemented
// in comply, since Gonum has a robust and
// tested rejection sampling implementation.
// This types interface is modeled on Gonum's
// distuv.Rejection type.
// This implementation is more simplistic than the
// the Gonum one, and panics on any errors it
// encounters. It also requires a random source
// and does not fall back to other sources of
// randomness if Rejection.Src == nil.
type Rejection struct {
	C        float64
	Target   distuv.LogProber
	Proposal distuv.RandLogProber
	Src      rand.Source
}

// Sample samples len(batch) numbers from the target
// distribution. This implements a distuv.Sampler.
func (r *Rejection) Sample(batch []float64) {
	// Validate inputs
	if r.C < 1 {
		panic("rejection constant must be >= 1")
	} else if r.Src == nil {
		panic("a random source is required")
	} else if r.Target == nil || r.Proposal == nil {
		panic("target and proposal functions are required")
	}

	// Fill batch with random numbers
	rng := rand.New(r.Src)
	for i := range batch {
		for {
			// Make proposal and calculate acceptance
			// probability
			prop := r.Proposal.Rand()
			probT := r.Target.LogProb(prop)
			probP := r.Proposal.LogProb(prop)
			probA := math.Exp(probT-probP) / r.C
			// Accept with probA probability
			if probA > 1 {
				panic("probability of acceptance should not exceed one")
			} else if probA > rng.Float64() {
				batch[i] = prop
				break
			}
		}
	}
}
