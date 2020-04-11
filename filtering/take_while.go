package filtering

import "github.com/r-the-thinker/piper"

// TakeWhile allows value to come through until the predicate returned false once!
// After that the predicate wont be called anymore and the pipefunc wont allow any more
// emissions. As always this does not mean that the input channel will be closed.
// predicate - function that evaluates the value and returns false if the value and all
//				following should not be emitted anymore
func TakeWhile(predicate func(interface{}) bool) piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// If we get a value but the predicate decided that we should stop taking more
		// then close the pipefunc. Otherwise we do as the predecessor said.
		if r.IsValue && !predicate(r.Value) {
			r.State = piper.Closed
		}
		return r, nil
	}, InitialStorage: nil, EventEmitter: nil}
}
