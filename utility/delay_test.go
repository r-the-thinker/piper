package utility_test

import (
	"testing"
	"time"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/utility"
)

func TestDelayOne(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(utility.Delay(time.Millisecond * 500)).Get().(chan int)

	// send the value that we want to have delayed
	inChan <- 1
	inChan <- 2

	start := time.Now()
	<-outChan
	end := time.Now()
	close(inChan)

	if diff := end.Sub(start) - (time.Millisecond * 500); diff > time.Millisecond*50 {
		t.Fatalf("The second delay was not within a acceptable range. Diff: %v", diff)
	}
}

func TestDelayTwo(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(utility.Delay(time.Millisecond * 500)).Get().(chan int)

	// send the value that we want to have delayed
	inChan <- 1
	inChan <- 2

	start := time.Now()
	<-outChan
	middle := time.Now()
	<-outChan
	end := time.Now()
	close(inChan)

	if diff := middle.Sub(start) - (time.Millisecond * 500); diff > time.Millisecond*50 {
		t.Fatalf("The first delay was not within a acceptable range. Diff: %v", diff)
	}
	if diff := end.Sub(middle) - (time.Millisecond * 500); diff > time.Millisecond*50 {
		t.Fatalf("The second delay was not within a acceptable range. Diff: %v", diff)
	}
}
