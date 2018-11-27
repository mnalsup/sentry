package task

import (
	"io/ioutil"
	"log"
)

// LoadTasksFromFile reads the /etc/sentry/tasks directory
func LoadTasksFromFile(directory string) ([]*Task, error) {
	log.Println("Reading in tasks from /etc/sentry/tasks...")
	// TODO: use ReadDir and get some number of files
	dir, err := ioutil.ReadDir("/etc/sentry/tasks")
	if err != nil {
		return nil, err
	}
	tasks := make([]*Task, len(dir))
	for i, f := range dir {
		log.Printf("Found %s.", f.Name())
		dat, err := ioutil.ReadFile(directory)
		if err != nil {
			return nil, err
		}
		t, err := NewFromJSON(dat)
		tasks[i] = t
		log.Println(t)
	}
	return tasks, nil
}
