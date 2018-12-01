package philipshue

import (
	"encoding/json"
	"fmt"

	"github.com/amimof/huego"
	"github.com/mnalsup/sentry/core/config"
	log "github.com/sirupsen/logrus"
)

// LightsOnID is the identifier for this action
const LightsOnID string = "PHILIPS_HUE_LIGHTS_ON"

// LightsOn specifies an action to turn on Philips Hue Lights
type LightsOn struct {
	ID     string
	Source string
	HubURI string
	User   string
}

type LightsOnConfig struct {
	ID     string
	Source string
}

// NewLightsOn returns a *LightsOn
func NewLightsOn(id string, source string, apiURI string, user string) *LightsOn {
	return &LightsOn{
		ID:     id,
		Source: source,
		HubURI: apiURI,
		User:   user,
	}
}

// NewLightsOnFromJSON converts json to LightsOn
func NewLightsOnFromJSON(conf *config.Configuration, dat []byte) (*LightsOn, error) {
	var loac = LightsOnConfig{}
	err := json.Unmarshal(dat, &loac)
	if err != nil {
		return nil, err
	}
	var api string
	var user string
	for _, conn := range conf.PhilipsHue.Sources {
		if conn.Name == loac.Source {
			api = conn.APIURI
			user = conn.User
		}
	}
	if api == "" || user == "" {
		return nil, fmt.Errorf("no sources configured matched: %s", loac.Source)
	}
	return &LightsOn{
		ID:     loac.ID,
		Source: loac.Source,
		HubURI: api,
		User:   user,
	}, nil
}

// Exec fulfills the Action interface, executes the action specified by the action
func (ph *LightsOn) Exec() error {
	log.WithFields(log.Fields{
		"file": "lightson.go",
		"func": "Exec",
	}).Tracef("Lights on exec with HubUri:  %s\n", ph.HubURI)
	bridge := huego.New(ph.HubURI, ph.User)
	lights, err := bridge.GetLights()
	if err != nil {
		return err
	}
	for _, light := range lights {
		fmt.Println(light)
		err = light.Off()
		if err != nil {
			return err
		}
	}
	return nil
}
