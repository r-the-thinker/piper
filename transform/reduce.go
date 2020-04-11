package transform

import "github.com/r-the-thinker/piper"

// Reduce uses the supplied reducer to reduce all values that are received to a single output which will
// be emitted once the input channel closes.
// - reducer: (accumulator, currentValue) newValue
// - seedValue: the value used to initialize the accumulator with
func Reduce(reducer func(interface{}, interface{}) interface{}, seedValue interface{}) piper.PipeOperator {
	// called with the reducer directly because there is not need for type assertion
	return piper.PipeOperator{F: reducing(reducer), InitialStorage: seedValue, EventEmitter: nil}
}

// ReduceToInt uses the supplied reducer to reduce all values that are received to a single output of type int
// which will be emitted once the input channel closes.
// - reducer: (accumulator, currentValue) newValue
// - seedValue: the value used to initialize the accumulator with
func ReduceToInt(reducer func(int, interface{}) int, seedValue int) piper.PipeOperator {
	update := func(storage, val interface{}) interface{} {
		return reducer(storage.(int), val)
	}
	return piper.PipeOperator{F: reducing(update), InitialStorage: seedValue, EventEmitter: nil}
}

// ReduceToFloat32 uses the supplied reducer to reduce all values that are received to a single output of type Float32
// which will be emitted once the input channel closes.
// - reducer: (accumulator, currentValue) newValue
// - seedValue: the value used to initialize the accumulator with
func ReduceToFloat32(reducer func(float32, interface{}) float32, seedValue float32) piper.PipeOperator {
	update := func(storage, val interface{}) interface{} {
		return reducer(storage.(float32), val)
	}
	return piper.PipeOperator{F: reducing(update), InitialStorage: seedValue, EventEmitter: nil}
}

// ReduceToFloat64 uses the supplied reducer to reduce all values that are received to a single output of type Float64
// which will be emitted once the input channel closes.
// - reducer: (accumulator, currentValue) newValue
// - seedValue: the value used to initialize the accumulator with
func ReduceToFloat64(reducer func(float64, interface{}) float64, seedValue float64) piper.PipeOperator {
	update := func(storage, val interface{}) interface{} {
		return reducer(storage.(float64), val)
	}
	return piper.PipeOperator{F: reducing(update), InitialStorage: seedValue, EventEmitter: nil}
}

// ReduceToString uses the supplied reducer to reduce all values that are received to a single output of type string
// which will be emitted once the input channel closes.
// - reducer: (accumulator, currentValue) newValue
// - seedValue: the value used to initialize the accumulator with
func ReduceToString(reducer func(string, interface{}) string, seedValue string) piper.PipeOperator {
	update := func(storage, val interface{}) interface{} {
		return reducer(storage.(string), val)
	}
	return piper.PipeOperator{F: reducing(update), InitialStorage: seedValue, EventEmitter: nil}
}

// Performs the reducing operation but consults the update function for the type assertion.
func reducing(update func(interface{}, interface{}) interface{}) piper.Pipefunc {
	return func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// If we get a value we have to process it by updating the storage
		if r.IsValue {
			storage = update(storage, r.Value)
		}

		// If the predecessor is open, then we break the chain
		if r.State != piper.Closed {
			return piper.PipeResult{IsValue: false, State: piper.Break}, storage
		}
		// if the predecessor closed, then we emit our storage and close outself
		return piper.PipeResult{Value: storage, IsValue: true, IsEvent: false, State: piper.Closed}, nil
	}
}
