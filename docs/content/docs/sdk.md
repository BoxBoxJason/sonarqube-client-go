---
title: Go SDK
description: "Create a SonarQube client in Go, authenticate, and call V1 and V2 web API services with strongly typed request options and responses."
weight: 30
---

## Create a client

The root type is [`Client`](https://pkg.go.dev/github.com/boxboxjason/sonarqube-client-go/v2/sonar#Client). Create one with `NewClient`, configured with functional options:

{{% example "Example" %}}

or with a `ClientCreateOptions` struct (both can be combined, the options are applied after the struct):

{{% example "ExampleNewClient" %}}

{{< callout type="info" >}}
The base URL is the URL of the API, ending with `/api/`: `https://sonarqube.example.com/api/`.
{{< /callout >}}

## Authentication

A [user token](https://docs.sonarsource.com/sonarqube-server/latest/user-guide/managing-tokens/) is recommended (`WithToken`, above). Basic authentication is also supported:

{{% example "ExampleWithBasicAuth" %}}

### From the environment

In CI/CD, `NewClientFromEnv` reads the configuration from the environment:

| Environment variable | Description |
| -------------------- | ----------- |
| `SONAR_URL` | Base API URL of the SonarQube instance |
| `SONAR_TOKEN` | Authentication token, takes precedence over basic authentication |
| `SONAR_USERNAME` | Username for basic authentication |
| `SONAR_PASSWORD` | Password for basic authentication |
| `SONAR_USER_AGENT` | Overrides the default `User-Agent` header |
| `SONAR_TIMEOUT` | HTTP request timeout, as a Go duration (`30s`, `2m`) |

{{% example "ExampleNewClientFromEnv" %}}

Options passed to `NewClientFromEnv` are applied after the environment is read, so they take precedence over it.

## Call the API

Each API area is a typed service field of the client: `client.Projects`, `client.Issues`, `client.Qualitygates`, ... (see the [services]({{< relref "services" >}})). Every method takes a `context.Context` and, when the endpoint has parameters, a `*<Service><Method>Options` struct. It returns the decoded response, the raw `*http.Response` and an error.

{{% example "ExampleQualitygatesService_ProjectStatus" %}}

{{% example "ExampleIssuesService_Search" %}}

## V1 and V2 APIs

SonarQube exposes two web APIs. The V1 services are fields of the client (`client.Issues`), and encode their parameters in the query string. The V2 services are grouped under `client.V2` (`client.V2.System`), use JSON request bodies, and get the `v2/` path prefix automatically:

{{% example "ExampleSystemServiceV2_GetHealth" %}}

The [services]({{< relref "services#v2-api-services" >}}) page lists both.
