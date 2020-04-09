package transform_test

import (
	"testing"

	piper "github.com/r-the-thinker/piper"
	transform "github.com/r-the-thinker/piper/transform"
)

func TestReduce(t *testing.T) {
	t.Parallel()

	reducer := func(accumulator, currentVal interface{}) interface{} {
		return accumulator.(int) + currentVal.(int)
	}

	inputChan := make(chan int, 3)
	outputChan := piper.From(inputChan).Pipe(transform.Reduce(reducer, int(0))).Get().(chan int)

	inputChan <- 1
	inputChan <- 2
	inputChan <- 3
	close(inputChan)
	if v := <-outputChan; v != 6 {
		t.Fatalf("Expected to receive 2, but got %v", v)
	}
}
func TestReduceToInt(t *testing.T) {
	t.Parallel()

	reducer := func(accumulator int, currentVal interface{}) int {
		return accumulator + currentVal.(int)
	}

	inputChan := make(chan int, 3)
	outputChan := piper.From(inputChan).Pipe(transform.ReduceToInt(reducer, 0)).Get().(chan int)

	inputChan <- 1
	inputChan <- 2
	inputChan <- 3
	close(inputChan)
	if v := <-outputChan; v != 6 {
		t.Fatalf("Expected to receive 2, but got %v", v)
	}
}

func TestReduceFloat32(t *testing.T) {
	t.Parallel()

	reducer := func(accumulator float32, currentVal interface{}) float32 {
		return accumulator + currentVal.(float32)
	}

	inputChan := make(chan float32, 3)
	outputChan := piper.From(inputChan).Pipe(transform.ReduceToFloat32(reducer, 0)).Get().(chan float32)

	inputChan <- 1
	inputChan <- 2
	inputChan <- 3
	close(inputChan)
	if v := <-outputChan; v != 6 {
		t.Fatalf("Expected to receive 2, but got %v", v)
	}
}

func TestReduceToFloat64(t *testing.T) {
	t.Parallel()

	reducer := func(accumulator float64, currentVal interface{}) float64 {
		return accumulator + currentVal.(float64)
	}

	inputChan := make(chan float64, 3)
	outputChan := piper.From(inputChan).Pipe(transform.ReduceToFloat64(reducer, 0)).Get().(chan float64)

	inputChan <- 1
	inputChan <- 2
	inputChan <- 3
	close(inputChan)
	if v := <-outputChan; v != 6 {
		t.Fatalf("Expected to receive 2, but got %v", v)
	}
}

func TestReduceToString(t *testing.T) {
	t.Parallel()

	reducer := func(accumulator string, currentVal interface{}) string {
		// fmt.Println(accumulator, currentVal, accumulator+currentVal)
		return accumulator + currentVal.(string)
	}

	inputChan := make(chan string, 3)
	outputChan := piper.From(inputChan).Pipe(transform.ReduceToString(reducer, "")).Get().(chan string)

	inputChan <- "T"
	inputChan <- "E"
	inputChan <- "S"
	inputChan <- "T"
	close(inputChan)
	if v := <-outputChan; v != "TEST" {
		t.Fatalf("Expected to receive TEST, but got %v", v)
	}
}
