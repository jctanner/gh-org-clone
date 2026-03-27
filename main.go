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
	listFlag := flag.Bool("list", false, "List repositories without cloning")

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

	// Count private vs public repos
	privateCount := 0
	for _, repo := range repos {
		if repo.Private {
			privateCount++
		}
	}
	fmt.Printf("Found %d repositories (%d public, %d private)\n\n", len(repos), len(repos)-privateCount, privateCount)

	// If list mode, just print repositories and exit
	if *listFlag {
		for i, repo := range repos {
			fmt.Printf("%d. %s\n", i+1, repo.Name)
			if repo.Private {
				fmt.Printf("   Private: yes\n")
				fmt.Printf("   Clone URL: %s\n", repo.SSHURL)
			} else {
				fmt.Printf("   Private: no\n")
				fmt.Printf("   Clone URL: %s\n", repo.CloneURL)
			}
		}
		return
	}

	// Clone all repositories
	targetDir := filepath.Join(*pathFlag, orgName)
	result := clone.CloneAll(repos, targetDir, *branchFlag, *sshFlag)

	// Print summary
	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Cloned: %d\n", result.Cloned)
	fmt.Printf("  Failed: %d\n", result.Failed)
	fmt.Printf("  Skipped: %d\n", result.Skipped)
}
