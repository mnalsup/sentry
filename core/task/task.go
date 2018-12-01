package task

import (
	"encoding/json"
	"regexp"

	"github.com/mnalsup/sentry/core/config"
	"github.com/mnalsup/sentry/core/task/action"
	"github.com/mnalsup/sentry/monitoring"
	log "github.com/sirupsen/logrus"
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
	Actions  []action.Action
	Status   int
}

// Trigger is a definition of a trigger to watch
type Trigger struct {
	Source        string
	EventMatchers []*regexp.Regexp
}

type jsonTask struct {
	Triggers []jsonTrigger
	Actions  []action.InitialAction
}

type jsonTrigger struct {
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

// NewFromJSON returns a new Task from a json byte array
func NewFromJSON(conf *config.Configuration, taskBytes []byte) (*Task, error) {
	log.Tracef("New task from JSON: %s\n", string(taskBytes))
	var t jsonTask
	err := json.Unmarshal(taskBytes, &t)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var triggers = make([]Trigger, len(t.Triggers))
	for i, v := range t.Triggers {
		triggers[i] = triggerMustCompile(v.Source, v.EventMatchers)
	}
	var actions = make([]action.Action, len(t.Actions))
	for i, ia := range t.Actions {
		a, err := ia.Transform(conf)
		if err != nil {
			return nil, err
		}
		actions[i] = a
	}
	var task = Task{triggers, actions, INIT}
	return &task, nil
}

// MatchEvent receives and event and decides whether it matches a trigger
func (task *Task) MatchEvent(event *monitoring.Event) bool {
	log.Tracef("Matching on event: %s %s.\n", event.Name, event.Source)
	log.Traceln(task.Triggers)
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
		for _, a := range task.Actions {
			if task.Status < ONGOING {
				task.Status = ONGOING
				log.Debugf("Task status set to %d", task.Status)
				a.Exec()
				task.Status = POST
				log.Debugf("Task status set to %d", task.Status)
			}
		}
	} else {
		task.Status = INIT
		log.Debugf("Task status set to %d", task.Status)
	}
}
