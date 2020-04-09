package transform_test

import (
	"testing"

	piper "github.com/r-the-thinker/piper"
	transform "github.com/r-the-thinker/piper/transform"
)

func TestMapOpen(t *testing.T) {
	t.Parallel()

	mapper := func(v interface{}) interface{} {
		return 2
	}

	inputChan := make(chan int)
	outputChan := piper.From(inputChan).Pipe(transform.Map(mapper)).Get().(chan int)

	inputChan <- 1
	if v := <-outputChan; v != 2 {
		t.Fatalf("Expected to receive 2, but got %v", v)
	}

	close(inputChan)
}

func TestMapClosed(t *testing.T) {
	t.Parallel()

	mapper := func(v interface{}) interface{} {
		return 2
	}

	inputChan := make(chan int)
	outputChan := piper.From(inputChan).Pipe(transform.Map(mapper)).Get().(chan int)

	close(inputChan)
	if _, ok := <-outputChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}

}

/*
func TestMapInt(t *testing.T) {
	t.Parallel()

	mapper := func(v int) int {
		return 2
	}

	inputChan := make(chan int)
	outputChan := piper.From(inputChan).Pipe(transform.MapInt(mapper)).Get().(chan int)

	inputChan <- 1
	if v := <-outputChan; v != 2 {
		t.Fatalf("Expected to receive 2, but got %v", v)
	}

	close(inputChan)
}

func TestMapInt8(t *testing.T) {
	t.Parallel()

	mapper := func(v int8) int8 {
		return 2
	}

	inputChan := make(chan int8)
	outputChan := piper.From(inputChan).Pipe(transform.MapInt8(mapper)).Get().(chan int8)

	inputChan <- 1
	if v := <-outputChan; v != 2 {
		t.Fatalf("Expected to receive 2, but got %v", v)
	}

	close(inputChan)
}

func TestMapInt16(t *testing.T) {
	t.Parallel()

	mapper := func(v int16) int16 {
		return 2
	}

	inputChan := make(chan int16)
	outputChan := piper.From(inputChan).Pipe(transform.MapInt16(mapper)).Get().(chan int16)

	inputChan <- 1
	if v := <-outputChan; v != 2 {
		t.Fatalf("Expected to receive 2, but got %v", v)
	}

	close(inputChan)
}

func TestMapInt32(t *testing.T) {
	t.Parallel()

	mapper := func(v int32) int32 {
		return 2
	}

	inputChan := make(chan int32)
	outputChan := piper.From(inputChan).Pipe(transform.MapInt32(mapper)).Get().(chan int32)

	inputChan <- 1
	if v := <-outputChan; v != 2 {
		t.Fatalf("Expected to receive 2, but got %v", v)
	}

	close(inputChan)
}

func TestMapInt64(t *testing.T) {
	t.Parallel()

	mapper := func(v int64) int64 {
		return 2
	}

	inputChan := make(chan int64)
	outputChan := piper.From(inputChan).Pipe(transform.MapInt64(mapper)).Get().(chan int64)

	inputChan <- 1
	if v := <-outputChan; v != 2 {
		t.Fatalf("Expected to receive 2, but got %v", v)
	}

	close(inputChan)
}

func TestMapFloat32(t *testing.T) {
	t.Parallel()

	mapper := func(v float32) float32 {
		return 2
	}

	inputChan := make(chan float32)
	outputChan := piper.From(inputChan).Pipe(transform.MapFloat32(mapper)).Get().(chan float32)

	inputChan <- 1
	if v := <-outputChan; v != 2.0 {
		t.Fatalf("Expected to receive 2, but got %v", v)
	}

	close(inputChan)
}

func TestMapFloat64(t *testing.T) {
	t.Parallel()

	mapper := func(v float64) float64 {
		return 2
	}

	inputChan := make(chan float64)
	outputChan := piper.From(inputChan).Pipe(transform.MapFloat64(mapper)).Get().(chan float64)

	inputChan <- 1
	if v := <-outputChan; v != 2.0 {
		t.Fatalf("Expected to receive 2, but got %v", v)
	}

	close(inputChan)
}

func TestMapString(t *testing.T) {
	t.Parallel()

	mapper := func(v string) string {
		return "2"
	}

	inputChan := make(chan string)
	outputChan := piper.From(inputChan).Pipe(transform.MapString(mapper)).Get().(chan string)

	inputChan <- "1"
	if v := <-outputChan; v != "2" {
		t.Fatalf("Expected to receive 2, but got %v", v)
	}

	close(inputChan)
}
*/
