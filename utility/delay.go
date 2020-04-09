package utility

import (
	"container/list"
	"time"

	piper "github.com/r-the-thinker/piper"
)

type delayStorage struct {
	isWaiting bool
	isClosed  bool
	items     *list.List
}

// Delay is used to delay a given input by the specified duration
// The delay is not exact excact as it competes with other channels.
// This can only be used with PipeSingle!
func Delay(delay time.Duration) piper.PipeOperator {

	// to signal the pipefunc that it can emit
	eventEmitter := make(chan interface{}, 1)
	// to send to the runner below that we want another delay scheduled
	scheduleChan := make(chan bool, 1)
	go func() {
		var timer <-chan time.Time
		var timeCtl *time.Timer = nil

		for {
			select {
			// when the delay has finished we sent out the notification and ensure that
			// the timer channel is blocking
			case <-timer:
				eventEmitter <- timeCtl.Stop()
			// if we receive from the schedule chan then the pipefunc wants a new delay
			// so create that. If the scheudle chan closes we have to kill the goroutine
			case _, ok := <-scheduleChan:
				if !ok {
					timeCtl.Stop()
					return
				}

				if timer == nil {
					timeCtl = time.NewTimer(delay)
				} else {
					timeCtl.Reset(delay)
				}
				timer = timeCtl.C
			}
		}
	}()

	f := func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		s := storage.(delayStorage)

		// we can close if the predecessor closed but we are not getting a value to
		// process and have nothing left to emit
		if r.State == piper.Closed && !r.IsValue && s.items.Len() == 0 {
			close(eventEmitter)
			close(scheduleChan)
			return piper.PipeResult{State: piper.Closed}, nil
		}

		// Handle event callback
		if r.IsEvent {
			// get the item that we want to emit and by default we want to close now
			item := s.items.Remove(s.items.Front())

			// if there are items left then we schedule a new delay
			if s.items.Len() > 0 {
				scheduleChan <- true
				s.isWaiting = true
			}
			state := getDelayReturnState(&s)
			if state == piper.Closed {
				close(eventEmitter)
				close(scheduleChan)
			}

			return piper.PipeResult{Value: item, IsValue: true, State: state}, s
		}

		// we only get the closing signal once at the end. We have to remember that so that
		// we can close if it is time. After this round we can only get event callbacks so there
		// wont be a modification to this variable.
		s.isClosed = r.State == piper.Closed
		// if we get a new value in we have to save it and schedule an event emitter if we are
		// not already in a waiting state.
		if r.IsValue {
			s.items.PushBack(r.Value)
			if !s.isWaiting {
				scheduleChan <- true
				s.isWaiting = true
			}
		}
		return piper.PipeResult{State: piper.Break}, s
	}

	return piper.PipeOperator{F: f, InitialStorage: delayStorage{false, false, list.New()}, EventEmitter: eventEmitter}
}

func getDelayReturnState(s *delayStorage) piper.PipeState {
	// if the predecessor is open or if we still have items then we need to stay open
	if !s.isClosed || s.items.Len() > 0 {
		return piper.Open
	}

	// There is no reason to stay open now so also close
	return piper.Closed
}
