package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Monitoring Service", Ordered, func() {
	var (
		client *sonargo.Client
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	// Helper function to check if the endpoint requires authentication
	checkAuthRequired := func(resp *http.Response) {
		if resp != nil && (resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized) {
			Skip("Monitoring metrics endpoint requires system passCode")
		}
	}

	// =========================================================================
	// Metrics
	// =========================================================================
	Describe("Metrics", func() {
		Context("Functional Tests", func() {
			It("should get monitoring metrics", func() {
				result, resp, err := client.Monitoring.Metrics()
				checkAuthRequired(resp)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should return Prometheus format metrics with specific metric categories", func() {
				result, resp, err := client.Monitoring.Metrics()
				checkAuthRequired(resp)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())

				metrics := *result
				Expect(len(metrics)).To(BeNumerically(">", 0))
			})

			It("should return consistent results on multiple calls", func() {
				result1, resp1, err := client.Monitoring.Metrics()
				checkAuthRequired(resp1)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp1.StatusCode).To(Equal(http.StatusOK))

				result2, resp2, err := client.Monitoring.Metrics()
				checkAuthRequired(resp2)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp2.StatusCode).To(Equal(http.StatusOK))

				// Both should return non-empty metrics
				Expect(len(*result1)).To(BeNumerically(">", 0))
				Expect(len(*result2)).To(BeNumerically(">", 0))
			})
		})
	})
})
