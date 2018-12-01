package action

import (
	"encoding/json"
	"fmt"

	"github.com/mnalsup/sentry/core/config"
	"github.com/mnalsup/sentry/core/task/action/philipshue"
	log "github.com/sirupsen/logrus"
)

// Action is an interface defining an action
type Action interface {
	Exec() error
}

// Initial Action holds an unspecified action to be transformed
type InitialAction struct {
	ID        string
	JSONBytes []byte
}

type actionID struct {
	ID     string
	Source string
}

// UnmarshalJSON retrieves the ID and byte array for an action definition
func (a *InitialAction) UnmarshalJSON(b []byte) error {
	var id actionID
	log.Tracef("UnmarshalJSON for: %s\n", id)
	err := json.Unmarshal(b, &id)
	if err != nil {
		return err
	}
	a.ID = id.ID
	a.JSONBytes = b
	return nil
}

// Transform unmarshals an initial action based on the ID
func (a *InitialAction) Transform(conf *config.Configuration) (Action, error) {
	switch a.ID {
	case philipshue.LightsOnID:
		action, err := philipshue.NewLightsOnFromJSON(conf, a.JSONBytes)
		if err != nil {
			return nil, err
		}
		return action, nil
	}
	return nil, fmt.Errorf("No Actions found for ID: %s", a.ID)
}
