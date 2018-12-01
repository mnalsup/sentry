package config

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

/*
config:
  googlecalendar:
    enabled: true
    sources:
      - name: alsupmn-cal-primary
        calendarId: primary
        token: /etc/sentry/token.json
        credentials: /etc/sentry/credentials.json
  philipshue:
    enabled: true
    sources:
      - name: hue-home
        apiUri: http://192.168.0.109/api
        user: uTzr9S5IYlzCv5BCs0j7l65Gyi0YO4GiURRmYACP
*/

// Configuration is the top level config
type Configuration struct {
	GoogleCalendar GoogleCalendarConfig
	PhilipsHue     PhilipsHueConfig
}

// PhilipsHueConfig has the enabled flag and connections
type PhilipsHueConfig struct {
	Enabled bool
	Sources []PhilipsHueConnection
}

// PhilipsHueConnection contains connection information for a philips hue hub
type PhilipsHueConnection struct {
	Name   string
	APIURI string
	User   string
}

// GoogleCalendarConfig contains enabled flag and connections
type GoogleCalendarConfig struct {
	Enabled bool
	Sources []GoogleCalendarConnection
}

// GoogleCalendarConnection contains information for connection
type GoogleCalendarConnection struct {
	Name            string
	CalendarID      string
	TokenFile       string
	CredentialsFile string
}

// EventConfig represents a configuration for an event
type EventConfig struct {
	Name          string
	EventMatchers []string
}

// Config is config

// New creates the config
func New() *Configuration {
	// Log configuration
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)

	// build config
	var Config Configuration
	fmt.Println("Initializing Configs...")
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/sentry")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.WithFields(log.Fields{
			"file":     "config.go",
			"function": "New",
		}).Fatalf("fatal error config file: %s", err)
	}
	err = viper.Unmarshal(&Config)
	if err != nil { // Handle errors reading the config file
		log.WithFields(log.Fields{
			"file":     "config.go",
			"function": "New",
		}).Fatalf("fatal error config file: %s", err)
	}
	return &Config
}
