# go.ngs.io

Go module vanity import path service for packages hosted at go.ngs.io. This repository manages a Hugo-based website that provides custom import paths for Go packages.

## Overview

This repository hosts a Hugo static site that enables Go modules to be installed using the custom domain `go.ngs.io`.

For example:
```bash
go install go.ngs.io/freecal@latest
```

## Prerequisites

- Go 1.22 or later
- Hugo (for site generation)
- GitHub CLI (`gh`) authenticated for API access

## Installation

Install the management tools:

```bash
# Install add-package command
go install ./cmd/add-package

# Install update-packages command  
go install ./cmd/update-packages
```

## Usage

### Adding a New Package

Use the `add-package` command to add a new Go package to the site:

```bash
# Add a package with automatic GitHub detection
add-package mypackage --repo https://github.com/ngs/mypackage

# Add with custom import path
add-package tools --import-path go.ngs.io/tools --repo https://github.com/ngs/go-tools

# Add with author information
add-package utils --repo https://github.com/ngs/utils --author "Atsushi Nagase"
```

The command will:
1. Fetch metadata from GitHub API (description, license, timestamps)
2. Detect the latest version/release
3. Create a markdown file in the `content/` directory
4. Validate the Hugo site builds correctly

### Updating Package Metadata

Use the `update-packages` command to update package metadata from GitHub:

```bash
# Update all packages
update-packages

# Update specific packages
update-packages freecal servedir

# Preview changes without updating (dry run)
update-packages --dry-run

# Update author information from GitHub
update-packages --update-author

# Update timestamps for missing/private repositories
update-packages --update-missing
```

### Manual Package Management

Package files are stored as markdown files in the `content/` directory with YAML frontmatter:

```yaml
---
title: "packagename"
import_path: "go.ngs.io/packagename"
repo_url: "https://github.com/ngs/packagename"
description: "Package description from GitHub"
version: "v1.0.0"
documentation_url: "https://pkg.go.dev/go.ngs.io/packagename"
license: "MIT"
author: "Atsushi Nagase"
created_at: 2024-01-01T00:00:00Z
updated_at: 2024-12-01T00:00:00Z
---
```

### Building the Site

After adding or updating packages, build the Hugo site:

```bash
# Build the site
hugo --gc --minify

# Serve locally for testing
hugo server
```

## Command Options

### add-package

```
Usage: add-package <package-name> [options]

Options:
  --import-path string   Custom import path (default: go.ngs.io/<package-name>)
  --repo string         GitHub repository URL
  --author string       Package author name
  -h, --help           Show help message
```

### update-packages

```
Usage: update-packages [package-names...] [options]

Options:
  --dry-run          Show what would be updated without making changes
  --update-author    Also update author information from GitHub
  --update-missing   Update timestamps for repositories that return 404
  -h, --help        Show help message
```

## How it Works

When Go tools fetch a module with a custom import path, they:

1. Request `https://go.ngs.io/freecal?go-get=1`
2. Parse the HTML response for `<meta name="go-import">` tags
3. Use the repository URL specified in the meta tag to fetch the actual code

## Project Structure

```
.
├── cmd/
│   ├── add-package/      # Command to add new packages
│   └── update-packages/   # Command to update package metadata
├── internal/
│   ├── github/           # GitHub API client
│   └── hugo/             # Hugo package file operations
├── content/              # Package markdown files
├── layouts/              # Hugo templates
├── static/               # Static assets
├── hugo.toml            # Hugo configuration
└── go.mod               # Go module definition
```

## Workflow

1. **Add a new package:**
   ```bash
   add-package myproject --repo https://github.com/ngs/myproject
   ```

2. **Review the generated file:**
   ```bash
   cat content/myproject.md
   ```

3. **Update all packages periodically:**
   ```bash
   update-packages
   ```

4. **Build and test locally:**
   ```bash
   hugo server
   ```

5. **Commit and push changes:**
   ```bash
   git add content/
   git commit -m "Add/Update packages"
   git push
   ```

## GitHub API Rate Limits

The tools use the GitHub API via the `gh` CLI tool. Rate limits:
- Authenticated: 5,000 requests/hour
- Unauthenticated: 60 requests/hour

Make sure `gh` is authenticated:
```bash
gh auth login
```

## Setup

This site is deployed via GitHub Pages with a custom domain.

### DNS Configuration

Add a CNAME record for `go.ngs.io` pointing to `ngs.github.io`.

### GitHub Pages Settings

1. Repository Settings > Pages
2. Source: Deploy from GitHub Actions
3. Custom domain: go.ngs.io

## License

MIT License - See the individual module repositories for their respective licenses.

## Author

Atsushi Nagase - https://ngs.io