package monitoring

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mnalsup/sentry/core"
	"github.com/mnalsup/sentry/core/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "/etc/sentry/token.json"
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

func getCalendarServiceFromFile() *calendar.Service {
	b, err := ioutil.ReadFile("/etc/sentry/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
	return srv
}

/*
Monitor starts the monitor of calendar events
*/
func Monitor(conf *config.Configuration) {
	sentryEvents := make([]*core.Event, len(conf.Events))
	for i, v := range conf.Events {
		sentryEvents[i] = core.EventMustCompile(v.Name, v.EventMatchers)
	}
	srv := getCalendarServiceFromFile()
	for {
		t := time.Now().Format(time.RFC3339)
		events, err := srv.Events.List("primary").ShowDeleted(false).
			SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
		if err != nil {
			log.Printf("Unable to retrieve next ten of the user's events: %v", err)
			log.Println("Attempting to get new credentials...")
			srv = getCalendarServiceFromFile()
		} else {
			fmt.Println("Upcoming events:")
			if len(events.Items) == 0 {
				fmt.Println("No upcoming events found.")
			} else {
				for _, item := range events.Items {
					/* start filter */
					for _, event := range sentryEvents {
						for _, matcher := range event.EventMatchers {
							found := matcher.MatchString(item.Summary)
							if found {
								date := item.Start.DateTime
								if date == "" {
									date = item.Start.Date
								}
								fmt.Printf("%v (%v)\n", item.Summary, date)
							}
						}
					}
					/* end filter */
				}
			}
		}
		time.Sleep(30 * time.Second)
	}
}
