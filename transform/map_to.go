package transform

import (
	"github.com/r-the-thinker/piper"
)

// MapTo replaces the received value with a static one
func MapTo(val interface{}) piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		r.Value = val
		return r, nil
	}, InitialStorage: nil, EventEmitter: nil}
}
