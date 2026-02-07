package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("ProjectDump Service", Ordered, func() {
	var (
		client  *sonar.Client
		cleanup *helpers.CleanupManager
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
		cleanup = helpers.NewCleanupManager(client)
	})

	AfterAll(func() {
		cleanup.Cleanup()
	})

	// =========================================================================
	// Export
	// =========================================================================
	Describe("Export", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.ProjectDump.Export(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required key", func() {
				result, resp, err := client.ProjectDump.Export(&sonar.ProjectDumpExportOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should fail for non-existent project", func() {
				result, resp, err := client.ProjectDump.Export(&sonar.ProjectDumpExportOption{
					Key: "non-existent-project-key",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				Expect(result).To(BeNil())
			})

			It("should successfully trigger export for existing project", func() {
				// Create a project for testing
				projectKey := helpers.UniqueResourceName("dump-export")
				_, resp, err := client.Projects.Create(&sonar.ProjectsCreateOption{
					Name:    projectKey,
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				cleanup.RegisterCleanup("project", projectKey, func() error {
					_, err := client.Projects.Delete(&sonar.ProjectsDeleteOption{Project: projectKey})
					return err
				})

				// Trigger export
				result, resp, err := client.ProjectDump.Export(&sonar.ProjectDumpExportOption{
					Key: projectKey,
				})

				// Export may succeed or fail depending on server configuration
				// Community edition may not support project export
				if err != nil {
					// Accept that export may not be available
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					Expect(result).NotTo(BeNil())
					Expect(result.ProjectKey).To(Equal(projectKey))
				}
			})
		})
	})

	// =========================================================================
	// Status
	// =========================================================================
	Describe("Status", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.ProjectDump.Status(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without id or key", func() {
				result, resp, err := client.ProjectDump.Status(&sonar.ProjectDumpStatusOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ID or Key"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should fail for non-existent project by key", func() {
				result, resp, err := client.ProjectDump.Status(&sonar.ProjectDumpStatusOption{
					Key: "non-existent-project-key",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				Expect(result).To(BeNil())
			})

			It("should get status for existing project by key", func() {
				// Create a project for testing
				projectKey := helpers.UniqueResourceName("dump-status")
				_, resp, err := client.Projects.Create(&sonar.ProjectsCreateOption{
					Name:    projectKey,
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				cleanup.RegisterCleanup("project", projectKey, func() error {
					_, err := client.Projects.Delete(&sonar.ProjectsDeleteOption{Project: projectKey})
					return err
				})

				// Get status
				result, resp, err := client.ProjectDump.Status(&sonar.ProjectDumpStatusOption{
					Key: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				// CanBeExported may be true or false depending on server edition
			})
		})
	})

	// =========================================================================
	// Full Workflow
	// =========================================================================
	Describe("Full Workflow", func() {
		Context("Project Export Status Lifecycle", func() {
			It("should handle project dump status workflow", func() {
				// Create a project
				projectKey := helpers.UniqueResourceName("dump-workflow")
				_, resp, err := client.Projects.Create(&sonar.ProjectsCreateOption{
					Name:    projectKey,
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				cleanup.RegisterCleanup("project", projectKey, func() error {
					_, err := client.Projects.Delete(&sonar.ProjectsDeleteOption{Project: projectKey})
					return err
				})

				// Check initial status
				status, resp, err := client.ProjectDump.Status(&sonar.ProjectDumpStatusOption{
					Key: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(status).NotTo(BeNil())

				// If export is available, try to trigger it
				if status.CanBeExported {
					result, resp, err := client.ProjectDump.Export(&sonar.ProjectDumpExportOption{
						Key: projectKey,
					})
					if err == nil {
						Expect(resp.StatusCode).To(Equal(http.StatusOK))
						Expect(result).NotTo(BeNil())
						Expect(result.ProjectKey).To(Equal(projectKey))
					}
				}
			})
		})
	})
})
