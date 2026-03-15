package sonar

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

	// V2 contains all V2 API services.
	V2 *ServicesV2
}

// ServicesV2 groups all SonarQube V2 API services.
type ServicesV2 struct {
	// Analysis provides methods for the Analysis V2 API.
	Analysis *AnalysisService
	// Authorizations provides methods for the Authorizations V2 API.
	Authorizations *AuthorizationsService
	// CleanCodePolicy provides methods for the Clean Code Policy V2 API.
	CleanCodePolicy *CleanCodePolicyService
	// DopTranslation provides methods for the Dop Translation V2 API.
	DopTranslation *DopTranslationService
	// Marketplace provides methods for the Marketplace V2 API.
	Marketplace *MarketplaceService
	// System provides methods for the System V2 API.
	System *SystemServiceV2
	// UsersManagement provides methods for the Users Management V2 API.
	UsersManagement *UsersManagementService
}

// ClientCreateOptions contains options for creating a new Client.
// Everything is optional and will not throw an error if not provided.
type ClientCreateOptions struct {
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
func NewClient(createOpts *ClientCreateOptions, options ...ClientOptionFunc) (*Client, error) {
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
func applyCreateOptions(client *Client, createOpts *ClientCreateOptions) error {
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

	initServicesV2(client)
}

// initServicesV2 initializes all V2 service instances for the client.
func initServicesV2(client *Client) {
	client.V2 = &ServicesV2{
		Analysis:        &AnalysisService{client: client},
		Authorizations:  &AuthorizationsService{client: client},
		CleanCodePolicy: &CleanCodePolicyService{client: client},
		DopTranslation:  &DopTranslationService{client: client},
		Marketplace:     &MarketplaceService{client: client},
		System:          &SystemServiceV2{client: client},
		UsersManagement: &UsersManagementService{client: client},
	}
}

// SonarAPIRequestParameters contains parameters for making a SonarQube API
// request. This struct is the unified way to provide parameters to
// NewSonarQubeAPIRequest. For most use cases, prefer the version-specific
// helpers NewSonarQubeV1APIRequest and NewSonarQubeV2APIRequest which handle
// query encoding and path prefixing automatically.
//
//nolint:govet // fieldalignment: keeping logical field grouping for readability
type SonarAPIRequestParameters struct {
	// Method is the HTTP method to use for the request (e.g. http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete).
	// Defaults to http.MethodGet if empty.
	Method string
	// Path is the URL path for the API endpoint, relative to the base URL
	// (e.g. "components/search" for V1 or "v2/system/health" for V2).
	// This field is required.
	Path string
	// RawQuery contains pre-encoded URL query parameters to include in the
	// request URL. Use the version-specific helpers which encode query structs
	// automatically based on the appropriate struct tag convention.
	RawQuery url.Values
	// Headers is a map of additional HTTP headers to include in the request.
	// Default is no additional headers beyond the standard Accept, Content-Type,
	// authentication and User-Agent headers.
	Headers map[string]string
	// Body is the request body to include in the API request. It will be
	// JSON-encoded if not nil.
	Body any
}

// NewSonarQubeAPIRequest creates a new API request based on the provided
// SonarAPIRequestParameters. It applies default values for any missing
// parameters, sets authentication and standard headers, and returns an
// http.Request ready to be sent.
//
// For most use cases, prefer the version-specific helpers:
//   - NewSonarQubeV1APIRequest for V1 API endpoints (go-querystring encoding).
//   - NewSonarQubeV2APIRequest for V2 API endpoints (JSON-tag encoding, body support).
func (c *Client) NewSonarQubeAPIRequest(params SonarAPIRequestParameters) (*http.Request, error) {
	method := http.MethodGet
	if params.Method != "" {
		method = params.Method
	}

	if params.Path == "" {
		return nil, errors.New("path is required in SonarAPIRequestParameters")
	}

	requestURL := c.buildRequestURL(params)

	bodyReader, err := marshalBody(params.Body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(context.Background(), method, requestURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setRequestHeaders(req, method, params.Headers)

	return req, nil
}

// NewSonarQubeV1APIRequest creates a V1 API request. The path is resolved
// relative to the client base URL. If opt is non-nil, it is encoded as URL
// query parameters using go-querystring struct tags and appended to the
// request URL.
func (c *Client) NewSonarQubeV1APIRequest(method, path string, opt any) (*http.Request, error) {
	var rawQuery url.Values

	if opt != nil {
		var err error

		rawQuery, err = query.Values(opt)
		if err != nil {
			return nil, fmt.Errorf("failed to encode query values: %w", err)
		}
	}

	//nolint:exhaustruct // Headers and Body intentionally unset for V1 requests
	return c.NewSonarQubeAPIRequest(SonarAPIRequestParameters{
		Method:   method,
		Path:     path,
		RawQuery: rawQuery,
	})
}

// NewSonarQubeV2APIRequest creates a V2 API request. The path is resolved
// relative to the client base URL with the "v2/" prefix automatically
// prepended. If queryOpt is non-nil, it is encoded as URL query parameters
// using JSON struct tags. If body is non-nil, it is JSON-marshaled and used
// as the request body.
func (c *Client) NewSonarQubeV2APIRequest(method, path string, queryOpt any, body any) (*http.Request, error) {
	var rawQuery url.Values

	if queryOpt != nil {
		var err error

		rawQuery, err = jsonStructToQueryValues(queryOpt)
		if err != nil {
			return nil, err
		}
	}

	//nolint:exhaustruct // Headers intentionally unset for V2 requests
	return c.NewSonarQubeAPIRequest(SonarAPIRequestParameters{
		Method:   method,
		Path:     v2BasePath + path,
		RawQuery: rawQuery,
		Body:     body,
	})
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) Do(req *http.Request, dest any) (*http.Response, error) {
	return Do(c.httpClient, req, dest)
}

// =============================================
// UNEXPORTED HELPERS
// =============================================

// buildRequestURL constructs the full URL from the base URL, path and query
// parameters provided in the request parameters.
func (c *Client) buildRequestURL(params SonarAPIRequestParameters) string {
	baseURLCopy := *c.baseURL
	baseURLCopy.Path = c.baseURL.Path + params.Path

	if params.RawQuery != nil {
		baseURLCopy.RawQuery = params.RawQuery.Encode()
	}

	return baseURLCopy.String()
}

// marshalBody JSON-encodes the request body if non-nil. Returns http.NoBody
// when body is nil.
func marshalBody(body any) (io.Reader, error) {
	if body == nil {
		return http.NoBody, nil
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	return bytes.NewReader(data), nil
}

// setRequestHeaders applies Content-Type, Accept, authentication and
// User-Agent headers to the request. Custom headers are applied last so
// callers can override defaults.
func (c *Client) setRequestHeaders(req *http.Request, method string, extraHeaders map[string]string) {
	// Set Content-Type based on HTTP method.
	switch method {
	case http.MethodPatch:
		req.Header.Set("Content-Type", "application/merge-patch+json")
	case http.MethodPost, http.MethodPut:
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")

	// Set authentication headers.
	switch c.authType {
	case basicAuth, oAuthToken:
		req.SetBasicAuth(c.username, c.password)
	case privateToken:
		req.SetBasicAuth(c.token, "")
	}

	req.Header.Set("User-Agent", c.userAgent)

	// Apply any additional custom headers.
	for key, value := range extraHeaders {
		req.Header.Set(key, value)
	}
}
