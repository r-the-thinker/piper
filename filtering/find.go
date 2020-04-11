package filtering

import "github.com/r-the-thinker/piper"

// Find returns the first item that matches the predicate and then there wont be any more
// value after that.
// predicate - func(val)found. A function that given a value, decides if that value is the one
//				that is searched for
func Find(predicate func(interface{}) bool) piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// If it is found then emit that value and close
		if r.IsValue && predicate(r.Value) {
			return piper.PipeResult{Value: r.Value, IsValue: true, State: piper.Closed}, nil
		}
		// if it not a value (closing) or the predicate said that it is not the value that should be found
		// then break if the pred state is open or close with it
		return piper.PipeResult{IsValue: false, State: r.State}, nil
	}, InitialStorage: nil, EventEmitter: nil}
}
