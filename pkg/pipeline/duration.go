package pipeline

import (
	"encoding/json"
	"time"
)

type Duration struct {
	time.Duration
}

// Send out in nanoseconds for cross-compatability
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Milliseconds())
}
