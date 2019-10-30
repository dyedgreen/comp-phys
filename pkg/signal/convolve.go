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
func FFTConvolve(a, b []float64) []float64 {
	// Translate b to be stored in wrap-around order
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
	// Compute convolution
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
