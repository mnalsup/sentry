package config

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Configuration is the top level config
type Configuration struct {
	Events []EventConfig
}

// EventConfig represents a configuration for an event
type EventConfig struct {
	Name          string
	EventMatchers []string
}

// Config is config

// New creates the config
func New() *Configuration {
	var Config Configuration
	fmt.Println("Initializing Configs...")
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/sentry")
	viper.SetConfigType("yaml")
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	err = viper.Unmarshal(&Config)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	fmt.Println(fmt.Sprintf("%s", Config.Events[0].Name))
	return &Config
}
