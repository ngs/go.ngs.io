package github

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	gh "github.com/cli/go-gh/v2"
)

type Repository struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	License     *License   `json:"license"`
	Topics      []string   `json:"topics"`
	Owner       Owner      `json:"owner"`
}

type License struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SPDXID string `json:"spdx_id"`
}

type Owner struct {
	Login string `json:"login"`
	Name  string `json:"name"`
}

type Release struct {
	TagName string `json:"tag_name"`
}

type Tag struct {
	Name string `json:"name"`
}

type Readme struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

func GetReadme(owner, repo string) (string, error) {
	args := []string{"api", fmt.Sprintf("repos/%s/%s/readme", owner, repo)}

	stdout, _, err := gh.Exec(args...)
	if err != nil {
		return "", nil // No README available
	}

	var readme Readme
	if err := json.Unmarshal(stdout.Bytes(), &readme); err != nil {
		return "", fmt.Errorf("failed to parse readme data: %w", err)
	}

	if readme.Encoding != "base64" {
		return "", fmt.Errorf("unexpected encoding: %s", readme.Encoding)
	}

	content, err := base64.StdEncoding.DecodeString(readme.Content)
	if err != nil {
		return "", fmt.Errorf("failed to decode readme: %w", err)
	}

	return string(content), nil
}

func GetRepository(owner, repo string) (*Repository, error) {
	args := []string{"api", fmt.Sprintf("repos/%s/%s", owner, repo)}
	
	stdout, _, err := gh.Exec(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository: %w", err)
	}
	
	var repository Repository
	if err := json.Unmarshal(stdout.Bytes(), &repository); err != nil {
		return nil, fmt.Errorf("failed to parse repository data: %w", err)
	}
	
	return &repository, nil
}

func GetLatestVersion(owner, repo string) (string, error) {
	// Try to get latest release first
	args := []string{"api", fmt.Sprintf("repos/%s/%s/releases/latest", owner, repo)}
	
	stdout, _, err := gh.Exec(args...)
	if err == nil {
		var release Release
		if err := json.Unmarshal(stdout.Bytes(), &release); err == nil && release.TagName != "" {
			return release.TagName, nil
		}
	}
	
	// If no releases, try tags
	args = []string{"api", fmt.Sprintf("repos/%s/%s/tags", owner, repo)}
	
	stdout, _, err = gh.Exec(args...)
	if err != nil {
		return "", nil // No version available
	}
	
	var tags []Tag
	if err := json.Unmarshal(stdout.Bytes(), &tags); err != nil {
		return "", nil
	}
	
	if len(tags) > 0 {
		return tags[0].Name, nil
	}
	
	return "", nil
}

func ParseRepoURL(url string) (owner, repo string, err error) {
	// Parse GitHub URL formats:
	// https://github.com/owner/repo
	// https://github.com/owner/repo.git
	// git@github.com:owner/repo.git
	
	// Try HTTPS URL with path parsing
	if len(url) > 19 && url[:19] == "https://github.com/" {
		path := url[19:]
		parts := strings.Split(path, "/")
		if len(parts) >= 2 {
			owner = parts[0]
			repo = parts[1]
			// Remove .git suffix if present
			if strings.HasSuffix(repo, ".git") {
				repo = repo[:len(repo)-4]
			}
			return owner, repo, nil
		}
	}
	
	// Try git SSH URL
	if len(url) > 15 && url[:15] == "git@github.com:" {
		path := url[15:]
		parts := strings.Split(path, "/")
		if len(parts) >= 2 {
			owner = parts[0]
			repo = parts[1]
			// Remove .git suffix if present
			if strings.HasSuffix(repo, ".git") {
				repo = repo[:len(repo)-4]
			}
			return owner, repo, nil
		}
	}
	
	return "", "", fmt.Errorf("invalid GitHub URL format: %s", url)
}