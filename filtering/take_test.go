package filtering_test

import (
	"testing"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/filtering"
)

func TestTake(t *testing.T) {
	// We only need a buffer size of 2 because the third wont be emitted
	inChan := make(chan int, 2)
	outChan := piper.From(inChan).Pipe(filtering.Take(2)).Get().(chan int)

	inChan <- 0
	inChan <- 1
	inChan <- 2
	if first, second := <-outChan, <-outChan; first != 0 || second != 1 {
		t.Fatalf("Expected to receive 0 and 1 but got %v and %v", first, second)
	}

	close(inChan)
}

func TestTakeClosedAfter(t *testing.T) {
	// We only need a buffer size of 2 because the third wont be emitted
	inChan := make(chan int, 2)
	outChan := piper.From(inChan).Pipe(filtering.Take(2)).Get().(chan int)

	inChan <- 0
	inChan <- 1
	close(inChan)
	if first, second := <-outChan, <-outChan; first != 0 || second != 1 {
		t.Fatalf("Expected to receive 0 and 1 but got %v and %v", first, second)
	}
}

func TestTakeClosedBefore(t *testing.T) {
	// We only need a buffer size of 2 because the third wont be emitted
	inChan := make(chan int, 2)
	outChan := piper.From(inChan).Pipe(filtering.Take(2)).Get().(chan int)

	inChan <- 0
	close(inChan)
	if first := <-outChan; first != 0 {
		t.Fatalf("Expected to receive 0 but got %v", first)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}
