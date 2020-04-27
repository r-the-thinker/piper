package transform_test

import (
	"testing"
	"time"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/transform"
)

func TestBufferTimeNone(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(transform.BufferTime(time.Millisecond * 500)).Get().(chan int)
	before := time.Now()
	close(inChan)

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}

	if diff := time.Now().Sub(before) - (time.Millisecond * 500); diff > time.Millisecond*50 {
		t.Fatalf("Expected the output channel to be closed right away with an error of 50ms but it took %v", diff.Milliseconds())
	}
}

func TestBufferTimeOneValueinDuration(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(transform.BufferTime(time.Millisecond*500), sliceToIntMapper).Get().(chan int)

	before := time.Now()
	inChan <- 1
	close(inChan)

	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*500); val != 1 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 1 in less than 50ms but got %v with an error of %v", val, diff.Milliseconds())
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
	if diff := time.Now().Sub(before) - (time.Millisecond * 500); diff > time.Millisecond*50 {
		t.Fatalf("Expected the output channel to be closed right away with an error of 50ms but it took %v", diff.Milliseconds())
	}
}

func TestBufferTimeMultipleValuesInDuration(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(transform.BufferTime(time.Millisecond*500), sliceToIntMapper).Get().(chan int)

	before := time.Now()
	inChan <- 1
	inChan <- 2
	close(inChan)

	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*500); val != 3 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 3 in less than 50ms but got %v with an error of %v", val, diff.Milliseconds())
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
	if diff := time.Now().Sub(before) - (time.Millisecond * 500); diff > time.Millisecond*50 {
		t.Fatalf("Expected the output channel to be closed right away with an error of 50ms but it took %v", diff.Milliseconds())
	}
}

func TestBufferTimeOneValueOutsideDuration(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(transform.BufferTime(time.Millisecond*500), sliceToIntMapper).Get().(chan int)

	before := time.Now()
	inChan <- 1
	time.Sleep(time.Millisecond * 500)
	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*500); val != 1 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 1 in less than 50ms but got %v with an error of %v", val, diff.Milliseconds())
	}

	before = time.Now()
	inChan <- 2
	close(inChan)

	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*500); val != 2 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 2 in less than 50ms but got %v with an error of %v", val, diff.Milliseconds())
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
	if diff := time.Now().Sub(before) - (time.Millisecond * 500); diff > time.Millisecond*50 {
		t.Fatalf("Expected the output channel to be closed right away with an error of 50ms but it took %v", diff.Milliseconds())
	}
}

func TestBufferTimeMultipleValuesOutsideDuration(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(transform.BufferTime(time.Millisecond*500), sliceToIntMapper).Get().(chan int)

	before := time.Now()
	inChan <- 1
	inChan <- 2
	time.Sleep(time.Millisecond * 500)
	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*500); val != 3 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 3 in less than 50ms but got %v with an error of %v", val, diff.Milliseconds())
	}

	before = time.Now()
	inChan <- 2
	inChan <- 3
	close(inChan)

	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*500); val != 5 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 5 in less than 50ms but got %v with an error of %v", val, diff.Milliseconds())
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
	if diff := time.Now().Sub(before) - (time.Millisecond * 500); diff > time.Millisecond*50 {
		t.Fatalf("Expected the output channel to be closed right away with an error of 50ms but it took %v", diff.Milliseconds())
	}
}

func TestBufferTimeClosedWithValue(t *testing.T) {
	t.Parallel()

	closer := piper.PipeOperator{F: func(r piper.PipeResult, _ interface{}) (piper.PipeResult, interface{}) {
		return piper.PipeResult{Value: 1, IsValue: true, State: piper.Closed}, nil
	}}

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(closer, transform.BufferTime(time.Millisecond*500), sliceToIntMapper).Get().(chan int)

	before := time.Now()
	inChan <- 1
	close(inChan)

	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*500); val != 1 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 1 in less than 50ms but got %v with an error of %v", val, diff.Milliseconds())
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
	if diff := time.Now().Sub(before) - (time.Millisecond * 500); diff > time.Millisecond*50 {
		t.Fatalf("Expected the output channel to be closed right away with an error of 50ms but it took %v", diff.Milliseconds())
	}
}
