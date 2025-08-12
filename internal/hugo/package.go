package hugo

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Package struct {
	Title            string    `yaml:"title"`
	ImportPath       string    `yaml:"import_path"`
	RepoURL          string    `yaml:"repo_url"`
	Description      string    `yaml:"description"`
	Version          string    `yaml:"version"`
	DocumentationURL string    `yaml:"documentation_url"`
	License          string    `yaml:"license"`
	Author           string    `yaml:"author"`
	CreatedAt        time.Time `yaml:"created_at"`
	UpdatedAt        time.Time `yaml:"updated_at"`
}

func ReadPackage(filePath string) (*Package, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	// Extract frontmatter
	frontmatter, err := extractFrontmatter(data)
	if err != nil {
		return nil, err
	}
	
	var pkg Package
	if err := yaml.Unmarshal(frontmatter, &pkg); err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}
	
	return &pkg, nil
}

func WritePackage(filePath string, pkg *Package) error {
	// Marshal package to YAML
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(pkg); err != nil {
		return fmt.Errorf("failed to encode package: %w", err)
	}
	
	// Create markdown content with frontmatter
	content := fmt.Sprintf("---\n%s---\n", buf.String())
	
	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Write file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

func ListPackages(contentDir string) ([]string, error) {
	var packages []string
	
	entries, err := os.ReadDir(contentDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read content directory: %w", err)
	}
	
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		name := entry.Name()
		if strings.HasSuffix(name, ".md") && name != "_index.md" {
			packages = append(packages, filepath.Join(contentDir, name))
		}
	}
	
	return packages, nil
}

func extractFrontmatter(data []byte) ([]byte, error) {
	content := string(data)
	
	// Check for frontmatter delimiters
	if !strings.HasPrefix(content, "---\n") {
		return nil, fmt.Errorf("no frontmatter found")
	}
	
	// Find the closing delimiter
	endIndex := strings.Index(content[4:], "\n---")
	if endIndex == -1 {
		return nil, fmt.Errorf("invalid frontmatter format")
	}
	
	// Extract frontmatter content (without delimiters)
	frontmatter := content[4 : endIndex+4]
	return []byte(frontmatter), nil
}