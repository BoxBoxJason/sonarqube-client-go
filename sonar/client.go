package sonargo

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

// =============================================
// TYPES AND STRUCTS
// =============================================

// Client is the SonarQube API client.
//
//nolint:govet // fieldalignment: keeping logical field grouping for readability
type Client struct {
	baseURL    *url.URL
	username   string
	password   string
	token      string
	authType   authType
	httpClient *http.Client
	userAgent  string

	AlmIntegrations    *AlmIntegrationsService
	AlmSettings        *AlmSettingsService
	AnalysisCache      *AnalysisCacheService
	AnalysisReports    *AnalysisReportsService
	Authentication     *AuthenticationService
	Batch              *BatchService
	Ce                 *CeService
	Components         *ComponentsService
	Developers         *DevelopersService
	DismissMessage     *DismissMessageService
	Duplications       *DuplicationsService
	Emails             *EmailsService
	Favorites          *FavoritesService
	Features           *FeaturesService
	GithubProvisioning *GithubProvisioningService
	Hotspots           *HotspotsService
	Issues             *IssuesService
	L10N               *L10NService
	Languages          *LanguagesService
	Measures           *MeasuresService
	Metrics            *MetricsService
	Monitoring         *MonitoringService
	Navigation         *NavigationService
	NewCodePeriods     *NewCodePeriodsService
	Notifications      *NotificationsService
	Permissions        *PermissionsService
	Plugins            *PluginsService
	ProjectAnalyses    *ProjectAnalysesService
	ProjectBadges      *ProjectBadgesService
	ProjectBranches    *ProjectBranchesService
	ProjectDump        *ProjectDumpService
	ProjectLinks       *ProjectLinksService
	ProjectTags        *ProjectTagsService
	Projects           *ProjectsService
	Push               *PushService
	Qualitygates       *QualitygatesService
	Qualityprofiles    *QualityprofilesService
	Rules              *RulesService
	Server             *ServerService
	Settings           *SettingsService
	Sources            *SourcesService
	System             *SystemService
	UserGroups         *UserGroupsService
	UserTokens         *UserTokensService
	Users              *UsersService
	Webhooks           *WebhooksService
	Webservices        *WebservicesService
}

// ClientCreateOption contains options for creating a new Client.
// Everything is optional and will not throw an error if not provided.
type ClientCreateOption struct {
	// URL is the base API URL for the SonarQube instance.
	URL *string
	// Username is the username for basic authentication.
	Username *string
	// Password is the password for basic authentication.
	Password *string
	// Token is the token for private token authentication.
	Token *string
	// HttpClient is the HTTP client to use for API requests.
	HttpClient *http.Client
	// UserAgent is the User-Agent header to use for API requests.
	UserAgent *string
}

// ClientOptionFunc can be used to customize a new SonarQube API client.
type ClientOptionFunc func(*Client) error

// =============================================
// CLIENT INITIALIZATION
// =============================================

// NewClient creates a new SonarQube API client. createOpts can be used to
// provide initial configuration options. Additional functional options can be
// provided via options.
//
//nolint:exhaustruct // Fields initialized dynamically via options and initServices
func NewClient(createOpts *ClientCreateOption, options ...ClientOptionFunc) (*Client, error) {
	client := &Client{}

	err := applyCreateOptions(client, createOpts)
	if err != nil {
		return nil, err
	}

	err = applyFunctionalOptions(client, options)
	if err != nil {
		return nil, err
	}

	err = setDefaults(client)
	if err != nil {
		return nil, err
	}

	initServices(client)

	return client, nil
}

// applyCreateOptions applies initial configuration options to the client.
func applyCreateOptions(client *Client, createOpts *ClientCreateOption) error {
	if createOpts == nil {
		return nil
	}

	if createOpts.Token != nil {
		client.token = *createOpts.Token
		client.authType = privateToken
	}

	if createOpts.Username != nil && createOpts.Password != nil {
		client.username = *createOpts.Username
		client.password = *createOpts.Password
		client.authType = basicAuth
	}

	if createOpts.URL != nil {
		err := client.SetBaseURL(createOpts.URL)
		if err != nil {
			return err
		}
	}

	assignPtrIfNotNil(&client.httpClient, createOpts.HttpClient)
	assignIfNotNil(&client.userAgent, createOpts.UserAgent)

	return nil
}

// applyFunctionalOptions applies functional options to the client.
func applyFunctionalOptions(client *Client, options []ClientOptionFunc) error {
	for _, option := range options {
		err := option(client)
		if err != nil {
			return err
		}
	}

	return nil
}

// setDefaults sets default values for the client if not provided.
func setDefaults(client *Client) error {
	if client.baseURL == nil {
		defaultURL := defaultBaseURL

		err := client.SetBaseURL(&defaultURL)
		if err != nil {
			return err
		}
	}

	if client.httpClient == nil {
		client.httpClient = http.DefaultClient
	}

	if client.userAgent == "" {
		client.userAgent = defaultUserAgent
	}

	return nil
}

// WithToken is a ClientOptionFunc that sets the token for private token authentication.
func WithToken(token string) ClientOptionFunc {
	return func(c *Client) error {
		c.token = token
		c.authType = privateToken

		return nil
	}
}

// WithBasicAuth is a ClientOptionFunc that sets the username and password for basic authentication.
func WithBasicAuth(username, password string) ClientOptionFunc {
	return func(c *Client) error {
		c.username = username
		c.password = password
		c.authType = basicAuth

		return nil
	}
}

// WithBaseURL is a ClientOptionFunc that sets the base URL for API requests to a custom endpoint.
// urlStr should always be specified with a trailing slash.
func WithBaseURL(urlStr string) ClientOptionFunc {
	return func(c *Client) error {
		return c.SetBaseURL(&urlStr)
	}
}

// WithHTTPClient is a ClientOptionFunc that sets the HTTP client for API requests.
func WithHTTPClient(httpClient *http.Client) ClientOptionFunc {
	return func(c *Client) error {
		c.httpClient = httpClient

		return nil
	}
}

// WithUserAgent is a ClientOptionFunc that sets the User-Agent header for API requests.
func WithUserAgent(userAgent string) ClientOptionFunc {
	return func(c *Client) error {
		c.userAgent = userAgent

		return nil
	}
}

// =============================================
// SETTERS
// =============================================

// SetBaseURL sets the base URL for API requests to a custom endpoint.
// urlStr should always be specified with a trailing slash.
func (c *Client) SetBaseURL(urlStr *string) error {
	if urlStr == nil {
		return errors.New("urlStr cannot be nil")
	}

	// Work on a local copy to avoid mutating the caller-provided string.
	value := *urlStr

	// Make sure the given URL ends with a slash.
	if !strings.HasSuffix(value, "/") {
		value += "/"
	}

	baseURL, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("failed to parse URL: %w", err)
	}

	c.baseURL = baseURL

	return nil
}

// SetHTTPClient sets the HTTP client for API requests.
func (c *Client) SetHTTPClient(httpClient *http.Client) {
	c.httpClient = httpClient
}

// SetBasicAuth sets the username and password for basic authentication.
func (c *Client) SetBasicAuth(username, password string) {
	c.username = username
	c.password = password
	c.authType = basicAuth
}

// SetPrivateToken sets the token for private token authentication.
func (c *Client) SetPrivateToken(token string) {
	c.token = token
	c.authType = privateToken
}

// =============================================
// GETTERS
// =============================================

// BaseURL returns a copy of the base URL.
func (c *Client) BaseURL() *url.URL {
	baseURLCopy := *c.baseURL

	return &baseURLCopy
}

// =============================================
// HELPER FUNCTIONS
// =============================================

// initServices initializes all service instances for the client.
//
//nolint:funlen // necessary to initialize all services
func initServices(client *Client) {
	client.AlmIntegrations = &AlmIntegrationsService{client: client}
	client.AlmSettings = &AlmSettingsService{client: client}
	client.AnalysisCache = &AnalysisCacheService{client: client}
	client.AnalysisReports = &AnalysisReportsService{client: client}
	client.Authentication = &AuthenticationService{client: client}
	client.Batch = &BatchService{client: client}
	client.Ce = &CeService{client: client}
	client.Components = &ComponentsService{client: client}
	client.Developers = &DevelopersService{client: client}
	client.DismissMessage = &DismissMessageService{client: client}
	client.Duplications = &DuplicationsService{client: client}
	client.Emails = &EmailsService{client: client}
	client.Favorites = &FavoritesService{client: client}
	client.Features = &FeaturesService{client: client}
	client.GithubProvisioning = &GithubProvisioningService{client: client}
	client.Hotspots = &HotspotsService{client: client}
	client.Issues = &IssuesService{client: client}
	client.L10N = &L10NService{client: client}
	client.Languages = &LanguagesService{client: client}
	client.Measures = &MeasuresService{client: client}
	client.Metrics = &MetricsService{client: client}
	client.Monitoring = &MonitoringService{client: client}
	client.Navigation = &NavigationService{client: client}
	client.NewCodePeriods = &NewCodePeriodsService{client: client}
	client.Notifications = &NotificationsService{client: client}
	client.Permissions = &PermissionsService{client: client}
	client.Plugins = &PluginsService{client: client}
	client.ProjectAnalyses = &ProjectAnalysesService{client: client}
	client.ProjectBadges = &ProjectBadgesService{client: client}
	client.ProjectBranches = &ProjectBranchesService{client: client}
	client.ProjectDump = &ProjectDumpService{client: client}
	client.ProjectLinks = &ProjectLinksService{client: client}
	client.ProjectTags = &ProjectTagsService{client: client}
	client.Projects = &ProjectsService{client: client}
	client.Push = &PushService{client: client}
	client.Qualitygates = &QualitygatesService{client: client}
	client.Qualityprofiles = &QualityprofilesService{client: client}
	client.Rules = &RulesService{client: client}
	client.Server = &ServerService{client: client}
	client.Settings = &SettingsService{client: client}
	client.Sources = &SourcesService{client: client}
	client.System = &SystemService{client: client}
	client.UserGroups = &UserGroupsService{client: client}
	client.UserTokens = &UserTokensService{client: client}
	client.Users = &UsersService{client: client}
	client.Webhooks = &WebhooksService{client: client}
	client.Webservices = &WebservicesService{client: client}
}

// NewRequest creates an API request. A relative URL path can be provided in
// path, in which case it is resolved relative to the base URL of the Client.
// Relative URL paths should always be specified without a preceding slash.
// If opt is non-nil, it is encoded as URL query parameters using go-querystring
// and appended to the request URL. The request body is not populated.
func (c *Client) NewRequest(method, path string, opt any) (*http.Request, error) {
	baseURLCopy := *c.baseURL
	baseURLCopy.Path = c.baseURL.Path + path

	if opt != nil {
		queryValues, err := query.Values(opt)
		if err != nil {
			return nil, fmt.Errorf("failed to encode query values: %w", err)
		}

		baseURLCopy.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequestWithContext(context.Background(), method, baseURLCopy.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if method == http.MethodPost || method == http.MethodPut {
		// SonarQube uses RawQuery even when method is POST
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")

	switch c.authType {
	case basicAuth, oAuthToken:
		req.SetBasicAuth(c.username, c.password)
	case privateToken:
		req.SetBasicAuth(c.token, "")
	}

	req.Header.Set("User-Agent", c.userAgent)

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) Do(req *http.Request, dest any) (*http.Response, error) {
	return Do(c.httpClient, req, dest)
}
