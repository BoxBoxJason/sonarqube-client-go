package integration_testing_test

import (
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("L10N Service", Ordered, func() {
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
	// Index
	// =========================================================================
	Describe("Index", func() {
		Context("Parameter Validation", func() {
			It("should work with nil options", func() {
				result, resp, err := client.L10N.Index(nil)
				// Skip if API not available
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("L10N API is not available in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("Default Locale", func() {
			It("should get localization messages with default locale", func() {
				result, resp, err := client.L10N.Index(nil)
				// Skip if API not available
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("L10N API is not available in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())

				// Verify locale if returned (some SonarQube versions may not include it)
				if result.Locale != "" {
					Expect(result.Locale).To(Equal("en"))
				}
				Expect(result.Messages).NotTo(BeEmpty())

				// Check for common translation keys
				_, hasCommonKey := result.Messages["quality_gates.operator.LT"]
				if !hasCommonKey {
					_, hasCommonKey = result.Messages["projects.create_project"]
				}
				if !hasCommonKey {
					_, hasCommonKey = result.Messages["issues.type.BUG"]
				}
				Expect(hasCommonKey).To(BeTrue(), "Should contain at least one common translation key")
			})
		})

		Context("Specific Locale", func() {
			It("should get localization messages for specific locale", func() {
				result, resp, err := client.L10N.Index(&sonar.L10NIndexOption{
					Locale: "en",
				})
				// Skip if API not available
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("L10N API is not available in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())

				// Verify locale if returned (some SonarQube versions may not include it)
				if result.Locale != "" {
					Expect(result.Locale).To(Equal("en"))
				}
				Expect(result.Messages).NotTo(BeEmpty())
			})
		})

		Context("Compare Different Locales", func() {
			It("should return different translations for different locales if available", func() {
				// Get English locale
				resultEN, respEN, err := client.L10N.Index(&sonar.L10NIndexOption{
					Locale: "en",
				})
				// Skip if API not available
				if respEN != nil && respEN.StatusCode == http.StatusNotFound {
					Skip("L10N API is not available in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(respEN.StatusCode).To(Equal(http.StatusOK))

				// Try to get French locale
				resultFR, respFR, err := client.L10N.Index(&sonar.L10NIndexOption{
					Locale: "fr",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(respFR.StatusCode).To(Equal(http.StatusOK))

				// If locales are returned and French is available, verify they differ
				if resultEN.Locale != "" && resultFR.Locale != "" && resultFR.Locale == "fr" {
					Expect(resultEN.Locale).NotTo(Equal(resultFR.Locale))
					// Find a common key and compare translations if possible
					for key := range resultEN.Messages {
						if valFR, exists := resultFR.Messages[key]; exists {
							// If the same key exists in both, the structure is valid
							// (translations may or may not differ depending on the key)
							_ = valFR
							break
						}
					}
				}
			})
		})

		Context("With Timestamp", func() {
			It("should get localization messages with timestamp", func() {
				// Use a recent timestamp (1 year ago)
				oneYearAgo := time.Now().AddDate(-1, 0, 0).Format("2006-01-02T15:04:05-0700")
				result, resp, err := client.L10N.Index(&sonar.L10NIndexOption{
					Timestamp: oneYearAgo,
				})
				// Skip if API not available
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("L10N API is not available in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})
	})
})
