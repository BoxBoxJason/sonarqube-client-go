package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("V2 System Service", Ordered, func() {
	var client *sonar.Client

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	// =========================================================================
	// GetMigrationsStatus
	// =========================================================================
	Describe("GetMigrationsStatus", func() {
		It("should return migration status", func() {
			result, resp, err := client.V2.System.GetMigrationsStatus()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Status).NotTo(BeEmpty())
		})
	})

	// =========================================================================
	// CheckLiveness
	// =========================================================================
	Describe("CheckLiveness", func() {
		Context("without passcode", func() {
			It("should return a successful liveness check", func() {
				resp, err := client.V2.System.CheckLiveness(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(BeNumerically(">=", 200))
				Expect(resp.StatusCode).To(BeNumerically("<", 300))
			})
		})

		Context("with empty passcode", func() {
			It("should return a successful liveness check", func() {
				resp, err := client.V2.System.CheckLiveness(&sonar.SystemPasscodeOptionV2{})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(BeNumerically(">=", 200))
				Expect(resp.StatusCode).To(BeNumerically("<", 300))
			})
		})
	})

	// =========================================================================
	// GetHealth
	// =========================================================================
	Describe("GetHealth", func() {
		Context("without passcode", func() {
			It("should return system health", func() {
				result, resp, err := client.V2.System.GetHealth(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Status).NotTo(BeEmpty())
				Expect(result.Status).To(BeElementOf("GREEN", "YELLOW", "RED"))
			})
		})

		Context("with empty passcode", func() {
			It("should return system health", func() {
				result, resp, err := client.V2.System.GetHealth(&sonar.SystemPasscodeOptionV2{})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Status).NotTo(BeEmpty())
			})
		})
	})
})
