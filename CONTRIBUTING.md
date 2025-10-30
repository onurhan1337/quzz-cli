# Contributing to Quzz CLI

Thank you for considering contributing to Quzz CLI! This document provides guidelines and instructions for contributing.

## Getting Started

### Prerequisites

- Go 1.24 or later
- Git

### Setting Up Development Environment

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/quzz-cli.git
   cd quzz-cli
   ```
3. Add upstream remote:
   ```bash
   git remote add upstream https://github.com/onurhan1337/quzz-cli.git
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```

## Development Workflow

### Building

```bash
make build
```

### Running

```bash
make run
# or
./quzz visualize examples/traces.json
```

### Testing

```bash
make test
```

### Code Style

- Follow standard Go conventions and idioms
- Use `gofmt` to format your code
- Run `go vet` to check for common mistakes
- Add comments for exported functions and types

### Project Structure

```
quzz-cli/
├── cmd/                    # Command implementations
│   ├── root.go            # Root command
│   ├── init.go            # Init command
│   ├── visualize.go       # Visualize command
│   └── version.go         # Version command
├── internal/              # Internal packages
│   ├── config/            # Configuration management
│   │   ├── types.go
│   │   └── generator.go
│   ├── trace/             # Trace loading and filtering
│   │   ├── types.go
│   │   └── loader.go
│   └── ui/                # Terminal UI components
│       ├── styles.go
│       └── table.go
├── examples/              # Example files
└── main.go               # Entry point
```

## Making Changes

### Creating a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

### Commit Messages

Follow conventional commit format:

- `feat: add new feature`
- `fix: resolve bug`
- `docs: update documentation`
- `style: format code`
- `refactor: restructure code`
- `test: add tests`
- `chore: update dependencies`

### Example:

```bash
git commit -m "feat: add support for filtering by custom fields"
```

## Submitting Changes

### Pull Request Process

1. Update your fork with the latest upstream changes:
   ```bash
   git fetch upstream
   git rebase upstream/master
   ```

2. Push your changes:
   ```bash
   git push origin your-branch-name
   ```

3. Create a Pull Request on GitHub with:
   - Clear title describing the change
   - Detailed description of what changed and why
   - Reference any related issues
   - Screenshots for UI changes (if applicable)

4. Respond to review feedback

### Pull Request Guidelines

- Keep PRs focused on a single feature or fix
- Update documentation if needed
- Add tests for new features
- Ensure all tests pass
- Keep commits clean and atomic

## Adding New Features

### Adding a New Command

1. Create a new file in `cmd/` (e.g., `cmd/mycommand.go`)
2. Implement the command using Cobra patterns
3. Register the command in `init()` function
4. Add tests and documentation
5. Update README.md with usage examples

### Adding Internal Packages

1. Create package in `internal/` directory
2. Follow Go package naming conventions
3. Export only necessary types and functions
4. Add package documentation
5. Write unit tests

## Testing

### Running Tests

```bash
go test ./...
```

### Writing Tests

- Place test files next to source files (`*_test.go`)
- Use table-driven tests where appropriate
- Aim for good test coverage
- Test edge cases and error conditions

Example:

```go
func TestFilterTraces(t *testing.T) {
    tests := []struct {
        name     string
        traces   []TraceEntry
        opts     FilterOptions
        expected int
    }{
        // test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

## Documentation

### Code Documentation

- Document all exported functions, types, and packages
- Use complete sentences in comments
- Provide examples in doc comments when helpful

### README Updates

- Add new features to the feature list
- Update usage examples
- Add new command documentation

## Release Process

Releases are automated through GitHub Actions:

1. Ensure all tests pass
2. Update version in relevant files
3. Create and push a tag:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
4. GitHub Actions will automatically build and create a release

## Getting Help

- Open an issue for bugs or feature requests
- Join discussions in existing issues
- Reach out to maintainers for guidance

## Code of Conduct

### Our Standards

- Be respectful and inclusive
- Welcome newcomers
- Accept constructive criticism gracefully
- Focus on what's best for the community

### Unacceptable Behavior

- Harassment or discriminatory language
- Trolling or insulting comments
- Personal or political attacks
- Publishing others' private information

## Recognition

Contributors will be recognized in:
- GitHub contributors list
- Release notes
- Project documentation

Thank you for contributing to Quzz CLI!
