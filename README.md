# SonarQube Go Client

[![Go Reference](https://pkg.go.dev/badge/github.com/boxboxjason/sonarqube-client-go/v2.svg)](https://pkg.go.dev/github.com/boxboxjason/sonarqube-client-go/v2)
[![codecov](https://codecov.io/gh/BoxBoxJason/sonarqube-client-go/graph/badge.svg)](https://codecov.io/gh/BoxBoxJason/sonarqube-client-go)
[![License](https://img.shields.io/github/license/BoxBoxJason/sonarqube-client-go)](LICENSE)

A comprehensive, type-safe Go client library **and command-line interface** for the SonarQube Web API. Whether you're building automation tools, integrating SonarQube into your CI/CD pipeline, or managing your SonarQube instance from the terminal - this project has you covered.

📖 **Documentation: <https://boxboxjason.github.io/sonarqube-client-go/>** · 📚 **API reference: [pkg.go.dev](https://pkg.go.dev/github.com/boxboxjason/sonarqube-client-go/v2/sonar)**

## Features

- **Go SDK**: strongly typed options and responses for every SonarQube V1 and V2 API service, token or basic authentication, opt-in retries with backoff, HTTP middleware for logging and tracing.
- **CLI (`sonar-cli`)**: every service and method as a subcommand, JSON / YAML / table output, automatic pagination with `--all`, shell completion.
- **Tested** with unit tests and integration tests against real SonarQube Community and Enterprise instances.

## Quick start

### CLI

```bash
go install github.com/boxboxjason/sonarqube-client-go/v2/cmd/sonar-cli@latest

export SONAR_CLI_URL=https://sonar.example.com
export SONAR_CLI_TOKEN=mytoken
sonar-cli issues search --projects my-project --severities CRITICAL,MAJOR --all
```

### Go SDK

```bash
go get github.com/boxboxjason/sonarqube-client-go/v2/sonar
```

```go
client, err := sonar.NewClient(nil,
	sonar.WithBaseURL("https://sonarqube.example.com/api/"),
	sonar.WithToken("my-token"),
)
if err != nil {
	log.Fatal(err)
}

projects, _, err := client.Projects.Search(context.Background(), &sonar.ProjectsSearchOptions{
	Query: "my-service",
})
```

The [documentation](https://boxboxjason.github.io/sonarqube-client-go/) covers the [CLI](https://boxboxjason.github.io/sonarqube-client-go/docs/cli/), the [SDK](https://boxboxjason.github.io/sonarqube-client-go/docs/sdk/), [client configuration](https://boxboxjason.github.io/sonarqube-client-go/docs/configuration/), [errors and pagination](https://boxboxjason.github.io/sonarqube-client-go/docs/errors-and-pagination/) and the list of [services](https://boxboxjason.github.io/sonarqube-client-go/docs/services/). It is built from the [`docs`](./docs) directory, and updated together with the code.

## Contributing

Contributions are welcome: read the [Contributing Guide](CONTRIBUTING.md), which covers the development setup, the tests and the documentation.

## Reporting Issues

- [Report a bug](https://github.com/BoxBoxJason/sonarqube-client-go/issues/new) with the steps to reproduce, the Go, OS and SonarQube versions, and a minimal reproducer.
- [Request a feature](https://github.com/BoxBoxJason/sonarqube-client-go/issues/new).
- **Do not** open a public issue for a security vulnerability: email the maintainers privately with the details.

## License

Licensed under the **Apache License 2.0**, see [LICENSE](LICENSE).
