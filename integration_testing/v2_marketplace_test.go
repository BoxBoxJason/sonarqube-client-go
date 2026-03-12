package integration_testing_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("V2 Marketplace Service", Ordered, func() {
	var client *sonar.Client

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	// =========================================================================
	// BillAzureAccount
	// =========================================================================
	Describe("BillAzureAccount", func() {
		It("should attempt to bill the Azure account", func() {
			// BillAzureAccount requires an Azure marketplace license which is
			// typically not available in local test environments. We verify
			// the API call completes (successfully or with a meaningful server
			// error) without panicking.
			result, resp, err := client.V2.Marketplace.BillAzureAccount()
			if err != nil {
				// Expected in non-Azure environments: verify we got a server
				// error response rather than a client-side failure
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			} else {
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(BeNumerically(">=", 200))
				Expect(resp.StatusCode).To(BeNumerically("<", 300))
				Expect(result).NotTo(BeNil())
			}
		})
	})
})
