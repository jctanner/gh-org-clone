package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Client handles GitHub API requests
type Client struct {
	httpClient *http.Client
	token      string
}

// NewClient creates a new GitHub API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
		token:      os.Getenv("GITHUB_TOKEN"),
	}
}

// ListRepositories fetches all repositories for an organization or user
func (c *Client) ListRepositories(orgOrUser string) ([]Repository, error) {
	var allRepos []Repository
	page := 1

	for {
		repos, hasMore, err := c.fetchRepositoriesPage(orgOrUser, page)
		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos...)

		if !hasMore {
			break
		}
		page++
	}

	return allRepos, nil
}

// fetchRepositoriesPage fetches a single page of repositories
func (c *Client) fetchRepositoriesPage(orgOrUser string, page int) ([]Repository, bool, error) {
	// Try organization endpoint first
	repos, hasMore, err := c.fetchPage(fmt.Sprintf("https://api.github.com/orgs/%s/repos?per_page=100&page=%d", orgOrUser, page))
	if err == nil {
		return repos, hasMore, nil
	}

	// If org endpoint fails with 404, try user endpoint
	if isNotFoundError(err) {
		return c.fetchPage(fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=100&page=%d", orgOrUser, page))
	}

	return nil, false, err
}

// fetchPage fetches a single page from the given URL
func (c *Client) fetchPage(url string) ([]Repository, bool, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, false, err
	}

	// Add authentication if token is available
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	// Check rate limit
	if err := c.checkRateLimit(resp); err != nil {
		return nil, false, err
	}

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode == http.StatusNotFound {
			return nil, false, &NotFoundError{Message: string(body)}
		}
		return nil, false, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response body
	var repos []Repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, false, err
	}

	// Check if there are more pages (simplified: if we got 100 results, there might be more)
	hasMore := len(repos) == 100

	return repos, hasMore, nil
}

// checkRateLimit monitors rate limit headers and waits if necessary
func (c *Client) checkRateLimit(resp *http.Response) error {
	remainingStr := resp.Header.Get("X-RateLimit-Remaining")
	resetStr := resp.Header.Get("X-RateLimit-Reset")

	if remainingStr == "" || resetStr == "" {
		return nil
	}

	remaining, err := strconv.Atoi(remainingStr)
	if err != nil {
		return nil
	}

	reset, err := strconv.ParseInt(resetStr, 10, 64)
	if err != nil {
		return nil
	}

	// If we're out of requests or about to be, wait
	if remaining == 0 || resp.StatusCode == http.StatusForbidden {
		resetTime := time.Unix(reset, 0)
		waitDuration := time.Until(resetTime)

		if waitDuration > 0 {
			fmt.Printf("Rate limit exceeded. Waiting until %s (%d seconds)\n",
				resetTime.Format("15:04:05"),
				int(waitDuration.Seconds()))
			time.Sleep(waitDuration + time.Second) // Add 1 second buffer
		}
	}

	return nil
}

// NotFoundError represents a 404 response from the API
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

// isNotFoundError checks if an error is a 404 response
func isNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}
