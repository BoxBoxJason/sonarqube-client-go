package enterprise_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

var _ = Describe("SCIM Management Service", Ordered, func() {
	var client *sonar.Client

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	// =========================================================================
	// Status / Disable: live-verified to work unconditionally, regardless of
	// whether SAML is configured or SCIM was ever enabled.
	// =========================================================================
	Describe("Status", func() {
		Context("Functional Tests", func() {
			It("should report the current SCIM provisioning status", func() {
				result, resp, err := client.ScimManagement.Status(context.Background())
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})
	})

	Describe("Disable", func() {
		Context("Functional Tests", func() {
			It("should succeed even when SCIM was never enabled", func() {
				// Live-verified: Disable returns 204 unconditionally, even
				// against a server where SCIM has never been configured.
				resp, err := client.ScimManagement.Disable(context.Background())
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))

				result, statusResp, err := client.ScimManagement.Status(context.Background())
				Expect(err).NotTo(HaveOccurred())
				Expect(statusResp.StatusCode).To(Equal(http.StatusOK))
				Expect(result.Enabled).To(BeFalse())
			})
		})
	})

	// =========================================================================
	// Enable: live-verified via the SonarQube container logs to fail with 500
	// and java.lang.IllegalStateException: "SAML must be enabled to enable
	// SCIM." - a genuine server-side precondition (SCIM provisioning rides on
	// top of SAML SSO), not a license gate or an SDK bug. This suite has no
	// SAML IdP configured, so Enable is expected to fail with exactly that
	// error rather than a guessed-at success.
	// =========================================================================
	Describe("Enable", func() {
		Context("Functional Tests", func() {
			It("should fail without a configured SAML identity provider", func() {
				resp, err := client.ScimManagement.Enable(context.Background())
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))

				result, statusResp, statusErr := client.ScimManagement.Status(context.Background())
				Expect(statusErr).NotTo(HaveOccurred())
				Expect(statusResp.StatusCode).To(Equal(http.StatusOK))
				Expect(result.Enabled).To(BeFalse())
			})
		})
	})
})
