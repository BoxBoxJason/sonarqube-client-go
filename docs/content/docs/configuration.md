---
title: Client configuration
weight: 40
---

Every setting is a functional option of `NewClient` (and of `NewClientFromEnv`). Their full list is in the [API reference](https://pkg.go.dev/github.com/boxboxjason/sonarqube-client-go/v2/sonar#ClientOptionFunc).

## Timeout

The HTTP client managed by the SDK times requests out after 30 seconds by default:

{{% example "ExampleWithTimeout" %}}

## HTTP client

To share an HTTP client with the rest of an application, pass it with `WithHTTPClient`. The timeout of that client then applies, `WithTimeout` is ignored.

{{% example "ExampleWithHTTPClient" %}}

## Transport

`WithTransportConfig` tunes the connection pool and the TLS settings of the SDK-managed transport:

{{% example "ExampleWithTransportConfig" %}}

{{< callout type="error" >}}
Never set `InsecureSkipVerify` outside of a trusted local instance: it disables the verification of the server certificate, which exposes the token and the source code metadata exchanged with SonarQube to man-in-the-middle attacks.
{{< /callout >}}

## Retries

Retries are disabled by default. `WithRetry` enables them, with exponential backoff and full jitter, for the listed status codes and for network errors:

{{% example "ExampleWithRetry" %}}

Only idempotent methods (`GET`, `HEAD`, `PUT`, `DELETE`, ...) are retried. `RetryNonIdempotent` extends retries to `POST` and `PATCH`: enable it only for endpoints that are safe to repeat, as a retried request can otherwise create a resource twice.

## Middleware and observability

`WithMiddleware` wraps the transport with `http.RoundTripper` middleware, the place for request logging, metrics or tracing. Middleware sits inside the retry layer, so it observes every attempt, not only the first one:

{{% example "ExampleWithMiddleware" %}}

For OpenTelemetry tracing, wrap [`otelhttp.NewTransport`](https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp#NewTransport) the same way:

```go
otelMiddleware := func(next http.RoundTripper) http.RoundTripper {
	return otelhttp.NewTransport(next)
}

client, err := sonar.NewClient(nil,
	sonar.WithBaseURL("https://sonar.example.com/api/"),
	sonar.WithToken(token),
	sonar.WithMiddleware(otelMiddleware),
)
```

## User agent

{{% example "ExampleWithUserAgent" %}}
