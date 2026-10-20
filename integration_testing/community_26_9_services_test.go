package integration_testing_test

import (
	"context"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

// Specs for the endpoints that SonarQube Community Build 26.9 gained from
// Enterprise Edition 2026.5. They assert the behavior shared by both editions,
// so they run unchanged against either one.
var _ = Describe("SonarQube 26.9 Community Build services", Ordered, func() {
	var (
		client     *sonar.Client
		cleanup    *helpers.CleanupManager
		projectKey string
		branchID   string
		startDate  string
	)

	// A syntactically valid (but non-existent) 36-character entity UUID.
	const dummyEntityID = "00000000-0000-0000-0000-000000000000"

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		cleanup = helpers.NewCleanupManager(client)

		projectKey = helpers.UniqueResourceName("ce-26-9")
		_, _, err = client.Projects.Create(context.Background(), &sonar.ProjectsCreateOptions{Name: projectKey, Project: projectKey})
		Expect(err).NotTo(HaveOccurred())
		cleanup.RegisterCleanup("project", projectKey, func() error {
			_, err := client.Projects.Delete(context.Background(), &sonar.ProjectsDeleteOptions{Project: projectKey})
			return err
		})

		branches, _, err := client.ProjectBranches.List(context.Background(), &sonar.ProjectBranchesListOptions{Project: projectKey})
		Expect(err).NotTo(HaveOccurred())
		Expect(branches.Branches).NotTo(BeEmpty())
		branchID = branches.Branches[0].BranchID

		startDate = time.Now().UTC().AddDate(0, 0, -7).Format(time.RFC3339)
	})

	AfterAll(func() {
		for _, err := range cleanup.Cleanup() {
			GinkgoWriter.Printf("Cleanup error: %v\n", err)
		}
	})

	// =========================================================================
	// Built-in dashboards
	// =========================================================================
	It("should list built-in dashboards and get one", func() {
		list, resp, err := client.V2.Dashboards.ListBuiltIn(context.Background(), &sonar.DashboardsListBuiltInOptions{ResourceType: "project"})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(list.Dashboards).NotTo(BeEmpty())

		builtIn, _, err := client.V2.Dashboards.GetBuiltIn(context.Background(), list.Dashboards[0].Key)
		Expect(err).NotTo(HaveOccurred())
		Expect(builtIn.Key).To(Equal(list.Dashboards[0].Key))
		Expect(builtIn.Layout).NotTo(BeEmpty())
	})

	// =========================================================================
	// History
	// =========================================================================
	Describe("History", func() {
		It("should return issue density history for a branch", func() {
			result, resp, err := client.V2.History.GetIssueDensityHistory(context.Background(), &sonar.HistoryGetIssueDensityHistoryOptions{
				EntityId: branchID, EntityType: "PROJECT_BRANCH", StartDate: startDate, SliceBy: "SEVERITY",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		It("should return issue resolution history for a branch", func() {
			for _, statistic := range []string{"MTTR", "RECENT_MTTR", "RESOLVED_ISSUES"} {
				result, resp, err := client.V2.History.GetIssueResolutionHistory(context.Background(), &sonar.HistoryGetIssueResolutionHistoryOptions{
					EntityId: branchID, EntityType: "PROJECT_BRANCH", Statistic: statistic, StartDate: startDate,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result.Statistic).To(Equal(statistic))
			}
		})

		It("should accept type severities on issue count history", func() {
			result, resp, err := client.V2.History.GetIssueCountHistory(context.Background(), &sonar.HistoryIssueCountHistoryOptions{
				EntityID: branchID, EntityType: sonar.HistoryEntityTypeProjectBranch, StartDate: startDate, TypeSeverities: []string{"MAJOR"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		It("should accept the APPLICATION entity type on measures history", func() {
			_, resp, err := client.V2.History.GetMeasuresHistory(context.Background(), &sonar.HistoryMeasuresHistoryOptions{
				EntityType: sonar.HistoryEntityTypeApplication, EntityID: dummyEntityID, MetricKeys: []string{"ncloc"}, StartDate: startDate,
			})
			Expect(err).To(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("should reject an unknown slice dimension client-side", func() {
			result, resp, err := client.V2.History.GetIssueCountHistory(context.Background(), &sonar.HistoryIssueCountHistoryOptions{
				EntityID: branchID, EntityType: sonar.HistoryEntityTypeProjectBranch, StartDate: startDate, SliceBy: "issueType",
			})
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("should return not found for the project issue resolution of an unknown portfolio", func() {
			_, resp, err := client.V2.History.GetProjectIssueResolution(context.Background(), &sonar.HistoryGetProjectIssueResolutionOptions{
				Statistic: "MTTR", PortfolioId: helpers.UniqueResourceName("portfolio"),
			})
			Expect(err).To(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})
	})

	// =========================================================================
	// Onboarding
	// =========================================================================
	It("should return the onboarding overview, projects and statistics", func() {
		overview, resp, err := client.V2.Onboarding.GetOverview(context.Background(), nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(overview).NotTo(BeNil())

		projects, _, err := client.V2.Onboarding.GetProjects(context.Background(), &sonar.OnboardingGetProjectsOptions{Q: projectKey})
		Expect(err).NotTo(HaveOccurred())
		Expect(projects.Projects).To(HaveLen(1))
		Expect(projects.Projects[0].Key).To(Equal(projectKey))
		Expect(projects.Projects[0].ScanStatus).To(Equal("NOT_SCANNED"))

		stats, _, err := client.V2.Onboarding.GetStatistics(context.Background(), nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(stats.DevopsPlatforms.Total).To(BeNumerically(">", 0))
	})

	// =========================================================================
	// Agentic jobs
	// =========================================================================
	It("should search agentic jobs", func() {
		result, resp, err := client.V2.Agentic.SearchJobs(context.Background(), &sonar.AgenticSearchJobsOptions{Type: "HUNTER"})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(result).NotTo(BeNil())
	})

	// =========================================================================
	// DevOps Platform permissions and SCM access tokens
	// =========================================================================
	Describe("DevOps Platform permissions", func() {
		It("should check DevOps Platform permissions", func() {
			result, resp, err := client.V2.DopTranslation.CheckPermissions(context.Background(), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())

			result, _, err = client.V2.DopTranslation.CheckPermissions(context.Background(), &sonar.DopTranslationCheckPermissionsOptions{Project: projectKey})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
		})

		It("should reject a refresh without a configuration client-side", func() {
			result, resp, err := client.V2.DopTranslation.CheckPermissions(context.Background(), &sonar.DopTranslationCheckPermissionsOptions{Refresh: new(true)})
			Expect(err).To(HaveOccurred())
			Expect(resp).To(BeNil())
			Expect(result).To(BeNil())
		})

		It("should not mint an SCM access token for an unbound project", func() {
			_, resp, err := client.V2.DopTranslation.GenerateScmAccessToken(context.Background(), &sonar.DopTranslationGenerateScmAccessTokenOptions{Project: projectKey})
			Expect(err).To(HaveOccurred())
			Expect(resp).NotTo(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})
	})

	// =========================================================================
	// Hotspots migration
	// =========================================================================
	Describe("Hotspots migration", func() {
		It("should report the migration status of the instance", func() {
			result, resp, err := client.Hotspots.MigrationStatus(context.Background(), nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.RemainingHotspots).To(BeNumerically(">=", 0))
		})

		It("should report a complete migration for a project without hotspots", func() {
			result, resp, err := client.Hotspots.MigrationStatus(context.Background(), &sonar.HotspotsMigrationStatusOptions{Project: projectKey})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result.RemainingHotspots).To(BeZero())
			Expect(result.Complete).To(BeTrue())
		})

		It("should dry-run a migration without writing anything", func() {
			result, resp, err := client.Hotspots.MigrateToIssues(context.Background(), &sonar.HotspotsMigrateToIssuesOptions{
				Project: projectKey,
				DryRun:  true,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result.DryRun).To(BeTrue())
			Expect(result.Skipped).To(BeZero())
			Expect(result.Projects).To(BeEmpty())
		})
	})
})
