package filtering

import (
	"time"

	"github.com/r-the-thinker/piper"
)

// ThrottleTime emits the value upon receiving it and then ignores
// all following value for the specified duration
func ThrottleTime(dur time.Duration) piper.PipeOperator {
	return Throttle(func(_ interface{}) time.Duration {
		return dur
	})
}
