package main

import (
	"fmt"
	"os"

	"github.com/jctanner/gh-org-clone/clone"
	"github.com/jctanner/gh-org-clone/github"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <org-name>\n", os.Args[0])
		os.Exit(1)
	}

	orgName := os.Args[1]

	fmt.Printf("Fetching repositories for: %s\n", orgName)

	// Create GitHub client and fetch repositories
	client := github.NewClient()
	repos, err := client.ListRepositories(orgName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching repositories: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d repositories\n\n", len(repos))

	// Clone all repositories
	targetDir := orgName
	result := clone.CloneAll(repos, targetDir)

	// Print summary
	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Cloned: %d\n", result.Cloned)
	fmt.Printf("  Failed: %d\n", result.Failed)
	fmt.Printf("  Skipped: %d\n", result.Skipped)
}
