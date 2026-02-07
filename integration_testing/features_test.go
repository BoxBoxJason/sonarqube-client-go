package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("Features Service", Ordered, func() {
	var (
		client *sonar.Client
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	// =========================================================================
	// List
	// =========================================================================
	Describe("List", func() {
		Context("Valid Requests", func() {
			It("should list supported features with non-empty values", func() {
				result, resp, err := client.Features.List()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())

				// Validate the returned list
				features := *result
				// Some SonarQube instances may have no features (e.g., community edition)
				// We just validate that when features exist, they are non-empty strings
				for _, feature := range features {
					Expect(feature).NotTo(BeEmpty(), "Feature name should not be empty")
				}
			})

			It("should return consistent results on multiple calls", func() {
				result1, resp1, err := client.Features.List()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp1.StatusCode).To(Equal(http.StatusOK))

				result2, resp2, err := client.Features.List()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp2.StatusCode).To(Equal(http.StatusOK))

				// Compare actual contents, not just length
				Expect(*result1).To(ConsistOf(*result2), "Features list should be consistent across calls")
			})
		})
	})
})
