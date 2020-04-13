package filtering

import (
	"time"

	"github.com/r-the-thinker/piper"
)

// DebounceTime emits a value after the specified duration. Every new value that comes
// in, resets the timer and replaces the old value, therefore only if there arn't emissions
// in that duration the value will be emitted
func DebounceTime(dur time.Duration) piper.PipeOperator {
	return Debounce(func(_ interface{}) time.Duration { return dur })
}
