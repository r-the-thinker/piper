package filtering

import (
	"time"

	"github.com/r-the-thinker/piper"
)

// AuditTime waits the specified duration upon receiving an element.
// It will output the last value it received during that time.
// dur - the duration to wait
func AuditTime(dur time.Duration) piper.PipeOperator {
	return Audit(func(_ interface{}) time.Duration {
		return dur
	})
}
