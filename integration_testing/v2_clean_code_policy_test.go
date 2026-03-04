package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("V2 Clean Code Policy Service", Ordered, func() {
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
		errors := cleanup.Cleanup()
		for _, err := range errors {
			GinkgoWriter.Printf("Cleanup error: %v\n", err)
		}
	})

	// =========================================================================
	// CreateRule
	// =========================================================================
	Describe("CreateRule", func() {
		Context("parameter validation", func() {
			It("should fail with nil request", func() {
				result, resp, err := client.V2.CleanCodePolicy.CreateRule(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing key", func() {
				result, resp, err := client.V2.CleanCodePolicy.CreateRule(&sonar.CleanCodePolicyCreateRuleOptions{
					TemplateKey:         "java:S100",
					Name:                "Test Rule",
					MarkdownDescription: "Test description",
					Impacts:             []sonar.RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: "HIGH"}},
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing template key", func() {
				result, resp, err := client.V2.CleanCodePolicy.CreateRule(&sonar.CleanCodePolicyCreateRuleOptions{
					Key:                 "custom:test_rule",
					Name:                "Test Rule",
					MarkdownDescription: "Test description",
					Impacts:             []sonar.RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: "HIGH"}},
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing name", func() {
				result, resp, err := client.V2.CleanCodePolicy.CreateRule(&sonar.CleanCodePolicyCreateRuleOptions{
					Key:                 "custom:test_rule",
					TemplateKey:         "java:S100",
					MarkdownDescription: "Test description",
					Impacts:             []sonar.RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: "HIGH"}},
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing markdown description", func() {
				result, resp, err := client.V2.CleanCodePolicy.CreateRule(&sonar.CleanCodePolicyCreateRuleOptions{
					Key:         "custom:test_rule",
					TemplateKey: "java:S100",
					Name:        "Test Rule",
					Impacts:     []sonar.RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: "HIGH"}},
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with empty impacts", func() {
				result, resp, err := client.V2.CleanCodePolicy.CreateRule(&sonar.CleanCodePolicyCreateRuleOptions{
					Key:                 "custom:test_rule",
					TemplateKey:         "java:S100",
					Name:                "Test Rule",
					MarkdownDescription: "Test description",
					Impacts:             []sonar.RuleImpact{},
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with invalid clean code attribute", func() {
				result, resp, err := client.V2.CleanCodePolicy.CreateRule(&sonar.CleanCodePolicyCreateRuleOptions{
					Key:                 "custom:test_rule",
					TemplateKey:         "java:S100",
					Name:                "Test Rule",
					MarkdownDescription: "Test description",
					Impacts:             []sonar.RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: "HIGH"}},
					CleanCodeAttribute:  "INVALID",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with invalid impact severity", func() {
				result, resp, err := client.V2.CleanCodePolicy.CreateRule(&sonar.CleanCodePolicyCreateRuleOptions{
					Key:                 "custom:test_rule",
					TemplateKey:         "java:S100",
					Name:                "Test Rule",
					MarkdownDescription: "Test description",
					Impacts:             []sonar.RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: "INVALID"}},
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Context("successful creation", func() {
			var templateRule sonar.RuleDetails

			BeforeAll(func() {
				// Find a template rule to use for creating custom rules.
				result, _, err := client.Rules.Search(&sonar.RulesSearchOption{
					IsTemplate: true,
					Languages:  []string{"java"},
					PaginationArgs: sonar.PaginationArgs{
						PageSize: 1,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Rules).NotTo(BeEmpty(), "No template rules found for Java")
				templateRule = result.Rules[0]
			})

			It("should create a custom rule from a template", func() {
				customKey := helpers.UniqueResourceName("rule")
				ruleKey := templateRule.Repo + ":" + customKey

				result, resp, err := client.V2.CleanCodePolicy.CreateRule(&sonar.CleanCodePolicyCreateRuleOptions{
					Key:                 ruleKey,
					TemplateKey:         templateRule.Key,
					Name:                "E2E V2 Test Rule " + customKey,
					MarkdownDescription: "This is a test rule created by V2 e2e tests",
					Impacts:             []sonar.RuleImpact{{SoftwareQuality: "MAINTAINABILITY", Severity: "HIGH"}},
					Status:              "READY",
					CleanCodeAttribute:  "CONVENTIONAL",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Name).To(Equal("E2E V2 Test Rule " + customKey))
				Expect(result.Key).To(Equal(ruleKey))
				Expect(result.Status).To(Equal("READY"))
				Expect(result.CleanCodeAttribute).To(Equal("CONVENTIONAL"))
				Expect(result.Impacts).NotTo(BeEmpty())

				// Register cleanup using V1 Rules.Delete.
				cleanup.RegisterCleanup("rule", ruleKey, func() error {
					_, err := client.Rules.Delete(&sonar.RulesDeleteOption{
						Key: ruleKey,
					})
					return err
				})
			})
		})
	})
})
