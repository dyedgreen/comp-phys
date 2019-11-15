package casino

import (
	"golang.org/x/exp/rand"
)

// Transformer takes a uniform
// deviate and returns a
// value distributed according
// to a probability function.
type Transformer interface {
	Transform(float64) float64
}

// Prober calculates the probability
// density.
type Prober interface {
	Prob(float64) float64
}

// Distribution is sufficient to
// sample from a probability distribution
type Distribution interface {
	Transformer
	Prober
}

// Sampler allows to sample from a random distribution
type Sampler struct {
	Distribution
	rng *rand.Rand
}

// Create a new sampler. The underlying randomness is provided
// by PCG. These sources are light-weight and it is reasonable
// to have separate sources for every go routine.
// Note that the provided dist should be safe for concurrent
// use.
func NewSampler(dist Distribution, seed uint64) *Sampler {
	if dist == nil {
		// If no distribution is specified, this is a fancy
		// rand.Rand, with less capabilities.
		dist = UnitDist{}
	}
	return &Sampler{dist, rand.New(rand.NewSource(seed))}
}

func NewUniformSampler(a, b float64, seed uint64) *Sampler {
	return NewSampler(UnitDistAB{a, b}, seed)
}

// Returns a random number from the underlying source,
// transformed according the the contained distribution.
func (s *Sampler) Sample() float64 {
	return s.Transform(s.rng.Float64())
}

// Implements a uniform distribution [0,1)
type UnitDist struct{}

func (_ UnitDist) Transform(x float64) float64 {
	return x
}

func (_ UnitDist) Prob(x float64) float64 {
	return 1
}

// Implements a uniform distribution [A, B)
type UnitDistAB struct {
	A, B float64
}

func (d UnitDistAB) Transform(x float64) float64 {
	return d.A + (d.B-d.A)*x
}

func (d UnitDistAB) Prob(x float64) float64 {
	if x < d.A || x > d.B {
		return 0
	}
	return 1 / (d.B - d.A)
}

// TODO: Linear probability, Normal probability
