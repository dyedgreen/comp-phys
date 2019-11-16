package casino

import (
	"errors"
	"math"

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
		dist = UniDist{}
	}
	return &Sampler{dist, rand.New(rand.NewSource(seed))}
}

func NewUniformSampler(a, b float64, seed uint64) *Sampler {
	return NewSampler(UniDistAB{a, b}, seed)
}

// Returns a random number from the underlying source,
// transformed according the the contained distribution.
func (s *Sampler) Sample() float64 {
	return s.Transform(s.rng.Float64())
}

// Implements a uniform distribution [0,1)
type UniDist struct{}

func (_ UniDist) Transform(x float64) float64 {
	return x
}

func (_ UniDist) Prob(x float64) float64 {
	return 1
}

// Implements a uniform distribution [A, B)
type UniDistAB struct {
	A, B float64
}

func (d UniDistAB) Transform(x float64) float64 {
	return d.A + (d.B-d.A)*x
}

func (d UniDistAB) Prob(x float64) float64 {
	if x < d.A || x > d.B {
		return 0
	}
	return 1 / (d.B - d.A)
}

// TODO: Linear probability, Normal probability
type linearDist struct {
	a, b, alpha, beta, gamma float64
}

// Create a linear distribution gamma * (alpha * x + beta), with
// support [a,b]. (The distribution will be scaled by gamma)
// automatically to ensure normalization.
func NewLinearDist(a, b, alpha, beta float64) (Distribution, error) {
	if a >= b {
		return nil, errors.New("invalid range")
	}

	dist := linearDist{
		a, b, alpha, beta, 1 / (beta*(b-a) + alpha/2*(b*b-a*a)),
	}

	if dist.Prob(a) < 0 || dist.Prob(b) < 0 {
		return nil, errors.New("the linear function can not be negative within [a, b]")
	}
	return dist, nil
}

func (d linearDist) Transform(x float64) float64 {
	p := d.beta / d.alpha
	return -p + math.Sqrt(p*p+d.a*d.a+2*d.a*p+2*d.gamma*x/d.alpha)
}

func (d linearDist) Prob(x float64) float64 {
	if x < d.a || x > d.b {
		return 0
	}
	return d.gamma * (d.alpha*x + d.beta)
}

// Approximately implements a normal distribution
// with given mean and variance. Sigma must be positive.
type NormalDist struct {
	Mu, Sigma float64
}

func (d NormalDist) Transform(x float64) float64 {
	// \sqrt{2} \sigma \erf^{-1}(2x - 1) + \mu = y
	return math.Erfinv(2*x-1)*math.Sqrt2*d.Sigma + d.Mu
}

func (d NormalDist) Prob(x float64) float64 {
	return 1 / (math.SqrtPi * math.Sqrt2 * d.Sigma) * math.Exp(-(x-d.Mu)*(x-d.Mu)/(2*d.Sigma*d.Sigma))
}