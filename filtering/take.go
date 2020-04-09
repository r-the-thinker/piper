package filtering

import piper "github.com/r-the-thinker/piper"

// Take grabs the first >number< of items then stops by not allowing any more values
// to bubble through. The input channel will still be open because we don't close
// channels on the receiving side.
// number - how many values to take n is in range [1, MAX_INT]
func Take(number int) piper.PipeOperator {
	// check the value before continuing
	if number < 1 {
		number = 1
	}
	return piper.PipeOperator{F: func(r piper.PipeResult, numLeftStorage interface{}) (piper.PipeResult, interface{}) {
		numLeft := numLeftStorage.(int)
		state := r.State
		// if we only have one left but the predecessor is still open then this means
		// that this is our last emission
		if numLeft == 1 && r.State == piper.Open {
			state = piper.Closed
		}
		return piper.PipeResult{Value: r.Value, IsValue: r.IsValue, State: state}, numLeft - 1
	}, InitialStorage: number, EventEmitter: nil}
}
