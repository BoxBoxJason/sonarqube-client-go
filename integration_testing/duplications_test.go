package integration_testing_test

import (
"net/http"

. "github.com/onsi/ginkgo/v2"
. "github.com/onsi/gomega"

sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Duplications Service", Ordered, func() {
	var (
client     *sonargo.Client
cleanup    *helpers.CleanupManager
projectKey string
)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
		cleanup = helpers.NewCleanupManager(client)

		// Create a test project for duplications-related operations
		projectKey = helpers.UniqueResourceName("dup")
		_, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
			Name:    "Duplications Test Project",
			Project: projectKey,
		})
		Expect(err).NotTo(HaveOccurred())

		cleanup.RegisterCleanup("project", projectKey, func() error {
			_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
				Project: projectKey,
			})
			return err
		})
	})

	AfterAll(func() {
		errors := cleanup.Cleanup()
		for _, err := range errors {
			GinkgoWriter.Printf("Cleanup error: %v\n", err)
		}
	})

	// =========================================================================
	// Show
	// =========================================================================
	Describe("Show", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Duplications.Show(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required key", func() {
				result, resp, err := client.Duplications.Show(&sonargo.DuplicationsShowOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent File", func() {
			It("should fail for non-existent file key", func() {
				result, resp, err := client.Duplications.Show(&sonargo.DuplicationsShowOption{
					Key: "non-existent-file-key",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})

		Context("Project Key", func() {
			It("should work with project key (empty duplications)", func() {
				result, resp, err := client.Duplications.Show(&sonargo.DuplicationsShowOption{
					Key: projectKey,
				})
				// May succeed with empty result or fail if the project doesn't have analyzed files
if err == nil {
Expect(resp.StatusCode).To(Equal(http.StatusOK))
Expect(result).NotTo(BeNil())
}
})
})
})
})
