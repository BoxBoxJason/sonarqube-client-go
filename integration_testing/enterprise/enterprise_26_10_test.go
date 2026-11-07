package enterprise_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

// Specs for the Enterprise endpoints added or covered with SonarQube Enterprise 2026.5.3.
// They run against an unlicensed instance, so nothing here needs an analysis or a license.
var _ = Describe("SonarQube Enterprise 2026.5.3 endpoints", Ordered, func() {
	var (
		client      *sonar.Client
		ctx         = context.Background()
		projectKey  string
		appKey      string
		gateName    string
		almKey      string
		dopSettings string
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())

		projectKey = helpers.UniqueResourceName("ent-26-10")
		_, _, err = client.Projects.Create(ctx, &sonar.ProjectsCreateOptions{Name: projectKey, Project: projectKey})
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			_, _ = client.Projects.Delete(ctx, &sonar.ProjectsDeleteOptions{Project: projectKey})
		})

		appKey = helpers.UniqueResourceName("ent-app")
		_, _, err = client.Applications.Create(ctx, &sonar.ApplicationsCreateOptions{Name: appKey, Key: appKey})
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			_, _ = client.Applications.Delete(ctx, &sonar.ApplicationsDeleteOptions{Application: appKey})
		})

		gateName = helpers.UniqueResourceName("ent-gate")
		_, _, err = client.Qualitygates.Create(ctx, &sonar.QualitygatesCreateOptions{Name: gateName})
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			_, _ = client.Qualitygates.Delete(ctx, &sonar.QualitygatesDeleteOptions{Name: gateName})
		})

		almKey = helpers.UniqueResourceName("ent-gh")
		_, err = client.AlmSettings.CreateGithub(ctx, &sonar.AlmSettingsCreateGithubOptions{
			Key: almKey, AppID: "123", ClientID: "cid", ClientSecret: "secret", PrivateKey: "pk", URL: "https://api.github.com",
		})
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			_, _ = client.AlmSettings.Delete(ctx, &sonar.AlmSettingsDeleteOptions{Key: almKey})
		})

		_, err = client.AlmSettings.SetGithubBinding(ctx, &sonar.AlmSettingsSetGithubBindingOptions{
			AlmSetting: almKey, Project: projectKey, Repository: "org/repo",
		})
		Expect(err).NotTo(HaveOccurred())

		settings, _, err := client.V2.DopTranslation.GetDopSettings(ctx)
		Expect(err).NotTo(HaveOccurred())

		for _, setting := range settings.DopSettings {
			if setting.Key == almKey {
				dopSettings = setting.Id
			}
		}

		Expect(dopSettings).NotTo(BeEmpty())
	})

	Describe("AI code assurance", func() {
		It("should flag a project as containing AI code", func() {
			resp, err := client.Projects.SetContainsAiCode(ctx, &sonar.ProjectsSetContainsAiCodeOptions{Project: projectKey, ContainsAiCode: true})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			contains, _, err := client.Projects.GetContainsAiCode(ctx, &sonar.ProjectsProjectOptions{Project: projectKey})
			Expect(err).NotTo(HaveOccurred())
			Expect(contains.ContainsAiCode).To(BeTrue())

			detected, _, err := client.Projects.GetDetectedAiCode(ctx, &sonar.ProjectsProjectOptions{Project: projectKey})
			Expect(err).NotTo(HaveOccurred())
			Expect(detected).NotTo(BeNil())

			assurance, _, err := client.ProjectBranches.GetAiCodeAssurance(ctx, &sonar.ProjectBranchesGetAiCodeAssuranceOptions{Project: projectKey})
			Expect(err).NotTo(HaveOccurred())
			Expect(assurance.AiCodeAssurance).NotTo(BeEmpty())
		})

		It("should generate the AI code assurance badge", func() {
			badge, _, err := client.ProjectBadges.AiCodeAssurance(ctx, &sonar.ProjectBadgesAiCodeAssuranceOptions{Project: projectKey})
			Expect(err).NotTo(HaveOccurred())
			Expect(*badge).To(ContainSubstring("<svg"))
		})

		It("should qualify and disqualify a quality gate as AI Code Assured", func() {
			for _, enabled := range []bool{true, false} {
				resp, err := client.Qualitygates.SetAiCodeAssurance(ctx, &sonar.QualitygatesSetAiCodeAssuranceOptions{GateName: gateName, AiCodeAssurance: enabled})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			}
		})
	})

	Describe("Findings exports", func() {
		It("should export the GitLab SAST report", func() {
			result, _, err := client.Issues.GitlabSastExport(ctx, &sonar.IssuesGitlabSastExportOptions{ProjectKey: projectKey})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Version).NotTo(BeEmpty())
			Expect(result.Scan.Type).To(Equal("sast"))
		})

		It("should export all findings", func() {
			result, _, err := client.Projects.ExportFindings(ctx, &sonar.ProjectsExportFindingsOptions{Project: projectKey})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
		})
	})

	Describe("Misc V1 endpoints", func() {
		It("should return the application quality gate status", func() {
			result, _, err := client.Qualitygates.ApplicationStatus(ctx, &sonar.QualitygatesApplicationStatusOptions{Application: appKey})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Status).NotTo(BeEmpty())
		})

		It("should validate the project binding", func() {
			// the fake GitHub App private key is invalid, so the server answers 400 with the validation errors
			_, resp, err := client.AlmSettings.ValidateBinding(ctx, &sonar.AlmSettingsValidateBindingOptions{Project: projectKey})
			Expect(err).To(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			Expect(err.Error()).To(ContainSubstring("private key"))
		})

		It("should report the GitHub provisioning status", func() {
			result, _, err := client.GithubProvisioning.Status(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Enabled).To(BeFalse())
		})

		It("should set the Compute Engine worker count", func() {
			current, _, err := client.Ce.WorkerCount(ctx)
			Expect(err).NotTo(HaveOccurred())

			resp, err := client.Ce.SetWorkerCount(ctx, &sonar.CeSetWorkerCountOptions{Count: current.Value})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})
	})

	Describe("Project bindings", func() {
		It("should search and fetch the project binding", func() {
			list, _, err := client.V2.DopTranslation.SearchProjectBindings(ctx, &sonar.DopTranslationProjectBindingsSearchOptions{Repository: "org/repo"})
			Expect(err).NotTo(HaveOccurred())

			var binding *sonar.DopTranslationProjectBinding

			for i := range list.ProjectBindings {
				if list.ProjectBindings[i].ProjectKey == projectKey {
					binding = &list.ProjectBindings[i]
				}
			}

			Expect(binding).NotTo(BeNil())
			Expect(binding.DevOpsPlatformSettingId).To(Equal(dopSettings))

			fetched, _, err := client.V2.DopTranslation.GetProjectBinding(ctx, binding.Id)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetched.ProjectKey).To(Equal(projectKey))
		})
	})

	Describe("GitHub configuration", Ordered, func() {
		It("should create, update, fetch and delete a configuration", func() {
			created, _, err := client.V2.DopTranslation.CreateGithubConfiguration(ctx, &sonar.DopTranslationGithubConfigurationCreateOptions{
				ClientId: "cid", ClientSecret: "secret", ApplicationId: "1", PrivateKey: "pk",
				ApiUrl: "https://api.github.com", WebUrl: "https://github.com",
				AllowedOrganizations: []string{"org"}, ProvisioningType: sonar.DopProvisioningTypeJIT,
			})
			Expect(err).NotTo(HaveOccurred())
			DeferCleanup(func() {
				_, _ = client.V2.DopTranslation.DeleteGithubConfiguration(ctx, created.Id)
			})

			updated, _, err := client.V2.DopTranslation.UpdateGithubConfiguration(ctx, created.Id, &sonar.DopTranslationGithubConfigurationUpdateOptions{
				AllowedOrganizations: &sonar.UpdateFieldListStringV2{Value: []string{"org", "other"}, Defined: true},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.AllowedOrganizations).To(ConsistOf("org", "other"))

			list, _, err := client.V2.DopTranslation.SearchGithubConfigurations(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(list.GithubConfigurations).NotTo(BeEmpty())

			fetched, _, err := client.V2.DopTranslation.GetGithubConfiguration(ctx, created.Id)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetched.ApplicationId).To(Equal("1"))

			resp, err := client.V2.DopTranslation.DeleteGithubConfiguration(ctx, created.Id)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})
	})

	Describe("GitLab configuration", Ordered, func() {
		It("should create, update, fetch and delete a configuration", func() {
			created, _, err := client.V2.DopTranslation.CreateGitlabConfiguration(ctx, &sonar.DopTranslationGitlabConfigurationCreateOptions{
				ApplicationId: "1", Url: "https://gitlab.com", Secret: "secret",
				AllowedGroups: []string{"group"}, ProvisioningType: sonar.DopProvisioningTypeJIT,
			})
			Expect(err).NotTo(HaveOccurred())
			DeferCleanup(func() {
				_, _ = client.V2.DopTranslation.DeleteGitlabConfiguration(ctx, created.Id)
			})

			updated, _, err := client.V2.DopTranslation.UpdateGitlabConfiguration(ctx, created.Id, &sonar.DopTranslationGitlabConfigurationUpdateOptions{
				AllowedGroups: &sonar.UpdateFieldListStringV2{Value: []string{"g1", "g2"}, Defined: true},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.AllowedGroups).To(ConsistOf("g1", "g2"))

			list, _, err := client.V2.DopTranslation.SearchGitlabConfigurations(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(list.GitlabConfigurations).NotTo(BeEmpty())

			fetched, _, err := client.V2.DopTranslation.GetGitlabConfiguration(ctx, created.Id)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetched.Url).To(Equal("https://gitlab.com"))

			resp, err := client.V2.DopTranslation.DeleteGitlabConfiguration(ctx, created.Id)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})
	})

	Describe("Permission mappings", func() {
		It("should manage custom GitHub and GitLab role mappings", func() {
			role := helpers.UniqueResourceName("role")
			scan := true

			githubList, _, err := client.V2.DopTranslation.SearchGithubPermissionMappings(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(githubList.PermissionMappings).NotTo(BeEmpty())

			gitlabList, _, err := client.V2.DopTranslation.SearchGitlabPermissionMappings(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(gitlabList.PermissionMappings).NotTo(BeEmpty())

			create := &sonar.DopTranslationPermissionMappingCreateOptions{Role: role, Permissions: sonar.DopTranslationPermissions{User: true, CodeViewer: true}}
			update := &sonar.DopTranslationPermissionMappingUpdateOptions{Permissions: &sonar.DopTranslationPermissionsUpdate{Scan: &scan}}

			github, _, err := client.V2.DopTranslation.CreateGithubPermissionMapping(ctx, create)
			Expect(err).NotTo(HaveOccurred())
			Expect(github.BaseRole).To(BeFalse())

			github, _, err = client.V2.DopTranslation.UpdateGithubPermissionMapping(ctx, role, update)
			Expect(err).NotTo(HaveOccurred())
			Expect(github.Permissions.Scan).To(BeTrue())

			_, err = client.V2.DopTranslation.DeleteGithubPermissionMapping(ctx, role)
			Expect(err).NotTo(HaveOccurred())

			gitlab, _, err := client.V2.DopTranslation.CreateGitlabPermissionMapping(ctx, create)
			Expect(err).NotTo(HaveOccurred())
			Expect(gitlab.BaseRole).To(BeFalse())

			gitlab, _, err = client.V2.DopTranslation.UpdateGitlabPermissionMapping(ctx, role, update)
			Expect(err).NotTo(HaveOccurred())
			Expect(gitlab.Permissions.Scan).To(BeTrue())

			_, err = client.V2.DopTranslation.DeleteGitlabPermissionMapping(ctx, role)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
