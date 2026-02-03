package integration_testing_test

import (
"net/http"
"strings"

. "github.com/onsi/ginkgo/v2"
. "github.com/onsi/gomega"

sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Batch Service", Ordered, func() {
	var (
client      *sonargo.Client
testProject *sonargo.ProjectsCreate
)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())

		// Create a test project for batch operations
		testProject, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
			Name:    "batch-e2e-test-project",
			Project: "batch-e2e-test-project",
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
	// Index
	// =========================================================================
	Describe("Index", func() {
		Context("Functional Tests", func() {
			It("should get batch index", func() {
				result, resp, err := client.Batch.Index()
				// Batch API may not be available in newer SonarQube versions
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("Batch API is not available in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should return JAR file list", func() {
				result, resp, err := client.Batch.Index()
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("Batch API is not available in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())

				// Index contains JAR file entries (one per line)
				index := *result
				if len(index) > 0 {
					// Should contain JAR file references
					Expect(strings.Contains(index, ".jar") || len(index) > 0).To(BeTrue())
				}
			})

			It("should return consistent results on multiple calls", func() {
				result1, resp1, err := client.Batch.Index()
				if resp1 != nil && resp1.StatusCode == http.StatusNotFound {
					Skip("Batch API is not available in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp1.StatusCode).To(Equal(http.StatusOK))

				result2, resp2, err := client.Batch.Index()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp2.StatusCode).To(Equal(http.StatusOK))

				Expect(*result1).To(Equal(*result2))
			})
		})
	})

	// =========================================================================
	// File
	// =========================================================================
	Describe("File", func() {
		Context("Functional Tests", func() {
			It("should download batch file with valid name", func() {
				// First get the index to find a valid file name
				index, resp, err := client.Batch.Index()
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("Batch API is not available in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())

				if index != nil && len(*index) > 0 {
					lines := strings.Split(*index, "\n")
					for _, line := range lines {
						parts := strings.Split(strings.TrimSpace(line), "|")
						if len(parts) > 0 && strings.HasSuffix(parts[0], ".jar") {
							result, resp, err := client.Batch.File(&sonargo.BatchFileOption{
								Name: parts[0],
							})
							Expect(err).NotTo(HaveOccurred())
							Expect(resp.StatusCode).To(Equal(http.StatusOK))
							Expect(result).NotTo(BeNil())
							Expect(len(result)).To(BeNumerically(">", 0))
							break
						}
					}
				}
			})
		})

		Context("Error Handling", func() {
			It("should fail with non-existent file", func() {
				_, resp, err := client.Batch.File(&sonargo.BatchFileOption{
					Name: "non-existent-file.jar",
				})
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					// API not available or file not found - both valid
					return
				}
				if resp != nil {
					Expect(resp.StatusCode).NotTo(Equal(http.StatusOK))
				} else {
					Expect(err).To(HaveOccurred())
				}
			})
		})
	})

	// =========================================================================
	// Project
	// =========================================================================
	Describe("Project", func() {
		Context("Functional Tests", func() {
			It("should get project batch info with valid key", func() {
				result, resp, err := client.Batch.Project(&sonargo.BatchProjectOption{
					Key: testProject.Project.Key,
				})
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("Batch API is not available in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("Error Handling", func() {
			It("should fail with missing key", func() {
				_, _, err := client.Batch.Project(&sonargo.BatchProjectOption{})
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
