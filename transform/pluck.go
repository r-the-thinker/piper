package transform

import (
	"fmt"
	"reflect"

	"github.com/r-the-thinker/piper"
)

// Pluck extracts the value to a key from a struct(Exported fields only). If the provided value
// is not a struct this will panic. This behaviour is caused by the reflect package
// and intentionally not caught. Essentially it helps finding errors because of the
// type uncertenty inside the pipe.
// If the value cannot be found, then it is regarded as no value if andFilter is true.
// -label the key of the value to extract
// -andFilter if a struct does not have the key filter it out. If set to false
// 	it will panic when the specified key is not present
func Pluck(label string, andFilter bool) piper.PipeOperator {
	return piper.PipeOperator{F: func(r piper.PipeResult, storage interface{}) (piper.PipeResult, interface{}) {
		if r.IsValue {
			// try getting the field of the struct. this might panic but as explained
			// above, this is intended. The value is only regarded as IsValue=true when
			// it is not the zero value of reflect.Value
			item := reflect.ValueOf(r.Value).FieldByName(label)
			valid := item.IsValid()
			if !andFilter && !valid {
				panic(fmt.Sprintf("Could not find key in the struct %v", r.Value))
			}
			if valid {
				return piper.PipeResult{Value: item.Interface(), IsValue: valid, State: r.State}, nil
			}
			return piper.PipeResult{IsValue: false, State: r.State}, nil
		}
		return r, nil
	}, InitialStorage: nil, EventEmitter: nil}
}
