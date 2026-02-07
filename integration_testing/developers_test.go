package integration_testing_test

import (
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("Developers Service", Ordered, func() {
	var (
		client         *sonar.Client
		cleanupManager *helpers.CleanupManager
		testProject    *sonar.ProjectsCreate
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())

		cleanupManager = helpers.NewCleanupManager(client)

		// Create a uniquely named test project for developer events to avoid collisions across runs.
		projectKey := helpers.UniqueResourceName("developers")
		testProject, _, err = client.Projects.Create(&sonar.ProjectsCreateOption{
			Name:    projectKey,
			Project: projectKey,
		})
		Expect(err).NotTo(HaveOccurred())
		cleanupManager.RegisterCleanup("project", projectKey, func() error {
			_, err := client.Projects.Delete(&sonar.ProjectsDeleteOption{
				Project: testProject.Project.Key,
			})
			return err
		})
	})

	AfterAll(func() {
		errors := cleanupManager.Cleanup()
		for _, err := range errors {
			GinkgoWriter.Printf("Cleanup error: %v\n", err)
		}
	})

	// =========================================================================
	// SearchEvents
	// =========================================================================
	Describe("SearchEvents", func() {
		Context("Functional Tests", func() {
			It("should search developer events with valid parameters", func() {
				fromDate := time.Now().UTC().Add(-24 * time.Hour).Format("2006-01-02T15:04:05-0700")
				result, resp, err := client.Developers.SearchEvents(&sonar.DevelopersSearchEventsOption{
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
				// Events might be empty for a new project with no analysis
				Expect(result.Events).NotTo(BeNil())
			})
		})

		Context("Error Handling", func() {
			It("should fail with missing from date", func() {
				_, resp, err := client.Developers.SearchEvents(&sonar.DevelopersSearchEventsOption{
					Projects: []string{testProject.Project.Key},
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing projects", func() {
				fromDate := time.Now().UTC().Add(-24 * time.Hour).Format("2006-01-02T15:04:05-0700")
				_, resp, err := client.Developers.SearchEvents(&sonar.DevelopersSearchEventsOption{
					From: []string{fromDate},
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with nil options", func() {
				_, resp, err := client.Developers.SearchEvents(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})
})
