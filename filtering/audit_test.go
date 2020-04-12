package filtering_test

import (
	"testing"
	"time"

	"github.com/r-the-thinker/piper"
	"github.com/r-the-thinker/piper/filtering"
)

// convert 1 to 100 Milliseconds
func auditCalcFunc(val interface{}) time.Duration {
	return time.Millisecond * time.Duration(val.(int)) * 100
}

func TestAudit(t *testing.T) {
	t.Parallel()

	inChan := make(chan int, 2)
	outChan := piper.From(inChan).Pipe(filtering.Audit(auditCalcFunc)).Get().(chan int)

	before := time.Now()
	inChan <- 5
	inChan <- 2
	close(inChan)
	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*500); val != 2 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive value: 2 after about 500 Milliseconds with less than 50 Millisecond Error but got %v with an error of %v", val, diff)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestAuditMultiple(t *testing.T) {
	t.Parallel()

	inChan := make(chan int, 2)
	outChan := piper.From(inChan).Pipe(filtering.Audit(auditCalcFunc)).Get().(chan int)

	before := time.Now()
	inChan <- 5
	inChan <- 2
	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*500); val != 2 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive value: 2 after about 500 Milliseconds with less than 50 Millisecond Error but got %v with an error of %v", val, diff)
	}

	before = time.Now()
	inChan <- 3
	inChan <- 9
	close(inChan)
	if val, diff := <-outChan, time.Now().Sub(before)-(time.Millisecond*300); val != 9 || diff > time.Millisecond*50 {
		t.Fatalf("Expected to receive value: 2 after about 500 Milliseconds with less than 50 Millisecond Error but got %v with an error of %v", val, diff)
	}
	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed but it's not")
	}
}

func TestAuditNone(t *testing.T) {
	t.Parallel()

	inChan := make(chan int)
	outChan := piper.From(inChan).Pipe(filtering.Audit(auditCalcFunc)).Get().(chan int)
	close(inChan)

	if _, ok := <-outChan; ok {
		t.Fatal("Expected the output channel to be closed")
	}
}
