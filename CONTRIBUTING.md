# Contributing to SonarQube Go Client

Thank you for your interest in contributing to the SonarQube Go Client! This guide will help you get started with development, testing, and submitting contributions.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Makefile Commands](#makefile-commands)
- [Pull Request Guidelines](#pull-request-guidelines)
- [Reporting Issues](#reporting-issues)

## Code of Conduct

### Common Courtesy Rules

We strive to maintain a welcoming and inclusive community. When contributing, please:

- **Be respectful**: Treat all contributors with respect and professionalism
- **Be constructive**: Provide helpful feedback and suggestions
- **Be patient**: Remember that everyone has different experience levels
- **Be collaborative**: Work together to improve the project
- **Stay on topic**: Keep discussions focused on the project
- **Follow conventions**: Respect the established coding style and patterns

## Getting Started

### Prerequisites

- **Go 1.25 or higher** - [Installation guide](https://go.dev/doc/install)
- **Docker or Podman** (for integration tests) - The Makefile automatically detects which is available
- **Git** - For version control
- **curl** - For API specification fetching

### Setting Up Your Development Environment

1. **Fork and clone the repository**:

   ```bash
   git clone https://github.com/boxboxjason/sonarqube-client-go.git
   cd sonarqube-client-go
   ```

2. **Install Go dependencies**:

   ```bash
   go mod download
   ```

3. **Install development tools** (optional, they will be installed automatically when needed):

   ```bash
   # Test runner with enhanced output
   go install gotest.tools/gotestsum@v1.13.0

   # Linter
   go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.8.0

   # Integration test framework
   go install github.com/onsi/ginkgo/v2/ginkgo@v2.28.1
   ```

4. **Verify your setup**:

   ```bash
   make lint
   make test
   ```

## Project Structure

### Directory Layout

```plaintext
sonarqube-client-go/
├── sonar/                      # SDK source code and unit tests
│   ├── client.go              # Main client implementation
│   ├── *_service.go           # Service implementations (e.g., projects_service.go)
│   ├── *_service_test.go      # Unit tests for services
│   ├── common.go              # Shared types and utilities
│   ├── errors.go              # Error handling
│   └── ...
├── integration_testing/        # End-to-end integration tests
│   ├── *_test.go              # Integration test files
│   ├── suite_test.go          # Test suite setup
│   └── helpers/               # Test utilities
│       ├── client.go          # Shared test client
│       ├── cleanup.go         # Resource cleanup helpers
│       └── wait.go            # Wait and retry utilities
├── assets/                     # API specifications and resources
├── codequality/               # Generated test reports and coverage
├── .github/                   # GitHub workflows and issue templates
├── Makefile                   # Build and development tasks
├── go.mod                     # Go module dependencies
└── README.md                  # Project documentation
```

### Key Directories

#### `sonar/` - SDK Implementation

This is the main package containing:

- **Service files**: Each `*_service.go` file implements a SonarQube API service (e.g., `projects_service.go`, `issues_service.go`)
- **Unit test files**: Corresponding `*_service_test.go` files contain unit tests with mocked API responses
- **Core client**: `client.go` provides the main client structure and HTTP communication
- **Common types**: Shared structs, constants, and utility functions

When adding or modifying SDK functionality:

- Place your code in the appropriate `*_service.go` file
- Add corresponding unit tests in `*_service_test.go`
- Use the `fake` subdirectories for test helpers and mock data

#### `integration_testing/` - End-to-End Tests

This directory contains integration tests that run against a real SonarQube instance:

- Each `*_test.go` file corresponds to a service being tested
- Tests use the Ginkgo testing framework
- Helper utilities are located in `integration_testing/helpers/`
- Tests validate actual API communication and responses

## Development Workflow

### Making Changes

1. **Create a branch** from `main`:

   ```bash
   git checkout -b feat/your-feature-name
   # or
   git checkout -b fix/your-bug-fix
   ```

2. **Follow conventional commit format** for your commit messages:
   - `feat:` - New feature
   - `fix:` - Bug fix
   - `docs:` - Documentation changes
   - `test:` - Test additions or modifications
   - `refactor:` - Code refactoring
   - `chore:` - Maintenance tasks
   - `perf:` - Performance improvements
   - `ci:` - CI configuration changes
   - `build:` - Build system changes

   Example:

   ```bash
   git commit -m "feat: add support for project analysis history"
   ```

3. **Write clean, idiomatic Go code**:
   - Follow the existing project style and naming conventions
   - Add godoc comments for all exported functions, types, and methods
   - Keep functions focused and reasonably sized
   - Use descriptive variable names

4. **Document complex logic**:
   - Add comments explaining the "why" for non-obvious code
   - Document any assumptions or edge cases
   - Reference related issues or SonarQube API documentation when relevant

### Code Style

All code must pass the project's linter. The configuration is in [.golangci.yml](.golangci.yml) and includes:

- **Formatters**: `gci`, `gofmt`, `gofumpt`, `goimports`
- **Linters**: `govet`, and many others (see `.golangci.yml` for the full list)
- **Import ordering**: Standard library → External packages → This project

Run the linter before submitting:

```bash
make lint
```

## Testing

### Unit Tests

Unit tests are located alongside the source code in the `sonar/` directory.

**Requirements**:

- Achieve **80-100% coverage** for new or modified code
- Use **table-driven tests** for functions with multiple scenarios
- **Mock API responses** - Do not make real HTTP calls in unit tests
- Use common test utilities from `test_helpers_test.go` or `fake/` subdirectories

**Writing unit tests**:

```go
func TestMyService_MyMethod(t *testing.T) {
    tests := []struct {
        name        string
        input       MyInput
        mockResp    string
        mockStatus  int
        expected    *MyOutput
        expectError bool
    }{
        {
            name:       "successful request",
            input:      MyInput{Field: "value"},
            mockResp:   `{"result": "success"}`,
            mockStatus: 200,
            expected:   &MyOutput{Result: "success"},
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

**Run unit tests**:

```bash
make test           # Run all unit tests
make coverage       # Run tests with coverage report
```

### Integration Tests (E2E)

Integration tests validate the SDK against a real SonarQube instance.

**Requirements**:

- Located in `integration_testing/` directory
- Use Ginkgo/Gomega testing framework
- Clean up resources after tests (use helpers from `integration_testing/helpers/cleanup.go`)
- Handle eventual consistency with retry logic (use helpers from `integration_testing/helpers/wait.go`)

**Run integration tests**:

```bash
make e2e            # Starts SonarQube and runs integration tests
```

The `make e2e` command automatically:

1. Checks if SonarQube is already running
2. Starts a SonarQube container if needed (Docker/Podman)
3. Waits for SonarQube to be ready
4. Runs all integration tests

To teardown the SonarQube instance:

```bash
make teardown.sonar
```

## Makefile Commands

The project includes several Make targets for common development tasks:

### Testing Commands

| Command | Description |
|---------|-------------|
| `make test` | Run all unit tests with pretty output and generate JUnit XML report |
| `make coverage` | Run unit tests with coverage report (generates `codequality/coverage.out`) |
| `make e2e` | Set up SonarQube and run end-to-end integration tests |

### Code Quality Commands

| Command | Description |
|---------|-------------|
| `make lint` | Run `golangci-lint` and generate a checkstyle report for CI |

### SonarQube Management Commands

| Command | Description |
|---------|-------------|
| `make setup.sonar` | Start a SonarQube instance for testing (Docker/Podman) |
| `make teardown.sonar` | Stop and remove the SonarQube container |
| `make api` | Fetch the SonarQube API specification from a running instance |

### Documentation Commands

| Command | Description |
|---------|-------------|
| `make changelog` | Generate CHANGELOG.md using git-cliff |
| `make changelog-check` | Verify CHANGELOG.md is up-to-date |

### Configuration Variables

You can customize the Makefile behavior with environment variables:

```bash
# Example: Use a different SonarQube instance for integration tests
make e2e endpoint=http://my-sonarqube:9000 username=admin password=mypass
```

Available variables:

- `endpoint` - SonarQube API endpoint (default: `http://127.0.0.1:9000`)
- `username` - SonarQube username (default: `admin`)
- `password` - SonarQube password (default: `admin`)
- `sonarqube_version` - Docker image version (default: `26.1.0.118079-community`)

## Pull Request Guidelines

### Before Submitting

Ensure your pull request meets these requirements:

- [ ] **Linting passes**: `make lint` completes without errors
- [ ] **Unit tests pass**: `make test` completes successfully
- [ ] **Integration tests pass**: `make e2e` completes successfully (when applicable)
- [ ] **Code coverage**: New code has at least 80% test coverage
- [ ] **Conventional commits**: All commits follow the conventional commit format
- [ ] **Issue reference**: PR description references the issue it addresses (e.g., "Fixes #123")
- [ ] **Documentation**: Code includes godoc comments for exported symbols
- [ ] **Tests included**: New features and bug fixes include appropriate tests

### PR Description Template

When you create a PR, GitHub will automatically populate the description with our template. Be sure to:

1. Describe what your PR does
2. Reference the issue it fixes: `Fixes #<issue-number>`
3. Check off all completed items in the checklist
4. Describe how you tested your changes

### Review Process

1. **Automated checks**: CI will run linting, unit tests, and integration tests
2. **Code review**: Maintainers will review your code for quality and correctness
3. **Feedback**: Address any comments or requested changes
4. **Approval**: Once approved, a maintainer will merge your PR

### Common Review Feedback

- **Missing tests**: Add unit tests covering your changes
- **Documentation gaps**: Add godoc comments for exported functions
- **Linting errors**: Run `make lint` and fix reported issues
- **Test failures**: Ensure `make test` and `make e2e` pass locally
- **Code complexity**: Consider breaking large functions into smaller helpers
- **Naming conventions**: Follow Go naming conventions and match the existing style

## Reporting Issues

### Found a Bug?

If you discover a bug:

1. **Check existing issues**: Search [GitHub Issues](https://github.com/BoxBoxJason/sonarqube-client-go/issues) to see if it's already reported
2. **Create a new issue**: If not found, [open a new issue](https://github.com/BoxBoxJason/sonarqube-client-go/issues/new)
3. **Provide details**:
   - Clear description of the bug
   - Steps to reproduce
   - Expected vs actual behavior
   - Go version and OS (if relevant)
   - SonarQube version (if relevant)
   - Code sample demonstrating the issue (if possible)

### Feature Requests

Have an idea for a new feature?

1. **Check existing issues**: See if someone else has requested it
2. **Open a feature request**: [Create a new issue](https://github.com/BoxBoxJason/sonarqube-client-go/issues/new) with the `enhancement` label
3. **Describe your use case**:
   - What problem would this solve?
   - How would you like it to work?
   - Are there alternatives you've considered?

### Questions and Discussions

For general questions or discussions:

- Open a [GitHub Discussion](https://github.com/BoxBoxJason/sonarqube-client-go/discussions) (if enabled)
- Create an issue with the `question` label
- Check the project documentation first

## Additional Resources

- [SonarQube Web API Documentation](https://docs.sonarsource.com/sonarqube/latest/extension-guide/web-api/)
- [Go Documentation](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Conventional Commits](https://www.conventionalcommits.org/)

---

Thank you for contributing to the SonarQube Go Client! Your efforts help make this SDK better for everyone. 🎉
