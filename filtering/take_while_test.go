package filtering_test

import (
	"testing"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/filtering"
)

func TestTakeWhile(t *testing.T) {
	inChan := make(chan int, 2)

	predicate := func(val interface{}) bool {
		return val.(int) != 2
	}
	outChan := piper.From(inChan).Pipe(filtering.TakeWhile(predicate)).Get().(chan int)

	// Send in the values, that should not block because
	// anything including the 2 won't be emitted anymore
	inChan <- 0
	inChan <- 1
	inChan <- 2
	inChan <- 3

	if first, second := <-outChan, <-outChan; first != 0 || second != 1 {
		t.Fatalf("Expected to get 0 and 1 but got %v and %v", first, second)
	}
}
