package integration_testing_test

import (
	"net/http"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("V2 DOP Translation Service", Ordered, func() {
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
	// FetchAllDopSettings
	// =========================================================================
	Describe("FetchAllDopSettings", func() {
		It("should return DevOps Platform settings", func() {
			result, resp, err := client.V2.DopTranslation.FetchAllDopSettings()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			// DopSettings might be empty if no ALM integrations are configured,
			// but the response structure should be valid
			Expect(result.Page.PageSize).To(BeNumerically(">=", 0))
		})
	})

	// =========================================================================
	// CreateBoundProject
	// =========================================================================
	Describe("CreateBoundProject", func() {
		Context("parameter validation", func() {
			It("should fail with nil request", func() {
				result, resp, err := client.V2.DopTranslation.CreateBoundProject(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing required fields", func() {
				result, resp, err := client.V2.DopTranslation.CreateBoundProject(&sonar.DopTranslationBoundProjectOptions{
					ProjectKey:  helpers.UniqueResourceName("v2dopproj"),
					ProjectName: "Test DOP Project",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without project key", func() {
				result, resp, err := client.V2.DopTranslation.CreateBoundProject(&sonar.DopTranslationBoundProjectOptions{
					DevOpsPlatformSettingId: "nonexistent-setting-id",
					ProjectName:             "Test DOP Project",
					RepositoryIdentifier:    "test/repo",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without repository identifier", func() {
				result, resp, err := client.V2.DopTranslation.CreateBoundProject(&sonar.DopTranslationBoundProjectOptions{
					DevOpsPlatformSettingId: "nonexistent-setting-id",
					ProjectKey:              helpers.UniqueResourceName("v2dopproj"),
					ProjectName:             "Test DOP Project",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Context("successful creation", func() {
			var (
				dopSettingID   string
				repoIdentifier string
			)

			BeforeAll(func() {
				settings, _, err := client.V2.DopTranslation.FetchAllDopSettings()
				Expect(err).NotTo(HaveOccurred())

				if len(settings.DopSettings) == 0 {
					Skip("No DOP platform settings configured, skipping successful CreateBoundProject tests")
				}

				dopSettingID = settings.DopSettings[0].Id

				repoIdentifier = os.Getenv("SONAR_DOP_REPOSITORY_IDENTIFIER")
				if repoIdentifier == "" {
					Skip("SONAR_DOP_REPOSITORY_IDENTIFIER env var not set, skipping successful CreateBoundProject tests")
				}
			})

			It("should create a bound project", func() {
				projectKey := helpers.UniqueResourceName("v2dopproj")

				result, resp, err := client.V2.DopTranslation.CreateBoundProject(&sonar.DopTranslationBoundProjectOptions{
					DevOpsPlatformSettingId: dopSettingID,
					ProjectKey:              projectKey,
					ProjectName:             "E2E V2 DOP Bound Project " + projectKey,
					RepositoryIdentifier:    repoIdentifier,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.ProjectId).NotTo(BeEmpty())
				Expect(result.NewProjectCreated).To(BeTrue())

				cleanup.RegisterCleanup("project", projectKey, func() error {
					_, err := client.Projects.Delete(&sonar.ProjectsDeleteOption{
						Project: projectKey,
					})
					return helpers.IgnoreNotFoundError(err)
				})
			})
		})
	})

	// =========================================================================
	// CreateOrUpdateBoundProject
	// =========================================================================
	Describe("CreateOrUpdateBoundProject", func() {
		Context("parameter validation", func() {
			It("should fail with nil request", func() {
				result, resp, err := client.V2.DopTranslation.CreateOrUpdateBoundProject(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing required fields", func() {
				result, resp, err := client.V2.DopTranslation.CreateOrUpdateBoundProject(&sonar.DopTranslationBoundProjectOptions{
					ProjectKey:  helpers.UniqueResourceName("v2dopproj"),
					ProjectName: "Test DOP Project",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Context("successful creation", func() {
			var (
				dopSettingID   string
				repoIdentifier string
			)

			BeforeAll(func() {
				settings, _, err := client.V2.DopTranslation.FetchAllDopSettings()
				Expect(err).NotTo(HaveOccurred())

				if len(settings.DopSettings) == 0 {
					Skip("No DOP platform settings configured, skipping successful CreateOrUpdateBoundProject tests")
				}

				dopSettingID = settings.DopSettings[0].Id

				repoIdentifier = os.Getenv("SONAR_DOP_REPOSITORY_IDENTIFIER")
				if repoIdentifier == "" {
					Skip("SONAR_DOP_REPOSITORY_IDENTIFIER env var not set, skipping successful CreateOrUpdateBoundProject tests")
				}
			})

			It("should create a bound project", func() {
				projectKey := helpers.UniqueResourceName("v2dopproj")

				result, resp, err := client.V2.DopTranslation.CreateOrUpdateBoundProject(&sonar.DopTranslationBoundProjectOptions{
					DevOpsPlatformSettingId: dopSettingID,
					ProjectKey:              projectKey,
					ProjectName:             "E2E V2 DOP CreateOrUpdate Project " + projectKey,
					RepositoryIdentifier:    repoIdentifier,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.ProjectId).NotTo(BeEmpty())
				Expect(result.NewProjectCreated).To(BeTrue())

				cleanup.RegisterCleanup("project", projectKey, func() error {
					_, err := client.Projects.Delete(&sonar.ProjectsDeleteOption{
						Project: projectKey,
					})
					return helpers.IgnoreNotFoundError(err)
				})
			})

			It("should update an existing bound project", func() {
				// Use the same project key to trigger an update rather than a new creation.
				projectKey := helpers.UniqueResourceName("v2dopupd")

				// First create the project.
				createResult, _, err := client.V2.DopTranslation.CreateOrUpdateBoundProject(&sonar.DopTranslationBoundProjectOptions{
					DevOpsPlatformSettingId: dopSettingID,
					ProjectKey:              projectKey,
					ProjectName:             "E2E V2 DOP Update Project " + projectKey,
					RepositoryIdentifier:    repoIdentifier,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(createResult.NewProjectCreated).To(BeTrue())

				cleanup.RegisterCleanup("project", projectKey, func() error {
					_, err := client.Projects.Delete(&sonar.ProjectsDeleteOption{
						Project: projectKey,
					})
					return helpers.IgnoreNotFoundError(err)
				})

				// Now update it (same project key, idempotent PUT).
				updateResult, resp, err := client.V2.DopTranslation.CreateOrUpdateBoundProject(&sonar.DopTranslationBoundProjectOptions{
					DevOpsPlatformSettingId: dopSettingID,
					ProjectKey:              projectKey,
					ProjectName:             "E2E V2 DOP Updated Project " + projectKey,
					RepositoryIdentifier:    repoIdentifier,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(updateResult).NotTo(BeNil())
				Expect(updateResult.NewProjectCreated).To(BeFalse())
			})
		})
	})
})
