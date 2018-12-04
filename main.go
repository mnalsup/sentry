package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	log "github.com/sirupsen/logrus"
	"github.com/mnalsup/sentry/core/task"
	"github.com/mnalsup/sentry/core/event"
	"github.com/mnalsup/sentry/core/config"
	"github.com/mnalsup/sentry/monitoring"
	"github.com/mnalsup/sentry/ui"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "/etc/sentry/token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	config := config.New()
	tasks, err := task.LoadTasksFromFile(config, "/etc/sentry/tasks/wakeup.json")
	if err != nil {        // Handle errors reading the config file
		log.WithFields(log.Fields{
			"file":     "main.go",
			"function": "main",
		}).Fatalln(err)
	}
	events := make(chan *event.Event)
	go monitoring.Monitor(config, events)
	go task.Subscribe(tasks, events)
	ui.ServeUI()
}
