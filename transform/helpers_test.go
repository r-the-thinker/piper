package transform_test

import (
	"fmt"

	"github.com/r-the-thinker/piper"
)

var sliceToIntMapper piper.PipeOperator = piper.PipeOperator{F: func(r piper.PipeResult, s interface{}) (piper.PipeResult, interface{}) {
	sum := 0
	fmt.Println(r.Value)
	for _, val := range r.Value.([]interface{}) {
		sum += val.(int)
	}
	return piper.PipeResult{Value: sum, IsValue: r.IsValue, State: r.State}, nil
}, InitialStorage: nil, EventEmitter: nil}
