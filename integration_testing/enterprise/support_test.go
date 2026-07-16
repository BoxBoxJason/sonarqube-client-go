package enterprise_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

var _ = Describe("Support Service", Ordered, func() {
	var client *sonar.Client

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	// =========================================================================
	// Info
	// =========================================================================
	Describe("Info", func() {
		Context("Functional Tests", func() {
			It("should return system support information, or a license-not-found error when unlicensed", func() {
				// Unlike most endpoints covered by this suite, Support.Info is
				// genuinely license-gated: it has been live-verified (see the
				// WARNING in sonar/support_service.go) to require an actual
				// installed commercial license independent of Enterprise Edition +
				// 'Administer System' permission, failing with HTTP 400 and body
				// {"errors":[{"msg":"License not found"}]} otherwise. This is one
				// of the "some other cases" the license-optional design
				// acknowledges, so branch on whether a license is actually active.
				result, resp, err := client.Support.Info(context.Background())

				if helpers.HasActiveLicense(client) {
					Expect(err).NotTo(HaveOccurred())
					Expect(resp).NotTo(BeNil())
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					Expect(result).NotTo(BeNil())

					// The System and SonarQube sections are always populated on a
					// real server (server ID, version, uptime, etc.), regardless of
					// which optional sections (e.g. Statistics) are included.
					Expect(result.System).NotTo(BeEmpty())
					Expect(result.SonarQube).NotTo(BeEmpty())

					GinkgoWriter.Printf("Support.Info succeeded with an active license\n")
				} else {
					Expect(err).To(HaveOccurred())
					Expect(resp).NotTo(BeNil())
					Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
					Expect(result).To(BeNil())

					GinkgoWriter.Printf("Support.Info failed with %q as expected: no active license\n", err.Error())
				}
			})
		})
	})
})
