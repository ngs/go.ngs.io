---
title: freecal
import_path: go.ngs.io/freecal
repo_url: https://github.com/ngs/freecal
description: A command-line tool to find free time slots in your Google Calendar.
version: v1.0.0
documentation_url: https://pkg.go.dev/go.ngs.io/freecal
license: MIT
author: Atsushi Nagase
created_at: 2025-08-12T07:37:35Z
updated_at: 2025-08-12T12:04:11Z
---

# FreeCal

A command-line tool to find free time slots in your Google Calendar. It fetches events from Google Calendar and outputs available time slots during business hours in Markdown format.

## Quick Install

```bash
go install go.ngs.io/freecal@latest
```

## Features

- Fetches events from Google Calendar using OAuth2 authentication
- Finds free time slots during configurable business hours
- Filters out weekends automatically
- Supports minimum duration filtering for free slots
- Outputs results in Markdown format with Japanese weekday names
- Automatic browser-based OAuth authentication flow

## Prerequisites

- Go 1.23 or higher
- Google Cloud Console account
- Google Calendar API enabled

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/ngs/freecal.git
cd freecal
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Build the application

```bash
go build -o freecal main.go
```

## Setup

### 1. Enable Google Calendar API

1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google Calendar API:
   - Navigate to "APIs & Services" > "Library"
   - Search for "Google Calendar API"
   - Click on it and press "Enable"

### 2. Create OAuth 2.0 credentials

1. In Google Cloud Console, go to "APIs & Services" > "Credentials"
2. Click "Create Credentials" > "OAuth client ID"
3. If prompted, configure the OAuth consent screen first:
   - Choose "External" user type (or "Internal" if using Google Workspace)
   - Fill in the required fields
   - Add your email to test users if using "External" type
4. For Application type, select "Desktop app"
5. Give it a name (e.g., "FreeCal")
6. Click "Create"
7. Download the credentials JSON file
8. Save it as `credentials.json` in the project directory

### 3. First-time authentication

When you run the application for the first time, it will:
1. Start a local server to handle OAuth callback
2. Open your default browser for Google authentication
3. Ask you to authorize access to your Google Calendar
4. Save the authentication token locally for future use

## Usage

### Basic usage

```bash
./freecal -credentials ./credentials.json -start 2025-01-13 -end 2025-01-17
```

### Full command with all options

```bash
./freecal \
  -credentials ./credentials.json \
  -token ./token.json \
  -calendar primary \
  -start 2025-01-13 \
  -end 2025-01-17 \
  -workstart 09:00 \
  -workend 17:00 \
  -min 60 \
  -tz Asia/Tokyo
```

### Command-line options

| Option | Description | Default |
|--------|-------------|---------|
| `-credentials` | Path to OAuth client credentials JSON file | (required) |
| `-token` | Path to save/load OAuth token | `token.json` |
| `-calendar` | Calendar ID (use "primary" for your main calendar) | `primary` |
| `-start` | Start date in YYYY-MM-DD format | (required) |
| `-end` | End date in YYYY-MM-DD format | (required) |
| `-workstart` | Business hours start time (HH:MM) | `09:00` |
| `-workend` | Business hours end time (HH:MM) | `17:00` |
| `-min` | Minimum free slot duration in minutes | `60` |
| `-tz` | IANA timezone (e.g., Asia/Tokyo, America/New_York) | `Asia/Tokyo` |

## Example output

```markdown
- 2025-01-13（月） 09:00~10:00, 14:00~15:30
- 2025-01-14（火） 10:30~12:00, 13:00~17:00
- 2025-01-15（水） 09:00~11:00
- 2025-01-16（木） 15:00~17:00
- 2025-01-17（金） 09:00~12:00, 14:00~17:00
```

## Security notes

- Never commit `credentials.json` or `token.json` to version control
- The `.gitignore` file is configured to exclude these sensitive files
- Tokens are stored locally and are specific to your machine

## Troubleshooting

### Browser doesn't open automatically

If the browser doesn't open automatically during authentication, manually copy and paste the URL shown in the terminal into your browser.

### Permission denied error

Make sure the built binary has execute permissions:

```bash
chmod +x freecal
```

### Token expired

If you get authentication errors after the tool was working, delete `token.json` and re-authenticate:

```bash
rm token.json
./freecal -credentials ./credentials.json -start 2025-01-13 -end 2025-01-17
```

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Contributing

Contributions are welcome! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.