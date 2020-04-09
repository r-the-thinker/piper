package utility

import piper "github.com/r-the-thinker/piper"

// Tap calls the provided function whenever there is a value comming in
// without affecting aynthing else
func Tap(f func(interface{})) piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// only apply the function if it is a actual value
		if r.IsValue {
			f(r.Value)
		}
		return piper.PipeResult{Value: r.Value, IsValue: r.IsValue, State: r.State, IsEvent: false}, nil
	}, InitialStorage: nil, EventEmitter: nil}
}

// TODO add more datatypes? or is it okay because a sideeffect like this can have anytype
