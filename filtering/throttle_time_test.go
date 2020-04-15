package filtering_test

import (
	"testing"
	"time"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/filtering"
)

func TestThrottleTimeNoItems(t *testing.T) {
	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(filtering.ThrottleTime(time.Millisecond * 500)).Get().(chan int)

	before := time.Now()
	close(inChan)

	if _, ok := <-outChan; ok {
		t.Fatal("Expected output channel to be closed but it's not")
	}
	if diff := time.Now().Sub(before); diff > time.Millisecond*50 {
		t.Fatalf("Expected the output channel to be closed right away(less than 50ms error allowence) but it took %v", diff.Milliseconds())
	}
}

func TestThrottleTimeOneItem(t *testing.T) {
	inChan := make(chan int, 1)
	outChan := piper.From(inChan).Pipe(filtering.ThrottleTime(time.Millisecond * 500)).Get().(chan int)

	before := time.Now()
	inChan <- 1
	close(inChan)

	if val, diff := <-outChan, time.Now().Sub(before); val != 1 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 1 after no more than 50 milliseconds but got %v after %v", val, diff.Milliseconds())
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected output channel to be closed but it's not")
	}
	if diff := time.Now().Sub(before); diff > time.Millisecond*50 {
		t.Fatalf("Expected the output channel to be closed right away(less than 50ms error allowence) but it took %v", diff.Milliseconds())
	}
}

func TestThrottleTimeMultipleNoWait(t *testing.T) {
	inChan := make(chan int, 1)
	outChan := piper.From(inChan).Pipe(filtering.ThrottleTime(time.Millisecond * 500)).Get().(chan int)

	before := time.Now()
	inChan <- 1
	inChan <- 2
	close(inChan)

	if val, diff := <-outChan, time.Now().Sub(before); val != 1 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 1 after no more than 50 milliseconds but got %v after %v", val, diff.Milliseconds())
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected output channel to be closed but it's not")
	}
	if diff := time.Now().Sub(before); diff > time.Millisecond*50 {
		t.Fatalf("Expected the output channel to be closed right away(less than 50ms error allowence) but it took %v", diff.Milliseconds())
	}
}

func TestThrottleTimeMultipleWait(t *testing.T) {
	inChan := make(chan int, 2)
	outChan := piper.From(inChan).Pipe(filtering.ThrottleTime(time.Millisecond * 500)).Get().(chan int)

	before := time.Now()
	inChan <- 1
	if val, diff := <-outChan, time.Now().Sub(before); val != 1 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 1 after no more than 50 milliseconds but got %v after %v", val, diff.Milliseconds())
	}

	inChan <- 2
	time.Sleep(time.Millisecond * 550)
	before = time.Now()
	inChan <- 3
	close(inChan)

	if val, diff := <-outChan, time.Now().Sub(before); val != 3 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive 1 after no more than 50 milliseconds but got %v after %v", val, diff.Milliseconds())
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected output channel to be closed but it's not")
	}
	if diff := time.Now().Sub(before); diff > time.Millisecond*50 {
		t.Fatalf("Expected the output channel to be closed right away(less than 50ms error allowence) but it took %v", diff.Milliseconds())
	}
}
