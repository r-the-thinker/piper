package filtering

import (
	"fmt"
	"time"

	"github.com/r-the-thinker/piper"
)

// Throttle emits a value and then ignores all following values based on the the
// duration that was calculated off of the value emitted.
func Throttle(throttleCalcFunc func(interface{}) time.Duration) piper.PipeOperator {
	eventEmitter := make(chan interface{}, 1)
	comChan := make(chan piper.PipeResult, 1)

	go throttleWorker(throttleCalcFunc, eventEmitter, comChan)

	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		if r.IsEvent {
			fmt.Println("got", r.Value.(piper.PipeResult), "from worker")
			return r.Value.(piper.PipeResult), nil
		}

		comChan <- r
		return piper.PipeResult{State: piper.Break}, nil
	}, InitialStorage: nil, EventEmitter: eventEmitter}
}

func throttleWorker(throttleCalcFunc func(interface{}) time.Duration, eventEmitter chan interface{}, comChan chan piper.PipeResult) {
	var timer <-chan time.Time
	var state piper.PipeState

loop:
	for {
		select {
		case nResult := <-comChan:
			fmt.Println("worker received", nResult, timer)
			if timer == nil {
				// start the timer based on the value, if it is not a value (the predecessor closed without a value)
				// then we emit it but do no start a new timer
				if nResult.IsValue {
					timer = time.After(throttleCalcFunc(nResult.Value))
				}
				// if the timer was nil then it is the first one,
				// therefore emit the item
				eventEmitter <- nResult

				if nResult.State == piper.Closed {
					break loop
				}
			}

			// whenever we get the closed signal but a timer is already active,
			// then the value is not going to be recognized. But we can issue a closing
			// signal right away because there is no need that someone waits for a closing
			// without a value.
			if nResult.State == piper.Closed {
				eventEmitter <- piper.PipeResult{State: piper.Closed, IsValue: false}
				break
			}
		case <-timer:
			timer = nil
			fmt.Println("Timer executed")

			// if it was a closed state then we break out of here and die
			if state == piper.Closed {
				eventEmitter <- piper.PipeResult{State: piper.Closed}
				break loop
			}
		}
	}
	close(comChan)
	close(eventEmitter)
}
