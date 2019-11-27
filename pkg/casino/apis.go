package casino

// APIS implements the Adaptive
// Population Importance Sampling
// algorithm as presented in:
//
//    Martino, L.; Elvira, V.; Luengo, D.; Corander, J. (2015-08-01).
//    "An Adaptive Population Importance Sampler: Learning From Uncertainty".
//    IEEE Transactions on Signal Processing. 63 (16): 4422–4437.
//
// The algorithm estimates:
//    I = 1/Z \int f(x) pi(x) dx
// and
//    Z = \int pi(x) dx
// using a set of Gaussian proposal
// functions.
//
// This implementation is not
// parallelized. The memory requirements
// are O(len(Mus)).
type APIS struct {
	// Functions used in estimation
	Function, Pi func(float64) float64

	// Number of epochs and iterations
	// per epoch.
	Epochs, Iterations int

	// Mean and variances for the
	// Gaussian family of distributions.
	Mus, Sigmas []float64

	// Seeds for the samplers
	// used to sample from the
	// distributions.
	Seeds []uint64
}

// Estimate returns an estimate for I and Z based
// on the APIS algorithm.
func (apis *APIS) Estimate() (I, Z float64) {
	if len(apis.Mus) != len(apis.Sigmas) || len(apis.Mus) != len(apis.Seeds) {
		panic("need to provide the same amount of mus, sigmas, and seeds")
	}

	// Initialize samplers / distribution family
	samplers := make([]sampler, len(apis.Mus))
	for i := range samplers {
		// We use raw-samplers and are able to access and modify the
		// Gaussian parameters in place.
		samplers[i] = *NewSampler(NormalDist{apis.Mus[i], apis.Sigmas[i]}, apis.Seeds[i]).(*sampler)
	}

	var L float64
	// as well as I, Z; declared as return values

	z := make([]float64, len(apis.Mus))
	w := make([]float64, len(apis.Mus))

	W := make([]float64, len(apis.Mus))
	eta := make([]float64, len(apis.Mus))

	for epoch := 0; epoch < apis.Epochs; epoch++ {
		for iteration := 0; iteration < apis.Iterations; iteration++ {
			// Importance Sampling part

			// sample z form every member of the family
			// O(N)
			for i := range z {
				z[i] = samplers[i].Sample()
			}

			// compute weights
			// O(N^2)
			// Note: This is a trade-of, could use extra
			// O(N) memory and do this in O(N) time instead.
			for i := range z {
				w[i] = apis.Pi(z[i])
				var norm float64
				for j := range samplers {
					norm += samplers[j].Prob(z[i])
				}
				w[i] *= float64(len(samplers)) / norm
			}

			// Normalize w
			var S float64
			for i := range w {
				S += w[i]
			}

			// Improve estimate for I
			var J float64
			for i := range z {
				J += w[i] / S * apis.Function(z[i])
			}

			I = (L*I + S*J) / (L + S)
			L += S

			// Learning - Collect information to update proposal
			for i := range z {
				rho := apis.Pi(z[i]) / samplers[i].Prob(z[i])
				if rho+W[i] != 0 {
					// This needs to be guarded in case Pi has finite
					// support
					eta[i] = (W[i]*eta[i] + rho*z[i]) / (W[i] + rho)
				}
				W[i] += rho
			}
		}
		// Proposal adaptation
		for i := range samplers {
			// Update samplers
			samplers[i].Distribution = NormalDist{eta[i], apis.Sigmas[i]}
			// Reset memory
			W[i] = 0
			eta[i] = 0
		}
	}
	// Compute Z
	Z = L / float64(len(samplers)) / float64(apis.Epochs*apis.Iterations)
	return
}

// APISFamily will generate a random
// APIS Gaussian family. This is a way
// to conveniently seed your APIS.
func APISFamily(sampler Sampler, size int) ([]float64, []float64) {
	mus, sigmas := make([]float64, size), make([]float64, size)
	for i := range mus {
		mus[i] = sampler.Sample()
		sigmas[i] = sampler.Sample()
		if sigmas[i] < 0 {
			sigmas[i] *= -1
		}
	}
	return mus, sigmas
}
