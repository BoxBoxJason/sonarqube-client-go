package sonargo

import "net/http"

const (
	// MaxLoginLength is the maximum length for a user login.
	MaxLoginLength = 255
	// MinLoginLength is the minimum length for a user login.
	MinLoginLength = 2
	// MinPasswordLength is the minimum required length for user passwords.
	MinPasswordLength = 12
	// MaxEmailLength is the maximum length for a user email.
	MaxEmailLength = 100
	// MaxNameLength is the maximum length for a user name.
	MaxNameLength = 200
)

// UsersService handles communication with the user management related methods
// of the SonarQube API.
// This service manages users, passwords, groups, and identity providers.
type UsersService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

//nolint:gochecknoglobals // these are constant sets of allowed values
var (
	// allowedHomepageTypes is the set of supported homepage types.
	allowedHomepageTypes = map[string]struct{}{
		"PROJECT":     {},
		"PROJECTS":    {},
		"ISSUES":      {},
		"PORTFOLIOS":  {},
		"PORTFOLIO":   {},
		"APPLICATION": {},
	}

	// allowedNoticeTypes is the set of supported notice types.
	allowedNoticeTypes = map[string]struct{}{
		"educationPrinciples":                   {},
		"sonarlintAd":                           {},
		"showDesignAndArchitectureBanner":       {},
		"showNewModesBanner":                    {},
		"showSandboxedIssuesIntro":              {},
		"issueCleanCodeGuide":                   {},
		"issueNewIssueStatusAndTransitionGuide": {},
		"showDesignAndArchitectureOptInBanner":  {},
		"overviewZeroNewIssuesSimplification":   {},
		"showDesignAndArchitectureTour":         {},
		"showEnableSca":                         {},
	}
)

// -----------------------------------------------------------------------------
// Shared Types
// -----------------------------------------------------------------------------

// User represents a user in the system.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type User struct {
	// Active indicates whether the user is active.
	Active bool `json:"active,omitempty"`
	// Email is the user's email address.
	Email string `json:"email,omitempty"`
	// Local indicates whether the user is authenticated locally.
	Local bool `json:"local,omitempty"`
	// Login is the user's login.
	Login string `json:"login,omitempty"`
	// Name is the user's display name.
	Name string `json:"name,omitempty"`
	// ScmAccounts is the list of SCM accounts associated with the user.
	ScmAccounts []string `json:"scmAccounts,omitempty"`
}

// SearchedUser represents a user returned in search results.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type SearchedUser struct {
	// Active indicates whether the user is active.
	Active bool `json:"active,omitempty"`
	// Avatar is the user's avatar URL.
	Avatar string `json:"avatar,omitempty"`
	// Email is the user's email address.
	Email string `json:"email,omitempty"`
	// ExternalIdentity is the user's external identity.
	ExternalIdentity string `json:"externalIdentity,omitempty"`
	// ExternalProvider is the user's external identity provider.
	ExternalProvider string `json:"externalProvider,omitempty"`
	// Groups is the list of groups the user belongs to.
	Groups []string `json:"groups,omitempty"`
	// LastConnectionDate is the user's last connection date.
	LastConnectionDate string `json:"lastConnectionDate,omitempty"`
	// Local indicates whether the user is authenticated locally.
	Local bool `json:"local,omitempty"`
	// Login is the user's login.
	Login string `json:"login,omitempty"`
	// Managed indicates whether the user is managed externally.
	Managed bool `json:"managed,omitempty"`
	// Name is the user's display name.
	Name string `json:"name,omitempty"`
	// ScmAccounts is the list of SCM accounts associated with the user.
	ScmAccounts []string `json:"scmAccounts,omitempty"`
	// SonarLintLastConnectionDate is the user's last SonarLint connection date.
	SonarLintLastConnectionDate string `json:"sonarLintLastConnectionDate,omitempty"`
	// TokensCount is the number of tokens owned by the user.
	TokensCount int64 `json:"tokensCount,omitempty"`
}

// DeactivatedUser represents a deactivated user.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type DeactivatedUser struct {
	// Active indicates whether the user is active (always false for deactivated users).
	Active bool `json:"active,omitempty"`
	// Groups is the list of groups (empty for deactivated users).
	Groups []any `json:"groups,omitempty"`
	// Local indicates whether the user was authenticated locally.
	Local bool `json:"local,omitempty"`
	// Login is the user's login.
	Login string `json:"login,omitempty"`
	// Name is the user's display name.
	Name string `json:"name,omitempty"`
	// ScmAccounts is the list of SCM accounts (empty for deactivated users).
	ScmAccounts []any `json:"scmAccounts,omitempty"`
}

// CurrentUser represents the currently authenticated user.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type CurrentUser struct {
	// Avatar is the user's avatar URL.
	Avatar string `json:"avatar,omitempty"`
	// DismissedNotices contains the notices dismissed by the user.
	DismissedNotices DismissedNotices `json:"dismissedNotices,omitzero"`
	// Email is the user's email address.
	Email string `json:"email,omitempty"`
	// ExternalIdentity is the user's external identity.
	ExternalIdentity string `json:"externalIdentity,omitempty"`
	// ExternalProvider is the user's external identity provider.
	ExternalProvider string `json:"externalProvider,omitempty"`
	// Groups is the list of groups the user belongs to.
	Groups []string `json:"groups,omitempty"`
	// Homepage is the user's homepage configuration.
	Homepage Homepage `json:"homepage,omitzero"`
	// ID is the user's unique identifier.
	ID string `json:"id,omitempty"`
	// IsLoggedIn indicates whether the user is currently logged in.
	IsLoggedIn bool `json:"isLoggedIn,omitempty"`
	// Local indicates whether the user is authenticated locally.
	Local bool `json:"local,omitempty"`
	// Login is the user's login.
	Login string `json:"login,omitempty"`
	// Name is the user's display name.
	Name string `json:"name,omitempty"`
	// Permissions contains the user's global permissions.
	Permissions UserPermissions `json:"permissions,omitzero"`
	// ScmAccounts is the list of SCM accounts associated with the user.
	ScmAccounts []string `json:"scmAccounts,omitempty"`
	// UsingSonarLintConnectedMode indicates whether the user is using SonarLint connected mode.
	UsingSonarLintConnectedMode bool `json:"usingSonarLintConnectedMode,omitempty"`
}

// DismissedNotices represents the notices dismissed by a user.
type DismissedNotices struct {
	// EducationPrinciples indicates whether the education principles notice was dismissed.
	EducationPrinciples bool `json:"educationPrinciples,omitempty"`
	// SonarlintAd indicates whether the SonarLint ad notice was dismissed.
	SonarlintAd bool `json:"sonarlintAd,omitempty"`
}

// Homepage represents a user's homepage configuration.
type Homepage struct {
	// Component is the project key for project homepages.
	Component string `json:"component,omitempty"`
	// Type is the homepage type.
	Type string `json:"type,omitempty"`
}

// UserPermissions represents a user's global permissions.
type UserPermissions struct {
	// Global is the list of global permissions.
	Global []string `json:"global,omitempty"`
}

// UserGroup represents a group that a user belongs to.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type UserGroup struct {
	// Default indicates whether this is a default group.
	Default bool `json:"default,omitempty"`
	// Description is the group description.
	Description string `json:"description,omitempty"`
	// ID is the group's unique identifier.
	ID string `json:"id,omitempty"`
	// Name is the group name.
	Name string `json:"name,omitempty"`
	// Selected indicates whether the user is a member of this group.
	Selected bool `json:"selected,omitempty"`
}

// IdentityProvider represents an external identity provider.
type IdentityProvider struct {
	// BackgroundColor is the background color for the provider icon.
	BackgroundColor string `json:"backgroundColor,omitempty"`
	// HelpMessage is a help message for the identity provider.
	HelpMessage string `json:"helpMessage,omitempty"`
	// IconPath is the path to the provider icon.
	IconPath string `json:"iconPath,omitempty"`
	// Key is the unique key of the identity provider.
	Key string `json:"key,omitempty"`
	// Name is the display name of the identity provider.
	Name string `json:"name,omitempty"`
}

// UsersPaging represents pagination information for user queries.
type UsersPaging struct {
	// PageIndex is the current page index (1-based).
	PageIndex int64 `json:"pageIndex,omitempty"`
	// PageSize is the number of items per page.
	PageSize int64 `json:"pageSize,omitempty"`
	// Total is the total number of items.
	Total int64 `json:"total,omitempty"`
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// UsersCreate represents the response from creating a user.
type UsersCreate struct {
	// User is the created user.
	User User `json:"user,omitzero"`
}

// UsersCurrent represents the response from getting the current user.
type UsersCurrent = CurrentUser

// UsersDeactivate represents the response from deactivating a user.
type UsersDeactivate struct {
	// User is the deactivated user.
	User DeactivatedUser `json:"user,omitzero"`
}

// UsersGroups represents the response from listing a user's groups.
type UsersGroups struct {
	// Groups is the list of groups.
	Groups []UserGroup `json:"groups,omitempty"`
	// Paging contains pagination information.
	Paging UsersPaging `json:"paging,omitzero"`
}

// UsersIdentityProviders represents the response from listing identity providers.
type UsersIdentityProviders struct {
	// IdentityProviders is the list of identity providers.
	IdentityProviders []IdentityProvider `json:"identityProviders,omitempty"`
}

// UsersSearch represents the response from searching users.
type UsersSearch struct {
	// Users is the list of users.
	Users []SearchedUser `json:"users,omitempty"`
	// Paging contains pagination information.
	Paging UsersPaging `json:"paging,omitzero"`
}

// UsersUpdate represents the response from updating a user.
type UsersUpdate struct {
	// User is the updated user.
	User User `json:"user,omitzero"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// UsersAnonymizeOption contains parameters for the Anonymize method.
type UsersAnonymizeOption struct {
	// Login is the user login.
	// This field is required.
	Login string `url:"login"`
}

// UsersChangePasswordOption contains parameters for the ChangePassword method.
type UsersChangePasswordOption struct {
	// Login is the user login.
	// This field is required.
	Login string `url:"login"`
	// Password is the new password.
	// This field is required. Must be at least 12 characters.
	Password string `url:"password"`
	// PreviousPassword is the previous password. Required when changing one's own password.
	PreviousPassword string `url:"previousPassword,omitempty"`
}

// UsersCreateOption contains parameters for the Create method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type UsersCreateOption struct {
	// Email is the user email address.
	Email string `url:"email,omitempty"`
	// Local specifies if the user should be authenticated from SonarQube server.
	// When true, the user is authenticated locally. When false, user is external.
	// Password should not be set when Local is false.
	Local bool `url:"local,omitempty"`
	// Login is the user login.
	// This field is required. Must be between 2 and 255 characters.
	Login string `url:"login"`
	// Name is the user's display name.
	// This field is required.
	Name string `url:"name"`
	// Password is the user password.
	// Only required when creating a local user.
	Password string `url:"password,omitempty"`
	// ScmAccounts is the list of SCM accounts.
	ScmAccounts []string `url:"scmAccount,omitempty"`
}

// UsersDeactivateOption contains parameters for the Deactivate method.
type UsersDeactivateOption struct {
	// Login is the user login.
	// This field is required.
	Login string `url:"login"`
	// Anonymize specifies whether to anonymize the user in addition to deactivating.
	Anonymize bool `url:"anonymize,omitempty"`
}

// UsersDismissNoticeOption contains parameters for the DismissNotice method.
type UsersDismissNoticeOption struct {
	// Notice is the notice key to dismiss.
	// This field is required.
	// Allowed values: educationPrinciples, sonarlintAd, showDesignAndArchitectureBanner,
	// showNewModesBanner, showSandboxedIssuesIntro, issueCleanCodeGuide,
	// issueNewIssueStatusAndTransitionGuide, showDesignAndArchitectureOptInBanner,
	// overviewZeroNewIssuesSimplification, showDesignAndArchitectureTour, showEnableSca.
	Notice string `url:"notice"`
}

// UsersGroupsOption contains parameters for the Groups method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type UsersGroupsOption struct {
	PaginationArgs

	// Login is the user login.
	// This field is required.
	Login string `url:"login"`
	// Q is a limit search to group names that contain the supplied string.
	Q string `url:"q,omitempty"`
	// Selected filters the selection status.
	// Allowed values: all, selected, deselected.
	Selected string `url:"selected,omitempty"`
}

// UsersSearchOption contains parameters for the Search method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type UsersSearchOption struct {
	PaginationArgs

	// Deactivated returns deactivated users instead of active users when true.
	Deactivated bool `url:"deactivated,omitempty"`
	// ExternalIdentity finds a user by their external identity (case-sensitive).
	// Only available with Administer System permission.
	ExternalIdentity string `url:"externalIdentity,omitempty"`
	// LastConnectedAfter filters users who connected at or after this date.
	// Format: ISO 8601 datetime (YYYY-MM-DDThh:mm:ss±hhmm).
	LastConnectedAfter string `url:"lastConnectedAfter,omitempty"`
	// LastConnectedBefore filters users who never connected or connected at or before this date.
	// Format: ISO 8601 datetime (YYYY-MM-DDThh:mm:ss±hhmm).
	LastConnectedBefore string `url:"lastConnectedBefore,omitempty"`
	// Managed returns managed or non-managed users.
	// Only available for managed instances.
	Managed bool `url:"managed,omitempty"`
	// Q filters on login, name and email (partial match, case insensitive).
	Q string `url:"q,omitempty"`
	// SlLastConnectedAfter filters users who connected via SonarLint at or after this date.
	// Format: ISO 8601 datetime (YYYY-MM-DDThh:mm:ss±hhmm).
	SlLastConnectedAfter string `url:"slLastConnectedAfter,omitempty"`
	// SlLastConnectedBefore filters users who never connected or connected via SonarLint at or before this date.
	// Format: ISO 8601 datetime (YYYY-MM-DDThh:mm:ss±hhmm).
	SlLastConnectedBefore string `url:"slLastConnectedBefore,omitempty"`
}

// UsersSetHomepageOption contains parameters for the SetHomepage method.
type UsersSetHomepageOption struct {
	// Branch is the branch key. Only used when Type is PROJECT.
	Branch string `url:"branch,omitempty"`
	// Component is the project key. Only used when Type is PROJECT.
	Component string `url:"component,omitempty"`
	// Type is the type of the requested homepage.
	// This field is required.
	// Allowed values: PROJECT, PROJECTS, ISSUES, PORTFOLIOS, PORTFOLIO, APPLICATION.
	Type string `url:"type"`
}

// UsersUpdateOption contains parameters for the Update method.
type UsersUpdateOption struct {
	// Email is the user's new email address.
	Email string `url:"email,omitempty"`
	// Login is the user login.
	// This field is required.
	Login string `url:"login"`
	// Name is the user's new display name.
	Name string `url:"name,omitempty"`
	// ScmAccounts is the list of SCM accounts.
	ScmAccounts []string `url:"scmAccount,omitempty"`
}

// UsersUpdateIdentityProviderOption contains parameters for the UpdateIdentityProvider method.
type UsersUpdateIdentityProviderOption struct {
	// Login is the user login.
	// This field is required.
	Login string `url:"login"`
	// NewExternalIdentity is the new external identity.
	// If not provided, the previous identity will be used.
	NewExternalIdentity string `url:"newExternalIdentity,omitempty"`
	// NewExternalProvider is the new external provider key.
	// This field is required.
	NewExternalProvider string `url:"newExternalProvider"`
}

// UsersUpdateLoginOption contains parameters for the UpdateLogin method.
type UsersUpdateLoginOption struct {
	// Login is the current login (case-sensitive).
	// This field is required.
	Login string `url:"login"`
	// NewLogin is the new login.
	// This field is required. Must be between 2 and 255 characters.
	NewLogin string `url:"newLogin"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateAnonymizeOpt validates the options for the Anonymize method.
func (s *UsersService) ValidateAnonymizeOpt(opt *UsersAnonymizeOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Login, "Login")
	if err != nil {
		return err
	}

	return nil
}

// ValidateChangePasswordOpt validates the options for the ChangePassword method.
func (s *UsersService) ValidateChangePasswordOpt(opt *UsersChangePasswordOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Login, "Login")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Password, "Password")
	if err != nil {
		return err
	}

	err = ValidateMinLength(opt.Password, MinPasswordLength, "Password")
	if err != nil {
		return err
	}

	return nil
}

// validateLogin validates the login field with required, min and max length checks.
func validateLogin(login string) error {
	err := ValidateRequired(login, "Login")
	if err != nil {
		return err
	}

	err = ValidateMinLength(login, MinLoginLength, "Login")
	if err != nil {
		return err
	}

	return ValidateMaxLength(login, MaxLoginLength, "Login")
}

// validateUserPassword validates password for user creation.
func validateUserPassword(password string, isLocal bool) error {
	// Password is required for local users
	if isLocal && password == "" {
		return NewValidationError("Password", "is required for local users", ErrMissingRequired)
	}

	if password != "" {
		return ValidateMinLength(password, MinPasswordLength, "Password")
	}

	return nil
}

// ValidateCreateOpt validates the options for the Create method.
func (s *UsersService) ValidateCreateOpt(opt *UsersCreateOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := validateLogin(opt.Login)
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Name, MaxNameLength, "Name")
	if err != nil {
		return err
	}

	if opt.Email != "" {
		err = ValidateMaxLength(opt.Email, MaxEmailLength, "Email")
		if err != nil {
			return err
		}
	}

	return validateUserPassword(opt.Password, opt.Local)
}

// ValidateDeactivateOpt validates the options for the Deactivate method.
func (s *UsersService) ValidateDeactivateOpt(opt *UsersDeactivateOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Login, "Login")
	if err != nil {
		return err
	}

	return nil
}

// ValidateDismissNoticeOpt validates the options for the DismissNotice method.
func (s *UsersService) ValidateDismissNoticeOpt(opt *UsersDismissNoticeOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Notice, "Notice")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Notice, allowedNoticeTypes, "Notice")
	if err != nil {
		return err
	}

	return nil
}

// ValidateGroupsOpt validates the options for the Groups method.
func (s *UsersService) ValidateGroupsOpt(opt *UsersGroupsOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Login, "Login")
	if err != nil {
		return err
	}

	err = opt.Validate()
	if err != nil {
		return err
	}

	if opt.Selected != "" {
		err = IsValueAuthorized(opt.Selected, allowedSelectedFilters, "Selected")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateSearchOpt validates the options for the Search method.
func (s *UsersService) ValidateSearchOpt(opt *UsersSearchOption) error {
	if opt == nil {
		// Search with no options is valid
		return nil
	}

	err := opt.Validate()
	if err != nil {
		return err
	}

	return nil
}

// ValidateSetHomepageOpt validates the options for the SetHomepage method.
func (s *UsersService) ValidateSetHomepageOpt(opt *UsersSetHomepageOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Type, "Type")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Type, allowedHomepageTypes, "Type")
	if err != nil {
		return err
	}

	return nil
}

// ValidateUpdateOpt validates the options for the Update method.
func (s *UsersService) ValidateUpdateOpt(opt *UsersUpdateOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Login, "Login")
	if err != nil {
		return err
	}

	if opt.Email != "" {
		err = ValidateMaxLength(opt.Email, MaxEmailLength, "Email")
		if err != nil {
			return err
		}
	}

	if opt.Name != "" {
		err = ValidateMaxLength(opt.Name, MaxNameLength, "Name")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateUpdateIdentityProviderOpt validates the options for the UpdateIdentityProvider method.
func (s *UsersService) ValidateUpdateIdentityProviderOpt(opt *UsersUpdateIdentityProviderOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Login, "Login")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.NewExternalProvider, "NewExternalProvider")
	if err != nil {
		return err
	}

	return nil
}

// ValidateUpdateLoginOpt validates the options for the UpdateLogin method.
func (s *UsersService) ValidateUpdateLoginOpt(opt *UsersUpdateLoginOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Login, "Login")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.NewLogin, "NewLogin")
	if err != nil {
		return err
	}

	err = ValidateMinLength(opt.NewLogin, MinLoginLength, "NewLogin")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.NewLogin, MaxLoginLength, "NewLogin")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Anonymize anonymizes a deactivated user.
// Requires Administer System permission.
//
// Deprecated: Since SonarQube 10.4.
// API endpoint: POST /api/users/anonymize.
// Since: 9.7.
func (s *UsersService) Anonymize(opt *UsersAnonymizeOption) (*http.Response, error) {
	err := s.ValidateAnonymizeOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "users/anonymize", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// ChangePassword updates a user's password.
// Authenticated users can change their own password, provided that the account is not
// linked to an external authentication system. Administer System permission is required
// to change another user's password.
//
// API endpoint: POST /api/users/change_password.
// Since: 5.2.
func (s *UsersService) ChangePassword(opt *UsersChangePasswordOption) (*http.Response, error) {
	err := s.ValidateChangePasswordOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "users/change_password", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Create creates a user.
// If a deactivated user account exists with the given login, it will be reactivated.
// Requires Administer System permission.
//
// Deprecated: Since SonarQube 10.4.
// API endpoint: POST /api/users/create.
// Since: 3.7.
func (s *UsersService) Create(opt *UsersCreateOption) (*UsersCreate, *http.Response, error) {
	err := s.ValidateCreateOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "users/create", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(UsersCreate)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Current gets the details of the current authenticated user.
//
// API endpoint: GET /api/users/current.
// Since: 5.2.
func (s *UsersService) Current() (*UsersCurrent, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "users/current", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(UsersCurrent)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Deactivate deactivates a user.
// Requires Administer System permission.
//
// Deprecated: Since SonarQube 10.4.
// API endpoint: POST /api/users/deactivate.
// Since: 3.7.
func (s *UsersService) Deactivate(opt *UsersDeactivateOption) (*UsersDeactivate, *http.Response, error) {
	err := s.ValidateDeactivateOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "users/deactivate", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(UsersDeactivate)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DismissNotice dismisses a notice for the current user.
// Silently ignores if the notice is already dismissed.
//
// API endpoint: POST /api/users/dismiss_notice.
// Since: 9.6.
func (s *UsersService) DismissNotice(opt *UsersDismissNoticeOption) (*http.Response, error) {
	err := s.ValidateDismissNoticeOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "users/dismiss_notice", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Groups lists the groups a user belongs to.
// Requires Administer System permission.
//
// Deprecated: Since SonarQube 10.4.
// API endpoint: GET /api/users/groups.
// Since: 5.2.
func (s *UsersService) Groups(opt *UsersGroupsOption) (*UsersGroups, *http.Response, error) {
	err := s.ValidateGroupsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "users/groups", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(UsersGroups)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// IdentityProviders lists the external identity providers.
//
// API endpoint: GET /api/users/identity_providers.
// Since: 5.5.
func (s *UsersService) IdentityProviders() (*UsersIdentityProviders, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "users/identity_providers", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(UsersIdentityProviders)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Search searches for users.
// By default, only active users are returned.
// The following fields are only returned when the user has Administer System permission
// or for the logged-in user: email, externalIdentity, externalProvider, groups,
// lastConnectionDate, sonarLintLastConnectionDate, tokensCount.
// Field 'lastConnectionDate' is only updated every hour, so it may not be accurate.
//
// Deprecated: Since SonarQube 10.4.
// API endpoint: GET /api/users/search.
// Since: 3.6.
func (s *UsersService) Search(opt *UsersSearchOption) (*UsersSearch, *http.Response, error) {
	err := s.ValidateSearchOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "users/search", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(UsersSearch)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SetHomepage sets the homepage of the current user.
// Requires authentication.
//
// API endpoint: POST /api/users/set_homepage.
// Since: 7.0.
func (s *UsersService) SetHomepage(opt *UsersSetHomepageOption) (*http.Response, error) {
	err := s.ValidateSetHomepageOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "users/set_homepage", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Update updates a user.
// Requires Administer System permission.
//
// Deprecated: Since SonarQube 10.4.
// API endpoint: POST /api/users/update.
// Since: 3.7.
func (s *UsersService) Update(opt *UsersUpdateOption) (*UsersUpdate, *http.Response, error) {
	err := s.ValidateUpdateOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "users/update", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(UsersUpdate)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateIdentityProvider updates the identity provider information for a user.
// It's only possible to migrate to an installed identity provider.
// Once updated, the user will only be able to authenticate on the new identity provider.
// It is not possible to migrate an external user to a local one.
// Requires Administer System permission.
//
// Deprecated: Since SonarQube 10.4.
// API endpoint: POST /api/users/update_identity_provider.
// Since: 8.7.
func (s *UsersService) UpdateIdentityProvider(opt *UsersUpdateIdentityProviderOption) (*http.Response, error) {
	err := s.ValidateUpdateIdentityProviderOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "users/update_identity_provider", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UpdateLogin updates a user login.
// A login can be updated many times.
// Requires Administer System permission.
//
// Deprecated: Since SonarQube 10.4.
// API endpoint: POST /api/users/update_login.
// Since: 7.6.
func (s *UsersService) UpdateLogin(opt *UsersUpdateLoginOption) (*http.Response, error) {
	err := s.ValidateUpdateLoginOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "users/update_login", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
