package task

import (
	"encoding/json"
	"log"
	"regexp"

	"github.com/mnalsup/sentry/monitoring"
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

type taskJSON struct {
	Triggers []triggerJSON
	Actions  []Action
}

type triggerJSON struct {
	Source        string
	EventMatchers []string
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
func NewFromJSON(taskBytes []byte) (*Task, error) {
	var t taskJSON
	err := json.Unmarshal(taskBytes, &t)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var triggers = make([]Trigger, len(t.Triggers))
	for i, v := range t.Triggers {
		triggers[i] = triggerMustCompile(v.Source, v.EventMatchers)
	}
	var task = Task{triggers, t.Actions, INIT}
	return &task, nil
}

// MatchEvent receives and event and decides whether it matches a trigger
func (task *Task) MatchEvent(event *monitoring.Event) bool {
	log.Printf("Matching on event: %s %s.\n", event.Name, event.Source)
	for _, trigger := range task.Triggers {
		if trigger.Source == event.Source {
			for _, matcher := range trigger.EventMatchers {
				log.Printf("To: %s\n", matcher.String())
				if matcher.MatchString(event.Name) {
					return true
				}
			}
		}
	}
	return false
}

// HandleEvent handles an event
func (task *Task) HandleEvent(event *monitoring.Event) {
	if task.MatchEvent(event) {
		log.Println(event)
	}
}
