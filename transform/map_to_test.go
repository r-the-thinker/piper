package transform_test

import (
	"testing"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/transform"
)

func TestMapToNoValue(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(transform.MapTo(2)).Get().(chan int)
	close(inChan)

	if _, ok := <-outChan; ok {
		t.Fatalf("Expected the output channel to be closed but it's not")
	}
}
func TestMapToOneValue(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(transform.MapTo(2)).Get().(chan int)
	inChan <- 1
	close(inChan)

	if val := <-outChan; val != 2 {
		t.Fatalf("Expected to receive 2 but got %v instead", val)
	}
	if _, ok := <-outChan; ok {
		t.Fatalf("Expected the output channel to be closed but it's not")
	}
}

func TestMapToMultipleValues(t *testing.T) {
	t.Parallel()

	inChan := make(chan int, 2)
	outChan := piper.From(inChan).Pipe(transform.MapTo(2)).Get().(chan int)
	inChan <- 1
	inChan <- 2
	close(inChan)

	if val1, val2 := <-outChan, <-outChan; val1 != 2 || val2 != 2 {
		t.Fatalf("Expected to receive 2 and 2 but got %v and %v instead", val1, val2)
	}
	if _, ok := <-outChan; ok {
		t.Fatalf("Expected the output channel to be closed but it's not")
	}
}

func TestMapToCloseWithValue(t *testing.T) {
	t.Parallel()

	closer := piper.PipeOperator{F: func(r piper.PipeResult, _ interface{}) (piper.PipeResult, interface{}) {
		return piper.PipeResult{Value: 1, IsValue: true, State: piper.Closed}, nil
	}}

	inChan := make(chan int, 2)
	outChan := piper.From(inChan).Pipe(closer, transform.MapTo(2)).Get().(chan int)
	inChan <- 9
	close(inChan)

	if val := <-outChan; val != 2 {
		t.Fatalf("Expected to receive 2 but got %v instead", val)
	}
	if _, ok := <-outChan; ok {
		t.Fatalf("Expected the output channel to be closed but it's not")
	}
}
