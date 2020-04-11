package utility_test

import (
	"testing"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/utility"
)

func TestTap(t *testing.T) {
	inChan := make(chan int, 1)

	logChan := make(chan interface{}, 1)
	piper.From(inChan).Pipe(utility.Tap(func(val interface{}) {
		logChan <- val
	}))

	inChan <- 1
	if val := <-logChan; val.(int) != 1 {
		t.Fatalf("Expected go get 1 but got %v", val)
	}

	close(inChan)
	close(logChan)
}
