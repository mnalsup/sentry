package philipshue

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/amimof/huego"
	"github.com/mnalsup/sentry/core/config"
	"github.com/mnalsup/sentry/core/event"
	log "github.com/sirupsen/logrus"
)

// SunUpID is the identifier for this action
const SunUpID string = "PHILIPS_HUE_SUN_UP"

// SunUp specifies an action to turn on Philips Hue Lights
type SunUp struct {
	ID     string
	Source string
	HubURI string
	User   string
}

type sunUpConfig struct {
	ID     string
	Source string
}

// NewSunUp returns a *SunUp
func NewSunUp(id string, source string, immediate bool, apiURI string, user string) *SunUp {
	return &SunUp{
		ID:     id,
		Source: source,
		HubURI: apiURI,
		User:   user,
	}
}

// NewSunUpFromJSON converts json to SunUp
func NewSunUpFromJSON(conf *config.Configuration, dat []byte) (*SunUp, error) {
	var loac = sunUpConfig{}
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
	return &SunUp{
		ID:     loac.ID,
		Source: loac.Source,
		HubURI: api,
		User:   user,
	}, nil
}

// Exec fulfills the Action interface, executes the action specified by the action
func (ph *SunUp) Exec(event *event.Event) {
	log.WithFields(log.Fields{
		"file": "lightson.go",
		"func": "Exec",
	}).Tracef("Sun up exec with HubUri:  %s\n", ph.HubURI)
	bridge := huego.New(ph.HubURI, ph.User)
	lights, err := bridge.GetLights()
	if err != nil {
		return
	}
	for _, light := range lights {
		err = light.Bri(0)
	}
	for _, light := range lights {
		err = light.On()
		if err != nil {
			return
		}
	}
	duration := event.Duration.Minutes()
	log.Tracef("Duration: %f", duration)
	interval := int(255 / duration)
	log.Tracef("Interval:  %f to %d", 255/duration, interval)
	for i := 0; i < int(duration); i++ {
		for _, light := range lights {
			log.Debugf("Setting brightness to: %d", uint8(i*interval))
			err := light.Bri(uint8(i * interval))
			if err != nil {
				return
			}
		}
		for i := 0; i < 60; i++ {
			fmt.Println("still here")
			time.Sleep(1 * time.Second)
		}
	}
	return
}
