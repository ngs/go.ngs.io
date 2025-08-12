# Add Go Package to go.ngs.io

Add a new Go package to the go.ngs.io custom import path site. This command automatically fetches package information from GitHub API and creates a new content file with the appropriate frontmatter.

## Usage

Simply provide the repository URL:

```
/add-package https://github.com/username/repository-name
```

Or override auto-detected values:

```
/add-package https://github.com/username/repository-name --import-path go.ngs.io/custom-name
```

## What this command does

1. Parses the GitHub repository URL to extract owner and repository name
2. Fetches repository metadata using GitHub API:
   - Description
   - Default branch
   - Topics/tags
   - License information
   - Latest release version
   - Repository creation and update dates
3. Fetches owner information for author name
4. Automatically determines the import path (removes 'go-' prefix if present)
5. Creates a new content file in `content/` directory with all metadata
6. Configures go-import and go-source meta tags for the package

## Parameters

- **Repository URL** (required): The GitHub repository URL for the Go package
- **--import-path** (optional): Override the auto-generated import path
- **--version** (optional): Override the auto-detected version (from latest release)
- **--description** (optional): Override the repository description from GitHub
- **--author** (optional): Override the auto-detected author name

## Examples

### Basic usage
```
/add-package https://github.com/ngs/go-freecal
```

This creates `content/go-freecal.md` with import path `go.ngs.io/go-freecal`

### With custom import path
```
/add-package https://github.com/ngs/go-freecal --import-path go.ngs.io/freecal
```

This creates `content/freecal.md` with import path `go.ngs.io/freecal`

### With full details
```
/add-package https://github.com/ngs/go-freecal \
  --import-path go.ngs.io/freecal \
  --version v2.1.0 \
  --description "Go client library for FreeCal calendar API" \
  --license MIT \
  --author "Atsushi Nagase"
```

## Notes

- The package name in the content file will be derived from the import path
- The command will check if a file already exists to avoid overwriting
- After adding a package, commit the changes and push to trigger the GitHub Actions deployment
- The package will be accessible at `https://go.ngs.io/[package-name]/`

## Implementation Instructions

When implementing this command:

1. **Parse the GitHub repository URL**
   - Extract owner and repository name from the URL
   - Validate that it's a valid GitHub URL

2. **Fetch repository information using GitHub API**
   ```bash
   # Use gh api or curl to fetch repo data
   gh api repos/{owner}/{repo}
   ```
   Extract:
   - `description`: Repository description
   - `default_branch`: For go-source meta tag
   - `topics`: Use as tags
   - `license.spdx_id`: License type (e.g., MIT, Apache-2.0)
   - `created_at`: Repository creation date
   - `updated_at`: Last update date

3. **Fetch latest release version**
   ```bash
   gh api repos/{owner}/{repo}/releases/latest
   ```
   Extract `tag_name` as the version. If no releases, default to "v0.0.0" or check tags.

4. **Fetch owner information**
   ```bash
   gh api users/{owner}
   ```
   Extract `name` as the author. If not available, use `login` (username).

5. **Determine package name and import path**
   - If `--import-path` is provided, use it
   - Otherwise, generate from repository name:
     - Remove 'go-' prefix if present (e.g., 'go-freecal' → 'freecal')
     - Use the result as package name
     - Full import path: `go.ngs.io/{package-name}`

6. **Create the content file**
   - File path: `content/{package-name}.md`
   - Check if file already exists to avoid overwriting
   - Generate frontmatter with all collected data:
   ```yaml
   ---
   title: "{package-name}"
   import_path: "go.ngs.io/{package-name}"
   repo_url: "{original-github-url}"
   description: "{fetched-description}"
   version: "{fetched-version}"
   tags: [fetched-topics-array]
   documentation_url: "https://pkg.go.dev/go.ngs.io/{package-name}"
   license: "{fetched-license}"
   author: "{fetched-author-name}"
   created_at: "{fetched-created-date}"
   updated_at: "{fetched-updated-date}"
   ---
   ```

7. **Validate and build**
   - Run `hugo --gc --minify` to ensure the site builds correctly
   - Report any errors to the user

8. **Provide feedback**
   - Confirm successful creation with file path
   - Show the import path that will be available
   - Remind user to commit and push changes