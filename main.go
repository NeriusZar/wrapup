package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	username := ""
	if len(os.Args) > 1 {
		username = os.Args[1]
	}

	if username == "" {
		log.Fatal("Github username was not provided")
		return
	}

	fmt.Println("Let's look what you did in the last 30 days!")
	activities, err := getGitHubActivities(username)
	if err != nil {
		log.Fatal("Failed to get your activities", err)
	}

	for _, acitvity := range activities {
		fmt.Printf("- %s\n", acitvity)
	}
}
