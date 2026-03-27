package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jctanner/gh-org-clone/clone"
	"github.com/jctanner/gh-org-clone/github"
)

func main() {
	// Define flags
	pathFlag := flag.String("path", ".", "Base directory for cloning repositories")
	branchFlag := flag.String("branch", "", "Specific branch to clone (skips repos without this branch)")
	sshFlag := flag.Bool("ssh", false, "Force SSH clone URLs for all repositories")

	// Custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <org-name>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	// Check for positional argument
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	orgName := flag.Arg(0)

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
	targetDir := filepath.Join(*pathFlag, orgName)
	result := clone.CloneAll(repos, targetDir, *branchFlag, *sshFlag)

	// Print summary
	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Cloned: %d\n", result.Cloned)
	fmt.Printf("  Failed: %d\n", result.Failed)
	fmt.Printf("  Skipped: %d\n", result.Skipped)
}
