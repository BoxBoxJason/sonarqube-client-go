package helpers

import (
	"fmt"
	"os"
	"strings"
	"time"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"
)

const (
	// DefaultBaseURL is the default SonarQube URL for e2e tests.
	DefaultBaseURL = "http://127.0.0.1:9000"

	// DefaultUsername is the default admin username.
	DefaultUsername = "admin"

	// DefaultPassword is the default admin password.
	DefaultPassword = "admin"

	// E2EResourcePrefix is the prefix for all e2e test resources.
	E2EResourcePrefix = "e2e-"

	// DefaultTimeout is the default timeout for polling operations.
	DefaultTimeout = 60 * time.Second

	// DefaultPollInterval is the default interval between poll attempts.
	DefaultPollInterval = 2 * time.Second
)

// Config holds configuration for e2e tests.
type Config struct {
	BaseURL  string
	Username string
	Password string
	Token    string
}

// LoadConfig loads e2e test configuration from environment variables.
func LoadConfig() *Config {
	cfg := &Config{
		BaseURL:  DefaultBaseURL,
		Username: DefaultUsername,
		Password: DefaultPassword,
		Token:    "",
	}

	if url := os.Getenv("SONAR_URL"); url != "" {
		cfg.BaseURL = url
	}

	if username := os.Getenv("SONAR_USERNAME"); username != "" {
		cfg.Username = username
	}

	if password := os.Getenv("SONAR_PASSWORD"); password != "" {
		cfg.Password = password
	}

	if token := os.Getenv("SONAR_TOKEN"); token != "" {
		cfg.Token = token
	}

	return cfg
}

// NewClient creates a new SonarQube client for e2e tests using the provided config.
func NewClient(cfg *Config) (*sonargo.Client, error) {
	baseURL := cfg.BaseURL
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	// Ensure the URL ends with /api/ as expected by the SDK
	if !strings.HasSuffix(baseURL, "api/") {
		baseURL += "api/"
	}

	// Prefer token auth if available
	if cfg.Token != "" {
		client, err := sonargo.NewClient(nil, sonargo.WithBaseURL(baseURL), sonargo.WithToken(cfg.Token))
		if err != nil {
			return nil, fmt.Errorf("failed to create client with token: %w", err)
		}

		return client, nil
	}

	// Fall back to basic auth
	client, err := sonargo.NewClient(nil, sonargo.WithBaseURL(baseURL), sonargo.WithBasicAuth(cfg.Username, cfg.Password))
	if err != nil {
		return nil, fmt.Errorf("failed to create client with basic auth: %w", err)
	}

	return client, nil
}

// NewDefaultClient creates a new SonarQube client with default configuration.
func NewDefaultClient() (*sonargo.Client, error) {
	cfg := LoadConfig()

	return NewClient(cfg)
}

// UniqueResourceName generates a unique name for e2e test resources.
func UniqueResourceName(prefix string) string {
	timestamp := time.Now().Format("20060102-150405")

	return fmt.Sprintf("%s%s-%s", E2EResourcePrefix, prefix, timestamp)
}

// UniqueResourceNameWithSuffix generates a unique name with an additional suffix.
func UniqueResourceNameWithSuffix(prefix, suffix string) string {
	timestamp := time.Now().Format("20060102-150405")

	return fmt.Sprintf("%s%s-%s-%s", E2EResourcePrefix, prefix, timestamp, suffix)
}
