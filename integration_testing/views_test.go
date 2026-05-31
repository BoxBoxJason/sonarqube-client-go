package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("Views Service", Ordered, func() {
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
	// Create
	// =========================================================================
	Describe("Create", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Views.Create(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required name", func() {
				resp, err := client.Views.Create(&sonar.ViewsCreateOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Name"))
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid visibility", func() {
				resp, err := client.Views.Create(&sonar.ViewsCreateOptions{
					Name:       "Test View",
					Visibility: "invalid",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should create or return an enterprise-only error", func() {
				viewKey := helpers.UniqueResourceName("view")
				resp, err := client.Views.Create(&sonar.ViewsCreateOptions{
					Name: viewKey,
					Key:  viewKey,
				})
				if err != nil {
					// Community edition returns 404 for enterprise-only endpoints.
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					cleanup.RegisterCleanup("view", viewKey, func() error {
						_, cleanupErr := client.Views.Delete(&sonar.ViewsDeleteOptions{
							Key: viewKey,
						})
						return cleanupErr
					})
				}
			})
		})
	})

	// =========================================================================
	// List
	// =========================================================================
	Describe("List", func() {
		Context("Functional Tests", func() {
			It("should return a list or an enterprise-only error", func() {
				result, resp, err := client.Views.List()
				if err != nil {
					// Community edition does not expose views endpoints.
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp).NotTo(BeNil())
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})

	// =========================================================================
	// Search
	// =========================================================================
	Describe("Search", func() {
		Context("Parameter Validation", func() {
			It("should fail with invalid pagination", func() {
				_, _, err := client.Views.Search(&sonar.ViewsSearchOptions{
					PaginationArgs: sonar.PaginationArgs{Page: -1},
				})
				Expect(err).To(HaveOccurred())
			})
		})

		Context("Functional Tests", func() {
			It("should return results or an enterprise-only error", func() {
				result, resp, err := client.Views.Search(nil)
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp).NotTo(BeNil())
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})

	// =========================================================================
	// Show
	// =========================================================================
	Describe("Show", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				_, _, err := client.Views.Show(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
			})

			It("should fail without key", func() {
				_, _, err := client.Views.Show(&sonar.ViewsShowOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
			})
		})

		Context("Functional Tests", func() {
			It("should fail for a non-existent view", func() {
				result, resp, err := client.Views.Show(&sonar.ViewsShowOptions{
					Key: "non-existent-view-key",
				})
				Expect(err).To(HaveOccurred())
				// Either 404 (not found) or 404 (enterprise endpoint not available).
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Update
	// =========================================================================
	Describe("Update", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Views.Update(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required key", func() {
				resp, err := client.Views.Update(&sonar.ViewsUpdateOptions{
					Name: "New Name",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required name", func() {
				resp, err := client.Views.Update(&sonar.ViewsUpdateOptions{
					Key: "my-view",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Name"))
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Delete
	// =========================================================================
	Describe("Delete", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Views.Delete(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without key", func() {
				resp, err := client.Views.Delete(&sonar.ViewsDeleteOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should fail for a non-existent view", func() {
				resp, err := client.Views.Delete(&sonar.ViewsDeleteOptions{
					Key: "non-existent-view",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
			})
		})
	})

	// =========================================================================
	// AddProject / RemoveProject
	// =========================================================================
	Describe("AddProject", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Views.AddProject(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without key", func() {
				resp, err := client.Views.AddProject(&sonar.ViewsAddProjectOptions{
					Project: "my-project",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(resp).To(BeNil())
			})

			It("should fail without project key", func() {
				resp, err := client.Views.AddProject(&sonar.ViewsAddProjectOptions{
					Key: "my-view",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Project"))
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("RemoveProject", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Views.RemoveProject(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without key", func() {
				resp, err := client.Views.RemoveProject(&sonar.ViewsRemoveProjectOptions{
					Project: "my-project",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(resp).To(BeNil())
			})

			It("should fail without project key", func() {
				resp, err := client.Views.RemoveProject(&sonar.ViewsRemoveProjectOptions{
					Key: "my-view",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Project"))
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// MoveOptions
	// =========================================================================
	Describe("MoveOptions", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				_, _, err := client.Views.MoveOptions(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
			})

			It("should fail without key", func() {
				_, _, err := client.Views.MoveOptions(&sonar.ViewsMoveOptionsOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
			})
		})
	})

	// =========================================================================
	// Move
	// =========================================================================
	Describe("Move", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Views.Move(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without key", func() {
				resp, err := client.Views.Move(&sonar.ViewsMoveOptions{
					Destination: "dest",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(resp).To(BeNil())
			})

			It("should fail without destination", func() {
				resp, err := client.Views.Move(&sonar.ViewsMoveOptions{
					Key: "my-view",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Destination"))
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Projects
	// =========================================================================
	Describe("Projects", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				_, _, err := client.Views.Projects(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
			})

			It("should fail without key", func() {
				_, _, err := client.Views.Projects(&sonar.ViewsProjectsOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
			})

			It("should fail with invalid selected value", func() {
				_, _, err := client.Views.Projects(&sonar.ViewsProjectsOptions{
					Key:      "my-view",
					Selected: "invalid",
				})
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// =========================================================================
	// AddApplication / RemoveApplication
	// =========================================================================
	Describe("AddApplication", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Views.AddApplication(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without portfolio key", func() {
				resp, err := client.Views.AddApplication(&sonar.ViewsAddApplicationOptions{
					Application: "my-app",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Portfolio"))
				Expect(resp).To(BeNil())
			})

			It("should fail without application key", func() {
				resp, err := client.Views.AddApplication(&sonar.ViewsAddApplicationOptions{
					Portfolio: "my-portfolio",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Application"))
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("RemoveApplication", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Views.RemoveApplication(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without portfolio key", func() {
				resp, err := client.Views.RemoveApplication(&sonar.ViewsRemoveApplicationOptions{
					Application: "my-app",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Portfolio"))
				Expect(resp).To(BeNil())
			})

			It("should fail without application key", func() {
				resp, err := client.Views.RemoveApplication(&sonar.ViewsRemoveApplicationOptions{
					Portfolio: "my-portfolio",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Application"))
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// AddPortfolio / RemovePortfolio
	// =========================================================================
	Describe("AddPortfolio", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Views.AddPortfolio(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without portfolio", func() {
				resp, err := client.Views.AddPortfolio(&sonar.ViewsAddPortfolioOptions{
					Reference: "ref-portfolio",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Portfolio"))
				Expect(resp).To(BeNil())
			})

			It("should fail without reference", func() {
				resp, err := client.Views.AddPortfolio(&sonar.ViewsAddPortfolioOptions{
					Portfolio: "parent-portfolio",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Reference"))
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("RemovePortfolio", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Views.RemovePortfolio(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})
		})
	})
})
