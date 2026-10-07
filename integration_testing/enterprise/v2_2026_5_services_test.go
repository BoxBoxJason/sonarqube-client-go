package enterprise_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

// newUUID returns a random RFC 4122 version 4 UUID.
func newUUID() string {
	GinkgoHelper()

	b := make([]byte, 16)
	_, err := rand.Read(b)
	Expect(err).NotTo(HaveOccurred())
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// Specs for the V2 services introduced with SonarQube 2026.5. They run against
// an unlicensed Enterprise Edition instance, so endpoints that need a license,
// a reachable Agent Orchestrator or system-to-system authentication are only
// exercised up to the error they return.
var _ = Describe("SonarQube 2026.5 V2 services", Ordered, func() {
	var (
		client       *sonar.Client
		cleanup      *helpers.CleanupManager
		projectKey   string
		projectID    string
		branchID     string
		portfolioKey string
		portfolioID  string
		from, to     string
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		cleanup = helpers.NewCleanupManager(client)

		projectKey = helpers.UniqueResourceName("v2-2026-5")
		_, _, err = client.Projects.Create(context.Background(), &sonar.ProjectsCreateOptions{Name: projectKey, Project: projectKey})
		Expect(err).NotTo(HaveOccurred())
		cleanup.RegisterCleanup("project", projectKey, func() error {
			_, err := client.Projects.Delete(context.Background(), &sonar.ProjectsDeleteOptions{Project: projectKey})
			return err
		})
		projectID = projectUUID(client, projectKey)

		branches, _, err := client.ProjectBranches.List(context.Background(), &sonar.ProjectBranchesListOptions{Project: projectKey})
		Expect(err).NotTo(HaveOccurred())
		Expect(branches.Branches).NotTo(BeEmpty())
		branchID = branches.Branches[0].BranchID

		portfolioKey = helpers.UniqueResourceName("v2-2026-5-pf")
		_, err = client.Views.Create(context.Background(), &sonar.ViewsCreateOptions{Key: portfolioKey, Name: portfolioKey})
		Expect(err).NotTo(HaveOccurred())
		cleanup.RegisterCleanup("portfolio", portfolioKey, func() error {
			_, err := client.Views.Delete(context.Background(), &sonar.ViewsDeleteOptions{Key: portfolioKey})
			return err
		})
		portfolio, _, err := client.Navigation.Component(context.Background(), &sonar.NavigationComponentOptions{Component: portfolioKey})
		Expect(err).NotTo(HaveOccurred())
		portfolioID = portfolio.ID

		now := time.Now().UTC()
		from = now.AddDate(0, 0, -7).Format(time.DateOnly)
		to = now.Format(time.DateOnly)
	})

	AfterAll(func() {
		for _, err := range cleanup.Cleanup() {
			GinkgoWriter.Printf("Cleanup error: %v\n", err)
		}
	})

	// =========================================================================
	// Dashboards
	// =========================================================================
	Describe("Dashboards", Ordered, func() {
		var dashboardID string

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

		It("should create, get, list, update and delete a custom dashboard", func() {
			created, resp, err := client.V2.Dashboards.Create(context.Background(), &sonar.DashboardsCreateOptions{
				Name:         "e2e dashboard",
				Layout:       `{"widgets":[]}`,
				ResourceType: "project",
				ResourceId:   projectID,
				Description:  "created by e2e",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))
			Expect(created.ResourceId).To(Equal(projectID))
			dashboardID = created.Id

			got, _, err := client.V2.Dashboards.Get(context.Background(), dashboardID)
			Expect(err).NotTo(HaveOccurred())
			Expect(got.Name).To(Equal("e2e dashboard"))

			list, _, err := client.V2.Dashboards.List(context.Background(), &sonar.DashboardsListOptions{ResourceId: projectID, ResourceType: "project"})
			Expect(err).NotTo(HaveOccurred())
			Expect(list.Dashboards).To(HaveLen(1))

			updated, _, err := client.V2.Dashboards.Update(context.Background(), dashboardID, &sonar.DashboardsUpdateOptions{Name: new("e2e dashboard v2")})
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.Name).To(Equal("e2e dashboard v2"))

			resp, err = client.V2.Dashboards.Delete(context.Background(), dashboardID)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})
	})

	// =========================================================================
	// Detection agent
	// =========================================================================
	Describe("Detection agent", func() {
		It("should return the instance configuration", func() {
			result, resp, err := client.V2.DetectionAgent.GetInstanceConfiguration(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result.MaxConcurrentJobs).To(BeNumerically(">", 0))
		})

		It("should read, enable with a full schedule, and disable the project configuration", func() {
			result, _, err := client.V2.DetectionAgent.GetProjectConfiguration(context.Background(), projectID)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.ProjectId).To(Equal(projectID))
			Expect(result.Scheduling.Enabled).To(BeFalse())

			result, resp, err := client.V2.DetectionAgent.UpsertProjectConfiguration(context.Background(), projectID, &sonar.DetectionAgentUpsertProjectConfigurationOptions{
				Scheduling: &sonar.DetectionAgentSchedulingRequest{
					Enabled:    new(true),
					BranchId:   branchID,
					Time:       &sonar.DetectionAgentScheduleTime{Hour: new(int32(0)), Timezone: "Europe/Paris"},
					Recurrence: &sonar.DetectionAgentScheduleRecurrence{Type: "WEEKLY", DayOfWeek: "MONDAY"},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result.Scheduling.Enabled).To(BeTrue())
			Expect(result.Scheduling.BranchId).To(Equal(branchID))
			Expect(result.Scheduling.NextRunTime).NotTo(BeEmpty())
			Expect(result.Scheduling.Recurrence).To(Equal(&sonar.DetectionAgentScheduleRecurrence{Type: "WEEKLY", DayOfWeek: "MONDAY"}))
			Expect(result.Scheduling.Time).NotTo(BeNil())
			Expect(*result.Scheduling.Time.Hour).To(BeEquivalentTo(0))
			Expect(result.Scheduling.Time.Timezone).To(Equal("Europe/Paris"))

			result, _, err = client.V2.DetectionAgent.UpsertProjectConfiguration(context.Background(), projectID, &sonar.DetectionAgentUpsertProjectConfigurationOptions{
				Scheduling: &sonar.DetectionAgentSchedulingRequest{Enabled: new(false)},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Scheduling.Enabled).To(BeFalse())
		})
	})

	// =========================================================================
	// History
	// =========================================================================
	Describe("History", func() {
		It("should return issue density and resolution history for a branch", func() {
			density, resp, err := client.V2.History.GetIssueDensityHistory(context.Background(), &sonar.HistoryGetIssueDensityHistoryOptions{
				EntityId: branchID, EntityType: "PROJECT_BRANCH", StartDate: from + "T00:00:00Z", SliceBy: "SEVERITY",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(density).NotTo(BeNil())

			resolution, _, err := client.V2.History.GetIssueResolutionHistory(context.Background(), &sonar.HistoryGetIssueResolutionHistoryOptions{
				EntityId: branchID, EntityType: "PROJECT_BRANCH", Statistic: "MTTR", StartDate: from + "T00:00:00Z",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resolution.Statistic).To(Equal("MTTR"))

			sca, _, err := client.V2.History.GetScaResolutionHistory(context.Background(), &sonar.HistoryGetScaResolutionHistoryOptions{
				EntityId: branchID, EntityType: "PROJECT_BRANCH", Statistic: "SCA_MTTR", StartDate: from + "T00:00:00Z",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(sca.Statistic).To(Equal("SCA_MTTR"))
			Expect(sca.ScaResolutionHistory).NotTo(BeEmpty())
		})

		It("should return project issue and SCA resolution for a portfolio", func() {
			issues, resp, err := client.V2.History.GetProjectIssueResolution(context.Background(), &sonar.HistoryGetProjectIssueResolutionOptions{
				Statistic: "RESOLVED_ISSUES", EntityType: "PORTFOLIO", EntityId: portfolioID,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(issues.Statistic).To(Equal("RESOLVED_ISSUES"))

			sca, _, err := client.V2.History.GetProjectScaResolution(context.Background(), &sonar.HistoryGetProjectScaResolutionOptions{
				Statistic: "SCA_MTTR", EntityType: "PORTFOLIO", EntityId: portfolioID,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(sca.Statistic).To(Equal("SCA_MTTR"))
		})
	})

	// =========================================================================
	// Software quality reports, onboarding, security alerts, agentic jobs
	// =========================================================================
	It("should return a WCAG compliance report with nested categories", func() {
		report, resp, err := client.V2.SoftwareQualityReports.GetComplianceReport(context.Background(), &sonar.SoftwareQualityReportsGetComplianceReportOptions{
			ComponentKey: projectKey, Standard: "wcag", Version: "2.1", Level: "AA",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(report.Categories).NotTo(BeEmpty())
		Expect(report.Categories[0].Categories).NotTo(BeEmpty())
	})

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

	It("should search security alerts", func() {
		result, resp, err := client.V2.SecurityAlerts.Search(context.Background(), &sonar.SecurityAlertsSearchOptions{
			AlertTypes: []string{"DEPENDENCY_RISK"}, Statuses: []string{"OPEN"}, Sort: "FIRST_DETECTED_AT", Direction: "ASC",
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(result).NotTo(BeNil())

		_, resp, err = client.V2.SecurityAlerts.Get(context.Background(), "does-not-exist")
		Expect(err).To(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
	})

	It("should search agentic jobs", func() {
		result, resp, err := client.V2.Agentic.SearchJobs(context.Background(), &sonar.AgenticSearchJobsOptions{Type: "HUNTER"})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(result).NotTo(BeNil())
	})

	It("should check DevOps Platform permissions", func() {
		result, resp, err := client.V2.DopTranslation.CheckPermissions(context.Background(), &sonar.DopTranslationCheckPermissionsOptions{})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(result).NotTo(BeNil())

		// The project is not bound to any DevOps Platform, so no token can be minted.
		_, resp, err = client.V2.DopTranslation.GenerateScmAccessToken(context.Background(), &sonar.DopTranslationGenerateScmAccessTokenOptions{Project: projectKey})
		Expect(err).To(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
	})

	// =========================================================================
	// LLM connectivity
	// =========================================================================
	It("should list LLM provider definitions, mappings and providers", func() {
		definitions, resp, err := client.V2.LlmConnectivity.GetProviderDefinitions(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(definitions.ProviderDefinitions).NotTo(BeEmpty())
		Expect(definitions.ProviderDefinitions[0].Fields).NotTo(BeEmpty())

		mappings, _, err := client.V2.LlmConnectivity.GetProviderMappings(context.Background(), &sonar.LlmConnectivityGetProviderMappingsOptions{AiCapability: "AI_CODEFIX"})
		Expect(err).NotTo(HaveOccurred())
		Expect(mappings.ProviderMappings).To(HaveLen(1))
		Expect(mappings.ProviderMappings[0].AiCapability).To(Equal("AI_CODEFIX"))

		providers, _, err := client.V2.LlmConnectivity.ListProviders(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(providers).NotTo(BeNil())

		_, resp, err = client.V2.LlmConnectivity.GetProvider(context.Background(), "does-not-exist")
		Expect(err).To(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
	})

	// =========================================================================
	// CAG, billing and A3S
	// =========================================================================
	It("should answer the CAG ping, entitlement and impact endpoints", func() {
		ping, resp, err := client.V2.Cag.Ping(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(ping.Status).To(Equal("ok"))

		entitlement, _, err := client.V2.Cag.GetEntitlement(context.Background(), defaultOrganizationID)
		Expect(err).NotTo(HaveOccurred())
		Expect(entitlement.FeatureKey).To(Equal("contextAugmentation"))

		usage, _, err := client.V2.Cag.GetUsageStats(context.Background(), defaultOrganizationID, &sonar.CagGetUsageStatsOptions{Period: "LAST_7_DAYS"})
		Expect(err).NotTo(HaveOccurred())
		Expect(usage.Id).To(Equal(defaultOrganizationID))

		guide, _, err := client.V2.Cag.GetGuideMetrics(context.Background(), &sonar.CagGetGuideMetricsOptions{OrganizationId: defaultOrganizationID, From: from, To: to})
		Expect(err).NotTo(HaveOccurred())
		Expect(guide).NotTo(BeNil())

		verify, _, err := client.V2.Cag.GetVerifyMetrics(context.Background(), &sonar.CagGetVerifyMetricsOptions{OrganizationId: defaultOrganizationID, From: from, To: to})
		Expect(err).NotTo(HaveOccurred())
		Expect(verify).NotTo(BeNil())

		activity, _, err := client.V2.Cag.GetProjectActivity(context.Background(), &sonar.CagGetProjectActivityOptions{OrganizationId: defaultOrganizationID, From: from, To: to})
		Expect(err).NotTo(HaveOccurred())
		Expect(activity.TotalProjectCount).To(BeNumerically(">", 0))

		impact, _, err := client.V2.Cag.GetProjectImpact(context.Background(), &sonar.CagGetProjectImpactOptions{OrganizationId: defaultOrganizationID, From: from, To: to})
		Expect(err).NotTo(HaveOccurred())
		Expect(impact.Projects).NotTo(BeEmpty())
	})

	It("should record CAG usage and impact events", func() {
		resp, err := client.V2.Cag.RecordUsageEvent(context.Background(), &sonar.CagRecordUsageEventOptions{
			OrganizationId: defaultOrganizationID, ProjectId: projectID, InvocationId: newUUID(),
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusAccepted))

		resp, err = client.V2.Cag.RecordImpactEvent(context.Background(), &sonar.CagRecordImpactEventOptions{
			OrganizationId: defaultOrganizationID,
			ProjectId:      projectID,
			InvocationId:   newUUID(),
			CagInstanceId:  "e2e",
			EventType:      "NAVIGATION",
			EventVersion:   1,
			Success:        true,
			Transport:      "CLI",
			Payload:        map[string]any{"toolName": "get_source_code", "resultBytes": 10, "resultCount": 1, "truncated": false},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusAccepted))
	})

	It("should check billing entitlements and record consumption", func() {
		entitlement, resp, err := client.V2.Billing.GetEntitlementCheck(context.Background(), &sonar.BillingGetEntitlementCheckOptions{FeatureKey: "contextAugmentation"})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(entitlement.FeatureKey).To(Equal("contextAugmentation"))

		record, resp, err := client.V2.Billing.CreateConsumptionRecord(context.Background(), newUUID(), &sonar.BillingCreateConsumptionRecordOptions{
			FeatureKey: "contextAugmentation", Consumed: new(int64(1)),
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusAccepted))
		Expect(record.Consumed).To(BeEquivalentTo(1))
		Expect(record.ResourceId).To(Equal(defaultOrganizationID))
	})

	It("should answer the public A3S statistics, collection config and entitlement endpoints", func() {
		stats, resp, err := client.V2.A3s.GetAnalysisStats(context.Background(), defaultOrganizationID, &sonar.A3sGetAnalysisStatsOptions{Period: "LAST_7_DAYS"})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
		Expect(stats.Projects).NotTo(BeEmpty())

		config, _, err := client.V2.A3s.GetPublicCollectionConfig(context.Background(), &sonar.A3sGetPublicCollectionConfigOptions{
			OrganizationKey: "default-organization", ProjectKey: projectKey,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(config).NotTo(BeNil())

		entitlement, _, err := client.V2.A3s.GetPublicOrgEntitlement(context.Background(), defaultOrganizationID)
		Expect(err).NotTo(HaveOccurred())
		Expect(entitlement.Id).To(Equal(defaultOrganizationID))
	})
})
