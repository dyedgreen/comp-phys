package casino

import "testing"

// Basic tests to make sure the API
// works as advertised. This does not
// test the quality of the random numbers
// themselves.
func TestSampler(t *testing.T) {
	sampler := NewUniformSampler(0, 1, 42)
	if sampler.Prob(-5) != 0 {
		t.Fail()
	}
	if sampler.Prob(0.5) != 1 {
		t.Fail()
	}
	for i := 0; i < 1000; i++ {
		if num := sampler.Sample(); num < 0 || num > 1 {
			t.Fail()
		}
	}
}

// Note: These distributions are tested, by running
// tests for the expected value calculations. These
// tests are fast to run and allow us to determine
// if the 1st and 2nd order moments are faithfully
// reproduced.
