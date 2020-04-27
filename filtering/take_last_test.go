package filtering_test

import (
	"testing"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/filtering"
)

func TestTakeUntilOne(t *testing.T) {
	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(filtering.TakeLast(1)).Get().(chan int)

	inChan <- 1
	close(inChan)

	if val := <-outChan; val != 1 {
		t.Fatalf("Expected to receive 1 from the ouput channel, but got %v instead", val)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestTakeUntilMultipleSameAsBuffer(t *testing.T) {
	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(filtering.TakeLast(2)).Get().(chan int)

	inChan <- 1
	inChan <- 2
	close(inChan)

	if val1, val2 := <-outChan, <-outChan; val1 != 1 || val2 != 2 {
		t.Fatalf("Expected to receive 1 and 2 from the ouput channel, but got %v and %v instead", val1, val2)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestTakeUntilMultipleLessThanBuffer(t *testing.T) {
	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(filtering.TakeLast(2)).Get().(chan int)

	inChan <- 1
	close(inChan)

	if val := <-outChan; val != 1 {
		t.Fatalf("Expected to receive 1 from the ouput channel, but got %v instead", val)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestTakeUntilMultipleMoreThanBuffer(t *testing.T) {
	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(filtering.TakeLast(2)).Get().(chan int)

	inChan <- 1
	inChan <- 2
	inChan <- 3
	close(inChan)

	if val1, val2 := <-outChan, <-outChan; val1 != 2 || val2 != 3 {
		t.Fatalf("Expected to receive 2 and 3 from the ouput channel, but got %v and %v instead", val1, val2)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestTakeUntilNone(t *testing.T) {
	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(filtering.TakeLast(2)).Get().(chan int)
	close(inChan)

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestTakeUntilClosedWithValue(t *testing.T) {
	closer := piper.PipeOperator{F: func(r piper.PipeResult, _ interface{}) (piper.PipeResult, interface{}) {
		return piper.PipeResult{Value: 1, IsValue: true, State: piper.Closed}, nil
	}}

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(closer, filtering.TakeLast(2)).Get().(chan int)

	close(inChan)

	if val := <-outChan; val != 1 {
		t.Fatalf("Expected to receive 1 from the ouput channel, but got %v instead", val)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}
