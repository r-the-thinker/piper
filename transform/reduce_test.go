package transform_test

import (
	"testing"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/transform"
)

func filterReducer(accumulator, currentVal interface{}) interface{} {
	return accumulator.(int) + currentVal.(int)
}

func TestReduce(t *testing.T) {
	t.Parallel()

	inputChan := make(chan int, 3)
	outputChan := piper.Clone(inputChan).Pipe(transform.Reduce(filterReducer, 0)).Get().(chan int)

	inputChan <- 1
	inputChan <- 2
	inputChan <- 3
	close(inputChan)
	if v := <-outputChan; v != 6 {
		t.Fatalf("Expected to receive 2, but got %v", v)
	}
}

func TestReduceNoValue(t *testing.T) {
	t.Parallel()

	inputChan := make(chan int)
	outputChan := piper.Clone(inputChan).Pipe(transform.Reduce(filterReducer, 0)).Get().(chan int)
	close(inputChan)

	if _, ok := <-outputChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestReduceCloedWithValue(t *testing.T) {
	t.Parallel()

	closer := piper.PipeOperator{F: func(r piper.PipeResult, _ interface{}) (piper.PipeResult, interface{}) {
		return piper.PipeResult{Value: 1, IsValue: true, State: piper.Closed}, nil
	}}

	inChan := make(chan int)
	outputChan := piper.Clone(inChan).Pipe(closer).Pipe(transform.Reduce(filterReducer, 0)).Get().(chan int)
	inChan <- 1
	close(inChan)

	if val := <-outputChan; val != 1 {
		t.Fatalf("Expected to receive 1 but got %v instead", val)
	}

	if _, ok := <-outputChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}
