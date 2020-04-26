package transform

import "github.com/r-the-thinker/piper"

// Scan uses a function to accumulate values but emitting the result of every accumulation.
// It works the same as reduce but reduce emits once the stream closed whereas Scan
// emits every new result.
// -accumulator the function that accumulates the values func(previousAccumulated, currentValue)newValue
// -what should the accumulator be initially
func Scan(accumulator func(interface{}, interface{}) interface{}, seed interface{}) piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// if we get a value then accumulate it and update the result
		if r.IsValue {
			storage = accumulator(storage, r.Value)
			r.Value = storage
		}
		// return the (possibly modiefied) previous result and the accumulated value
		return r, storage
	}, InitialStorage: seed, EventEmitter: nil}
}
