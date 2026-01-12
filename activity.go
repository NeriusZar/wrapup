package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

type Commit struct {
	Sha    string
	Commit CommitDetails
}

type CommitDetails struct {
	Message string
}

type PushEventPayload struct {
	Head         string
	Before       string
	RepositoryId int
	PushId       int
	Ref          string
}

func getGitHubActivities(username string) ([]string, error) {
	urlString := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	resp, err := http.Get(urlString)
	if err != nil {
		//proccess error
		return nil, nil
	}
	defer resp.Body.Close()

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		fmt.Println("Failed to decode the response", err)
	}

	formattedEvents := []string{}
	for _, e := range events {
		formattedEvents = append(formattedEvents, e.String())

		newCommitsHappened := e.Type == "PushEvent"
		if newCommitsHappened {
			payload, err := convertToPushPayload(e.Payload)
			if err != nil {
				return nil, err
			}

			if commits, err := getCommits(username, e.Repo.RepoNameOnly(), payload.Head); err == nil {
				for _, c := range commits {
					formattedEvents = append(formattedEvents, c.String())
				}
			}
		}
	}

	return formattedEvents, nil
}

func convertToPushPayload(payloadMap map[string]interface{}) (PushEventPayload, error) {
	bytes, err := json.Marshal(payloadMap)
	if err != nil {
		return PushEventPayload{}, err
	}

	var payload PushEventPayload
	if err := json.Unmarshal(bytes, &payload); err != nil {
		return PushEventPayload{}, err
	}

	return payload, nil
}

func getCommits(username, repo, startingCommit string) ([]Commit, error) {
	urlString := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", username, repo)

	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("Failed to prase url", err)
		return nil, err
	}

	queryParams := parsedUrl.Query()
	queryParams.Add("sha", startingCommit)
	parsedUrl.RawQuery = queryParams.Encode()

	resp, err := http.Get(parsedUrl.String())
	if err != nil {
		fmt.Println("Failed to fetch commits", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("Failed to fetch commits")
	}
	defer resp.Body.Close()

	var commits []Commit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		fmt.Println("Failed to decode commits", err)
		return nil, err
	}

	return commits, nil
}

func (e Event) String() string {
	switch e.Type {
	case "WatchEvent":
		return fmt.Sprintf("[%s] You stared a %s repository ⭐️. %s", e.CreatedAt.Format("Jan 2"), e.Repo.Name, e.Repo.PublicUrl())
	case "PushEvent":
		return fmt.Sprintf("[%s] You pushed one or multiple commits to a %s repository 🛠️", e.CreatedAt.Format("Jan 2"), e.Repo.Name)
	default:
		return e.Type + " event is not supported yet"
	}
}

func (c Commit) String() string {
	return " - " + c.Commit.Message + " ⚒️"
}

func (r Repo) PublicUrl() string {
	return "https://github.com/" + r.Name
}

func (r Repo) RepoNameOnly() string {
	parts := strings.Split(r.Name, "/")
	return parts[len(parts)-1]
}
