package filtering

import (
	"github.com/r-the-thinker/piper"
)

// First only allows one value to pass through
func First() piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// if there is IsValue=false the predecessor closes
		// if there is IsValue=true then we close anyways
		r.State = piper.Closed
		return r, nil
	}, InitialStorage: nil, EventEmitter: nil}
}
