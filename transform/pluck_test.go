package transform_test

import (
	"testing"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/transform"
)

type pluckStruct struct {
	Test  int
	Other string
}

var pluckStructFiller = piper.PipeOperator{F: func(r piper.PipeResult, s interface{}) (piper.PipeResult, interface{}) {
	if r.IsValue {
		r.Value = pluckStruct{Test: r.Value.(int), Other: "test"}
	}
	return r, nil
}, InitialStorage: nil, EventEmitter: nil}

func TestPluckNoValues(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(pluckStructFiller, transform.Pluck("Test", true)).Get().(chan int)
	close(inChan)

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestPluckMultipleValues(t *testing.T) {
	t.Parallel()

	inChan := make(chan int, 2)
	outChan := piper.Clone(inChan).Pipe(pluckStructFiller, transform.Pluck("Test", true)).Get().(chan int)
	inChan <- 1
	inChan <- 2
	close(inChan)

	if val1, val2 := <-outChan, <-outChan; val1 != 1 || val2 != 2 {
		t.Fatalf("Expected to receive 1 and 2 but got %v and %v instead", val1, val2)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestPluckNotFound(t *testing.T) {
	t.Parallel()

	inChan := make(chan int, 2)
	outChan := piper.Clone(inChan).Pipe(pluckStructFiller, transform.Pluck("notfound", true)).Get().(chan int)
	inChan <- 1
	close(inChan)

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestPluckCloedWithValue(t *testing.T) {
	t.Parallel()

	closer := piper.PipeOperator{F: func(r piper.PipeResult, _ interface{}) (piper.PipeResult, interface{}) {
		return piper.PipeResult{Value: 1, IsValue: true, State: piper.Closed}, nil
	}}

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(closer, pluckStructFiller, transform.Pluck("Test", true)).Get().(chan int)
	inChan <- 1
	close(inChan)

	if val := <-outChan; val != 1 {
		t.Fatalf("Expected to receive 1 but got %v instead", val)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

// TODO possible? the panic happens in the other goroutine
func TestPluckPanicIfNotFound(t *testing.T) {}
