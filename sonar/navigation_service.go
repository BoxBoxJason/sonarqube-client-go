package sonargo

import "net/http"

// NavigationService handles communication with the navigation related methods
// of the SonarQube API.
// This service provides information required to build navigation UI components.
//
// Since: 5.2.
type NavigationService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// NavigationComponent represents the response from getting component navigation info.
//
//nolint:govet // fieldalignment - structure kept for readability
type NavigationComponent struct {
	// AnalysisDate is the date of the last analysis.
	AnalysisDate string `json:"analysisDate,omitempty"`
	// Breadcrumbs is the list of breadcrumb items for the component.
	Breadcrumbs []NavigationBreadcrumb `json:"breadcrumbs,omitempty"`
	// CanBrowseAllChildProjects indicates if user can browse all child projects.
	CanBrowseAllChildProjects bool `json:"canBrowseAllChildProjects,omitempty"`
	// Configuration contains the component configuration.
	Configuration NavigationConfiguration `json:"configuration,omitzero"`
	// Description is the component description.
	Description string `json:"description,omitempty"`
	// Extensions is the list of available extensions.
	Extensions []NavigationExtension `json:"extensions,omitempty"`
	// ID is the component ID.
	ID string `json:"id,omitempty"`
	// IsFavorite indicates if the component is a favorite.
	IsFavorite bool `json:"isFavorite,omitempty"`
	// Key is the component key.
	Key string `json:"key,omitempty"`
	// Name is the component name.
	Name string `json:"name,omitempty"`
	// QualityGate contains the quality gate information.
	QualityGate NavigationQualityGate `json:"qualityGate,omitzero"`
	// QualityProfiles is the list of quality profiles.
	QualityProfiles []NavigationQualityProfile `json:"qualityProfiles,omitempty"`
	// Version is the component version.
	Version string `json:"version,omitempty"`
}

// NavigationBreadcrumb represents a breadcrumb item in navigation.
type NavigationBreadcrumb struct {
	// Key is the component key.
	Key string `json:"key,omitempty"`
	// Name is the component name.
	Name string `json:"name,omitempty"`
	// Qualifier is the component qualifier.
	Qualifier string `json:"qualifier,omitempty"`
}

// NavigationConfiguration represents navigation configuration options.
//
//nolint:govet // fieldalignment - structure kept for readability
type NavigationConfiguration struct {
	// CanBrowseProject indicates if user can browse the project.
	CanBrowseProject bool `json:"canBrowseProject,omitempty"`
	// Extensions is the list of configuration extensions.
	Extensions []NavigationExtension `json:"extensions,omitempty"`
	// ShowBackgroundTasks indicates if background tasks should be shown.
	ShowBackgroundTasks bool `json:"showBackgroundTasks,omitempty"`
	// ShowHistory indicates if history should be shown.
	ShowHistory bool `json:"showHistory,omitempty"`
	// ShowLinks indicates if links should be shown.
	ShowLinks bool `json:"showLinks,omitempty"`
	// ShowPermissions indicates if permissions should be shown.
	ShowPermissions bool `json:"showPermissions,omitempty"`
	// ShowQualityGates indicates if quality gates should be shown.
	ShowQualityGates bool `json:"showQualityGates,omitempty"`
	// ShowQualityProfiles indicates if quality profiles should be shown.
	ShowQualityProfiles bool `json:"showQualityProfiles,omitempty"`
	// ShowSettings indicates if settings should be shown.
	ShowSettings bool `json:"showSettings,omitempty"`
	// ShowUpdateKey indicates if update key option should be shown.
	ShowUpdateKey bool `json:"showUpdateKey,omitempty"`
}

// NavigationQualityGate represents quality gate navigation info.
//
//nolint:govet // fieldalignment - structure kept for readability
type NavigationQualityGate struct {
	// IsDefault indicates if this is the default quality gate.
	IsDefault bool `json:"isDefault,omitempty"`
	// Key is the quality gate key.
	Key string `json:"key,omitempty"`
	// Name is the quality gate name.
	Name string `json:"name,omitempty"`
}

// NavigationQualityProfile represents quality profile navigation info.
type NavigationQualityProfile struct {
	// Key is the quality profile key.
	Key string `json:"key,omitempty"`
	// Language is the programming language.
	Language string `json:"language,omitempty"`
	// Name is the quality profile name.
	Name string `json:"name,omitempty"`
}

// NavigationExtension represents a navigation extension.
type NavigationExtension struct {
	// Key is the extension key.
	Key string `json:"key,omitempty"`
	// Name is the extension name.
	Name string `json:"name,omitempty"`
}

// NavigationGlobal represents the response from getting global navigation info.
//
//nolint:govet,tagliatelle // fieldalignment - structure kept for readability, API-defined JSON fields
type NavigationGlobal struct {
	// CanAdmin indicates if user has admin permissions.
	CanAdmin bool `json:"canAdmin,omitempty"`
	// DocumentationURL is the URL to the documentation.
	DocumentationURL string `json:"documentationUrl,omitempty"`
	// Edition is the SonarQube edition.
	Edition string `json:"edition,omitempty"`
	// GlobalPages is the list of global pages.
	GlobalPages []NavigationExtension `json:"globalPages,omitempty"`
	// LogoURL is the URL to the logo.
	LogoURL string `json:"logoUrl,omitempty"`
	// LogoWidth is the logo width.
	LogoWidth string `json:"logoWidth,omitempty"`
	// ProductionDatabase indicates if this is a production database.
	ProductionDatabase bool `json:"productionDatabase,omitempty"`
	// Qualifiers is the list of supported qualifiers.
	Qualifiers []string `json:"qualifiers,omitempty"`
	// Settings contains global settings.
	Settings NavigationGlobalSettings `json:"settings,omitzero"`
	// Standalone indicates if running in standalone mode.
	Standalone bool `json:"standalone,omitempty"`
	// Version is the SonarQube version.
	Version string `json:"version,omitempty"`
	// VersionEOL is the end of life date for the installed version.
	VersionEOL string `json:"versionEOL,omitempty"`
}

// NavigationGlobalSettings represents global navigation settings.
//
//nolint:tagliatelle // API-defined JSON field names
type NavigationGlobalSettings struct {
	// EnableGravatar indicates if Gravatar is enabled.
	EnableGravatar string `json:"sonar.lf.enableGravatar,omitempty"`
	// GravatarServerURL is the Gravatar server URL.
	GravatarServerURL string `json:"sonar.lf.gravatarServerUrl,omitempty"`
	// LogoURL is the custom logo URL.
	LogoURL string `json:"sonar.lf.logoUrl,omitempty"`
	// LogoWidthPx is the logo width in pixels.
	LogoWidthPx string `json:"sonar.lf.logoWidthPx,omitempty"`
	// TechnicalDebtRatingGrid is the technical debt rating grid.
	TechnicalDebtRatingGrid string `json:"sonar.technicalDebt.ratingGrid,omitempty"`
	// UpdateCenterActivate indicates if update center is activated.
	UpdateCenterActivate string `json:"sonar.updatecenter.activate,omitempty"`
}

// NavigationMarketplace represents the response from getting marketplace info.
//
//nolint:govet // fieldalignment - structure kept for readability
type NavigationMarketplace struct {
	// Ncloc is the total number of lines of code.
	Ncloc int64 `json:"ncloc,omitempty"`
	// ServerID is the server identifier.
	ServerID string `json:"serverId,omitempty"`
}

// NavigationSettings represents the response from getting settings navigation.
type NavigationSettings struct {
	// Extensions is the list of settings extensions.
	Extensions []NavigationSettingsExtension `json:"extensions,omitempty"`
	// ShowUpdateCenter indicates if update center should be shown.
	ShowUpdateCenter bool `json:"showUpdateCenter,omitempty"`
}

// NavigationSettingsExtension represents a settings extension.
type NavigationSettingsExtension struct {
	// Name is the extension name.
	Name string `json:"name,omitempty"`
	// URL is the extension URL.
	URL string `json:"url,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// NavigationComponentOption represents options for getting component navigation.
type NavigationComponentOption struct {
	// Branch is the branch key (optional).
	Branch string `url:"branch,omitempty"`
	// Component is the component key (optional).
	Component string `url:"component,omitempty"`
	// PullRequest is the pull request ID (optional).
	PullRequest string `url:"pullRequest,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateComponentOpt validates the options for the Component method.
func (s *NavigationService) ValidateComponentOpt(opt *NavigationComponentOption) error {
	if opt == nil {
		return nil
	}

	// No required fields, all parameters are optional
	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Component returns navigation information for a component.
// Requires the 'Browse' permission on the component's project.
// For applications, it also requires 'Browse' permission on its child projects.
//
// API endpoint: GET /api/navigation/component.
// Since: 5.2.
// Internal: true.
func (s *NavigationService) Component(opt *NavigationComponentOption) (*NavigationComponent, *http.Response, error) {
	err := s.ValidateComponentOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "navigation/component", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(NavigationComponent)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Global returns global navigation information for the current user.
//
// API endpoint: GET /api/navigation/global.
// Since: 5.2.
// Internal: true.
func (s *NavigationService) Global() (*NavigationGlobal, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "navigation/global", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(NavigationGlobal)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Marketplace provides data to prefill license request forms.
// Returns the server ID and the total number of lines of code.
//
// API endpoint: GET /api/navigation/marketplace.
// Since: 7.2.
// Internal: true.
func (s *NavigationService) Marketplace() (*NavigationMarketplace, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "navigation/marketplace", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(NavigationMarketplace)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Settings returns configuration information for the settings page.
// This includes plugin-contributed pages and whether to show update center.
//
// API endpoint: GET /api/navigation/settings.
// Since: 5.2.
// Internal: true.
func (s *NavigationService) Settings() (*NavigationSettings, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "navigation/settings", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(NavigationSettings)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
