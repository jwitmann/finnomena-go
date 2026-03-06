# Contributing to finnomena-go

Thank you for your interest in contributing! This document provides guidelines for contributing to this project.

## Code of Conduct

Be respectful and constructive in all interactions.

## How to Contribute

### Reporting Bugs

If you find a bug, please open an issue with:
- A clear title and description
- Steps to reproduce the bug
- Expected vs actual behavior
- Go version and environment details
- Any error messages or logs

### Suggesting Enhancements

Enhancement suggestions are welcome! Please open an issue with:
- A clear use case
- Description of the proposed enhancement
- Any potential implementation approaches

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** with clear, focused commits
3. **Add tests** for any new functionality
4. **Ensure all tests pass**: `go test ./...`
5. **Run the linter**: `golangci-lint run`
6. **Update documentation** if needed
7. **Submit a pull request** with a clear description

## Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/finnomena-go.git
cd finnomena-go

# Install dependencies
go mod download

# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## Code Style

- Follow standard Go conventions (gofmt, go vet)
- Write clear, idiomatic Go code
- Add comments for exported functions and types
- Keep functions focused and small
- Use meaningful variable names

## Testing

- All new features must include tests
- Aim for high test coverage
- Tests should be clear and serve as documentation
- Use table-driven tests where appropriate

## Commit Messages

- Use clear, descriptive commit messages
- Start with a verb (Add, Fix, Update, Remove)
- Reference issues when applicable: "Fix #123"

Example:
```
Add retry logic for 429 rate limit responses

Implements exponential backoff with jitter to handle
rate limiting from the Finnomena API.

Fixes #45
```

## Questions?

Feel free to open an issue for any questions or discussion.

Thank you for contributing!
