package filtering_test

import (
	"testing"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/filtering"
)

func TestLatestWithItems(t *testing.T) {
	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(filtering.Last()).Get().(chan int)

	inChan <- 1
	inChan <- 2
	inChan <- 3
	close(inChan)

	if val := <-outChan; val != 3 {
		t.Fatalf("Expected to receive 3, but got %v", val)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestLatestWithoutItems(t *testing.T) {
	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(filtering.Last()).Get().(chan int)

	close(inChan)

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}
