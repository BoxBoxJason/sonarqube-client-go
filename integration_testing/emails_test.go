package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Emails Service", Ordered, func() {
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
	// Send
	// =========================================================================
	Describe("Send", func() {
		Context("Functional Tests", func() {
			It("should attempt to send test email with valid parameters", func() {
				// Email sending requires SMTP to be configured
				// If not configured, the API returns an error
				resp, err := client.Emails.Send(&sonargo.EmailsSendOption{
					To:      "test@example.com",
					Message: "Test message from e2e tests",
					Subject: "Test Subject",
				})
				// Skip if API not available
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("Emails API is not available in this SonarQube version")
				}
				// Email may fail if SMTP is not configured - this is expected
				// We just verify the endpoint exists and accepts valid parameters
				if resp != nil {
					// 200 = success, 400 = validation error, 500 = SMTP not configured
					Expect(resp.StatusCode).To(SatisfyAny(
						Equal(http.StatusOK),
						Equal(http.StatusNoContent),
						Equal(http.StatusBadRequest),
						Equal(http.StatusInternalServerError),
					))
				} else {
					// If no response, an error should have occurred
					Expect(err).To(HaveOccurred())
				}
			})
		})

		Context("Error Handling", func() {
			It("should fail with missing to address", func() {
				_, err := client.Emails.Send(&sonargo.EmailsSendOption{
					Message: "Test message",
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with missing message", func() {
				_, err := client.Emails.Send(&sonargo.EmailsSendOption{
					To: "test@example.com",
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, err := client.Emails.Send(nil)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
