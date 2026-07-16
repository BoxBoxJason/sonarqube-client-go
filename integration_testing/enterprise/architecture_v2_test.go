package enterprise_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

var _ = Describe("Architecture V2 Service", Ordered, func() {
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

		projectKey = helpers.UniqueResourceName("architecture-project")
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
	// FileGraph
	//
	// Architecture graphs are only populated once a scanner analysis of the
	// project has run. This suite has no scanner-analysis helper available, so
	// the real project created above is guaranteed to have no analysis data.
	// That means a "no graph data" style response (an empty/opaque payload, or
	// a not-found-ish error) is a legitimate, expected outcome here. Per
	// product decision, the Architecture service does NOT require a license to
	// respond meaningfully, so the one thing this test must never accept is a
	// license-gate response: not 402 Payment Required, and — since this
	// suite's BeforeSuite has already confirmed the server is Enterprise+ and
	// architecture analysis itself doesn't require a license — not 403
	// Forbidden either.
	// =========================================================================
	Describe("FileGraph", func() {
		Context("Functional Tests", func() {
			It("should return file graph data or a data-availability error, never a license gate", func() {
				result, resp, err := client.V2.Architecture.FileGraph(context.Background(), &sonar.ArchitectureFileGraphOptions{
					ProjectKey: projectKey,
					BranchKey:  "main",
					Source:     "java",
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
					// Data-availability condition (no analysis has ever run against this
					// project), not an enterprise/license gate: explicitly rule out the
					// license-gate status codes.
					Expect(resp.StatusCode).NotTo(Equal(http.StatusPaymentRequired))
					Expect(resp.StatusCode).NotTo(Equal(http.StatusForbidden))
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})
})
