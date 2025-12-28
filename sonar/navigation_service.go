// Get information required to build navigation UI components
package sonargo

import "net/http"

type NavigationService struct {
	client *Client
}

type NavigationComponentObject struct {
	AnalysisDate              string                           `json:"analysisDate,omitempty"`
	Breadcrumbs               []NavigationComponentObject_sub1 `json:"breadcrumbs,omitempty"`
	CanBrowseAllChildProjects bool                             `json:"canBrowseAllChildProjects,omitempty"`
	Configuration             NavigationComponentObject_sub3   `json:"configuration,omitempty"`
	Description               string                           `json:"description,omitempty"`
	Extensions                []NavigationComponentObject_sub2 `json:"extensions,omitempty"`
	ID                        string                           `json:"id,omitempty"`
	IsFavorite                bool                             `json:"isFavorite,omitempty"`
	Key                       string                           `json:"key,omitempty"`
	Name                      string                           `json:"name,omitempty"`
	QualityGate               NavigationComponentObject_sub4   `json:"qualityGate,omitempty"`
	QualityProfiles           []NavigationComponentObject_sub5 `json:"qualityProfiles,omitempty"`
	Version                   string                           `json:"version,omitempty"`
}

type NavigationComponentObject_sub3 struct {
	CanBrowseProject    bool                             `json:"canBrowseProject,omitempty"`
	Extensions          []NavigationComponentObject_sub2 `json:"extensions,omitempty"`
	ShowBackgroundTasks bool                             `json:"showBackgroundTasks,omitempty"`
	ShowHistory         bool                             `json:"showHistory,omitempty"`
	ShowLinks           bool                             `json:"showLinks,omitempty"`
	ShowPermissions     bool                             `json:"showPermissions,omitempty"`
	ShowQualityGates    bool                             `json:"showQualityGates,omitempty"`
	ShowQualityProfiles bool                             `json:"showQualityProfiles,omitempty"`
	ShowSettings        bool                             `json:"showSettings,omitempty"`
	ShowUpdateKey       bool                             `json:"showUpdateKey,omitempty"`
}

type NavigationComponentObject_sub4 struct {
	IsDefault bool   `json:"isDefault,omitempty"`
	Key       string `json:"key,omitempty"`
	Name      string `json:"name,omitempty"`
}

type NavigationComponentObject_sub5 struct {
	Key      string `json:"key,omitempty"`
	Language string `json:"language,omitempty"`
	Name     string `json:"name,omitempty"`
}

type NavigationComponentObject_sub1 struct {
	Key       string `json:"key,omitempty"`
	Name      string `json:"name,omitempty"`
	Qualifier string `json:"qualifier,omitempty"`
}

type NavigationComponentObject_sub2 struct {
	Key  string `json:"key,omitempty"`
	Name string `json:"name,omitempty"`
}

type NavigationGlobalObject struct {
	CanAdmin           bool                          `json:"canAdmin,omitempty"`
	DocumentationURL   string                        `json:"documentationUrl,omitempty"`
	Edition            string                        `json:"edition,omitempty"`
	GlobalPages        []NavigationGlobalObject_sub1 `json:"globalPages,omitempty"`
	LogoURL            string                        `json:"logoUrl,omitempty"`
	LogoWidth          string                        `json:"logoWidth,omitempty"`
	ProductionDatabase bool                          `json:"productionDatabase,omitempty"`
	Qualifiers         []string                      `json:"qualifiers,omitempty"`
	Settings           NavigationGlobalObject_sub2   `json:"settings,omitempty"`
	Standalone         bool                          `json:"standalone,omitempty"`
	Version            string                        `json:"version,omitempty"`
	VersionEOL         string                        `json:"versionEOL,omitempty"`
}

type NavigationGlobalObject_sub1 struct {
	Key  string `json:"key,omitempty"`
	Name string `json:"name,omitempty"`
}

type NavigationGlobalObject_sub2 struct {
	Sonar_lf_enableGravatar        string `json:"sonar.lf.enableGravatar,omitempty"`
	Sonar_lf_gravatarServerURL     string `json:"sonar.lf.gravatarServerUrl,omitempty"`
	Sonar_lf_logoURL               string `json:"sonar.lf.logoUrl,omitempty"`
	Sonar_lf_logoWidthPx           string `json:"sonar.lf.logoWidthPx,omitempty"`
	Sonar_technicalDebt_ratingGrid string `json:"sonar.technicalDebt.ratingGrid,omitempty"`
	Sonar_updatecenter_activate    string `json:"sonar.updatecenter.activate,omitempty"`
}

type NavigationMarketplaceObject struct {
	Ncloc    int64  `json:"ncloc,omitempty"`
	ServerID string `json:"serverId,omitempty"`
}

type NavigationSettingsObject struct {
	Extensions       []NavigationSettingsObject_sub1 `json:"extensions,omitempty"`
	ShowUpdateCenter bool                            `json:"showUpdateCenter,omitempty"`
}

type NavigationSettingsObject_sub1 struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

type NavigationComponentOption struct {
	Branch      string `url:"branch,omitempty"`      // Description:"Branch key",ExampleValue:"feature/my_branch"
	Component   string `url:"component,omitempty"`   // Description:"A component key.",ExampleValue:"my_project"
	PullRequest string `url:"pullRequest,omitempty"` // Description:"Pull request id",ExampleValue:"5461"
}

// Component Get information concerning component navigation for the current user. Requires the 'Browse' permission on the component's project. <br>For applications, it also requires 'Browse' permission on its child projects.
func (s *NavigationService) Component(opt *NavigationComponentOption) (v *NavigationComponentObject, resp *http.Response, err error) {
	err = s.ValidateComponentOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "navigation/component", opt)
	if err != nil {
		return
	}
	v = new(NavigationComponentObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

// Global Get information concerning global navigation for the current user.
func (s *NavigationService) Global() (v *NavigationGlobalObject, resp *http.Response, err error) {
	req, err := s.client.NewRequest("GET", "navigation/global", nil)
	if err != nil {
		return
	}
	v = new(NavigationGlobalObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

// Marketplace Provide data to prefill license request forms: the server ID and the total number of lines of code.
func (s *NavigationService) Marketplace() (v *NavigationMarketplaceObject, resp *http.Response, err error) {
	req, err := s.client.NewRequest("GET", "navigation/marketplace", nil)
	if err != nil {
		return
	}
	v = new(NavigationMarketplaceObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

// Settings Get configuration information for the settings page:<ul>  <li>List plugin-contributed pages</li>  <li>Show update center (or not)</li></ul>
func (s *NavigationService) Settings() (v *NavigationSettingsObject, resp *http.Response, err error) {
	req, err := s.client.NewRequest("GET", "navigation/settings", nil)
	if err != nil {
		return
	}
	v = new(NavigationSettingsObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}
