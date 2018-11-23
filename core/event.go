package core

import (
	"regexp"
)

type Event struct {
	Name          string
	EventMatchers []*regexp.Regexp
}

/*
MustCompile creates an Event and panics if the expression does not compile
*/
func EventMustCompile(name string, exprs []string) *Event {
	length := len(exprs)
	matchers := make([]*regexp.Regexp, length)
	for i, v := range exprs {
		matchers[i] = regexp.MustCompile(v)
	}
	return &Event{name, matchers}
}
