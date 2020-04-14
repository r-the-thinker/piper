package filtering

import (
	"container/list"
	"fmt"

	"github.com/r-the-thinker/piper"
)

type takeLastStorage struct {
	items    *list.List
	numItems int
}

// Add to items and remove an item if there were to many
func (s *takeLastStorage) addToItems(maxItems int, item interface{}) {
	s.items.PushBack(item)
	s.numItems++
	if s.numItems > maxItems {
		s.items.Remove(s.items.Front())
	}
}

// TakeLast only emits the last >num< items from the stream. Where num
// is in range [1, MAX_INT). Keep in mind that >num< items need to be
// buffered because there is no way to tell when the stream closes
// The >num< values or less will be emitted only after the predecessor
// pipefunc closed.
func TakeLast(num int) piper.PipeOperator {
	if num < 1 {
		num = 1
	}

	comChan := make(chan takeLastStorage, 1)
	eventEmitter := make(chan interface{})
	// start the worker that waits for the com chan to send the items over
	go takeLastWorker(eventEmitter, comChan)

	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		if r.IsEvent {
			return r.Value.(piper.PipeResult), nil
		}

		// whenever it is not an event we are getting a value from the predecessor
		// if it is actually a value then add it to the items. When the predecessor closes
		// give the storage to the worker and always break because we only emit what the
		// worker tells us to emit.
		s := storage.(takeLastStorage)
		if r.IsValue {
			(&s).addToItems(num, r.Value)
		}
		if r.State == piper.Closed {
			comChan <- s
		}
		return piper.PipeResult{IsValue: false, State: piper.Break}, s
	}, InitialStorage: takeLastStorage{items: list.New(), numItems: 0}, EventEmitter: eventEmitter}
}

// the worker waits for the comChan to send over the storage
// and then emits all values of it as Pipe Results
func takeLastWorker(eventEmitter chan interface{}, comChan chan takeLastStorage) {
	// wait for the items to be sent over
	itemStore := <-comChan
	// send all values over
	for e := itemStore.items.Front(); e != nil; e = e.Next() {
		fmt.Println("Emitting", e.Value)
		eventEmitter <- piper.PipeResult{Value: e.Value, IsValue: true, State: piper.Open}
	}
	// tell the pipefunc to close
	eventEmitter <- piper.PipeResult{State: piper.Closed}
	// cleanup
	close(comChan)
	close(eventEmitter)
}
