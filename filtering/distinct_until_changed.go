package filtering

import (
	"reflect"

	"github.com/r-the-thinker/piper"
)

// DistinctUntilChanged let's items through only if it is different
// from the one before
// [theSame] - func(old, new interface{}) returns true if the newValue is the same as old value
//				set nil if the comparision will be done with reflect.DeepEquals which might not
//				always be the best option depending on the data
func DistinctUntilChanged(theSame func(interface{}, interface{}) bool) piper.PipeOperator {
	if theSame == nil {
		theSame = reflect.DeepEqual
	}

	return piper.PipeOperator{F: func(r piper.PipeResult, lastItem interface{}) (piper.PipeResult, interface{}) {
		// if there was a previous item and it is the same => break
		if lastItem != nil && theSame(lastItem, r.Value) && r.IsValue {
			return piper.PipeResult{IsValue: false, State: r.State}, r.Value
		}
		// otherwise forward the element and save the new value. In case
		// IsValue is false then this means we are closing with the predecessor
		// so it does not matter if we save that zero value.
		return r, r.Value
	}, InitialStorage: nil, EventEmitter: nil}
}
