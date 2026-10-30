package sonar_test

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

// The examples below are rendered verbatim in the documentation site
// (docs/content), so they must stay short and self-explanatory: the docs build
// fails when one of them is renamed or removed.

// Create a client from ClientCreateOptions, with token authentication.
func ExampleNewClient() {
	url := "https://sonarqube.example.com/api/"
	token := "my-token"

	client, err := sonar.NewClient(&sonar.ClientCreateOptions{
		URL:   &url,
		Token: &token,
	})
	if err != nil {
		panic(err)
	}

	_ = client
}

// Authenticate with a username and a password instead of a token.
func ExampleWithBasicAuth() {
	client, err := sonar.NewClient(nil,
		sonar.WithBaseURL("https://sonarqube.example.com/api/"),
		sonar.WithBasicAuth("admin", "my-password"),
	)
	if err != nil {
		panic(err)
	}

	_ = client
}

// Override the default 30 second request timeout of the SDK-managed HTTP client.
func ExampleWithTimeout() {
	client, err := sonar.NewClient(nil,
		sonar.WithBaseURL("https://sonarqube.example.com/api/"),
		sonar.WithToken("my-token"),
		sonar.WithTimeout(time.Minute),
	)
	if err != nil {
		panic(err)
	}

	_ = client
}

// Bring your own HTTP client, for example one shared with the rest of an application.
func ExampleWithHTTPClient() {
	httpClient := &http.Client{Timeout: time.Minute}

	client, err := sonar.NewClient(nil,
		sonar.WithBaseURL("https://sonarqube.example.com/api/"),
		sonar.WithToken("my-token"),
		sonar.WithHTTPClient(httpClient),
	)
	if err != nil {
		panic(err)
	}

	_ = client
}

// Tune the connection pool and TLS settings of the SDK-managed transport.
func ExampleWithTransportConfig() {
	client, err := sonar.NewClient(nil,
		sonar.WithBaseURL("https://sonarqube.example.com/api/"),
		sonar.WithToken("my-token"),
		sonar.WithTransportConfig(sonar.TransportConfig{
			TLSClientConfig:     &tls.Config{MinVersion: tls.VersionTLS12},
			MaxIdleConns:        20,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
		}),
	)
	if err != nil {
		panic(err)
	}

	_ = client
}

// Identify your tool in the SonarQube access logs.
func ExampleWithUserAgent() {
	client, err := sonar.NewClient(nil,
		sonar.WithBaseURL("https://sonarqube.example.com/api/"),
		sonar.WithToken("my-token"),
		sonar.WithUserAgent("my-tool/1.0"),
	)
	if err != nil {
		panic(err)
	}

	_ = client
}

// Fail a pipeline when the quality gate of a project is not passing.
func ExampleQualitygatesService_ProjectStatus() {
	client, err := sonar.NewClientFromEnv()
	if err != nil {
		panic(err)
	}

	status, _, err := client.Qualitygates.ProjectStatus(context.Background(), &sonar.QualitygatesProjectStatusOptions{
		ProjectKey: "my-project",
	})
	if err != nil {
		panic(err)
	}

	if status.ProjectStatus.Status != "OK" {
		panic("quality gate failed: " + status.ProjectStatus.Status)
	}
}

// Read one page of open issues of a project.
func ExampleIssuesService_Search() {
	client, err := sonar.NewClientFromEnv()
	if err != nil {
		panic(err)
	}

	issues, _, err := client.Issues.Search(context.Background(), &sonar.IssuesSearchOptions{
		PaginationArgs: sonar.PaginationArgs{Page: 1, PageSize: 100},
		Projects:       []string{"my-project"},
		IssueStatuses:  []string{"OPEN", "CONFIRMED"},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("page %d: %d of %d issues\n", issues.Paging.PageIndex, len(issues.Issues), issues.Paging.Total)

	for _, issue := range issues.Issues {
		fmt.Println(issue.Key, issue.Message)
	}
}

// Fetch every matching user, across all pages.
func ExampleUsersService_SearchAll() {
	client, err := sonar.NewClientFromEnv()
	if err != nil {
		panic(err)
	}

	users, _, err := client.Users.SearchAll(context.Background(), &sonar.UsersSearchOptions{
		Query: "john",
	})
	if err != nil {
		panic(err)
	}

	for _, user := range users {
		fmt.Println(user.Login)
	}
}

// Call a V2 API service: they are grouped under client.V2.
func ExampleSystemServiceV2_GetHealth() {
	client, err := sonar.NewClientFromEnv()
	if err != nil {
		panic(err)
	}

	health, _, err := client.V2.System.GetHealth(context.Background(), &sonar.SystemPasscodeOptionV2{
		Passcode: "my-monitoring-passcode",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(health.Status)
}

// Options are validated before any request is sent: a missing or invalid
// parameter is reported as a *ValidationError, an API failure as a *ResponseError.
func ExampleValidationError() {
	client, err := sonar.NewClientFromEnv()
	if err != nil {
		panic(err)
	}

	_, _, err = client.Qualitygates.ProjectStatus(context.Background(), &sonar.QualitygatesProjectStatusOptions{})

	var validationErr *sonar.ValidationError
	if errors.As(err, &validationErr) {
		fmt.Println("invalid field:", validationErr.Field)
	}

	var responseErr *sonar.ResponseError
	if errors.As(err, &responseErr) {
		fmt.Println("SonarQube answered", responseErr.StatusCode)
	}
}
