package integration_testing_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
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
	// Disable
	// =========================================================================
	Describe("Disable", func() {
		Context("Functional Tests", func() {
			It("should succeed or return an enterprise-only error", func() {
				resp, err := client.ScimManagement.Disable(context.Background())
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
				}
			})
		})
	})

	// =========================================================================
	// Enable
	// =========================================================================
	Describe("Enable", func() {
		Context("Functional Tests", func() {
			It("should succeed or return an enterprise-only error", func() {
				resp, err := client.ScimManagement.Enable(context.Background())
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
				}
			})
		})
	})

	// =========================================================================
	// Status
	// =========================================================================
	Describe("Status", func() {
		Context("Functional Tests", func() {
			It("should succeed or return an enterprise-only error", func() {
				result, resp, err := client.ScimManagement.Status(context.Background())
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})
})
