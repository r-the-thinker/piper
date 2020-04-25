package transform

import (
	"fmt"

	"github.com/r-the-thinker/piper"
)

// Group emits once the predecessor closed a map of slices that has a label->slice
// mapping. A group without values will have a nil slice not a empty one len=0
// regardless.
func Group(labels []string, selector func(interface{}) int) piper.PipeOperator {
	if len(labels) == 0 {
		panic("There were no group labels provided")
	}
	if selector == nil {
		panic("The selector is nil but a function is expected")
	}

	groups := make(map[string][]interface{})
	for _, label := range labels {
		groups[label] = nil
	}

	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		groups := storage.(map[string][]interface{})
		if r.IsValue {
			slot := selector(r.Value)
			if slot >= len(labels) {
				panic(fmt.Sprintf("The selector chose %v but there are only %v slots", slot, len(labels)))
			}

			groups[labels[slot]] = append(groups[labels[slot]], r.Value)
		}

		if r.State != piper.Closed {
			return piper.PipeResult{State: piper.Break}, groups
		}
		return piper.PipeResult{Value: groups, IsValue: len(groups) > 0, State: piper.Closed}, nil
	}, InitialStorage: groups, EventEmitter: nil}
}
