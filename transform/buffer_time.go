package transform

import (
	"container/list"
	"fmt"
	"time"

	"github.com/r-the-thinker/piper"
)

// BufferTime collects all values received in the given duration and emits
// them as slice, there will only be an emission if there is data. Therefore
// Empty slices wont be emitted.
func BufferTime(dur time.Duration) piper.PipeOperator {
	eventEmitter := make(chan interface{})
	comChan := make(chan piper.PipeResult, 1)

	go bufferTimeWorker(dur, eventEmitter, comChan)

	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		fmt.Println("Got", r)
		if r.IsEvent {
			return r.Value.(piper.PipeResult), nil
		}
		comChan <- r
		return piper.PipeResult{State: piper.Break}, nil
	}, InitialStorage: nil, EventEmitter: eventEmitter}
}

func bufferTimeWorker(dur time.Duration, eventEmitter chan interface{}, comChan chan piper.PipeResult) {
	// items stores all the values that we receive during the buffer period
	// a list should be more efficient here because we do not know how many
	// values that we are going to receive. When an underlaying array is full
	// appending to a slice would cause the array size to be doubled. The
	// reallocation could happen several times and the access unneeded memory
	// can be saved using a list. Appending to a list is a constant time
	// operation we only need to iterate over the list once in order to
	// move the elements from the list to the array.
	values := list.New()
	state := piper.Open

	var timer <-chan time.Time = nil

loop:
	for {
		select {
		case nResult := <-comChan:
			// if there is no timer yet, start a new one. If the value that is received
			// has IsValue=false which is only possible if the predecessor closed, then
			// the new timer should fire immediately
			if timer == nil {
				if !nResult.IsValue {
					timer = time.After(0)
				} else {
					timer = time.After(dur)
				}
			}

			if nResult.IsValue {
				values.PushBack(nResult.Value)
			}
			state = nResult.State
		case <-timer:
			// Create the result array once and put in all values
			arr := make([]interface{}, values.Len())
			i := 0
			for e := values.Front(); e != nil; e = e.Next() {
				arr[i] = e.Value
				i++
			}

			eventEmitter <- piper.PipeResult{Value: arr, IsValue: len(arr) != 0, State: state}

			// break out of the loop once the closed state has been sent
			if state == piper.Closed {
				break loop
			}

			timer = nil
			values = list.New()
		}
	}

	close(eventEmitter)
	close(comChan)
}
