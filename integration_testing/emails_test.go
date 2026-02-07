package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("Emails Service", Ordered, func() {
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
	// Send
	// =========================================================================
	Describe("Send", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Emails.Send(nil)
				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
			})

			It("should fail with missing to address", func() {
				resp, err := client.Emails.Send(&sonar.EmailsSendOption{
					Message: "Test message",
				})
				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
			})

			It("should fail with missing message", func() {
				resp, err := client.Emails.Send(&sonar.EmailsSendOption{
					To: "test@example.com",
				})
				Expect(resp).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Message"))
			})
		})

		Context("Functional Tests", func() {
			It("should attempt to send test email with valid parameters", func() {
				// Email sending requires SMTP to be configured
				// If not configured, the API returns an error
				resp, err := client.Emails.Send(&sonar.EmailsSendOption{
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
	})
})
