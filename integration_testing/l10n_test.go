package integration_testing_test

import (
"net/http"

. "github.com/onsi/ginkgo/v2"
. "github.com/onsi/gomega"

sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("L10N Service", Ordered, func() {
	var client *sonargo.Client

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	// =========================================================================
	// Index
	// =========================================================================
	Describe("Index", func() {
		Context("Functional Tests", func() {
			It("should get localization messages with default locale", func() {
				result, resp, err := client.L10N.Index(nil)
				// Skip if API not available
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("L10N API is not available in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Messages).NotTo(BeEmpty())
			})

			It("should get localization messages for specific locale", func() {
				result, resp, err := client.L10N.Index(&sonargo.L10NIndexOption{
					Locale: "en",
				})
				// Skip if API not available
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("L10N API is not available in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Messages).NotTo(BeEmpty())
			})

			It("should get localization messages with timestamp", func() {
				result, resp, err := client.L10N.Index(&sonargo.L10NIndexOption{
					Timestamp: "2020-01-01T00:00:00+0000",
				})
				// Skip if API not available
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("L10N API is not available in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})
	})
})
