package integration_testing_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

var _ = Describe("Governance Reports Service", Ordered, func() {
	var client *sonar.Client

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	// =========================================================================
	// Download
	// =========================================================================
	Describe("Download", func() {
		Context("Functional Tests", func() {
			It("should download or return an enterprise-only error", func() {
				result, resp, err := client.GovernanceReports.Download(context.Background(), &sonar.GovernanceReportsDownloadOptions{
					ComponentKey: "nonexistent-portfolio",
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
	// Status
	// =========================================================================
	Describe("Status", func() {
		Context("Functional Tests", func() {
			It("should return status or an enterprise-only error", func() {
				result, resp, err := client.GovernanceReports.Status(context.Background(), &sonar.GovernanceReportsStatusOptions{
					ComponentKey: "nonexistent-portfolio",
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
	// Subscribe
	// =========================================================================
	Describe("Subscribe", func() {
		Context("Functional Tests", func() {
			It("should subscribe or return an enterprise-only error", func() {
				resp, err := client.GovernanceReports.Subscribe(context.Background(), &sonar.GovernanceReportsSubscribeOptions{
					ComponentKey: "nonexistent-portfolio",
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
				}
			})
		})
	})

	// =========================================================================
	// Unsubscribe
	// =========================================================================
	Describe("Unsubscribe", func() {
		Context("Functional Tests", func() {
			It("should unsubscribe or return an enterprise-only error", func() {
				resp, err := client.GovernanceReports.Unsubscribe(context.Background(), &sonar.GovernanceReportsUnsubscribeOptions{
					ComponentKey: "nonexistent-portfolio",
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
				}
			})
		})
	})

	// =========================================================================
	// UpdateFrequency
	// =========================================================================
	Describe("UpdateFrequency", func() {
		Context("Functional Tests", func() {
			It("should update frequency or return an enterprise-only error", func() {
				resp, err := client.GovernanceReports.UpdateFrequency(context.Background(), &sonar.GovernanceReportsUpdateFrequencyOptions{
					ComponentKey: "nonexistent-portfolio",
					Frequency:    "WEEKLY",
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
				}
			})
		})
	})

	// =========================================================================
	// UpdateRecipients
	// =========================================================================
	Describe("UpdateRecipients", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.GovernanceReports.UpdateRecipients(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required recipients", func() {
				resp, err := client.GovernanceReports.UpdateRecipients(context.Background(), &sonar.GovernanceReportsUpdateRecipientsOptions{
					ComponentKey: "my-portfolio",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Recipients"))
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should update recipients or return an enterprise-only error", func() {
				resp, err := client.GovernanceReports.UpdateRecipients(context.Background(), &sonar.GovernanceReportsUpdateRecipientsOptions{
					ComponentKey: "nonexistent-portfolio",
					Recipients:   "test@example.com",
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
				}
			})
		})
	})
})
