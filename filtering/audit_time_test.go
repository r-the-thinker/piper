package filtering_test

import (
	"testing"
	"time"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/filtering"
)

func TestAuditTime(t *testing.T) {
	t.Parallel()

	inChan := make(chan int, 2)
	outChan := piper.From(inChan).Pipe(filtering.AuditTime(time.Millisecond * 500)).Get().(chan int)

	before := time.Now()
	inChan <- 1
	inChan <- 2
	inChan <- 3
	after := time.Now()
	close(inChan)

	// check time difference
	if diff := after.Sub(before) - (time.Millisecond * 500); diff > time.Millisecond*50 {
		t.Fatalf("The emission was not in time. It took %v. Expected was 500", diff)
	}

	// check that there is only one emission, the last one
	if val := <-outChan; val != 3 {
		t.Fatalf("Expected to receive 3 but got %v", val)
	}

	// check that the output channel closed
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed, but it's not.")
	}
}

func TestAuditTimeMultiple(t *testing.T) {
	t.Parallel()

	inChan := make(chan int, 2)
	outChan := piper.From(inChan).Pipe(filtering.AuditTime(time.Millisecond * 500)).Get().(chan int)

	before := time.Now()
	inChan <- 1
	time.Sleep(time.Millisecond * 550)
	inChan <- 2
	after := time.Now()
	close(inChan)

	// check time difference (wait for auditTime twice plus the active sleep)
	if diff := after.Sub(before) - (time.Millisecond * (2*500 + 550)); diff > time.Millisecond*50 {
		t.Fatalf("The emission was not in time. It took %v. Expected was 500", diff)
	}

	// check that there is only one emission, the last one
	if val1, val2 := <-outChan, <-outChan; val1 != 1 || val2 != 2 {
		t.Fatalf("Expected to receive 1 and 2 but got %v and %v", val1, val2)
	}

	// check that the output channel closed
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed, but it's not.")
	}
}

func TestAuditTimeNone(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(filtering.AuditTime(time.Millisecond * 1000)).Get().(chan int)
	close(inChan)

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestAuditTimeCloseWithValue(t *testing.T) {
	t.Parallel()

	closer := piper.PipeOperator{F: func(r piper.PipeResult, _ interface{}) (piper.PipeResult, interface{}) {
		return piper.PipeResult{Value: 5, IsValue: true, State: piper.Closed}, nil
	}}

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(closer, filtering.AuditTime(time.Millisecond*500)).Get().(chan int)

	before := time.Now()
	close(inChan)

	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*500); val != 5 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive value: 5 after about 500 Milliseconds with less than 50 Millisecond Error but got %v with an error of %v", val, diff)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}
