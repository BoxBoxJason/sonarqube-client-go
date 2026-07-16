package enterprise_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

// This suite exercises the license-management API (client.V2.Entitlements)
// itself. Unlike the rest of the enterprise suite, a real license is
// *optional* here too, but this specific file is inherently about whether
// and how a license is provisioned. There is no way to know in this
// environment whether a real, valid license (e.g. via SONAR_LICENSE used by
// `make setup.sonar.enterprise`) will actually be installed when this suite
// runs, so every test below is written to be safe both with and without a
// license present.
//
// DeleteLicense, UpdateLicense and DeactivateOffline are intentionally never
// called anywhere in this file, not even for parameter-validation checks:
// DeactivateOffline and UpdateLicense take no options at all (nothing to
// validate client-side), and invoking any of the three for real would mutate
// license state on a possibly-shared server that may already carry a real,
// valid license. Only read-only calls (GetLicense, GetPurchasableFeatures)
// and client-side-only parameter validation (which never reaches the
// network) are exercised here.
var _ = Describe("Entitlements Service", func() {
	var client *sonar.Client

	BeforeEach(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	// =========================================================================
	// GetLicense: live-verified against an unlicensed instance - unlike its
	// name suggests, this is not an unconditionally-succeeding read. It fails
	// with 404 {"message":"License not found"} when no license is installed.
	// =========================================================================
	Describe("GetLicense", func() {
		It("should return license information when licensed, or a not-found error otherwise", func() {
			result, resp, err := client.V2.Entitlements.GetLicense(context.Background())

			if helpers.HasActiveLicense(client) {
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())

				GinkgoWriter.Printf("GetLicense: active license detected (edition=%q), asserting license-field values\n", result.Edition)
				Expect(result.Edition).NotTo(BeEmpty())
				Expect(result.ValidEdition).To(BeTrue())
			} else {
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				Expect(result).To(BeNil())

				GinkgoWriter.Printf("GetLicense: no active license detected, confirmed 404 as expected\n")
			}
		})
	})

	// =========================================================================
	// GetPurchasableFeatures: read-only, content depends on the specific
	// license/edition so only structural assertions are made.
	// =========================================================================
	Describe("GetPurchasableFeatures", func() {
		It("should return the list of purchasable features", func() {
			result, resp, err := client.V2.Entitlements.GetPurchasableFeatures(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})
	})

	// =========================================================================
	// Mutating operations: parameter validation only. These never reach the
	// network (validation happens client-side before the request is built),
	// so they are safe to run against a real, possibly-already-licensed
	// server. No functional/live calls are made for any of these methods.
	// =========================================================================
	Describe("Mutating operations", func() {
		Context("Parameter Validation", func() {
			Describe("ActivateOnline", func() {
				It("should fail with nil options", func() {
					resp, err := client.V2.Entitlements.ActivateOnline(context.Background(), nil)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("required"))
					Expect(resp).To(BeNil())
				})

				It("should fail without required license key", func() {
					resp, err := client.V2.Entitlements.ActivateOnline(context.Background(), &sonar.EntitlementsActivateOnlineOptions{})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("LicenseKey"))
					Expect(resp).To(BeNil())
				})
			})

			Describe("ActivateLegacy", func() {
				It("should fail with nil options", func() {
					resp, err := client.V2.Entitlements.ActivateLegacy(context.Background(), nil)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("required"))
					Expect(resp).To(BeNil())
				})

				It("should fail without required license key", func() {
					resp, err := client.V2.Entitlements.ActivateLegacy(context.Background(), &sonar.EntitlementsActivateLegacyOptions{})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("LicenseKey"))
					Expect(resp).To(BeNil())
				})
			})

			Describe("ActivateOffline", func() {
				It("should fail with nil options", func() {
					resp, err := client.V2.Entitlements.ActivateOffline(context.Background(), nil)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("required"))
					Expect(resp).To(BeNil())
				})

				It("should fail without required license content", func() {
					resp, err := client.V2.Entitlements.ActivateOffline(context.Background(), &sonar.EntitlementsActivateOfflineOptions{
						LicenseKey: "ABCD-EFGH-IJKL-MNOP",
					})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("License"))
					Expect(resp).To(BeNil())
				})

				It("should fail without required license key", func() {
					resp, err := client.V2.Entitlements.ActivateOffline(context.Background(), &sonar.EntitlementsActivateOfflineOptions{
						License: "some-license-file-content",
					})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("LicenseKey"))
					Expect(resp).To(BeNil())
				})
			})

			Describe("GetOfflineActivationRequest", func() {
				It("should fail with nil options", func() {
					result, resp, err := client.V2.Entitlements.GetOfflineActivationRequest(context.Background(), nil)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("required"))
					Expect(result).To(BeNil())
					Expect(resp).To(BeNil())
				})

				It("should fail without required license key", func() {
					result, resp, err := client.V2.Entitlements.GetOfflineActivationRequest(context.Background(), &sonar.EntitlementsGetOfflineActivationRequestOptions{})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("LicenseKey"))
					Expect(result).To(BeNil())
					Expect(resp).To(BeNil())
				})
			})
		})
	})
})
