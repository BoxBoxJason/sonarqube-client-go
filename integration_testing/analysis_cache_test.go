package integration_testing_test

import (
"net/http"

. "github.com/onsi/ginkgo/v2"
. "github.com/onsi/gomega"

sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("AnalysisCache Service", Ordered, func() {
	var (
client      *sonargo.Client
testProject *sonargo.ProjectsCreate
)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())

		// Create a test project for analysis cache operations
		testProject, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
			Name:    "analysis-cache-e2e-test-project",
			Project: "analysis-cache-e2e-test-project",
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
	// Clear
	// =========================================================================
	Describe("Clear", func() {
		Context("Functional Tests", func() {
			It("should clear all cached data without options", func() {
				resp, err := client.AnalysisCache.Clear(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})

			It("should clear all cached data with empty options", func() {
				resp, err := client.AnalysisCache.Clear(&sonargo.AnalysisCacheClearOption{})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})

			It("should clear cached data for a specific project", func() {
				resp, err := client.AnalysisCache.Clear(&sonargo.AnalysisCacheClearOption{
					Project: testProject.Project.Key,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})

			It("should clear cached data for a specific project branch", func() {
				resp, err := client.AnalysisCache.Clear(&sonargo.AnalysisCacheClearOption{
					Project: testProject.Project.Key,
					Branch:  "main",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})
		})

		Context("Error Handling", func() {
			It("should fail when branch is specified without project", func() {
				_, err := client.AnalysisCache.Clear(&sonargo.AnalysisCacheClearOption{
					Branch: "main",
				})
				Expect(err).To(HaveOccurred())
			})

			It("should succeed with non-existent project (idempotent)", func() {
				resp, err := client.AnalysisCache.Clear(&sonargo.AnalysisCacheClearOption{
					Project: "non-existent-project-12345",
				})
				// Clear is idempotent - should succeed even for non-existent project
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})
		})
	})

	// =========================================================================
	// Get
	// =========================================================================
	Describe("Get", func() {
		Context("Functional Tests", func() {
			It("should get cached data for a project", func() {
				resp, err := client.AnalysisCache.Get(&sonargo.AnalysisCacheGetOption{
					Project: testProject.Project.Key,
				})
				// May return 404 if no cache exists yet (project not analyzed)
				if err == nil {
					Expect(resp.StatusCode).To(SatisfyAny(
Equal(http.StatusOK),
Equal(http.StatusNoContent),
))
					if resp.Body != nil {
						_ = resp.Body.Close()
					}
				} else {
					Expect(resp).NotTo(BeNil())
					Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				}
			})

			It("should get cached data for a specific branch", func() {
				resp, err := client.AnalysisCache.Get(&sonargo.AnalysisCacheGetOption{
					Project: testProject.Project.Key,
					Branch:  "main",
				})
				// May return 404 if no cache exists yet (branch not analyzed)
				if err == nil {
					Expect(resp.StatusCode).To(SatisfyAny(
Equal(http.StatusOK),
Equal(http.StatusNoContent),
))
					if resp.Body != nil {
						_ = resp.Body.Close()
					}
				} else {
					Expect(resp).NotTo(BeNil())
					Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				}
			})
		})

		Context("Error Handling", func() {
			It("should fail with nil options", func() {
				_, err := client.AnalysisCache.Get(nil)
				Expect(err).To(HaveOccurred())
			})

			It("should fail with missing project", func() {
				_, err := client.AnalysisCache.Get(&sonargo.AnalysisCacheGetOption{})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with non-existent project", func() {
				resp, err := client.AnalysisCache.Get(&sonargo.AnalysisCacheGetOption{
					Project: "non-existent-project-12345",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})
})
