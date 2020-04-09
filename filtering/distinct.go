package filtering

import piper "github.com/r-the-thinker/piper"

// Distinct ensures that items are somewhat unique. When you pass 1, 2, 1 as three seperate
// values then the last 1 wont be emitted. This uses the "hashing" of the go map.
func Distinct() piper.PipeOperator {
	// is value known 	=> IsValue = false
	// if value unknown => safe that
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		items := storage.(map[interface{}]bool)
		if r.IsValue {
			if items[r.Value] {
				r.IsValue = false
			} else {
				items[r.Value] = true
				storage = items
			}
		}
		return r, storage
	}, InitialStorage: make(map[interface{}]bool), EventEmitter: nil}
}
