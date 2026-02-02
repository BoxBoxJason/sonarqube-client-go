package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Components Service", Ordered, func() {
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
	// Search
	// =========================================================================
	Describe("Search", func() {
		It("should search for projects", func() {
			result, resp, err := client.Components.Search(&sonargo.ComponentsSearchOption{
				Qualifiers: []string{"TRK"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			// There may be no projects yet
		})

		It("should search projects with pagination", func() {
			result, resp, err := client.Components.Search(&sonargo.ComponentsSearchOption{
				Qualifiers: []string{"TRK"},
				PaginationArgs: sonargo.PaginationArgs{
					PageSize: 5,
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		It("should search projects by query", func() {
			// First create a project to search for
			projectKey := helpers.UniqueResourceName("proj")
			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    projectKey,
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			// Search for it
			result, resp, err := client.Components.Search(&sonargo.ComponentsSearchOption{
				Qualifiers: []string{"TRK"},
				Query:      projectKey[:10], // First 10 chars
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.Components.Search(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing qualifiers", func() {
				_, resp, err := client.Components.Search(&sonargo.ComponentsSearchOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// SearchProjects
	// =========================================================================
	Describe("SearchProjects", func() {
		It("should search for projects", func() {
			result, resp, err := client.Components.SearchProjects(&sonargo.ComponentsSearchProjectsOption{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		It("should search projects with pagination", func() {
			result, resp, err := client.Components.SearchProjects(&sonargo.ComponentsSearchProjectsOption{
				PaginationArgs: sonargo.PaginationArgs{
					PageSize: 5,
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(len(result.Components)).To(BeNumerically("<=", 5))
		})

		It("should search projects with filter by name", func() {
			// First create a project to search for
			projectKey := helpers.UniqueResourceName("proj")
			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    projectKey,
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			// Search with filter
			result, resp, err := client.Components.SearchProjects(&sonargo.ComponentsSearchProjectsOption{
				Filter: "query = \"" + projectKey[:10] + "\"",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		It("should search projects with facets", func() {
			result, resp, err := client.Components.SearchProjects(&sonargo.ComponentsSearchProjectsOption{
				Facets: []string{"languages", "tags"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result.Facets).NotTo(BeEmpty())
		})

		It("should search projects sorted by name", func() {
			result, resp, err := client.Components.SearchProjects(&sonargo.ComponentsSearchProjectsOption{
				Sort:      "name",
				Ascending: true,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.Components.SearchProjects(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Show
	// =========================================================================
	Describe("Show", func() {
		var projectKey string

		BeforeAll(func() {
			// Create a project for Show tests
			projectKey = helpers.UniqueResourceName("proj")
			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    projectKey,
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

		It("should show component details", func() {
			result, resp, err := client.Components.Show(&sonargo.ComponentsShowOption{
				Component: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Component.Key).To(Equal(projectKey))
			Expect(result.Component.Qualifier).To(Equal("TRK"))
		})

		It("should return ancestors for project", func() {
			result, resp, err := client.Components.Show(&sonargo.ComponentsShowOption{
				Component: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			// Projects have no ancestors
			Expect(result.Ancestors).To(BeEmpty())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.Components.Show(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing component", func() {
				_, resp, err := client.Components.Show(&sonargo.ComponentsShowOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with non-existent component", func() {
				_, resp, err := client.Components.Show(&sonargo.ComponentsShowOption{
					Component: "non-existent-component-key",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	// =========================================================================
	// Tree
	// =========================================================================
	Describe("Tree", func() {
		var projectKey string

		BeforeAll(func() {
			// Create a project for Tree tests
			projectKey = helpers.UniqueResourceName("proj")
			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    projectKey,
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

		It("should get component tree", func() {
			result, resp, err := client.Components.Tree(&sonargo.ComponentsTreeOption{
				Component: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.BaseComponent.Key).To(Equal(projectKey))
		})

		It("should get tree with all strategy", func() {
			result, resp, err := client.Components.Tree(&sonargo.ComponentsTreeOption{
				Component: projectKey,
				Strategy:  "all",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		It("should get tree with children strategy", func() {
			result, resp, err := client.Components.Tree(&sonargo.ComponentsTreeOption{
				Component: projectKey,
				Strategy:  "children",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		It("should get tree with leaves strategy", func() {
			result, resp, err := client.Components.Tree(&sonargo.ComponentsTreeOption{
				Component: projectKey,
				Strategy:  "leaves",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		It("should get tree with pagination", func() {
			result, resp, err := client.Components.Tree(&sonargo.ComponentsTreeOption{
				Component: projectKey,
				PaginationArgs: sonargo.PaginationArgs{
					PageSize: 10,
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		It("should get tree with qualifiers filter", func() {
			result, resp, err := client.Components.Tree(&sonargo.ComponentsTreeOption{
				Component:  projectKey,
				Qualifiers: []string{"FIL", "DIR"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		It("should get tree sorted by name", func() {
			result, resp, err := client.Components.Tree(&sonargo.ComponentsTreeOption{
				Component: projectKey,
				Sort:      []string{"name"},
				Ascending: true,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.Components.Tree(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing component", func() {
				_, resp, err := client.Components.Tree(&sonargo.ComponentsTreeOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid strategy", func() {
				_, resp, err := client.Components.Tree(&sonargo.ComponentsTreeOption{
					Component: projectKey,
					Strategy:  "invalid-strategy",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Suggestions
	// =========================================================================
	Describe("Suggestions", func() {
		It("should get suggestions without search query", func() {
			result, resp, err := client.Components.Suggestions(&sonargo.ComponentsSuggestionsOption{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		It("should get suggestions with search query", func() {
			// First create a project
			projectKey := helpers.UniqueResourceName("proj")
			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    projectKey,
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			// Search for suggestions
			result, resp, err := client.Components.Suggestions(&sonargo.ComponentsSuggestionsOption{
				Search: projectKey[:8], // Min 2 chars
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		It("should get more suggestions for TRK qualifier", func() {
			result, resp, err := client.Components.Suggestions(&sonargo.ComponentsSuggestionsOption{
				More: "TRK",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.Components.Suggestions(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid more qualifier", func() {
				_, resp, err := client.Components.Suggestions(&sonargo.ComponentsSuggestionsOption{
					More: "INVALID",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// App
	// =========================================================================
	Describe("App", func() {
		var projectKey string

		BeforeAll(func() {
			// Create a project for App tests
			projectKey = helpers.UniqueResourceName("proj")
			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    projectKey,
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

		It("should get app data for component", func() {
			result, resp, err := client.Components.App(&sonargo.ComponentsAppOption{
				Component: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Key).To(Equal(projectKey))
		})

		It("should return component info", func() {
			result, resp, err := client.Components.App(&sonargo.ComponentsAppOption{
				Component: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result.Name).NotTo(BeEmpty())
			Expect(result.Q).To(Equal("TRK"))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.Components.App(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing component", func() {
				_, resp, err := client.Components.App(&sonargo.ComponentsAppOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with both branch and pullRequest", func() {
				_, resp, err := client.Components.App(&sonargo.ComponentsAppOption{
					Component:   projectKey,
					Branch:      "main",
					PullRequest: "123",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})
})
