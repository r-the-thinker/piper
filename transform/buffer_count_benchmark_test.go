package transform_test

import (
	"testing"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/transform"
)

func benchmarkBufferCount(loadCount int, b *testing.B) {
	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(transform.BufferCount(loadCount), sliceToIntMapper).Get().(chan int)

	// run N rounds wich is defined by the benchmark
	for round := 1; round < b.N; round++ {
		for i := 0; i < loadCount; i++ {
			inChan <- round
		}
		if val := <-outChan; val != loadCount*round {
			b.Fatalf("Expected the sum of the inputs to be %v but it is %v", loadCount*round, val)
		}
	}

	close(inChan)
	if _, ok := <-outChan; ok {
		b.Fatal("Expected the output channel to be closed but it's not")
	}
}

func BenchmarkBufferCount1(b *testing.B) {
	benchmarkBufferCount(1, b)
}

func BenchmarkBufferCount5(b *testing.B) {
	benchmarkBufferCount(5, b)
}

func BenchmarkBufferCount1000(b *testing.B) {
	benchmarkBufferCount(1000, b)
}
