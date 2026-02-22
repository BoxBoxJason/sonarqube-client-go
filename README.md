# SonarQube Go Client

[![Go Reference](https://pkg.go.dev/badge/github.com/boxboxjason/sonarqube-client-go.svg)](https://pkg.go.dev/github.com/boxboxjason/sonarqube-client-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/boxboxjason/sonarqube-client-go)](https://goreportcard.com/report/github.com/boxboxjason/sonarqube-client-go)
[![License](https://img.shields.io/github/license/BoxBoxJason/sonarqube-client-go)](LICENSE)

A comprehensive, type-safe Go client library **and command-line interface** for the SonarQube Web API. Whether you're building automation tools, integrating SonarQube into your CI/CD pipeline, or managing your SonarQube instance from the terminal — this project has you covered.

## Table of Contents

- [Features](#features)
- [CLI — sonar-cli](#cli--sonar-cli)
  - [Installation](#installation)
  - [Authentication](#authentication)
  - [Basic Usage](#basic-usage)
  - [Output Formats](#output-formats)
  - [Pagination](#pagination)
  - [Shell Completion](#shell-completion)
- [Go SDK](#go-sdk)
  - [SDK Installation](#sdk-installation)
  - [Quick Start](#quick-start)
  - [Authentication Options](#authentication-options)
  - [Advanced Usage](#advanced-usage)
- [Available Services](#available-services)
- [Dependencies](#dependencies)
- [Contributing](#contributing)
- [Reporting Issues](#reporting-issues)
- [License](#license)

## Features

### CLI (`sonar-cli`)

- ✅ **Full API Coverage from the Terminal**: Every SonarQube service and method available as a subcommand
- ✅ **Multiple Output Formats**: JSON, YAML, and ASCII table — pipe-friendly
- ✅ **Automatic Pagination**: Fetch all pages of results with a single `--all` flag
- ✅ **Shell Completion**: Tab completion for Bash, Zsh, Fish, and PowerShell
- ✅ **Flexible Authentication**: Token, username/password, or environment variables
- ✅ **Structured Error Logging**: Clear, context-rich error messages to stderr
- ✅ **Configurable Timeout**: Per-command HTTP timeout control

### Go SDK

- ✅ **Complete API Coverage**: Support for all major SonarQube API services
- ✅ **Type Safety**: Strongly-typed request options and response structures
- ✅ **Flexible Authentication**: Token-based and username/password authentication
- ✅ **Multiple Response Formats**: JSON, Protocol Buffers, text, and binary responses
- ✅ **Modern Go**: Built with Go 1.25+, using Go modules
- ✅ **Well Tested**: Comprehensive unit tests and integration tests against real SonarQube instances

---

## CLI — sonar-cli

`sonar-cli` is a fully featured command-line interface that wraps the entire SonarQube API. Every service and method available in the SDK is exposed as a CLI subcommand, making it trivial to interact with SonarQube from scripts, CI pipelines, or your terminal.

### Installation

**From source:**

```bash
go install github.com/boxboxjason/sonarqube-client-go/cmd/sonar-cli@latest
```

**From a release binary:**

Download the latest binary for your platform from the [Releases](https://github.com/BoxBoxJason/sonarqube-client-go/releases) page.

**Build locally:**

```bash
git clone https://github.com/BoxBoxJason/sonarqube-client-go.git
cd sonarqube-client-go
make build
# Binary: ./bin/sonar-cli
```

Build with a specific version:

```bash
make build version=1.2.3
./bin/sonar-cli --version  # sonar-cli version 1.2.3
```

### Authentication

Authentication can be provided via flags or environment variables (environment variables take precedence):

| Flag | Environment Variable | Description |
|------|---------------------|-------------|
| `--url` | `SONAR_CLI_URL` | SonarQube server URL |
| `--token` | `SONAR_CLI_TOKEN` | Authentication token (recommended) |
| `--username` | `SONAR_CLI_USERNAME` | Username for basic auth |
| `--password` | `SONAR_CLI_PASSWORD` | Password for basic auth |

```bash
# Using flags
sonar-cli --url https://sonar.example.com --token mytoken projects search

# Using environment variables (recommended for scripts and CI)
export SONAR_CLI_URL=https://sonar.example.com
export SONAR_CLI_TOKEN=mytoken
sonar-cli projects search
```

### Basic Usage

Commands follow the pattern: `sonar-cli [global flags] <service> <method> [flags]`

```bash
# List all projects
sonar-cli projects search

# Search for critical issues in a project
sonar-cli issues search --projects my-project --severities CRITICAL,MAJOR

# Get quality gate status
sonar-cli qualitygates get-project-status --project-key my-project

# Create a user token
sonar-cli user-tokens generate --login john --name "ci-token"

# Delete a project
sonar-cli projects delete --project my-old-project

# Search rules
sonar-cli rules search --languages go --severities MAJOR

# Get system health
sonar-cli system health
```

Use `--help` at any level to explore available commands:

```bash
sonar-cli --help
sonar-cli projects --help
sonar-cli issues search --help
```

### Output Formats

Control output format globally with `--output` (default: `json`):

```bash
# JSON output (default) — great for jq and scripting
sonar-cli projects search --output json | jq '.components[].key'

# YAML output — human-readable structured data
sonar-cli projects search --output yaml

# Table output — quick visual overview in the terminal
sonar-cli projects search --output table
```

**Table output example:**

```
KEY              NAME             QUALIFIER  VISIBILITY
my-project       My Project       TRK        public
another-project  Another Project  TRK        private
```

### Pagination

Endpoints that return paginated results support a `--all` flag to automatically fetch and merge every page:

```bash
# Fetch the first page (default)
sonar-cli issues search --projects my-project

# Fetch ALL issues across all pages automatically
sonar-cli issues search --projects my-project --all

# Manual pagination control
sonar-cli projects search --p 2 --ps 50
```

### Shell Completion

Enable tab-completion for your shell. Once set up, pressing `Tab` will autocomplete services, methods, and flags.

**Bash:**

```bash
sonar-cli completion bash > /etc/bash_completion.d/sonar-cli
# or for the current user:
sonar-cli completion bash > ~/.local/share/bash-completion/completions/sonar-cli
```

**Zsh:**

```bash
sonar-cli completion zsh > "${fpath[1]}/_sonar-cli"
# Then reload your shell or run:
source ~/.zshrc
```

**Fish:**

```bash
sonar-cli completion fish > ~/.config/fish/completions/sonar-cli.fish
```

**PowerShell:**

```powershell
sonar-cli completion powershell | Out-String | Invoke-Expression
```

---

## Go SDK

### SDK Installation

```bash
go get github.com/boxboxjason/sonarqube-client-go/sonar
```

**Requirements**: Go 1.25 or higher, access to a SonarQube instance (version 25+ recommended).

### Quick Start

```go
package main

import (
 "context"
 "fmt"
 "log"

 "github.com/boxboxjason/sonarqube-client-go/sonar"
)

func main() {
 url := "https://your-sonarqube-instance.com"
 token := "your-sonarqube-token"

 client, err := sonar.NewClient(&sonar.ClientCreateOption{
  URL:   &url,
  Token: &token,
 })
 if err != nil {
  log.Fatal(err)
 }

 // Search for projects
 projects, _, err := client.Projects.Search(context.Background(), &sonar.ProjectsSearchOption{
  Ps: sonar.Int(10),
 })
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("Found %d projects:\n", len(projects.Components))
 for _, project := range projects.Components {
  fmt.Printf("  - %s (%s)\n", project.Name, project.Key)
 }

 // Search for open issues
 issues, _, err := client.Issues.Search(context.Background(), &sonar.IssuesSearchOption{
  Projects: sonar.String("my-project-key"),
  Statuses: sonar.String("OPEN,CONFIRMED"),
 })
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("Found %d open issues\n", issues.Total)
}
```

### Authentication Options

**Token authentication (recommended):**

```go
client, err := sonar.NewClient(&sonar.ClientCreateOption{
 URL:   &url,
 Token: &token,
})
```

**Username/password authentication:**

```go
client, err := sonar.NewClient(&sonar.ClientCreateOption{
 URL:      &url,
 Username: &username,
 Password: &password,
})
```

### Advanced Usage

**Custom HTTP client:**

```go
import "net/http"
import "time"

httpClient := &http.Client{Timeout: 60 * time.Second}

client, err := sonar.NewClient(&sonar.ClientCreateOption{
 URL:        &url,
 Token:      &token,
 HttpClient: httpClient,
})
```

**Quality gate status:**

```go
status, _, err := client.Qualitygates.GetProjectStatus(ctx, &sonar.QualitygatesGetProjectStatusOption{
 ProjectKey: sonar.String("my-project"),
})
fmt.Printf("Quality Gate: %s\n", status.ProjectStatus.Status)
```

**User management:**

```go
users, _, err := client.Users.Search(ctx, &sonar.UsersSearchOption{
 Q: sonar.String("john"),
})
for _, user := range users.Users {
 fmt.Printf("User: %s (%s)\n", user.Name, user.Login)
}
```

---

## Available Services

Both the SDK and CLI expose all 50+ SonarQube API services:

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

## Dependencies

### Production Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| [google/go-querystring](https://github.com/google/go-querystring) | v1.2.0 | Type-safe URL query parameter encoding |
| [spf13/cobra](https://github.com/spf13/cobra) | v1.10.2 | CLI framework (commands, flags, shell completion) |
| [spf13/pflag](https://github.com/spf13/pflag) | v1.0.10 | POSIX-compatible flag parsing |
| [go.uber.org/zap](https://github.com/uber-go/zap) | v1.27.1 | Structured, high-performance logging |
| [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) | v3 | YAML output formatting |

### Development/Testing Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| [onsi/ginkgo](https://github.com/onsi/ginkgo) | v2.28.1 | BDD-style integration testing framework |
| [onsi/gomega](https://github.com/onsi/gomega) | v1.39.1 | Matcher library for Ginkgo |
| [stretchr/testify](https://github.com/stretchr/testify) | v1.11.1 | Unit test assertions and mocking |
| [Masterminds/semver](https://github.com/Masterminds/semver) | v3.4.0 | Semantic version comparison |

## Contributing

We welcome contributions! Whether you're fixing bugs, adding features, or improving documentation — your help is appreciated.

### How to Contribute

1. **Read the [Contributing Guide](CONTRIBUTING.md)** — covers development setup, testing, and PR guidelines
2. **Check existing issues** — look for labels `good first issue` or `help wanted`
3. **Fork and create a branch** — use conventional commit format (`feat:`, `fix:`, `docs:`, etc.)
4. **Write tests** — ensure 80%+ coverage for new code
5. **Submit a Pull Request** — reference the issue you're addressing

### Quick Contribution Checklist

- [ ] `make lint` passes without errors
- [ ] `make test` passes (unit tests)
- [ ] `make e2e` passes (integration tests, if applicable)
- [ ] Code has appropriate godoc comments
- [ ] Tests cover new functionality (80%+ coverage)
- [ ] Commit messages follow [conventional commits](https://www.conventionalcommits.org/)
- [ ] PR description references an issue (e.g., "Fixes #123")

## Reporting Issues

### Bug Reports

Found a bug? Provide:

- Clear description and steps to reproduce
- Expected vs actual behavior
- Go version, OS, and SonarQube version
- Minimal code or command reproducing the issue

[Report a bug →](https://github.com/BoxBoxJason/sonarqube-client-go/issues/new)

### Feature Requests

[Request a feature →](https://github.com/BoxBoxJason/sonarqube-client-go/issues/new)

### Security Issues

**DO NOT** open a public issue for security vulnerabilities. Email the maintainers privately with full details.

## License

This project is licensed under the **Apache License 2.0** — see the [LICENSE](LICENSE) file for details.

- ✅ Commercial use, modification, and distribution
- ✅ Patent use
- ⚠️ Include license and copyright notice in distributions

---

## Additional Resources

- [SonarQube Web API Documentation](https://docs.sonarsource.com/sonarqube/latest/extension-guide/web-api/)
- [Go Package Documentation](https://pkg.go.dev/github.com/boxboxjason/sonarqube-client-go)
- [Contributing Guide](CONTRIBUTING.md)
- [GitHub Issues](https://github.com/BoxBoxJason/sonarqube-client-go/issues)

---

**Made with ❤️ by [BoxBoxJason](https://github.com/BoxBoxJason)**

If this project helps you, please consider giving it a ⭐ on GitHub!
