# Quzz CLI

A powerful CLI companion for [Quzz](https://github.com/onurhan1337/quzz) - the React Server Components debugging tool. Built with Go and designed for performance, featuring rich terminal UI and powerful filtering capabilities.

## Features

- **Trace Visualization**: Beautiful terminal output with statistics and detailed trace information
- **Powerful Filtering**: Filter by component, operation, duration, date range, and more
- **Interactive Setup**: Guided configuration wizard for quick project setup
- **Statistical Analysis**: Get insights into your application's performance
- **JSON Export**: Machine-readable output for programmatic processing
- **Rich Terminal UI**: Built with Charm libraries for an excellent CLI experience

## Installation

### Go Users

```bash
go install github.com/onurhan1337/quzz-cli@latest
```

### From Source

```bash
git clone https://github.com/onurhan1337/quzz-cli.git
cd quzz-cli
go build -o quzz
```

### Binaries

Download pre-built binaries from [Releases](https://github.com/onurhan1337/quzz-cli/releases).

## Quick Start

### 1. Initialize Your Project

Run the interactive setup wizard to create a configuration file:

```bash
quzz init
```

For TypeScript projects:

```bash
quzz init --typescript
```

Skip prompts and use defaults:

```bash
quzz init --skip-prompts
```

### 2. Visualize Traces

View all traces from your application:

```bash
quzz visualize traces.json
```

Show only statistics:

```bash
quzz visualize traces.json --stats
```

## Usage Examples

### Basic Visualization

```bash
# View all traces with statistics
quzz visualize traces.json

# Show only statistics
quzz visualize traces.json --stats

# Limit displayed traces
quzz visualize traces.json --limit 100
```

### Filtering

Filter by component:

```bash
quzz visualize traces.json --component UserProfile
```

Filter by duration range:

```bash
quzz visualize traces.json --min-duration 100 --max-duration 500
```

Filter by date range:

```bash
quzz visualize traces.json --start-date 2024-01-01T00:00:00Z --end-date 2024-01-31T23:59:59Z
```

Filter by log level:

```bash
quzz visualize traces.json --level error
```

Show only errors:

```bash
quzz visualize traces.json --errors
```

Show only performance warnings:

```bash
quzz visualize traces.json --warnings
```

Filter by operation:

```bash
quzz visualize traces.json --operation render
```

Use regex to filter components:

```bash
quzz visualize traces.json --component-regex "^(Blog|Product)"
```

### Combining Filters

```bash
quzz visualize traces.json \
  --component-regex "^Product" \
  --min-duration 200 \
  --level warn \
  --limit 20
```

### JSON Output

Export filtered results as JSON for further processing:

```bash
quzz visualize traces.json --json > filtered-traces.json
```

## Commands

### `quzz init`

Initialize Quzz configuration for your project with an interactive setup wizard.

**Flags:**
- `--skip-prompts`: Skip prompts and use default configuration
- `--typescript`: Generate TypeScript config file

**Example:**
```bash
quzz init --typescript
```

### `quzz visualize [path]`

Visualize and explore traces with powerful filtering options.

**Flags:**
- `-s, --stats`: Show only statistics
- `-j, --json`: Output as JSON
- `-c, --component string`: Filter by component name
- `-o, --operation string`: Filter by operation
- `--category string`: Filter by category
- `-l, --level string`: Filter by log level (info, warn, error, debug)
- `--min-duration int`: Filter by minimum duration (ms)
- `--max-duration int`: Filter by maximum duration (ms)
- `--start-date string`: Filter by start date (RFC3339 format)
- `--end-date string`: Filter by end date (RFC3339 format)
- `--errors`: Show only traces with errors
- `--warnings`: Show only traces with warnings
- `--component-regex string`: Filter components by regex pattern
- `--limit int`: Limit number of traces displayed (default 50, 0 for all)

**Example:**
```bash
quzz visualize traces.json --component UserProfile --min-duration 100
```

## Configuration

The `quzz.config.js` or `quzz.config.ts` file configures the Quzz npm package behavior. Example configuration:

```javascript
module.exports = {
  logLevel: "info",
  outputFormat: "compact",
  performance: {
    warnThreshold: 500,
  },
  componentFilter: /^(Blog|Product)/,
  sensitiveKeys: ["apiKey", "secretToken", "password"],
};
```

## Examples

Check the [examples](./examples) directory for:
- Sample `traces.json` file
- Example configuration files (JS and TS)

## Development

### Building

```bash
go build -o quzz
```

### Running Tests

```bash
go test ./...
```

### Project Structure

```
quzz-cli/
├── cmd/                    # Command implementations
│   ├── root.go            # Root command
│   ├── init.go            # Init command
│   └── visualize.go       # Visualize command
├── internal/              # Internal packages
│   ├── config/            # Configuration management
│   ├── trace/             # Trace loading and filtering
│   └── ui/                # Terminal UI components
├── examples/              # Example files
│   ├── traces.json
│   ├── quzz.config.js
│   └── quzz.config.ts
└── main.go               # Entry point
```

## Technologies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions for terminal output
- [tablewriter](https://github.com/olekukonko/tablewriter) - Table rendering
- [promptui](https://github.com/manifoldco/promptui) - Interactive prompts

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Related Projects

- [Quzz](https://github.com/onurhan1337/quzz) - The main Quzz npm package
- [Quzz Documentation](https://github.com/onurhan1337/quzz#readme)

## Support

- GitHub Issues: [Report a bug](https://github.com/onurhan1337/quzz-cli/issues)
- Quzz Main Project: [onurhan1337/quzz](https://github.com/onurhan1337/quzz)
