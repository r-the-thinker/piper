package transform

import (
	"container/list"

	"github.com/r-the-thinker/piper"
)

// BufferWhen collects values and asks the provided closer function if the
// buffer should be emitted. The buffer will include the value itself when
// emitted. Empty buffers are not seen as values and therefore not being emitted.
// emitOnClose determines if the buffer will be emitted regardless
// of the closer having decided that it is not yet time to do so.
func BufferWhen(closer func(interface{}) bool, emitOnClose bool) piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		l := storage.(*list.List)
		// only if emitOnClose is active the closed state will play
		// a role in emitting the buffer
		emit := emitOnClose && r.State == piper.Closed

		if r.IsValue {
			l.PushBack(r.Value)
			emit = closer(r.Value)
		}

		// if emit is set, the list will be transformed to an Slice and emitted if not then the value is
		// only appended to the list
		if emit {
			return piper.PipeResult{Value: listToSlice(l), IsValue: l.Len() > 0, State: r.State}, list.New()
		}
		return piper.PipeResult{IsValue: false, State: r.State}, l

	}, InitialStorage: list.New(), EventEmitter: nil}
}

func listToSlice(l *list.List) []interface{} {
	arr := make([]interface{}, l.Len())
	counter := 0
	for e := l.Front(); e != nil; e = e.Next() {
		arr[counter] = e.Value
		counter++
	}
	return arr
}
