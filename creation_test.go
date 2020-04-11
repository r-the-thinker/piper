package piper_test

import (
	"reflect"
	"testing"

	"github.com/r-the-thinker/piper"
)

// Method: From
// Knowledge: if this this test passes, then we know that from creates a channel of the correct data type
func TestFromType(t *testing.T) {
	inChan := make(chan int)
	outChan := piper.From(inChan).Get().(chan int)
	if c := cap(outChan); c != 0 {
		t.Fatalf("The capacity of the channel should have been 0, but it is %v", c)
	}
	if cType := reflect.TypeOf(outChan).Elem().Kind(); cType != reflect.Int {
		t.Fatalf("From created a channel of type %v, but expected it to be of type int", cType)
	}

	inChan = make(chan int, 10)
	outChan = piper.From(inChan).Get().(chan int)
	if c := cap(outChan); c != 10 {
		t.Fatalf("The capacity of the channel should have been 10, but it is %v", c)
	}
	if cType := reflect.TypeOf(outChan).Elem().Kind(); cType != reflect.Int {
		t.Fatalf("From created a channel of type %v, but expected it to be of type int", cType)
	}
	close(inChan)
}

// Method: From
// Knowledge: if this tests passes we know that From creates a working unbuffered channel
func TestFromUnbuffered(t *testing.T) {
	inChan := make(chan int)
	outChan := piper.From(inChan).Get().(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			inChan <- i
		}
	}()

	for i := 0; i < 10; i++ {
		if val := <-outChan; i != val {
			t.Fatalf("Expected %v received %v", i, val)
		}
	}

	close(inChan)
}

// Method: From
// Knowledge: if this tests passes we know that From creates a working buffered channel
func TestFromBuffered(t *testing.T) {
	inChan := make(chan int, 10)
	outChan := piper.From(inChan).Get().(chan int)

	for i := 0; i < 10; i++ {
		inChan <- i
	}

	for i := 0; i < 10; i++ {
		if val := <-outChan; i != val {
			t.Fatalf("Expected %v received %v", i, val)
		}
	}

	close(inChan)
}
