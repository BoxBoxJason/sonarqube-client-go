package enterprise_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

var _ = Describe("Governance Reports Service", Ordered, func() {
	var (
		client       *sonar.Client
		cleanup      *helpers.CleanupManager
		portfolioKey string
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
		cleanup = helpers.NewCleanupManager(client)

		portfolioKey = helpers.UniqueResourceName("governance-portfolio")
		_, err = client.Views.Create(context.Background(), &sonar.ViewsCreateOptions{
			Key:  portfolioKey,
			Name: portfolioKey,
		})
		Expect(err).NotTo(HaveOccurred())
		cleanup.RegisterCleanup("portfolio", portfolioKey, func() error {
			_, err := client.Views.Delete(context.Background(), &sonar.ViewsDeleteOptions{Key: portfolioKey})

			return err
		})

		// Subscribe is live-verified to fail with 400 "User 'admin' has no
		// email" on a fresh instance where the admin account has no email
		// configured. Since a subscription is inherently tied to a real
		// recipient, this is a genuine prerequisite (not a guess to work
		// around), so it is satisfied here rather than tolerated as an error.
		_, _, err = client.Users.Update(context.Background(), &sonar.UsersUpdateOptions{
			Login: "admin",
			Email: "e2e-governance-reports@example.com",
		})
		Expect(err).NotTo(HaveOccurred())
	})

	AfterAll(func() {
		errors := cleanup.Cleanup()
		for _, err := range errors {
			GinkgoWriter.Printf("Cleanup error: %v\n", err)
		}
	})

	// =========================================================================
	// Status
	// =========================================================================
	Describe("Status", func() {
		Context("Functional Tests", func() {
			It("should return real status for a real portfolio", func() {
				// Status reflects report generation rights/metadata, it does not
				// require the portfolio to have been analyzed, so this must
				// succeed for real on a licensed Enterprise instance.
				result, resp, err := client.GovernanceReports.Status(context.Background(), &sonar.GovernanceReportsStatusOptions{
					ComponentKey: portfolioKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})
	})

	// =========================================================================
	// Subscribe / Unsubscribe
	// =========================================================================
	Describe("Subscribe", func() {
		Context("Functional Tests", func() {
			It("should subscribe the current user to reports for the portfolio", func() {
				resp, err := client.GovernanceReports.Subscribe(context.Background(), &sonar.GovernanceReportsSubscribeOptions{
					ComponentKey: portfolioKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))
			})
		})
	})

	Describe("Unsubscribe", func() {
		Context("Functional Tests", func() {
			It("should unsubscribe the current user from reports for the portfolio", func() {
				// Subscribe first so there is something real to unsubscribe from.
				_, err := client.GovernanceReports.Subscribe(context.Background(), &sonar.GovernanceReportsSubscribeOptions{
					ComponentKey: portfolioKey,
				})
				Expect(err).NotTo(HaveOccurred())

				resp, err := client.GovernanceReports.Unsubscribe(context.Background(), &sonar.GovernanceReportsUnsubscribeOptions{
					ComponentKey: portfolioKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))
			})
		})
	})

	// =========================================================================
	// UpdateFrequency
	// =========================================================================
	Describe("UpdateFrequency", func() {
		Context("Functional Tests", func() {
			It("should update the report frequency for the portfolio", func() {
				// Live-verified: the API rejects uppercase frequency values
				// (e.g. "WEEKLY") with 400 "must be one of: [daily, weekly,
				// monthly]" and only accepts lowercase.
				resp, err := client.GovernanceReports.UpdateFrequency(context.Background(), &sonar.GovernanceReportsUpdateFrequencyOptions{
					ComponentKey: portfolioKey,
					Frequency:    "weekly",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))
			})
		})
	})

	// =========================================================================
	// UpdateRecipients
	// =========================================================================
	Describe("UpdateRecipients", func() {
		Context("Functional Tests", func() {
			It("should update the report recipients for the portfolio", func() {
				resp, err := client.GovernanceReports.UpdateRecipients(context.Background(), &sonar.GovernanceReportsUpdateRecipientsOptions{
					ComponentKey: portfolioKey,
					Recipients:   "e2e-governance-reports@example.com",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(BeElementOf(http.StatusOK, http.StatusNoContent))
			})
		})
	})

	// =========================================================================
	// Download
	// =========================================================================
	Describe("Download", func() {
		Context("Functional Tests", func() {
			It("should fail with a specific data-availability error, never a license-gate error", func() {
				// Live-verified against an unanalyzed portfolio: unlike
				// regulatory/security reports (which succeed unconditionally
				// even with zero analysis data), governance report generation
				// genuinely requires analysis and fails with exactly
				// 400 {"errors":[{"msg":"No analysis has been done on this
				// component"}]} - never a license-gate status.
				result, resp, err := client.GovernanceReports.Download(context.Background(), &sonar.GovernanceReportsDownloadOptions{
					ComponentKey: portfolioKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
				Expect(err.Error()).To(ContainSubstring("No analysis has been done"))
			})
		})
	})

	// =========================================================================
	// Nonexistent component
	// =========================================================================
	Describe("Nonexistent portfolio", func() {
		Context("Functional Tests", func() {
			It("should return a not-found style error, never a license-gate error", func() {
				// Even on a confirmed Enterprise server, a nonexistent component
				// key must be reported as "not found", not "you need a license".
				_, resp, err := client.GovernanceReports.Status(context.Background(), &sonar.GovernanceReportsStatusOptions{
					ComponentKey: "e2e-nonexistent-governance-portfolio",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).NotTo(Equal(http.StatusPaymentRequired))
				Expect(resp.StatusCode).NotTo(Equal(http.StatusForbidden))
			})
		})
	})
})
