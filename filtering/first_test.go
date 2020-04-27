package filtering_test

import (
	"testing"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/filtering"
)

func TestFirstOneValue(t *testing.T) {
	inChan := make(chan int, 2)
	outChan := piper.Clone(inChan).Pipe(filtering.First()).Get().(chan int)

	inChan <- 1
	close(inChan)
	if val := <-outChan; val != 1 {
		t.Fatalf("Expected to receive 1 but got %v", val)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestFirstTwoValue(t *testing.T) {
	inChan := make(chan int, 2)
	outChan := piper.Clone(inChan).Pipe(filtering.First()).Get().(chan int)

	inChan <- 1
	inChan <- 2
	close(inChan)
	if val := <-outChan; val != 1 {
		t.Fatalf("Expected to receive 1 but got %v", val)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestFirstNoValue(t *testing.T) {
	inChan := make(chan int, 2)
	outChan := piper.Clone(inChan).Pipe(filtering.First()).Get().(chan int)
	close(inChan)

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestFirstWithClosedFound(t *testing.T) {
	closer := piper.PipeOperator{F: func(r piper.PipeResult, _ interface{}) (piper.PipeResult, interface{}) {
		return piper.PipeResult{Value: 5, IsValue: true, State: piper.Closed}, nil
	}}

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(closer, filtering.First()).Get().(chan int)
	close(inChan)

	if val := <-outChan; val != 5 {
		t.Fatalf("Expected to receive 5 but got %v", val)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}
