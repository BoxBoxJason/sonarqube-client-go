package integration_testing_test

import (
"net/http"
"regexp"

. "github.com/onsi/ginkgo/v2"
. "github.com/onsi/gomega"

sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Server Service", Ordered, func() {
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
	// Version
	// =========================================================================
	Describe("Version", func() {
		Context("Functional Tests", func() {
			It("should return server version", func() {
				version, resp, err := client.Server.Version()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(version).NotTo(BeNil())
				Expect(*version).NotTo(BeEmpty())
			})

			It("should return version in expected format", func() {
				version, resp, err := client.Server.Version()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(version).NotTo(BeNil())

				// SonarQube version format: X.Y.Z or X.Y.Z.BUILDNUMBER
				versionPattern := regexp.MustCompile(`^\d+\.\d+(\.\d+)?(\.\d+)?$`)
				Expect(versionPattern.MatchString(*version)).To(BeTrue(),
					"Version %s should match pattern X.Y.Z or X.Y.Z.BUILD", *version)
			})

			It("should return consistent version on multiple calls", func() {
				version1, resp1, err := client.Server.Version()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp1.StatusCode).To(Equal(http.StatusOK))

				version2, resp2, err := client.Server.Version()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp2.StatusCode).To(Equal(http.StatusOK))

				Expect(*version1).To(Equal(*version2))
			})

			It("should return version with major version >= 9", func() {
				version, resp, err := client.Server.Version()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(version).NotTo(BeNil())

				// Parse major version
				majorVersionPattern := regexp.MustCompile(`^(\d+)\.`)
				matches := majorVersionPattern.FindStringSubmatch(*version)
				Expect(matches).To(HaveLen(2), "Should extract major version from %s", *version)
			})
		})
	})
})
