---
title: Documentation
description: "Documentation for sonarqube-client-go, the Go SDK and sonar-cli command-line tool for automating SonarQube through its web API."
weight: 1
---

sonarqube-client-go is a type-safe Go client library **and command-line interface** for the [SonarQube web API](https://docs.sonarsource.com/sonarqube-server/latest/extension-guide/web-api/).

- **Go SDK** (`sonar` package): strongly typed request options and responses, token or basic authentication, V1 and V2 APIs, opt-in retries, pluggable HTTP middleware.
- **CLI** (`sonar-cli`): every SDK service and method as a subcommand, JSON / YAML / table output, automatic pagination, shell completion.

These guides explain how to use both. The details of every type and method are in the [API reference on pkg.go.dev](https://pkg.go.dev/github.com/boxboxjason/sonarqube-client-go/v2/sonar).

{{< cards >}}
  {{< card link="installation" title="Installation" subtitle="Install the CLI or add the SDK to a Go module." >}}
  {{< card link="cli" title="CLI" subtitle="Authentication, output formats, pagination and completion of sonar-cli." >}}
  {{< card link="sdk" title="Go SDK" subtitle="Create a client, authenticate and call the API." >}}
  {{< card link="configuration" title="Client configuration" subtitle="Timeouts, HTTP client, transport, retries and middleware." >}}
  {{< card link="errors-and-pagination" title="Errors and pagination" subtitle="Handle API errors and walk through paginated results." >}}
  {{< card link="services" title="Services" subtitle="Every service, its CLI command and its methods (generated)." >}}
  {{< card link="compatibility" title="Compatibility" subtitle="Supported SonarQube versions and editions." >}}
{{< /cards >}}
