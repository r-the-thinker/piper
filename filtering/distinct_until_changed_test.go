package filtering_test

import (
	"fmt"
	"testing"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/filtering"
)

func TestDistinctUntilChangedWithComperator(t *testing.T) {
	same := func(old, new interface{}) bool {
		return old.(int) == new.(int)
	}

	inChan := make(chan int, 3)
	outChan := piper.Clone(inChan).Pipe(filtering.DistinctUntilChanged(same)).Get().(chan int)

	inChan <- 1
	inChan <- 1
	inChan <- 2
	inChan <- 1
	close(inChan)

	if val1, val2, val3 := <-outChan, <-outChan, <-outChan; val1 != 1 || val2 != 2 || val3 != 1 {
		t.Fatalf("Expected to receive 1, 2 and 1 but got %v, %v and %v", val1, val2, val3)
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestDistinctUntilChangedWithoutComperator(t *testing.T) {
	inChan := make(chan int, 3)
	outChan := piper.Clone(inChan).Pipe(filtering.DistinctUntilChanged(nil)).Get().(chan int)

	inChan <- 1
	inChan <- 1
	inChan <- 2
	inChan <- 1
	close(inChan)

	if val1, val2, val3 := <-outChan, <-outChan, <-outChan; val1 != 1 || val2 != 2 || val3 != 1 {
		t.Fatalf("Expected to receive 1, 2 and 1 but got %v, %v and %v", val1, val2, val3)
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

type customType struct {
	name string
	age  int
}

func (c customType) String() string {
	return fmt.Sprintf("(%s, %d)", c.name, c.age)
}
func (c customType) same(o customType) bool {
	return c.name == o.name && c.age == o.age
}
func TestDistinctUntilChangedWithoutComperatorCustomType(t *testing.T) {

	inChan := make(chan customType, 3)
	outChan := piper.Clone(inChan).Pipe(filtering.DistinctUntilChanged(nil)).Get().(chan customType)

	inChan <- customType{"Robin", 22}
	inChan <- customType{"Robin", 22}
	inChan <- customType{"Nibor", 22}
	inChan <- customType{"Robin", 22}
	close(inChan)

	if val1, val2, val3 := <-outChan, <-outChan, <-outChan; !val1.same(customType{"Robin", 22}) ||
		!val2.same(customType{"Nibor", 22}) ||
		!val3.same(customType{"Robin", 22}) {
		t.Fatalf("Expected to receive (Robin, 22), (Nibor, 22) and (Robin, 22) but got %v, %v and %v", val1, val2, val3)
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}
