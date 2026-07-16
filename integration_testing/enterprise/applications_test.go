package enterprise_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

var _ = Describe("Applications Service", Ordered, func() {
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

		projectKey = helpers.UniqueResourceName("app-project")
		_, resp, err := client.Projects.Create(context.Background(), &sonar.ProjectsCreateOptions{
			Name:    projectKey,
			Project: projectKey,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		cleanup.RegisterCleanup("project", projectKey, func() error {
			_, err := client.Projects.Delete(context.Background(), &sonar.ProjectsDeleteOptions{Project: projectKey})
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
	// Full lifecycle: Create -> Show -> AddProject -> SearchProjects ->
	// SetTags -> Update -> RemoveProject -> Delete
	// =========================================================================
	Describe("Lifecycle", func() {
		It("should create a new application", func() {
			appKey := helpers.UniqueResourceName("application")

			result, resp, err := client.Applications.Create(context.Background(), &sonar.ApplicationsCreateOptions{
				Name:        appKey,
				Key:         appKey,
				Description: "created by e2e enterprise suite",
				Visibility:  sonar.ProjectVisibilityPrivate,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Application.Key).To(Equal(appKey))
			Expect(result.Application.Name).To(Equal(appKey))
			Expect(result.Application.Visibility).To(Equal(sonar.ProjectVisibilityPrivate))

			cleanup.RegisterCleanup("application", appKey, func() error {
				_, err := client.Applications.Delete(context.Background(), &sonar.ApplicationsDeleteOptions{Application: appKey})
				return err
			})

			By("showing the created application")
			showResult, showResp, err := client.Applications.Show(context.Background(), &sonar.ApplicationsShowOptions{
				Application: appKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(showResp.StatusCode).To(Equal(http.StatusOK))
			Expect(showResult).NotTo(BeNil())
			Expect(showResult.Application.Key).To(Equal(appKey))
			Expect(showResult.Application.Projects).To(BeEmpty())

			By("adding the real project to the application")
			addResp, err := client.Applications.AddProject(context.Background(), &sonar.ApplicationsAddProjectOptions{
				Application: appKey,
				Project:     projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(addResp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))

			By("finding the project via SearchProjects")
			searchResult, searchResp, err := client.Applications.SearchProjects(context.Background(), &sonar.ApplicationsSearchProjectsOptions{
				Application: appKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(searchResp.StatusCode).To(Equal(http.StatusOK))
			Expect(searchResult).NotTo(BeNil())

			var found bool
			for _, p := range searchResult.Projects {
				if p.Key == projectKey {
					found = true

					break
				}
			}
			Expect(found).To(BeTrue(), "expected project %q to be listed in application %q", projectKey, appKey)

			By("confirming the project now appears on Show")
			showAfterAdd, _, err := client.Applications.Show(context.Background(), &sonar.ApplicationsShowOptions{
				Application: appKey,
			})
			Expect(err).NotTo(HaveOccurred())

			var projectListed bool
			for _, p := range showAfterAdd.Application.Projects {
				if p.Key == projectKey {
					projectListed = true

					break
				}
			}
			Expect(projectListed).To(BeTrue())

			By("setting tags on the application")
			tagsResp, err := client.Applications.SetTags(context.Background(), &sonar.ApplicationsSetTagsOptions{
				Application: appKey,
				Tags:        []string{"e2e", "enterprise"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(tagsResp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))

			By("verifying tags were applied")
			showAfterTags, _, err := client.Applications.Show(context.Background(), &sonar.ApplicationsShowOptions{
				Application: appKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(showAfterTags.Application.Tags).To(ConsistOf("e2e", "enterprise"))

			By("updating the application name and description")
			newName := appKey + "-updated"
			updateResp, err := client.Applications.Update(context.Background(), &sonar.ApplicationsUpdateOptions{
				Application: appKey,
				Name:        newName,
				Description: "updated by e2e enterprise suite",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(updateResp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))

			showAfterUpdate, _, err := client.Applications.Show(context.Background(), &sonar.ApplicationsShowOptions{
				Application: appKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(showAfterUpdate.Application.Name).To(Equal(newName))
			Expect(showAfterUpdate.Application.Description).To(Equal("updated by e2e enterprise suite"))

			By("removing the project from the application")
			removeResp, err := client.Applications.RemoveProject(context.Background(), &sonar.ApplicationsRemoveProjectOptions{
				Application: appKey,
				Project:     projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(removeResp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))

			showAfterRemove, _, err := client.Applications.Show(context.Background(), &sonar.ApplicationsShowOptions{
				Application: appKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(showAfterRemove.Application.Projects).To(BeEmpty())

			By("deleting the application")
			deleteResp, err := client.Applications.Delete(context.Background(), &sonar.ApplicationsDeleteOptions{
				Application: appKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(deleteResp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))

			_, showAfterDelete, err := client.Applications.Show(context.Background(), &sonar.ApplicationsShowOptions{
				Application: appKey,
			})
			Expect(err).To(HaveOccurred())
			Expect(showAfterDelete.StatusCode).To(Equal(http.StatusNotFound))
		})
	})

	// =========================================================================
	// Show / Delete: not-found behavior must be a real 404, never a
	// license-gate error, once the suite is running against Enterprise.
	// =========================================================================
	Describe("Show", func() {
		Context("Functional Tests", func() {
			It("should return 404 for a nonexistent application", func() {
				result, resp, err := client.Applications.Show(context.Background(), &sonar.ApplicationsShowOptions{
					Application: "nonexistent-application-xyz",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("Delete", func() {
		Context("Functional Tests", func() {
			It("should return 404 for a nonexistent application", func() {
				resp, err := client.Applications.Delete(context.Background(), &sonar.ApplicationsDeleteOptions{
					Application: "nonexistent-application-xyz",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})
})
