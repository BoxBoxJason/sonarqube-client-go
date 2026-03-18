package integration_testing_test

import (
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
				_, resp, err := client.AlmIntegrations.CheckPat(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required almSetting", func() {
				_, resp, err := client.AlmIntegrations.CheckPat(&sonar.AlmIntegrationsCheckPatOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("AlmSetting"))
				Expect(resp).To(BeNil())
			})

			It("should fail with almSetting too long", func() {
				longKey := string(make([]byte, 201))
				_, resp, err := client.AlmIntegrations.CheckPat(&sonar.AlmIntegrationsCheckPatOptions{
					AlmSetting: longKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("AlmSetting"))
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// GetGithubClientId
	// =========================================================================
	Describe("GetGithubClientId", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.AlmIntegrations.GetGithubClientId(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required almSetting", func() {
				result, resp, err := client.AlmIntegrations.GetGithubClientId(&sonar.AlmIntegrationsGetGithubClientIdOptions{})
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
				resp, err := client.AlmIntegrations.ImportAzureProject(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required projectName", func() {
				resp, err := client.AlmIntegrations.ImportAzureProject(&sonar.AlmIntegrationsImportAzureProjectOptions{
					RepositoryName: "test-repo",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ProjectName"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required repositoryName", func() {
				resp, err := client.AlmIntegrations.ImportAzureProject(&sonar.AlmIntegrationsImportAzureProjectOptions{
					ProjectName: "test-project",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("RepositoryName"))
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid newCodeDefinitionType", func() {
				resp, err := client.AlmIntegrations.ImportAzureProject(&sonar.AlmIntegrationsImportAzureProjectOptions{
					ProjectName:           "test-project",
					RepositoryName:        "test-repo",
					NewCodeDefinitionType: "INVALID_TYPE",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("NewCodeDefinitionType"))
				Expect(resp).To(BeNil())
			})

			It("should fail with NUMBER_OF_DAYS type and invalid days value", func() {
				resp, err := client.AlmIntegrations.ImportAzureProject(&sonar.AlmIntegrationsImportAzureProjectOptions{
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
				resp, err := client.AlmIntegrations.ImportAzureProject(&sonar.AlmIntegrationsImportAzureProjectOptions{
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
				resp, err := client.AlmIntegrations.ImportBitbucketCloudRepo(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required repositorySlug", func() {
				resp, err := client.AlmIntegrations.ImportBitbucketCloudRepo(&sonar.AlmIntegrationsImportBitbucketCloudRepoOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("RepositorySlug"))
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid newCodeDefinitionType", func() {
				resp, err := client.AlmIntegrations.ImportBitbucketCloudRepo(&sonar.AlmIntegrationsImportBitbucketCloudRepoOptions{
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
				resp, err := client.AlmIntegrations.ImportBitbucketServerProject(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required projectKey", func() {
				resp, err := client.AlmIntegrations.ImportBitbucketServerProject(&sonar.AlmIntegrationsImportBitbucketServerProjectOptions{
					RepositorySlug: "test-slug",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ProjectKey"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required repositorySlug", func() {
				resp, err := client.AlmIntegrations.ImportBitbucketServerProject(&sonar.AlmIntegrationsImportBitbucketServerProjectOptions{
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
				resp, err := client.AlmIntegrations.ImportGithubProject(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required repositoryKey", func() {
				resp, err := client.AlmIntegrations.ImportGithubProject(&sonar.AlmIntegrationsImportGithubProjectOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("RepositoryKey"))
				Expect(resp).To(BeNil())
			})

			It("should fail with repositoryKey too long", func() {
				longKey := string(make([]byte, 257))
				resp, err := client.AlmIntegrations.ImportGithubProject(&sonar.AlmIntegrationsImportGithubProjectOptions{
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
				resp, err := client.AlmIntegrations.ImportGitlabProject(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required gitlabProjectId", func() {
				resp, err := client.AlmIntegrations.ImportGitlabProject(&sonar.AlmIntegrationsImportGitlabProjectOptions{})
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
				result, resp, err := client.AlmIntegrations.ListAzureProjects(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required almSetting", func() {
				result, resp, err := client.AlmIntegrations.ListAzureProjects(&sonar.AlmIntegrationsListAzureProjectsOptions{})
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
				result, resp, err := client.AlmIntegrations.ListBitbucketServerProjects(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required almSetting", func() {
				result, resp, err := client.AlmIntegrations.ListBitbucketServerProjects(&sonar.AlmIntegrationsListBitbucketServerProjectsOptions{})
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
				result, resp, err := client.AlmIntegrations.ListGithubOrganizations(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required almSetting", func() {
				result, resp, err := client.AlmIntegrations.ListGithubOrganizations(&sonar.AlmIntegrationsListGithubOrganizationsOptions{})
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
				result, resp, err := client.AlmIntegrations.ListGithubRepositories(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required almSetting", func() {
				result, resp, err := client.AlmIntegrations.ListGithubRepositories(&sonar.AlmIntegrationsListGithubRepositoriesOptions{
					Organization: "test-org",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("AlmSetting"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required organization", func() {
				result, resp, err := client.AlmIntegrations.ListGithubRepositories(&sonar.AlmIntegrationsListGithubRepositoriesOptions{
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
				result, resp, err := client.AlmIntegrations.SearchAzureRepos(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required almSetting", func() {
				result, resp, err := client.AlmIntegrations.SearchAzureRepos(&sonar.AlmIntegrationsSearchAzureReposOptions{})
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
				result, resp, err := client.AlmIntegrations.SearchBitbucketCloudRepos(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required almSetting", func() {
				result, resp, err := client.AlmIntegrations.SearchBitbucketCloudRepos(&sonar.AlmIntegrationsSearchBitbucketCloudReposOptions{})
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
				result, resp, err := client.AlmIntegrations.SearchBitbucketServerRepos(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required almSetting", func() {
				result, resp, err := client.AlmIntegrations.SearchBitbucketServerRepos(&sonar.AlmIntegrationsSearchBitbucketServerReposOptions{})
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
				result, resp, err := client.AlmIntegrations.SearchGitlabRepos(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail without required almSetting", func() {
				result, resp, err := client.AlmIntegrations.SearchGitlabRepos(&sonar.AlmIntegrationsSearchGitlabReposOptions{})
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
				resp, err := client.AlmIntegrations.SetPat(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required pat", func() {
				resp, err := client.AlmIntegrations.SetPat(&sonar.AlmIntegrationsSetPatOptions{
					AlmSetting: "test-alm",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Pat"))
				Expect(resp).To(BeNil())
			})

			It("should fail with pat too long", func() {
				longPat := string(make([]byte, 2001))
				resp, err := client.AlmIntegrations.SetPat(&sonar.AlmIntegrationsSetPatOptions{
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
