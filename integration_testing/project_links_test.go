package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("ProjectLinks Service", Ordered, func() {
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
		var testProjectKey string

		BeforeEach(func() {
			testProjectKey = helpers.UniqueResourceName("proj-link")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    "Link Test Project",
				Project: testProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", testProjectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: testProjectKey,
				})
				return err
			})
		})

		It("should create a project link", func() {
			result, resp, err := client.ProjectLinks.Create(&sonargo.ProjectLinksCreateOption{
				ProjectKey: testProjectKey,
				Name:       "Homepage",
				URL:        "https://example.com",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Link.Name).To(Equal("Homepage"))
			Expect(result.Link.URL).To(Equal("https://example.com"))
			Expect(result.Link.ID).NotTo(BeEmpty())

			linkID := result.Link.ID
			cleanup.RegisterCleanup("project-link", linkID, func() error {
				_, err := client.ProjectLinks.Delete(&sonargo.ProjectLinksDeleteOption{
					ID: linkID,
				})
				return err
			})
		})

		It("should create multiple links for same project", func() {
			result1, _, err := client.ProjectLinks.Create(&sonargo.ProjectLinksCreateOption{
				ProjectKey: testProjectKey,
				Name:       "Homepage",
				URL:        "https://example.com",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result1.Link.ID).NotTo(BeEmpty())

			cleanup.RegisterCleanup("project-link", result1.Link.ID, func() error {
				_, err := client.ProjectLinks.Delete(&sonargo.ProjectLinksDeleteOption{
					ID: result1.Link.ID,
				})
				return err
			})

			result2, _, err := client.ProjectLinks.Create(&sonargo.ProjectLinksCreateOption{
				ProjectKey: testProjectKey,
				Name:       "CI",
				URL:        "https://ci.example.com",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result2.Link.ID).NotTo(BeEmpty())
			Expect(result2.Link.ID).NotTo(Equal(result1.Link.ID))

			cleanup.RegisterCleanup("project-link", result2.Link.ID, func() error {
				_, err := client.ProjectLinks.Delete(&sonargo.ProjectLinksDeleteOption{
					ID: result2.Link.ID,
				})
				return err
			})

			searchResult, _, err := client.ProjectLinks.Search(&sonargo.ProjectLinksSearchOption{
				ProjectKey: testProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(len(searchResult.Links)).To(BeNumerically(">=", 2))
		})

		It("should create links with different types", func() {
			// Create a generic custom link (won't match predefined type patterns)
			customResult, _, err := client.ProjectLinks.Create(&sonargo.ProjectLinksCreateOption{
				ProjectKey: testProjectKey,
				Name:       "Documentation",
				URL:        "https://docs.example.com",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(customResult.Link.Type).NotTo(BeEmpty())
			Expect(customResult.Link.Name).To(Equal("Documentation"))

			cleanup.RegisterCleanup("project-link", customResult.Link.ID, func() error {
				_, err := client.ProjectLinks.Delete(&sonargo.ProjectLinksDeleteOption{
					ID: customResult.Link.ID,
				})
				return err
			})

			// Verify the link type is set (SonarQube infers types based on name/URL patterns)
			// Types can be: homepage, issue, scm, ci, custom, etc.
			// Note: Links with certain inferred types (homepage, issue, scm, ci) may be
			// provided by SonarQube and cannot be deleted
			searchResult, _, err := client.ProjectLinks.Search(&sonargo.ProjectLinksSearchOption{
				ProjectKey: testProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			foundLink := false
			for _, link := range searchResult.Links {
				if link.ID == customResult.Link.ID {
					foundLink = true
					Expect(link.Type).NotTo(BeEmpty(), "Link type should be set")
					break
				}
			}
			Expect(foundLink).To(BeTrue(), "Created link should be found in search results")
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.ProjectLinks.Create(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing name", func() {
				_, resp, err := client.ProjectLinks.Create(&sonargo.ProjectLinksCreateOption{
					ProjectKey: testProjectKey,
					URL:        "https://example.com",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing URL", func() {
				_, resp, err := client.ProjectLinks.Create(&sonargo.ProjectLinksCreateOption{
					ProjectKey: testProjectKey,
					Name:       "Homepage",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing project identifier", func() {
				_, resp, err := client.ProjectLinks.Create(&sonargo.ProjectLinksCreateOption{
					Name: "Homepage",
					URL:  "https://example.com",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})

		Context("error cases", func() {
			It("should fail for non-existent project", func() {
				_, resp, err := client.ProjectLinks.Create(&sonargo.ProjectLinksCreateOption{
					ProjectKey: "non-existent-project-12345",
					Name:       "Homepage",
					URL:        "https://example.com",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})

			// Note: SonarQube is very lenient with URL validation and accepts most formats.
			// Testing with empty URL is covered in parameter validation.
			// Invalid URL format testing is skipped as SonarQube accepts various URL formats.
		})
	})

	// =========================================================================
	// Search
	// =========================================================================
	Describe("Search", func() {
		var testProjectKey string
		var testLinkID string

		BeforeEach(func() {
			testProjectKey = helpers.UniqueResourceName("proj-search-link")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    "Search Link Test Project",
				Project: testProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", testProjectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: testProjectKey,
				})
				return err
			})

			result, _, err := client.ProjectLinks.Create(&sonargo.ProjectLinksCreateOption{
				ProjectKey: testProjectKey,
				Name:       "Test Link",
				URL:        "https://test.example.com",
			})
			Expect(err).NotTo(HaveOccurred())
			testLinkID = result.Link.ID

			cleanup.RegisterCleanup("project-link", testLinkID, func() error {
				_, err := client.ProjectLinks.Delete(&sonargo.ProjectLinksDeleteOption{
					ID: testLinkID,
				})
				return err
			})
		})

		It("should search links by project key", func() {
			result, resp, err := client.ProjectLinks.Search(&sonargo.ProjectLinksSearchOption{
				ProjectKey: testProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Links).NotTo(BeEmpty())

			found := false
			for _, link := range result.Links {
				if link.ID == testLinkID {
					found = true
					Expect(link.Name).To(Equal("Test Link"))
					Expect(link.URL).To(Equal("https://test.example.com"))
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		It("should return empty list for project without links", func() {
			newProjectKey := helpers.UniqueResourceName("proj-no-links")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    "No Links Project",
				Project: newProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", newProjectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: newProjectKey,
				})
				return err
			})

			result, resp, err := client.ProjectLinks.Search(&sonargo.ProjectLinksSearchOption{
				ProjectKey: newProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Links).To(BeEmpty())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.ProjectLinks.Search(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing project identifier", func() {
				_, resp, err := client.ProjectLinks.Search(&sonargo.ProjectLinksSearchOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})

		Context("error cases", func() {
			It("should fail for non-existent project", func() {
				_, resp, err := client.ProjectLinks.Search(&sonargo.ProjectLinksSearchOption{
					ProjectKey: "non-existent-project-12345",
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
		It("should delete a project link", func() {
			projectKey := helpers.UniqueResourceName("proj-del-link")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    "Delete Link Test Project",
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			result, _, err := client.ProjectLinks.Create(&sonargo.ProjectLinksCreateOption{
				ProjectKey: projectKey,
				Name:       "To Delete",
				URL:        "https://delete.example.com",
			})
			Expect(err).NotTo(HaveOccurred())
			linkID := result.Link.ID

			// Register cleanup in case deletion fails
			cleanup.RegisterCleanup("project-link", linkID, func() error {
				_, err := client.ProjectLinks.Delete(&sonargo.ProjectLinksDeleteOption{
					ID: linkID,
				})
				return err
			})

			resp, err := client.ProjectLinks.Delete(&sonargo.ProjectLinksDeleteOption{
				ID: linkID,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			searchResult, _, err := client.ProjectLinks.Search(&sonargo.ProjectLinksSearchOption{
				ProjectKey: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			for _, link := range searchResult.Links {
				Expect(link.ID).NotTo(Equal(linkID))
			}
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.ProjectLinks.Delete(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing ID", func() {
				resp, err := client.ProjectLinks.Delete(&sonargo.ProjectLinksDeleteOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})

		Context("error cases", func() {
			It("should fail for non-existent link", func() {
				resp, err := client.ProjectLinks.Delete(&sonargo.ProjectLinksDeleteOption{
					ID: "99999999",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	// =========================================================================
	// Link Lifecycle
	// =========================================================================
	Describe("Link Lifecycle", func() {
		It("should complete full link lifecycle", func() {
			projectKey := helpers.UniqueResourceName("proj-link-lifecycle")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    "Link Lifecycle Test Project",
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			createResult, _, err := client.ProjectLinks.Create(&sonargo.ProjectLinksCreateOption{
				ProjectKey: projectKey,
				Name:       "Lifecycle Link",
				URL:        "https://lifecycle.example.com",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(createResult.Link.ID).NotTo(BeEmpty())
			linkID := createResult.Link.ID

			// Register cleanup in case test fails before deletion
			cleanup.RegisterCleanup("project-link", linkID, func() error {
				_, err := client.ProjectLinks.Delete(&sonargo.ProjectLinksDeleteOption{
					ID: linkID,
				})
				return err
			})

			searchResult, _, err := client.ProjectLinks.Search(&sonargo.ProjectLinksSearchOption{
				ProjectKey: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, link := range searchResult.Links {
				if link.ID == linkID {
					found = true
					Expect(link.Name).To(Equal("Lifecycle Link"))
					Expect(link.URL).To(Equal("https://lifecycle.example.com"))
					break
				}
			}
			Expect(found).To(BeTrue())

			_, err = client.ProjectLinks.Delete(&sonargo.ProjectLinksDeleteOption{
				ID: linkID,
			})
			Expect(err).NotTo(HaveOccurred())

			searchResult, _, err = client.ProjectLinks.Search(&sonargo.ProjectLinksSearchOption{
				ProjectKey: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			for _, link := range searchResult.Links {
				Expect(link.ID).NotTo(Equal(linkID))
			}
		})
	})
})
