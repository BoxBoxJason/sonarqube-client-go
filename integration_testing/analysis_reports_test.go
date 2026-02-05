package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("AnalysisReports Service", Ordered, func() {
	var (
		client *sonargo.Client
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	// =========================================================================
	// IsQueueEmpty
	// =========================================================================
	Describe("IsQueueEmpty", func() {
		Context("Functional Tests", func() {
			It("should check if compute engine queue is empty", func() {
				result, resp, err := client.AnalysisReports.IsQueueEmpty()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should return consistent results on multiple calls", func() {
				// Call twice and verify both return valid results
				result1, resp1, err1 := client.AnalysisReports.IsQueueEmpty()
				Expect(err1).NotTo(HaveOccurred())
				Expect(resp1.StatusCode).To(Equal(http.StatusOK))
				Expect(result1).NotTo(BeNil())

				result2, resp2, err2 := client.AnalysisReports.IsQueueEmpty()
				Expect(err2).NotTo(HaveOccurred())
				Expect(resp2.StatusCode).To(Equal(http.StatusOK))
				Expect(result2).NotTo(BeNil())
			})
		})
	})
})
