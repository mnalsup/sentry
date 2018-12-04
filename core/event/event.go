package event

import "time"

// Event is exported to match triggers for tasks
type Event struct {
	Name     string
	Source   string
	Duration time.Duration
}
