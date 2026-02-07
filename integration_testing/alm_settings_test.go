package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("AlmSettings Service", Ordered, func() {
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

		// Create a test project
		projectKey = helpers.UniqueResourceName("alm")
		_, _, err = client.Projects.Create(&sonar.ProjectsCreateOption{
			Name:    "AlmSettings Test Project",
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
	// ListDefinitions
	// =========================================================================
	Describe("ListDefinitions", func() {
		It("should list ALM setting definitions", func() {
			result, resp, err := client.AlmSettings.ListDefinitions()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})
	})

	// =========================================================================
	// List
	// =========================================================================
	Describe("List", func() {
		Context("Valid Requests", func() {
			It("should list ALM settings with nil options", func() {
				result, resp, err := client.AlmSettings.List(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should list ALM settings for a project", func() {
				result, resp, err := client.AlmSettings.List(&sonar.AlmSettingsListOption{
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})
	})

	// =========================================================================
	// GetBinding
	// =========================================================================
	Describe("GetBinding", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.AlmSettings.GetBinding(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required project", func() {
				result, resp, err := client.AlmSettings.GetBinding(&sonar.AlmSettingsGetBindingOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Project"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Bound Project", func() {
			It("should fail for project without ALM binding", func() {
				result, resp, err := client.AlmSettings.GetBinding(&sonar.AlmSettingsGetBindingOption{
					Project: projectKey,
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
	// CountBinding
	// =========================================================================
	Describe("CountBinding", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.AlmSettings.CountBinding(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required almSetting", func() {
				result, resp, err := client.AlmSettings.CountBinding(&sonar.AlmSettingsCountBindingOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("AlmSetting"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent ALM Setting", func() {
			It("should fail for non-existent ALM setting", func() {
				result, resp, err := client.AlmSettings.CountBinding(&sonar.AlmSettingsCountBindingOption{
					AlmSetting: "non-existent-alm-setting",
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
	// Delete
	// =========================================================================
	Describe("Delete", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.AlmSettings.Delete(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required key", func() {
				resp, err := client.AlmSettings.Delete(&sonar.AlmSettingsDeleteOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent ALM Setting", func() {
			It("should fail for non-existent ALM setting", func() {
				resp, err := client.AlmSettings.Delete(&sonar.AlmSettingsDeleteOption{
					Key: "non-existent-alm-setting",
				})
				Expect(err).To(HaveOccurred())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// CreateAzure
	// =========================================================================
	Describe("CreateAzure", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.AlmSettings.CreateAzure(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required key", func() {
				resp, err := client.AlmSettings.CreateAzure(&sonar.AlmSettingsCreateAzureOption{
					PersonalAccessToken: "test-token",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required personalAccessToken", func() {
				resp, err := client.AlmSettings.CreateAzure(&sonar.AlmSettingsCreateAzureOption{
					Key: "test-azure",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("PersonalAccessToken"))
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// CreateGithub
	// =========================================================================
	Describe("CreateGithub", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.AlmSettings.CreateGithub(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required key", func() {
				resp, err := client.AlmSettings.CreateGithub(&sonar.AlmSettingsCreateGithubOption{
					AppID:        "test-app-id",
					ClientID:     "test-client-id",
					ClientSecret: "test-client-secret",
					PrivateKey:   "test-private-key",
					URL:          "https://api.github.com",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required appId", func() {
				resp, err := client.AlmSettings.CreateGithub(&sonar.AlmSettingsCreateGithubOption{
					Key:          "test-github",
					ClientID:     "test-client-id",
					ClientSecret: "test-client-secret",
					PrivateKey:   "test-private-key",
					URL:          "https://api.github.com",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("AppID"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required privateKey", func() {
				resp, err := client.AlmSettings.CreateGithub(&sonar.AlmSettingsCreateGithubOption{
					Key:          "test-github",
					AppID:        "test-app-id",
					ClientID:     "test-client-id",
					ClientSecret: "test-client-secret",
					URL:          "https://api.github.com",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("PrivateKey"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required URL", func() {
				resp, err := client.AlmSettings.CreateGithub(&sonar.AlmSettingsCreateGithubOption{
					Key:          "test-github",
					AppID:        "test-app-id",
					ClientID:     "test-client-id",
					ClientSecret: "test-client-secret",
					PrivateKey:   "test-private-key",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("URL"))
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// CreateGitlab
	// =========================================================================
	Describe("CreateGitlab", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.AlmSettings.CreateGitlab(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required key", func() {
				resp, err := client.AlmSettings.CreateGitlab(&sonar.AlmSettingsCreateGitlabOption{
					PersonalAccessToken: "test-token",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required personalAccessToken", func() {
				resp, err := client.AlmSettings.CreateGitlab(&sonar.AlmSettingsCreateGitlabOption{
					Key: "test-gitlab",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("PersonalAccessToken"))
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// CreateBitbucket
	// =========================================================================
	Describe("CreateBitbucket", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.AlmSettings.CreateBitbucket(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required key", func() {
				resp, err := client.AlmSettings.CreateBitbucket(&sonar.AlmSettingsCreateBitbucketOption{
					PersonalAccessToken: "test-token",
					URL:                 "https://bitbucket.example.com",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required URL", func() {
				resp, err := client.AlmSettings.CreateBitbucket(&sonar.AlmSettingsCreateBitbucketOption{
					Key:                 "test-bitbucket",
					PersonalAccessToken: "test-token",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("URL"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required personalAccessToken", func() {
				resp, err := client.AlmSettings.CreateBitbucket(&sonar.AlmSettingsCreateBitbucketOption{
					Key: "test-bitbucket",
					URL: "https://bitbucket.example.com",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("PersonalAccessToken"))
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// CreateBitbucketCloud
	// =========================================================================
	Describe("CreateBitbucketCloud", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.AlmSettings.CreateBitbucketCloud(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required key", func() {
				resp, err := client.AlmSettings.CreateBitbucketCloud(&sonar.AlmSettingsCreateBitbucketCloudOption{
					ClientID:     "test-client-id",
					ClientSecret: "test-client-secret",
					Workspace:    "test-workspace",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required clientId", func() {
				resp, err := client.AlmSettings.CreateBitbucketCloud(&sonar.AlmSettingsCreateBitbucketCloudOption{
					Key:          "test-bitbucket-cloud",
					ClientSecret: "test-client-secret",
					Workspace:    "test-workspace",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ClientID"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required clientSecret", func() {
				resp, err := client.AlmSettings.CreateBitbucketCloud(&sonar.AlmSettingsCreateBitbucketCloudOption{
					Key:       "test-bitbucket-cloud",
					ClientID:  "test-client-id",
					Workspace: "test-workspace",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ClientSecret"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required workspace", func() {
				resp, err := client.AlmSettings.CreateBitbucketCloud(&sonar.AlmSettingsCreateBitbucketCloudOption{
					Key:          "test-bitbucket-cloud",
					ClientID:     "test-client-id",
					ClientSecret: "test-client-secret",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Workspace"))
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// UpdateAzure
	// =========================================================================
	Describe("UpdateAzure", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.AlmSettings.UpdateAzure(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required key", func() {
				resp, err := client.AlmSettings.UpdateAzure(&sonar.AlmSettingsUpdateAzureOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// UpdateGithub
	// =========================================================================
	Describe("UpdateGithub", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.AlmSettings.UpdateGithub(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required key", func() {
				resp, err := client.AlmSettings.UpdateGithub(&sonar.AlmSettingsUpdateGithubOption{
					AppID:    "test-app-id",
					ClientID: "test-client-id",
					URL:      "https://api.github.com",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required appId", func() {
				resp, err := client.AlmSettings.UpdateGithub(&sonar.AlmSettingsUpdateGithubOption{
					Key:      "test-github",
					ClientID: "test-client-id",
					URL:      "https://api.github.com",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("AppID"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required URL", func() {
				resp, err := client.AlmSettings.UpdateGithub(&sonar.AlmSettingsUpdateGithubOption{
					Key:      "test-github",
					AppID:    "test-app-id",
					ClientID: "test-client-id",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("URL"))
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// UpdateGitlab
	// =========================================================================
	Describe("UpdateGitlab", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.AlmSettings.UpdateGitlab(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required key", func() {
				resp, err := client.AlmSettings.UpdateGitlab(&sonar.AlmSettingsUpdateGitlabOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// UpdateBitbucket
	// =========================================================================
	Describe("UpdateBitbucket", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.AlmSettings.UpdateBitbucket(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required key", func() {
				resp, err := client.AlmSettings.UpdateBitbucket(&sonar.AlmSettingsUpdateBitbucketOption{
					URL: "https://bitbucket.example.com",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required URL", func() {
				resp, err := client.AlmSettings.UpdateBitbucket(&sonar.AlmSettingsUpdateBitbucketOption{
					Key: "test-bitbucket",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("URL"))
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// UpdateBitbucketCloud
	// =========================================================================
	Describe("UpdateBitbucketCloud", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.AlmSettings.UpdateBitbucketCloud(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required key", func() {
				resp, err := client.AlmSettings.UpdateBitbucketCloud(&sonar.AlmSettingsUpdateBitbucketCloudOption{
					ClientID:  "test-client-id",
					Workspace: "test-workspace",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Key"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required workspace", func() {
				resp, err := client.AlmSettings.UpdateBitbucketCloud(&sonar.AlmSettingsUpdateBitbucketCloudOption{
					Key:      "test-bitbucket-cloud",
					ClientID: "test-client-id",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Workspace"))
				Expect(resp).To(BeNil())
			})
		})
	})
})
