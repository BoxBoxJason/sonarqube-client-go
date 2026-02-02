package integration_testing_test

import (
"net/http"

. "github.com/onsi/ginkgo/v2"
. "github.com/onsi/gomega"

sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Measures Service", Ordered, func() {
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

		// Create a test project for measures-related operations
		projectKey = helpers.UniqueResourceName("msr")
		_, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
			Name:    "Measures Test Project",
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
	// Component
	// =========================================================================
	Describe("Component", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Measures.Component(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("option struct is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required component parameter", func() {
				result, resp, err := client.Measures.Component(&sonargo.MeasuresComponentOption{
					MetricKeys: []string{"ncloc"},
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Component"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required metric keys", func() {
				result, resp, err := client.Measures.Component(&sonargo.MeasuresComponentOption{
					Component: projectKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("MetricKeys"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Valid Requests", func() {
			It("should get component measures for an existing project", func() {
				result, resp, err := client.Measures.Component(&sonargo.MeasuresComponentOption{
					Component:  projectKey,
					MetricKeys: []string{"ncloc", "complexity", "coverage"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Component.Key).To(Equal(projectKey))
			})

			It("should get measures with single metric", func() {
				result, resp, err := client.Measures.Component(&sonargo.MeasuresComponentOption{
					Component:  projectKey,
					MetricKeys: []string{"lines"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Component.Key).To(Equal(projectKey))
			})

			It("should get measures with additional fields", func() {
				result, resp, err := client.Measures.Component(&sonargo.MeasuresComponentOption{
					Component:        projectKey,
					MetricKeys:       []string{"ncloc"},
					AdditionalFields: []string{"metrics"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("Non-Existent Component", func() {
			It("should fail with non-existent component", func() {
				result, resp, err := client.Measures.Component(&sonargo.MeasuresComponentOption{
					Component:  "non-existent-project",
					MetricKeys: []string{"ncloc"},
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
				result, resp, err := client.Measures.Component(&sonargo.MeasuresComponentOption{
					Component:  projectKey,
					MetricKeys: []string{"ncloc"},
					Branch:     "non-existent-branch",
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
				result, resp, err := client.Measures.Component(&sonargo.MeasuresComponentOption{
					Component:   projectKey,
					MetricKeys:  []string{"ncloc"},
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
	// ComponentTree
	// =========================================================================
	Describe("ComponentTree", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Measures.ComponentTree(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("option struct is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required component parameter", func() {
				result, resp, err := client.Measures.ComponentTree(&sonargo.MeasuresComponentTreeOption{
					MetricKeys: []string{"ncloc"},
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Component"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required metric keys", func() {
				result, resp, err := client.Measures.ComponentTree(&sonargo.MeasuresComponentTreeOption{
					Component: projectKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("MetricKeys"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid metric sort filter", func() {
				result, resp, err := client.Measures.ComponentTree(&sonargo.MeasuresComponentTreeOption{
					Component:        projectKey,
					MetricKeys:       []string{"ncloc"},
					MetricSortFilter: "invalid",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("MetricSortFilter"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid strategy", func() {
				result, resp, err := client.Measures.ComponentTree(&sonargo.MeasuresComponentTreeOption{
					Component:  projectKey,
					MetricKeys: []string{"ncloc"},
					Strategy:   "invalid",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Strategy"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Valid Requests", func() {
			It("should get component tree measures for an existing project", func() {
				result, resp, err := client.Measures.ComponentTree(&sonargo.MeasuresComponentTreeOption{
					Component:  projectKey,
					MetricKeys: []string{"ncloc", "complexity"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.BaseComponent.Key).To(Equal(projectKey))
			})

			It("should get component tree with children strategy", func() {
				result, resp, err := client.Measures.ComponentTree(&sonargo.MeasuresComponentTreeOption{
					Component:  projectKey,
					MetricKeys: []string{"ncloc"},
					Strategy:   "children",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should get component tree with leaves strategy", func() {
				result, resp, err := client.Measures.ComponentTree(&sonargo.MeasuresComponentTreeOption{
					Component:  projectKey,
					MetricKeys: []string{"ncloc"},
					Strategy:   "leaves",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should get component tree with all strategy", func() {
				result, resp, err := client.Measures.ComponentTree(&sonargo.MeasuresComponentTreeOption{
					Component:  projectKey,
					MetricKeys: []string{"ncloc"},
					Strategy:   "all",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should get component tree with metric sort filter", func() {
				result, resp, err := client.Measures.ComponentTree(&sonargo.MeasuresComponentTreeOption{
					Component:        projectKey,
					MetricKeys:       []string{"ncloc"},
					MetricSortFilter: "withMeasuresOnly",
					MetricSort:       "ncloc",
					Sort:             []string{"metric"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should get component tree with pagination", func() {
				result, resp, err := client.Measures.ComponentTree(&sonargo.MeasuresComponentTreeOption{
					Component:  projectKey,
					MetricKeys: []string{"ncloc"},
					PaginationArgs: sonargo.PaginationArgs{
						PageSize: 10,
						Page:     1,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Paging.PageSize).To(BeNumerically("<=", 10))
			})

			It("should get component tree with qualifiers filter", func() {
				result, resp, err := client.Measures.ComponentTree(&sonargo.MeasuresComponentTreeOption{
					Component:  projectKey,
					MetricKeys: []string{"ncloc"},
					Qualifiers: []string{"FIL"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should get component tree with additional fields", func() {
				result, resp, err := client.Measures.ComponentTree(&sonargo.MeasuresComponentTreeOption{
					Component:        projectKey,
					MetricKeys:       []string{"ncloc"},
					AdditionalFields: []string{"metrics"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("Non-Existent Component", func() {
			It("should fail with non-existent component", func() {
				result, resp, err := client.Measures.ComponentTree(&sonargo.MeasuresComponentTreeOption{
					Component:  "non-existent-project",
					MetricKeys: []string{"ncloc"},
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
	// Search
	// =========================================================================
	Describe("Search", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Measures.Search(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("option struct is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required metric keys", func() {
				result, resp, err := client.Measures.Search(&sonargo.MeasuresSearchOption{
					ProjectKeys: []string{projectKey},
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("MetricKeys"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required project keys", func() {
				result, resp, err := client.Measures.Search(&sonargo.MeasuresSearchOption{
					MetricKeys: []string{"ncloc"},
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ProjectKeys"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Valid Requests", func() {
			It("should search measures for an existing project", func() {
				result, resp, err := client.Measures.Search(&sonargo.MeasuresSearchOption{
					MetricKeys:  []string{"ncloc"},
					ProjectKeys: []string{projectKey},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should search measures with multiple metrics", func() {
				result, resp, err := client.Measures.Search(&sonargo.MeasuresSearchOption{
					MetricKeys:  []string{"ncloc", "complexity", "coverage"},
					ProjectKeys: []string{projectKey},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("Multiple Projects", func() {
			var secondProjectKey string

			BeforeAll(func() {
				secondProjectKey = helpers.UniqueResourceName("msr2")
				_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
					Name:    "Second Measures Test Project",
					Project: secondProjectKey,
				})
				Expect(err).NotTo(HaveOccurred())

				cleanup.RegisterCleanup("project", secondProjectKey, func() error {
					_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
						Project: secondProjectKey,
					})
					return err
				})
			})

			It("should search measures for multiple projects", func() {
				result, resp, err := client.Measures.Search(&sonargo.MeasuresSearchOption{
					MetricKeys:  []string{"ncloc"},
					ProjectKeys: []string{projectKey, secondProjectKey},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("Non-Existent Project", func() {
			It("should return empty measures for non-existent project", func() {
				// The API doesn't fail for non-existent projects, it just returns empty measures
				result, resp, err := client.Measures.Search(&sonargo.MeasuresSearchOption{
					MetricKeys:  []string{"ncloc"},
					ProjectKeys: []string{"non-existent-project"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				// Empty measures for non-existent project
				Expect(result.Measures).To(BeEmpty())
			})
		})
	})

	// =========================================================================
	// SearchHistory
	// =========================================================================
	Describe("SearchHistory", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Measures.SearchHistory(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("option struct is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required component parameter", func() {
				result, resp, err := client.Measures.SearchHistory(&sonargo.MeasuresSearchHistoryOption{
					Metrics: []string{"ncloc"},
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Component"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required metrics", func() {
				result, resp, err := client.Measures.SearchHistory(&sonargo.MeasuresSearchHistoryOption{
					Component: projectKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Metrics"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Valid Requests", func() {
			It("should get measure history for an existing project", func() {
				result, resp, err := client.Measures.SearchHistory(&sonargo.MeasuresSearchHistoryOption{
					Component: projectKey,
					Metrics:   []string{"ncloc"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should get history with multiple metrics", func() {
				result, resp, err := client.Measures.SearchHistory(&sonargo.MeasuresSearchHistoryOption{
					Component: projectKey,
					Metrics:   []string{"ncloc", "complexity", "coverage"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should get history with pagination", func() {
				result, resp, err := client.Measures.SearchHistory(&sonargo.MeasuresSearchHistoryOption{
					Component: projectKey,
					Metrics:   []string{"ncloc"},
					PaginationArgs: sonargo.PaginationArgs{
						PageSize: 10,
						Page:     1,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Paging.PageSize).To(BeNumerically("<=", 10))
			})

			It("should get history with date range", func() {
				result, resp, err := client.Measures.SearchHistory(&sonargo.MeasuresSearchHistoryOption{
					Component: projectKey,
					Metrics:   []string{"ncloc"},
					From:      "2020-01-01",
					To:        "2030-12-31",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("Non-Existent Component", func() {
			It("should fail with non-existent component", func() {
				result, resp, err := client.Measures.SearchHistory(&sonargo.MeasuresSearchHistoryOption{
					Component: "non-existent-project",
					Metrics:   []string{"ncloc"},
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
				result, resp, err := client.Measures.SearchHistory(&sonargo.MeasuresSearchHistoryOption{
					Component: projectKey,
					Metrics:   []string{"ncloc"},
					Branch:    "non-existent-branch",
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
				result, resp, err := client.Measures.SearchHistory(&sonargo.MeasuresSearchHistoryOption{
					Component:   projectKey,
					Metrics:     []string{"ncloc"},
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
})
