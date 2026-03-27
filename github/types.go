package github

// Repository represents a GitHub repository from the API response
type Repository struct {
	Name     string `json:"name"`
	CloneURL string `json:"clone_url"`
	SSHURL   string `json:"ssh_url"`
	Private  bool   `json:"private"`
}

// RateLimitInfo contains rate limit details from response headers
type RateLimitInfo struct {
	Remaining int
	Reset     int64 // Unix timestamp
}
