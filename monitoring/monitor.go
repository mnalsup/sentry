package monitoring

import (
	"fmt"
	"log"
	"time"

	"github.com/mnalsup/sentry/core/config"
	"github.com/mnalsup/sentry/monitoring/googlecalendar"
)

// Event is exported to match triggers for tasks
type Event struct {
	Name   string
	Source string
}

/*
Monitor starts the monitor of calendar events
*/
func Monitor(conf *config.Configuration, c chan *Event) {
	log.Println(conf)
	// TODO: Generalize for multiple calendars
	calSvc := googlecalendar.New(&conf.GoogleCalendar.Sources[0])
	for {
		events, err := calSvc.GetOnGoingEvents()
		log.Println(events)
		if err != nil {
			log.Printf("Unable to retrieve next ten of the user's events: %v", err)
		} else {
			fmt.Println("OnGoing events:")
			if len(events) == 0 {
				fmt.Println("No ongoing events found.")
			} else {
				for _, event := range events {
					// fmt.Printf("%v (Start: %v)(End: %v)\n", event.Summary, event.Start, event.End)
					// TODO: Generalize for multiple calendars
					c <- &Event{event.Summary, conf.GoogleCalendar.Sources[0].Name}
				}
			}
		}
		time.Sleep(30 * time.Second)
	}
	// close(c)
}
