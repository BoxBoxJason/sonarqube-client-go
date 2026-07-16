package enterprise_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

var _ = Describe("SCA V2 Service", Ordered, func() {
	var (
		client                     *sonar.Client
		cleanup                    *helpers.CleanupManager
		projectKey                 string
		originalEnablement         bool
		originalEnablementCaptured bool
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
		cleanup = helpers.NewCleanupManager(client)

		projectKey = helpers.UniqueResourceName("sca-project")
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
		// Restore SCA's original enablement state so this suite doesn't leave
		// the server's SCA feature flag permanently flipped for other tests /
		// other suites run against the same instance. Only meaningful when a
		// license was active, since SetEnablement itself is license-gated.
		if originalEnablementCaptured {
			_, resp, err := client.V2.Sca.SetEnablement(context.Background(), &sonar.ScaSetEnablementOptions{
				Enablement: originalEnablement,
			})
			if err != nil {
				GinkgoWriter.Printf("Failed to restore original SCA enablement (%v): %v\n", originalEnablement, err)
			} else {
				Expect(resp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))
			}
		}

		errors := cleanup.Cleanup()
		for _, err := range errors {
			GinkgoWriter.Printf("Cleanup error: %v\n", err)
		}
	})

	// =========================================================================
	// Enablement toggle. Live-verified against an unlicensed instance: unlike
	// what the endpoint's shape suggests, both GetEnablement and SetEnablement
	// require an active license and fail with 403 {"message":"Not available"}
	// without one - this is a genuine license gate, not a plain feature flag.
	// =========================================================================
	Describe("Enablement", func() {
		Context("Functional Tests", func() {
			It("should get the current enablement state, or fail with 403 when unlicensed", func() {
				result, resp, err := client.V2.Sca.GetEnablement(context.Background())

				if helpers.HasActiveLicense(client) {
					Expect(err).NotTo(HaveOccurred())
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					Expect(result).NotTo(BeNil())

					originalEnablement = result.Enablement
					originalEnablementCaptured = true
				} else {
					Expect(err).To(HaveOccurred())
					Expect(resp).NotTo(BeNil())
					Expect(resp.StatusCode).To(Equal(http.StatusForbidden))
					Expect(result).To(BeNil())
				}
			})

			It("should enable SCA, or fail with 403 when unlicensed", func() {
				result, resp, err := client.V2.Sca.SetEnablement(context.Background(), &sonar.ScaSetEnablementOptions{
					Enablement: true,
				})

				if helpers.HasActiveLicense(client) {
					Expect(originalEnablementCaptured).To(BeTrue(), "GetEnablement must run before SetEnablement")
					Expect(err).NotTo(HaveOccurred())
					Expect(resp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))
					Expect(result).NotTo(BeNil())
				} else {
					Expect(err).To(HaveOccurred())
					Expect(resp).NotTo(BeNil())
					Expect(resp.StatusCode).To(Equal(http.StatusForbidden))
					Expect(result).To(BeNil())
				}
			})

			It("should reflect the enabled state on a subsequent read, or fail with 403 when unlicensed", func() {
				result, resp, err := client.V2.Sca.GetEnablement(context.Background())

				if helpers.HasActiveLicense(client) {
					Expect(err).NotTo(HaveOccurred())
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					Expect(result).NotTo(BeNil())
					Expect(result.Enablement).To(BeTrue())
				} else {
					Expect(err).To(HaveOccurred())
					Expect(resp).NotTo(BeNil())
					Expect(resp.StatusCode).To(Equal(http.StatusForbidden))
					Expect(result).To(BeNil())
				}
			})
		})
	})

	// =========================================================================
	// ListClis: a static CLI listing, independent of any project's analysis
	// state or a license.
	// =========================================================================
	Describe("ListClis", func() {
		Context("Functional Tests", func() {
			It("should list available SCA CLI downloads", func() {
				result, resp, err := client.V2.Sca.ListClis(context.Background(), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})
	})

	// =========================================================================
	// SearchDependencyRisks / SearchReleases / GetSbomReport against the real
	// (but unanalyzed) project. Live-verified against an unlicensed instance:
	// since SCA can never be successfully enabled without a license (per the
	// Enablement block above), these downstream endpoints consistently fail
	// with 403 {"message":"SCA feature is not enabled for '<project>'"}. On a
	// licensed instance where SCA was actually enabled, no scan has run
	// against this project (this suite has no analysis helper), so an empty
	// result set is the expected data-availability outcome, and a license
	// gate (402/403) must never occur.
	// =========================================================================
	Describe("SearchDependencyRisks", func() {
		Context("Functional Tests", func() {
			It("should search dependency risks for the real project, or fail with 403 when SCA isn't enabled", func() {
				result, resp, err := client.V2.Sca.SearchDependencyRisks(context.Background(), &sonar.ScaDependencyRisksSearchOptions{
					ProjectKey: projectKey,
				})

				if !helpers.HasActiveLicense(client) {
					Expect(err).To(HaveOccurred())
					Expect(resp).NotTo(BeNil())
					Expect(resp.StatusCode).To(Equal(http.StatusForbidden))
					Expect(result).To(BeNil())
				} else if err != nil {
					Expect(resp).NotTo(BeNil())
					Expect(resp.StatusCode).NotTo(Equal(http.StatusPaymentRequired))
					Expect(resp.StatusCode).NotTo(Equal(http.StatusForbidden))
				} else {
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})

	Describe("SearchReleases", func() {
		Context("Functional Tests", func() {
			It("should search releases for the real project, or fail with 403 when SCA isn't enabled", func() {
				result, resp, err := client.V2.Sca.SearchReleases(context.Background(), &sonar.ScaReleasesSearchOptions{
					ProjectKey: projectKey,
				})

				if !helpers.HasActiveLicense(client) {
					Expect(err).To(HaveOccurred())
					Expect(resp).NotTo(BeNil())
					Expect(resp.StatusCode).To(Equal(http.StatusForbidden))
					Expect(result).To(BeNil())
				} else if err != nil {
					Expect(resp).NotTo(BeNil())
					Expect(resp.StatusCode).NotTo(Equal(http.StatusPaymentRequired))
					Expect(resp.StatusCode).NotTo(Equal(http.StatusForbidden))
				} else {
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})

	Describe("GetSbomReport", func() {
		Context("Functional Tests", func() {
			It("should return an SBOM report for the real project, or fail with 403 when SCA isn't enabled", func() {
				result, resp, err := client.V2.Sca.GetSbomReport(context.Background(), &sonar.ScaSbomReportOptions{
					Component: projectKey,
					Type:      sonar.ScaSbomReportTypeCycloneDX,
					Format:    sonar.ScaSbomReportFormatJSON,
				})

				if !helpers.HasActiveLicense(client) {
					Expect(err).To(HaveOccurred())
					Expect(resp).NotTo(BeNil())
					Expect(resp.StatusCode).To(Equal(http.StatusForbidden))
					Expect(result).To(BeNil())
				} else if err != nil {
					Expect(resp).NotTo(BeNil())
					Expect(resp.StatusCode).NotTo(Equal(http.StatusPaymentRequired))
					Expect(resp.StatusCode).NotTo(Equal(http.StatusForbidden))
				} else {
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})
})
