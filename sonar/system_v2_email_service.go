package sonar

import (
	"context"
	"net/http"
)

const (
	// EmailAuthMethodBasic is the basic (username/password) email authentication method.
	EmailAuthMethodBasic = "BASIC"
	// EmailAuthMethodOAuth is the OAuth email authentication method.
	EmailAuthMethodOAuth = "OAUTH"
)

//nolint:gochecknoglobals // allowed-value sets
var (
	// allowedEmailSecurityProtocols is the set of allowed email security protocols.
	allowedEmailSecurityProtocols = map[string]struct{}{"NONE": {}, "SSLTLS": {}, "STARTTLS": {}}
	// allowedEmailAuthMethods is the set of allowed email authentication methods.
	allowedEmailAuthMethods = map[string]struct{}{EmailAuthMethodBasic: {}, EmailAuthMethodOAuth: {}}
)

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// SystemEmailConfiguration represents an email configuration.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type SystemEmailConfiguration struct {
	// Id is the email configuration identifier.
	Id string `json:"id,omitempty"`
	// Host is the SMTP host.
	Host string `json:"host,omitempty"`
	// Port is the SMTP port.
	Port string `json:"port,omitempty"`
	// SecurityProtocol is the security protocol. Allowed values: NONE, SSLTLS, STARTTLS.
	SecurityProtocol string `json:"securityProtocol,omitempty"`
	// FromAddress is the sender address.
	FromAddress string `json:"fromAddress,omitempty"`
	// FromName is the sender name.
	FromName string `json:"fromName,omitempty"`
	// SubjectPrefix is the prefix added to email subjects.
	SubjectPrefix string `json:"subjectPrefix,omitempty"`
	// AuthMethod is the authentication method. Allowed values: BASIC, OAUTH.
	AuthMethod string `json:"authMethod,omitempty"`
	// Username is the SMTP username.
	Username *string `json:"username,omitempty"`
	// IsBasicPasswordSet indicates whether a basic password is set.
	IsBasicPasswordSet bool `json:"isBasicPasswordSet,omitempty"`
	// OauthAuthenticationHost is the OAuth authentication host.
	OauthAuthenticationHost *string `json:"oauthAuthenticationHost,omitempty"`
	// IsOauthClientIdSet indicates whether an OAuth client id is set.
	IsOauthClientIdSet bool `json:"isOauthClientIdSet,omitempty"`
	// IsOauthClientSecretSet indicates whether an OAuth client secret is set.
	IsOauthClientSecretSet bool `json:"isOauthClientSecretSet,omitempty"`
	// OauthTenant is the OAuth tenant.
	OauthTenant *string `json:"oauthTenant,omitempty"`
	// OauthScope is the scope of the OAuth token request. Since 2026.5.
	OauthScope *string `json:"oauthScope,omitempty"`
}

// SystemEmailConfigurationsSearch represents the response from searching email configurations.
type SystemEmailConfigurationsSearch struct {
	// EmailConfigurations is the list of email configurations.
	EmailConfigurations []SystemEmailConfiguration `json:"emailConfigurations,omitempty"`
	// Page contains pagination information.
	Page PageResponseV2 `json:"page,omitzero"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// SystemEmailConfigurationCreateOptions contains parameters for the CreateEmailConfiguration method.
type SystemEmailConfigurationCreateOptions struct {
	// Host is the SMTP host. This field is required.
	Host string `json:"host"`
	// Port is the SMTP port. This field is required.
	Port string `json:"port"`
	// SecurityProtocol is the security protocol. Allowed values: NONE, SSLTLS, STARTTLS. This field is required.
	SecurityProtocol string `json:"securityProtocol"`
	// FromAddress is the sender address. This field is required.
	FromAddress string `json:"fromAddress"`
	// FromName is the sender name. This field is required.
	FromName string `json:"fromName"`
	// SubjectPrefix is the prefix added to email subjects. This field is required.
	SubjectPrefix string `json:"subjectPrefix"`
	// AuthMethod is the authentication method. Allowed values: BASIC, OAUTH. This field is required.
	AuthMethod string `json:"authMethod"`
	// Username is the SMTP username. This field is required.
	Username string `json:"username"`
	// BasicPassword is the password, for BASIC authentication.
	BasicPassword string `json:"basicPassword,omitempty"`
	// OauthAuthenticationHost is the OAuth authentication host, for OAUTH authentication.
	OauthAuthenticationHost string `json:"oauthAuthenticationHost,omitempty"`
	// OauthClientId is the OAuth client id, for OAUTH authentication.
	OauthClientId string `json:"oauthClientId,omitempty"`
	// OauthClientSecret is the OAuth client secret, for OAUTH authentication.
	OauthClientSecret string `json:"oauthClientSecret,omitempty"`
	// OauthTenant is the OAuth tenant, for OAUTH authentication.
	OauthTenant string `json:"oauthTenant,omitempty"`
	// OauthScope is the scope of the OAuth token request.
	// Defaults to 'https://outlook.office365.com/.default'. Since 2026.5.
	OauthScope string `json:"oauthScope,omitempty"`
}

// SystemEmailConfigurationUpdateOptions contains parameters for the UpdateEmailConfiguration method.
// All fields are optional (PATCH merge semantics).
type SystemEmailConfigurationUpdateOptions struct {
	// Host is the SMTP host.
	Host string `json:"host,omitempty"`
	// Port is the SMTP port.
	Port string `json:"port,omitempty"`
	// SecurityProtocol is the security protocol. Allowed values: NONE, SSLTLS, STARTTLS.
	SecurityProtocol string `json:"securityProtocol,omitempty"`
	// FromAddress is the sender address.
	FromAddress string `json:"fromAddress,omitempty"`
	// FromName is the sender name.
	FromName string `json:"fromName,omitempty"`
	// SubjectPrefix is the prefix added to email subjects.
	SubjectPrefix string `json:"subjectPrefix,omitempty"`
	// AuthMethod is the authentication method. Allowed values: BASIC, OAUTH.
	AuthMethod string `json:"authMethod,omitempty"`
	// Username is the SMTP username.
	Username string `json:"username,omitempty"`
	// BasicPassword is the password, for BASIC authentication.
	BasicPassword string `json:"basicPassword,omitempty"`
	// OauthAuthenticationHost is the OAuth authentication host.
	OauthAuthenticationHost string `json:"oauthAuthenticationHost,omitempty"`
	// OauthClientId is the OAuth client id.
	OauthClientId string `json:"oauthClientId,omitempty"`
	// OauthClientSecret is the OAuth client secret.
	OauthClientSecret string `json:"oauthClientSecret,omitempty"`
	// OauthTenant is the OAuth tenant.
	OauthTenant string `json:"oauthTenant,omitempty"`
	// OauthScope is the scope of the OAuth token request. Since 2026.5.
	OauthScope string `json:"oauthScope,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation
// -----------------------------------------------------------------------------

// ValidateCreateEmailConfigurationOpt validates the options for the CreateEmailConfiguration method.
func (s *SystemServiceV2) ValidateCreateEmailConfigurationOpt(opt *SystemEmailConfigurationCreateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	for _, field := range []struct{ value, name string }{
		{opt.Host, "Host"},
		{opt.Port, "Port"},
		{opt.SecurityProtocol, "SecurityProtocol"},
		{opt.FromAddress, "FromAddress"},
		{opt.FromName, "FromName"},
		{opt.AuthMethod, "AuthMethod"},
		{opt.Username, "Username"},
	} {
		err := ValidateRequired(field.value, field.name)
		if err != nil {
			return err
		}
	}

	err := IsValueAuthorized(opt.SecurityProtocol, allowedEmailSecurityProtocols, "SecurityProtocol")
	if err != nil {
		return err
	}

	return IsValueAuthorized(opt.AuthMethod, allowedEmailAuthMethods, "AuthMethod")
}

// ValidateUpdateEmailConfigurationOpt validates the options for the UpdateEmailConfiguration method.
func (s *SystemServiceV2) ValidateUpdateEmailConfigurationOpt(opt *SystemEmailConfigurationUpdateOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	if opt.SecurityProtocol != "" {
		err := IsValueAuthorized(opt.SecurityProtocol, allowedEmailSecurityProtocols, "SecurityProtocol")
		if err != nil {
			return err
		}
	}

	if opt.AuthMethod != "" {
		return IsValueAuthorized(opt.AuthMethod, allowedEmailAuthMethods, "AuthMethod")
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// SearchEmailConfigurations lists the email configurations.
//
// API endpoint: GET /api/v2/system/email-configurations.
func (s *SystemServiceV2) SearchEmailConfigurations(ctx context.Context) (*SystemEmailConfigurationsSearch, *http.Response, error) {
	return callV2[SystemEmailConfigurationsSearch](ctx, s.client, http.MethodGet, "system/email-configurations", nil, nil)
}

// CreateEmailConfiguration creates an email configuration.
// Requires the 'Administer System' permission.
//
// API endpoint: POST /api/v2/system/email-configurations.
func (s *SystemServiceV2) CreateEmailConfiguration(ctx context.Context, opt *SystemEmailConfigurationCreateOptions) (*SystemEmailConfiguration, *http.Response, error) {
	err := s.ValidateCreateEmailConfigurationOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	return callV2[SystemEmailConfiguration](ctx, s.client, http.MethodPost, "system/email-configurations", nil, opt)
}

// GetEmailConfiguration fetches an email configuration.
//
// API endpoint: GET /api/v2/system/email-configurations/{id}.
func (s *SystemServiceV2) GetEmailConfiguration(ctx context.Context, id string) (*SystemEmailConfiguration, *http.Response, error) {
	path, err := v2PathWithID("system/email-configurations", id, "Id")
	if err != nil {
		return nil, nil, err
	}

	return callV2[SystemEmailConfiguration](ctx, s.client, http.MethodGet, path, nil, nil)
}

// UpdateEmailConfiguration updates an email configuration.
// Requires the 'Administer System' permission.
//
// API endpoint: PATCH /api/v2/system/email-configurations/{id}.
func (s *SystemServiceV2) UpdateEmailConfiguration(ctx context.Context, id string, opt *SystemEmailConfigurationUpdateOptions) (*SystemEmailConfiguration, *http.Response, error) {
	path, err := v2PathWithID("system/email-configurations", id, "Id")
	if err != nil {
		return nil, nil, err
	}

	err = s.ValidateUpdateEmailConfigurationOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	return callV2[SystemEmailConfiguration](ctx, s.client, http.MethodPatch, path, nil, opt)
}

// DeleteEmailConfiguration deletes an email configuration.
// Requires the 'Administer System' permission.
//
// API endpoint: DELETE /api/v2/system/email-configurations/{id}.
func (s *SystemServiceV2) DeleteEmailConfiguration(ctx context.Context, id string) (*http.Response, error) {
	path, err := v2PathWithID("system/email-configurations", id, "Id")
	if err != nil {
		return nil, err
	}

	return deleteV2(ctx, s.client, path)
}
