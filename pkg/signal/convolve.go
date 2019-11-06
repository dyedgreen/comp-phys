package signal

import (
	"gonum.org/v1/gonum/fourier"
)

// Convolve computes the convolution of
// a with b, i.e. (a * b).
// Note that this function is NOT normalized
// like a continuous convolution would be. To
// approximate an integral, you need to multiply
// element-wise by your step-size (ie by approx dx).
//
// The input for the response function should be given
// as centered on zero, i.e. b[len(b)/2] = b(0).
// The runtime complexity of this function is O(N^2).
func Convolve(a, b []float64) []float64 {
	c := make([]float64, len(a), len(a))
	off := len(b) / 2
	for i := range c {
		for j := range b {
			idx := j - i + off
			if idx < 0 || idx >= len(a) {
				continue
			}
			c[i] += a[idx] * b[j]
		}
	}
	return c
}

// FFTConvolve is equivalent to Convolve, but is
// implemented using Fourier transforms. I.e.
// this function computes (a * b) = F^-1(F(a) * F(b)),
// where the multiplication is understood to be
// element-wise.
//
// The runtime complexity of this function is O(N*ln(N))
func FFTConvolve(a, b []float64) []float64 {
	// Translate b to be stored in wrap-around order.
	// This is equivalent to assuming that b[len(b)/2] = b(0),
	// as we expect for Convolve.
	bWrap := make([]float64, len(a), len(a))
	off := len(b) / 2
	for i := 0; i < off; i++ {
		bWrap[i] = b[off+i]
		bWrap[off+i] = b[off-i]
	}
	// Compute Fourier transforms
	fft := fourier.NewFFT(len(a))
	c := fft.Coefficients(nil, a)
	d := fft.Coefficients(nil, bWrap)
	// Compute convolution (reusing c for the result)
	for i := range c {
		c[i] = c[i] * d[i]
	}
	// Obtain convolution from inverse transform
	res := fft.Sequence(nil, c)
	// Normalize transform result
	for i := range res {
		res[i] /= float64(len(a))
	}
	return res
}