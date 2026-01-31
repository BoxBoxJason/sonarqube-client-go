package integration_testing_test

import (
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Rules Service", Ordered, func() {
	var (
		client  *sonargo.Client
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
	// App
	// =========================================================================
	Describe("App", func() {
		It("should return rules application data", func() {
			result, resp, err := client.Rules.App()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Languages).NotTo(BeEmpty())
			Expect(result.Repositories).NotTo(BeEmpty())
		})

		It("should include canWrite flag", func() {
			result, resp, err := client.Rules.App()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			// CanWrite is a boolean indicating if the user can modify rules
			// Admin user should have write permission
			Expect(result.CanWrite).To(BeTrue())
		})
	})

	// =========================================================================
	// Repositories
	// =========================================================================
	Describe("Repositories", func() {
		It("should list all rule repositories", func() {
			result, resp, err := client.Rules.Repositories(&sonargo.RulesRepositoriesOption{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Repositories).NotTo(BeEmpty())
		})

		It("should filter by language", func() {
			result, resp, err := client.Rules.Repositories(&sonargo.RulesRepositoriesOption{
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			for _, repo := range result.Repositories {
				Expect(repo.Language).To(Equal("java"))
			}
		})

		It("should filter by query", func() {
			result, resp, err := client.Rules.Repositories(&sonargo.RulesRepositoriesOption{
				Query: "sonar",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			for _, repo := range result.Repositories {
				Expect(strings.Contains(strings.ToLower(repo.Key), "sonar") ||
					strings.Contains(strings.ToLower(repo.Name), "sonar")).To(BeTrue())
			}
		})

		It("should allow nil options", func() {
			result, resp, err := client.Rules.Repositories(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})
	})

	// =========================================================================
	// Tags
	// =========================================================================
	Describe("Tags", func() {
		It("should list all rule tags", func() {
			result, resp, err := client.Rules.Tags(&sonargo.RulesTagsOption{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Tags).NotTo(BeEmpty())
		})

		It("should limit by page size", func() {
			result, resp, err := client.Rules.Tags(&sonargo.RulesTagsOption{
				PageSize: 5,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(len(result.Tags)).To(BeNumerically("<=", 5))
		})

		It("should filter by query", func() {
			result, resp, err := client.Rules.Tags(&sonargo.RulesTagsOption{
				Query: "sec",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			// Tags containing "sec" should be returned (e.g., "security")
			for _, tag := range result.Tags {
				Expect(strings.Contains(strings.ToLower(tag), "sec")).To(BeTrue())
			}
		})

		It("should allow nil options", func() {
			result, resp, err := client.Rules.Tags(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		Context("parameter validation", func() {
			It("should fail with invalid page size (negative)", func() {
				_, resp, err := client.Rules.Tags(&sonargo.RulesTagsOption{
					PageSize: -1,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with page size exceeding max", func() {
				_, resp, err := client.Rules.Tags(&sonargo.RulesTagsOption{
					PageSize: 501,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Search
	// =========================================================================
	Describe("Search", func() {
		It("should search for all rules", func() {
			result, resp, err := client.Rules.Search(&sonargo.RulesSearchOption{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Rules).NotTo(BeEmpty())
			Expect(result.Paging.Total).To(BeNumerically(">", 0))
		})

		It("should filter by language", func() {
			result, resp, err := client.Rules.Search(&sonargo.RulesSearchOption{
				Languages: []string{"java"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			for _, rule := range result.Rules {
				Expect(rule.Lang).To(Equal("java"))
			}
		})

		It("should filter by repository", func() {
			// First get available repositories
			repos, _, err := client.Rules.Repositories(&sonargo.RulesRepositoriesOption{
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(repos.Repositories).NotTo(BeEmpty())

			repoKey := repos.Repositories[0].Key

			result, resp, err := client.Rules.Search(&sonargo.RulesSearchOption{
				Repositories: []string{repoKey},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			for _, rule := range result.Rules {
				Expect(rule.Repo).To(Equal(repoKey))
			}
		})

		It("should filter by severity", func() {
			result, resp, err := client.Rules.Search(&sonargo.RulesSearchOption{
				Severities: []string{"CRITICAL"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			for _, rule := range result.Rules {
				Expect(rule.Severity).To(Equal("CRITICAL"))
			}
		})

		It("should filter by type", func() {
			result, resp, err := client.Rules.Search(&sonargo.RulesSearchOption{
				Types: []string{"BUG"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			for _, rule := range result.Rules {
				Expect(rule.Type).To(Equal("BUG"))
			}
		})

		It("should filter by status", func() {
			result, resp, err := client.Rules.Search(&sonargo.RulesSearchOption{
				Statuses: []string{"READY"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			for _, rule := range result.Rules {
				Expect(rule.Status).To(Equal("READY"))
			}
		})

		It("should filter by tags", func() {
			// First get available tags
			tags, _, err := client.Rules.Tags(&sonargo.RulesTagsOption{
				PageSize: 1,
			})
			Expect(err).NotTo(HaveOccurred())
			if len(tags.Tags) > 0 {
				result, resp, err := client.Rules.Search(&sonargo.RulesSearchOption{
					Tags: []string{tags.Tags[0]},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				// Results may still be empty if no rules have this tag
			}
		})

		It("should filter template rules", func() {
			result, resp, err := client.Rules.Search(&sonargo.RulesSearchOption{
				IsTemplate: true,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			for _, rule := range result.Rules {
				Expect(rule.IsTemplate).To(BeTrue())
			}
		})

		It("should paginate results", func() {
			result1, resp, err := client.Rules.Search(&sonargo.RulesSearchOption{
				PaginationArgs: sonargo.PaginationArgs{
					Page:     1,
					PageSize: 5,
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(len(result1.Rules)).To(BeNumerically("<=", 5))

			if result1.Paging.Total > 5 {
				result2, resp, err := client.Rules.Search(&sonargo.RulesSearchOption{
					PaginationArgs: sonargo.PaginationArgs{
						Page:     2,
						PageSize: 5,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				// Different pages should have different rules
				if len(result1.Rules) > 0 && len(result2.Rules) > 0 {
					Expect(result1.Rules[0].Key).NotTo(Equal(result2.Rules[0].Key))
				}
			}
		})

		It("should include facets", func() {
			result, resp, err := client.Rules.Search(&sonargo.RulesSearchOption{
				Facets: []string{"languages", "repositories", "tags"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result.Facets).NotTo(BeEmpty())
		})

		It("should search by query", func() {
			result, resp, err := client.Rules.Search(&sonargo.RulesSearchOption{
				Query: "null",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		It("should include external rules", func() {
			result, resp, err := client.Rules.Search(&sonargo.RulesSearchOption{
				IncludeExternal: true,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		It("should allow nil options", func() {
			result, resp, err := client.Rules.Search(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		Context("parameter validation", func() {
			It("should fail with search query too short", func() {
				_, resp, err := client.Rules.Search(&sonargo.RulesSearchOption{
					Query: "a",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Show
	// =========================================================================
	Describe("Show", func() {
		var existingRuleKey string

		BeforeAll(func() {
			// Get an existing rule key for show tests
			result, _, err := client.Rules.Search(&sonargo.RulesSearchOption{
				PaginationArgs: sonargo.PaginationArgs{
					PageSize: 1,
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Rules).NotTo(BeEmpty())
			existingRuleKey = result.Rules[0].Key
		})

		It("should show rule details", func() {
			result, resp, err := client.Rules.Show(&sonargo.RulesShowOption{
				Key: existingRuleKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Rule.Key).To(Equal(existingRuleKey))
			Expect(result.Rule.Name).NotTo(BeEmpty())
		})

		It("should include actives when requested", func() {
			result, resp, err := client.Rules.Show(&sonargo.RulesShowOption{
				Key:     existingRuleKey,
				Actives: true,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			// Actives may be empty if rule is not activated in any profile
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.Rules.Show(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing key", func() {
				_, resp, err := client.Rules.Show(&sonargo.RulesShowOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with non-existent rule", func() {
				_, resp, err := client.Rules.Show(&sonargo.RulesShowOption{
					Key: "nonexistent:rule-key",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	// =========================================================================
	// Create, Update, Delete (Custom Rule Lifecycle)
	// =========================================================================
	Describe("Custom Rule Lifecycle", func() {
		var templateRule sonargo.RuleDetails

		BeforeAll(func() {
			// Find a template rule to use for creating custom rules
			result, _, err := client.Rules.Search(&sonargo.RulesSearchOption{
				IsTemplate: true,
				Languages:  []string{"java"},
				PaginationArgs: sonargo.PaginationArgs{
					PageSize: 1,
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Rules).NotTo(BeEmpty(), "No template rules found for Java")
			templateRule = result.Rules[0]
		})

		// =====================================================================
		// Create
		// =====================================================================
		Describe("Create", func() {
			It("should create a custom rule from template", func() {
				customKey := helpers.UniqueResourceName("rule")

				result, resp, err := client.Rules.Create(&sonargo.RulesCreateOption{
					CustomKey:           customKey,
					Name:                "E2E Test Rule " + customKey,
					MarkdownDescription: "This is a test rule created by e2e tests",
					TemplateKey:         templateRule.Key,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Rule.Name).To(Equal("E2E Test Rule " + customKey))
				Expect(result.Rule.TemplateKey).To(Equal(templateRule.Key))

				// Register cleanup
				cleanup.RegisterCleanup("rule", result.Rule.Key, func() error {
					_, err := client.Rules.Delete(&sonargo.RulesDeleteOption{
						Key: result.Rule.Key,
					})
					return err
				})
			})

			It("should create a custom rule with severity", func() {
				customKey := helpers.UniqueResourceName("rule")

				result, resp, err := client.Rules.Create(&sonargo.RulesCreateOption{
					CustomKey:           customKey,
					Name:                "E2E Severity Rule " + customKey,
					MarkdownDescription: "Rule with custom severity",
					TemplateKey:         templateRule.Key,
					Severity:            "CRITICAL",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result.Rule.Severity).To(Equal("CRITICAL"))

				cleanup.RegisterCleanup("rule", result.Rule.Key, func() error {
					_, err := client.Rules.Delete(&sonargo.RulesDeleteOption{
						Key: result.Rule.Key,
					})
					return err
				})
			})

			It("should create a custom rule with status", func() {
				customKey := helpers.UniqueResourceName("rule")

				result, resp, err := client.Rules.Create(&sonargo.RulesCreateOption{
					CustomKey:           customKey,
					Name:                "E2E Status Rule " + customKey,
					MarkdownDescription: "Rule with custom status",
					TemplateKey:         templateRule.Key,
					Status:              "BETA",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result.Rule.Status).To(Equal("BETA"))

				cleanup.RegisterCleanup("rule", result.Rule.Key, func() error {
					_, err := client.Rules.Delete(&sonargo.RulesDeleteOption{
						Key: result.Rule.Key,
					})
					return err
				})
			})

			Context("parameter validation", func() {
				It("should fail with nil options", func() {
					_, resp, err := client.Rules.Create(nil)
					Expect(err).To(HaveOccurred())
					Expect(resp).To(BeNil())
				})

				It("should fail with missing custom key", func() {
					_, resp, err := client.Rules.Create(&sonargo.RulesCreateOption{
						Name:                "Test Rule",
						MarkdownDescription: "Description",
						TemplateKey:         templateRule.Key,
					})
					Expect(err).To(HaveOccurred())
					Expect(resp).To(BeNil())
				})

				It("should fail with missing name", func() {
					_, resp, err := client.Rules.Create(&sonargo.RulesCreateOption{
						CustomKey:           "test-rule",
						MarkdownDescription: "Description",
						TemplateKey:         templateRule.Key,
					})
					Expect(err).To(HaveOccurred())
					Expect(resp).To(BeNil())
				})

				It("should fail with missing markdown description", func() {
					_, resp, err := client.Rules.Create(&sonargo.RulesCreateOption{
						CustomKey:   "test-rule",
						Name:        "Test Rule",
						TemplateKey: templateRule.Key,
					})
					Expect(err).To(HaveOccurred())
					Expect(resp).To(BeNil())
				})

				It("should fail with missing template key", func() {
					_, resp, err := client.Rules.Create(&sonargo.RulesCreateOption{
						CustomKey:           "test-rule",
						Name:                "Test Rule",
						MarkdownDescription: "Description",
					})
					Expect(err).To(HaveOccurred())
					Expect(resp).To(BeNil())
				})

				It("should fail with custom key too long", func() {
					longKey := strings.Repeat("a", 201)
					_, resp, err := client.Rules.Create(&sonargo.RulesCreateOption{
						CustomKey:           longKey,
						Name:                "Test Rule",
						MarkdownDescription: "Description",
						TemplateKey:         templateRule.Key,
					})
					Expect(err).To(HaveOccurred())
					Expect(resp).To(BeNil())
				})

				It("should fail with name too long", func() {
					longName := strings.Repeat("a", 201)
					_, resp, err := client.Rules.Create(&sonargo.RulesCreateOption{
						CustomKey:           "test-rule",
						Name:                longName,
						MarkdownDescription: "Description",
						TemplateKey:         templateRule.Key,
					})
					Expect(err).To(HaveOccurred())
					Expect(resp).To(BeNil())
				})
			})
		})

		// =====================================================================
		// Update
		// =====================================================================
		Describe("Update", func() {
			It("should update a custom rule name", func() {
				customKey := helpers.UniqueResourceName("rule")

				// Create rule
				createResult, _, err := client.Rules.Create(&sonargo.RulesCreateOption{
					CustomKey:           customKey,
					Name:                "Original Name",
					MarkdownDescription: "Original description",
					TemplateKey:         templateRule.Key,
				})
				Expect(err).NotTo(HaveOccurred())

				cleanup.RegisterCleanup("rule", createResult.Rule.Key, func() error {
					_, err := client.Rules.Delete(&sonargo.RulesDeleteOption{
						Key: createResult.Rule.Key,
					})
					return err
				})

				// Update rule
				updatedName := "Updated Name " + customKey
				updateResult, resp, err := client.Rules.Update(&sonargo.RulesUpdateOption{
					Key:                 createResult.Rule.Key,
					Name:                updatedName,
					MarkdownDescription: "Updated description",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(updateResult.Rule.Name).To(Equal(updatedName))
			})

			It("should update a custom rule severity", func() {
				customKey := helpers.UniqueResourceName("rule")

				createResult, _, err := client.Rules.Create(&sonargo.RulesCreateOption{
					CustomKey:           customKey,
					Name:                "Severity Update Rule",
					MarkdownDescription: "Testing severity update",
					TemplateKey:         templateRule.Key,
					Severity:            "MINOR",
				})
				Expect(err).NotTo(HaveOccurred())

				cleanup.RegisterCleanup("rule", createResult.Rule.Key, func() error {
					_, err := client.Rules.Delete(&sonargo.RulesDeleteOption{
						Key: createResult.Rule.Key,
					})
					return err
				})

				updateResult, resp, err := client.Rules.Update(&sonargo.RulesUpdateOption{
					Key:      createResult.Rule.Key,
					Severity: "BLOCKER",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(updateResult.Rule.Severity).To(Equal("BLOCKER"))
			})

			It("should update a custom rule status", func() {
				customKey := helpers.UniqueResourceName("rule")

				createResult, _, err := client.Rules.Create(&sonargo.RulesCreateOption{
					CustomKey:           customKey,
					Name:                "Status Update Rule",
					MarkdownDescription: "Testing status update",
					TemplateKey:         templateRule.Key,
					Status:              "READY",
				})
				Expect(err).NotTo(HaveOccurred())

				cleanup.RegisterCleanup("rule", createResult.Rule.Key, func() error {
					_, err := client.Rules.Delete(&sonargo.RulesDeleteOption{
						Key: createResult.Rule.Key,
					})
					return err
				})

				updateResult, resp, err := client.Rules.Update(&sonargo.RulesUpdateOption{
					Key:    createResult.Rule.Key,
					Status: "DEPRECATED",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(updateResult.Rule.Status).To(Equal("DEPRECATED"))
			})

			It("should update tags on a rule", func() {
				// Get an existing rule key
				result, _, err := client.Rules.Search(&sonargo.RulesSearchOption{
					Languages: []string{"java"},
					PaginationArgs: sonargo.PaginationArgs{
						PageSize: 1,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Rules).NotTo(BeEmpty())
				ruleKey := result.Rules[0].Key

				// Add a custom tag
				updateResult, resp, err := client.Rules.Update(&sonargo.RulesUpdateOption{
					Key:  ruleKey,
					Tags: []string{"e2e-test-tag"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(updateResult.Rule.Tags).To(ContainElement("e2e-test-tag"))

				// Clean up - remove the tag
				_, _, _ = client.Rules.Update(&sonargo.RulesUpdateOption{
					Key:  ruleKey,
					Tags: []string{},
				})
			})

			It("should add a note to a rule", func() {
				// Get an existing rule key
				result, _, err := client.Rules.Search(&sonargo.RulesSearchOption{
					Languages: []string{"java"},
					PaginationArgs: sonargo.PaginationArgs{
						PageSize: 1,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Rules).NotTo(BeEmpty())
				ruleKey := result.Rules[0].Key

				// Add a note
				updateResult, resp, err := client.Rules.Update(&sonargo.RulesUpdateOption{
					Key:          ruleKey,
					MarkdownNote: "E2E test note: This is a test note added during e2e testing",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(updateResult.Rule.MdNote).To(ContainSubstring("E2E test note"))

				// Clean up - remove the note
				_, _, _ = client.Rules.Update(&sonargo.RulesUpdateOption{
					Key:          ruleKey,
					MarkdownNote: "",
				})
			})

			Context("parameter validation", func() {
				It("should fail with nil options", func() {
					_, resp, err := client.Rules.Update(nil)
					Expect(err).To(HaveOccurred())
					Expect(resp).To(BeNil())
				})

				It("should fail with missing key", func() {
					_, resp, err := client.Rules.Update(&sonargo.RulesUpdateOption{
						Name: "Updated Name",
					})
					Expect(err).To(HaveOccurred())
					Expect(resp).To(BeNil())
				})

				It("should fail with key too long", func() {
					longKey := strings.Repeat("a", 201)
					_, resp, err := client.Rules.Update(&sonargo.RulesUpdateOption{
						Key:  longKey,
						Name: "Updated Name",
					})
					Expect(err).To(HaveOccurred())
					Expect(resp).To(BeNil())
				})
			})
		})

		// =====================================================================
		// Delete
		// =====================================================================
		Describe("Delete", func() {
			It("should delete a custom rule", func() {
				customKey := helpers.UniqueResourceName("rule")

				// Create rule
				createResult, _, err := client.Rules.Create(&sonargo.RulesCreateOption{
					CustomKey:           customKey,
					Name:                "Rule To Delete",
					MarkdownDescription: "This rule will be deleted",
					TemplateKey:         templateRule.Key,
				})
				Expect(err).NotTo(HaveOccurred())

				// Delete rule
				resp, err := client.Rules.Delete(&sonargo.RulesDeleteOption{
					Key: createResult.Rule.Key,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				// Verify deletion - rule is marked as REMOVED
				showResult, _, err := client.Rules.Show(&sonargo.RulesShowOption{
					Key: createResult.Rule.Key,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(showResult.Rule.Status).To(Equal("REMOVED"))
			})

			Context("parameter validation", func() {
				It("should fail with nil options", func() {
					resp, err := client.Rules.Delete(nil)
					Expect(err).To(HaveOccurred())
					Expect(resp).To(BeNil())
				})

				It("should fail with missing key", func() {
					resp, err := client.Rules.Delete(&sonargo.RulesDeleteOption{})
					Expect(err).To(HaveOccurred())
					Expect(resp).To(BeNil())
				})
			})
		})
	})

	// =========================================================================
	// List
	// =========================================================================
	Describe("List", func() {
		It("should list rules", func() {
			result, resp, err := client.Rules.List(&sonargo.RulesListOption{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			// List returns protobuf format, so the response is a string
		})

		It("should paginate results", func() {
			result, resp, err := client.Rules.List(&sonargo.RulesListOption{
				PaginationArgs: sonargo.PaginationArgs{
					PageSize: 10,
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		It("should allow nil options", func() {
			result, resp, err := client.Rules.List(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})
	})
})
