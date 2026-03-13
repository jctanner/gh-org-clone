package clone

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/jctanner/gh-org-clone/github"
)

// Result tracks the outcome of cloning operations
type Result struct {
	Cloned  int
	Failed  int
	Skipped int
}

// CloneAll clones all repositories to the target directory
func CloneAll(repos []github.Repository, targetDir string, branch string) Result {
	result := Result{}

	// Create target directory
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", targetDir, err)
		return result
	}

	for i, repo := range repos {
		fmt.Printf("Cloning repository %d/%d: %s\n", i+1, len(repos), repo.Name)
		fmt.Printf("  URL: %s\n", repo.CloneURL)
		if branch != "" {
			fmt.Printf("  Branch: %s\n", branch)
		}

		repoPath := filepath.Join(targetDir, repo.Name)

		// Check if repository already exists
		if _, err := os.Stat(repoPath); err == nil {
			fmt.Printf("  Skipped (already exists)\n")
			result.Skipped++
			continue
		}

		// Build git clone command
		var cmd *exec.Cmd
		if branch != "" {
			cmd = exec.Command("git", "clone", "-b", branch, repo.CloneURL, repoPath)
		} else {
			cmd = exec.Command("git", "clone", repo.CloneURL, repoPath)
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			// If a specific branch was requested and clone failed, skip the repo
			if branch != "" {
				fmt.Printf("  Skipped (branch '%s' not found)\n", branch)
				result.Skipped++
			} else {
				fmt.Printf("  Failed: %v\n", err)
				result.Failed++
			}
			continue
		}

		result.Cloned++
	}

	return result
}
