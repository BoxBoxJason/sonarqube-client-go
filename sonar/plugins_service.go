package sonargo

import "net/http"

// PluginsService handles communication with the plugins related methods
// of the SonarQube API.
// This service provides management of plugins on the server.
//
// Since: 5.2.
type PluginsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// PluginsAvailable represents the response from listing available plugins.
//
//nolint:govet // fieldalignment - structure kept for readability
type PluginsAvailable struct {
	// Plugins is the list of available plugins.
	Plugins []PluginAvailable `json:"plugins,omitempty"`
	// UpdateCenterRefresh is the timestamp when Update Center was last refreshed.
	UpdateCenterRefresh string `json:"updateCenterRefresh,omitempty"`
}

// PluginAvailable represents an available plugin for installation.
//
//nolint:govet // fieldalignment - structure kept for readability
type PluginAvailable struct {
	// Category is the plugin category.
	Category string `json:"category,omitempty"`
	// Description is the plugin description.
	Description string `json:"description,omitempty"`
	// EditionBundled indicates if the plugin is bundled with an edition.
	EditionBundled bool `json:"editionBundled,omitempty"`
	// Key is the plugin key.
	Key string `json:"key,omitempty"`
	// License is the plugin license.
	License string `json:"license,omitempty"`
	// Name is the plugin name.
	Name string `json:"name,omitempty"`
	// OrganizationName is the organization name.
	OrganizationName string `json:"organizationName,omitempty"`
	// OrganizationURL is the organization URL.
	OrganizationURL string `json:"organizationUrl,omitempty"`
	// Release contains release information.
	Release PluginRelease `json:"release,omitzero"`
	// TermsAndConditionsURL is the terms and conditions URL.
	TermsAndConditionsURL string `json:"termsAndConditionsUrl,omitempty"`
	// Update contains update information.
	Update PluginUpdateInfo `json:"update,omitzero"`
}

// PluginRelease represents plugin release information.
type PluginRelease struct {
	// ChangeLogURL is the changelog URL.
	ChangeLogURL string `json:"changeLogUrl,omitempty"`
	// Date is the release date.
	Date string `json:"date,omitempty"`
	// Description is the release description.
	Description string `json:"description,omitempty"`
	// Version is the release version.
	Version string `json:"version,omitempty"`
}

// PluginUpdateInfo represents plugin update information.
//
//nolint:govet // fieldalignment - structure kept for readability
type PluginUpdateInfo struct {
	// Requires is the list of required plugins.
	Requires []PluginRequirement `json:"requires,omitempty"`
	// Status is the update status.
	Status string `json:"status,omitempty"`
}

// PluginRequirement represents a plugin requirement.
type PluginRequirement struct {
	// Description is the requirement description.
	Description string `json:"description,omitempty"`
	// Key is the required plugin key.
	Key string `json:"key,omitempty"`
	// Name is the required plugin name.
	Name string `json:"name,omitempty"`
}

// PluginsInstalled represents the response from listing installed plugins.
type PluginsInstalled struct {
	// Plugins is the list of installed plugins.
	Plugins []PluginInstalled `json:"plugins,omitempty"`
}

// PluginInstalled represents an installed plugin.
//
//nolint:govet // fieldalignment - structure kept for readability
type PluginInstalled struct {
	// Description is the plugin description.
	Description string `json:"description,omitempty"`
	// DocumentationPath is the path to documentation.
	//
	// Deprecated: Since 9.8.
	DocumentationPath string `json:"documentationPath,omitempty"`
	// EditionBundled indicates if the plugin is bundled with an edition.
	EditionBundled bool `json:"editionBundled,omitempty"`
	// Filename is the plugin filename.
	Filename string `json:"filename,omitempty"`
	// Hash is the file hash.
	Hash string `json:"hash,omitempty"`
	// HomepageURL is the plugin homepage URL.
	HomepageURL string `json:"homepageUrl,omitempty"`
	// ImplementationBuild is the implementation build.
	ImplementationBuild string `json:"implementationBuild,omitempty"`
	// IssueTrackerURL is the issue tracker URL.
	IssueTrackerURL string `json:"issueTrackerUrl,omitempty"`
	// Key is the plugin key.
	Key string `json:"key,omitempty"`
	// License is the plugin license.
	License string `json:"license,omitempty"`
	// Name is the plugin name.
	Name string `json:"name,omitempty"`
	// OrganizationName is the organization name.
	OrganizationName string `json:"organizationName,omitempty"`
	// OrganizationURL is the organization URL.
	OrganizationURL string `json:"organizationUrl,omitempty"`
	// RequiredForLanguages is the list of languages requiring this plugin.
	RequiredForLanguages []string `json:"requiredForLanguages,omitempty"`
	// SonarLintSupported indicates if SonarLint is supported.
	SonarLintSupported bool `json:"sonarLintSupported,omitempty"`
	// UpdatedAt is the timestamp when the plugin was updated.
	UpdatedAt int64 `json:"updatedAt,omitempty"`
	// Version is the plugin version.
	Version string `json:"version,omitempty"`
}

// PluginsPending represents the response from listing pending plugin changes.
type PluginsPending struct {
	// Installing is the list of plugins being installed.
	Installing []PluginPending `json:"installing,omitempty"`
	// Removing is the list of plugins being removed.
	Removing []PluginPending `json:"removing,omitempty"`
	// Updating is the list of plugins being updated.
	Updating []PluginPendingUpdate `json:"updating,omitempty"`
}

// PluginPending represents a pending plugin change.
type PluginPending struct {
	// Category is the plugin category.
	Category string `json:"category,omitempty"`
	// Description is the plugin description.
	Description string `json:"description,omitempty"`
	// DocumentationPath is the path to documentation.
	//
	// Deprecated: Since 9.8.
	DocumentationPath string `json:"documentationPath,omitempty"`
	// HomepageURL is the plugin homepage URL.
	HomepageURL string `json:"homepageUrl,omitempty"`
	// ImplementationBuild is the implementation build.
	ImplementationBuild string `json:"implementationBuild,omitempty"`
	// IssueTrackerURL is the issue tracker URL.
	IssueTrackerURL string `json:"issueTrackerUrl,omitempty"`
	// Key is the plugin key.
	Key string `json:"key,omitempty"`
	// License is the plugin license.
	License string `json:"license,omitempty"`
	// Name is the plugin name.
	Name string `json:"name,omitempty"`
	// OrganizationName is the organization name.
	OrganizationName string `json:"organizationName,omitempty"`
	// OrganizationURL is the organization URL.
	OrganizationURL string `json:"organizationUrl,omitempty"`
	// Version is the plugin version.
	Version string `json:"version,omitempty"`
}

// PluginPendingUpdate represents a pending plugin update.
type PluginPendingUpdate struct {
	// Category is the plugin category.
	Category string `json:"category,omitempty"`
	// Description is the plugin description.
	Description string `json:"description,omitempty"`
	// HomepageURL is the plugin homepage URL.
	HomepageURL string `json:"homepageUrl,omitempty"`
	// ImplementationBuild is the implementation build.
	ImplementationBuild string `json:"implementationBuild,omitempty"`
	// IssueTrackerURL is the issue tracker URL.
	IssueTrackerURL string `json:"issueTrackerUrl,omitempty"`
	// Key is the plugin key.
	Key string `json:"key,omitempty"`
	// License is the plugin license.
	License string `json:"license,omitempty"`
	// Name is the plugin name.
	Name string `json:"name,omitempty"`
	// OrganizationName is the organization name.
	OrganizationName string `json:"organizationName,omitempty"`
	// OrganizationURL is the organization URL.
	OrganizationURL string `json:"organizationUrl,omitempty"`
	// Version is the plugin version.
	Version string `json:"version,omitempty"`
}

// PluginsUpdates represents the response from listing plugin updates.
type PluginsUpdates struct {
	// Plugins is the list of plugins with available updates.
	Plugins []PluginWithUpdates `json:"plugins,omitempty"`
}

// PluginWithUpdates represents a plugin with available updates.
//
//nolint:govet // fieldalignment - structure kept for readability
type PluginWithUpdates struct {
	// Category is the plugin category.
	Category string `json:"category,omitempty"`
	// Description is the plugin description.
	Description string `json:"description,omitempty"`
	// EditionBundled indicates if the plugin is bundled with an edition.
	EditionBundled bool `json:"editionBundled,omitempty"`
	// Key is the plugin key.
	Key string `json:"key,omitempty"`
	// License is the plugin license.
	License string `json:"license,omitempty"`
	// Name is the plugin name.
	Name string `json:"name,omitempty"`
	// OrganizationName is the organization name.
	OrganizationName string `json:"organizationName,omitempty"`
	// OrganizationURL is the organization URL.
	OrganizationURL string `json:"organizationUrl,omitempty"`
	// TermsAndConditionsURL is the terms and conditions URL.
	TermsAndConditionsURL string `json:"termsAndConditionsUrl,omitempty"`
	// Updates is the list of available updates.
	Updates []PluginUpdateDetail `json:"updates,omitempty"`
}

// PluginUpdateDetail represents details of a plugin update.
//
//nolint:govet // fieldalignment - structure kept for readability
type PluginUpdateDetail struct {
	// Release contains release information.
	Release PluginRelease `json:"release,omitzero"`
	// Requires is the list of required plugins.
	Requires []PluginRequirement `json:"requires,omitempty"`
	// Status is the update status.
	Status string `json:"status,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// PluginsDownloadOption represents options for downloading a plugin.
type PluginsDownloadOption struct {
	// Plugin is the key identifying the plugin to download (required).
	Plugin string `url:"plugin,omitempty"`
}

// PluginsInstallOption represents options for installing a plugin.
type PluginsInstallOption struct {
	// Key is the key identifying the plugin to install (required).
	Key string `url:"key,omitempty"`
}

// PluginsInstalledOption represents options for listing installed plugins.
//
//nolint:govet // fieldalignment - structure kept for readability
type PluginsInstalledOption struct {
	// Fields is the list of additional fields to return.
	// Possible values: category.
	Fields []string `url:"f,omitempty,comma"`
	// Type filters plugins by type (internal).
	// Possible values: BUNDLED, EXTERNAL.
	Type string `url:"type,omitempty"`
}

// PluginsUninstallOption represents options for uninstalling a plugin.
type PluginsUninstallOption struct {
	// Key is the key identifying the plugin to uninstall (required).
	Key string `url:"key,omitempty"`
}

// PluginsUpdateOption represents options for updating a plugin.
type PluginsUpdateOption struct {
	// Key is the key identifying the plugin to update (required).
	Key string `url:"key,omitempty"`
}

// -----------------------------------------------------------------------------
// Allowed Values
// -----------------------------------------------------------------------------

//nolint:gochecknoglobals // these are constant sets of allowed values
var (
	allowedPluginTypes = map[string]struct{}{
		"BUNDLED":  {},
		"EXTERNAL": {},
	}

	allowedPluginFields = map[string]struct{}{
		"category": {},
	}
)

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateDownloadOpt validates the options for the Download method.
func (s *PluginsService) ValidateDownloadOpt(opt *PluginsDownloadOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Plugin, "Plugin")
}

// ValidateInstallOpt validates the options for the Install method.
func (s *PluginsService) ValidateInstallOpt(opt *PluginsInstallOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Key, "Key")
}

// ValidateInstalledOpt validates the options for the Installed method.
func (s *PluginsService) ValidateInstalledOpt(opt *PluginsInstalledOption) error {
	if opt == nil {
		return nil
	}

	if len(opt.Fields) > 0 {
		err := AreValuesAuthorized(opt.Fields, allowedPluginFields, "Fields")
		if err != nil {
			return err
		}
	}

	if opt.Type != "" {
		err := IsValueAuthorized(opt.Type, allowedPluginTypes, "Type")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateUninstallOpt validates the options for the Uninstall method.
func (s *PluginsService) ValidateUninstallOpt(opt *PluginsUninstallOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Key, "Key")
}

// ValidateUpdateOpt validates the options for the Update method.
func (s *PluginsService) ValidateUpdateOpt(opt *PluginsUpdateOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Key, "Key")
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Available returns the list of all available plugins for installation.
// Plugin information is retrieved from Update Center.
// Requires 'Administer System' permission.
//
// API endpoint: GET /api/plugins/available.
// Since: 5.2.
func (s *PluginsService) Available() (*PluginsAvailable, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "plugins/available", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(PluginsAvailable)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CancelAll cancels any operation pending on any plugin.
// Requires user to be authenticated with Administer System permissions.
//
// API endpoint: POST /api/plugins/cancel_all.
// Since: 5.2.
func (s *PluginsService) CancelAll() (*http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "plugins/cancel_all", nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Download downloads plugin JAR for usage by scanner engine.
//
// API endpoint: GET /api/plugins/download.
// Since: 7.2.
// Internal: true.
func (s *PluginsService) Download(opt *PluginsDownloadOption) (*string, *http.Response, error) {
	err := s.ValidateDownloadOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "plugins/download", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(string)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Install installs the latest version of a plugin specified by its key.
// Plugin information is retrieved from Update Center.
// Fails if used on commercial editions or plugin risk consent has not been accepted.
// Requires user to be authenticated with Administer System permissions.
//
// API endpoint: POST /api/plugins/install.
// Since: 5.2.
func (s *PluginsService) Install(opt *PluginsInstallOption) (*http.Response, error) {
	err := s.ValidateInstallOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "plugins/install", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Installed returns the list of all installed plugins.
// Requires authentication.
//
// API endpoint: GET /api/plugins/installed.
// Since: 5.2.
func (s *PluginsService) Installed(opt *PluginsInstalledOption) (*PluginsInstalled, *http.Response, error) {
	err := s.ValidateInstalledOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "plugins/installed", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(PluginsInstalled)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Pending returns the list of plugins pending installation/removal/update.
// Requires 'Administer System' permission.
//
// API endpoint: GET /api/plugins/pending.
// Since: 5.2.
func (s *PluginsService) Pending() (*PluginsPending, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "plugins/pending", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(PluginsPending)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Uninstall uninstalls the plugin specified by its key.
// Requires user to be authenticated with Administer System permissions.
//
// API endpoint: POST /api/plugins/uninstall.
// Since: 5.2.
func (s *PluginsService) Uninstall(opt *PluginsUninstallOption) (*http.Response, error) {
	err := s.ValidateUninstallOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "plugins/uninstall", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Update updates a plugin to the latest version compatible with the SonarQube instance.
// Plugin information is retrieved from Update Center.
// Requires user to be authenticated with Administer System permissions.
//
// API endpoint: POST /api/plugins/update.
// Since: 5.2.
func (s *PluginsService) Update(opt *PluginsUpdateOption) (*http.Response, error) {
	err := s.ValidateUpdateOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "plugins/update", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Updates lists plugins installed that have at least one newer version available.
// Plugin information is retrieved from Update Center.
// Requires 'Administer System' permission.
//
// API endpoint: GET /api/plugins/updates.
// Since: 5.2.
func (s *PluginsService) Updates() (*PluginsUpdates, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "plugins/updates", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(PluginsUpdates)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
