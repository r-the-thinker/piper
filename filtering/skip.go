package filtering

import "github.com/r-the-thinker/piper"

// Skip skips n items, so it lets items pass through after n items have already passed.
// where the amount that should be skipped is in range [1, MAX_INT)
func Skip(amout int) piper.PipeOperator {
	if amout < 1 {
		amout = 1
	}

	// we do not have to watch out for closing because we always use the predecessor state. This means we also do not have to
	// care for the IsValue flag of the predecessor
	return piper.PipeOperator{F: func(r piper.PipeResult, itemsLeftStorage interface{}) (piper.PipeResult, interface{}) {
		itemsLeft := itemsLeftStorage.(int)
		// if we are still skipping then mark it as no value
		if itemsLeft > 0 {
			r.IsValue = false
			return r, itemsLeft - 1
		}
		return r, itemsLeft
	}, InitialStorage: amout, EventEmitter: nil}
}
