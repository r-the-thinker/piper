package filtering

import (
	"time"

	"github.com/r-the-thinker/piper"
)

// Debounce does upon receiving a value call the supplied function that calculates
// the time window that is going to be waited before emitting that value. But every new
// value that comes in resets the timer, therefore only if there are no new values comming
// in during 	the duration the value will be emitted. The window of time is only calculated at
// the beginning (if there weren't any emsissions yet) and after a successful emission.
// When a value comes in during the waiting period there won't be a recalculation of
// that window of time.
func Debounce(durCalc func(interface{}) time.Duration) piper.PipeOperator {
	eventEmitter := make(chan interface{}, 1)
	comChan := make(chan piper.PipeResult)

	// start the worker which takes care of the waiting etc.
	go debounceWorker(durCalc, eventEmitter, comChan)

	f := func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		if r.IsEvent {
			// when we receive a command from the worker then emit that
			return r.Value.(piper.PipeResult), nil
		}
		// if it is not an event then relay the recieved value to the worker
		// so that it can decide how to handle it
		comChan <- r
		return piper.PipeResult{IsValue: false, State: piper.Break}, nil
	}

	return piper.PipeOperator{F: f, InitialStorage: nil, EventEmitter: eventEmitter}
}

func debounceWorker(durCalc func(interface{}) time.Duration, eventEmitter chan interface{}, comChan chan piper.PipeResult) {
	var timer <-chan time.Time
	var mostRecentResult piper.PipeResult = piper.PipeResult{}
	var currentWaitingDuration time.Duration = time.Duration(0)

loop:
	for {
		select {
		case nResult := <-comChan:

			// if we get a value without there being an active wait time then
			// schedule a new timer
			if timer == nil {
				// if we only get the closing command without it being a value
				// and there is no timer yet that we are waiting on then end it
				// right away. Otherwise start the timer for the specified time.
				if nResult.State == piper.Closed && !nResult.IsValue {
					timer = time.After(time.Duration(0))
				} else {
					currentWaitingDuration = durCalc(nResult.Value)
					timer = time.After(currentWaitingDuration)
				}
			} else {
				// if there is an active timer then restart the timer with
				// the old waiting duration
				timer = time.After(currentWaitingDuration)
			}

			// Only override our current value if we actually receive a value
			// this must be done because the Closed emission has IsValue=false
			// so be do not want to override our current value with the null value
			if nResult.IsValue {
				mostRecentResult.Value = nResult.Value
				mostRecentResult.IsValue = true
			}
			// Always remember the lastest state. Closed cannot be overridden because
			// it's always the last emission
			mostRecentResult.State = nResult.State
		case <-timer:
			// if the timer successfully ended then make the emission and reset the values
			eventEmitter <- mostRecentResult
			timer = nil

			// check if we are going to end now
			if mostRecentResult.State == piper.Closed {
				close(comChan)
				break loop
			}
		}
	}
}
