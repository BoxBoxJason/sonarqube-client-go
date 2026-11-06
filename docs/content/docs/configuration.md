---
title: Client configuration
description: "Configure the sonarqube-client-go client: timeouts, custom HTTP clients and transports, retries with backoff, middleware and the user agent."
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

Two options bound how long a request can take when the server misbehaves:

- `RetryableError` chooses which transport errors are retried, all of them when unset. Each attempt can take as long as the transport's own timeouts allow, so a server that stops answering costs that much per attempt. Returning `false` for timeouts fails the request after the first one:

  ```go
  sonar.WithRetry(sonar.RetryOptions{
  	MaxAttempts:          3,
  	InitialDelay:         500 * time.Millisecond,
  	MaxDelay:             5 * time.Second,
  	RetryableStatusCodes: []int{502, 503, 504},
  	RetryableError: func(err error) bool {
  		var netErr net.Error

  		return !errors.As(err, &netErr) || !netErr.Timeout()
  	},
  })
  ```

- `MaxRetryAfter` limits the `Retry-After` wait honoured on a `429`, which `MaxDelay` does not cap. A response asking for longer is returned to the caller right away, rather than retried early or waited out. There is no limit when unset.

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
