package quad

// Helper function that spawns integral workers and computes
// successive trapezoidal steps. The results are sent to out,
// and the workers are terminated if done is closed.
// Termination happens after every step.
//
// This will compute a least the first step, and at most as
// many steps as are requested by sending true to next, plus one.
// (The first result is reported immediately, further can be requested)
// If next is closed or false is sent, this terminates.
func trap_stepper(workers int, fn func(float64) float64, a, b float64, out chan<- float64, next <-chan bool) {
	// Channels used to gather results
	results := make(chan float64, workers)
	work := make(chan float64, 2) // TODO: What is good for capacity?
	defer close(work)

	// Spawn workers, workers are killed by
	// closing the work channel
	for i := 0; i < workers; i++ {
		go func() {
			for x := range work {
				results <- fn(x)
			}
		}()
	}

	// Gather new points along the function, until done
	// is closed.
	h := b - a
	work <- a
	work <- b
	integral := 0.5 * h * (<-results + <-results)

	for n := 1; true; n *= 2 {
		// Report last result
		out <- integral
		// Only produce next integral step if wanted
		if want, ok := <-next; !want || !ok {
			return
		}

		// Half step distance used and update integral
		h *= 0.5
		integral *= 0.5 // Adjust `previously used' h to be half
		// Fill in the missing evaluations
		stp := (b - a) / float64(n)
		x, s, r := a+0.5*stp, 0, 0 // x, sent, received
		for r < n {
			if s < n {
				select {
				case work <- x:
					s++
					x += stp
				case y := <-results:
					r++
					integral += h * y // inner ones are all factor 1
				}
			} else {
				integral += h * <-results
				r++
			}
		}
	}
}
