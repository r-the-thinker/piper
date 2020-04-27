package filtering_test

import (
	"testing"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/filtering"
)

func TestSkip(t *testing.T) {
	inChan := make(chan int, 2)
	outChan := piper.Clone(inChan).Pipe(filtering.Skip(2)).Get().(chan int)

	inChan <- 1
	inChan <- 2
	inChan <- 3
	inChan <- 4
	close(inChan)

	if val1, val2 := <-outChan, <-outChan; val1 != 3 || val2 != 4 {
		t.Fatalf("Expected to receive 3 and 4 but got %v and %v", val1, val2)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestSkipCloseWithValue(t *testing.T) {
	closer := piper.PipeOperator{F: func(r piper.PipeResult, _ interface{}) (piper.PipeResult, interface{}) {
		return piper.PipeResult{Value: 5, IsValue: true, State: piper.Closed}, nil
	}}

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(closer, filtering.Skip(1)).Get().(chan int)
	close(inChan)

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}
