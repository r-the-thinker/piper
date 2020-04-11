package filtering

import (
	"fmt"
	"time"

	"github.com/r-the-thinker/piper"
)

// AuditTime waits the specified duration upon receiving an element.
// It will output the last value it received during that time.
// dur - the duration to wait
func AuditTime(dur time.Duration) piper.PipeOperator {
	emitter := make(chan interface{})
	comChan := make(chan piper.PipeResult, 1)

	// start the worker which handles the waiting
	go worker(dur, emitter, comChan)

	f := func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		fmt.Println(r)
		// if it is an event then the worker sent what we should do
		if r.IsEvent {
			return r.Value.(piper.PipeResult), nil
		}

		// if it is not an event then send the piperesult to the worker and break
		// only the worker decides what we do.
		comChan <- r
		return piper.PipeResult{IsValue: false, State: piper.Break}, nil
	}

	return piper.PipeOperator{F: f, InitialStorage: nil, EventEmitter: emitter}
}

func worker(dur time.Duration, emitter chan interface{}, comChan chan piper.PipeResult) {
	var timer <-chan time.Time
	var mostRecentResult piper.PipeResult = piper.PipeResult{}
	var valueExists bool = false

	for {
		select {
		case nResult := <-comChan:
			// if we dont have an item yet then start the timer
			if !valueExists {
				timer = time.After(dur)
				valueExists = true
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
			emitter <- mostRecentResult
			timer = nil
			valueExists = false

			// check if we are going to end now
			if mostRecentResult.State == piper.Closed {
				close(comChan)
				break
			}
		}
	}
}
