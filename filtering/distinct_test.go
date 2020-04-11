package filtering_test

import (
	"testing"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/filtering"
)

func TestDistinct(t *testing.T) {
	inChan := make(chan int, 4)
	outChan := piper.From(inChan).Pipe(filtering.Distinct()).Get().(chan int)

	inChan <- 1
	inChan <- 1
	inChan <- 2
	inChan <- 1
	close(inChan)

	if val1, val2 := <-outChan, <-outChan; val1 != 1 || val2 != 2 {
		t.Fatalf("Expected to receive 1 and 2 but got %v and %v", val1, val2)
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}
