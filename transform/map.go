package transform

import piper "github.com/r-the-thinker/piper"

// Map can be used to transform the input into some arbitary output
// It should be noted that the end result that is getting passed to the output channel
// should be the assignable to the type of the output channel.
func Map(mapper func(interface{}) interface{}) piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		nVal := r.Value
		// If we get a value then transform it, it might not be if it is the closing signal
		if r.IsValue {
			nVal = mapper(nVal)
		}

		// return our new value but otherwise use what we have been given
		return piper.PipeResult{Value: nVal, IsValue: r.IsValue, State: r.State}, nil
	}, InitialStorage: nil, EventEmitter: nil}
}

// TODO is this really necessary?
/*


// MapToInt takes a mapper function and does type assertion before calling the mapper function
// If the value that was passed is not of the correct type the mapper wont be called and the
// value won't be emitted. (break)
func MapToInt(mapper func(interface{}) int) piper.PipeOperator {
	return piper.PipeOperator{func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// start by using the same values as the predecessor
		nVal := r.Value

		// if we dont get a value which should only happen when the channel is closed
		// then we dont to anything but mirror the behaviour of the predecessor
		if !r.IsValue {
			return r, nil
		}

		// because it is a value, we can use type assertion to check if it is of the right type
		v, ok := nVal.(int)
		// If it is not of the right type but the predecessor is still open so it was sent on
		// purpose, then we break
		if !ok && r.State != piper.Closed {
			return piper.PipeResult{Value: nVal, IsValue: r.IsValue, State: piper.Break}, nil
		}

		// it's a emittable value, so use it. If the predecessor closed then we should too
		return piper.PipeResult{Value: mapper(v), IsValue: true, State: r.State}, nil
	}, nil, nil}
}

// MapToFloat32 takes a mapper function and does type assertion before calling the mapper function
// If the value that was passed is not of the correct type the mapper wont be called and the
// value won't be emitted. (break)
func MapToFloat32(mapper func(float32) float32) piper.PipeOperator {
	return piper.PipeOperator{func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// start by using the same values as the predecessor
		nVal := r.Value

		// if we dont get a value which should only happen when the channel is closed
		// then we dont to anything but mirror the behaviour of the predecessor
		if !r.IsValue {
			return r, nil
		}

		// because it is a value, we can use type assertion to check if it is of the right type
		v, ok := nVal.(float32)
		// If it is not of the right type but the predecessor is still open so it was sent on
		// purpose, then we break
		if !ok && r.State != piper.Closed {
			return piper.PipeResult{Value: nVal, IsValue: r.IsValue, State: piper.Break}, nil
		}

		// it's a emittable value, so use it. If the predecessor closed then we should too
		return piper.PipeResult{Value: mapper(v), IsValue: true, State: r.State}, nil
	}, nil, nil}
}

// MapFloat64 takes a mapper function and does type assertion before calling the mapper function
// If the value that was passed is not of the correct type the mapper wont be called and the
// value won't be emitted. (break)
func MapToFloat64(mapper func(float64) float64) piper.PipeOperator {
	return piper.PipeOperator{func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// start by using the same values as the predecessor
		nVal := r.Value

		// if we dont get a value which should only happen when the channel is closed
		// then we dont to anything but mirror the behaviour of the predecessor
		if !r.IsValue {
			return r, nil
		}

		// because it is a value, we can use type assertion to check if it is of the right type
		v, ok := nVal.(float64)
		// If it is not of the right type but the predecessor is still open so it was sent on
		// purpose, then we break
		if !ok && r.State != piper.Closed {
			return piper.PipeResult{Value: nVal, IsValue: r.IsValue, State: piper.Break}, nil
		}

		// it's a emittable value, so use it. If the predecessor closed then we should too
		return piper.PipeResult{Value: mapper(v), IsValue: true, State: r.State}, nil
	}, nil, nil}
}

// MapToString takes a mapper function and does type assertion before calling the mapper function
// If the value that was passed is not of the correct type the mapper wont be called and the
// value won't be emitted. (break)
func MapToString(mapper func(string) string) piper.PipeOperator {
	return piper.PipeOperator{func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// start by using the same values as the predecessor
		nVal := r.Value

		// if we dont get a value which should only happen when the channel is closed
		// then we dont to anything but mirror the behaviour of the predecessor
		if !r.IsValue {
			return r, nil
		}

		// because it is a value, we can use type assertion to check if it is of the right type
		v, ok := nVal.(string)
		// If it is not of the right type but the predecessor is still open so it was sent on
		// purpose, then we break
		if !ok && r.State != piper.Closed {
			return piper.PipeResult{Value: nVal, IsValue: r.IsValue, State: piper.Break}, nil
		}

		// it's a emittable value, so use it. If the predecessor closed then we should too
		return piper.PipeResult{Value: mapper(v), IsValue: true, State: r.State}, nil
	}, nil, nil}

}
*/
