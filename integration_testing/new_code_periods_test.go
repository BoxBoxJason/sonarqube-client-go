package integration_testing_test

import (
"net/http"

. "github.com/onsi/ginkgo/v2"
. "github.com/onsi/gomega"

sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("NewCodePeriods Service", Ordered, func() {
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
	// Show
	// =========================================================================
	Describe("Show", func() {
		It("should show global new code period definition", func() {
			result, resp, err := client.NewCodePeriods.Show(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Type).NotTo(BeEmpty())
		})

		It("should show project-level new code period definition", func() {
			projectKey := helpers.UniqueResourceName("ncp-show-proj")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    "NCP Show Test Project",
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			result, resp, err := client.NewCodePeriods.Show(&sonargo.NewCodePeriodsShowOption{
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.ProjectKey).To(Equal(projectKey))
		})

		It("should show branch-level new code period definition", func() {
			projectKey := helpers.UniqueResourceName("ncp-show-branch")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "NCP Branch Show Test Project",
				Project:    projectKey,
				MainBranch: "main",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			result, resp, err := client.NewCodePeriods.Show(&sonargo.NewCodePeriodsShowOption{
				Project: projectKey,
				Branch:  "main",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.BranchKey).To(Equal("main"))
		})
	})

	// =========================================================================
	// List
	// =========================================================================
	Describe("List", func() {
		It("should list new code periods for project", func() {
			projectKey := helpers.UniqueResourceName("ncp-list-proj")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "NCP List Test Project",
				Project:    projectKey,
				MainBranch: "main",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			result, resp, err := client.NewCodePeriods.List(&sonargo.NewCodePeriodsListOption{
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.NewCodePeriods).NotTo(BeNil())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.NewCodePeriods.List(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing project key", func() {
				result, resp, err := client.NewCodePeriods.List(&sonargo.NewCodePeriodsListOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Context("error cases", func() {
			It("should fail for non-existent project", func() {
				result, resp, err := client.NewCodePeriods.List(&sonargo.NewCodePeriodsListOption{
					Project: "non-existent-project-12345",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Set
	// =========================================================================
	Describe("Set", func() {
		It("should set project-level new code period with PREVIOUS_VERSION", func() {
			projectKey := helpers.UniqueResourceName("ncp-prevver")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "NCP Set PreviousVersion Test",
				Project:    projectKey,
				MainBranch: "main",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			resp, err := client.NewCodePeriods.Set(&sonargo.NewCodePeriodsSetOption{
				Project: projectKey,
				Type:    "PREVIOUS_VERSION",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("should set project-level new code period with NUMBER_OF_DAYS", func() {
			projectKey := helpers.UniqueResourceName("ncp-numdays")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "NCP Set Days Test",
				Project:    projectKey,
				MainBranch: "main",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			resp, err := client.NewCodePeriods.Set(&sonargo.NewCodePeriodsSetOption{
				Project: projectKey,
				Type:    "NUMBER_OF_DAYS",
				Value:   "30",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("should set project-level new code period with REFERENCE_BRANCH", func() {
			projectKey := helpers.UniqueResourceName("ncp-refbranch")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "NCP Set RefBranch Test",
				Project:    projectKey,
				MainBranch: "main",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			resp, err := client.NewCodePeriods.Set(&sonargo.NewCodePeriodsSetOption{
				Project: projectKey,
				Type:    "REFERENCE_BRANCH",
				Value:   "main",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("should set branch-level new code period", func() {
			projectKey := helpers.UniqueResourceName("ncp-branchlvl")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "NCP Set Branch Test",
				Project:    projectKey,
				MainBranch: "main",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			resp, err := client.NewCodePeriods.Set(&sonargo.NewCodePeriodsSetOption{
				Project: projectKey,
				Branch:  "main",
				Type:    "NUMBER_OF_DAYS",
				Value:   "15",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.NewCodePeriods.Set(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing type", func() {
				resp, err := client.NewCodePeriods.Set(&sonargo.NewCodePeriodsSetOption{
					Project: "some-project",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid type", func() {
				resp, err := client.NewCodePeriods.Set(&sonargo.NewCodePeriodsSetOption{
					Project: "some-project",
					Type:    "INVALID_TYPE",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with NUMBER_OF_DAYS and invalid value", func() {
				resp, err := client.NewCodePeriods.Set(&sonargo.NewCodePeriodsSetOption{
					Project: "some-project",
					Type:    "NUMBER_OF_DAYS",
					Value:   "invalid",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with NUMBER_OF_DAYS exceeding max", func() {
				resp, err := client.NewCodePeriods.Set(&sonargo.NewCodePeriodsSetOption{
					Project: "some-project",
					Type:    "NUMBER_OF_DAYS",
					Value:   "100",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with REFERENCE_BRANCH missing project", func() {
				resp, err := client.NewCodePeriods.Set(&sonargo.NewCodePeriodsSetOption{
					Type:  "REFERENCE_BRANCH",
					Value: "main",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with SPECIFIC_ANALYSIS missing branch", func() {
				resp, err := client.NewCodePeriods.Set(&sonargo.NewCodePeriodsSetOption{
					Project: "some-project",
					Type:    "SPECIFIC_ANALYSIS",
					Value:   "some-analysis-id",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Unset
	// =========================================================================
	Describe("Unset", func() {
		It("should unset project-level new code period", func() {
			projectKey := helpers.UniqueResourceName("ncp-unsetproj")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    "NCP Unset Test Project",
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			// First, set a new code period
			_, err = client.NewCodePeriods.Set(&sonargo.NewCodePeriodsSetOption{
				Project: projectKey,
				Type:    "PREVIOUS_VERSION",
			})
			Expect(err).NotTo(HaveOccurred())

			// Unset it
			resp, err := client.NewCodePeriods.Unset(&sonargo.NewCodePeriodsUnsetOption{
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("should unset branch-level new code period", func() {
			projectKey := helpers.UniqueResourceName("ncp-unsetbranch")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "NCP Unset Branch Test Project",
				Project:    projectKey,
				MainBranch: "main",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			// First, set a branch-level new code period
			_, err = client.NewCodePeriods.Set(&sonargo.NewCodePeriodsSetOption{
				Project: projectKey,
				Branch:  "main",
				Type:    "NUMBER_OF_DAYS",
				Value:   "45",
			})
			Expect(err).NotTo(HaveOccurred())

			// Unset it
			resp, err := client.NewCodePeriods.Unset(&sonargo.NewCodePeriodsUnsetOption{
				Project: projectKey,
				Branch:  "main",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		It("should succeed when unsetting with nil options", func() {
			resp, err := client.NewCodePeriods.Unset(nil)
			// Should get either success or a valid error
			if err != nil {
				Expect(resp).NotTo(BeNil())
			} else {
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			}
		})
	})

	// =========================================================================
	// NewCodePeriods Lifecycle
	// =========================================================================
	Describe("NewCodePeriods Lifecycle", func() {
		It("should complete full new code period lifecycle", func() {
			projectKey := helpers.UniqueResourceName("ncp-lifecycle")

			// Step 1: Create project
			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "NCP Lifecycle Test Project",
				Project:    projectKey,
				MainBranch: "main",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			// Step 2: Show project-level (inherits from global initially)
			result, _, err := client.NewCodePeriods.Show(&sonargo.NewCodePeriodsShowOption{
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())

			// Step 3: Set project-level new code period
			_, err = client.NewCodePeriods.Set(&sonargo.NewCodePeriodsSetOption{
				Project: projectKey,
				Type:    "NUMBER_OF_DAYS",
				Value:   "30",
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 4: Set branch-level new code period
			_, err = client.NewCodePeriods.Set(&sonargo.NewCodePeriodsSetOption{
				Project: projectKey,
				Branch:  "main",
				Type:    "PREVIOUS_VERSION",
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 5: List all new code periods for project
			listResult, _, err := client.NewCodePeriods.List(&sonargo.NewCodePeriodsListOption{
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(listResult.NewCodePeriods).NotTo(BeNil())

			// Step 6: Show branch-level
			result, _, err = client.NewCodePeriods.Show(&sonargo.NewCodePeriodsShowOption{
				Project: projectKey,
				Branch:  "main",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.BranchKey).To(Equal("main"))

			// Step 7: Unset branch-level
			_, err = client.NewCodePeriods.Unset(&sonargo.NewCodePeriodsUnsetOption{
				Project: projectKey,
				Branch:  "main",
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 8: Unset project-level
			_, err = client.NewCodePeriods.Unset(&sonargo.NewCodePeriodsUnsetOption{
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
