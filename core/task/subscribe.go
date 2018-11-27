package task

import (
	"github.com/mnalsup/sentry/monitoring"
)

// Subscribe subscribes a task list to an event stream
func Subscribe(tasks []*Task, events chan *monitoring.Event) {
	for e := range events {
		for _, task := range tasks {
			task.HandleEvent(e)
		}
	}
}
