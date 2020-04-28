package piper_test

import (
	"reflect"
	"testing"

	"github.com/r-the-thinker/piper"
)

func TestCloneType(t *testing.T) {
	inChan := make(chan int)
	outChan := piper.Clone(inChan).Get().(chan int)
	if c := cap(outChan); c != 0 {
		t.Fatalf("The capacity of the channel should have been 0, but it is %v", c)
	}
	if cType := reflect.TypeOf(outChan).Elem().Kind(); cType != reflect.Int {
		t.Fatalf("From created a channel of type %v, but expected it to be of type int", cType)
	}

	inChan = make(chan int, 10)
	outChan = piper.Clone(inChan).Get().(chan int)
	if c := cap(outChan); c != 10 {
		t.Fatalf("The capacity of the channel should have been 10, but it is %v", c)
	}
	if cType := reflect.TypeOf(outChan).Elem().Kind(); cType != reflect.Int {
		t.Fatalf("From created a channel of type %v, but expected it to be of type int", cType)
	}
	close(inChan)
}

func TestCloneUnbuffered(t *testing.T) {
	inChan := make(chan int)
	outChan := piper.Clone(inChan).Get().(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			inChan <- i
		}
	}()

	for i := 0; i < 10; i++ {
		if val := <-outChan; i != val {
			t.Fatalf("Expected %v received %v", i, val)
		}
	}

	close(inChan)
}

func TestCloneBuffered(t *testing.T) {
	inChan := make(chan int, 10)
	outChan := piper.Clone(inChan).Get().(chan int)

	for i := 0; i < 10; i++ {
		inChan <- i
	}

	for i := 0; i < 10; i++ {
		if val := <-outChan; i != val {
			t.Fatalf("Expected %v received %v", i, val)
		}
	}

	close(inChan)
}

func TestFromTo(t *testing.T) {
	op := piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		if r.IsValue {
			r.Value = float32(r.Value.(int))
		}
		return r, nil
	}, InitialStorage: nil, EventEmitter: nil}

	inChan := make(chan int)
	outChan := make(chan float32, 1)
	piper.FromTo(inChan, outChan).Pipe(op)

	inChan <- 1
	close(inChan)

	if val := <-outChan; val != 1.0 {
		t.Fatalf("Expected to receive 1.0 but got %v instead", val)
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}
