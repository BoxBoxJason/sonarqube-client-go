// Package sonar is a Go client for the SonarQube web API.
//
// The root type is Client. Create one with NewClient and reach each API area
// through a typed service field, for example client.Projects, client.Issues or
// client.Qualitygates. V2 API services are grouped under client.V2.
//
// # Authentication
//
// SonarQube supports token authentication (recommended) and basic auth. Both
// can be supplied at construction via ClientCreateOptions or with the WithToken
// and WithBasicAuth options:
//
//	client, err := sonar.NewClient(nil,
//		sonar.WithBaseURL("https://sonarqube.example.com/api/"),
//		sonar.WithToken("my-token"),
//	)
//
// In CI/CD environments the SONAR_URL, SONAR_TOKEN, SONAR_USERNAME, and
// SONAR_PASSWORD environment variables can be used instead:
//
//	client, err := sonar.NewClientFromEnv()
//
// # Requests and responses
//
// Every service method takes a context.Context and, where applicable, a typed
// *<Service><Method>Options struct whose fields are validated before the request
// is sent. Methods return the decoded response, the raw *http.Response, and an
// error.
//
// # V1 vs V2
//
// V1 endpoints encode parameters as URL query values (url:"" struct tags). V2
// endpoints use JSON tags, support a request body, and live under the V2 field;
// the "v2/" path prefix is added automatically.
//
// # Errors
//
// API errors are returned as *ResponseError. Use the sentinel helpers
// (IsNotFound, IsUnauthorized, IsForbidden, IsConflict, IsRateLimited,
// IsServerError) to branch on the HTTP status without unwrapping by hand.
//
// # Pagination
//
// Paginated V1 endpoints expose both a single-page method (Search) and a
// convenience method that fetches every page (SearchAll / ListAll).
//
// # Resilience
//
// Retries with exponential backoff and jitter are opt-in via WithRetry. The
// transport can be customized with WithTransportConfig, and arbitrary
// http.RoundTripper middleware (logging, tracing, metrics) can be attached with
// WithMiddleware.
package sonar
