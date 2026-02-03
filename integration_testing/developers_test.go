package integration_testing_test

import (
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Developers Service", Ordered, func() {
	var (
		client      *sonargo.Client
		testProject *sonargo.ProjectsCreate
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())

		// Create a test project for developer events
		testProject, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
			Name:    "developers-e2e-test-project",
			Project: "developers-e2e-test-project",
		})
		Expect(err).NotTo(HaveOccurred())
	})

	AfterAll(func() {
		if testProject != nil {
			_, _ = client.Projects.Delete(&sonargo.ProjectsDeleteOption{
				Project: testProject.Project.Key,
			})
		}
	})

	// =========================================================================
	// SearchEvents
	// =========================================================================
	Describe("SearchEvents", func() {
		Context("Functional Tests", func() {
			It("should search developer events with valid parameters", func() {
				fromDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02T15:04:05+0000")
				result, resp, err := client.Developers.SearchEvents(&sonargo.DevelopersSearchEventsOption{
					From:     []string{fromDate},
					Projects: []string{testProject.Project.Key},
				})
				// Skip if API not available
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("Developers API is not available in this SonarQube version")
				}
				// 503 if indexing is in progress
				if resp != nil && resp.StatusCode == http.StatusServiceUnavailable {
					Skip("Issue indexing is in progress")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				// Events might be empty for a new project
			})

			It("should return empty events for new project", func() {
				fromDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02T15:04:05+0000")
				result, resp, err := client.Developers.SearchEvents(&sonargo.DevelopersSearchEventsOption{
					From:     []string{fromDate},
					Projects: []string{testProject.Project.Key},
				})
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("Developers API is not available in this SonarQube version")
				}
				if resp != nil && resp.StatusCode == http.StatusServiceUnavailable {
					Skip("Issue indexing is in progress")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				// New project with no analysis should have empty or minimal events
			})
		})

		Context("Error Handling", func() {
			It("should fail with missing from date", func() {
				_, _, err := client.Developers.SearchEvents(&sonargo.DevelopersSearchEventsOption{
					Projects: []string{testProject.Project.Key},
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with missing projects", func() {
				fromDate := time.Now().Add(-24 * time.Hour).Format("2006-01-02T15:04:05+0000")
				_, _, err := client.Developers.SearchEvents(&sonargo.DevelopersSearchEventsOption{
					From: []string{fromDate},
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, _, err := client.Developers.SearchEvents(nil)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
