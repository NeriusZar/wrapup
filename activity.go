package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Event struct {
	ID        string
	Type      string
	CreatedAt time.Time `json:"created_at"`
	Repo      Repo
	Payload   map[string]interface{}
}

type Repo struct {
	ID   int
	Name string
	Url  string
}

func getGitHubActivities(username string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	resp, err := http.Get(url)
	if err != nil {
		//proccess error
		return nil, nil
	}

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		fmt.Println("Failed to decode the response", err)
	}

	formattedEvents := make([]string, len(events))
	for i, e := range events {
		formattedEvents[i] = e.String()
	}

	return formattedEvents, nil
}

func (e Event) String() string {
	switch e.Type {
	case "WatchEvent":
		return fmt.Sprintf("[%s] You stared a %s repository ⭐️. %s", e.CreatedAt.Format("Jan 2"), e.Repo.Name, e.Repo.PublicUrl())
	default:
		return "This event is not supported yet"
	}
}

func (r Repo) PublicUrl() string {
	return "https://github.com/" + r.Name
}
