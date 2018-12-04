package task

import (
	"github.com/mnalsup/sentry/core/event"
)

// Subscribe subscribes a task list to an event stream
func Subscribe(tasks []*Task, events chan *event.Event) {
	for e := range events {
		for _, task := range tasks {
			task.HandleEvent(e)
		}
	}
}
