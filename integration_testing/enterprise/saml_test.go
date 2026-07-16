package enterprise_test

import (
	"context"
	"encoding/base64"
	"net/http"

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
		Context("Functional Tests", func() {
			It("should return a genuine validation outcome for a well-formed but bogus SAML response", func() {
				// Valid base64, but not a real IdP-signed assertion. This is
				// well-formed enough to reach the server's actual SAML validation
				// logic instead of being rejected purely on transport/encoding
				// malformation, so the server should give a real, definitive answer.
				fakeAssertion := base64.StdEncoding.EncodeToString([]byte("<samlp:Response>not-a-real-assertion</samlp:Response>"))

				result, resp, err := client.Saml.Validation(context.Background(), &sonar.SamlValidationOptions{
					SAMLResponse: fakeAssertion,
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
					// SAML is not itself a license-gated feature: a 402/403 here
					// would indicate an (incorrect) license/edition gate rather than
					// a legitimate validation-failure response. It may, however,
					// legitimately be unconfigured on the test server, which surfaces
					// as some other 4xx - a distinct condition from a license gate,
					// so it is tolerated narrowly here (4xx only, license codes excluded).
					Expect(resp.StatusCode).NotTo(BeElementOf(http.StatusPaymentRequired, http.StatusForbidden))
					Expect(resp.StatusCode).To(BeNumerically(">=", http.StatusBadRequest))
					Expect(resp.StatusCode).To(BeNumerically("<", http.StatusInternalServerError))
				} else {
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					Expect(result).NotTo(BeNil())
					Expect(result).NotTo(BeEmpty())
				}
			})
		})
	})

	// =========================================================================
	// ValidationInit
	// =========================================================================
	Describe("ValidationInit", func() {
		Context("Functional Tests", func() {
			It("should initiate the SAML validation flow or report a non-license-related condition", func() {
				resp, err := client.Saml.ValidationInit(context.Background())
				if err != nil {
					Expect(resp).NotTo(BeNil())
					// Same reasoning as Validation above: SAML is not license-gated,
					// so a license/edition-gate status here would be a genuine defect,
					// while an "unconfigured" style 4xx is a legitimate, distinct
					// condition on a test server that has no SAML identity provider set up.
					Expect(resp.StatusCode).NotTo(BeElementOf(http.StatusPaymentRequired, http.StatusForbidden))
				} else {
					Expect(resp.StatusCode).To(BeNumerically(">=", http.StatusOK))
					Expect(resp.StatusCode).To(BeNumerically("<", http.StatusMultipleChoices))
				}
			})
		})
	})
})
