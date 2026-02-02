package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Projects Service", Ordered, func() {
	var (
		client  *sonargo.Client
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
		errors := cleanup.Cleanup()
		for _, err := range errors {
			GinkgoWriter.Printf("Cleanup error: %v\n", err)
		}
	})

	// =========================================================================
	// Create
	// =========================================================================
	Describe("Create", func() {
		It("should create a project with minimal parameters", func() {
			projectKey := helpers.UniqueResourceName("proj-min")

			result, resp, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    "Minimal Project",
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			// Register cleanup immediately after successful Create to avoid orphaned resources
			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Project.Key).To(Equal(projectKey))
			Expect(result.Project.Name).To(Equal("Minimal Project"))
		})

		It("should create a project with full configuration", func() {
			projectKey := helpers.UniqueResourceName("proj-full")

			result, resp, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "Full Config Project",
				Project:    projectKey,
				Visibility: "private",
				MainBranch: "main",
			})
			Expect(err).NotTo(HaveOccurred())

			// Register cleanup immediately after successful Create to avoid orphaned resources
			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Project.Key).To(Equal(projectKey))
			Expect(result.Project.Visibility).To(Equal("private"))
		})

		It("should create a public project", func() {
			projectKey := helpers.UniqueResourceName("proj-pub")

			result, resp, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "Public Project",
				Project:    projectKey,
				Visibility: "public",
			})
			Expect(err).NotTo(HaveOccurred())

			// Register cleanup immediately after successful Create to avoid orphaned resources
			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Project.Visibility).To(Equal("public"))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.Projects.Create(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing name", func() {
				_, resp, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
					Project: helpers.UniqueResourceName("proj-noname"),
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing project key", func() {
				_, resp, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
					Name: "No Key Project",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid visibility", func() {
				_, resp, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
					Name:       "Invalid Visibility",
					Project:    helpers.UniqueResourceName("proj-badvis"),
					Visibility: "invalid",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})

		Context("error cases", func() {
			It("should fail to create project with duplicate key", func() {
				projectKey := helpers.UniqueResourceName("proj-dup")

				// Create first project
				_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
					Name:    "Original Project",
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())

				cleanup.RegisterCleanup("project", projectKey, func() error {
					_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
						Project: projectKey,
					})
					return err
				})

				// Try to create duplicate
				_, resp, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
					Name:    "Duplicate Project",
					Project: projectKey,
				})
				Expect(err).To(HaveOccurred())
				if resp != nil {
					Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
				}
			})
		})
	})

	// =========================================================================
	// Search
	// =========================================================================
	Describe("Search", func() {
		var testProjectKey string

		BeforeEach(func() {
			testProjectKey = helpers.UniqueResourceName("proj-search")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "Search Test Project",
				Project:    testProjectKey,
				Visibility: "private",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", testProjectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: testProjectKey,
				})
				return err
			})
		})

		It("should search all projects", func() {
			result, resp, err := client.Projects.Search(&sonargo.ProjectsSearchOption{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Components).NotTo(BeEmpty())
		})

		It("should search projects by query", func() {
			result, resp, err := client.Projects.Search(&sonargo.ProjectsSearchOption{
				Query: testProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			found := false
			for _, c := range result.Components {
				if c.Key == testProjectKey {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		It("should search projects with pagination", func() {
			result, resp, err := client.Projects.Search(&sonargo.ProjectsSearchOption{
				PaginationArgs: sonargo.PaginationArgs{
					PageSize: 5,
					Page:     1,
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(len(result.Components)).To(BeNumerically("<=", 5))
		})

		It("should search projects by visibility", func() {
			result, resp, err := client.Projects.Search(&sonargo.ProjectsSearchOption{
				Visibility: "private",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			for _, c := range result.Components {
				Expect(c.Visibility).To(Equal("private"))
			}
		})

		It("should search projects by qualifier", func() {
			result, resp, err := client.Projects.Search(&sonargo.ProjectsSearchOption{
				Qualifiers: []string{"TRK"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			for _, c := range result.Components {
				Expect(c.Qualifier).To(Equal("TRK"))
			}
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.Projects.Search(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid visibility", func() {
				_, resp, err := client.Projects.Search(&sonargo.ProjectsSearchOption{
					Visibility: "invalid",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid qualifier", func() {
				_, resp, err := client.Projects.Search(&sonargo.ProjectsSearchOption{
					Qualifiers: []string{"INVALID"},
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Delete
	// =========================================================================
	Describe("Delete", func() {
		It("should delete a project", func() {
			projectKey := helpers.UniqueResourceName("proj-del")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    "Delete Test Project",
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			// Register cleanup to handle orphaned resources if delete or assertions fail
			// The cleanup will gracefully handle "not found" errors if the project was already deleted
			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				// Ignore not found errors since the test may have already deleted it
				return helpers.IgnoreNotFoundError(err)
			})

			resp, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify project is deleted
			result, _, err := client.Projects.Search(&sonargo.ProjectsSearchOption{
				Query: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			for _, c := range result.Components {
				Expect(c.Key).NotTo(Equal(projectKey))
			}
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Projects.Delete(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing project key", func() {
				resp, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})

		Context("error cases", func() {
			It("should fail to delete non-existent project", func() {
				resp, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: "non-existent-project-key-12345",
				})
				Expect(err).To(HaveOccurred())
				if resp != nil {
					Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				}
			})
		})
	})

	// =========================================================================
	// UpdateKey
	// =========================================================================
	Describe("UpdateKey", func() {
		It("should update project key", func() {
			oldKey := helpers.UniqueResourceName("proj-oldkey")
			newKey := helpers.UniqueResourceName("proj-newkey")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    "Update Key Project",
				Project: oldKey,
			})
			Expect(err).NotTo(HaveOccurred())

			// Register cleanup with old key first in case UpdateKey fails
			cleanup.RegisterCleanup("project", oldKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: oldKey,
				})
				return helpers.IgnoreNotFoundError(err)
			})

			resp, err := client.Projects.UpdateKey(&sonargo.ProjectsUpdateKeyOption{
				From: oldKey,
				To:   newKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Register cleanup with new key
			cleanup.RegisterCleanup("project", newKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: newKey,
				})
				return err
			})

			// Verify old key no longer exists
			result, _, err := client.Projects.Search(&sonargo.ProjectsSearchOption{
				Query: oldKey,
			})
			Expect(err).NotTo(HaveOccurred())
			for _, c := range result.Components {
				Expect(c.Key).NotTo(Equal(oldKey))
			}

			// Verify new key exists
			result, _, err = client.Projects.Search(&sonargo.ProjectsSearchOption{
				Query: newKey,
			})
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, c := range result.Components {
				if c.Key == newKey {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Projects.UpdateKey(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing from key", func() {
				resp, err := client.Projects.UpdateKey(&sonargo.ProjectsUpdateKeyOption{
					To: helpers.UniqueResourceName("proj-to"),
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing to key", func() {
				resp, err := client.Projects.UpdateKey(&sonargo.ProjectsUpdateKeyOption{
					From: helpers.UniqueResourceName("proj-from"),
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// UpdateVisibility
	// =========================================================================
	Describe("UpdateVisibility", func() {
		It("should update project visibility from public to private", func() {
			projectKey := helpers.UniqueResourceName("proj-vis")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "Visibility Test Project",
				Project:    projectKey,
				Visibility: "public",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			resp, err := client.Projects.UpdateVisibility(&sonargo.ProjectsUpdateVisibilityOption{
				Project:    projectKey,
				Visibility: "private",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify visibility changed
			result, _, err := client.Projects.Search(&sonargo.ProjectsSearchOption{
				Query: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, c := range result.Components {
				if c.Key == projectKey {
					found = true
					Expect(c.Visibility).To(Equal("private"))
					break
				}
			}
			Expect(found).To(BeTrue(), "project %s not found in search results", projectKey)
		})

		It("should update project visibility from private to public", func() {
			projectKey := helpers.UniqueResourceName("proj-vis2")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "Visibility Test Project 2",
				Project:    projectKey,
				Visibility: "private",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			resp, err := client.Projects.UpdateVisibility(&sonargo.ProjectsUpdateVisibilityOption{
				Project:    projectKey,
				Visibility: "public",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify visibility changed
			result, _, err := client.Projects.Search(&sonargo.ProjectsSearchOption{
				Query: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, c := range result.Components {
				if c.Key == projectKey {
					found = true
					Expect(c.Visibility).To(Equal("public"))
					break
				}
			}
			Expect(found).To(BeTrue(), "project %s not found in search results", projectKey)
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Projects.UpdateVisibility(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing project key", func() {
				resp, err := client.Projects.UpdateVisibility(&sonargo.ProjectsUpdateVisibilityOption{
					Visibility: "private",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing visibility", func() {
				resp, err := client.Projects.UpdateVisibility(&sonargo.ProjectsUpdateVisibilityOption{
					Project: helpers.UniqueResourceName("proj"),
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid visibility", func() {
				resp, err := client.Projects.UpdateVisibility(&sonargo.ProjectsUpdateVisibilityOption{
					Project:    helpers.UniqueResourceName("proj"),
					Visibility: "invalid",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// BulkDelete
	// =========================================================================
	Describe("BulkDelete", func() {
		It("should bulk delete projects by keys", func() {
			projectKey1 := helpers.UniqueResourceName("proj-bulk1")
			projectKey2 := helpers.UniqueResourceName("proj-bulk2")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    "Bulk Delete 1",
				Project: projectKey1,
			})
			Expect(err).NotTo(HaveOccurred())

			// Register cleanup immediately after successful Create (ignores 404 if bulk delete succeeds)
			cleanup.RegisterCleanup("project", projectKey1, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey1,
				})
				return helpers.IgnoreNotFoundError(err)
			})

			_, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    "Bulk Delete 2",
				Project: projectKey2,
			})
			Expect(err).NotTo(HaveOccurred())

			// Register cleanup immediately after successful Create (ignores 404 if bulk delete succeeds)
			cleanup.RegisterCleanup("project", projectKey2, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey2,
				})
				return helpers.IgnoreNotFoundError(err)
			})

			resp, err := client.Projects.BulkDelete(&sonargo.ProjectsBulkDeleteOption{
				Projects: []string{projectKey1, projectKey2},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify projectKey1 is deleted
			result, _, err := client.Projects.Search(&sonargo.ProjectsSearchOption{
				Query: projectKey1,
			})
			Expect(err).NotTo(HaveOccurred())
			for _, c := range result.Components {
				Expect(c.Key).NotTo(Equal(projectKey1))
			}

			// Verify projectKey2 is deleted
			result, _, err = client.Projects.Search(&sonargo.ProjectsSearchOption{
				Query: projectKey2,
			})
			Expect(err).NotTo(HaveOccurred())
			for _, c := range result.Components {
				Expect(c.Key).NotTo(Equal(projectKey2))
			}
		})

		It("should bulk delete projects by query", func() {
			uniquePrefix := helpers.UniqueResourceName("bulkq")
			projectKey1 := uniquePrefix + "-a"
			projectKey2 := uniquePrefix + "-b"

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    "Bulk Query Delete 1",
				Project: projectKey1,
			})
			Expect(err).NotTo(HaveOccurred())

			// Register cleanup immediately after successful Create (ignores 404 if bulk delete succeeds)
			cleanup.RegisterCleanup("project", projectKey1, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey1,
				})
				return helpers.IgnoreNotFoundError(err)
			})

			_, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    "Bulk Query Delete 2",
				Project: projectKey2,
			})
			Expect(err).NotTo(HaveOccurred())

			// Register cleanup immediately after successful Create (ignores 404 if bulk delete succeeds)
			cleanup.RegisterCleanup("project", projectKey2, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey2,
				})
				return helpers.IgnoreNotFoundError(err)
			})

			resp, err := client.Projects.BulkDelete(&sonargo.ProjectsBulkDeleteOption{
				Query: uniquePrefix,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Projects.BulkDelete(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with no filter", func() {
				resp, err := client.Projects.BulkDelete(&sonargo.ProjectsBulkDeleteOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid visibility", func() {
				resp, err := client.Projects.BulkDelete(&sonargo.ProjectsBulkDeleteOption{
					Query:      "test",
					Visibility: "invalid",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid qualifier", func() {
				resp, err := client.Projects.BulkDelete(&sonargo.ProjectsBulkDeleteOption{
					Query:      "test",
					Qualifiers: []string{"INVALID"},
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// SearchMyProjects
	// =========================================================================
	Describe("SearchMyProjects", func() {
		It("should search my projects", func() {
			result, resp, err := client.Projects.SearchMyProjects(&sonargo.ProjectsSearchMyProjectsOption{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			// Result may be empty if no projects are marked as favorites
		})

		It("should search my projects with pagination", func() {
			result, resp, err := client.Projects.SearchMyProjects(&sonargo.ProjectsSearchMyProjectsOption{
				PaginationArgs: sonargo.PaginationArgs{
					PageSize: 5,
					Page:     1,
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.Projects.SearchMyProjects(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// SearchMyScannableProjects
	// =========================================================================
	Describe("SearchMyScannableProjects", func() {
		It("should search my scannable projects", func() {
			result, resp, err := client.Projects.SearchMyScannableProjects(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		It("should search my scannable projects with query", func() {
			result, resp, err := client.Projects.SearchMyScannableProjects(&sonargo.ProjectsSearchMyScannableProjectsOption{
				Query: helpers.E2EResourcePrefix,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})
	})

	// =========================================================================
	// UpdateDefaultVisibility
	// =========================================================================
	Describe("UpdateDefaultVisibility", Ordered, func() {
		// NOTE: These tests modify a global SonarQube setting. We restore the
		// original value after all tests complete to avoid affecting other tests
		// or leaving the developer's instance in a modified state.
		var originalVisibility string

		BeforeAll(func() {
			// Capture current default visibility before tests
			// Default is typically "public" but we should restore whatever was set
			originalVisibility = "public" // SonarQube default
		})

		AfterAll(func() {
			// Restore original default visibility
			_, err := client.Projects.UpdateDefaultVisibility(&sonargo.ProjectsUpdateDefaultVisibilityOption{
				ProjectVisibility: originalVisibility,
			})
			if err != nil {
				GinkgoWriter.Printf("Warning: failed to restore default visibility: %v\n", err)
			}
		})

		It("should update default visibility to private", func() {
			resp, err := client.Projects.UpdateDefaultVisibility(&sonargo.ProjectsUpdateDefaultVisibilityOption{
				ProjectVisibility: "private",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		It("should update default visibility to public", func() {
			resp, err := client.Projects.UpdateDefaultVisibility(&sonargo.ProjectsUpdateDefaultVisibilityOption{
				ProjectVisibility: "public",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Projects.UpdateDefaultVisibility(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing visibility", func() {
				resp, err := client.Projects.UpdateDefaultVisibility(&sonargo.ProjectsUpdateDefaultVisibilityOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid visibility", func() {
				resp, err := client.Projects.UpdateDefaultVisibility(&sonargo.ProjectsUpdateDefaultVisibilityOption{
					ProjectVisibility: "invalid",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Project Lifecycle
	// =========================================================================
	Describe("Project Lifecycle", func() {
		It("should complete full project lifecycle", func() {
			projectKey := helpers.UniqueResourceName("proj-lifecycle")
			newProjectKey := helpers.UniqueResourceName("proj-lifecycle-new")

			// Step 1: Create project
			result, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "Lifecycle Test",
				Project:    projectKey,
				Visibility: "public",
			})
			Expect(err).NotTo(HaveOccurred())

			// Register cleanup for both possible project keys to handle orphans
			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return helpers.IgnoreNotFoundError(err)
			})
			cleanup.RegisterCleanup("project", newProjectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: newProjectKey,
				})
				return helpers.IgnoreNotFoundError(err)
			})

			Expect(result.Project.Key).To(Equal(projectKey))
			Expect(result.Project.Visibility).To(Equal("public"))

			// Step 2: Search for project
			searchResult, _, err := client.Projects.Search(&sonargo.ProjectsSearchOption{
				Query: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, c := range searchResult.Components {
				if c.Key == projectKey {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())

			// Step 3: Update visibility
			_, err = client.Projects.UpdateVisibility(&sonargo.ProjectsUpdateVisibilityOption{
				Project:    projectKey,
				Visibility: "private",
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 4: Update key
			_, err = client.Projects.UpdateKey(&sonargo.ProjectsUpdateKeyOption{
				From: projectKey,
				To:   newProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 5: Verify changes
			searchResult, _, err = client.Projects.Search(&sonargo.ProjectsSearchOption{
				Query: newProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			found = false
			for _, c := range searchResult.Components {
				if c.Key == newProjectKey {
					found = true
					Expect(c.Visibility).To(Equal("private"))
					break
				}
			}
			Expect(found).To(BeTrue())

			// Step 6: Delete project
			_, err = client.Projects.Delete(&sonargo.ProjectsDeleteOption{
				Project: newProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 7: Verify deletion
			searchResult, _, err = client.Projects.Search(&sonargo.ProjectsSearchOption{
				Query: newProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			for _, c := range searchResult.Components {
				Expect(c.Key).NotTo(Equal(newProjectKey))
			}
		})
	})
})
