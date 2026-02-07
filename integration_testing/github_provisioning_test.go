package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("GithubProvisioning Service", Ordered, func() {
	var client *sonar.Client

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	// =========================================================================
	// Check
	// =========================================================================
	// Note: The GithubProvisioning service only exposes the Check() method.
	// Enable/disable flows mentioned in issue #121 are not available via the API.
	Describe("Check", func() {
		Context("Functional Tests", func() {
			It("should check GitHub provisioning configuration", func() {
				result, resp, err := client.GithubProvisioning.Check()
				// Skip if GitHub integration is not configured or API not available
				if resp != nil && (resp.StatusCode == http.StatusNotFound ||
					resp.StatusCode == http.StatusBadRequest ||
					resp.StatusCode == http.StatusUnauthorized ||
					resp.StatusCode == http.StatusForbidden) {
					Skip("GitHub provisioning is not configured or API not available in this SonarQube instance")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should return valid provisioning status structure", func() {
				result, resp, err := client.GithubProvisioning.Check()
				// Skip if GitHub integration is not configured
				if resp != nil && (resp.StatusCode == http.StatusNotFound ||
					resp.StatusCode == http.StatusBadRequest ||
					resp.StatusCode == http.StatusUnauthorized ||
					resp.StatusCode == http.StatusForbidden) {
					Skip("GitHub provisioning is not configured or API not available in this SonarQube instance")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())

				// Verify Application level status exists and has proper structure
				Expect(result.Application.AutoProvisioning.Status).NotTo(BeNil())
				Expect(result.Application.Jit.Status).NotTo(BeNil())

				// If installations are present, verify their structure
				if len(result.Installations) > 0 {
					for _, installation := range result.Installations {
						Expect(installation.Organization).NotTo(BeEmpty())
						Expect(installation.AutoProvisioning.Status).NotTo(BeNil())
						Expect(installation.Jit.Status).NotTo(BeNil())
					}
				}
			})
		})
	})
})
