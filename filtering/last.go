package filtering

import (
	"github.com/r-the-thinker/piper"
)

// Last only outputs the last element of the stream if there is one
func Last() piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, lastValueStorage interface{}) (piper.PipeResult, interface{}) {
		lastValue := lastValueStorage.(piper.PipeResult)

		if r.IsValue {
			lastValue = r
		}

		// if the predecessor now closed then close with it
		// we use our most recent value whatever it is. It might not be a value but
		// the successor will know that.
		if r.State == piper.Closed {
			lastValue.State = piper.Closed
			return lastValue, lastValue
		}
		// break and save our lastest value
		return piper.PipeResult{IsValue: false, State: piper.Break}, lastValue
	}, InitialStorage: piper.PipeResult{IsValue: false, State: piper.Open}, EventEmitter: nil}
}
