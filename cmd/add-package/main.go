package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/pflag"
	"go.ngs.io/internal/github"
	"go.ngs.io/internal/hugo"
)

func main() {
	var (
		importPath string
		repoURL    string
		author     string
		help       bool
	)

	pflag.StringVar(&importPath, "import-path", "", "Custom import path (e.g., go.ngs.io/package)")
	pflag.StringVar(&repoURL, "repo", "", "GitHub repository URL")
	pflag.StringVar(&author, "author", "", "Package author name")
	pflag.BoolVarP(&help, "help", "h", false, "Show help message")
	pflag.Parse()

	if help || pflag.NArg() < 1 {
		printUsage()
		os.Exit(0)
	}

	packageName := pflag.Arg(0)
	if err := addPackage(packageName, importPath, repoURL, author); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func printUsage() {
	fmt.Println("Usage: add-package <package-name> [options]")
	fmt.Println("\nAdd a new Go package to go.ngs.io")
	fmt.Println("\nOptions:")
	pflag.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("  add-package mypackage --repo https://github.com/username/mypackage")
	fmt.Println("  add-package tools --import-path go.ngs.io/tools --repo https://github.com/ngs/tools")
}

func addPackage(packageName, importPath, repoURL, author string) error {
	// Validate package name
	if packageName == "" {
		return fmt.Errorf("package name is required")
	}

	// Set default import path if not provided
	if importPath == "" {
		importPath = fmt.Sprintf("go.ngs.io/%s", packageName)
	}

	// Parse repository URL if provided
	var owner, repo string
	var err error
	if repoURL != "" {
		owner, repo, err = github.ParseRepoURL(repoURL)
		if err != nil {
			return fmt.Errorf("invalid repository URL: %w", err)
		}
	} else {
		// Try to guess from package name
		owner = "ngs"
		repo = packageName
		repoURL = fmt.Sprintf("https://github.com/%s/%s", owner, repo)
	}

	fmt.Printf("Adding package '%s' from %s...\n", packageName, repoURL)

	// Create package struct
	pkg := &hugo.Package{
		Title:            packageName,
		ImportPath:       importPath,
		RepoURL:          repoURL,
		DocumentationURL: fmt.Sprintf("https://pkg.go.dev/%s", importPath),
		Author:           author,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	// Fetch metadata from GitHub
	fmt.Println("Fetching repository metadata from GitHub...")
	ghRepo, err := github.GetRepository(owner, repo)
	if err != nil {
		// Exit with error if repository doesn't exist
		return fmt.Errorf("failed to fetch repository metadata from %s: %w", repoURL, err)
	}
	
	// Update package with GitHub data
	pkg.Description = ghRepo.Description
	pkg.CreatedAt = ghRepo.CreatedAt
	pkg.UpdatedAt = ghRepo.UpdatedAt
	
	if ghRepo.License != nil {
		pkg.License = ghRepo.License.SPDXID
	}
	
	// Get author from GitHub if not provided
	if author == "" && ghRepo.Owner.Name != "" {
		pkg.Author = ghRepo.Owner.Name
	} else if author == "" {
		pkg.Author = ghRepo.Owner.Login
	}

	// Fetch latest version
	version, err := github.GetLatestVersion(owner, repo)
	if err != nil {
		fmt.Printf("Warning: Could not fetch version information: %v\n", err)
	} else if version != "" {
		pkg.Version = version
		fmt.Printf("Found version: %s\n", version)
	}

	// Create file path
	filePath := filepath.Join("content", fmt.Sprintf("%s.md", packageName))

	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		return fmt.Errorf("package file already exists: %s", filePath)
	}

	// Write package file
	if err := hugo.WritePackage(filePath, pkg); err != nil {
		return fmt.Errorf("failed to write package file: %w", err)
	}

	fmt.Printf("✓ Created %s\n", filePath)

	// Build site to validate
	fmt.Println("Validating site build...")
	cmd := exec.Command("hugo", "--gc", "--minify")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Warning: Site build validation failed: %v\n", err)
		fmt.Printf("Output: %s\n", output)
	} else {
		fmt.Println("✓ Site builds successfully")
	}

	// Print summary
	fmt.Println("\n=== Package Added Successfully ===")
	fmt.Printf("Name: %s\n", pkg.Title)
	fmt.Printf("Import Path: %s\n", pkg.ImportPath)
	fmt.Printf("Repository: %s\n", pkg.RepoURL)
	if pkg.Version != "" {
		fmt.Printf("Version: %s\n", pkg.Version)
	}
	if pkg.Description != "" {
		fmt.Printf("Description: %s\n", pkg.Description)
	}
	
	fmt.Println("\nNext steps:")
	fmt.Println("1. Review the generated file:", filePath)
	fmt.Println("2. Commit the changes: git add", filePath, "&& git commit -m \"Add", packageName, "package\"")
	fmt.Println("3. Push to deploy: git push")

	return nil
}