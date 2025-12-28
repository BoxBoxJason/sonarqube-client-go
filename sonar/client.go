package sonargo

import (
	"net/http"
	"net/url"
)

type Client struct {
	baseURL                   *url.URL
	username, password, token string
	authType                  authType
	httpClient                *http.Client
	AlmIntegrations           *AlmIntegrationsService
	AlmSettings               *AlmSettingsService
	AnalysisCache             *AnalysisCacheService
	AnalysisReports           *AnalysisReportsService
	Authentication            *AuthenticationService
	Ce                        *CeService
	Components                *ComponentsService
	Developers                *DevelopersService
	DismissMessage            *DismissMessageService
	Duplications              *DuplicationsService
	Emails                    *EmailsService
	Favorites                 *FavoritesService
	Features                  *FeaturesService
	GithubProvisioning        *GithubProvisioningService
	Hotspots                  *HotspotsService
	Issues                    *IssuesService
	L10N                      *L10NService
	Languages                 *LanguagesService
	Measures                  *MeasuresService
	Metrics                   *MetricsService
	Monitoring                *MonitoringService
	Navigation                *NavigationService
	NewCodePeriods            *NewCodePeriodsService
	Notifications             *NotificationsService
	Permissions               *PermissionsService
	Plugins                   *PluginsService
	ProjectAnalyses           *ProjectAnalysesService
	ProjectBadges             *ProjectBadgesService
	ProjectBranches           *ProjectBranchesService
	ProjectDump               *ProjectDumpService
	ProjectLinks              *ProjectLinksService
	ProjectTags               *ProjectTagsService
	Projects                  *ProjectsService
	Push                      *PushService
	Qualitygates              *QualitygatesService
	Qualityprofiles           *QualityprofilesService
	Rules                     *RulesService
	Server                    *ServerService
	Settings                  *SettingsService
	Sources                   *SourcesService
	System                    *SystemService
	UserGroups                *UserGroupsService
	UserTokens                *UserTokensService
	Users                     *UsersService
	Webhooks                  *WebhooksService
	Webservices               *WebservicesService
	H                         *HService
}

func NewClient(endpoint, username, password string) (*Client, error) {
	c := &Client{username: username, password: password, authType: basicAuth, httpClient: http.DefaultClient}
	if endpoint == "" {
		c.SetBaseURL(defaultBaseURL)
	} else {
		if err := c.SetBaseURL(endpoint); err != nil {
			return nil, err
		}
	}
	c.AlmIntegrations = &AlmIntegrationsService{client: c}
	c.AlmSettings = &AlmSettingsService{client: c}
	c.AnalysisCache = &AnalysisCacheService{client: c}
	c.AnalysisReports = &AnalysisReportsService{client: c}
	c.Authentication = &AuthenticationService{client: c}
	c.Ce = &CeService{client: c}
	c.Components = &ComponentsService{client: c}
	c.Developers = &DevelopersService{client: c}
	c.DismissMessage = &DismissMessageService{client: c}
	c.Duplications = &DuplicationsService{client: c}
	c.Emails = &EmailsService{client: c}
	c.Favorites = &FavoritesService{client: c}
	c.Features = &FeaturesService{client: c}
	c.GithubProvisioning = &GithubProvisioningService{client: c}
	c.Hotspots = &HotspotsService{client: c}
	c.Issues = &IssuesService{client: c}
	c.L10N = &L10NService{client: c}
	c.Languages = &LanguagesService{client: c}
	c.Measures = &MeasuresService{client: c}
	c.Metrics = &MetricsService{client: c}
	c.Monitoring = &MonitoringService{client: c}
	c.Navigation = &NavigationService{client: c}
	c.NewCodePeriods = &NewCodePeriodsService{client: c}
	c.Notifications = &NotificationsService{client: c}
	c.Permissions = &PermissionsService{client: c}
	c.Plugins = &PluginsService{client: c}
	c.ProjectAnalyses = &ProjectAnalysesService{client: c}
	c.ProjectBadges = &ProjectBadgesService{client: c}
	c.ProjectBranches = &ProjectBranchesService{client: c}
	c.ProjectDump = &ProjectDumpService{client: c}
	c.ProjectLinks = &ProjectLinksService{client: c}
	c.ProjectTags = &ProjectTagsService{client: c}
	c.Projects = &ProjectsService{client: c}
	c.Push = &PushService{client: c}
	c.Qualitygates = &QualitygatesService{client: c}
	c.Qualityprofiles = &QualityprofilesService{client: c}
	c.Rules = &RulesService{client: c}
	c.Server = &ServerService{client: c}
	c.Settings = &SettingsService{client: c}
	c.Sources = &SourcesService{client: c}
	c.System = &SystemService{client: c}
	c.UserGroups = &UserGroupsService{client: c}
	c.UserTokens = &UserTokensService{client: c}
	c.Users = &UsersService{client: c}
	c.Webhooks = &WebhooksService{client: c}
	c.Webservices = &WebservicesService{client: c}
	c.H = &HService{client: c}
	return c, nil
}
func NewClientWithToken(endpoint, token string) (*Client, error) {
	c := &Client{token: token, authType: privateToken, httpClient: http.DefaultClient}
	if endpoint == "" {
		c.SetBaseURL(defaultBaseURL)
	} else {
		if err := c.SetBaseURL(endpoint); err != nil {
			return nil, err
		}
	}
	c.AlmIntegrations = &AlmIntegrationsService{client: c}
	c.AlmSettings = &AlmSettingsService{client: c}
	c.AnalysisCache = &AnalysisCacheService{client: c}
	c.AnalysisReports = &AnalysisReportsService{client: c}
	c.Authentication = &AuthenticationService{client: c}
	c.Ce = &CeService{client: c}
	c.Components = &ComponentsService{client: c}
	c.Developers = &DevelopersService{client: c}
	c.DismissMessage = &DismissMessageService{client: c}
	c.Duplications = &DuplicationsService{client: c}
	c.Emails = &EmailsService{client: c}
	c.Favorites = &FavoritesService{client: c}
	c.Features = &FeaturesService{client: c}
	c.GithubProvisioning = &GithubProvisioningService{client: c}
	c.Hotspots = &HotspotsService{client: c}
	c.Issues = &IssuesService{client: c}
	c.L10N = &L10NService{client: c}
	c.Languages = &LanguagesService{client: c}
	c.Measures = &MeasuresService{client: c}
	c.Metrics = &MetricsService{client: c}
	c.Monitoring = &MonitoringService{client: c}
	c.Navigation = &NavigationService{client: c}
	c.NewCodePeriods = &NewCodePeriodsService{client: c}
	c.Notifications = &NotificationsService{client: c}
	c.Permissions = &PermissionsService{client: c}
	c.Plugins = &PluginsService{client: c}
	c.ProjectAnalyses = &ProjectAnalysesService{client: c}
	c.ProjectBadges = &ProjectBadgesService{client: c}
	c.ProjectBranches = &ProjectBranchesService{client: c}
	c.ProjectDump = &ProjectDumpService{client: c}
	c.ProjectLinks = &ProjectLinksService{client: c}
	c.ProjectTags = &ProjectTagsService{client: c}
	c.Projects = &ProjectsService{client: c}
	c.Push = &PushService{client: c}
	c.Qualitygates = &QualitygatesService{client: c}
	c.Qualityprofiles = &QualityprofilesService{client: c}
	c.Rules = &RulesService{client: c}
	c.Server = &ServerService{client: c}
	c.Settings = &SettingsService{client: c}
	c.Sources = &SourcesService{client: c}
	c.System = &SystemService{client: c}
	c.UserGroups = &UserGroupsService{client: c}
	c.UserTokens = &UserTokensService{client: c}
	c.Users = &UsersService{client: c}
	c.Webhooks = &WebhooksService{client: c}
	c.Webservices = &WebservicesService{client: c}
	c.H = &HService{client: c}
	return c, nil
}
