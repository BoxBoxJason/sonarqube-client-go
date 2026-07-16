package enterprise_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

var _ = Describe("Security Reports Service", Ordered, func() {
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

		projectKey = helpers.UniqueResourceName("security-project")
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
	// Show
	// =========================================================================
	Describe("Show", func() {
		Context("Functional Tests", func() {
			It("should return a real (zero-valued) report even for an unanalyzed project", func() {
				// Live-verified: this endpoint succeeds unconditionally (200,
				// full OWASP Top 10 category structure, zero-valued since no
				// analysis has run) - there is no data-availability failure
				// mode to tolerate here.
				result, resp, err := client.SecurityReports.Show(context.Background(), &sonar.SecurityReportsShowOptions{
					Project:  projectKey,
					Standard: "owaspTop10",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})
	})

	// =========================================================================
	// Download
	// =========================================================================
	Describe("Download", func() {
		Context("Functional Tests", func() {
			It("should download a real report even for an unanalyzed project", func() {
				// Live-verified: like Show, this endpoint succeeds
				// unconditionally (200, real ~17MB PDF) even without any
				// scanner analysis - there is no data-availability failure
				// mode to tolerate here. Unlike Show, Download only accepts a
				// restricted subset of standards ("sonarsourceSecurity",
				// "casa", "owaspMasvs-v2"); "owaspTop10" - valid for Show -
				// is rejected here with 400 "Standard is not supported".
				result, resp, err := client.SecurityReports.Download(context.Background(), &sonar.SecurityReportsDownloadOptions{
					Project:   projectKey,
					Standards: []string{"sonarsourceSecurity"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeEmpty())
			})
		})
	})

	// =========================================================================
	// Nonexistent project
	// =========================================================================
	Describe("Nonexistent project", func() {
		Context("Functional Tests", func() {
			It("should return a not-found style error, never a license-gate error", func() {
				// Even on a confirmed Enterprise server, a nonexistent project
				// key must be reported as "not found", not "you need a license".
				_, resp, err := client.SecurityReports.Show(context.Background(), &sonar.SecurityReportsShowOptions{
					Project:  "e2e-nonexistent-security-project",
					Standard: "owaspTop10",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).NotTo(Equal(http.StatusPaymentRequired))
				Expect(resp.StatusCode).NotTo(Equal(http.StatusForbidden))
			})
		})
	})
})
