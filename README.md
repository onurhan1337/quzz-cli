# Quzz CLI

A powerful CLI for debugging React Server Components. Built with Go and Bubble Tea for a modern, interactive terminal experience.

## Features

- **Interactive TUI** - Full-screen terminal interface with smooth navigation
- **Keyboard Shortcuts** - Vim-style (j/k) and arrow keys
- **Smart Pagination** - Handles large trace files efficiently
- **Rich Filtering** - By component, duration, date, errors, and more
- **Statistics Dashboard** - Toggle detailed stats with Tab key
- **Zero-Config** - Works out of the box

## Installation

### Homebrew (macOS/Linux)
```bash
brew tap onurhan1337/tap
brew install quzz-cli
```

### Go Install
```bash
go install github.com/onurhan1337/quzz-cli@latest
```
Note: The binary will be named `quzz-cli`, so run commands as `quzz-cli` instead of `quzz`.

### Pre-built Binaries
Download for your platform from [Releases](https://github.com/onurhan1337/quzz-cli/releases).

**Linux/macOS:**
```bash
# Download and extract (replace VERSION and OS/ARCH)
curl -LO https://github.com/onurhan1337/quzz-cli/releases/download/v1.0.0/quzz-cli_Linux_x86_64.tar.gz
tar -xzf quzz-cli_Linux_x86_64.tar.gz
sudo mv quzz /usr/local/bin/
```

**Windows:**
Download the `.zip` file from releases and add `quzz.exe` to your PATH.

## Quick Start

Initialize configuration:
```bash
quzz init
```

Visualize traces:
```bash
quzz visualize traces.json
```

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| Up / k | Scroll up |
| Down / j | Scroll down |
| PgUp / b | Page up |
| PgDn / f | Page down |
| Home / g | Jump to top |
| End / G | Jump to bottom |
| Tab | Toggle statistics |
| ? | Show help |
| q / Esc | Quit |

## Commands

### `quzz init`

Interactive setup wizard for project configuration.

```bash
quzz init                 # Interactive mode
quzz init --typescript    # Generate TypeScript config
quzz init --skip-prompts  # Use defaults
```

### `quzz visualize [path]`

View and analyze traces with interactive TUI.

```bash
quzz visualize traces.json
```

**Filtering:**
```bash
quzz visualize traces.json --component UserProfile
quzz visualize traces.json --min-duration 100 --max-duration 500
quzz visualize traces.json --errors
quzz visualize traces.json --component-regex "^(Blog|Product)"
```

**Flags:**
- `-c, --component` - Filter by component name
- `-o, --operation` - Filter by operation
- `-l, --level` - Filter by log level (error, warn, info, debug)
- `--min-duration` - Minimum duration in ms
- `--max-duration` - Maximum duration in ms
- `--errors` - Show only errors
- `--warnings` - Show only warnings
- `--component-regex` - Filter by regex pattern
- `--limit` - Limit number of traces (default: 50)
- `-j, --json` - Output as JSON
- `-s, --stats` - Show only statistics

## Configuration

Generated config file (`quzz.config.js` or `quzz.config.ts`):

```javascript
module.exports = {
  logLevel: "info",
  outputFormat: "pretty",
  performance: {
    enabled: true,
    warnThreshold: 500,
  },
  componentFilter: /^(Blog|Product)/,
};
```

## Technologies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [Cobra](https://github.com/spf13/cobra) - CLI framework

## Development

```bash
go build -o quzz
go test ./...
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Links

- [Quzz npm package](https://github.com/onurhan1337/quzz)
- [Report Issues](https://github.com/onurhan1337/quzz-cli/issues)
