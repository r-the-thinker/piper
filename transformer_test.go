package piper_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/r-the-thinker/piper"
)

// Pipefuncs can be triggered with an event emitter and their result is retuned normally
func TestEventChain(t *testing.T) {
	inChan := make(chan int, 10)

	// the channel we wait on, which will be filled by the pipefunc
	valueCheck := make(chan interface{})
	p1 := func(r piper.PipeResult, s interface{}) (piper.PipeResult, interface{}) {
		return piper.PipeResult{Value: 1, IsValue: true, State: piper.Open}, nil
	}
	p2 := func(r piper.PipeResult, s interface{}) (piper.PipeResult, interface{}) {
		valueCheck <- r.Value
		return piper.PipeResult{Value: 99, IsValue: true, State: piper.Open}, nil
	}

	// The emitter the Piper should listen to
	eventEmitter := make(chan interface{})
	outChan := piper.Clone(inChan).Pipe(piper.PipeOperator{p1, nil, nil}).Pipe(piper.PipeOperator{p2, nil, eventEmitter}).Get().(chan int)

	eventEmitter <- true
	// Read what value was  received it should be true because this is what the EventEmitter said
	if val := <-valueCheck; val.(bool) != true {
		t.Fatalf("The event pipefunc was called with %v as value but it should have been true from the emitter", val)
	}

	if val := <-outChan; val != 99 {
		t.Fatalf("The value that was returned by the through an Event triggered function was %v not 99", val)
	}

	close(inChan)
}

// The storage is safed and can be used even if not registered with pipe
func TestPipeNoInitialStorage(t *testing.T) {
	inChan := make(chan int, 2)
	resultChan := make(chan error, 2)
	// The pipefunc that does the testing
	testPipefunc := func(r piper.PipeResult, s interface{}) (piper.PipeResult, interface{}) {
		// At the beginning there should be no value in the storage
		if r.Value.(int) == 0 {
			if s != nil {
				resultChan <- fmt.Errorf("Expected the storage to be empty when it was initialized with nil")
			} else {
				resultChan <- nil
			}
		}
		// After the first iteration there should be a value in the storage
		if r.Value.(int) == 1 {
			if s == nil {
				resultChan <- fmt.Errorf("Expected the storage not to be empty because we returned a value")
			} else {
				resultChan <- nil
			}
		}
		return piper.PipeResult{Value: r.Value, IsValue: true, State: piper.Open}, 1
	}

	// Create the Piper with nil as initial storage
	piper.Clone(inChan).Pipe(piper.PipeOperator{testPipefunc, nil, nil}).Get()

	// Send the values in so the pipefunc knows which state it is in
	inChan <- 0
	inChan <- 1

	// Read from the outchan to make sure that the pipefunc was called
	result1 := <-resultChan
	result2 := <-resultChan
	if result1 != nil {
		t.Error(result1)
	}
	if result2 != nil {
		t.Error(result2)
	}
	close(inChan)
}

// The initial storage is passed to the PIpefunc
func TestPipeWithInitialStorage(t *testing.T) {
	inChan := make(chan int, 2)
	resultChan := make(chan error, 2)
	// The pipefunc that does the testing
	testPipefunc := func(r piper.PipeResult, s interface{}) (piper.PipeResult, interface{}) {
		// At the beginning there should be the initial value
		if r.Value.(int) == 0 {
			if s.(int) != 1 {
				resultChan <- fmt.Errorf("Expected the storage to contain the value 1 but it was %v", s)
			} else {
				resultChan <- nil
			}
		}
		// After the first iteration there should be a value in the storage
		if r.Value.(int) == 1 {
			if s == nil {
				resultChan <- fmt.Errorf("Expected the storage to contain the value 2 but it was %v", s)
			} else {
				resultChan <- nil
			}
		}
		return piper.PipeResult{Value: r.Value, IsValue: true, State: piper.Open}, 1
	}

	// Create the Piper with nil as initial storage
	piper.Clone(inChan).Pipe(piper.PipeOperator{testPipefunc, 1, nil}).Get()

	// Send the values in so the pipefunc knows which state it is in
	inChan <- 0
	inChan <- 1

	// Read from the outchan to make sure that the pipefunc was called
	result1 := <-resultChan
	result2 := <-resultChan
	if result1 != nil {
		t.Error(result1)
	}
	if result2 != nil {
		t.Error(result2)
	}
	close(inChan)
}

func TestClosingNoFuncs(t *testing.T) {
	inChan := make(chan int)
	outChan := piper.Clone(inChan).Get().(chan int)

	close(inChan)
	if _, ok := <-outChan; ok {
		t.Fatal("The output channel should have been closed")
	}
}

func TestClosingOneFunc(t *testing.T) {
	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(piper.PipeOperator{
		func(r piper.PipeResult, s interface{}) (piper.PipeResult, interface{}) {
			return piper.PipeResult{State: piper.Closed}, nil
		}, nil, nil}).Get().(chan int)

	close(inChan)
	if _, ok := <-outChan; ok {
		t.Fatal("The output channel should have been closed")
	}
}

func TestClosingTwoFunc(t *testing.T) {
	inChan := make(chan int)

	op := piper.PipeOperator{
		func(r piper.PipeResult, s interface{}) (piper.PipeResult, interface{}) {
			return piper.PipeResult{State: piper.Closed}, nil
		}, nil, nil}
	outChan := piper.Clone(inChan).Pipe(op, op).Get().(chan int)

	close(inChan)
	if _, ok := <-outChan; ok {
		t.Fatal("The output channel should have been closed")
	}
}

func TestClosingKeepOpen(t *testing.T) {
	inChan := make(chan int)

	op1 := piper.PipeOperator{
		func(r piper.PipeResult, s interface{}) (piper.PipeResult, interface{}) {
			return piper.PipeResult{State: piper.Closed}, nil
		}, nil, nil}

	emitter := make(chan interface{})
	op2 := piper.PipeOperator{
		func(r piper.PipeResult, s interface{}) (piper.PipeResult, interface{}) {
			// when getting the event callback then close
			if r.IsEvent {
				close(emitter)
				return piper.PipeResult{Value: 1, IsValue: true, State: piper.Closed}, nil
			}
			// keep open no matter what unless event...
			return piper.PipeResult{State: piper.Open}, nil
		}, nil, emitter}

	outChan := piper.Clone(inChan).Pipe(op1, op1, op2).Get().(chan int)

	close(inChan)
	emitter <- true
	// we get one emission after the event before it closes
	if val := <-outChan; val != 1 {
		t.Fatalf("Expected to read 1 from the output channel but got %v", val)
	}
	if _, ok := <-outChan; ok {
		t.Fatalf("Expected the output channel to be closed but it's not")
	}
}

func TestOnlyOneIsEvent(t *testing.T) {
	inChan := make(chan int)

	eventEmitter := make(chan interface{})
	checkChan := make(chan bool, 2)
	f := func(r piper.PipeResult, s interface{}) (piper.PipeResult, interface{}) {
		checkChan <- r.IsEvent
		return piper.PipeResult{Value: 1, IsValue: true, State: piper.Open}, nil
	}
	o1 := piper.PipeOperator{f, nil, eventEmitter}
	o2 := piper.PipeOperator{f, nil, nil}
	piper.Clone(inChan).Pipe(o1, o2)

	eventEmitter <- 1
	if val := <-checkChan; !val {
		t.Fatal("Expected the first pipefunc to have been called as Event. But it's not")
	}
	if val := <-checkChan; val {
		t.Fatal("Expected the second pipefunc to not have been called as Event. But it was")
	}
}

// Closing an EventEmitter should have no influcence on the sent values and states
func TestEventEmitterInflucence(t *testing.T) {
	inChan := make(chan int, 2)

	eventEmitter := make(chan interface{})
	o := piper.PipeOperator{func(r piper.PipeResult, s interface{}) (piper.PipeResult, interface{}) {
		return r, s
	}, nil, eventEmitter}
	outChan := piper.Clone(inChan).Pipe(o).Get().(chan int)

	inChan <- 1
	close(eventEmitter)
	inChan <- 2

	if v1, v2 := <-outChan, <-outChan; v1 != 1 || v2 != 2 {
		t.Fatalf("Expected 1 and 2, got %v and %v", v1, v2)
	}
}

func TestClosingEarly(t *testing.T) {
	// this is the mapper that we want to be closed when the 2nd closes
	f1 := func(r piper.PipeResult, s interface{}) (piper.PipeResult, interface{}) {
		return piper.PipeResult{Value: 5, IsValue: r.IsValue, State: r.State}, nil
	}

	// this func closes when getting value 2
	f2 := func(r piper.PipeResult, s interface{}) (piper.PipeResult, interface{}) {
		if s.(int) == 1 {
			return piper.PipeResult{Value: r.Value, IsValue: true, State: piper.Open}, 2
		}
		return piper.PipeResult{Value: r.Value, IsValue: true, State: piper.Closed}, nil
	}

	inChan := make(chan int, 3)
	outChan := piper.Clone(inChan).Pipe(
		piper.PipeOperator{F: f1, InitialStorage: nil, EventEmitter: nil},
		piper.PipeOperator{F: f2, InitialStorage: 1, EventEmitter: nil}).Get().(chan int)

	inChan <- 1 // should be mapped to 2
	inChan <- 2 // should be mapped to 4
	inChan <- 3 // should not be registered

	if val1, val2 := <-outChan, <-outChan; val1 != 5 || val2 != 5 {
		t.Fatalf("Expected to receive 5 and 5 but got %v and %v", val1, val2)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it is not")
	}
}

// We know that Get gives us a channel
func TestGet(t *testing.T) {
	c := make(chan int)
	v := reflect.ValueOf(piper.Clone(c).Get())

	if v.Kind() != reflect.Chan {
		t.Fatalf("Get did not return a channel")
	}

	close(c)
}
