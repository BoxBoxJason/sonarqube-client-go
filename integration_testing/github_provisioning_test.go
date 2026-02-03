package integration_testing_test

import (
"net/http"

. "github.com/onsi/ginkgo/v2"
. "github.com/onsi/gomega"

sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("GithubProvisioning Service", Ordered, func() {
	var client *sonargo.Client

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	// =========================================================================
	// Check
	// =========================================================================
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
		})
	})
})
