package sonar_test

import (
	"context"
	"fmt"
	"time"

	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

// Create a client with token authentication and search for projects.
func Example() {
	client, err := sonar.NewClient(nil,
		sonar.WithBaseURL("https://sonarqube.example.com/api/"),
		sonar.WithToken("my-token"),
	)
	if err != nil {
		panic(err)
	}

	result, _, err := client.Projects.Search(context.Background(), &sonar.ProjectsSearchOptions{
		Query: "my-service",
	})
	if err != nil {
		panic(err)
	}

	for _, project := range result.Components {
		fmt.Println(project.Key)
	}
}

// Fetch every page of a paginated endpoint with the SearchAll helper.
func ExampleProjectsService_SearchAll() {
	client, err := sonar.NewClient(nil,
		sonar.WithBaseURL("https://sonarqube.example.com/api/"),
		sonar.WithToken("my-token"),
	)
	if err != nil {
		panic(err)
	}

	//nolint:exhaustruct // only the filter we care about is set
	projects, _, err := client.Projects.SearchAll(context.Background(), &sonar.ProjectsSearchOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("found %d projects\n", len(projects))
}

// Enable retries with exponential backoff for transient server errors.
func ExampleWithRetry() {
	client, err := sonar.NewClient(nil,
		sonar.WithBaseURL("https://sonarqube.example.com/api/"),
		sonar.WithToken("my-token"),
		sonar.WithRetry(sonar.RetryOptions{
			MaxAttempts:          3,
			InitialDelay:         100 * time.Millisecond,
			MaxDelay:             2 * time.Second,
			RetryableStatusCodes: []int{502, 503, 504},
		}),
	)
	if err != nil {
		panic(err)
	}

	_ = client
	// Output:
}

// Branch on the HTTP status of an API error using the sentinel helpers.
func ExampleIsNotFound() {
	client, err := sonar.NewClient(nil,
		sonar.WithBaseURL("https://sonarqube.example.com/api/"),
		sonar.WithToken("my-token"),
	)
	if err != nil {
		panic(err)
	}

	_, err = client.Projects.Delete(context.Background(), &sonar.ProjectsDeleteOptions{Project: "missing-project"})
	if sonar.IsNotFound(err) {
		fmt.Println("project not found")
	}
}

// Configure a client from SONAR_URL and SONAR_TOKEN environment variables.
func ExampleNewClientFromEnv() {
	client, err := sonar.NewClientFromEnv()
	if err != nil {
		panic(err)
	}

	_ = client
	// Output:
}
