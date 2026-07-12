package integration_testing_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

var _ = Describe("SAML Service", Ordered, func() {
	var client *sonar.Client

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	// =========================================================================
	// Validation
	// =========================================================================
	Describe("Validation", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Saml.Validation(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required SAMLResponse", func() {
				result, resp, err := client.Saml.Validation(context.Background(), &sonar.SamlValidationOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("SAMLResponse"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should validate or return an expected error", func() {
				result, resp, err := client.Saml.Validation(context.Background(), &sonar.SamlValidationOptions{
					SAMLResponse: "invalid-saml-response",
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})

	// =========================================================================
	// ValidationInit
	// =========================================================================
	Describe("ValidationInit", func() {
		Context("Functional Tests", func() {
			It("should initiate validation or return an expected error", func() {
				resp, err := client.Saml.ValidationInit(context.Background())
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
				}
			})
		})
	})
})
