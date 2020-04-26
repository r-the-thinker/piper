package transform

import (
	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/filtering"
)

// Reduce uses the supplied reducer to reduce all values that are received to a single output which will
// be emitted once the input channel closes.
// - reducer: (accumulated, currentValue) newValue
// - seedValue: the value used to initialize the accumulator with
func Reduce(reducer func(interface{}, interface{}) interface{}, seedValue interface{}) (piper.PipeOperator, piper.PipeOperator) {
	return Scan(reducer, seedValue), filtering.Last()
}
