package filtering_test

import (
	"testing"
	"time"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/filtering"
)

func debounceTimeCalc(val interface{}) time.Duration {
	return time.Duration(val.(int)) * time.Millisecond * 100
}

func TestDebouceNoReset(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(filtering.Debounce(debounceTimeCalc)).Get().(chan int)

	before := time.Now()
	inChan <- 5
	close(inChan)

	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*500); val != 5 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 5 after 500 Milliseconds with less than Milliseconds error, but got %v with a %v error", val, diff)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed")
	}
}

func TestDebouceReset(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(filtering.Debounce(debounceTimeCalc)).Get().(chan int)

	before := time.Now()
	inChan <- 5
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

func TestDebouceMultipleNoReset(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(filtering.Debounce(debounceTimeCalc)).Get().(chan int)

	before := time.Now()
	inChan <- 5

	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*500); val != 5 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 5 after 500 Milliseconds with less than Milliseconds error, but got %v with a %v error", val, diff)
	}

	before = time.Now()
	inChan <- 3
	close(inChan)

	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*500); val != 3 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 3 after 300 Milliseconds with less than Milliseconds error, but got %v with a %v error", val, diff)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed")
	}
}

func TestDebouceMultipleReset(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(filtering.Debounce(debounceTimeCalc)).Get().(chan int)

	before := time.Now()
	inChan <- 5
	time.Sleep(time.Millisecond * 250)
	inChan <- 2

	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*750); val != 2 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 2 after 750 Milliseconds with less than Milliseconds error, but got %v with a %v error", val, diff)
	}

	before = time.Now()
	inChan <- 3
	time.Sleep(time.Millisecond * 200)
	inChan <- 1
	close(inChan)

	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*500); val != 1 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 1 after 500 Milliseconds with less than Milliseconds error, but got %v with a %v error", val, diff)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed")
	}
}

func TestDebouceNone(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(filtering.Debounce(debounceTimeCalc)).Get().(chan int)
	close(inChan)

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed")
	}
}
