package integration_testing_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("Audit Logs Service", Ordered, func() {
	var (
		client  *sonar.Client
		cleanup *helpers.CleanupManager
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
		cleanup = helpers.NewCleanupManager(client)
	})

	AfterAll(func() {
		cleanup.Cleanup()
	})

	Describe("Download", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.AuditLogs.Download(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required From", func() {
				result, resp, err := client.AuditLogs.Download(context.Background(), &sonar.AuditLogsDownloadOptions{
					To: "2024-12-31T23:59:59+00:00",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required To", func() {
				result, resp, err := client.AuditLogs.Download(context.Background(), &sonar.AuditLogsDownloadOptions{
					From: "2024-01-01T00:00:00+00:00",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should download or return an enterprise-only error", func() {
				result, resp, err := client.AuditLogs.Download(context.Background(), &sonar.AuditLogsDownloadOptions{
					From: "2024-01-01T00:00:00+00:00",
					To:   "2024-12-31T23:59:59+00:00",
				})
				if err != nil {
					// Enterprise-only endpoints return 4xx on Community Edition.
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})
})
