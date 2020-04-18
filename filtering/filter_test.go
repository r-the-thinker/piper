package filtering_test

import (
	"strings"
	"testing"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/filtering"
)

func TestFilter(t *testing.T) {
	t.Parallel()

	filterer := func(val interface{}) bool {
		return val.(int) >= 1
	}

	// This needs only size one so that we can send in val 1 without blocking
	inChan := make(chan int, 1)
	outChan := piper.From(inChan).Pipe(filtering.Filter(filterer)).Get().(chan int)

	inChan <- 0
	inChan <- 1
	if val := <-outChan; val != 1 {
		t.Fatalf("Expected to get 1 from the filtered channel but instead got %v", val)
	}

	close(inChan)
}

func TestFilterInt(t *testing.T) {
	t.Parallel()

	filterer := func(val int) bool {
		return val >= 1
	}

	// This needs only size one so that we can send in val 1 without blocking
	inChan := make(chan int, 1)
	outChan := piper.From(inChan).Pipe(filtering.FilterInt(filterer)).Get().(chan int)

	inChan <- 0
	inChan <- 1
	if val := <-outChan; val != 1 {
		t.Fatalf("Expected to get 1 from the filtered channel but instead got %v", val)
	}

	close(inChan)
}

func TestFilterInt8(t *testing.T) {
	t.Parallel()

	filterer := func(val int8) bool {
		return val >= 1
	}

	// This needs only size one so that we can send in val 1 without blocking
	inChan := make(chan int8, 1)
	outChan := piper.From(inChan).Pipe(filtering.FilterInt8(filterer)).Get().(chan int8)

	inChan <- 0
	inChan <- 1
	if val := <-outChan; val != 1 {
		t.Fatalf("Expected to get 1 from the filtered channel but instead got %v", val)
	}

	close(inChan)
}

func TestFilterInt16(t *testing.T) {
	t.Parallel()

	filterer := func(val int16) bool {
		return val >= 1
	}

	// This needs only size one so that we can send in val 1 without blocking
	inChan := make(chan int16, 1)
	outChan := piper.From(inChan).Pipe(filtering.FilterInt16(filterer)).Get().(chan int16)

	inChan <- 0
	inChan <- 1
	if val := <-outChan; val != 1 {
		t.Fatalf("Expected to get 1 from the filtered channel but instead got %v", val)
	}

	close(inChan)
}

func TestFilterInt32(t *testing.T) {
	t.Parallel()

	filterer := func(val int32) bool {
		return val >= 1
	}

	// This needs only size one so that we can send in val 1 without blocking
	inChan := make(chan int32, 1)
	outChan := piper.From(inChan).Pipe(filtering.FilterInt32(filterer)).Get().(chan int32)

	inChan <- 0
	inChan <- 1
	if val := <-outChan; val != 1 {
		t.Fatalf("Expected to get 1 from the filtered channel but instead got %v", val)
	}

	close(inChan)
}

func TestFilterInt64(t *testing.T) {
	t.Parallel()

	filterer := func(val int64) bool {
		return val >= 1
	}

	// This needs only size one so that we can send in val 1 without blocking
	inChan := make(chan int64, 1)
	outChan := piper.From(inChan).Pipe(filtering.FilterInt64(filterer)).Get().(chan int64)

	inChan <- 0
	inChan <- 1
	if val := <-outChan; val != 1 {
		t.Fatalf("Expected to get 1 from the filtered channel but instead got %v", val)
	}

	close(inChan)
}

func TestFilterFloat32(t *testing.T) {
	t.Parallel()

	filterer := func(val float32) bool {
		return val >= 1
	}

	// This needs only size one so that we can send in val 1 without blocking
	inChan := make(chan float32, 1)
	outChan := piper.From(inChan).Pipe(filtering.FilterFloat32(filterer)).Get().(chan float32)

	inChan <- 0
	inChan <- 1
	if val := <-outChan; val != 1 {
		t.Fatalf("Expected to get 1 from the filtered channel but instead got %v", val)
	}

	close(inChan)
}

func TestFilterFloat64(t *testing.T) {
	t.Parallel()

	filterer := func(val float64) bool {
		return val >= 1
	}

	// This needs only size one so that we can send in val 1 without blocking
	inChan := make(chan float64, 1)
	outChan := piper.From(inChan).Pipe(filtering.FilterFloat64(filterer)).Get().(chan float64)

	inChan <- 0
	inChan <- 1
	if val := <-outChan; val != 1 {
		t.Fatalf("Expected to get 1 from the filtered channel but instead got %v", val)
	}

	close(inChan)
}
func TestFilterString(t *testing.T) {
	t.Parallel()

	filterer := func(val string) bool {
		// test if val > b
		return strings.HasPrefix(val, "Tes")
	}

	// This needs only size one so that we can send in val 1 without blocking
	inChan := make(chan string, 1)
	outChan := piper.From(inChan).Pipe(filtering.FilterString(filterer)).Get().(chan string)

	inChan <- "a"
	inChan <- "Test"
	if val := <-outChan; val != "Test" {
		t.Fatalf("Expected to get >Test< from the filtered channel but instead got %v", val)
	}

	close(inChan)
}

func TestFilterBool(t *testing.T) {
	t.Parallel()

	filterer := func(val bool) bool {
		return !val
	}

	// This needs only size one so that we can send in val 1 without blocking
	inChan := make(chan bool, 1)
	outChan := piper.From(inChan).Pipe(filtering.FilterBool(filterer)).Get().(chan bool)

	inChan <- true
	inChan <- false
	if val := <-outChan; val {
		t.Fatal("Expected to get false from the filtered channel but instead got true")
	}

	close(inChan)
}

func TestFilterWithClosedAllowed(t *testing.T) {
	closer := piper.PipeOperator{F: func(r piper.PipeResult, _ interface{}) (piper.PipeResult, interface{}) {
		return piper.PipeResult{Value: 5, IsValue: true, State: piper.Closed}, nil
	}}
	filterer := func(val interface{}) bool {
		return val.(int) >= 1
	}

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(closer, filtering.Filter(filterer)).Get().(chan int)

	close(inChan)

	if val := <-outChan; val != 5 {
		t.Fatalf("Expected to receive 5 but got %v", val)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestFilterWithClosedNotAllowed(t *testing.T) {
	closer := piper.PipeOperator{F: func(r piper.PipeResult, _ interface{}) (piper.PipeResult, interface{}) {
		return piper.PipeResult{Value: 5, IsValue: true, State: piper.Closed}, nil
	}}
	filterer := func(val interface{}) bool {
		return val.(int) >= 11
	}

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(closer, filtering.Filter(filterer)).Get().(chan int)

	close(inChan)
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}
