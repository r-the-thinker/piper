package transform_test

import (
	"testing"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/transform"
)

var bufferWhenFunc = func(val interface{}) bool {
	return val.(int) > 2
}

func TestBufferWhenNoValues(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(transform.BufferWhen(bufferWhenFunc, true), sliceToIntMapper).Get().(chan int)
	close(inChan)

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestBufferWhenEmitOnClose(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(transform.BufferWhen(bufferWhenFunc, true), sliceToIntMapper).Get().(chan int)
	// the one is not enough for the buffer to be emitted but the
	// emit on close flag will make sure that we get the 1
	inChan <- 1
	close(inChan)

	if val := <-outChan; val != 1 {
		t.Fatalf("Expected the output channel to receive 1 but got %v instead", val)
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestBufferWhenDontEmitOnClose(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(transform.BufferWhen(bufferWhenFunc, false), sliceToIntMapper).Get().(chan int)
	// the one is not enough for the buffer to be emitted but the
	// emit on close flag will make sure that we get the 1
	inChan <- 1
	close(inChan)

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestBufferWhenOne(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(transform.BufferWhen(bufferWhenFunc, true), sliceToIntMapper).Get().(chan int)
	// the one is not enough for the buffer to be emitted but the
	// emit on close flag will make sure that we get the 1
	inChan <- 1
	inChan <- 2
	inChan <- 3
	close(inChan)

	if val := <-outChan; val != 6 {
		t.Fatalf("Expected the output channel to receive 6 but got %v instead", val)
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestBufferWhenTwo(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(transform.BufferWhen(bufferWhenFunc, true), sliceToIntMapper).Get().(chan int)
	// the one is not enough for the buffer to be emitted but the
	// emit on close flag will make sure that we get the 1
	inChan <- 1
	inChan <- 2
	inChan <- 3

	if val := <-outChan; val != 6 {
		t.Fatalf("Expected the output channel to receive 6 but got %v instead", val)
	}

	inChan <- 1
	inChan <- 2
	inChan <- 3
	close(inChan)

	if val := <-outChan; val != 6 {
		t.Fatalf("Expected the output channel to receive 6 but got %v instead", val)
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestBufferWhenOneAndHalfEmitOnClose(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(transform.BufferWhen(bufferWhenFunc, true), sliceToIntMapper).Get().(chan int)
	// the one is not enough for the buffer to be emitted but the
	// emit on close flag will make sure that we get the 1
	inChan <- 1
	inChan <- 2
	inChan <- 3

	if val := <-outChan; val != 6 {
		t.Fatalf("Expected the output channel to receive 6 but got %v instead", val)
	}

	inChan <- 1
	inChan <- 2
	close(inChan)

	if val := <-outChan; val != 3 {
		t.Fatalf("Expected the output channel to receive 3 but got %v instead", val)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestBufferWhenOneAndHalfDontEmitOnClose(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(transform.BufferWhen(bufferWhenFunc, false), sliceToIntMapper).Get().(chan int)
	// the one is not enough for the buffer to be emitted but the
	// emit on close flag will make sure that we get the 1
	inChan <- 1
	inChan <- 2
	inChan <- 3

	if val := <-outChan; val != 6 {
		t.Fatalf("Expected the output channel to receive 6 but got %v instead", val)
	}

	close(inChan)
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}
