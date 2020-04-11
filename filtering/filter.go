package filtering

import "github.com/r-the-thinker/piper"

// Filter can be used to filter out values.
// The filterer should return true when the value should be kept and false if not.
func Filter(filterer func(interface{}) bool) piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// otherwise return the value but only if it passes the filterer
		return piper.PipeResult{
			Value:   r.Value,
			IsValue: r.IsValue && filterer(r.Value),
			State:   r.State,
			IsEvent: false,
		}, nil
	}, InitialStorage: nil, EventEmitter: nil}
}

// FilterInt can be used to filter out values.
// The filterer should return true when the value should be kept and false if not.
// If the value passed is not of the correct type the filterer wont be called and the
// value won't be emitted (break)
func FilterInt(filterer func(int) bool) piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// Check if the value can be used
		var val int
		var save bool = r.IsValue
		if save {
			val, save = r.Value.(int)
		}

		// Use the state of the predecessor. But if the value is not save to use and we do not
		// get the closed signal then we break out of the chain.
		state := r.State
		if !save && state != piper.Closed {
			state = piper.Break
		}

		return piper.PipeResult{
			Value:   val,
			IsValue: save && filterer(val),
			State:   state,
			IsEvent: false,
		}, nil
	}, InitialStorage: nil, EventEmitter: nil}
}

// FilterInt8 can be used to filter out values.
// The filterer should return true when the value should be kept and false if not.
// If the value passed is not of the correct type the filterer wont be called and the
// value won't be emitted (break)
func FilterInt8(filterer func(int8) bool) piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// Check if the value can be used
		var val int8
		var save bool = r.IsValue
		if save {
			val, save = r.Value.(int8)
		}

		// Use the state of the predecessor. But if the value is not save to use and we do not
		// get the closed signal then we break out of the chain.
		state := r.State
		if !save && state != piper.Closed {
			state = piper.Break
		}

		return piper.PipeResult{
			Value:   val,
			IsValue: save && filterer(val),
			State:   state,
			IsEvent: false,
		}, nil
	}, InitialStorage: nil, EventEmitter: nil}
}

// FilterInt16 can be used to filter out values.
// The filterer should return true when the value should be kept and false if not.
// If the value passed is not of the correct type the filterer wont be called and the
// value won't be emitted (break)
func FilterInt16(filterer func(int16) bool) piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// Check if the value can be used
		var val int16
		var save bool = r.IsValue
		if save {
			val, save = r.Value.(int16)
		}

		// Use the state of the predecessor. But if the value is not save to use and we do not
		// get the closed signal then we break out of the chain.
		state := r.State
		if !save && state != piper.Closed {
			state = piper.Break
		}

		return piper.PipeResult{
			Value:   val,
			IsValue: save && filterer(val),
			State:   state,
			IsEvent: false,
		}, nil
	}, InitialStorage: nil, EventEmitter: nil}
}

// FilterInt32 can be used to filter out values.
// The filterer should return true when the value should be kept and false if not.
// If the value passed is not of the correct type the filterer wont be called and the
// value won't be emitted (break)
func FilterInt32(filterer func(int32) bool) piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// Check if the value can be used
		var val int32
		var save bool = r.IsValue
		if save {
			val, save = r.Value.(int32)
		}

		// Use the state of the predecessor. But if the value is not save to use and we do not
		// get the closed signal then we break out of the chain.
		state := r.State
		if !save && state != piper.Closed {
			state = piper.Break
		}

		return piper.PipeResult{
			Value:   val,
			IsValue: save && filterer(val),
			State:   state,
			IsEvent: false,
		}, nil
	}, InitialStorage: nil, EventEmitter: nil}
}

// FilterInt64 can be used to filter out values.
// The filterer should return true when the value should be kept and false if not.
// If the value passed is not of the correct type the filterer wont be called and the
// value won't be emitted (break)
func FilterInt64(filterer func(int64) bool) piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// Check if the value can be used
		var val int64
		var save bool = r.IsValue
		if save {
			val, save = r.Value.(int64)
		}

		// Use the state of the predecessor. But if the value is not save to use and we do not
		// get the closed signal then we break out of the chain.
		state := r.State
		if !save && state != piper.Closed {
			state = piper.Break
		}

		return piper.PipeResult{
			Value:   val,
			IsValue: save && filterer(val),
			State:   state,
			IsEvent: false,
		}, nil
	}, InitialStorage: nil, EventEmitter: nil}
}

// FilterFloat32 can be used to filter out values.
// The filterer should return true when the value should be kept and false if not.
// If the value passed is not of the correct type the filterer wont be called and the
// value won't be emitted (break)
func FilterFloat32(filterer func(float32) bool) piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// Check if the value can be used
		var val float32
		var save bool = r.IsValue
		if save {
			val, save = r.Value.(float32)
		}

		// Use the value of the predecessor. But if the value is not save to use and we do not
		// get the closed signal then we break out of the chain.
		state := r.State
		if !save && state != piper.Closed {
			state = piper.Break
		}

		return piper.PipeResult{
			Value:   val,
			IsValue: save && filterer(val),
			State:   state,
			IsEvent: false,
		}, nil
	}, InitialStorage: nil, EventEmitter: nil}
}

// FilterFloat64 can be used to filter out values.
// The filterer should return true when the value should be kept and false if not.
// If the value passed is not of the correct type the filterer wont be called and the
// value won't be emitted (break)
func FilterFloat64(filterer func(float64) bool) piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// Check if the value can be used
		var val float64
		var save bool = r.IsValue
		if save {
			val, save = r.Value.(float64)
		}

		// Use the value of the predecessor. But if the value is not save to use and we do not
		// get the closed signal then we break out of the chain.
		state := r.State
		if !save && state != piper.Closed {
			state = piper.Break
		}

		return piper.PipeResult{
			Value:   val,
			IsValue: save && filterer(val),
			State:   state,
			IsEvent: false,
		}, nil
	}, InitialStorage: nil, EventEmitter: nil}
}

// FilterString can be used to filter out values.
// The filterer should return true when the value should be kept and false if not.
// If the value passed is not of the correct type the filterer wont be called and the
// value won't be emitted (break)
func FilterString(filterer func(string) bool) piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// Check if the value can be used
		var val string
		var save bool = r.IsValue
		if save {
			val, save = r.Value.(string)
		}

		// Use the value of the predecessor. But if the value is not save to use and we do not
		// get the closed signal then we break out of the chain.
		state := r.State
		if !save && state != piper.Closed {
			state = piper.Break
		}

		return piper.PipeResult{
			Value:   val,
			IsValue: save && filterer(val),
			State:   state,
			IsEvent: false,
		}, nil
	}, InitialStorage: nil, EventEmitter: nil}
}

// FilterBool can be used to filter out values.
// The filterer should return true when the value should be kept and false if not.
// If the value passed is not of the correct type the filterer wont be called and the
// value won't be emitted (break)
func FilterBool(filterer func(bool) bool) piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// Check if the value can be used
		var val bool
		var save bool = r.IsValue
		if save {
			val, save = r.Value.(bool)
		}

		// Use the value of the predecessor. But if the value is not save to use and we do not
		// get the closed signal then we break out of the chain.
		state := r.State
		if !save && state != piper.Closed {
			state = piper.Break
		}

		return piper.PipeResult{
			Value:   val,
			IsValue: save && filterer(val),
			State:   state,
			IsEvent: false,
		}, nil
	}, InitialStorage: nil, EventEmitter: nil}
}

// TODO add more datatypes and refactor the duplicate code but keep the clarity
