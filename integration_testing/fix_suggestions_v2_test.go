package integration_testing_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("Fix Suggestions V2 Service", Ordered, func() {
	var client *sonar.Client

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	Describe("CreateSuggestion", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.FixSuggestions.CreateSuggestion(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})
		Context("Functional Tests", func() {
			It("should create suggestion or return an expected error", func() {
				result, resp, err := client.V2.FixSuggestions.CreateSuggestion(context.Background(), &sonar.FixSuggestionsCreateOptions{
					IssueId: "nonexistent-issue",
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

	Describe("GetEnablement", func() {
		Context("Functional Tests", func() {
			It("should return enablement or an expected error", func() {
				result, resp, err := client.V2.FixSuggestions.GetEnablement(context.Background())
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})

	Describe("SetEnablement", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.V2.FixSuggestions.SetEnablement(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with an invalid enablement value", func() {
				resp, err := client.V2.FixSuggestions.SetEnablement(context.Background(), &sonar.FixSuggestionsSetEnablementOptions{
					Enablement: "NOT_A_VALID_STATE",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
		Context("Functional Tests", func() {
			It("should set enablement or return an expected error", func() {
				resp, err := client.V2.FixSuggestions.SetEnablement(context.Background(), &sonar.FixSuggestionsSetEnablementOptions{
					Enablement: sonar.FixSuggestionsEnablementDisabled,
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
				}
			})
		})
	})

	Describe("AwarenessBannerInteraction", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.FixSuggestions.AwarenessBannerInteraction(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail with an invalid banner type", func() {
				result, resp, err := client.V2.FixSuggestions.AwarenessBannerInteraction(context.Background(), &sonar.FixSuggestionsAwarenessBannerOptions{
					BannerType: "NOT_A_VALID_TYPE",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})
		Context("Functional Tests", func() {
			It("should record the interaction or return an expected error", func() {
				result, resp, err := client.V2.FixSuggestions.AwarenessBannerInteraction(context.Background(), &sonar.FixSuggestionsAwarenessBannerOptions{
					BannerType: sonar.FixSuggestionsBannerTypeEnable,
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

	Describe("GetIssueAvailability", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.FixSuggestions.GetIssueAvailability(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})
		Context("Functional Tests", func() {
			It("should return availability or an expected error", func() {
				result, resp, err := client.V2.FixSuggestions.GetIssueAvailability(context.Background(), &sonar.FixSuggestionsIssueOptions{
					IssueId: "nonexistent-issue",
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

	Describe("GetServiceInfo", func() {
		Context("Functional Tests", func() {
			It("should return service info or an expected error", func() {
				result, resp, err := client.V2.FixSuggestions.GetServiceInfo(context.Background())
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})

	Describe("GetSubscriptionType", func() {
		Context("Functional Tests", func() {
			It("should return subscription type or an expected error", func() {
				result, resp, err := client.V2.FixSuggestions.GetSubscriptionType(context.Background())
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})

	Describe("GetSupportedLlmProviders", func() {
		Context("Functional Tests", func() {
			It("should return providers or an expected error", func() {
				result, resp, err := client.V2.FixSuggestions.GetSupportedLlmProviders(context.Background())
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})
})
