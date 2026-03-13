# gh-org-clone

Command-line tool to clone all public repositories from a GitHub organization or user account.

## Requirements

- Go 1.16 or later
- git

## Building

```bash
make build
```

Binary will be created at `bin/gh-org-clone`.

## Usage

```bash
./bin/gh-org-clone [options] <org-or-username>
```

### Options

- `--path <directory>`: Base directory for cloning repositories (default: current directory)
- `--branch <branch-name>`: Clone only the specified branch; skips repositories that don't have this branch

### Examples

Basic usage:
```bash
./bin/gh-org-clone opendatahub-io
```

This creates a directory named after the organization and clones all repositories into it.

Clone to a specific directory:
```bash
./bin/gh-org-clone --path=checkouts opendatahub-io
```

This creates `checkouts/opendatahub-io/` and clones all repositories there.

Clone only a specific branch:
```bash
./bin/gh-org-clone --branch=v2.0 opendatahub-io
```

This clones only the `v2.0` branch from each repository. Repositories without this branch are skipped.

Combine both options:
```bash
./bin/gh-org-clone --path=checkouts --branch=main opendatahub-io
```

## Authentication

The tool works without authentication but is subject to GitHub's rate limit of 60 requests per hour.

To use authenticated requests (5000 requests per hour), set the GITHUB_TOKEN environment variable:

```bash
export GITHUB_TOKEN=your_token_here
./bin/gh-org-clone opendatahub-io
```

## Rate Limiting

The tool monitors GitHub API rate limits. When the limit is reached, it waits until the rate limit resets before continuing.

## Installation

```bash
make install
```

Copies the binary to `~/bin/gh-org-clone`.

## Behavior

- Tries organization endpoint first, falls back to user endpoint on 404
- Skips repositories that already exist locally
- When `--branch` is specified, skips repositories that don't have that branch
- Reports summary of cloned, failed, and skipped repositories
- Shows git clone output in real-time
