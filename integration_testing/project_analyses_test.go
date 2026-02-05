package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("ProjectAnalyses Service", Ordered, func() {
	var (
		client         *sonargo.Client
		cleanupManager *helpers.CleanupManager
		testProject    *sonargo.ProjectsCreate
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())

		cleanupManager = helpers.NewCleanupManager(client)

		// Create a test project for project analyses operations
		projectKey := helpers.UniqueResourceName("proj-analyses")
		testProject, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
			Name:    projectKey,
			Project: projectKey,
		})
		Expect(err).NotTo(HaveOccurred())
		cleanupManager.RegisterCleanup("project", testProject.Project.Key, func() error {
			_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
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
	// Search
	// =========================================================================
	Describe("Search", func() {
		Context("Functional Tests", func() {
			It("should search project analyses", func() {
				result, resp, err := client.ProjectAnalyses.Search(&sonargo.ProjectAnalysesSearchOption{
					Project: testProject.Project.Key,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				// Project may not have analyses yet, but the search should work
			})

			It("should search project analyses with pagination", func() {
				result, resp, err := client.ProjectAnalyses.Search(&sonargo.ProjectAnalysesSearchOption{
					Project: testProject.Project.Key,
					PaginationArgs: sonargo.PaginationArgs{
						Page:     1,
						PageSize: 10,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should search project analyses with category filter", func() {
				result, resp, err := client.ProjectAnalyses.Search(&sonargo.ProjectAnalysesSearchOption{
					Project:  testProject.Project.Key,
					Category: "VERSION",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should search project analyses with date range", func() {
				result, resp, err := client.ProjectAnalyses.Search(&sonargo.ProjectAnalysesSearchOption{
					Project: testProject.Project.Key,
					From:    "2020-01-01",
					To:      "2030-12-31",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("Parameter Validation", func() {
			It("should fail with missing project", func() {
				_, _, err := client.ProjectAnalyses.Search(&sonargo.ProjectAnalysesSearchOption{})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, _, err := client.ProjectAnalyses.Search(nil)
				Expect(err).To(HaveOccurred())
			})

			It("should fail with invalid category", func() {
				_, _, err := client.ProjectAnalyses.Search(&sonargo.ProjectAnalysesSearchOption{
					Project:  testProject.Project.Key,
					Category: "INVALID_CATEGORY",
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with invalid date format", func() {
				_, _, err := client.ProjectAnalyses.Search(&sonargo.ProjectAnalysesSearchOption{
					Project: testProject.Project.Key,
					From:    "invalid-date",
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with non-existent project", func() {
				_, resp, err := client.ProjectAnalyses.Search(&sonargo.ProjectAnalysesSearchOption{
					Project: "non-existent-project-12345",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	// =========================================================================
	// SearchAll
	// =========================================================================
	Describe("SearchAll", func() {
		Context("Functional Tests", func() {
			It("should search all project analyses", func() {
				result, resp, err := client.ProjectAnalyses.SearchAll(&sonargo.ProjectAnalysesSearchOption{
					Project: testProject.Project.Key,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				// Result may be empty/nil if no analyses exist
				_ = result
			})
		})

		Context("Parameter Validation", func() {
			It("should fail with missing project", func() {
				_, _, err := client.ProjectAnalyses.SearchAll(&sonargo.ProjectAnalysesSearchOption{})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, _, err := client.ProjectAnalyses.SearchAll(nil)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// =========================================================================
	// CreateEvent
	// =========================================================================
	Describe("CreateEvent", func() {
		Context("Parameter Validation", func() {
			It("should fail with missing analysis", func() {
				_, _, err := client.ProjectAnalyses.CreateEvent(&sonargo.ProjectAnalysesCreateEventOption{
					Name: "test-event",
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with missing name", func() {
				_, _, err := client.ProjectAnalyses.CreateEvent(&sonargo.ProjectAnalysesCreateEventOption{
					Analysis: "some-analysis-key",
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, _, err := client.ProjectAnalyses.CreateEvent(nil)
				Expect(err).To(HaveOccurred())
			})

			It("should fail with invalid category", func() {
				_, _, err := client.ProjectAnalyses.CreateEvent(&sonargo.ProjectAnalysesCreateEventOption{
					Analysis: "some-analysis-key",
					Name:     "test-event",
					Category: "INVALID_CATEGORY",
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with non-existent analysis", func() {
				_, resp, err := client.ProjectAnalyses.CreateEvent(&sonargo.ProjectAnalysesCreateEventOption{
					Analysis: "non-existent-analysis-12345",
					Name:     "test-event",
					Category: "VERSION",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	// =========================================================================
	// UpdateEvent
	// =========================================================================
	Describe("UpdateEvent", func() {
		Context("Parameter Validation", func() {
			It("should fail with missing event", func() {
				_, _, err := client.ProjectAnalyses.UpdateEvent(&sonargo.ProjectAnalysesUpdateEventOption{
					Name: "updated-name",
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with missing name", func() {
				_, _, err := client.ProjectAnalyses.UpdateEvent(&sonargo.ProjectAnalysesUpdateEventOption{
					Event: "some-event-key",
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, _, err := client.ProjectAnalyses.UpdateEvent(nil)
				Expect(err).To(HaveOccurred())
			})

			It("should fail with non-existent event", func() {
				_, resp, err := client.ProjectAnalyses.UpdateEvent(&sonargo.ProjectAnalysesUpdateEventOption{
					Event: "non-existent-event-12345",
					Name:  "updated-name",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	// =========================================================================
	// DeleteEvent
	// =========================================================================
	Describe("DeleteEvent", func() {
		Context("Parameter Validation", func() {
			It("should fail with missing event", func() {
				_, err := client.ProjectAnalyses.DeleteEvent(&sonargo.ProjectAnalysesDeleteEventOption{})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, err := client.ProjectAnalyses.DeleteEvent(nil)
				Expect(err).To(HaveOccurred())
			})

			It("should fail with non-existent event", func() {
				resp, err := client.ProjectAnalyses.DeleteEvent(&sonargo.ProjectAnalysesDeleteEventOption{
					Event: "non-existent-event-12345",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	// =========================================================================
	// Delete
	// =========================================================================
	Describe("Delete", func() {
		Context("Parameter Validation", func() {
			It("should fail with missing analysis", func() {
				_, err := client.ProjectAnalyses.Delete(&sonargo.ProjectAnalysesDeleteOption{})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, err := client.ProjectAnalyses.Delete(nil)
				Expect(err).To(HaveOccurred())
			})

			It("should fail with non-existent analysis", func() {
				resp, err := client.ProjectAnalyses.Delete(&sonargo.ProjectAnalysesDeleteOption{
					Analysis: "non-existent-analysis-12345",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})
})
