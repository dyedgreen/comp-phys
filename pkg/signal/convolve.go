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
	// Compute convolution
	for i := range c {
		for j := range b {
			idx := j - i + off
			if idx < 0 || idx >= len(a) {
				// This is equivalent to zero-padding a
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
// element-wise. Note that this function expects
// len(a) == len(b).
//
// The runtime complexity of this function is O(N*ln(N))
func FFTConvolve(a, b []float64) []float64 {
	if len(a) != len(b) {
		panic("a and b must be of equal length")
	}
	// Translate b to be stored in wrap-around order.
	// This is equivalent to assuming that b[len(b)/2] = b(0),
	// as we expect for Convolve.
	res := make([]float64, len(a), len(a)) // we can reuse this for the output
	off := len(b) / 2
	for i := 0; i < off; i++ {
		res[i] = b[off+i]
		res[off+i] = b[off-i]
	}
	// Compute Fourier transforms
	fft := fourier.NewFFT(len(a))
	c := fft.Coefficients(nil, a)
	d := fft.Coefficients(nil, res)
	// Compute convolution (reusing c for the result)
	for i := range c {
		c[i] = c[i] * d[i]
	}
	// Obtain convolution from inverse transform
	fft.Sequence(res, c)
	// Normalize transform result
	for i := range res {
		res[i] /= float64(len(a))
	}
	return res
}
