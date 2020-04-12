package filtering

import (
	"fmt"
	"time"

	"github.com/r-the-thinker/piper"
)

// Audit takes a function that calculates the time that is going to be waited before emitting
// the value. The last item that was received during the waiting period will be the one that
// is going to be emitted. The calculation will only be done for the first item that was sent
// without a timer being active.
func Audit(durCalc func(interface{}) time.Duration) piper.PipeOperator {
	emitter := make(chan interface{})
	comChan := make(chan piper.PipeResult, 1)

	// start the worker which handles the waiting
	go auditWorker(durCalc, emitter, comChan)

	f := func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		// if it is an event then the worker sent what we should do
		if r.IsEvent {
			return r.Value.(piper.PipeResult), nil
		}

		// if it is not an event then send the piperesult to the worker and break
		// only the worker decides what we do.
		fmt.Println("Sending to worker", r)
		comChan <- r
		return piper.PipeResult{IsValue: false, State: piper.Break}, nil
	}

	return piper.PipeOperator{F: f, InitialStorage: nil, EventEmitter: emitter}
}

func auditWorker(durCalc func(interface{}) time.Duration, emitter chan interface{}, comChan chan piper.PipeResult) {
	var timer <-chan time.Time
	var mostRecentResult piper.PipeResult = piper.PipeResult{}

runner:
	for {
		select {
		case nResult := <-comChan:
			// if we are not waiting, start the timer
			if timer == nil {
				timer = time.After(durCalc(nResult.Value))
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

			// check if we are going to end now
			if mostRecentResult.State == piper.Closed {
				close(comChan)
				break runner
			}
		}
	}
}
