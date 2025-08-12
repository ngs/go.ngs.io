package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"go.ngs.io/internal/github"
	"go.ngs.io/internal/hugo"
)

type updateResult struct {
	name    string
	status  string
	message string
	err     error
}

func main() {
	var (
		dryRun        bool
		updateAuthor  bool
		updateMissing bool
		help          bool
	)

	pflag.BoolVar(&dryRun, "dry-run", false, "Show what would be updated without making changes")
	pflag.BoolVar(&updateAuthor, "update-author", false, "Also update author information from GitHub")
	pflag.BoolVar(&updateMissing, "update-missing", false, "Update timestamps to current date for repositories that return 404")
	pflag.BoolVarP(&help, "help", "h", false, "Show help message")
	pflag.Parse()

	if help {
		printUsage()
		os.Exit(0)
	}

	// Get specific packages from arguments or update all
	packages := pflag.Args()
	if err := updatePackages(packages, dryRun, updateAuthor, updateMissing); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func printUsage() {
	fmt.Println("Usage: update-packages [package-names...] [options]")
	fmt.Println("\nUpdate Go packages metadata from GitHub API")
	fmt.Println("\nIf no package names are provided, all packages will be updated.")
	fmt.Println("\nOptions:")
	pflag.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("  update-packages                    # Update all packages")
	fmt.Println("  update-packages freecal servedir   # Update specific packages")
	fmt.Println("  update-packages --dry-run          # Preview changes without updating")
	fmt.Println("  update-packages --update-missing   # Update timestamps for missing repos")
}

func updatePackages(specificPackages []string, dryRun, updateAuthor, updateMissing bool) error {
	fmt.Println("Updating packages from GitHub...")
	if dryRun {
		fmt.Println("(DRY RUN - no changes will be made)")
	}
	fmt.Println()

	// Get list of package files
	packageFiles, err := hugo.ListPackages("content")
	if err != nil {
		return fmt.Errorf("failed to list packages: %w", err)
	}

	// Filter packages if specific ones requested
	if len(specificPackages) > 0 {
		filtered := []string{}
		for _, file := range packageFiles {
			base := filepath.Base(file)
			name := strings.TrimSuffix(base, ".md")
			
			for _, requested := range specificPackages {
				if name == requested {
					filtered = append(filtered, file)
					break
				}
			}
		}
		packageFiles = filtered
		
		if len(packageFiles) == 0 {
			return fmt.Errorf("no matching packages found")
		}
	}

	// Process each package
	results := []updateResult{}
	updatedCount := 0
	skippedCount := 0
	errorCount := 0

	for _, filePath := range packageFiles {
		name := strings.TrimSuffix(filepath.Base(filePath), ".md")
		result := processPackage(filePath, name, dryRun, updateAuthor, updateMissing)
		results = append(results, result)

		// Print result immediately
		switch result.status {
		case "updated":
			fmt.Printf("✓ %s - %s\n", result.name, result.message)
			updatedCount++
		case "skipped":
			fmt.Printf("○ %s - %s\n", result.name, result.message)
			skippedCount++
		case "error":
			fmt.Printf("✗ %s - %s\n", result.name, result.message)
			errorCount++
		case "missing":
			if updateMissing {
				fmt.Printf("⚠ %s - %s\n", result.name, result.message)
				updatedCount++
			} else {
				fmt.Printf("○ %s - %s\n", result.name, result.message)
				skippedCount++
			}
		}
	}

	// Validate site build if not dry run and changes were made
	if !dryRun && updatedCount > 0 {
		fmt.Println("\nValidating site build...")
		cmd := exec.Command("hugo", "--gc", "--minify")
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Warning: Site build validation failed: %v\n", err)
			fmt.Printf("Output: %s\n", output)
		} else {
			fmt.Println("✓ Site builds successfully")
		}
	}

	// Print summary
	fmt.Printf("\nSummary: %d updated", updatedCount)
	if skippedCount > 0 {
		fmt.Printf(", %d skipped", skippedCount)
	}
	if errorCount > 0 {
		fmt.Printf(", %d errors", errorCount)
	}
	fmt.Println()

	// Return error if any packages failed
	if errorCount > 0 {
		return fmt.Errorf("%d packages failed to update", errorCount)
	}

	return nil
}

func processPackage(filePath, name string, dryRun, updateAuthor, updateMissing bool) updateResult {
	// Read existing package
	pkg, err := hugo.ReadPackage(filePath)
	if err != nil {
		return updateResult{
			name:    name,
			status:  "error",
			message: fmt.Sprintf("failed to read package: %v", err),
			err:     err,
		}
	}

	// Skip if no repo URL
	if pkg.RepoURL == "" {
		return updateResult{
			name:    name,
			status:  "skipped",
			message: "no repository URL",
		}
	}

	// Parse repository URL
	owner, repo, err := github.ParseRepoURL(pkg.RepoURL)
	if err != nil {
		return updateResult{
			name:    name,
			status:  "error",
			message: fmt.Sprintf("invalid repository URL: %v", err),
			err:     err,
		}
	}

	// Fetch GitHub metadata
	ghRepo, err := github.GetRepository(owner, repo)
	if err != nil {
		// Always return error for non-existent repositories
		return updateResult{
			name:    name,
			status:  "error",
			message: fmt.Sprintf("failed to fetch repository %s: %v", pkg.RepoURL, err),
			err:     err,
		}
	}

	// Check what needs updating
	changes := []string{}

	// Always update timestamps from GitHub
	if pkg.CreatedAt != ghRepo.CreatedAt {
		pkg.CreatedAt = ghRepo.CreatedAt
		changes = append(changes, "created_at")
	}
	if pkg.UpdatedAt != ghRepo.UpdatedAt {
		pkg.UpdatedAt = ghRepo.UpdatedAt
		changes = append(changes, "updated_at")
	}

	// Update description
	if pkg.Description != ghRepo.Description {
		pkg.Description = ghRepo.Description
		if ghRepo.Description == "" {
			changes = append(changes, "description cleared")
		} else {
			changes = append(changes, "description")
		}
	}

	// Update license
	if ghRepo.License != nil {
		if pkg.License != ghRepo.License.SPDXID {
			pkg.License = ghRepo.License.SPDXID
			changes = append(changes, "license")
		}
	}

	// Update author if requested
	if updateAuthor {
		newAuthor := ghRepo.Owner.Name
		if newAuthor == "" {
			newAuthor = ghRepo.Owner.Login
		}
		if pkg.Author != newAuthor {
			pkg.Author = newAuthor
			changes = append(changes, "author")
		}
	}

	// Fetch and update version
	version, err := github.GetLatestVersion(owner, repo)
	if err == nil && version != "" && pkg.Version != version {
		oldVersion := pkg.Version
		pkg.Version = version
		if oldVersion == "" {
			changes = append(changes, fmt.Sprintf("version: %s", version))
		} else {
			changes = append(changes, fmt.Sprintf("version: %s → %s", oldVersion, version))
		}
	}

	// Check if any changes were made
	if len(changes) == 0 {
		return updateResult{
			name:    name,
			status:  "skipped",
			message: "already up to date",
		}
	}

	// Write changes if not dry run
	if !dryRun {
		if err := hugo.WritePackage(filePath, pkg); err != nil {
			return updateResult{
				name:    name,
				status:  "error",
				message: fmt.Sprintf("failed to write package: %v", err),
				err:     err,
			}
		}
	}

	return updateResult{
		name:    name,
		status:  "updated",
		message: strings.Join(changes, ", "),
	}
}