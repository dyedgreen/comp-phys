package signal

import (
	"fmt"
	"math"
	"testing"
)

// Known convolution result
func g(t float64) float64 {
	return math.Exp(-t*t/2) / math.Sqrt(2*math.Pi)
}
func h(t float64) float64 {
	if 3 <= t && t <= 5 {
		return 4
	}
	return 0
}
func gh(t float64) float64 {
	return -math.Erf(math.Sqrt(2)*(3-t)/2)*2 + math.Erf(math.Sqrt(2)*(5-t)/2)*2
}

// Helper for test cases
func makeSeries() (xs, gs, hs []float64, step float64) {
	const res = 1 << 8
	const min float64 = -10
	const max float64 = +10
	const stp float64 = (max - min) / res
	step = stp

	xs = make([]float64, res)
	gs = make([]float64, res)
	hs = make([]float64, res)

	for i := 0; i < res; i++ {
		xs[i] = min + stp*float64(i)
		gs[i] = g(xs[i])
		hs[i] = h(xs[i])
	}

	return
}

// Test Convolve function
func TestConvolve(t *testing.T) {
	xs, gs, hs, step := makeSeries()
	ghs := Convolve(gs, hs)

	for i := range ghs {
		ghs[i] *= step
	}

	for i := range xs {
		if math.Abs(gh(xs[i])-ghs[i]) >= 1e-1 {
			t.Fatal(fmt.Sprintf("convolution is incorrect, have %v, want %v", ghs[i], gh(xs[i])))
		}
	}
}

// Test FFTConvolve function
func TestFFTConvolve(t *testing.T) {
	xs, gs, hs, step := makeSeries()
	ghs := FFTConvolve(gs, hs)

	for i := range ghs {
		ghs[i] *= step
	}

	for i := range xs {
		if math.Abs(gh(xs[i])-ghs[i]) >= 1e-1 {
			t.Fatal(fmt.Sprintf("convolution is incorrect, have %v, want %v", ghs[i], gh(xs[i])))
		}
	}
}
