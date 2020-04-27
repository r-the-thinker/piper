package filtering_test

import (
	"testing"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/filtering"
)

func TestFindOpen(t *testing.T) {
	inChan := make(chan int, 1)

	predicate := func(val interface{}) bool {
		return val.(int) > 2
	}
	outChan := piper.Clone(inChan).Pipe(filtering.Find(predicate)).Get().(chan int)

	// this should not block because we only get 3 and then nothing
	inChan <- 1
	inChan <- 2
	inChan <- 3
	inChan <- 4
	close(inChan)

	if val := <-outChan; val != 3 {
		t.Fatalf("Expected to receive 3 but got %v", val)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the channel to be closed")
	}
}

func TestFindClosed(t *testing.T) {
	inChan := make(chan int, 1)

	predicate := func(val interface{}) bool {
		return val.(int) > 2
	}
	outChan := piper.Clone(inChan).Pipe(filtering.Find(predicate)).Get().(chan int)

	// this should not block because we only get 3 and then nothing
	inChan <- 1
	inChan <- 2
	close(inChan)

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the channel to be closed")
	}
}

func TestFindWithClosedFound(t *testing.T) {
	closer := piper.PipeOperator{F: func(r piper.PipeResult, _ interface{}) (piper.PipeResult, interface{}) {
		return piper.PipeResult{Value: 5, IsValue: true, State: piper.Closed}, nil
	}}
	predicate := func(val interface{}) bool {
		return val.(int) > 2
	}

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(closer, filtering.Filter(predicate)).Get().(chan int)
	close(inChan)

	if val := <-outChan; val != 5 {
		t.Fatalf("Expected to receive 5 but got %v", val)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestFindWithClosedNotFound(t *testing.T) {
	closer := piper.PipeOperator{F: func(r piper.PipeResult, _ interface{}) (piper.PipeResult, interface{}) {
		return piper.PipeResult{Value: 5, IsValue: true, State: piper.Closed}, nil
	}}
	predicate := func(val interface{}) bool {
		return val.(int) > 11
	}

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(closer, filtering.Filter(predicate)).Get().(chan int)
	close(inChan)

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}
