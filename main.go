package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	username := ""

	fmt.Println("Please enter your Github username:")
	for {
		scanner.Scan()
		result := cleanInput(scanner.Text())
		if len(result) > 0 {
			username = result[0]
			break
		}

		log.Fatal("Failed to provide Github username")
	}

	fmt.Printf("Retrieving the Github activity of past 30 days of %s!\n\n", username)

	activities, err := getGitHubActivities(username)
	if err != nil {
		log.Fatal("Failed to get your activities", err)
	}

	for _, acitvity := range activities {
		fmt.Println(acitvity)
	}
}
