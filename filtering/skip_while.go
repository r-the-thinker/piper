package filtering

import (
	"github.com/r-the-thinker/piper"
)

// SkipWhile skips all items until the predicate returns false
// as in "wanna skip?"
func SkipWhile(skipItFunc func(interface{}) bool) piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// if we dont skip anymore then forward
		if !storage.(bool) {
			return r, false
		}

		// check if this item should be skipped
		skipNow := true
		if r.IsValue {
			skipNow = skipItFunc(r.Value)
		}
		r.IsValue = !skipNow
		return r, skipNow
	}, InitialStorage: true, EventEmitter: nil}
}
