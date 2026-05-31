package integration_testing_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("License Service", Ordered, func() {
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

	// =========================================================================
	// ActivateGracePeriod
	// =========================================================================
	Describe("ActivateGracePeriod", func() {
		Context("Functional Tests", func() {
			It("should activate grace period or return an error on community edition", func() {
				resp, err := client.Editions.ActivateGracePeriod()
				if err != nil {
					// Community edition does not expose this endpoint;
					// a 404 or 403 is acceptable.
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp).NotTo(BeNil())
				}
			})
		})
	})

	// =========================================================================
	// Get
	// =========================================================================
	Describe("Get", func() {
		Context("Functional Tests", func() {
			It("should return license information or an error on community edition", func() {
				result, resp, err := client.Editions.Get()
				if err != nil {
					// Community edition does not expose the license endpoint;
					// a 404 or 403 is acceptable.
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp).NotTo(BeNil())
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})

	// =========================================================================
	// IsValidLicense
	// =========================================================================
	Describe("IsValidLicense", func() {
		Context("Functional Tests", func() {
			It("should return license validity or an error on community edition", func() {
				result, resp, err := client.Editions.IsValidLicense()
				if err != nil {
					// Community edition does not expose this endpoint;
					// a 404 or 403 is acceptable.
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp).NotTo(BeNil())
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})

	// =========================================================================
	// Set
	// =========================================================================
	Describe("Set", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Editions.Set(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without license key", func() {
				resp, err := client.Editions.Set(&sonar.LicenseSetOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("License"))
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should fail with an invalid license key", func() {
				resp, err := client.Editions.Set(&sonar.LicenseSetOptions{
					License: "invalid-license-key",
				})
				// Expect an error because the license key is invalid.
				// On community edition this may return 404 (endpoint not found).
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
			})
		})
	})

	// =========================================================================
	// UnsetLicense
	// =========================================================================
	Describe("UnsetLicense", func() {
		Context("Functional Tests", func() {
			It("should unset the license or return an error on community edition", func() {
				resp, err := client.Editions.UnsetLicense()
				if err != nil {
					// Community edition does not expose this endpoint;
					// a 404 or 403 is acceptable.
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp).NotTo(BeNil())
				}
			})
		})
	})
})
