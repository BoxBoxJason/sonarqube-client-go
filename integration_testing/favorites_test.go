package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Favorites Service", Ordered, func() {
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

		// Create a test project for favorites-related operations
		projectKey = helpers.UniqueResourceName("fav")
		_, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
			Name:    "Favorites Test Project",
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
	// Search
	// =========================================================================
	Describe("Search", func() {
		Context("Valid Requests", func() {
			It("should search favorites with nil options", func() {
				result, resp, err := client.Favorites.Search(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should search favorites with pagination", func() {
				result, resp, err := client.Favorites.Search(&sonargo.FavoritesSearchOption{
					PaginationArgs: sonargo.PaginationArgs{
						PageSize: 10,
						Page:     1,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})
	})

	// =========================================================================
	// Add
	// =========================================================================
	Describe("Add", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Favorites.Add(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required component", func() {
				resp, err := client.Favorites.Add(&sonargo.FavoritesAddOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Component"))
				Expect(resp).To(BeNil())
			})
		})

		Context("Valid Requests", func() {
			It("should add a project as favorite", func() {
				resp, err := client.Favorites.Add(&sonargo.FavoritesAddOption{
					Component: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// Clean up: remove the favorite
				_, _ = client.Favorites.Remove(&sonargo.FavoritesRemoveOption{
					Component: projectKey,
				})
			})
		})

		Context("Non-Existent Component", func() {
			It("should fail for non-existent component", func() {
				resp, err := client.Favorites.Add(&sonargo.FavoritesAddOption{
					Component: "non-existent-component",
				})
				Expect(err).To(HaveOccurred())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// Remove
	// =========================================================================
	Describe("Remove", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Favorites.Remove(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required component", func() {
				resp, err := client.Favorites.Remove(&sonargo.FavoritesRemoveOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Component"))
				Expect(resp).To(BeNil())
			})
		})

		Context("Valid Requests", func() {
			It("should remove a project from favorites", func() {
				// First add the project as favorite
				_, err := client.Favorites.Add(&sonargo.FavoritesAddOption{
					Component: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())

				// Now remove it
				resp, err := client.Favorites.Remove(&sonargo.FavoritesRemoveOption{
					Component: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})
		})

		Context("Non-Existent Component", func() {
			It("should fail for non-existent component", func() {
				resp, err := client.Favorites.Remove(&sonargo.FavoritesRemoveOption{
					Component: "non-existent-component",
				})
				Expect(err).To(HaveOccurred())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})

		Context("Not Favorited Component", func() {
			It("should fail when removing a component that is not favorited", func() {
				// Ensure it's not favorited first
				_, _ = client.Favorites.Remove(&sonargo.FavoritesRemoveOption{
					Component: projectKey,
				})

				// Try to remove again
				resp, err := client.Favorites.Remove(&sonargo.FavoritesRemoveOption{
					Component: projectKey,
				})
				Expect(err).To(HaveOccurred())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// Full Workflow
	// =========================================================================
	Describe("Full Workflow", func() {
		It("should add, verify, and remove a favorite", func() {
			// Add favorite
			resp, err := client.Favorites.Add(&sonargo.FavoritesAddOption{
				Component: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Search and verify it's there
			searchResult, resp, err := client.Favorites.Search(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(searchResult).NotTo(BeNil())

			found := false
			for _, fav := range searchResult.Favorites {
				if fav.Key == projectKey {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue(), "Project should be in favorites")

			// Remove favorite
			resp, err = client.Favorites.Remove(&sonargo.FavoritesRemoveOption{
				Component: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify it's gone
			searchResult, resp, err = client.Favorites.Search(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			found = false
			for _, fav := range searchResult.Favorites {
				if fav.Key == projectKey {
					found = true
					break
				}
			}
			Expect(found).To(BeFalse(), "Project should not be in favorites")
		})
	})
})
