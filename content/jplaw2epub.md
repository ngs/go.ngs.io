---
title: jplaw2epub
import_path: go.ngs.io/jplaw2epub
repo_url: https://github.com/ngs/jplaw2epub
description: Go library and command line tool to convert Japanese Standard Law XML Schema (法令標準XMLスキーマ) into EPUB files.
version: v1.2.0
documentation_url: https://pkg.go.dev/go.ngs.io/jplaw2epub
license: MIT
author: ngs
created_at: 2025-02-09T00:46:45Z
updated_at: 2025-08-20T08:25:49Z
---

# jplaw2epub

Go library and command line tool to convert [Japanese Standard Law XML Schema (法令標準XMLスキーマ)][xmldoc] into EPUB files.

## Overview

jplaw2epub is a comprehensive tool for converting Japanese law documents from the official XML format to EPUB e-books. It supports the full Japanese Standard Law XML Schema specification, including all document structures, appendixes, tables, and figures.

## Library Usage

```go
package main

import (
    "fmt"
    "os"
    
    "go.ngs.io/jplaw2epub"
)

func main() {
    // Method 1: Create EPUB from XML file path
    book, err := jplaw2epub.CreateEPUBFromXMLPath("law.xml")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    if err := jplaw2epub.WriteEPUB(book, "output.epub"); err != nil {
        fmt.Printf("Error writing EPUB: %v\n", err)
        return
    }

    // Method 2: Create EPUB from io.Reader
    xmlFile, err := os.Open("law.xml")
    if err != nil {
        fmt.Printf("Error opening file: %v\n", err)
        return
    }
    defer xmlFile.Close()

    book2, err := jplaw2epub.CreateEPUBFromXMLFile(xmlFile)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    if err := jplaw2epub.WriteEPUB(book2, "output2.epub"); err != nil {
        fmt.Printf("Error writing EPUB: %v\n", err)
        return
    }

    fmt.Println("EPUB files created successfully!")
}
```

### API Functions

- `CreateEPUBFromXMLPath(xmlPath string) (*epub.Epub, error)` - Creates an EPUB from a file path
- `CreateEPUBFromXMLFile(xmlFile io.Reader) (*epub.Epub, error)` - Creates an EPUB from an io.Reader
- `WriteEPUB(book *epub.Epub, destPath string) error` - Writes an EPUB book to a file

## Command Line Usage

### Installation

Install the CLI tool:

```sh
go install go.ngs.io/jplaw2epub/cmd/jplaw2epub@latest
```

Or build from source:

```sh
git clone https://github.com/ngs/jplaw2epub.git
cd jplaw2epub
go build -o jplaw2epub ./cmd/jplaw2epub
```

### Basic Usage

```sh
jplaw2epub -d output.epub input.xml
```

### Command Line Options

```
-d string
    Destination file path (required)
-no-images
    Skip downloading and embedding images
-max-image-height string
    Maximum image height (e.g., '300px', '80vh', '50%') (default "80vh")
```

### Examples

Convert a law XML file to EPUB:
```sh
jplaw2epub -d mylaw.epub path/to/law.xml
```

Convert without downloading images:
```sh
jplaw2epub -no-images -d mylaw.epub path/to/law.xml
```

Convert with custom image height limit:
```sh
jplaw2epub -max-image-height "500px" -d mylaw.epub path/to/law.xml
```

## Installation as Go Library

Add to your Go project:

```sh
go get go.ngs.io/jplaw2epub
```

## Features

### Document Structure Support
- **Main Provisions**: Chapters, sections, subsections, divisions, and articles
- **Supplementary Provisions**: Full support with chapters, articles, and appendixes
- **Paragraph Hierarchy**: Proper handling of numbered and unnumbered paragraphs
- **Item Structure**: Support for Items, Subitem1, Subitem2, and Subitem3
- **List Elements**: Native list support with proper nesting (List, Sublist1-3)

### Appendix Support
- **AppdxTable**: Appendix tables with full table structure
- **AppdxNote**: Appendix notes with structured content
- **AppdxStyle**: Appendix styles with formatting
- **AppdxFormat**: Appendix formats with embedded figures
- **AppdxFig**: Appendix figures with image support

### Advanced Features
- **Ruby Annotations**: Full support for Japanese phonetic guides (ルビ)
- **Table Processing**: Complex tables with headers, spans, and borders
- **Image Embedding**: Automatic download and embedding of referenced images
- **Figure Support**: FigStruct and Fig element processing
- **Style Management**: StyleStruct and Format element handling
- **Dynamic List Styling**: Automatic detection (CJK ideographic, katakana-iroha, hiragana-iroha)

### Technical Features
- **High Test Coverage**: 71.5% code coverage with comprehensive test suite
- **Clean Code**: Passes all Go linters (gofmt, go vet, golangci-lint)
- **Modular Architecture**: Well-organized code with separate processors for each element type
- **Error Handling**: Robust error handling throughout the conversion process
- **EPUB Metadata**: Full support with Japanese era information
- **Clean Go API**: Simple library interface for programmatic usage
- **Cross-platform CLI**: Command-line tool with binary releases

## Data Sources

The tool processes XML files conforming to the Japanese Standard Law XML Schema. You can obtain these files from:

- [e-Gov Laws Database (e-Gov法令検索)](https://laws.e-gov.go.jp/)
- [e-Gov Law API](https://laws.e-gov.go.jp/api.html)

## Requirements

- Go 1.18 or higher (for building from source)
- XML files conforming to Japanese Standard Law XML Schema

## Testing

Run the test suite:

```sh
go test -v ./...
```

Generate coverage report:

```sh
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

Current test coverage: **71.5%**

## Development

### Building from Source

```sh
git clone https://github.com/ngs/jplaw2epub.git
cd jplaw2epub
go build -o jplaw2epub ./cmd/jplaw2epub
```

### Running Linters

```sh
# Format code
gofmt -w .

# Check for issues
go vet ./...

# Run comprehensive linting (requires golangci-lint)
golangci-lint run ./...

# Clean up dependencies
go mod tidy
```

### Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure:
- All tests pass (`go test ./...`)
- Code coverage remains above 70%
- All linters pass
- New features include appropriate tests

## Author

[Atsushi Nagase]

## License

MIT License

Copyright &copy; 2025 [Atsushi Nagase]. See [LICENSE.md](LICENSE.md) for details.

[Atsushi Nagase]: https://ngs.io/
[xmldoc]: https://laws.e-gov.go.jp/docs/law-data-basic/419a603-xml-schema-for-japanese-law/