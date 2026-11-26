package enterprise_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

var _ = Describe("Request encoding and response fields", Ordered, func() {
	var (
		client  *sonar.Client
		cleanup *helpers.CleanupManager
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())

		edition, err := helpers.GetEdition(client)
		Expect(err).NotTo(HaveOccurred())

		if !helpers.IsEnterpriseOrAbove(edition) {
			Skip("requires Enterprise Edition")
		}

		cleanup = helpers.NewCleanupManager(client)
	})

	AfterAll(func() {
		if cleanup == nil {
			return
		}

		for _, err := range cleanup.Cleanup() {
			GinkgoWriter.Printf("Cleanup error: %v\n", err)
		}
	})

	It("should update and delete a permission mapping whose role contains a space", func() {
		role := helpers.UniqueResourceName("custom role")

		_, _, err := client.V2.DopTranslation.CreateGithubPermissionMapping(context.Background(), &sonar.DopTranslationPermissionMappingCreateOptions{
			Role:        role,
			Permissions: sonar.DopTranslationPermissions{User: true},
		})
		Expect(err).NotTo(HaveOccurred())
		cleanup.RegisterCleanup("permission mapping", role, func() error {
			_, err := client.V2.DopTranslation.DeleteGithubPermissionMapping(context.Background(), role)
			return helpers.IgnoreNotFoundError(err)
		})

		scan := true
		updated, _, err := client.V2.DopTranslation.UpdateGithubPermissionMapping(context.Background(), role, &sonar.DopTranslationPermissionMappingUpdateOptions{
			Permissions: &sonar.DopTranslationPermissionsUpdate{Scan: &scan},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(updated.Role).To(Equal(role))

		resp, err := client.V2.DopTranslation.DeleteGithubPermissionMapping(context.Background(), role)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
	})

	It("should send the instance-wide scheduled agent config ID in a form the server accepts", func() {
		_, resp, err := client.V2.RemediationAgent.GetScheduledAgentConfig(context.Background(), "*")
		// Without a Remediation Agent entitlement the server answers 403. A
		// doubly escaped ID is rejected earlier, with a 400.
		if err != nil {
			Expect(resp).NotTo(BeNil())
			Expect(resp.StatusCode).NotTo(Equal(http.StatusBadRequest))
		}
	})

	It("should remove every tag of an application", func() {
		key := helpers.UniqueResourceName("fix-application")

		_, _, err := client.Applications.Create(context.Background(), &sonar.ApplicationsCreateOptions{Name: key, Key: key})
		Expect(err).NotTo(HaveOccurred())
		cleanup.RegisterCleanup("application", key, func() error {
			_, err := client.Applications.Delete(context.Background(), &sonar.ApplicationsDeleteOptions{Application: key})
			return helpers.IgnoreNotFoundError(err)
		})

		for _, tags := range [][]string{{"finance"}, {}} {
			_, err = client.Applications.SetTags(context.Background(), &sonar.ApplicationsSetTagsOptions{Application: key, Tags: tags})
			Expect(err).NotTo(HaveOccurred())

			shown, _, err := client.Applications.Show(context.Background(), &sonar.ApplicationsShowOptions{Application: key})
			Expect(err).NotTo(HaveOccurred())
			Expect(shown.Application.Tags).To(HaveLen(len(tags)))
		}
	})

	It("should report the portfolio frequency and application branches, and clear recipients", func() {
		portfolio := helpers.UniqueResourceName("fix-portfolio")
		application := helpers.UniqueResourceName("fix-member")
		project := helpers.UniqueResourceName("fix-project")

		_, _, err := client.Projects.Create(context.Background(), &sonar.ProjectsCreateOptions{Name: project, Project: project})
		Expect(err).NotTo(HaveOccurred())
		cleanup.RegisterCleanup("project", project, func() error {
			_, err := client.Projects.Delete(context.Background(), &sonar.ProjectsDeleteOptions{Project: project})
			return helpers.IgnoreNotFoundError(err)
		})

		_, err = client.Views.Create(context.Background(), &sonar.ViewsCreateOptions{Name: portfolio, Key: portfolio})
		Expect(err).NotTo(HaveOccurred())
		cleanup.RegisterCleanup("portfolio", portfolio, func() error {
			_, err := client.Views.Delete(context.Background(), &sonar.ViewsDeleteOptions{Key: portfolio})
			return helpers.IgnoreNotFoundError(err)
		})

		_, _, err = client.Applications.Create(context.Background(), &sonar.ApplicationsCreateOptions{Name: application, Key: application})
		Expect(err).NotTo(HaveOccurred())
		cleanup.RegisterCleanup("application", application, func() error {
			_, err := client.Applications.Delete(context.Background(), &sonar.ApplicationsDeleteOptions{Application: application})
			return helpers.IgnoreNotFoundError(err)
		})

		_, err = client.Applications.AddProject(context.Background(), &sonar.ApplicationsAddProjectOptions{Application: application, Project: project})
		Expect(err).NotTo(HaveOccurred())
		_, err = client.Views.AddApplication(context.Background(), &sonar.ViewsAddApplicationOptions{Portfolio: portfolio, Application: application})
		Expect(err).NotTo(HaveOccurred())
		_, err = client.Views.AddApplicationBranch(context.Background(), &sonar.ViewsAddApplicationBranchOptions{Key: portfolio, Application: application, Branch: "main"})
		Expect(err).NotTo(HaveOccurred())

		details, _, err := client.Views.Show(context.Background(), &sonar.ViewsShowOptions{Key: portfolio})
		Expect(err).NotTo(HaveOccurred())
		Expect(details.SubViews).To(HaveLen(1))
		Expect(details.SubViews[0].SelectedBranches).To(ConsistOf("main"))

		_, err = client.GovernanceReports.UpdateFrequency(context.Background(), &sonar.GovernanceReportsUpdateFrequencyOptions{ComponentKey: portfolio, Frequency: "weekly"})
		Expect(err).NotTo(HaveOccurred())

		for _, recipients := range []string{"reports@example.com", ""} {
			_, err = client.GovernanceReports.UpdateRecipients(context.Background(), &sonar.GovernanceReportsUpdateRecipientsOptions{ComponentKey: portfolio, Recipients: recipients})
			Expect(err).NotTo(HaveOccurred())
		}

		status, _, err := client.GovernanceReports.Status(context.Background(), &sonar.GovernanceReportsStatusOptions{ComponentKey: portfolio})
		Expect(err).NotTo(HaveOccurred())
		Expect(status.ComponentFrequency).To(Equal("weekly"))
		Expect(status.ComponentRecipients).To(BeEmpty())
	})

	It("should update the sandbox deferral days", func() {
		before, _, err := client.V2.Issues.GetSandboxSettings(context.Background())
		Expect(err).NotTo(HaveOccurred())

		days := int64(3)
		_, _, err = client.V2.Issues.UpdateSandboxSettings(context.Background(), &sonar.IssuesV2UpdateSandboxSettingsOptions{DeferralDays: &days})
		Expect(err).NotTo(HaveOccurred())
		cleanup.RegisterCleanup("sandbox settings", "instance", func() error {
			_, _, err := client.V2.Issues.UpdateSandboxSettings(context.Background(), &sonar.IssuesV2UpdateSandboxSettingsOptions{DeferralDays: &before.DeferralDays})
			return err
		})

		after, _, err := client.V2.Issues.GetSandboxSettings(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(after.DeferralDays).To(Equal(days))
	})

	It("should send the AI CodeFix provider with the enablement", func() {
		resp, err := client.V2.FixSuggestions.SetEnablement(context.Background(), &sonar.FixSuggestionsSetEnablementOptions{
			Enablement: "ENABLED_FOR_ALL_PROJECTS",
			Provider: &sonar.FixSuggestionsProviderOptions{
				Type:   "CUSTOM",
				Config: map[string]string{"endpoint": "https://llm.invalid", "modelId": "model"},
			},
		})
		// The server tries to reach the endpoint, which proves it read the
		// provider. Without it the request fails on the missing field.
		Expect(err).To(HaveOccurred())
		Expect(resp).NotTo(BeNil())
		Expect(err.Error()).NotTo(ContainSubstring("must not be empty"))
	})
})
