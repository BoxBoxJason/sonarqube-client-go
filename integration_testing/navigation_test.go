package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("Navigation Service", Ordered, func() {
	var (
		client     *sonar.Client
		cleanup    *helpers.CleanupManager
		projectKey string
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
		cleanup = helpers.NewCleanupManager(client)

		// Create a test project for component navigation
		projectKey = helpers.UniqueResourceName("nav")
		_, _, err = client.Projects.Create(&sonar.ProjectsCreateOption{
			Name:    "Navigation Test Project",
			Project: projectKey,
		})
		Expect(err).NotTo(HaveOccurred())

		cleanup.RegisterCleanup("project", projectKey, func() error {
			_, err := client.Projects.Delete(&sonar.ProjectsDeleteOption{
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
	// Component
	// =========================================================================
	Describe("Component", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, _, err := client.Navigation.Component(nil)
				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
			})

			It("should fail with empty component key", func() {
				resp, _, err := client.Navigation.Component(&sonar.NavigationComponentOption{
					Component: "",
				})
				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
			})
		})

		Context("Valid Requests", func() {
			It("should get component navigation with valid component key and breadcrumbs", func() {
				result, resp, err := client.Navigation.Component(&sonar.NavigationComponentOption{
					Component: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Key).To(Equal(projectKey))
				Expect(result.Breadcrumbs).NotTo(BeEmpty())
			})
		})

		Context("Error Cases", func() {
			It("should fail with non-existent component key", func() {
				_, resp, err := client.Navigation.Component(&sonar.NavigationComponentOption{
					Component: "non-existent-component-key-12345",
				})
				Expect(err).To(HaveOccurred())
				if resp != nil {
					Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				}
			})
		})
	})

	// =========================================================================
	// Global
	// =========================================================================
	Describe("Global", func() {
		Context("Valid Requests", func() {
			It("should get global navigation with version and edition information", func() {
				result, resp, err := client.Navigation.Global()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Version).NotTo(BeEmpty())
				// Edition could be empty for community edition
			})

			It("should return consistent results on multiple calls", func() {
				result1, resp1, err := client.Navigation.Global()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp1.StatusCode).To(Equal(http.StatusOK))

				result2, resp2, err := client.Navigation.Global()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp2.StatusCode).To(Equal(http.StatusOK))

				Expect(result1.Version).To(Equal(result2.Version))
			})
		})
	})

	// =========================================================================
	// Marketplace
	// =========================================================================
	Describe("Marketplace", func() {
		Context("Valid Requests", func() {
			It("should get marketplace navigation", func() {
				result, resp, err := client.Navigation.Marketplace()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				// ServerID might be empty in some installations
			})
		})
	})

	// =========================================================================
	// Settings
	// =========================================================================
	Describe("Settings", func() {
		Context("Valid Requests", func() {
			It("should get settings navigation", func() {
				result, resp, err := client.Navigation.Settings()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				// Extensions list might be empty
			})
		})
	})
})
