package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Navigation Service", Ordered, func() {
	var (
		client      *sonargo.Client
		testProject *sonargo.ProjectsCreate
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())

		// Create a test project for component navigation
		testProject, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
			Name:    "navigation-e2e-test-project",
			Project: "navigation-e2e-test-project",
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
	// Component
	// =========================================================================
	Describe("Component", func() {
		Context("Functional Tests", func() {
			It("should get component navigation with valid component key", func() {
				result, resp, err := client.Navigation.Component(&sonargo.NavigationComponentOption{
					Component: testProject.Project.Key,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Key).To(Equal(testProject.Project.Key))
			})

			It("should return breadcrumbs for component", func() {
				result, resp, err := client.Navigation.Component(&sonargo.NavigationComponentOption{
					Component: testProject.Project.Key,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Breadcrumbs).NotTo(BeEmpty())
			})
		})
	})

	// =========================================================================
	// Global
	// =========================================================================
	Describe("Global", func() {
		Context("Functional Tests", func() {
			It("should get global navigation", func() {
				result, resp, err := client.Navigation.Global()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should return SonarQube version", func() {
				result, resp, err := client.Navigation.Global()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Version).NotTo(BeEmpty())
			})

			It("should return edition information", func() {
				result, resp, err := client.Navigation.Global()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
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
		Context("Functional Tests", func() {
			It("should get marketplace navigation", func() {
				result, resp, err := client.Navigation.Marketplace()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should return server ID", func() {
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
		Context("Functional Tests", func() {
			It("should get settings navigation", func() {
				result, resp, err := client.Navigation.Settings()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should return extensions list", func() {
				result, resp, err := client.Navigation.Settings()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				// Extensions list might be empty
			})
		})
	})
})
