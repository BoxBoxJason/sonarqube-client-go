---
title: CLI
description: "Use sonar-cli to call every SonarQube web API service from the terminal, with JSON, YAML or table output, automatic pagination and shell completion."
weight: 20
---

`sonar-cli` exposes every service and method of the SDK as a subcommand:

```text
sonar-cli [global flags] <service> <method> [flags]
```

The list of services, and of the methods of each one, is on the [services]({{< relref "services" >}}) page. `--help` works at every level:

```bash
sonar-cli --help
sonar-cli projects --help
sonar-cli issues search --help
```

## Authentication

Authentication is passed through global flags, or through environment variables (the recommended way in scripts and CI). The `--url` is mandatory.

| Flag | Environment variable | Description |
| ---- | -------------------- | ----------- |
| `--url` | `SONAR_CLI_URL` | SonarQube server URL |
| `--token` | `SONAR_CLI_TOKEN` | Authentication token (recommended) |
| `--username` | `SONAR_CLI_USERNAME` | Username for basic authentication |
| `--password` | `SONAR_CLI_PASSWORD` | Password for basic authentication |

```bash
# Using flags
sonar-cli --url https://sonar.example.com --token mytoken projects search

# Using environment variables
export SONAR_CLI_URL=https://sonar.example.com
export SONAR_CLI_TOKEN=mytoken
sonar-cli projects search
```

{{< callout type="warning" >}}
Values passed on the command line are visible to the other users of the machine (in `ps`, in the shell history): prefer the environment variables for the token and the password.
{{< /callout >}}

## Other global flags

| Flag | Default | Description |
| ---- | ------- | ----------- |
| `--output` | `json` | Output format: `json`, `yaml` or `table` |
| `--timeout` | `30s` | HTTP request timeout (Go duration) |

## Examples

```bash
# List all projects
sonar-cli projects search

# Search for critical issues in a project
sonar-cli issues search --projects my-project --severities CRITICAL,MAJOR

# Get the quality gate status of a project
sonar-cli qualitygates project-status --project-key my-project

# Create a user token
sonar-cli user-tokens generate --login john --name "ci-token"

# Search rules
sonar-cli rules search --languages go --severities MAJOR

# Get the system health
sonar-cli system health
```

## Output formats

```bash
# JSON (default), for jq and scripting
sonar-cli projects search --output json | jq '.components[].key'

# YAML
sonar-cli projects search --output yaml

# Table, for a quick look in the terminal
sonar-cli projects search --output table
```

```text
KEY              NAME             QUALIFIER  VISIBILITY
my-project       My Project       TRK        public
another-project  Another Project  TRK        private
```

## Pagination

Methods returning paginated results accept `--all`, which fetches every page and merges the results:

```bash
# First page only (default)
sonar-cli issues search --projects my-project

# Every page
sonar-cli issues search --projects my-project --all

# Manual pagination
sonar-cli projects search --page 2 --page-size 50
```

Method flags are named after the fields of the SDK options (`PageSize` becomes `--page-size`), so the [API reference](https://pkg.go.dev/github.com/boxboxjason/sonarqube-client-go/v2/sonar) also documents the flags of every method.

## Shell completion

Services, methods and flags can be completed with `Tab`:

```bash
# Bash
sonar-cli completion bash > ~/.local/share/bash-completion/completions/sonar-cli

# Zsh
sonar-cli completion zsh > "${fpath[1]}/_sonar-cli"

# Fish
sonar-cli completion fish > ~/.config/fish/completions/sonar-cli.fish
```

PowerShell:

```powershell
sonar-cli completion powershell | Out-String | Invoke-Expression
```
