package sonargo

import "net/http"

const (
	// MaxComponentSearchQueryLength is the maximum length for search queries (truncated by API if exceeded).
	MaxComponentSearchQueryLength = 15
	// MinComponentSearchQueryLength is the minimum length for search queries.
	MinComponentSearchQueryLength = 2
	// MinComponentTreeQueryLength is the minimum length for tree search queries.
	MinComponentTreeQueryLength = 3
	// MaxRecentlyBrowsedItems is the maximum number of recently browsed items.
	MaxRecentlyBrowsedItems = 50
	// MinComponentFilterLength is the minimum length for project filter.
	MinComponentFilterLength = 2
)

// ComponentsService handles communication with the Components related methods
// of the SonarQube API.
// Get information about a component (file, directory, project, ...) and its ancestors or descendants.
type ComponentsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// =============================================================================
// Allowed Values
// =============================================================================

//nolint:gochecknoglobals // constant set of allowed values
var (
	// allowedComponentSearchQualifiers is the set of allowed qualifiers for component search.
	allowedComponentSearchQualifiers = map[string]struct{}{
		"TRK": {},
	}

	// allowedComponentTreeQualifiers is the set of allowed qualifiers for component tree.
	allowedComponentTreeQualifiers = map[string]struct{}{
		"UTS": {},
		"FIL": {},
		"DIR": {},
		"TRK": {},
	}

	// allowedComponentTreeStrategies is the set of allowed strategies for tree navigation.
	allowedComponentTreeStrategies = map[string]struct{}{
		"all":      {},
		"children": {},
		"leaves":   {},
	}

	// allowedComponentTreeSortFields is the set of allowed sort fields for tree.
	allowedComponentTreeSortFields = map[string]struct{}{
		"name":      {},
		"path":      {},
		"qualifier": {},
	}

	// allowedComponentSuggestionsMore is the set of allowed values for the more parameter.
	allowedComponentSuggestionsMore = map[string]struct{}{
		"VW":  {},
		"SVW": {},
		"APP": {},
		"TRK": {},
	}

	// allowedSearchProjectsFields is the set of allowed fields for search projects response.
	allowedSearchProjectsFields = map[string]struct{}{
		"analysisDate":   {},
		"leakPeriodDate": {},
		"_all":           {},
	}

	// allowedSearchProjectsFacets is the set of allowed facets for search projects.
	allowedSearchProjectsFacets = map[string]struct{}{
		"alert_status":                   {},
		"coverage":                       {},
		"duplicated_lines_density":       {},
		"languages":                      {},
		"ncloc":                          {},
		"new_coverage":                   {},
		"new_duplicated_lines_density":   {},
		"new_lines":                      {},
		"new_maintainability_rating":     {},
		"new_reliability_rating":         {},
		"new_security_hotspots_reviewed": {},
		"new_security_rating":            {},
		"new_security_review_rating":     {},
		"new_software_quality_maintainability_rating": {},
		"new_software_quality_reliability_rating":     {},
		"new_software_quality_security_rating":        {},
		"qualifier":                                   {},
		"reliability_rating":                          {},
		"security_hotspots_reviewed":                  {},
		"security_rating":                             {},
		"security_review_rating":                      {},
		"software_quality_maintainability_rating":     {},
		"software_quality_reliability_rating":         {},
		"software_quality_security_rating":            {},
		"sqale_rating":                                {},
		"tags":                                        {},
	}

	// allowedSearchProjectsSortFields is the set of allowed sort fields for search projects.
	allowedSearchProjectsSortFields = map[string]struct{}{
		"alert_status":                   {},
		"analysisDate":                   {},
		"coverage":                       {},
		"creationDate":                   {},
		"duplicated_lines_density":       {},
		"lines":                          {},
		"name":                           {},
		"ncloc":                          {},
		"ncloc_language_distribution":    {},
		"new_coverage":                   {},
		"new_duplicated_lines_density":   {},
		"new_lines":                      {},
		"new_maintainability_rating":     {},
		"new_reliability_rating":         {},
		"new_security_hotspots_reviewed": {},
		"new_security_rating":            {},
		"new_security_review_rating":     {},
		"new_software_quality_maintainability_rating": {},
		"new_software_quality_reliability_rating":     {},
		"new_software_quality_security_rating":        {},
		"reliability_rating":                          {},
		"security_hotspots_reviewed":                  {},
		"security_rating":                             {},
		"security_review_rating":                      {},
		"software_quality_maintainability_rating":     {},
		"software_quality_reliability_rating":         {},
		"software_quality_security_rating":            {},
		"sqale_rating":                                {},
	}
)

// =============================================================================
// Shared Types
// =============================================================================

// ComponentMeasures represents measures for a component (used in App response).
type ComponentMeasures struct {
	// Debt is the technical debt of the component.
	Debt string `json:"debt,omitempty"`
	// DebtRatio is the technical debt ratio of the component.
	DebtRatio string `json:"debtRatio,omitempty"`
	// DuplicationDensity is the duplication density.
	DuplicationDensity string `json:"duplicationDensity,omitempty"`
	// Issues is the number of issues.
	Issues string `json:"issues,omitempty"`
	// Lines is the number of lines.
	Lines string `json:"lines,omitempty"`
	// SqaleRating is the SQALE maintainability rating.
	SqaleRating string `json:"sqaleRating,omitempty"`
}

// ComponentSearchItem represents a component in search results.
type ComponentSearchItem struct {
	// Key is the component key.
	Key string `json:"key,omitempty"`
	// Name is the component name.
	Name string `json:"name,omitempty"`
	// Project is the project key containing this component.
	Project string `json:"project,omitempty"`
	// Qualifier is the component qualifier (TRK, DIR, FIL, etc.).
	Qualifier string `json:"qualifier,omitempty"`
}

// ComponentProject represents a project in search projects results.
//
//nolint:govet // Field order maintained for API response consistency
type ComponentProject struct {
	// AiCodeAssurance is the AI code assurance status.
	AiCodeAssurance string `json:"aiCodeAssurance,omitempty"`
	// ContainsAiCode indicates if the project contains AI-generated code.
	ContainsAiCode bool `json:"containsAiCode,omitempty"`
	// IsAiCodeFixEnabled indicates if AI code fix is enabled.
	IsAiCodeFixEnabled bool `json:"isAiCodeFixEnabled,omitempty"`
	// IsFavorite indicates if the project is a favorite.
	IsFavorite bool `json:"isFavorite,omitempty"`
	// Key is the project key.
	Key string `json:"key,omitempty"`
	// Name is the project name.
	Name string `json:"name,omitempty"`
	// Qualifier is the project qualifier.
	Qualifier string `json:"qualifier,omitempty"`
	// Tags is the list of tags.
	Tags []string `json:"tags,omitempty"`
	// UUID is the project UUID.
	UUID string `json:"uuid,omitempty"`
	// Visibility is the project visibility.
	Visibility string `json:"visibility,omitempty"`
}

// ComponentFacetValue represents a value in a facet.
//
//nolint:govet // Field order maintained for API response consistency
type ComponentFacetValue struct {
	// Count is the number of items matching this value.
	Count int64 `json:"count,omitempty"`
	// Val is the facet value.
	Val string `json:"val,omitempty"`
}

// ComponentFacet represents a facet in search projects results.
type ComponentFacet struct {
	// Property is the facet property name.
	Property string `json:"property,omitempty"`
	// Values is the list of facet values.
	Values []ComponentFacetValue `json:"values,omitempty"`
}

// ComponentAncestor represents an ancestor component in show results.
//
//nolint:govet // Field order maintained for API response consistency
type ComponentAncestor struct {
	// AnalysisDate is the last analysis date.
	AnalysisDate string `json:"analysisDate,omitempty"`
	// Description is the component description.
	Description string `json:"description,omitempty"`
	// Key is the component key.
	Key string `json:"key,omitempty"`
	// Name is the component name.
	Name string `json:"name,omitempty"`
	// Path is the path to the component.
	Path string `json:"path,omitempty"`
	// Qualifier is the component qualifier.
	Qualifier string `json:"qualifier,omitempty"`
	// Tags is the list of tags.
	Tags []string `json:"tags,omitempty"`
	// Version is the component version.
	Version string `json:"version,omitempty"`
	// Visibility is the component visibility.
	Visibility string `json:"visibility,omitempty"`
}

// ComponentDetails represents detailed component information in show results.
type ComponentDetails struct {
	// AnalysisDate is the last analysis date.
	AnalysisDate string `json:"analysisDate,omitempty"`
	// Key is the component key.
	Key string `json:"key,omitempty"`
	// Language is the component language.
	Language string `json:"language,omitempty"`
	// LeakPeriodDate is the leak period date.
	LeakPeriodDate string `json:"leakPeriodDate,omitempty"`
	// Name is the component name.
	Name string `json:"name,omitempty"`
	// Path is the path to the component.
	Path string `json:"path,omitempty"`
	// Qualifier is the component qualifier.
	Qualifier string `json:"qualifier,omitempty"`
	// Version is the component version.
	Version string `json:"version,omitempty"`
}

// ComponentSuggestionItem represents an item in suggestions results.
//
//nolint:govet // Field order maintained for API response consistency
type ComponentSuggestionItem struct {
	// IsFavorite indicates if the component is a favorite.
	IsFavorite bool `json:"isFavorite,omitempty"`
	// IsRecentlyBrowsed indicates if the component was recently browsed.
	IsRecentlyBrowsed bool `json:"isRecentlyBrowsed,omitempty"`
	// Key is the component key.
	Key string `json:"key,omitempty"`
	// Match is the display name with emphasis on matching characters.
	Match string `json:"match,omitempty"`
	// Name is the component name.
	Name string `json:"name,omitempty"`
	// Project is the project key.
	Project string `json:"project,omitempty"`
}

// ComponentSuggestionGroup represents a group of suggestions by qualifier.
//
//nolint:govet // Field order maintained for API response consistency
type ComponentSuggestionGroup struct {
	// Items is the list of suggestions in this group.
	Items []ComponentSuggestionItem `json:"items,omitempty"`
	// More is the count of additional results available.
	More int64 `json:"more,omitempty"`
	// Q is the qualifier for this group.
	Q string `json:"q,omitempty"`
}

// ComponentTreeBase represents the base component in tree results.
//
//nolint:govet // Field order maintained for API response consistency
type ComponentTreeBase struct {
	// Description is the component description.
	Description string `json:"description,omitempty"`
	// Key is the component key.
	Key string `json:"key,omitempty"`
	// Qualifier is the component qualifier.
	Qualifier string `json:"qualifier,omitempty"`
	// Tags is the list of tags.
	Tags []string `json:"tags,omitempty"`
	// Visibility is the component visibility.
	Visibility string `json:"visibility,omitempty"`
}

// ComponentTreeItem represents a component in tree results.
type ComponentTreeItem struct {
	// Key is the component key.
	Key string `json:"key,omitempty"`
	// Language is the component language.
	Language string `json:"language,omitempty"`
	// Name is the component name.
	Name string `json:"name,omitempty"`
	// Path is the path to the component.
	Path string `json:"path,omitempty"`
	// Qualifier is the component qualifier.
	Qualifier string `json:"qualifier,omitempty"`
}

// =============================================================================
// Response Types
// =============================================================================

// ComponentsApp represents the response from the App method.
// Contains coverage data required for rendering the component viewer.
//
//nolint:govet // Field order maintained for API response consistency
type ComponentsApp struct {
	// CanCreateManualIssue indicates if manual issues can be created.
	CanCreateManualIssue bool `json:"canCreateManualIssue,omitempty"`
	// CanMarkAsFavorite indicates if the component can be marked as favorite.
	CanMarkAsFavorite bool `json:"canMarkAsFavorite,omitempty"`
	// Fav indicates if the component is a favorite.
	Fav bool `json:"fav,omitempty"`
	// Key is the component key.
	Key string `json:"key,omitempty"`
	// LongName is the long name of the component.
	LongName string `json:"longName,omitempty"`
	// Measures contains component measures.
	Measures ComponentMeasures `json:"measures,omitzero"`
	// Name is the component name.
	Name string `json:"name,omitempty"`
	// Project is the project key.
	Project string `json:"project,omitempty"`
	// ProjectName is the project name.
	ProjectName string `json:"projectName,omitempty"`
	// Q is the component qualifier.
	Q string `json:"q,omitempty"`
	// UUID is the component UUID.
	UUID string `json:"uuid,omitempty"`
}

// ComponentsSearch represents the response from the Search method.
// Contains a list of components and pagination information.
type ComponentsSearch struct {
	// Components is the list of found components.
	Components []ComponentSearchItem `json:"components,omitempty"`
	// Paging contains pagination information.
	Paging Paging `json:"paging"`
}

// ComponentsSearchProjects represents the response from the SearchProjects method.
// Contains a list of projects, facets, and pagination information.
type ComponentsSearchProjects struct {
	// Components is the list of found projects.
	Components []ComponentProject `json:"components,omitempty"`
	// Facets is the list of computed facets.
	Facets []ComponentFacet `json:"facets,omitempty"`
	// Paging contains pagination information.
	Paging Paging `json:"paging"`
}

// ComponentsShow represents the response from the Show method.
// Contains a component and its ancestors.
//
//nolint:govet // Field order maintained for API response consistency
type ComponentsShow struct {
	// Ancestors is the list of ancestor components, ordered from parent to root.
	Ancestors []ComponentAncestor `json:"ancestors,omitempty"`
	// Component contains detailed component information.
	Component ComponentDetails `json:"component"`
}

// ComponentsSuggestions represents the response from the Suggestions method.
// Contains search results grouped by qualifier.
type ComponentsSuggestions struct {
	// Projects is the list of project suggestions (deprecated, use Results).
	Projects []any `json:"projects,omitempty"`
	// Results is the list of suggestion groups by qualifier.
	Results []ComponentSuggestionGroup `json:"results,omitempty"`
}

// ComponentsTree represents the response from the Tree method.
// Contains the base component, its descendants, and pagination information.
type ComponentsTree struct {
	// BaseComponent is the base component from which the tree navigation starts.
	BaseComponent ComponentTreeBase `json:"baseComponent"`
	// Components is the list of descendant components.
	Components []ComponentTreeItem `json:"components,omitempty"`
	// Paging contains pagination information.
	Paging Paging `json:"paging"`
}

// =============================================================================
// Option Types
// =============================================================================

// ComponentsAppOption contains parameters for the App method.
type ComponentsAppOption struct {
	// Branch is the branch key. Not available in the community edition.
	// Either branch or pullRequest can be provided, not both.
	Branch string `url:"branch,omitempty"`
	// Component is the component key.
	// This field is required.
	Component string `url:"component"`
	// PullRequest is the pull request id. Not available in the community edition.
	// Either branch or pullRequest can be provided, not both.
	PullRequest string `url:"pullRequest,omitempty"`
}

// ComponentsSearchOption contains parameters for the Search method.
//
//nolint:govet // Field order maintained for API parameter consistency
type ComponentsSearchOption struct {
	// PaginationArgs contains the pagination parameters.
	PaginationArgs `url:",inline"`

	// Query limits search to component names that contain the supplied string
	// or component keys that are exactly the same as the supplied string.
	// The value length must be between 2 and 15 (inclusive) characters.
	// Longer values will be truncated.
	Query string `url:"q,omitempty"`
	// Qualifiers is the list of component qualifiers to filter by.
	// Allowed values: TRK (Projects).
	// This field is required.
	Qualifiers []string `url:"qualifiers,comma"`
}

// ComponentsSearchProjectsOption contains parameters for the SearchProjects method.
//
//nolint:govet // Field order maintained for API parameter consistency
type ComponentsSearchProjectsOption struct {
	// PaginationArgs contains the pagination parameters.
	PaginationArgs `url:",inline"`

	// Ascending indicates ascending sort order.
	// Default: true.
	Ascending bool `url:"asc,omitempty"`
	// Fields is the list of fields to be returned in response.
	// Allowed values: analysisDate, leakPeriodDate, _all.
	Fields []string `url:"f,omitempty,comma"`
	// Facets is the list of facets to be computed.
	// No facet is computed by default.
	Facets []string `url:"facets,omitempty,comma"`
	// Filter is the filter expression for projects on name, key, measure value,
	// quality gate, language, tag or whether a project is a favorite.
	// The filter must be URL-encoded.
	Filter string `url:"filter,omitempty"`
	// Sort is the sort field.
	// Allowed values: alert_status, analysisDate, coverage, creationDate, etc.
	Sort string `url:"s,omitempty"`
}

// ComponentsShowOption contains parameters for the Show method.
type ComponentsShowOption struct {
	// Branch is the branch key. Not available in the community edition.
	Branch string `url:"branch,omitempty"`
	// Component is the component key.
	// This field is required.
	Component string `url:"component"`
	// PullRequest is the pull request id. Not available in the community edition.
	PullRequest string `url:"pullRequest,omitempty"`
}

// ComponentsSuggestionsOption contains parameters for the Suggestions method.
//
//nolint:govet // Field order maintained for API parameter consistency
type ComponentsSuggestionsOption struct {
	// More is the category for which to display the next 20 results.
	// Allowed values: VW, SVW, APP, TRK.
	More string `url:"more,omitempty"`
	// RecentlyBrowsed is the list of component keys that have recently been browsed.
	// Only the first 50 items will be used. Order is not taken into account.
	RecentlyBrowsed []string `url:"recentlyBrowsed,omitempty,comma"`
	// Search is the search query. Can contain several search tokens separated by spaces.
	// Minimum length: 2 characters.
	Search string `url:"s,omitempty"`
}

// ComponentsTreeOption contains parameters for the Tree method.
//
//nolint:govet // Field order maintained for API parameter consistency
type ComponentsTreeOption struct {
	// PaginationArgs contains the pagination parameters.
	PaginationArgs `url:",inline"`

	// Ascending indicates ascending sort order.
	// Default: true.
	Ascending bool `url:"asc,omitempty"`
	// Branch is the branch key. Not available in the community edition.
	Branch string `url:"branch,omitempty"`
	// Component is the base component key. The search is based on this component.
	// This field is required.
	Component string `url:"component"`
	// PullRequest is the pull request id. Not available in the community edition.
	PullRequest string `url:"pullRequest,omitempty"`
	// Query limits search to component names that contain the supplied string
	// or component keys that are exactly the same as the supplied string.
	// Minimum length: 3 characters.
	Query string `url:"q,omitempty"`
	// Qualifiers is the list of component qualifiers to filter by.
	// Allowed values: UTS, FIL, DIR, TRK.
	Qualifiers []string `url:"qualifiers,omitempty,comma"`
	// Sort is the list of sort fields.
	// Allowed values: name, path, qualifier.
	// Default: name.
	Sort []string `url:"s,omitempty,comma"`
	// Strategy is the strategy to search for base component descendants.
	// Allowed values: all, children, leaves.
	// Default: all.
	Strategy string `url:"strategy,omitempty"`
}

// =============================================================================
// Validation Functions
// =============================================================================

// ValidateAppOpt validates the options for the App method.
func (s *ComponentsService) ValidateAppOpt(opt *ComponentsAppOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Component, "Component")
	if err != nil {
		return err
	}

	if opt.Branch != "" && opt.PullRequest != "" {
		return NewValidationError("Branch", "branch and pullRequest are mutually exclusive", ErrInvalidValue)
	}

	return nil
}

// ValidateSearchOpt validates the options for the Search method.
func (s *ComponentsService) ValidateSearchOpt(opt *ComponentsSearchOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := opt.Validate()
	if err != nil {
		return err
	}

	if len(opt.Qualifiers) == 0 {
		return NewValidationError("Qualifiers", "is required", ErrMissingRequired)
	}

	err = AreValuesAuthorized(opt.Qualifiers, allowedComponentSearchQualifiers, "Qualifiers")
	if err != nil {
		return err
	}

	err = ValidateMinLength(opt.Query, MinComponentSearchQueryLength, "Query")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSearchProjectsOpt validates the options for the SearchProjects method.
func (s *ComponentsService) ValidateSearchProjectsOpt(opt *ComponentsSearchProjectsOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := opt.Validate()
	if err != nil {
		return err
	}

	err = AreValuesAuthorized(opt.Fields, allowedSearchProjectsFields, "Fields")
	if err != nil {
		return err
	}

	err = AreValuesAuthorized(opt.Facets, allowedSearchProjectsFacets, "Facets")
	if err != nil {
		return err
	}

	err = ValidateMinLength(opt.Filter, MinComponentFilterLength, "Filter")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Sort, allowedSearchProjectsSortFields, "Sort")
	if err != nil {
		return err
	}

	return nil
}

// ValidateShowOpt validates the options for the Show method.
func (s *ComponentsService) ValidateShowOpt(opt *ComponentsShowOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Component, "Component")
	if err != nil {
		return err
	}

	if opt.Branch != "" && opt.PullRequest != "" {
		return NewValidationError("Branch", "branch and pullRequest are mutually exclusive", ErrInvalidValue)
	}

	return nil
}

// ValidateSuggestionsOpt validates the options for the Suggestions method.
func (s *ComponentsService) ValidateSuggestionsOpt(opt *ComponentsSuggestionsOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := IsValueAuthorized(opt.More, allowedComponentSuggestionsMore, "More")
	if err != nil {
		return err
	}

	if len(opt.RecentlyBrowsed) > MaxRecentlyBrowsedItems {
		return NewValidationError("RecentlyBrowsed", "cannot exceed 50 items", ErrOutOfRange)
	}

	err = ValidateMinLength(opt.Search, MinComponentSearchQueryLength, "Search")
	if err != nil {
		return err
	}

	return nil
}

// ValidateTreeOpt validates the options for the Tree method.
func (s *ComponentsService) ValidateTreeOpt(opt *ComponentsTreeOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := opt.Validate()
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Component, "Component")
	if err != nil {
		return err
	}

	err = ValidateMinLength(opt.Query, MinComponentTreeQueryLength, "Query")
	if err != nil {
		return err
	}

	err = AreValuesAuthorized(opt.Qualifiers, allowedComponentTreeQualifiers, "Qualifiers")
	if err != nil {
		return err
	}

	err = AreValuesAuthorized(opt.Sort, allowedComponentTreeSortFields, "Sort")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Strategy, allowedComponentTreeStrategies, "Strategy")
	if err != nil {
		return err
	}

	if opt.Branch != "" && opt.PullRequest != "" {
		return NewValidationError("Branch", "branch and pullRequest are mutually exclusive", ErrInvalidValue)
	}

	return nil
}

// =============================================================================
// Service Methods
// =============================================================================

// App returns coverage data required for rendering the component viewer.
// Either branch or pull request can be provided, not both.
// Requires the following permission: 'Browse'.
//
// This is an internal API and may change without notice.
//
// Since: 4.4.
func (s *ComponentsService) App(opt *ComponentsAppOption) (*ComponentsApp, *http.Response, error) {
	err := s.ValidateAppOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "components/app", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ComponentsApp)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Search searches for components.
//
// Since: 6.3.
func (s *ComponentsService) Search(opt *ComponentsSearchOption) (*ComponentsSearch, *http.Response, error) {
	err := s.ValidateSearchOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "components/search", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ComponentsSearch)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SearchProjects searches for projects.
//
// This is an internal API and may change without notice.
//
// Since: 6.2.
func (s *ComponentsService) SearchProjects(opt *ComponentsSearchProjectsOption) (*ComponentsSearchProjects, *http.Response, error) {
	err := s.ValidateSearchProjectsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "components/search_projects", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ComponentsSearchProjects)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Show returns a component (file, directory, project, portfolio…) and its ancestors.
// The ancestors are ordered from the parent to the root project.
// Requires the following permission: 'Browse' on the project of the specified component.
//
// Since: 5.4.
func (s *ComponentsService) Show(opt *ComponentsShowOption) (*ComponentsShow, *http.Response, error) {
	err := s.ValidateShowOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "components/show", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ComponentsShow)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Suggestions returns component suggestions for the top-right search engine.
// The result contains component search results grouped by their qualifiers.
// Each result contains the component key, name, and optionally a display name
// with emphasis on matching characters.
//
// This is an internal API and may change without notice.
//
// Since: 4.2.
func (s *ComponentsService) Suggestions(opt *ComponentsSuggestionsOption) (*ComponentsSuggestions, *http.Response, error) {
	err := s.ValidateSuggestionsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "components/suggestions", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ComponentsSuggestions)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Tree navigates through components based on the chosen strategy.
// Requires the following permission: 'Browse' on the specified project.
// When limiting search with the q parameter, directories are not returned.
//
// Since: 5.4.
func (s *ComponentsService) Tree(opt *ComponentsTreeOption) (*ComponentsTree, *http.Response, error) {
	err := s.ValidateTreeOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "components/tree", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ComponentsTree)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
