package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

const (
	nonExistentProjectKey = "non-existent-project"
	nonExistentIssueKey   = "AXxxxxxxxxxxxxxxxxxx"
)

var _ = Describe("Issues Service", Ordered, func() {
	var (
		client     *sonargo.Client
		cleanup    *helpers.CleanupManager
		projectKey string
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
		cleanup = helpers.NewCleanupManager(client)

		// Create a test project for issues-related operations
		projectKey = helpers.UniqueResourceName("iss")
		_, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
			Name:    "Issues Test Project",
			Project: projectKey,
		})
		Expect(err).NotTo(HaveOccurred())

		cleanup.RegisterCleanup("project", projectKey, func() error {
			_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
				Project: projectKey,
			})
			return err
		})
	})

	AfterAll(func() {
		errors := cleanup.Cleanup()
		for _, err := range errors {
			GinkgoWriter.Printf("Cleanup error: %v\n", err)
		}
	})

	// =========================================================================
	// Search
	// =========================================================================
	Describe("Search", func() {
		Context("Valid Requests", func() {
			It("should search issues with nil options", func() {
				result, resp, err := client.Issues.Search(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should search issues for a project", func() {
				result, resp, err := client.Issues.Search(&sonargo.IssuesSearchOption{
					Projects: []string{projectKey},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should search issues with pagination", func() {
				result, resp, err := client.Issues.Search(&sonargo.IssuesSearchOption{
					Projects: []string{projectKey},
					PaginationArgs: sonargo.PaginationArgs{
						PageSize: 10,
						Page:     1,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should search issues with impact severities filter", func() {
				result, resp, err := client.Issues.Search(&sonargo.IssuesSearchOption{
					Projects:         []string{projectKey},
					ImpactSeverities: []string{"HIGH", "MEDIUM"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should search issues with issue statuses filter", func() {
				result, resp, err := client.Issues.Search(&sonargo.IssuesSearchOption{
					Projects:      []string{projectKey},
					IssueStatuses: []string{"OPEN", "CONFIRMED"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should search issues with clean code categories filter", func() {
				result, resp, err := client.Issues.Search(&sonargo.IssuesSearchOption{
					Projects:                     []string{projectKey},
					CleanCodeAttributeCategories: []string{"INTENTIONAL", "CONSISTENT"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should search issues with additional fields", func() {
				result, resp, err := client.Issues.Search(&sonargo.IssuesSearchOption{
					Projects:         []string{projectKey},
					AdditionalFields: []string{"rules", "users"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should search issues with facets", func() {
				result, resp, err := client.Issues.Search(&sonargo.IssuesSearchOption{
					Projects: []string{projectKey},
					Facets:   []string{"impactSeverities", "issueStatuses"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("Non-Existent Project", func() {
			It("should return empty results for non-existent project", func() {
				result, resp, err := client.Issues.Search(&sonargo.IssuesSearchOption{
					Projects: []string{nonExistentProjectKey},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Issues).To(BeEmpty())
			})
		})
	})

	// =========================================================================
	// List
	// =========================================================================
	Describe("List", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Issues.List(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot be nil"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without project or component", func() {
				result, resp, err := client.Issues.List(&sonargo.IssuesListOption{})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Valid Requests", func() {
			It("should list issues for a project", func() {
				result, resp, err := client.Issues.List(&sonargo.IssuesListOption{
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should list issues with pagination", func() {
				result, resp, err := client.Issues.List(&sonargo.IssuesListOption{
					Project: projectKey,
					PaginationArgs: sonargo.PaginationArgs{
						PageSize: 10,
						Page:     1,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})
	})

	// =========================================================================
	// Authors
	// =========================================================================
	Describe("Authors", func() {
		Context("Valid Requests", func() {
			It("should list authors with nil options", func() {
				result, resp, err := client.Issues.Authors(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should list authors for a project", func() {
				result, resp, err := client.Issues.Authors(&sonargo.IssuesAuthorsOption{
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should list authors with query filter", func() {
				result, resp, err := client.Issues.Authors(&sonargo.IssuesAuthorsOption{
					Query: "admin",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should list authors with pagination", func() {
				result, resp, err := client.Issues.Authors(&sonargo.IssuesAuthorsOption{
					PageSize: 10,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})
	})

	// =========================================================================
	// Tags
	// =========================================================================
	Describe("Tags", func() {
		Context("Valid Requests", func() {
			It("should list tags with nil options", func() {
				result, resp, err := client.Issues.Tags(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should list tags for a project", func() {
				result, resp, err := client.Issues.Tags(&sonargo.IssuesTagsOption{
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should list tags with query filter", func() {
				result, resp, err := client.Issues.Tags(&sonargo.IssuesTagsOption{
					Query: "security",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})
	})

	// =========================================================================
	// AddComment (requires existing issue)
	// =========================================================================
	Describe("AddComment", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Issues.AddComment(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot be nil"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required issue key", func() {
				result, resp, err := client.Issues.AddComment(&sonargo.IssuesAddCommentOption{
					Text: "Test comment",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Issue"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required text", func() {
				result, resp, err := client.Issues.AddComment(&sonargo.IssuesAddCommentOption{
					Issue: nonExistentIssueKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Text"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent Issue", func() {
			It("should fail with non-existent issue key", func() {
				result, resp, err := client.Issues.AddComment(&sonargo.IssuesAddCommentOption{
					Issue: nonExistentIssueKey,
					Text:  "Test comment",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// Assign (requires existing issue)
	// =========================================================================
	Describe("Assign", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Issues.Assign(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot be nil"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required issue key", func() {
				result, resp, err := client.Issues.Assign(&sonargo.IssuesAssignOption{
					Assignee: "admin",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Issue"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent Issue", func() {
			It("should fail with non-existent issue key", func() {
				result, resp, err := client.Issues.Assign(&sonargo.IssuesAssignOption{
					Issue:    nonExistentIssueKey,
					Assignee: "admin",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// Changelog (requires existing issue)
	// =========================================================================
	Describe("Changelog", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Issues.Changelog(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot be nil"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required issue key", func() {
				result, resp, err := client.Issues.Changelog(&sonargo.IssuesChangelogOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Issue"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent Issue", func() {
			It("should fail with non-existent issue key", func() {
				result, resp, err := client.Issues.Changelog(&sonargo.IssuesChangelogOption{
					Issue: nonExistentIssueKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// DoTransition (requires existing issue)
	// =========================================================================
	Describe("DoTransition", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Issues.DoTransition(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot be nil"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required issue key", func() {
				result, resp, err := client.Issues.DoTransition(&sonargo.IssuesDoTransitionOption{
					Transition: "confirm",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Issue"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required transition", func() {
				result, resp, err := client.Issues.DoTransition(&sonargo.IssuesDoTransitionOption{
					Issue: nonExistentIssueKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Transition"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent Issue", func() {
			It("should fail with non-existent issue key", func() {
				result, resp, err := client.Issues.DoTransition(&sonargo.IssuesDoTransitionOption{
					Issue:      nonExistentIssueKey,
					Transition: "confirm",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// BulkChange (requires existing issues)
	// =========================================================================
	Describe("BulkChange", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Issues.BulkChange(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot be nil"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required issues", func() {
				result, resp, err := client.Issues.BulkChange(&sonargo.IssuesBulkChangeOption{
					AddTags: []string{"test-tag"},
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Issues"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent Issues", func() {
			It("should handle bulk change on non-existent issues", func() {
				result, resp, err := client.Issues.BulkChange(&sonargo.IssuesBulkChangeOption{
					Issues:  []string{nonExistentIssueKey},
					AddTags: []string{"test-tag"},
				})
				// API may succeed with 0 failures or may fail
				if err == nil {
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					Expect(result).NotTo(BeNil())
				} else {
					Expect(result).To(BeNil())
				}
			})
		})
	})

	// =========================================================================
	// SetSeverity (requires existing issue)
	// =========================================================================
	Describe("SetSeverity", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Issues.SetSeverity(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot be nil"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required issue key", func() {
				result, resp, err := client.Issues.SetSeverity(&sonargo.IssuesSetSeverityOption{
					Severity: "MAJOR",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Issue"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent Issue", func() {
			It("should fail with non-existent issue key", func() {
				result, resp, err := client.Issues.SetSeverity(&sonargo.IssuesSetSeverityOption{
					Issue:    nonExistentIssueKey,
					Severity: "MAJOR",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// SetTags (requires existing issue)
	// =========================================================================
	Describe("SetTags", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Issues.SetTags(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot be nil"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required issue key", func() {
				result, resp, err := client.Issues.SetTags(&sonargo.IssuesSetTagsOption{
					Tags: []string{"test-tag"},
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Issue"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent Issue", func() {
			It("should fail with non-existent issue key", func() {
				result, resp, err := client.Issues.SetTags(&sonargo.IssuesSetTagsOption{
					Issue: nonExistentIssueKey,
					Tags:  []string{"test-tag"},
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// SetType (requires existing issue)
	// =========================================================================
	Describe("SetType", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Issues.SetType(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot be nil"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required issue key", func() {
				result, resp, err := client.Issues.SetType(&sonargo.IssuesSetTypeOption{
					Type: "BUG",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Issue"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required type", func() {
				result, resp, err := client.Issues.SetType(&sonargo.IssuesSetTypeOption{
					Issue: nonExistentIssueKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Type"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent Issue", func() {
			It("should fail with non-existent issue key", func() {
				result, resp, err := client.Issues.SetType(&sonargo.IssuesSetTypeOption{
					Issue: nonExistentIssueKey,
					Type:  "BUG",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// DeleteComment (requires existing comment)
	// =========================================================================
	Describe("DeleteComment", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Issues.DeleteComment(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot be nil"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required comment key", func() {
				result, resp, err := client.Issues.DeleteComment(&sonargo.IssuesDeleteCommentOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Comment"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent Comment", func() {
			It("should fail with non-existent comment key", func() {
				result, resp, err := client.Issues.DeleteComment(&sonargo.IssuesDeleteCommentOption{
					Comment: nonExistentIssueKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// EditComment (requires existing comment)
	// =========================================================================
	Describe("EditComment", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Issues.EditComment(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot be nil"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required comment key", func() {
				result, resp, err := client.Issues.EditComment(&sonargo.IssuesEditCommentOption{
					Text: "Updated comment",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Comment"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required text", func() {
				result, resp, err := client.Issues.EditComment(&sonargo.IssuesEditCommentOption{
					Comment: nonExistentIssueKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Text"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent Comment", func() {
			It("should fail with non-existent comment key", func() {
				result, resp, err := client.Issues.EditComment(&sonargo.IssuesEditCommentOption{
					Comment: nonExistentIssueKey,
					Text:    "Updated comment",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// Reindex
	// =========================================================================
	Describe("Reindex", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Issues.Reindex(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot be nil"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required project key", func() {
				resp, err := client.Issues.Reindex(&sonargo.IssuesReindexOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Project"))
				Expect(resp).To(BeNil())
			})
		})

		Context("Valid Requests", func() {
			It("should trigger reindex for a project", func() {
				resp, err := client.Issues.Reindex(&sonargo.IssuesReindexOption{
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})
		})

		Context("Non-Existent Project", func() {
			It("should fail for non-existent project", func() {
				resp, err := client.Issues.Reindex(&sonargo.IssuesReindexOption{
					Project: nonExistentProjectKey,
				})
				Expect(err).To(HaveOccurred())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// Pull
	// =========================================================================
	Describe("Pull", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Issues.Pull(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot be nil"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required project key", func() {
				result, resp, err := client.Issues.Pull(&sonargo.IssuesPullOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ProjectKey"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required branch name", func() {
				result, resp, err := client.Issues.Pull(&sonargo.IssuesPullOption{
					ProjectKey: projectKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				// API returns 400 when branchName is missing
				if resp != nil {
					Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
				}
			})
		})

		Context("Non-Existent Project", func() {
			It("should fail for non-existent project", func() {
				result, resp, err := client.Issues.Pull(&sonargo.IssuesPullOption{
					ProjectKey: nonExistentProjectKey,
					BranchName: "main",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// PullTaint
	// =========================================================================
	Describe("PullTaint", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Issues.PullTaint(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot be nil"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required project key", func() {
				result, resp, err := client.Issues.PullTaint(&sonargo.IssuesPullTaintOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ProjectKey"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required branch name", func() {
				result, resp, err := client.Issues.PullTaint(&sonargo.IssuesPullTaintOption{
					ProjectKey: projectKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				// API returns 400 when branchName is missing
				if resp != nil {
					Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
				}
			})
		})

		Context("Non-Existent Project", func() {
			It("should fail for non-existent project", func() {
				result, resp, err := client.Issues.PullTaint(&sonargo.IssuesPullTaintOption{
					ProjectKey: nonExistentProjectKey,
					BranchName: "main",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// ComponentTags
	// =========================================================================
	Describe("ComponentTags", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Issues.ComponentTags(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot be nil"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required component UUID", func() {
				result, resp, err := client.Issues.ComponentTags(&sonargo.IssuesComponentTagsOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ComponentUuid"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent Component", func() {
			It("should fail with non-existent component UUID", func() {
				result, resp, err := client.Issues.ComponentTags(&sonargo.IssuesComponentTagsOption{
					ComponentUuid: "non-existent-uuid",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})
})
