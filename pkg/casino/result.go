package casino

import "time"

// Holds statistics on the results
// of a complete computation.
type Stats struct {
	// These are totaled over all run
	// experiments
	Burn, Trials int
	// Time taken for computation
	Time time.Time
}

// Contains a single scalar
// result.
type Result struct {
	Value    float64
	Variance float64
	Stats
}
