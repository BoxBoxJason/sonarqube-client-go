package enterprise_test

import (
	"context"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

var _ = Describe("Audit Logs Service", Ordered, func() {
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
			It("should download audit logs for a recent time range", func() {
				to := time.Now().UTC()
				from := to.Add(-24 * time.Hour)

				result, resp, err := client.AuditLogs.Download(context.Background(), &sonar.AuditLogsDownloadOptions{
					From: from.Format(time.RFC3339),
					To:   to.Format(time.RFC3339),
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})
	})
})
