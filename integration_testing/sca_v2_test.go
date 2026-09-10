package integration_testing_test

import (
	"context"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

var _ = Describe("SCA V2 Service", Ordered, func() {
	var client *sonar.Client

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	Describe("ListClis", func() {
		Context("Functional Tests", func() {
			It("should list CLIs or return an expected error", func() {
				result, resp, err := client.V2.Sca.ListClis(context.Background(), nil)
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
				result, resp, err := client.V2.Sca.GetEnablement(context.Background())
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
				result, resp, err := client.V2.Sca.SetEnablement(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})
		Context("Functional Tests", func() {
			It("should set enablement or return an expected error", func() {
				result, resp, err := client.V2.Sca.SetEnablement(context.Background(), &sonar.ScaSetEnablementOptions{
					Enablement: false,
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

	Describe("SearchDependencyRisks", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.Sca.SearchDependencyRisks(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})
		Context("Functional Tests", func() {
			It("should search or return an expected error", func() {
				result, resp, err := client.V2.Sca.SearchDependencyRisks(context.Background(), &sonar.ScaDependencyRisksSearchOptions{
					ProjectKey: "nonexistent-project",
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

	Describe("SearchReleases", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.Sca.SearchReleases(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})
		Context("Functional Tests", func() {
			It("should search or return an expected error", func() {
				result, resp, err := client.V2.Sca.SearchReleases(context.Background(), &sonar.ScaReleasesSearchOptions{
					ProjectKey: "nonexistent-project",
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

	Describe("GetSbomReport", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.Sca.GetSbomReport(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})
		Context("Functional Tests", func() {
			It("should return SBOM or an expected error", func() {
				result, resp, err := client.V2.Sca.GetSbomReport(context.Background(), &sonar.ScaSbomReportOptions{
					Component: "nonexistent-project",
					Type:      sonar.ScaSbomReportTypeCycloneDX,
					Format:    sonar.ScaSbomReportFormatJSON,
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

	Describe("BulkChangeIssueReleases", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil body", func() {
				result, resp, err := client.V2.Sca.BulkChangeIssueReleases(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail with no issue-release keys", func() {
				result, resp, err := client.V2.Sca.BulkChangeIssueReleases(context.Background(), &sonar.ScaBulkIssueReleaseChangeRequest{})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail with an unknown transition key", func() {
				result, resp, err := client.V2.Sca.BulkChangeIssueReleases(context.Background(), &sonar.ScaBulkIssueReleaseChangeRequest{
					IssueReleaseKeys: []string{"nonexistent-issue-release"},
					TransitionKey:    "NOT_A_TRANSITION",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should apply a bulk change or return an expected error", func() {
				result, resp, err := client.V2.Sca.BulkChangeIssueReleases(context.Background(), &sonar.ScaBulkIssueReleaseChangeRequest{
					IssueReleaseKeys: []string{"nonexistent-issue-release"},
					TransitionKey:    sonar.ScaTransitionConfirm,
					Comment:          "integration test",
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

	Describe("ParseDependencyFiles", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.Sca.ParseDependencyFiles(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without any files", func() {
				result, resp, err := client.V2.Sca.ParseDependencyFiles(context.Background(), &sonar.ScaParseDependencyFilesOptions{
					ProjectKey: "nonexistent-project",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should parse dependency files or return an expected error", func() {
				result, resp, err := client.V2.Sca.ParseDependencyFiles(context.Background(), &sonar.ScaParseDependencyFilesOptions{
					ProjectKey: "nonexistent-project",
					Files: []sonar.ScaDependencyFile{
						{Filename: "package-lock.json", Content: strings.NewReader(`{"lockfileVersion":3,"packages":{}}`)},
					},
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

	Describe("ListReachabilityDefinitions", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.Sca.ListReachabilityDefinitions(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without a language key", func() {
				result, resp, err := client.V2.Sca.ListReachabilityDefinitions(context.Background(), &sonar.ScaReachabilityDefinitionsOptions{})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should return reachability definitions or an expected error", func() {
				result, resp, err := client.V2.Sca.ListReachabilityDefinitions(context.Background(), &sonar.ScaReachabilityDefinitionsOptions{
					LanguageKey: "java",
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
})
