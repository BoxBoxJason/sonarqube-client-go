package sonar

import "net/http"

// UserTokensService handles communication with the user tokens related methods
// of the SonarQube API.
// This service lists, creates, and deletes a user's access tokens.
type UserTokensService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

//nolint:gochecknoglobals // these are constant sets of allowed values
var (
	// allowedTokenTypes is the set of supported token types.
	allowedTokenTypes = map[string]struct{}{
		"USER_TOKEN":             {},
		"GLOBAL_ANALYSIS_TOKEN":  {},
		"PROJECT_ANALYSIS_TOKEN": {},
	}
)

// -----------------------------------------------------------------------------
// Shared Types
// -----------------------------------------------------------------------------

// UserTokenProject represents a project associated with a token.
type UserTokenProject struct {
	// Key is the project key.
	Key string `json:"key,omitempty"`
	// Name is the project name.
	Name string `json:"name,omitempty"`
}

// UserToken represents a user access token.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type UserToken struct {
	// CreatedAt is the creation date of the token.
	CreatedAt string `json:"createdAt,omitempty"`
	// ExpirationDate is the expiration date of the token.
	ExpirationDate string `json:"expirationDate,omitempty"`
	// IsExpired indicates whether the token has expired.
	IsExpired bool `json:"isExpired,omitempty"`
	// Name is the name of the token.
	Name string `json:"name,omitempty"`
	// Project is the project associated with a PROJECT_ANALYSIS_TOKEN.
	Project UserTokenProject `json:"project,omitzero"`
	// Type is the type of the token (USER_TOKEN, GLOBAL_ANALYSIS_TOKEN, PROJECT_ANALYSIS_TOKEN).
	Type string `json:"type,omitempty"`
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// UserTokensGenerate represents the response from generating a token.
type UserTokensGenerate struct {
	// CreatedAt is the creation date of the token.
	CreatedAt string `json:"createdAt,omitempty"`
	// ExpirationDate is the expiration date of the token.
	ExpirationDate string `json:"expirationDate,omitempty"`
	// Login is the user login.
	Login string `json:"login,omitempty"`
	// Name is the name of the token.
	Name string `json:"name,omitempty"`
	// Token is the generated token value (only returned once).
	Token string `json:"token,omitempty"`
	// Type is the type of the token.
	Type string `json:"type,omitempty"`
}

// UserTokensSearch represents the response from searching tokens.
type UserTokensSearch struct {
	// Login is the user login.
	Login string `json:"login,omitempty"`
	// UserTokens is the list of tokens.
	UserTokens []UserToken `json:"userTokens,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// UserTokensGenerateOption contains parameters for the Generate method.
type UserTokensGenerateOption struct {
	// ExpirationDate is the expiration date of the token in ISO 8601 format (YYYY-MM-DD).
	// If not set, defaults to no expiration.
	ExpirationDate string `url:"expirationDate,omitempty"`
	// Login is the user login. If not set, the token is generated for the authenticated user.
	Login string `url:"login,omitempty"`
	// Name is the token name.
	// This field is required. Maximum length is 100 characters.
	Name string `url:"name"`
	// ProjectKey is the key of the only project that can be analyzed by the PROJECT_ANALYSIS_TOKEN.
	// Required when Type is set to PROJECT_ANALYSIS_TOKEN.
	ProjectKey string `url:"projectKey,omitempty"`
	// Type is the token type.
	// Allowed values: USER_TOKEN, GLOBAL_ANALYSIS_TOKEN, PROJECT_ANALYSIS_TOKEN.
	// Default: USER_TOKEN.
	// If set to PROJECT_ANALYSIS_TOKEN, ProjectKey must be provided.
	Type string `url:"type,omitempty"`
}

// UserTokensRevokeOption contains parameters for the Revoke method.
type UserTokensRevokeOption struct {
	// Login is the user login.
	Login string `url:"login,omitempty"`
	// Name is the token name.
	// This field is required.
	Name string `url:"name"`
}

// UserTokensSearchOption contains parameters for the Search method.
type UserTokensSearchOption struct {
	// Login is the user login.
	Login string `url:"login,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateGenerateOpt validates the options for the Generate method.
func (s *UserTokensService) ValidateGenerateOpt(opt *UserTokensGenerateOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Name, MaxTokenNameLength, "Name")
	if err != nil {
		return err
	}

	if opt.Type != "" {
		err := IsValueAuthorized(opt.Type, allowedTokenTypes, "Type")
		if err != nil {
			return err
		}

		// If Type is PROJECT_ANALYSIS_TOKEN, ProjectKey is required
		if opt.Type == "PROJECT_ANALYSIS_TOKEN" && opt.ProjectKey == "" {
			return NewValidationError("ProjectKey", "is required when Type is PROJECT_ANALYSIS_TOKEN", ErrMissingRequired)
		}
	}

	return nil
}

// ValidateRevokeOpt validates the options for the Revoke method.
func (s *UserTokensService) ValidateRevokeOpt(opt *UserTokensRevokeOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSearchOpt validates the options for the Search method.
func (s *UserTokensService) ValidateSearchOpt(opt *UserTokensSearchOption) error {
	// Options are optional; nothing to validate.
	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Generate generates a user access token.
// Please keep your tokens secret. They enable authentication and project analysis.
// It requires administration permissions to specify a 'login' and generate a token for another user.
// Otherwise, a token is generated for the current user.
//
// API endpoint: POST /api/user_tokens/generate.
// Since: 5.3.
func (s *UserTokensService) Generate(opt *UserTokensGenerateOption) (*UserTokensGenerate, *http.Response, error) {
	err := s.ValidateGenerateOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "user_tokens/generate", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(UserTokensGenerate)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Revoke revokes a user access token.
// It requires administration permissions to specify a 'login' and revoke a token for another user.
// Otherwise, the token for the current user is revoked.
//
// API endpoint: POST /api/user_tokens/revoke.
// Since: 5.3.
func (s *UserTokensService) Revoke(opt *UserTokensRevokeOption) (*http.Response, error) {
	err := s.ValidateRevokeOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "user_tokens/revoke", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Search lists the access tokens of a user.
// The login must exist and be active.
// Field 'lastConnectionDate' is only updated every hour, so it may not be accurate.
// It requires administration permissions to specify a 'login' and list the tokens of another user.
// Otherwise, tokens for the current user are listed.
// Authentication is required for this API endpoint.
//
// API endpoint: GET /api/user_tokens/search.
// Since: 5.3.
func (s *UserTokensService) Search(opt *UserTokensSearchOption) (*UserTokensSearch, *http.Response, error) {
	err := s.ValidateSearchOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "user_tokens/search", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(UserTokensSearch)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
