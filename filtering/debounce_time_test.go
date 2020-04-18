package filtering_test

import (
	"testing"
	"time"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/filtering"
)

func TestDebouceTimeNoReset(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(filtering.DebounceTime(time.Microsecond * 500)).Get().(chan int)

	before := time.Now()
	inChan <- 1
	close(inChan)

	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*500); val != 1 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 1 after 500 Milliseconds with less than Milliseconds error, but got %v with a %v error", val, diff)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed")
	}
}

func TestDebouceTimeReset(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(filtering.DebounceTime(time.Millisecond * 500)).Get().(chan int)

	before := time.Now()
	inChan <- 1
	time.Sleep(time.Millisecond * 250)
	inChan <- 2
	close(inChan)

	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*750); val != 2 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 2 after 750 Milliseconds with less than Milliseconds error, but got %v with a %v error", val, diff)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed")
	}
}

func TestDebouceTimeMultipleNoReset(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(filtering.DebounceTime(time.Millisecond * 500)).Get().(chan int)

	before := time.Now()
	inChan <- 1

	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*500); val != 1 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 1 after 500 Milliseconds with less than Milliseconds error, but got %v with a %v error", val, diff)
	}

	before = time.Now()
	inChan <- 2
	close(inChan)

	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*500); val != 2 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 2 after 500 Milliseconds with less than Milliseconds error, but got %v with a %v error", val, diff)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed")
	}
}

func TestDebouceTimeMultipleReset(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(filtering.DebounceTime(time.Millisecond * 500)).Get().(chan int)

	before := time.Now()
	inChan <- 1
	time.Sleep(time.Millisecond * 250)
	inChan <- 2

	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*750); val != 2 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 2 after 750 Milliseconds with less than Milliseconds error, but got %v with a %v error", val, diff)
	}

	before = time.Now()
	inChan <- 2
	time.Sleep(time.Millisecond * 250)
	inChan <- 3
	close(inChan)

	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*750); val != 3 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 3 after 750 Milliseconds with less than Milliseconds error, but got %v with a %v error", val, diff)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed")
	}
}

func TestDebouceTimeNone(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(filtering.DebounceTime(time.Millisecond * 500)).Get().(chan int)
	close(inChan)

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed")
	}
}

func TestDebounceTimeCloseWithValue(t *testing.T) {
	t.Parallel()

	closer := piper.PipeOperator{F: func(r piper.PipeResult, _ interface{}) (piper.PipeResult, interface{}) {
		return piper.PipeResult{Value: 5, IsValue: true, State: piper.Closed}, nil
	}}

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(closer, filtering.DebounceTime(time.Millisecond*500)).Get().(chan int)

	before := time.Now()
	close(inChan)

	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*500); val != 5 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive value: 5 after about 500 Milliseconds with less than 50 Millisecond Error but got %v with an error of %v", val, diff)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}
