# SonarQube Go Client

[![Go Reference](https://pkg.go.dev/badge/github.com/boxboxjason/sonarqube-client-go.svg)](https://pkg.go.dev/github.com/boxboxjason/sonarqube-client-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/boxboxjason/sonarqube-client-go)](https://goreportcard.com/report/github.com/boxboxjason/sonarqube-client-go)
[![License](https://img.shields.io/github/license/BoxBoxJason/sonarqube-client-go)](LICENSE)

A comprehensive, type-safe Go client library for the SonarQube Web API. This SDK provides idiomatic Go access to SonarQube's extensive REST API, enabling seamless integration with SonarQube for code quality analysis, project management, and continuous inspection workflows.

## Table of Contents

- [Project Goals](#project-goals)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Available Services](#available-services)
- [Dependencies](#dependencies)
- [Contributing](#contributing)
- [Reporting Issues](#reporting-issues)
- [License](#license)

## Project Goals

The SonarQube Go Client aims to:

- **Provide complete API coverage**: Support all SonarQube Web API endpoints, including projects, issues, quality gates, quality profiles, users, permissions, and more
- **Type-safe interactions**: Offer strongly-typed Go structs for all API requests and responses to catch errors at compile time
- **Idiomatic Go design**: Follow Go best practices and conventions for a natural developer experience
- **Maintainability**: Keep the codebase clean, well-documented, and easy to extend
- **Production-ready**: Deliver a reliable SDK with comprehensive testing and error handling

This SDK is ideal for:

- Building automation tools for SonarQube
- Integrating code quality analysis into CI/CD pipelines
- Creating custom dashboards and reporting tools
- Managing SonarQube resources programmatically
- Extending SonarQube functionality with custom applications

## Features

### Core Capabilities

- ✅ **Complete API Coverage**: Support for all major SonarQube API services
- ✅ **Type Safety**: Strongly-typed request options and response structures
- ✅ **Flexible Authentication**: Support for token-based and username/password authentication
- ✅ **Multiple Response Formats**: Handle JSON, Protocol Buffers, text, and binary responses
- ✅ **Error Handling**: Detailed error responses with status codes and messages
- ✅ **Modern Go**: Built with Go 1.25+, using Go modules (no GOPATH required)
- ✅ **Well Tested**: Comprehensive unit tests and integration tests against real SonarQube instances
- ✅ **Production Ready**: Used in production environments with extensive edge case handling

### Response Type Support

The client intelligently handles different response formats:

- **JSON**: Automatically unmarshals JSON responses into Go structs
- **Protocol Buffers**: Returns raw bytes for protobuf endpoints
- **Text/CSV**: Returns plain text responses as strings
- **Binary**: Handles binary data (e.g., file downloads)

## Installation

Install the SDK using `go get`:

```bash
go get github.com/boxboxjason/sonarqube-client-go/sonar
```

**Requirements**:

- Go 1.25 or higher
- Access to a SonarQube instance (version 25+ recommended)

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/boxboxjason/sonarqube-client-go/sonar"
)

func main() {
    // Create a new client with token authentication
    client, err := sonar.NewClient(
        "https://your-sonarqube-instance.com",
        sonar.WithToken("your-sonarqube-token"),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Example 1: Search for projects
    projects, _, err := client.Projects.Search(context.Background(), &sonar.ProjectsSearchOption{
        Ps: sonar.Int(10), // Page size: 10 results
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d projects:\n", len(projects.Components))
    for _, project := range projects.Components {
        fmt.Printf("  - %s (%s)\n", project.Name, project.Key)
    }

    // Example 2: Get project issues
    issues, _, err := client.Issues.Search(context.Background(), &sonar.IssuesSearchOption{
        Projects: sonar.String("my-project-key"),
        Statuses: sonar.String("OPEN,CONFIRMED"),
        Ps:       sonar.Int(50),
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("\nFound %d open issues\n", issues.Total)
}
```

### Authentication Options

#### Token Authentication (Recommended)

```go
client, err := sonar.NewClient(
    "https://your-sonarqube-instance.com",
    sonar.WithToken("your-sonarqube-token"),
)
```

#### Username/Password Authentication

```go
client, err := sonar.NewClient(
    "https://your-sonarqube-instance.com",
    sonar.WithBasicAuth("username", "password"),
)
```

### Advanced Usage

#### Custom HTTP Client

```go
import "net/http"

httpClient := &http.Client{
    Timeout: 30 * time.Second,
}

client, err := sonar.NewClient(
    "https://your-sonarqube-instance.com",
    sonar.WithHTTPClient(httpClient),
    sonar.WithToken("your-token"),
)
```

#### Working with Quality Gates

```go
// Get quality gate status for a project
status, _, err := client.Qualitygates.GetProjectStatus(context.Background(), &sonar.QualitygatesGetProjectStatusOption{
    ProjectKey: sonar.String("my-project"),
})
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Quality Gate Status: %s\n", status.ProjectStatus.Status)
```

#### Managing Users

```go
// Search for users
users, _, err := client.Users.Search(context.Background(), &sonar.UsersSearchOption{
    Q: sonar.String("john"),
})
if err != nil {
    log.Fatal(err)
}

for _, user := range users.Users {
    fmt.Printf("User: %s (%s)\n", user.Name, user.Login)
}
```

## Available Services

The SDK provides access to the following SonarQube API services:

<details>
<summary><strong>Project Management</strong></summary>

- **Projects** - Create, search, and manage projects
- **Project Analyses** - Retrieve and delete project analysis history
- **Project Badges** - Generate quality badges for projects
- **Project Branches** - Manage project branches and pull requests
- **Project Dump** - Import and export project data
- **Project Links** - Manage project external links
- **Project Tags** - Add and manage project tags

</details>

<details>
<summary><strong>Code Analysis</strong></summary>

- **Issues** - Search, assign, and manage code issues
- **Hotspots** - Security hotspot management
- **Duplications** - Code duplication detection
- **Measures** - Retrieve project and component metrics
- **Metrics** - Manage custom metrics
- **Sources** - Access source code and blame information
- **Analysis Cache** - Manage analysis cache
- **Analysis Reports** - Access analysis reports

</details>

<details>
<summary><strong>Quality Management</strong></summary>

- **Quality Gates** - Define and manage quality gates
- **Quality Profiles** - Manage quality profiles and rules
- **Rules** - Search and manage coding rules
- **New Code Periods** - Configure new code period definitions

</details>

<details>
<summary><strong>User & Permission Management</strong></summary>

- **Users** - Create and manage users
- **User Groups** - Manage user groups
- **User Tokens** - Generate and revoke user tokens
- **Permissions** - Manage project and global permissions
- **Authentication** - Validate authentication and sessions

</details>

<details>
<summary><strong>Integration & ALM</strong></summary>

- **ALM Integrations** - Application Lifecycle Management integrations
- **ALM Settings** - Configure GitHub, GitLab, Bitbucket, and Azure DevOps
- **GitHub Provisioning** - GitHub organization and user provisioning
- **Webhooks** - Configure and manage webhooks

</details>

<details>
<summary><strong>System & Administration</strong></summary>

- **System** - System health, status, and information
- **Server** - Server version and status
- **Settings** - Global and project settings management
- **Monitoring** - System monitoring endpoints
- **Plugins** - Install and manage plugins
- **CE (Compute Engine)** - Background task management

</details>

<details>
<summary><strong>Developer Tools</strong></summary>

- **Developers** - Developer-specific metrics and issues
- **Favorites** - Manage favorite projects and filters
- **Notifications** - User notification preferences
- **Navigation** - UI navigation and component information
- **Web Services** - API documentation and metadata

</details>

<details>
<summary><strong>Other Services</strong></summary>

- **Batch** - Batch operations for IDE integrations
- **Components** - Component tree and search
- **Dismiss Message** - User message preferences
- **Emails** - Email configuration testing
- **Features** - Feature flag management
- **L10n** - Localization and internationalization
- **Languages** - Supported programming languages
- **Push** - Push events for live updates

</details>

**Total**: 50+ services covering all SonarQube API endpoints

## Dependencies

This project uses carefully selected dependencies to provide a robust and maintainable SDK:

### Production Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| [google/go-querystring](https://github.com/google/go-querystring) | v1.2.0 | Encoding Go structs into URL query parameters for API requests. Provides type-safe URL generation. |

### Development/Testing Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| [onsi/ginkgo](https://github.com/onsi/ginkgo) | v2.28.1 | BDD-style testing framework used for integration tests. Provides expressive test structure and excellent readability. |
| [onsi/gomega](https://github.com/onsi/gomega) | v1.39.1 | Matcher library for Ginkgo. Offers rich assertion syntax for testing. |
| [stretchr/testify](https://github.com/stretchr/testify) | v1.11.1 | Assertion toolkit for unit tests. Provides simple and intuitive test assertions. |
| [Masterminds/semver](https://github.com/Masterminds/semver) | v3.4.0 | Semantic versioning library. Used for version comparison and validation. |

### Standard Library Usage

The SDK extensively uses Go's standard library, including:

- `net/http` - HTTP client and server functionality
- `encoding/json` - JSON encoding/decoding
- `context` - Request context and cancellation
- `io` - Input/output primitives

### Why These Dependencies?

- **Minimal external dependencies**: Only one production dependency reduces supply chain risk and maintenance burden
- **Well-maintained packages**: All dependencies are actively maintained with strong community support
- **Type safety**: `go-querystring` ensures compile-time safety for API request parameters
- **Testing excellence**: Ginkgo/Gomega provide superior test readability and reporting
- **No unnecessary bloat**: Each dependency serves a specific, irreplaceable purpose

## Contributing

We welcome contributions from the community! Whether you're fixing bugs, adding features, or improving documentation, your help is appreciated.

### How to Contribute

1. **Read the [Contributing Guide](CONTRIBUTING.md)** - Comprehensive guide covering development setup, testing, and PR guidelines
2. **Check existing issues** - Look for issues labeled `good first issue` or `help wanted`
3. **Fork and create a branch** - Use conventional commit format (`feat:`, `fix:`, `docs:`, etc.)
4. **Write tests** - Ensure 80%+ coverage for new code
5. **Submit a Pull Request** - Reference the issue you're addressing

### Quick Contribution Checklist

Before submitting a PR, ensure:

- [ ] `make lint` passes without errors
- [ ] `make test` passes (unit tests)
- [ ] `make e2e` passes (integration tests, if applicable)
- [ ] Code has appropriate godoc comments
- [ ] Tests cover new functionality (80%+ coverage)
- [ ] Commit messages follow [conventional commits](https://www.conventionalcommits.org/)
- [ ] PR description references an issue (e.g., "Fixes #123")

For detailed guidelines, see [CONTRIBUTING.md](CONTRIBUTING.md).

## Reporting Issues

### Bug Reports

Found a bug? Help us fix it by providing:

- Clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Go version, OS, and SonarQube version
- Code sample demonstrating the issue (if possible)

[Report a bug →](https://github.com/BoxBoxJason/sonarqube-client-go/issues/new)

### Feature Requests

Have an idea for a new feature or improvement?

- Describe your use case and how it would benefit users
- Explain why existing alternatives don't meet your needs
- Provide examples of how you'd like the API to work

[Request a feature →](https://github.com/BoxBoxJason/sonarqube-client-go/issues/new)

### Security Issues

If you discover a security vulnerability:

- **DO NOT** open a public issue
- Email security concerns privately to the maintainers
- Provide detailed information to help us address the issue quickly

For non-security bugs and features, please use GitHub Issues.

## License

This project is licensed under the **Apache License 2.0** - see the [LICENSE](LICENSE) file for details.

### What This Means

- ✅ **Commercial use** - Use this SDK in commercial projects
- ✅ **Modification** - Modify the source code
- ✅ **Distribution** - Distribute the SDK
- ✅ **Patent use** - Use any patents that cover the SDK
- ⚠️ **License and copyright notice** - Include the license and copyright notice in distributions
- ⚠️ **State changes** - Document significant changes you make

---

## Additional Resources

- [SonarQube Web API Documentation](https://docs.sonarsource.com/sonarqube/latest/extension-guide/web-api/)
- [Go Documentation](https://pkg.go.dev/github.com/boxboxjason/sonarqube-client-go)
- [Contributing Guide](CONTRIBUTING.md)
- [GitHub Issues](https://github.com/BoxBoxJason/sonarqube-client-go/issues)

---

**Made with ❤️ by [BoxBoxJason](https://github.com/BoxBoxJason)**

If this SDK helps you, please consider giving it a ⭐ on GitHub!
