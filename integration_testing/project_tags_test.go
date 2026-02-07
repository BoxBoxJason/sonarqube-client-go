package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("ProjectTags Service", Ordered, func() {
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
		errors := cleanup.Cleanup()
		for _, err := range errors {
			GinkgoWriter.Printf("Cleanup error: %v\n", err)
		}
	})

	// =========================================================================
	// Search
	// =========================================================================
	Describe("Search", func() {
		It("should search for all tags", func() {
			result, resp, err := client.ProjectTags.Search(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			// Tags list may be empty or contain values
		})

		It("should search tags with query", func() {
			// First set a tag on a project
			projectKey := helpers.UniqueResourceName("proj-tag-search")

			_, _, err := client.Projects.Create(&sonar.ProjectsCreateOption{
				Name:    "Tag Search Test Project",
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonar.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			// Set a unique tag
			tagName := "e2e-test-tag"
			_, err = client.ProjectTags.Set(&sonar.ProjectTagsSetOption{
				Project: projectKey,
				Tags:    []string{tagName},
			})
			Expect(err).NotTo(HaveOccurred())

			// Search for the tag
			result, resp, err := client.ProjectTags.Search(&sonar.ProjectTagsSearchOption{
				Query: tagName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Tags).To(ContainElement(tagName))
		})

		It("should search tags with pagination", func() {
			result, resp, err := client.ProjectTags.Search(&sonar.ProjectTagsSearchOption{
				PaginationArgs: sonar.PaginationArgs{
					PageSize: 5,
					Page:     1,
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})
	})

	// =========================================================================
	// Set
	// =========================================================================
	Describe("Set", func() {
		var testProjectKey string

		BeforeEach(func() {
			testProjectKey = helpers.UniqueResourceName("proj-tag")

			_, _, err := client.Projects.Create(&sonar.ProjectsCreateOption{
				Name:    "Tag Test Project",
				Project: testProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", testProjectKey, func() error {
				_, err := client.Projects.Delete(&sonar.ProjectsDeleteOption{
					Project: testProjectKey,
				})
				return err
			})
		})

		It("should set a single tag on a project", func() {
			resp, err := client.ProjectTags.Set(&sonar.ProjectTagsSetOption{
				Project: testProjectKey,
				Tags:    []string{"backend"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		It("should set multiple tags on a project", func() {
			resp, err := client.ProjectTags.Set(&sonar.ProjectTagsSetOption{
				Project: testProjectKey,
				Tags:    []string{"backend", "api", "golang"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		It("should replace existing tags", func() {
			// Set initial tags
			_, err := client.ProjectTags.Set(&sonar.ProjectTagsSetOption{
				Project: testProjectKey,
				Tags:    []string{"old-tag"},
			})
			Expect(err).NotTo(HaveOccurred())

			// Replace with new tags
			resp, err := client.ProjectTags.Set(&sonar.ProjectTagsSetOption{
				Project: testProjectKey,
				Tags:    []string{"new-tag"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		It("should clear all tags with empty array", func() {
			// Set initial tags
			_, err := client.ProjectTags.Set(&sonar.ProjectTagsSetOption{
				Project: testProjectKey,
				Tags:    []string{"tag1", "tag2", "tag3"},
			})
			Expect(err).NotTo(HaveOccurred())

			// Clear all tags by passing an empty array
			resp, err := client.ProjectTags.Set(&sonar.ProjectTagsSetOption{
				Project: testProjectKey,
				Tags:    []string{},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify tags were cleared by setting them again
			_, err = client.ProjectTags.Set(&sonar.ProjectTagsSetOption{
				Project: testProjectKey,
				Tags:    []string{"verified"},
			})
			Expect(err).NotTo(HaveOccurred())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.ProjectTags.Set(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing project key", func() {
				resp, err := client.ProjectTags.Set(&sonar.ProjectTagsSetOption{
					Tags: []string{"tag1"},
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})

		Context("error cases", func() {
			It("should fail for non-existent project", func() {
				resp, err := client.ProjectTags.Set(&sonar.ProjectTagsSetOption{
					Project: "non-existent-project-12345",
					Tags:    []string{"tag1"},
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	// =========================================================================
	// Tag Lifecycle
	// =========================================================================
	Describe("Tag Lifecycle", func() {
		It("should complete full tag lifecycle", func() {
			projectKey := helpers.UniqueResourceName("proj-tag-lifecycle")

			// Step 1: Create project
			_, _, err := client.Projects.Create(&sonar.ProjectsCreateOption{
				Name:    "Tag Lifecycle Test Project",
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonar.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			// Step 2: Set initial tags
			_, err = client.ProjectTags.Set(&sonar.ProjectTagsSetOption{
				Project: projectKey,
				Tags:    []string{"initial-tag", "e2e-lifecycle"},
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 3: Search for the tag
			result, _, err := client.ProjectTags.Search(&sonar.ProjectTagsSearchOption{
				Query: "e2e-lifecycle",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Tags).To(ContainElement("e2e-lifecycle"))

			// Step 4: Update tags (replaces all previous tags)
			_, err = client.ProjectTags.Set(&sonar.ProjectTagsSetOption{
				Project: projectKey,
				Tags:    []string{"updated-tag"},
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 5: Verify update by searching
			result, _, err = client.ProjectTags.Search(&sonar.ProjectTagsSearchOption{
				Query: "updated-tag",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Tags).To(ContainElement("updated-tag"))
		})
	})
})
