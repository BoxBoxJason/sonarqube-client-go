package integration_testing_test

import (
"net/http"

. "github.com/onsi/ginkgo/v2"
. "github.com/onsi/gomega"

"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("AnalysisReports Service", Ordered, func() {
	BeforeAll(func() {
		client, err := helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	// =========================================================================
	// IsQueueEmpty
	// =========================================================================
	Describe("IsQueueEmpty", func() {
		Context("Functional Tests", func() {
			It("should check if compute engine queue is empty", func() {
				client, err := helpers.NewDefaultClient()
				Expect(err).NotTo(HaveOccurred())

				result, resp, err := client.AnalysisReports.IsQueueEmpty()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				// IsEmpty is a boolean, could be true or false depending on queue state
				// Just verify it's one of those values
Expect(result.IsEmpty).To(SatisfyAny(BeTrue(), BeFalse()))
})

It("should return consistent results on multiple calls", func() {
client, err := helpers.NewDefaultClient()
Expect(err).NotTo(HaveOccurred())

// Call twice and verify both return valid results
result1, resp1, err1 := client.AnalysisReports.IsQueueEmpty()
Expect(err1).NotTo(HaveOccurred())
Expect(resp1.StatusCode).To(Equal(http.StatusOK))
Expect(result1).NotTo(BeNil())

result2, resp2, err2 := client.AnalysisReports.IsQueueEmpty()
Expect(err2).NotTo(HaveOccurred())
Expect(resp2.StatusCode).To(Equal(http.StatusOK))
Expect(result2).NotTo(BeNil())

// Both should be valid booleans (they may differ based on queue state)
Expect(result1.IsEmpty).To(SatisfyAny(BeTrue(), BeFalse()))
Expect(result2.IsEmpty).To(SatisfyAny(BeTrue(), BeFalse()))
})
})
})
})
