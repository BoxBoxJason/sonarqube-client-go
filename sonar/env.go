package sonar

import (
	"fmt"
	"os"
	"time"
)

// Environment variable names read by NewClientFromEnv.
const (
	// EnvURL is the base API URL of the SonarQube instance.
	EnvURL = "SONAR_URL"
	// EnvToken is the authentication token (takes precedence over basic auth).
	EnvToken = "SONAR_TOKEN"
	// EnvUsername is the username for basic authentication.
	EnvUsername = "SONAR_USERNAME"
	// EnvPassword is the password for basic authentication.
	EnvPassword = "SONAR_PASSWORD"
	// EnvUserAgent overrides the default User-Agent header.
	EnvUserAgent = "SONAR_USER_AGENT"
	// EnvTimeout sets the HTTP request timeout (parsed as a Go duration string, e.g. "30s", "2m").
	EnvTimeout = "SONAR_TIMEOUT"
)

// NewClientFromEnv creates a Client configured from SONAR_* environment
// variables, the convention used across the SDK and its integration tests:
//
//   - SONAR_URL        base API URL (defaults to the SDK default when unset)
//   - SONAR_TOKEN      token authentication (takes precedence over basic auth)
//   - SONAR_USERNAME   basic-auth username (used with SONAR_PASSWORD)
//   - SONAR_PASSWORD   basic-auth password
//   - SONAR_USER_AGENT optional User-Agent override
//   - SONAR_TIMEOUT    HTTP request timeout as a Go duration string (e.g. "30s", "2m")
//
// Additional functional options can be supplied and are applied after the
// environment is read, so they take precedence over the environment.
func NewClientFromEnv(options ...ClientOptionFunc) (*Client, error) {
	createOpts := &ClientCreateOptions{} //nolint:exhaustruct // populated from the environment below

	if baseURL := os.Getenv(EnvURL); baseURL != "" {
		createOpts.URL = &baseURL
	}

	if userAgent := os.Getenv(EnvUserAgent); userAgent != "" {
		createOpts.UserAgent = &userAgent
	}

	if raw := os.Getenv(EnvTimeout); raw != "" {
		d, err := time.ParseDuration(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid %s %q: %w", EnvTimeout, raw, err)
		}

		createOpts.Timeout = &d
	}

	// Token authentication takes precedence over basic auth when both are set.
	if token := os.Getenv(EnvToken); token != "" {
		createOpts.Token = &token
	} else {
		username := os.Getenv(EnvUsername)
		password := os.Getenv(EnvPassword)

		if username != "" && password != "" {
			createOpts.Username = &username
			createOpts.Password = &password
		}
	}

	return NewClient(createOpts, options...)
}
