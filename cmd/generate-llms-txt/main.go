package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/spf13/pflag"
	"go.ngs.io/internal/hugo"
)

const llmsTxtTemplate = `# go.ngs.io

> Go Module Vanity Import Path Service

go.ngs.io provides custom vanity import paths for Go modules developed by Atsushi Nagase.
You can install packages using short, memorable URLs like ` + "`go.ngs.io/package`" + `.

## Installation

` + "```bash" + `
go install go.ngs.io/<package-name>@latest
` + "```" + `

## Available Packages

{{range .Packages}}### {{.Title}}
{{if .Description}}
{{.Description}}
{{end}}
- **Import**: ` + "`{{.ImportPath}}`" + `{{if .Version}}
- **Version**: {{.Version}}{{end}}{{if .License}}
- **License**: {{.License}}{{end}}{{if .DocumentationURL}}
- **Documentation**: {{.DocumentationURL}}{{end}}{{if .RepoURL}}
- **Repository**: {{.RepoURL}}{{end}}

` + "```bash" + `
go install {{.ImportPath}}@latest
` + "```" + `

{{end}}`

type templateData struct {
	Packages []*hugo.Package
}

func main() {
	var (
		outputFile string
		help       bool
	)

	pflag.StringVarP(&outputFile, "output", "o", "", "Output file path (default: stdout)")
	pflag.BoolVarP(&help, "help", "h", false, "Show help message")
	pflag.Parse()

	if help {
		printUsage()
		os.Exit(0)
	}

	if err := generateLLMsTxt(outputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: generate-llms-txt [options]")
	fmt.Println("\nGenerate llms.txt file from package data")
	fmt.Println("\nOptions:")
	pflag.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("  generate-llms-txt                    # Output to stdout")
	fmt.Println("  generate-llms-txt -o llms.txt        # Output to file")
	fmt.Println("  generate-llms-txt -o public/llms.txt # Output to public directory")
}

func generateLLMsTxt(outputFile string) error {
	// Get list of packages
	packageFiles, err := hugo.ListPackages("content")
	if err != nil {
		return fmt.Errorf("failed to list packages: %w", err)
	}

	// Read all packages
	packages := make([]*hugo.Package, 0, len(packageFiles))
	for _, filePath := range packageFiles {
		pkg, err := hugo.ReadPackage(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to read %s: %v\n", filePath, err)
			continue
		}
		packages = append(packages, pkg)
	}

	// Sort packages by title
	sort.Slice(packages, func(i, j int) bool {
		return strings.ToLower(packages[i].Title) < strings.ToLower(packages[j].Title)
	})

	// Create template
	tmpl, err := template.New("llms.txt").Parse(llmsTxtTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// Prepare template data
	data := templateData{
		Packages: packages,
	}

	// Determine output destination
	var output *os.File
	if outputFile == "" {
		output = os.Stdout
	} else {
		// Ensure directory exists
		dir := filepath.Dir(outputFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		output, err = os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer output.Close()
	}

	// Execute template
	if err := tmpl.Execute(output, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if outputFile != "" {
		fmt.Fprintf(os.Stderr, "Generated %s\n", outputFile)
	}

	return nil
}
