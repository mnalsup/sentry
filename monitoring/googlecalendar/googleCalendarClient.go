package googlecalendar

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mnalsup/sentry/core/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// Client is a struct for interacting with the google calendar api
type Client struct {
	Name            string
	CalendarID      string
	TokenFile       string
	CredentialsFile string
	Service         *calendar.Service
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config, tokFile string) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		/*
			tok = getTokenFromWeb(config)
			saveToken(tokFile, tok)
		*/
		log.Println("Not authorized for google calendar.")
	}
	return config.Client(context.Background(), tok)
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

func getCalendarServiceFromFile(credsFile string, tokFile string) *calendar.Service {
	b, err := ioutil.ReadFile(credsFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config, tokFile)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
	return srv
}

// New returns a Client
func New(conf *config.GoogleCalendarConnection) *Client {
	service := getCalendarServiceFromFile(conf.CredentialsFile, conf.TokenFile)
	return &Client{
		conf.Name,
		conf.CalendarID,
		conf.TokenFile,
		conf.CredentialsFile,
		service,
	}
}

// GetOnGoingEvents gets events that are on going
func (client *Client) GetOnGoingEvents() ([]*calendar.Event, error) {
	timeMax := time.Now().Add(12 * time.Hour).Format(time.RFC3339)
	// Covers a 25 hour period for all day events. Will miss multiday events
	timeMin := time.Now().Add(-(13 * time.Hour)).Format(time.RFC3339)
	log.Printf("TimeMin all events ending before: %s", timeMin)
	log.Printf("TimeMax all events starting after: %s", timeMax)
	events, err := client.Service.Events.List(client.CalendarID).ShowDeleted(false).
		SingleEvents(true).TimeMax(timeMax).TimeMin(timeMin).MaxResults(20).OrderBy("startTime").Do()
	if err != nil {
		log.Printf("Unable to retrieve next ten of the user's events: %v", err)
		log.Println("Attempting to get new credentials...")
		client.Service = getCalendarServiceFromFile(client.CredentialsFile, client.TokenFile)
		return nil, err
	}
	var ongoingEvents []*calendar.Event
	for _, event := range events.Items {
		start, err := time.Parse(time.RFC3339, event.Start.DateTime)
		if err != nil {
			log.Fatal(err)
		}
		end, err := time.Parse(time.RFC3339, event.End.DateTime)
		if err != nil {
			log.Fatal(err)
		}
		if time.Until(start) <= 0 && time.Until(end) >= 0 {
			ongoingEvents = append(ongoingEvents, event)
		}
	}
	return ongoingEvents, nil
}
