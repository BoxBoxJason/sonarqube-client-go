package integration_testing_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("AlmIntegrations Service", Ordered, func() {
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
	// CheckPat
	// =========================================================================
	Describe("CheckPat", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.AlmIntegrations.CheckPat(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required almSetting", func() {
				_, resp, err := client.AlmIntegrations.CheckPat(context.Background(), &sonar.AlmIntegrationsCheckPatOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("AlmSetting"))
				Expect(resp).To(BeNil())
			})

			It("should fail with almSetting too long", func() {
				longKey := string(make([]byte, 201))
				_, resp, err := client.AlmIntegrations.CheckPat(context.Background(), &sonar.AlmIntegrationsCheckPatOptions{
					AlmSetting: longKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("AlmSetting"))
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// GetGithubClientID
	// =========================================================================
	Describe("GetGithubClientID", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.AlmIntegrations.GetGithubClientID(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required almSetting", func() {
				result, resp, err := client.AlmIntegrations.GetGithubClientID(context.Background(), &sonar.AlmIntegrationsGetGithubClientIDOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("AlmSetting"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// ImportAzureProject
	// =========================================================================
	Describe("ImportAzureProject", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.AlmIntegrations.ImportAzureProject(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required projectName", func() {
				resp, err := client.AlmIntegrations.ImportAzureProject(context.Background(), &sonar.AlmIntegrationsImportAzureProjectOptions{
					RepositoryName: "test-repo",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ProjectName"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required repositoryName", func() {
				resp, err := client.AlmIntegrations.ImportAzureProject(context.Background(), &sonar.AlmIntegrationsImportAzureProjectOptions{
					ProjectName: "test-project",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("RepositoryName"))
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid newCodeDefinitionType", func() {
				resp, err := client.AlmIntegrations.ImportAzureProject(context.Background(), &sonar.AlmIntegrationsImportAzureProjectOptions{
					ProjectName:           "test-project",
					RepositoryName:        "test-repo",
					NewCodeDefinitionType: "INVALID_TYPE",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("NewCodeDefinitionType"))
				Expect(resp).To(BeNil())
			})

			It("should fail with NUMBER_OF_DAYS type and invalid days value", func() {
				resp, err := client.AlmIntegrations.ImportAzureProject(context.Background(), &sonar.AlmIntegrationsImportAzureProjectOptions{
					ProjectName:            "test-project",
					RepositoryName:         "test-repo",
					NewCodeDefinitionType:  sonar.NewCodePeriodTypeNumberOfDays,
					NewCodeDefinitionValue: 0,
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("NewCodeDefinitionValue"))
				Expect(resp).To(BeNil())
			})

			It("should fail with NUMBER_OF_DAYS type and days value too high", func() {
				resp, err := client.AlmIntegrations.ImportAzureProject(context.Background(), &sonar.AlmIntegrationsImportAzureProjectOptions{
					ProjectName:            "test-project",
					RepositoryName:         "test-repo",
					NewCodeDefinitionType:  sonar.NewCodePeriodTypeNumberOfDays,
					NewCodeDefinitionValue: 100,
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("NewCodeDefinitionValue"))
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// ImportBitbucketCloudRepo
	// =========================================================================
	Describe("ImportBitbucketCloudRepo", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.AlmIntegrations.ImportBitbucketCloudRepo(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required repositorySlug", func() {
				resp, err := client.AlmIntegrations.ImportBitbucketCloudRepo(context.Background(), &sonar.AlmIntegrationsImportBitbucketCloudRepoOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("RepositorySlug"))
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid newCodeDefinitionType", func() {
				resp, err := client.AlmIntegrations.ImportBitbucketCloudRepo(context.Background(), &sonar.AlmIntegrationsImportBitbucketCloudRepoOptions{
					RepositorySlug:        "test-slug",
					NewCodeDefinitionType: "INVALID_TYPE",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("NewCodeDefinitionType"))
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// ImportBitbucketServerProject
	// =========================================================================
	Describe("ImportBitbucketServerProject", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.AlmIntegrations.ImportBitbucketServerProject(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required projectKey", func() {
				resp, err := client.AlmIntegrations.ImportBitbucketServerProject(context.Background(), &sonar.AlmIntegrationsImportBitbucketServerProjectOptions{
					RepositorySlug: "test-slug",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ProjectKey"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required repositorySlug", func() {
				resp, err := client.AlmIntegrations.ImportBitbucketServerProject(context.Background(), &sonar.AlmIntegrationsImportBitbucketServerProjectOptions{
					ProjectKey: "test-project",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("RepositorySlug"))
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// ImportGithubProject
	// =========================================================================
	Describe("ImportGithubProject", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.AlmIntegrations.ImportGithubProject(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required repositoryKey", func() {
				resp, err := client.AlmIntegrations.ImportGithubProject(context.Background(), &sonar.AlmIntegrationsImportGithubProjectOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("RepositoryKey"))
				Expect(resp).To(BeNil())
			})

			It("should fail with repositoryKey too long", func() {
				longKey := string(make([]byte, 257))
				resp, err := client.AlmIntegrations.ImportGithubProject(context.Background(), &sonar.AlmIntegrationsImportGithubProjectOptions{
					RepositoryKey: longKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("RepositoryKey"))
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// ImportGitlabProject
	// =========================================================================
	Describe("ImportGitlabProject", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.AlmIntegrations.ImportGitlabProject(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required gitlabProjectId", func() {
				resp, err := client.AlmIntegrations.ImportGitlabProject(context.Background(), &sonar.AlmIntegrationsImportGitlabProjectOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("GitlabProjectId"))
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// ListAzureProjects
	// =========================================================================
	Describe("ListAzureProjects", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.AlmIntegrations.ListAzureProjects(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required almSetting", func() {
				result, resp, err := client.AlmIntegrations.ListAzureProjects(context.Background(), &sonar.AlmIntegrationsListAzureProjectsOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("AlmSetting"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// ListBitbucketServerProjects
	// =========================================================================
	Describe("ListBitbucketServerProjects", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.AlmIntegrations.ListBitbucketServerProjects(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required almSetting", func() {
				result, resp, err := client.AlmIntegrations.ListBitbucketServerProjects(context.Background(), &sonar.AlmIntegrationsListBitbucketServerProjectsOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("AlmSetting"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// ListGithubOrganizations
	// =========================================================================
	Describe("ListGithubOrganizations", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.AlmIntegrations.ListGithubOrganizations(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required almSetting", func() {
				result, resp, err := client.AlmIntegrations.ListGithubOrganizations(context.Background(), &sonar.AlmIntegrationsListGithubOrganizationsOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("AlmSetting"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// ListGithubRepositories
	// =========================================================================
	Describe("ListGithubRepositories", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.AlmIntegrations.ListGithubRepositories(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required almSetting", func() {
				result, resp, err := client.AlmIntegrations.ListGithubRepositories(context.Background(), &sonar.AlmIntegrationsListGithubRepositoriesOptions{
					Organization: "test-org",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("AlmSetting"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required organization", func() {
				result, resp, err := client.AlmIntegrations.ListGithubRepositories(context.Background(), &sonar.AlmIntegrationsListGithubRepositoriesOptions{
					AlmSetting: "test-github",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Organization"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// SearchAzureRepos
	// =========================================================================
	Describe("SearchAzureRepos", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.AlmIntegrations.SearchAzureRepos(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required almSetting", func() {
				result, resp, err := client.AlmIntegrations.SearchAzureRepos(context.Background(), &sonar.AlmIntegrationsSearchAzureReposOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("AlmSetting"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// SearchBitbucketCloudRepos
	// =========================================================================
	Describe("SearchBitbucketCloudRepos", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.AlmIntegrations.SearchBitbucketCloudRepos(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required almSetting", func() {
				result, resp, err := client.AlmIntegrations.SearchBitbucketCloudRepos(context.Background(), &sonar.AlmIntegrationsSearchBitbucketCloudReposOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("AlmSetting"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// SearchBitbucketServerRepos
	// =========================================================================
	Describe("SearchBitbucketServerRepos", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.AlmIntegrations.SearchBitbucketServerRepos(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required almSetting", func() {
				result, resp, err := client.AlmIntegrations.SearchBitbucketServerRepos(context.Background(), &sonar.AlmIntegrationsSearchBitbucketServerReposOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("AlmSetting"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// SearchGitlabRepos
	// =========================================================================
	Describe("SearchGitlabRepos", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.AlmIntegrations.SearchGitlabRepos(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required almSetting", func() {
				result, resp, err := client.AlmIntegrations.SearchGitlabRepos(context.Background(), &sonar.AlmIntegrationsSearchGitlabReposOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("AlmSetting"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// SetPat
	// =========================================================================
	Describe("SetPat", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.AlmIntegrations.SetPat(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required pat", func() {
				resp, err := client.AlmIntegrations.SetPat(context.Background(), &sonar.AlmIntegrationsSetPatOptions{
					AlmSetting: "test-alm",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Pat"))
				Expect(resp).To(BeNil())
			})

			It("should fail with pat too long", func() {
				longPat := string(make([]byte, 2001))
				resp, err := client.AlmIntegrations.SetPat(context.Background(), &sonar.AlmIntegrationsSetPatOptions{
					AlmSetting: "test-alm",
					Pat:        longPat,
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Pat"))
				Expect(resp).To(BeNil())
			})
		})
	})
})
