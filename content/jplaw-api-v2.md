---
title: jplaw-api-v2
import_path: go.ngs.io/jplaw-api-v2
repo_url: https://github.com/ngs/go-jplaw-api-v2
description: Go Client library for Japanese Laws API Version 2
version: v0.0.3
documentation_url: https://pkg.go.dev/go.ngs.io/jplaw-api-v2
license: MIT
author: ngs
created_at: 2025-08-12T21:59:41Z
updated_at: 2025-08-12T23:40:52Z
---

# Japan Law API v2 Go Client Library

Go client library for the [Japan Law API v2] (法令API v2), automatically generated from the OpenAPI specification.

## Features

- Type-safe Go client for all Japan Law API v2 endpoints
- Automatic handling of path parameters and query parameters
- Custom date/time types for proper date format handling (YYYY-MM-DD)
- Support for both JSON and raw content (XML/HTML) responses
- Helper functions for creating pointer values
- Comprehensive type definitions for all API responses
- Full English documentation and comments

## Installation

```bash
go get go.ngs.io/jplaw-api-v2
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "log"
    
    lawapi "go.ngs.io/jplaw-api-v2"
)

func main() {
    // Create API client
    client := lawapi.NewClient()
    
    // Search for laws by title
    params := &lawapi.GetLawsParams{
        LawTitle: lawapi.StringPtr("電波法"),
        Limit:    lawapi.Int32Ptr(10),
    }
    
    result, err := client.GetLaws(params)
    if err != nil {
        log.Fatal(err)
    }
    
    // Process results
    for _, law := range result.Laws {
        if law.LawInfo != nil {
            fmt.Printf("Law ID: %s\n", law.LawInfo.LawId)
        }
        if law.RevisionInfo != nil {
            fmt.Printf("Law Title: %s\n", law.RevisionInfo.LawTitle)
        }
    }
}
```

## API Methods

### GetLaws
Search and retrieve laws based on various criteria.

```go
params := &lawapi.GetLawsParams{
    LawTitle:     lawapi.StringPtr("電波法"),
    LawType:      &[]lawapi.LawType{lawapi.LawTypeAct},
    Limit:        lawapi.Int32Ptr(100),
}
result, err := client.GetLaws(params)
```

### GetLawData
Retrieve full law data including the law text.

```go
lawID := "325AC0000000131"
params := &lawapi.GetLawDataParams{
    ResponseFormat: (*lawapi.ResponseFormat)(lawapi.StringPtr("json")),
}
lawData, err := client.GetLawData(lawID, params)
```

### GetLawFile
Retrieve law file in various formats (XML, JSON, HTML, RTF, DOCX).

```go
lawID := "325AC0000000131"
fileType := "xml" // or "json", "html", "rtf", "docx"
params := &lawapi.GetLawFileParams{}
content, err := client.GetLawFile(lawID, fileType, params)
```

### GetRevisions
Get revision history for a specific law.

```go
lawID := "325AC0000000131"
params := &lawapi.GetRevisionsParams{}
revisions, err := client.GetRevisions(lawID, params)
```

### GetKeyword
Search laws by keyword.

```go
params := &lawapi.GetKeywordParams{
    Keyword: "個人情報",
    Limit:   lawapi.Int32Ptr(50),
}
result, err := client.GetKeyword(params)
```

### GetAttachment
Retrieve attachments from law documents.

```go
revisionID := "325AC0000000131_20260527"
params := &lawapi.GetAttachmentParams{
    Src: lawapi.StringPtr("./pict/example.jpg"),
}
attachment, err := client.GetAttachment(revisionID, params)
```

## Code Generation

This library is automatically generated from the OpenAPI specification. To regenerate the client:

### Prerequisites

```bash
go get gopkg.in/yaml.v3
```

### Download OpenAPI Specification

```bash
wget https://laws.e-gov.go.jp/api/2/swagger-ui/lawapi-v2.yaml
```

### Build the Generator

```bash
cd cmd/clientgen
go build -o ../../clientgen
```

### Generate the Client

```bash
./clientgen -input lawapi-v2.yaml -output . -package lawapi
```

### Generator Options

- `-input`: Path to the OpenAPI specification file (default: "lawapi-v2.yaml")
- `-output`: Output directory for generated files (default: ".")
- `-package`: Package name for generated code (default: "lawapi")

## Project Structure

- `cmd/clientgen/` - Code generation tool
  - `main.go` - Entry point for the generator
  - `openapi.go` - OpenAPI specification structures
  - `generator.go` - Code generation logic
- `types.go` - Generated type definitions
- `client.go` - Generated HTTP client and API methods
- `example/` - Usage examples
  - `main.go` - Basic usage example
  - `test_path_params.go` - Path parameters example

## Date Handling

The library includes custom `Date` and `DateTime` types to handle the API's date formats:

- `Date`: Handles dates in "YYYY-MM-DD" format
- `DateTime`: Handles both RFC3339 and "YYYY-MM-DD" formats

## Helper Functions

The library provides helper functions for creating pointer values:

```go
lawapi.StringPtr("value")     // *string
lawapi.IntPtr(42)             // *int
lawapi.Int32Ptr(100)          // *int32
lawapi.Int64Ptr(1000)         // *int64
lawapi.BoolPtr(true)          // *bool
lawapi.Float32Ptr(3.14)       // *float32
lawapi.Float64Ptr(2.718)      // *float64
```

## Type Definitions

The library includes comprehensive type definitions for all API responses:

- `LawsResponse` - Response for law search
- `LawDataResponse` - Response for law data retrieval
- `KeywordResponse` - Response for keyword search
- `LawRevisionsResponse` - Response for revision history
- `LawInfo` - Basic law information
- `RevisionInfo` - Law revision information
- `LawItem` - Individual law entry
- `KeywordItem` - Keyword search result item
- And many more...

## Enumerations

Type-safe enumerations for various API parameters:

- `LawType` - Constitution, Act, CabinetOrder, etc.
- `LawNumEra` - Meiji, Taisho, Showa, Heisei, Reiwa
- `ResponseFormat` - json, xml
- `FileType` - xml, json, html, rtf, docx
- `RepealStatus` - None, Repeal, Expire, Suspend, LossOfEffectiveness
- And more...

## Error Handling

All API methods return an error as the second return value:

```go
result, err := client.GetLaws(params)
if err != nil {
    log.Printf("API error: %v", err)
    return
}
```

## Custom HTTP Client

You can provide a custom HTTP client for advanced configurations:

```go
client := lawapi.NewClient()

customHTTPClient := &http.Client{
    Timeout: 60 * time.Second,
    // Add custom transport, etc.
}

client.SetHTTPClient(customHTTPClient)
```

## License

This client library is generated from the public Japan Law API v2 specification. Please refer to the official API terms of use for usage guidelines.

## Contributing

Issues and pull requests are welcome. To regenerate the client after updating the generator:

1. Make changes to the generator in `cmd/clientgen/`
2. Rebuild the generator: `go build -o clientgen cmd/clientgen/*.go`
3. Regenerate the client: `./clientgen -input lawapi-v2.yaml -output . -package lawapi`
4. Test the changes with examples in `example/`

[Japan Law API v2]: https://laws.e-gov.go.jp/api/2/swagger-ui#/
