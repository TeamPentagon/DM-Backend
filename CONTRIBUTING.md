# Contributing to DM-Backend

Thank you for your interest in contributing to DM-Backend! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Process](#development-process)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for everyone.

## Getting Started

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/DM-Backend.git
   cd DM-Backend
   ```

3. Add the upstream remote:
   ```bash
   git remote add upstream https://github.com/TeamPentagon/DM-Backend.git
   ```

### Set Up Development Environment

1. Install Go 1.22.2 or higher

2. Install required tools:
   ```bash
   # Protocol Buffers compiler
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   
   # Linting tools (optional but recommended)
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Generate Protocol Buffer files:
   ```bash
   cd internal/model
   protoc --go_out=. --go_opt=paths=source_relative *.proto
   ```

## Development Process

### Branching Strategy

- `main` - Production-ready code
- `develop` - Integration branch for features
- `feature/*` - New features
- `bugfix/*` - Bug fixes
- `hotfix/*` - Urgent production fixes

### Creating a Feature Branch

```bash
# Sync with upstream
git fetch upstream
git checkout main
git merge upstream/main

# Create feature branch
git checkout -b feature/your-feature-name
```

## Pull Request Process

### Before Submitting

1. **Run tests:**
   ```bash
   go test ./... -v
   ```

2. **Run linter:**
   ```bash
   go vet ./...
   golangci-lint run  # if installed
   ```

3. **Format code:**
   ```bash
   go fmt ./...
   ```

4. **Update documentation** if you've changed APIs or added features

### Submitting a PR

1. Push your branch to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. Create a Pull Request on GitHub

3. Fill in the PR template with:
   - Description of changes
   - Related issues
   - Testing performed
   - Screenshots (if UI changes)

### PR Requirements

- [ ] All tests pass
- [ ] Code is formatted with `go fmt`
- [ ] No linting errors
- [ ] Documentation updated
- [ ] Meaningful commit messages

## Coding Standards

### Go Style Guide

Follow the official [Effective Go](https://golang.org/doc/effective_go) guidelines.

### Naming Conventions

```go
// Package names: lowercase, single word
package database

// Exported functions/types: PascalCase
func CreateDatabase() {}
type UserModel struct {}

// Unexported: camelCase
func validateInput() {}
var internalCounter int

// Constants: PascalCase for exported, camelCase for unexported
const MaxRetries = 3
const defaultTimeout = 30
```

### Error Handling

```go
// Always check and handle errors
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// Define custom errors for packages
var (
    ErrNotFound = errors.New("resource not found")
    ErrInvalid  = errors.New("invalid input")
)
```

### Comments

```go
// Package database provides utilities for database operations.
package database

// CreateConnection establishes a new database connection.
// It returns a connection handle and any error encountered.
//
// Example:
//
//	conn, err := CreateConnection(config)
//	if err != nil {
//	    log.Fatal(err)
//	}
func CreateConnection(config Config) (*Connection, error) {
    // Implementation
}
```

### File Organization

```go
// 1. Package clause
package mypackage

// 2. Imports (grouped: stdlib, external, internal)
import (
    "errors"
    "fmt"
    
    "github.com/gin-gonic/gin"
    
    "github.com/TeamPentagon/DM-Backend/internal/database"
)

// 3. Constants
const (
    MaxSize = 100
)

// 4. Variables
var (
    ErrInvalid = errors.New("invalid")
)

// 5. Types
type MyStruct struct {}

// 6. Functions
func MyFunction() {}
```

## Testing Guidelines

### Test File Naming

- Unit tests: `filename_test.go` in the same package
- Integration tests: in the `test/` directory

### Writing Tests

```go
func TestFunctionName(t *testing.T) {
    // Arrange
    input := "test"
    expected := "result"
    
    // Act
    result := FunctionName(input)
    
    // Assert
    if result != expected {
        t.Errorf("FunctionName(%s) = %s, want %s", input, result, expected)
    }
}
```

### Table-Driven Tests

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid input", "test@example.com", false},
        {"invalid email", "invalid", true},
        {"empty input", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := Validate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Test Coverage

Aim for at least 80% test coverage on new code:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

## Documentation

### Code Documentation

- Add package-level documentation in a `doc.go` file or at the top of the main package file
- Document all exported types, functions, and methods
- Include examples for complex functionality

### README Updates

Update the README when:
- Adding new features
- Changing API endpoints
- Modifying configuration options
- Adding new dependencies

### Changelog

Add entries to CHANGELOG.md for:
- New features
- Bug fixes
- Breaking changes
- Deprecations

## Questions?

If you have questions, feel free to:
- Open an issue for discussion
- Reach out to maintainers

Thank you for contributing! 🎉
