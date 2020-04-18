package transform

import (
	"github.com/r-the-thinker/piper"
)

// BufferCount buffers >num< elements and emits them as slice if the predecessor closed
// without the slice being full then the len of the array will be less then >num<.
// Also the pipefunc will only return IsValue as true if there is at least one element
// in the slice.
// num -  in range [1, MAX_INT)
func BufferCount(num int) piper.PipeOperator {
	if num < 1 {
		num = 1
	}
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		arr := storage.([]interface{})
		var nStorage []interface{} = nil

		// if we get a value then safe that to the buffer
		// and if necessary allocate a new array
		if r.IsValue {
			arr = append(arr, r.Value)
			if len(arr) == num && r.State != piper.Closed {
				nStorage = make([]interface{}, 0, num)
			}
		}

		// if we have enough elements or the piper closed then use the nStorage which will be
		// nil (when the predecessor closed) or it has been initialized to a new array
		if len(arr) == num || r.State == piper.Closed {
			// fmt.Println("Emitting the buffer", arr, len(arr))
			return piper.PipeResult{Value: arr, IsValue: len(arr) != 0, State: r.State}, nStorage
		}

		// If we don't emit then just save the new array in the storage
		return piper.PipeResult{IsValue: false, State: piper.Break}, arr

	}, InitialStorage: make([]interface{}, 0, num), EventEmitter: nil}
}
