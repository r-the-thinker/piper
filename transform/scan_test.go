package transform_test

import (
	"testing"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/transform"
)

func scanAccumulator(prevAccumulated, currentValue interface{}) interface{} {
	return prevAccumulated.(int) + currentValue.(int)
}

func TestScanNoValue(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(transform.Scan(scanAccumulator, 0)).Get().(chan int)
	close(inChan)

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel tot be closed but it's not")
	}
}

func TestScanOneValue(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(transform.Scan(scanAccumulator, 0)).Get().(chan int)
	inChan <- 1
	close(inChan)

	if val := <-outChan; val != 1 {
		t.Fatalf("Expected to receive 1 but got %v instead", val)
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel tot be closed but it's not")
	}
}

func TestScanMultipleValues(t *testing.T) {
	t.Parallel()

	inChan := make(chan int, 2)
	outChan := piper.Clone(inChan).Pipe(transform.Scan(scanAccumulator, 0)).Get().(chan int)
	inChan <- 1
	inChan <- 2
	close(inChan)

	if val1, val2 := <-outChan, <-outChan; val1 != 1 || val2 != 3 {
		t.Fatalf("Expected to receive 1 and 3 but got %v and %v instead", val1, val2)
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel tot be closed but it's not")
	}
}

func TestScanCloseWithValue(t *testing.T) {
	t.Parallel()

	closer := piper.PipeOperator{F: func(r piper.PipeResult, _ interface{}) (piper.PipeResult, interface{}) {
		return piper.PipeResult{Value: 1, IsValue: true, State: piper.Closed}, nil
	}}

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(closer, transform.Scan(scanAccumulator, 0)).Get().(chan int)
	inChan <- 1
	close(inChan)

	if val := <-outChan; val != 1 {
		t.Fatalf("Expected to receive 1 but got %v instead", val)
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel tot be closed but it's not")
	}
}
