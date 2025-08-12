# Update All Packages Metadata

Update all existing Go packages in go.ngs.io with the latest metadata from GitHub API. This command fetches current repository information including descriptions, timestamps, versions, and other metadata for all packages.

## Usage

```
/update-packages
```

Or update specific packages only:

```
/update-packages freecal servedir
```

## What this command does

1. Scans the `content/` directory for all existing package markdown files
2. For each package:
   - Extracts the GitHub repository URL from frontmatter
   - Fetches latest repository metadata via GitHub API
   - Fetches latest release/tag information
   - Updates the frontmatter with current data
   - Preserves custom fields not from GitHub
3. Validates the site builds correctly after updates
4. Reports which packages were updated and any errors

## Updated Fields

The following fields are updated from GitHub API:
- **description**: Repository description
- **version**: Latest release tag (if available)
- **license**: License type (e.g., MIT, Apache-2.0)
- **created_at**: Repository creation date
- **updated_at**: Last repository update date
- **topics**: Repository topics as tags (if any)

## Preserved Fields

The following fields are NOT modified:
- **title**: Package name
- **import_path**: Custom import path
- **repo_url**: Original repository URL
- **documentation_url**: Package documentation link
- **author**: Author name (unless --update-author flag is used)
- Any custom fields added manually

## Options

- **--dry-run**: Show what would be updated without making changes
- **--update-author**: Also update author information from GitHub
- **--parallel**: Update multiple packages in parallel (faster but may hit rate limits)

## Examples

### Update all packages
```
/update-packages
```

### Update specific packages
```
/update-packages freecal servedir
```

### Dry run to preview changes
```
/update-packages --dry-run
```

### Update including author information
```
/update-packages --update-author
```

## Implementation Instructions

When implementing this command:

1. **Scan content directory**
   ```bash
   # Find all .md files in content/
   find content/ -name "*.md" -type f
   ```

2. **For each markdown file:**
   - Read the file and parse frontmatter
   - Extract `repo_url` field
   - Parse GitHub owner and repo from URL

3. **Fetch repository data**
   ```bash
   gh api repos/{owner}/{repo}
   ```
   Extract updated fields:
   - description
   - license.spdx_id
   - created_at
   - updated_at
   - topics (for tags)
   - default_branch

4. **Fetch latest version**
   ```bash
   # Try releases first
   gh api repos/{owner}/{repo}/releases/latest
   # If no releases, check tags
   gh api repos/{owner}/{repo}/tags | head -1
   ```

5. **Update frontmatter**
   - Parse existing frontmatter
   - Update only the fields from GitHub
   - Preserve all other fields
   - Write back to file

6. **Batch processing considerations**
   - Process files sequentially by default to avoid rate limits
   - With --parallel flag, process up to 5 files concurrently
   - Add delay between API calls if needed
   - Handle API rate limit errors gracefully

7. **Error handling**
   - Skip files without repo_url
   - Log errors for repositories that can't be accessed
   - Continue processing other files on error
   - Report summary at the end

8. **Validation**
   - After all updates, run `hugo --gc --minify`
   - Report any build errors

9. **Output format**
   ```
   Updating packages from GitHub...
   
   ✓ freecal - description updated, version: v2.1.0 → v2.2.0
   ✓ servedir - timestamps updated
   ✗ example - Error: repository not found
   
   Summary: 2 updated, 1 error, 0 skipped
   ```

## Notes

- GitHub API rate limits: 60 requests/hour unauthenticated, 5000/hour authenticated
- The `gh` CLI tool must be authenticated for private repositories
- Large repositories with many releases may take longer to process
- Always commit changes before running this command for easy rollback
- Consider running with --dry-run first to preview changes

## Error Recovery

If the command fails partway through:
1. Files already updated will have their changes preserved
2. Run the command again to continue updating remaining files
3. Use git diff to review changes before committing
4. Use git checkout to revert unwanted changes