package transform_test

import (
	"testing"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/transform"
)

var sliceToIntMapper piper.PipeOperator = piper.PipeOperator{F: func(r piper.PipeResult, s interface{}) (piper.PipeResult, interface{}) {
	sum := 0
	for _, val := range r.Value.([]interface{}) {
		sum += val.(int)
	}
	return piper.PipeResult{Value: sum, IsValue: r.IsValue, State: r.State}, nil
}, InitialStorage: nil, EventEmitter: nil}

func TestBufferCountNoValues(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(transform.BufferCount(1), sliceToIntMapper).Get().(chan int)
	close(inChan)

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestBufferCountBufferManyValues(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(transform.BufferCount(1), sliceToIntMapper).Get().(chan int)

	inChan <- 1
	close(inChan)

	if val := <-outChan; val != 1 {
		t.Fatalf("Expected the sum of the inputs to be 1 but it is %v", val)
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestBufferCountLessThenBufferManyValues(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(transform.BufferCount(2), sliceToIntMapper).Get().(chan int)

	inChan <- 1
	close(inChan)

	if val := <-outChan; val != 1 {
		t.Fatalf("Expected the sum of the inputs to be 1 but it is %v", val)
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestBufferCountBufferManyValuesTwice(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(transform.BufferCount(2), sliceToIntMapper).Get().(chan int)

	inChan <- 1
	inChan <- 2

	if val := <-outChan; val != 3 {
		t.Fatalf("Expected the sum of the inputs to be 3 but it is %v", val)
	}

	inChan <- 3
	inChan <- 4
	close(inChan)

	if val := <-outChan; val != 7 {
		t.Fatalf("Expected the sum of the inputs to be 7 but it is %v", val)
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestBufferCountMixedCount(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(transform.BufferCount(2), sliceToIntMapper).Get().(chan int)

	inChan <- 1
	inChan <- 2

	if val := <-outChan; val != 3 {
		t.Fatalf("Expected the sum of the inputs to be 3 but it is %v", val)
	}

	inChan <- 4
	close(inChan)

	if val := <-outChan; val != 4 {
		t.Fatalf("Expected the sum of the inputs to be 4 but it is %v", val)
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestBufferCountCloseWithValue(t *testing.T) {
	t.Parallel()

	closer := piper.PipeOperator{F: func(r piper.PipeResult, _ interface{}) (piper.PipeResult, interface{}) {
		return piper.PipeResult{Value: 1, IsValue: true, State: piper.Closed}, nil
	}}

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(closer, transform.BufferCount(2), sliceToIntMapper).Get().(chan int)
	inChan <- 1

	if val := <-outChan; val != 1 {
		t.Fatalf("Expected the sum of the inputs to be 1 but it is %v", val)
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}

	close(inChan)
}
