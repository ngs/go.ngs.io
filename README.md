# go.ngs.io

GitHub Pages site for hosting Go modules with custom import paths.

## Overview

This repository hosts the static HTML files that enable Go modules to be installed using the custom domain `go.ngs.io`.

For example:
```bash
go install go.ngs.io/freecal@latest
```

## How it works

When Go tools fetch a module with a custom import path, they:

1. Request `https://go.ngs.io/freecal?go-get=1`
2. Parse the HTML response for `<meta name="go-import">` tags
3. Use the repository URL specified in the meta tag to fetch the actual code

## Modules

### freecal

A command-line tool to find free time slots in your Google Calendar.

- **Import path:** `go.ngs.io/freecal`
- **Source:** https://github.com/ngs/freecal
- **Documentation:** https://pkg.go.dev/go.ngs.io/freecal

## Setup

This site is hosted on GitHub Pages with a custom domain.

### DNS Configuration

Add a CNAME record for `go.ngs.io` pointing to `ngs.github.io`.

### GitHub Pages Settings

1. Repository Settings > Pages
2. Source: Deploy from a branch
3. Branch: main / root
4. Custom domain: go.ngs.io

## Adding a New Module

To add a new Go module:

1. Create a directory with the module name
2. Add an `index.html` file with the appropriate meta tags:

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="go-import" content="go.ngs.io/MODULE_NAME git https://github.com/ngs/MODULE_NAME">
    <meta name="go-source" content="go.ngs.io/MODULE_NAME https://github.com/ngs/MODULE_NAME https://github.com/ngs/MODULE_NAME/tree/master{/dir} https://github.com/ngs/MODULE_NAME/blob/master{/dir}/{file}#L{line}">
</head>
</html>
```

3. Update the main `index.html` to list the new module
4. Commit and push the changes

## License

MIT License - See the individual module repositories for their respective licenses.

## Author

Atsushi Nagase - https://ngs.io