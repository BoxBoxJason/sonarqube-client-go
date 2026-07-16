package enterprise_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

var _ = Describe("Views Service", Ordered, func() {
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

		projectKey = helpers.UniqueResourceName("view-project")
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
	// Full lifecycle: Create -> Show -> AddProject -> Projects -> Search ->
	// Update -> Create (sub) -> AddPortfolio -> Show -> RemovePortfolio ->
	// Delete (sub) -> RemoveProject -> Projects -> Delete -> Show (404)
	// =========================================================================
	Describe("Lifecycle", func() {
		It("should run a full portfolio lifecycle", func() {
			portfolioKey := helpers.UniqueResourceName("portfolio")

			By("creating a new root portfolio")
			resp, err := client.Views.Create(context.Background(), &sonar.ViewsCreateOptions{
				Name: portfolioKey,
				Key:  portfolioKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))
			cleanup.RegisterCleanup("portfolio", portfolioKey, func() error {
				_, delErr := client.Views.Delete(context.Background(), &sonar.ViewsDeleteOptions{Key: portfolioKey})
				return helpers.IgnoreNotFoundError(delErr)
			})

			By("showing the freshly created portfolio")
			show, _, err := client.Views.Show(context.Background(), &sonar.ViewsShowOptions{
				Key: portfolioKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(show.Key).To(Equal(portfolioKey))
			Expect(show.Name).To(Equal(portfolioKey))
			Expect(show.SubViews).To(BeEmpty())

			By("adding the real project to the portfolio")
			addProjectResp, err := client.Views.AddProject(context.Background(), &sonar.ViewsAddProjectOptions{
				Key:     portfolioKey,
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(addProjectResp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))

			By("confirming the project is now listed in the portfolio")
			projects, _, err := client.Views.Projects(context.Background(), &sonar.ViewsProjectsOptions{
				Key: portfolioKey,
			})
			Expect(err).NotTo(HaveOccurred())

			var projectListed bool
			for _, p := range projects.Projects {
				if p.Key == projectKey {
					projectListed = true

					break
				}
			}
			Expect(projectListed).To(BeTrue(), "expected project %q to be listed in portfolio %q", projectKey, portfolioKey)

			By("searching for the portfolio by key")
			search, _, err := client.Views.Search(context.Background(), &sonar.ViewsSearchOptions{
				Query: portfolioKey,
			})
			Expect(err).NotTo(HaveOccurred())

			var portfolioFound bool
			for _, c := range search.Components {
				if c.Key == portfolioKey {
					portfolioFound = true

					break
				}
			}
			Expect(portfolioFound).To(BeTrue(), "expected portfolio %q to be listed in search results", portfolioKey)

			By("updating the portfolio name and description")
			newName := portfolioKey + "-updated"
			updateResp, err := client.Views.Update(context.Background(), &sonar.ViewsUpdateOptions{
				Key:         portfolioKey,
				Name:        newName,
				Description: "updated by e2e enterprise suite",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(updateResp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))

			By("verifying the update took effect")
			showAfterUpdate, _, err := client.Views.Show(context.Background(), &sonar.ViewsShowOptions{
				Key: portfolioKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(showAfterUpdate.Name).To(Equal(newName))
			Expect(showAfterUpdate.Description).To(Equal("updated by e2e enterprise suite"))

			By("creating a second portfolio to use as a sub-portfolio")
			subPortfolioKey := helpers.UniqueResourceName("sub-portfolio")
			subResp, err := client.Views.Create(context.Background(), &sonar.ViewsCreateOptions{
				Name: subPortfolioKey,
				Key:  subPortfolioKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(subResp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))
			cleanup.RegisterCleanup("portfolio", subPortfolioKey, func() error {
				_, delErr := client.Views.Delete(context.Background(), &sonar.ViewsDeleteOptions{Key: subPortfolioKey})
				return helpers.IgnoreNotFoundError(delErr)
			})

			By("adding the second portfolio as a sub-portfolio of the first")
			addPortfolioResp, err := client.Views.AddPortfolio(context.Background(), &sonar.ViewsAddPortfolioOptions{
				Portfolio: portfolioKey,
				Reference: subPortfolioKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(addPortfolioResp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))

			By("confirming the sub-portfolio is nested under the parent")
			showWithSub, _, err := client.Views.Show(context.Background(), &sonar.ViewsShowOptions{
				Key: portfolioKey,
			})
			Expect(err).NotTo(HaveOccurred())

			// Live-verified: a nested sub-portfolio's Key is a composite
			// "<parent>:<child>" key, not its own key - OriginalKey carries
			// the sub-portfolio's own key.
			var subFound bool
			for _, sv := range showWithSub.SubViews {
				if sv.OriginalKey == subPortfolioKey {
					subFound = true

					break
				}
			}
			Expect(subFound).To(BeTrue(), "expected sub-portfolio %q to be nested under %q", subPortfolioKey, portfolioKey)

			By("removing the sub-portfolio from the parent")
			removePortfolioResp, err := client.Views.RemovePortfolio(context.Background(), &sonar.ViewsRemovePortfolioOptions{
				Portfolio: portfolioKey,
				Reference: subPortfolioKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(removePortfolioResp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))

			By("confirming the sub-portfolio is no longer nested under the parent")
			showAfterRemoveSub, _, err := client.Views.Show(context.Background(), &sonar.ViewsShowOptions{
				Key: portfolioKey,
			})
			Expect(err).NotTo(HaveOccurred())

			subFound = false
			for _, sv := range showAfterRemoveSub.SubViews {
				if sv.OriginalKey == subPortfolioKey {
					subFound = true

					break
				}
			}
			Expect(subFound).To(BeFalse())

			By("deleting the now-standalone sub-portfolio")
			delSubResp, err := client.Views.Delete(context.Background(), &sonar.ViewsDeleteOptions{
				Key: subPortfolioKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(delSubResp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))

			By("removing the project from the portfolio")
			removeProjectResp, err := client.Views.RemoveProject(context.Background(), &sonar.ViewsRemoveProjectOptions{
				Key:     portfolioKey,
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(removeProjectResp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))

			By("confirming the project no longer appears in the portfolio")
			projectsAfterRemove, _, err := client.Views.Projects(context.Background(), &sonar.ViewsProjectsOptions{
				Key: portfolioKey,
			})
			Expect(err).NotTo(HaveOccurred())

			projectListed = false
			for _, p := range projectsAfterRemove.Projects {
				if p.Key == projectKey {
					projectListed = true

					break
				}
			}
			Expect(projectListed).To(BeFalse())

			By("deleting the portfolio")
			delResp, err := client.Views.Delete(context.Background(), &sonar.ViewsDeleteOptions{
				Key: portfolioKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(delResp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))

			By("confirming the portfolio no longer exists")
			_, showResp, err := client.Views.Show(context.Background(), &sonar.ViewsShowOptions{
				Key: portfolioKey,
			})
			Expect(err).To(HaveOccurred())
			Expect(showResp).NotTo(BeNil())
			Expect(showResp.StatusCode).To(Equal(http.StatusNotFound))
		})
	})

	// =========================================================================
	// Show
	// =========================================================================
	Describe("Show", func() {
		Context("Functional Tests", func() {
			It("should return 404 nonexistent portfolio", func() {
				result, resp, err := client.Views.Show(context.Background(), &sonar.ViewsShowOptions{
					Key: "nonexistent-portfolio-xyz",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Delete
	// =========================================================================
	Describe("Delete", func() {
		Context("Functional Tests", func() {
			It("should return 404 nonexistent portfolio", func() {
				resp, err := client.Views.Delete(context.Background(), &sonar.ViewsDeleteOptions{
					Key: "nonexistent-portfolio-xyz",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})
})
