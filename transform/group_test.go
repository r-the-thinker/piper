package transform_test

import (
	"fmt"
	"testing"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/transform"
)

func evenOddGrouper() ([]string, func(interface{}) int) {

	f := func(val interface{}) int {
		v, ok := val.(int)
		if !ok {
			return -1
		}
		if v%2 == 0 {
			return 0
		}
		return 1
	}

	return []string{"Even", "Odd"}, f
}

func evenOddChecker(evenCount, oddCount int) piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, s interface{}) (piper.PipeResult, interface{}) {
		groups := r.Value.(map[string][]interface{})
		fmt.Println(groups)
		ok := true
		for _, val := range groups["Even"] {
			ok = ok && val.(int)%2 == 0
		}
		for _, val := range groups["Odd"] {
			ok = ok && val.(int)%2 == 1
		}
		ok = ok && len(groups) == 2
		ok = ok && len(groups["Even"]) == evenCount
		ok = ok && len(groups["Odd"]) == oddCount

		var result int
		if ok {
			result = 1
		} else {
			result = 0
		}

		return piper.PipeResult{Value: result, IsValue: r.IsValue, State: r.State}, nil
	}, InitialStorage: nil, EventEmitter: nil}
}

func TestGroupNone(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(transform.Group(evenOddGrouper()), evenOddChecker(0, 0)).Get().(chan int)
	close(inChan)

	if isOk := <-outChan; isOk != 1 {
		t.Fatal("Expected even and odd to be sorted correctly")
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestGroupEvenOddOne(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(transform.Group(evenOddGrouper()), evenOddChecker(1, 1)).Get().(chan int)

	inChan <- 1
	inChan <- 2
	close(inChan)

	if isOk := <-outChan; isOk != 1 {
		t.Fatal("Expected even and odd to be sorted correctly")
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestGroupEvenOddMultipleSame(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(transform.Group(evenOddGrouper()), evenOddChecker(2, 2)).Get().(chan int)

	inChan <- 1
	inChan <- 2
	inChan <- 3
	inChan <- 4
	close(inChan)

	if isOk := <-outChan; isOk != 1 {
		t.Fatal("Expected even and odd to be sorted correctly")
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestGroupEvenOddMultipleDifferent(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(transform.Group(evenOddGrouper()), evenOddChecker(1, 2)).Get().(chan int)

	inChan <- 1
	inChan <- 2
	inChan <- 3
	close(inChan)

	if isOk := <-outChan; isOk != 1 {
		t.Fatal("Expected even and odd to be sorted correctly")
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestGroupEvenOddOneCloseWithValue(t *testing.T) {
	t.Parallel()

	closer := piper.PipeOperator{F: func(r piper.PipeResult, _ interface{}) (piper.PipeResult, interface{}) {
		return piper.PipeResult{Value: 1, IsValue: true, State: piper.Closed}, nil
	}}

	inChan := make(chan int)
	outChan := piper.Clone(inChan).Pipe(closer, transform.Group(evenOddGrouper()), evenOddChecker(0, 1)).Get().(chan int)

	inChan <- 1
	close(inChan)

	if isOk := <-outChan; isOk != 1 {
		t.Fatal("Expected even and odd to be sorted correctly")
	}

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestGroupNilLabels(t *testing.T) {
	t.Parallel()

	defer func() {
		if res := recover(); res == nil {
			t.Fatal("Expected Group to have paniced without the labels")
		}
	}()

	inChan := make(chan int)
	piper.Clone(inChan).Pipe(transform.Group(nil, func(interface{}) int { return 1 }), evenOddChecker(0, 1))
	close(inChan)
}

func TestGroupNilSelector(t *testing.T) {
	t.Parallel()

	defer func() {
		if res := recover(); res == nil {
			t.Fatal("Expected Group to have paniced without the the selector")
		}
	}()

	inChan := make(chan int)
	piper.Clone(inChan).Pipe(transform.Group([]string{"Bla"}, nil), evenOddChecker(0, 1))
	close(inChan)
}

func TestGroupNilLabelsAndSelector(t *testing.T) {
	t.Parallel()

	defer func() {
		if res := recover(); res == nil {
			t.Fatal("Expected Group to have paniced without the the labels and the selector")
		}
	}()

	inChan := make(chan int)
	piper.Clone(inChan).Pipe(transform.Group(nil, nil), evenOddChecker(0, 1))
	close(inChan)
}
