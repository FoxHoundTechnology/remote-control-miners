package domain

import (
	"time"
)

// TODO: validation
// TODO: create custom response objects for domain entities
type Alert struct {
	ID        string
	Name      string
	Location  string
	Timestamp time.Time // UTC
}
