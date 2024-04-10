package alert

import (
	"time"
)

// TODO: Validation
type Alert struct {
	ID        string
	Name      string
	Location  string
	Timestamp time.Time // UTC
}
