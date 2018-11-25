package task

import (
	"encoding/json"
	"log"
	"regexp"

	"github.com/mnalsup/sentry/monitoring"
	"github.com/mnalsup/sentry/core"
)

// INIT is the state before actions begin or triggerQ
// ONGOING is the state when actions are started but not finished
// POST is the state after actions stop
const (
	INIT = iota
	ONGOING
	POST
)

// Task is a set of triggers and actions and state information
type Task struct {
	Triggers []Trigger
	Actions  []Action
	Status   int
}

// Trigger is a definition of a trigger to watch
type Trigger struct {
	Source        string
	EventMatchers []*regexp.Regexp
}

// Action is a definition of an action to take
type Action struct {
	sourceName string
	action     string
}

// MustCompile creates an Event and panics if the expression does not compile
func triggerMustCompile(source string, exprs []string) Trigger {
	length := len(exprs)
	matchers := make([]*regexp.Regexp, length)
	for i, v := range exprs {
		matchers[i] = regexp.MustCompile(v)
	}
	return Trigger{source, matchers}
}

// New returns a new Task
func New(taskJSON []byte) (*Task, error) {
	var task Task
	err := json.Unmarshal(taskJSON, &task)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &task, nil
}

// TestTrigger receives and event and decides whether it matches a trigger
func (task *Task) TestTrigger(event *monitoring.Event) bool {
	for _, trigger := range task.Triggers {
		if trigger.Source == event.Source {
			for _, matcher := range trigger.EventMatchers {
				if matcher.MatchString(event.Name) {
					return true
				}
			}
		}
	}
	return false
}
