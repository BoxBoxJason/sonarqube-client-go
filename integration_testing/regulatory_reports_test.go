package integration_testing_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

var _ = Describe("Regulatory Reports Service", Ordered, func() {
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
				result, resp, err := client.RegulatoryReports.Download(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required Project", func() {
				result, resp, err := client.RegulatoryReports.Download(context.Background(), &sonar.RegulatoryReportsDownloadOptions{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should download or return an enterprise-only error", func() {
				result, resp, err := client.RegulatoryReports.Download(context.Background(), &sonar.RegulatoryReportsDownloadOptions{
					Project: "non-existent-project",
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
