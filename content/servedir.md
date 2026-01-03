---
title: servedir
import_path: go.ngs.io/servedir
repo_url: https://github.com/ngs/servedir
description: Simple HTTP server that hosts local directory statically.
version: v1.0.2
documentation_url: https://pkg.go.dev/go.ngs.io/servedir
license: MIT
author: ngs
created_at: 2018-04-16T01:22:33Z
updated_at: 2025-08-12T12:17:30Z
---

# servedir

[![CI](https://github.com/ngs/servedir/actions/workflows/ci.yml/badge.svg)](https://github.com/ngs/servedir/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/go.ngs.io/servedir)](https://goreportcard.com/report/go.ngs.io/servedir)
[![Go Reference](https://pkg.go.dev/badge/go.ngs.io/servedir.svg)](https://pkg.go.dev/go.ngs.io/servedir)

A simple HTTP file server with Apache-style logging for serving static files from any directory.

## Features

- 🚀 Simple and lightweight HTTP file server
- 📝 Apache Combined Log Format output
- 🔧 Configurable port
- 🖥️ Auto-open browser (cross-platform)
- 📁 Serve any directory

## Installation

### Using Go Install

```sh
go install go.ngs.io/servedir@latest
```

### Download Binary

Download the latest binary from the [releases page](https://github.com/ngs/servedir/releases).

### Build from Source

```sh
git clone https://github.com/ngs/servedir.git
cd servedir
make build
```

## Usage

### Basic Usage

Serve current directory on any available port:

```sh
servedir
```

### Serve Specific Directory

```sh
servedir /path/to/directory
```

### Custom Port

```sh
servedir -port 8080 .
```

### Auto-open Browser

```sh
servedir -open -port 8123 .
```

### Command Line Options

```
-port int    HTTP Port to Listen (default 0, any available port)
-open        Open browser on started
```

## Development

### Prerequisites

- Go 1.21 or later
- Make (optional)

### Building

```sh
make build
```

### Testing

```sh
make test
```

### Running Locally

```sh
make run
```

### Available Make Commands

```sh
make help       # Show all available commands
make build      # Build the binary
make test       # Run tests
make coverage   # Generate coverage report
make lint       # Run golangci-lint
make fmt        # Format code
make clean      # Remove build artifacts
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Author

[Atsushi Nagase](https://ngs.io/)

## License

Copyright &copy; 2018 [Atsushi Nagase](https://ngs.io/). All rights reserved.

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
