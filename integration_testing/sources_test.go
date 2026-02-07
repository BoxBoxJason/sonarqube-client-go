package integration_testing_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("Sources Service", Ordered, func() {
	var (
		client     *sonar.Client
		cleanup    *helpers.CleanupManager
		projectKey string
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
		cleanup = helpers.NewCleanupManager(client)

		// Create a test project for source-related operations
		projectKey = helpers.UniqueResourceName("src")
		_, _, err = client.Projects.Create(&sonar.ProjectsCreateOption{
			Name:    "Sources Test Project",
			Project: projectKey,
		})
		Expect(err).NotTo(HaveOccurred())

		cleanup.RegisterCleanup("project", projectKey, func() error {
			_, err := client.Projects.Delete(&sonar.ProjectsDeleteOption{
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
	// Index (Deprecated)
	// =========================================================================
	Describe("Index", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Sources.Index(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("option struct is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required resource parameter", func() {
				result, resp, err := client.Sources.Index(&sonar.SourcesIndexOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Resource"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent Resource", func() {
			It("should fail with non-existent file key", func() {
				result, resp, err := client.Sources.Index(&sonar.SourcesIndexOption{
					Resource: "non-existent:src/main.go",
				})
				// API should return an error for non-existent resource
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})

		Context("With Valid Project (No Analysis)", func() {
			It("should fail for project without analyzed files", func() {
				// Projects without analysis have no source files indexed
				result, resp, err := client.Sources.Index(&sonar.SourcesIndexOption{
					Resource: projectKey + ":src/main.go",
				})
				// Should fail because there's no analyzed source
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})

		Context("With Line Range", func() {
			It("should fail for non-existent file with line range", func() {
				result, resp, err := client.Sources.Index(&sonar.SourcesIndexOption{
					Resource: "non-existent:src/main.go",
					From:     1,
					To:       10,
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
	// IssueSnippets
	// =========================================================================
	Describe("IssueSnippets", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Sources.IssueSnippets(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("option struct is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required issue key", func() {
				result, resp, err := client.Sources.IssueSnippets(&sonar.SourcesIssueSnippetsOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("IssueKey"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent Issue", func() {
			It("should fail with non-existent issue key", func() {
				result, resp, err := client.Sources.IssueSnippets(&sonar.SourcesIssueSnippetsOption{
					IssueKey: "AXxxxxxxxxxxxxxxxxxx",
				})
				// API should return an error for non-existent issue
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// Lines
	// =========================================================================
	Describe("Lines", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Sources.Lines(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("option struct is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required key parameter", func() {
				result, resp, err := client.Sources.Lines(&sonar.SourcesLinesOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent File", func() {
			It("should fail with non-existent file key", func() {
				result, resp, err := client.Sources.Lines(&sonar.SourcesLinesOption{
					Key: "non-existent:src/main.go",
				})
				// API should return an error for non-existent file
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})

		Context("With Valid Project (No Analysis)", func() {
			It("should fail for project without analyzed files", func() {
				result, resp, err := client.Sources.Lines(&sonar.SourcesLinesOption{
					Key: projectKey + ":src/main.go",
				})
				// Should fail because there's no analyzed source
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})

		Context("With Line Range", func() {
			It("should fail for non-existent file with line range", func() {
				result, resp, err := client.Sources.Lines(&sonar.SourcesLinesOption{
					Key:  "non-existent:src/main.go",
					From: 1,
					To:   10,
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})

		Context("With Branch", func() {
			It("should fail for non-existent branch", func() {
				result, resp, err := client.Sources.Lines(&sonar.SourcesLinesOption{
					Key:    projectKey + ":src/main.go",
					Branch: "non-existent-branch",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})

		Context("With Pull Request", func() {
			It("should fail for non-existent pull request", func() {
				result, resp, err := client.Sources.Lines(&sonar.SourcesLinesOption{
					Key:         projectKey + ":src/main.go",
					PullRequest: "99999",
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
	// Raw
	// =========================================================================
	Describe("Raw", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Sources.Raw(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("option struct is required"))
				Expect(result).To(BeEmpty())
				Expect(resp).To(BeNil())
			})

			It("should fail without required key parameter", func() {
				result, resp, err := client.Sources.Raw(&sonar.SourcesRawOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(result).To(BeEmpty())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent File", func() {
			It("should fail with non-existent file key", func() {
				result, resp, err := client.Sources.Raw(&sonar.SourcesRawOption{
					Key: "non-existent:src/main.go",
				})
				// API should return an error for non-existent file
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeEmpty())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})

		Context("With Valid Project (No Analysis)", func() {
			It("should fail for project without analyzed files", func() {
				result, resp, err := client.Sources.Raw(&sonar.SourcesRawOption{
					Key: projectKey + ":src/main.go",
				})
				// Should fail because there's no analyzed source
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeEmpty())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})

		Context("With Branch", func() {
			It("should fail for non-existent branch", func() {
				result, resp, err := client.Sources.Raw(&sonar.SourcesRawOption{
					Key:    projectKey + ":src/main.go",
					Branch: "non-existent-branch",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeEmpty())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})

		Context("With Pull Request", func() {
			It("should fail for non-existent pull request", func() {
				result, resp, err := client.Sources.Raw(&sonar.SourcesRawOption{
					Key:         projectKey + ":src/main.go",
					PullRequest: "99999",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeEmpty())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// Scm
	// =========================================================================
	Describe("Scm", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Sources.Scm(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("option struct is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required key parameter", func() {
				result, resp, err := client.Sources.Scm(&sonar.SourcesScmOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent File", func() {
			It("should fail with non-existent file key", func() {
				result, resp, err := client.Sources.Scm(&sonar.SourcesScmOption{
					Key: "non-existent:src/main.go",
				})
				// API should return an error for non-existent file
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})

		Context("With Valid Project (No Analysis)", func() {
			It("should fail for project without analyzed files", func() {
				result, resp, err := client.Sources.Scm(&sonar.SourcesScmOption{
					Key: projectKey + ":src/main.go",
				})
				// Should fail because there's no analyzed source
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})

		Context("With Line Range", func() {
			It("should fail for non-existent file with line range", func() {
				result, resp, err := client.Sources.Scm(&sonar.SourcesScmOption{
					Key:  "non-existent:src/main.go",
					From: 1,
					To:   10,
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})

		Context("With CommitsByLine Option", func() {
			It("should fail for non-existent file with commits_by_line", func() {
				result, resp, err := client.Sources.Scm(&sonar.SourcesScmOption{
					Key:           "non-existent:src/main.go",
					CommitsByLine: true,
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
	// Show
	// =========================================================================
	Describe("Show", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Sources.Show(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("option struct is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required key parameter", func() {
				result, resp, err := client.Sources.Show(&sonar.SourcesShowOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent File", func() {
			It("should fail with non-existent file key", func() {
				result, resp, err := client.Sources.Show(&sonar.SourcesShowOption{
					Key: "non-existent:src/main.go",
				})
				// API should return an error for non-existent file
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})

		Context("With Valid Project (No Analysis)", func() {
			It("should fail for project without analyzed files", func() {
				result, resp, err := client.Sources.Show(&sonar.SourcesShowOption{
					Key: projectKey + ":src/main.go",
				})
				// Should fail because there's no analyzed source
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})

		Context("With Line Range", func() {
			It("should fail for non-existent file with line range", func() {
				result, resp, err := client.Sources.Show(&sonar.SourcesShowOption{
					Key:  "non-existent:src/main.go",
					From: 1,
					To:   10,
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
