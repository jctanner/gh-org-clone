# GitHub Organization Cloner - Design Document

## Overview

CLI tool to clone all public repositories from a GitHub organization or user account.

## Project Name

- Repository: `gh-org-clone`
- Binary: `gh-org-clone`
- Go module: `github.com/jctanner/gh-org-clone`

Naming rationale: Follows Unix philosophy of doing one thing well. Focused specifically on cloning organization/user repositories rather than expanding scope.

## Requirements

- Accept single argument: GitHub organization or username
- Discover all public repositories for the target
- Clone each repository to local filesystem
- Handle errors gracefully (rate limits, network failures, invalid targets)

## CLI Interface

```
github-org-cloner <org-name>
```

Example:
```
github-org-cloner opendatahub-io
```

## Technical Approach

### GitHub API Integration

Use GitHub REST API v3 to list repositories:
- Endpoint: `GET /orgs/{org}/repos` for organizations
- Endpoint: `GET /users/{username}/repos` for users
- Pagination required (default 30 items per page)
- Authentication via GITHUB_TOKEN environment variable if present
- Rate limits: 60/hour unauthenticated, 5000/hour authenticated

#### Rate Limit Handling

- Monitor response headers: `X-RateLimit-Remaining`, `X-RateLimit-Reset`
- On rate limit exceeded (403 with rate limit message or remaining = 0):
  - Calculate wait time from `X-RateLimit-Reset` header (Unix timestamp)
  - Log message: "Rate limit exceeded. Waiting until HH:MM:SS (N seconds)"
  - Sleep until rate limit reset time
  - Resume operations

### Cloning Strategy

- Create target directory: `./<org-name>/`
- Clone each repository into subdirectory by repo name
- Use `git clone` via exec
- Skip repositories that already exist locally

### Error Handling

- 404: Invalid organization/username
- 403: Rate limit exceeded
- Network errors: Report and continue with next repository
- Clone failures: Log error, continue with remaining repos

## Implementation

### Structure

```
main.go           - CLI entry point
github/
  client.go       - GitHub API client
  types.go        - Response structures
clone/
  cloner.go       - Git clone operations
```

### Dependencies

- Standard library: `net/http`, `os/exec`, `encoding/json`
- No external dependencies required

### Flow

1. Parse command line argument
2. Attempt to fetch from `/orgs/{name}/repos`
3. If 404, try `/users/{name}/repos`
4. Paginate through all results
5. Create target directory
6. Clone each repository in sequence
7. Report summary (success count, failure count)

## Output Format

```
Fetching repositories for: opendatahub-io
Found 47 repositories

Cloning repository 1/47: opendatahub-operator
Cloning repository 2/47: kubeflow
...

Summary:
  Cloned: 45
  Failed: 2
  Skipped: 0
```
