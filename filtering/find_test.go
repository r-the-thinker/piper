package filtering_test

import (
	"testing"

	piper "github.com/r-the-thinker/piper"
	filtering "github.com/r-the-thinker/piper/filtering"
)

func TestFindOpen(t *testing.T) {
	inChan := make(chan int, 1)

	predicate := func(val interface{}) bool {
		return val.(int) > 2
	}
	outChan := piper.From(inChan).Pipe(filtering.Find(predicate)).Get().(chan int)

	// this should not block because we only get 3 and then nothing
	inChan <- 1
	inChan <- 2
	inChan <- 3
	inChan <- 4
	close(inChan)

	if val := <-outChan; val != 3 {
		t.Fatalf("Expected to receive 3 but got %v", val)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the channel to be closed")
	}
}

func TestFindClosed(t *testing.T) {
	inChan := make(chan int, 1)

	predicate := func(val interface{}) bool {
		return val.(int) > 2
	}
	outChan := piper.From(inChan).Pipe(filtering.Find(predicate)).Get().(chan int)

	// this should not block because we only get 3 and then nothing
	inChan <- 1
	inChan <- 2
	close(inChan)

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the channel to be closed")
	}
}
