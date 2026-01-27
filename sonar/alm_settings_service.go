package sonargo

import "net/http"

const (
	// MaxAlmKeyLength is the maximum length for DevOps Platform setting keys.
	MaxAlmKeyLength = 200
	// MaxAlmURLLength is the maximum length for DevOps Platform URLs.
	MaxAlmURLLength = 2000
	// MaxPersonalAccessTokenLength is the maximum length for Personal Access Tokens.
	MaxPersonalAccessTokenLength = 2000
	// MaxGitHubAppIDLength is the maximum length for GitHub App ID.
	MaxGitHubAppIDLength = 80
	// MaxGitHubClientIDLength is the maximum length for GitHub Client ID.
	MaxGitHubClientIDLength = 80
	// MaxGitHubClientSecretLength is the maximum length for GitHub Client Secret.
	MaxGitHubClientSecretLength = 160
	// MaxGitHubPrivateKeyLength is the maximum length for GitHub App private key.
	MaxGitHubPrivateKeyLength = 2500
	// MaxGitHubWebhookSecretLength is the maximum length for GitHub App Webhook Secret.
	MaxGitHubWebhookSecretLength = 160
	// MaxBitbucketCloudClientIDLength is the maximum length for Bitbucket Cloud Client ID (create).
	MaxBitbucketCloudClientIDLength = 2000
	// MaxBitbucketCloudClientSecretLength is the maximum length for Bitbucket Cloud Client Secret (create).
	MaxBitbucketCloudClientSecretLength = 2000
	// MaxBitbucketCloudClientIDUpdateLength is the maximum length for Bitbucket Cloud Client ID (update).
	MaxBitbucketCloudClientIDUpdateLength = 80
	// MaxBitbucketCloudClientSecretUpdateLength is the maximum length for Bitbucket Cloud Client Secret (update).
	MaxBitbucketCloudClientSecretUpdateLength = 160
	// MaxBitbucketCloudWorkspaceUpdateLength is the maximum length for Bitbucket Cloud Workspace (update).
	MaxBitbucketCloudWorkspaceUpdateLength = 80
)

// AlmSettingsService handles communication with the DevOps Platform Settings related methods
// of the SonarQube API.
// This service manages configuration of Azure DevOps, Bitbucket, Bitbucket Cloud, GitHub, and GitLab
// integrations for project bindings.
type AlmSettingsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// AlmSettingsCountBinding represents the response from counting project bindings.
type AlmSettingsCountBinding struct {
	// Key is the DevOps Platform setting key.
	Key string `json:"key,omitempty"`
	// Projects is the number of projects bound to this setting.
	Projects int64 `json:"projects,omitempty"`
}

// AlmSettingsGetBinding represents the response from getting a project binding.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type AlmSettingsGetBinding struct {
	// Alm is the type of DevOps Platform (azure, bitbucket, bitbucketcloud, github, gitlab).
	Alm string `json:"alm,omitempty"`
	// InlineAnnotationsEnabled indicates if inline annotations are enabled (Azure only, since 2025.1).
	InlineAnnotationsEnabled bool `json:"inlineAnnotationsEnabled,omitempty"`
	// Key is the DevOps Platform setting key.
	Key string `json:"key,omitempty"`
	// Monorepo indicates if monorepo feature is enabled (Azure only, since 8.7).
	Monorepo bool `json:"monorepo,omitempty"`
	// Repository is the repository identifier.
	Repository string `json:"repository,omitempty"`
	// RepositoryURL is the URL to the repository (GitHub, GitLab, Azure since 2025.6).
	RepositoryURL string `json:"repositoryUrl,omitempty"`
	// Slug is the repository slug (Bitbucket).
	Slug string `json:"slug,omitempty"`
	// SummaryCommentEnabled indicates if summary comments are enabled.
	SummaryCommentEnabled bool `json:"summaryCommentEnabled,omitempty"`
	// URL is the DevOps Platform URL.
	URL string `json:"url,omitempty"`
}

// AlmSettingsList represents the response from listing available ALM settings.
type AlmSettingsList struct {
	// AlmSettings is the list of available DevOps Platform settings.
	AlmSettings []AlmSetting `json:"almSettings,omitempty"`
}

// AlmSetting represents a single DevOps Platform setting in the list.
type AlmSetting struct {
	// Alm is the type of DevOps Platform (azure, bitbucket, bitbucketcloud, github, gitlab).
	Alm string `json:"alm,omitempty"`
	// Key is the unique setting key.
	Key string `json:"key,omitempty"`
	// URL is the DevOps Platform URL.
	URL string `json:"url,omitempty"`
}

// AlmSettingsListDefinitions represents the response from listing all ALM setting definitions.
type AlmSettingsListDefinitions struct {
	// Azure contains Azure DevOps settings.
	Azure []AzureDefinition `json:"azure,omitempty"`
	// Bitbucket contains Bitbucket Server settings.
	Bitbucket []BitbucketDefinition `json:"bitbucket,omitempty"`
	// BitbucketCloud contains Bitbucket Cloud settings.
	BitbucketCloud []BitbucketCloudDefinition `json:"bitbucketcloud,omitempty"`
	// Github contains GitHub settings.
	Github []GithubDefinition `json:"github,omitempty"`
	// Gitlab contains GitLab settings.
	Gitlab []GitlabDefinition `json:"gitlab,omitempty"`
}

// AzureDefinition represents an Azure DevOps setting definition.
type AzureDefinition struct {
	// Key is the unique setting key.
	Key string `json:"key,omitempty"`
	// URL is the Azure DevOps API URL (since 8.6).
	URL string `json:"url,omitempty"`
}

// BitbucketDefinition represents a Bitbucket Server setting definition.
type BitbucketDefinition struct {
	// Key is the unique setting key.
	Key string `json:"key,omitempty"`
	// URL is the Bitbucket Server API URL.
	URL string `json:"url,omitempty"`
}

// BitbucketCloudDefinition represents a Bitbucket Cloud setting definition.
type BitbucketCloudDefinition struct {
	// ClientID is the OAuth client ID.
	ClientID string `json:"clientId,omitempty"`
	// Key is the unique setting key.
	Key string `json:"key,omitempty"`
	// Workspace is the Bitbucket Cloud workspace ID.
	Workspace string `json:"workspace,omitempty"`
}

// GithubDefinition represents a GitHub setting definition.
type GithubDefinition struct {
	// AppID is the GitHub App ID.
	AppID string `json:"appId,omitempty"`
	// ClientID is the GitHub App Client ID.
	ClientID string `json:"clientId,omitempty"`
	// Key is the unique setting key.
	Key string `json:"key,omitempty"`
	// URL is the GitHub API URL.
	URL string `json:"url,omitempty"`
}

// GitlabDefinition represents a GitLab setting definition.
type GitlabDefinition struct {
	// Key is the unique setting key.
	Key string `json:"key,omitempty"`
	// URL is the GitLab API URL (since 8.2).
	URL string `json:"url,omitempty"`
}

// AlmSettingsValidation represents the response from validating a DevOps Platform setting.
type AlmSettingsValidation struct {
	// Errors contains validation error messages, if any.
	Errors []AlmValidationError `json:"errors,omitempty"`
}

// AlmValidationError represents a single validation error.
type AlmValidationError struct {
	// Msg is the error message.
	Msg string `json:"msg,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// AlmSettingsCountBindingOption contains parameters for the CountBinding method.
type AlmSettingsCountBindingOption struct {
	// AlmSetting is the DevOps Platform setting key.
	// This field is required.
	AlmSetting string `url:"almSetting"`
}

// AlmSettingsCreateAzureOption contains parameters for the CreateAzure method.
type AlmSettingsCreateAzureOption struct {
	// Key is the unique key of the Azure DevOps instance setting.
	// This field is required. Maximum length: 200 characters.
	Key string `url:"key"`
	// PersonalAccessToken is the Azure DevOps personal access token.
	// This field is required. Maximum length: 2000 characters.
	PersonalAccessToken string `url:"personalAccessToken"`
	// URL is the Azure API URL.
	// This field is required. Maximum length: 2000 characters.
	URL string `url:"url"`
}

// AlmSettingsCreateBitbucketOption contains parameters for the CreateBitbucket method.
type AlmSettingsCreateBitbucketOption struct {
	// Key is the unique key of the Bitbucket instance setting.
	// This field is required. Maximum length: 200 characters.
	Key string `url:"key"`
	// PersonalAccessToken is the Bitbucket personal access token.
	// This field is required. Maximum length: 2000 characters.
	PersonalAccessToken string `url:"personalAccessToken"`
	// URL is the BitBucket Server API URL.
	// This field is required. Maximum length: 2000 characters.
	URL string `url:"url"`
}

// AlmSettingsCreateBitbucketCloudOption contains parameters for the CreateBitbucketCloud method.
type AlmSettingsCreateBitbucketCloudOption struct {
	// ClientID is the Bitbucket Cloud Client ID.
	// This field is required. Maximum length: 2000 characters.
	ClientID string `url:"clientId"`
	// ClientSecret is the Bitbucket Cloud Client Secret.
	// This field is required. Maximum length: 2000 characters.
	ClientSecret string `url:"clientSecret"`
	// Key is the unique key of the Bitbucket Cloud setting.
	// This field is required. Maximum length: 200 characters.
	Key string `url:"key"`
	// Workspace is the Bitbucket Cloud workspace ID.
	// This field is required.
	Workspace string `url:"workspace"`
}

// AlmSettingsCreateGithubOption contains parameters for the CreateGithub method.
type AlmSettingsCreateGithubOption struct {
	// AppID is the GitHub App ID.
	// This field is required. Maximum length: 80 characters.
	AppID string `url:"appId"`
	// ClientID is the GitHub App Client ID.
	// This field is required. Maximum length: 80 characters.
	ClientID string `url:"clientId"`
	// ClientSecret is the GitHub App Client Secret.
	// This field is required. Maximum length: 160 characters.
	ClientSecret string `url:"clientSecret"`
	// Key is the unique key of the GitHub instance setting.
	// This field is required. Maximum length: 200 characters.
	Key string `url:"key"`
	// PrivateKey is the GitHub App private key.
	// This field is required. Maximum length: 2500 characters.
	PrivateKey string `url:"privateKey"`
	// URL is the GitHub API URL.
	// This field is required. Maximum length: 2000 characters.
	URL string `url:"url"`
	// WebhookSecret is the GitHub App Webhook Secret.
	// This field is optional. Maximum length: 160 characters.
	WebhookSecret string `url:"webhookSecret,omitempty"`
}

// AlmSettingsCreateGitlabOption contains parameters for the CreateGitlab method.
type AlmSettingsCreateGitlabOption struct {
	// Key is the unique key of the GitLab instance setting.
	// This field is required. Maximum length: 200 characters.
	Key string `url:"key"`
	// PersonalAccessToken is the GitLab personal access token.
	// This field is required. Maximum length: 2000 characters.
	PersonalAccessToken string `url:"personalAccessToken"`
	// URL is the GitLab API URL.
	// This field is required. Maximum length: 2000 characters.
	URL string `url:"url"`
}

// AlmSettingsDeleteOption contains parameters for the Delete method.
type AlmSettingsDeleteOption struct {
	// Key is the DevOps Platform Setting key.
	// This field is required.
	Key string `url:"key"`
}

// AlmSettingsGetBindingOption contains parameters for the GetBinding method.
type AlmSettingsGetBindingOption struct {
	// Project is the project key.
	// This field is required.
	Project string `url:"project"`
}

// AlmSettingsListOption contains parameters for the List method.
type AlmSettingsListOption struct {
	// Project is the project key.
	// This field is optional.
	Project string `url:"project,omitempty"`
}

// AlmSettingsUpdateAzureOption contains parameters for the UpdateAzure method.
type AlmSettingsUpdateAzureOption struct {
	// Key is the unique key of the Azure instance setting.
	// This field is required. Maximum length: 200 characters.
	Key string `url:"key"`
	// NewKey is the optional new value for the unique key.
	// Maximum length: 200 characters.
	NewKey string `url:"newKey,omitempty"`
	// PersonalAccessToken is the Azure DevOps personal access token.
	// This field is optional (since 8.7). Maximum length: 2000 characters.
	PersonalAccessToken string `url:"personalAccessToken,omitempty"`
	// URL is the Azure API URL.
	// This field is required. Maximum length: 2000 characters.
	URL string `url:"url"`
}

// AlmSettingsUpdateBitbucketOption contains parameters for the UpdateBitbucket method.
type AlmSettingsUpdateBitbucketOption struct {
	// Key is the unique key of the Bitbucket instance setting.
	// This field is required. Maximum length: 200 characters.
	Key string `url:"key"`
	// NewKey is the optional new value for the unique key.
	// Maximum length: 200 characters.
	NewKey string `url:"newKey,omitempty"`
	// PersonalAccessToken is the Bitbucket personal access token.
	// This field is optional (since 8.7). Maximum length: 2000 characters.
	PersonalAccessToken string `url:"personalAccessToken,omitempty"`
	// URL is the Bitbucket API URL.
	// This field is required. Maximum length: 2000 characters.
	URL string `url:"url"`
}

// AlmSettingsUpdateBitbucketCloudOption contains parameters for the UpdateBitbucketCloud method.
type AlmSettingsUpdateBitbucketCloudOption struct {
	// ClientID is the Bitbucket Cloud Client ID.
	// This field is required. Maximum length: 80 characters.
	ClientID string `url:"clientId"`
	// ClientSecret is the optional new value for the Bitbucket Cloud client secret.
	// Maximum length: 160 characters.
	ClientSecret string `url:"clientSecret,omitempty"`
	// Key is the unique key of the Bitbucket Cloud setting.
	// This field is required. Maximum length: 200 characters.
	Key string `url:"key"`
	// NewKey is the optional new value for the unique key.
	// Maximum length: 200 characters.
	NewKey string `url:"newKey,omitempty"`
	// Workspace is the Bitbucket Cloud workspace ID.
	// This field is required. Maximum length: 80 characters.
	Workspace string `url:"workspace"`
}

// AlmSettingsUpdateGithubOption contains parameters for the UpdateGithub method.
type AlmSettingsUpdateGithubOption struct {
	// AppID is the GitHub API ID.
	// This field is required. Maximum length: 80 characters.
	AppID string `url:"appId"`
	// ClientID is the GitHub App Client ID.
	// This field is required. Maximum length: 80 characters.
	ClientID string `url:"clientId"`
	// ClientSecret is the GitHub App Client Secret.
	// This field is optional (since 8.7). Maximum length: 160 characters.
	ClientSecret string `url:"clientSecret,omitempty"`
	// Key is the unique key of the GitHub instance setting.
	// This field is required. Maximum length: 200 characters.
	Key string `url:"key"`
	// NewKey is the optional new value for the unique key.
	// Maximum length: 200 characters.
	NewKey string `url:"newKey,omitempty"`
	// PrivateKey is the GitHub App private key.
	// This field is optional (since 8.7). Maximum length: 2500 characters.
	PrivateKey string `url:"privateKey,omitempty"`
	// URL is the GitHub API URL.
	// This field is required. Maximum length: 2000 characters.
	URL string `url:"url"`
	// WebhookSecret is the GitHub App Webhook Secret.
	// This field is optional. Maximum length: 160 characters.
	WebhookSecret string `url:"webhookSecret,omitempty"`
}

// AlmSettingsUpdateGitlabOption contains parameters for the UpdateGitlab method.
type AlmSettingsUpdateGitlabOption struct {
	// Key is the unique key of the GitLab instance setting.
	// This field is required. Maximum length: 200 characters.
	Key string `url:"key"`
	// NewKey is the optional new value for the unique key.
	// Maximum length: 200 characters.
	NewKey string `url:"newKey,omitempty"`
	// PersonalAccessToken is the GitLab personal access token.
	// This field is optional (since 8.7). Maximum length: 2000 characters.
	PersonalAccessToken string `url:"personalAccessToken,omitempty"`
	// URL is the GitLab API URL.
	// This field is required. Maximum length: 2000 characters.
	URL string `url:"url"`
}

// AlmSettingsValidateOption contains parameters for the Validate method.
type AlmSettingsValidateOption struct {
	// Key is the unique key of the DevOps Platform settings.
	// This field is required. Maximum length: 200 characters.
	Key string `url:"key"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateCountBindingOpt validates the options for the CountBinding method.
func (s *AlmSettingsService) ValidateCountBindingOpt(opt *AlmSettingsCountBindingOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AlmSetting, "AlmSetting")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCreateAzureOpt validates the options for the CreateAzure method.
func (s *AlmSettingsService) ValidateCreateAzureOpt(opt *AlmSettingsCreateAzureOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Key, MaxAlmKeyLength, "Key")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.PersonalAccessToken, "PersonalAccessToken")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.PersonalAccessToken, MaxPersonalAccessTokenLength, "PersonalAccessToken")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.URL, "URL")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.URL, MaxAlmURLLength, "URL")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCreateBitbucketOpt validates the options for the CreateBitbucket method.
func (s *AlmSettingsService) ValidateCreateBitbucketOpt(opt *AlmSettingsCreateBitbucketOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Key, MaxAlmKeyLength, "Key")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.PersonalAccessToken, "PersonalAccessToken")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.PersonalAccessToken, MaxPersonalAccessTokenLength, "PersonalAccessToken")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.URL, "URL")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.URL, MaxAlmURLLength, "URL")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCreateBitbucketCloudOpt validates the options for the CreateBitbucketCloud method.
func (s *AlmSettingsService) ValidateCreateBitbucketCloudOpt(opt *AlmSettingsCreateBitbucketCloudOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ClientID, "ClientID")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.ClientID, MaxBitbucketCloudClientIDLength, "ClientID")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ClientSecret, "ClientSecret")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.ClientSecret, MaxBitbucketCloudClientSecretLength, "ClientSecret")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Key, MaxAlmKeyLength, "Key")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Workspace, "Workspace")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCreateGithubOpt validates the options for the CreateGithub method.
//
//nolint:cyclop,funlen // Validation functions are naturally complex due to multiple checks
func (s *AlmSettingsService) ValidateCreateGithubOpt(opt *AlmSettingsCreateGithubOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AppID, "AppID")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.AppID, MaxGitHubAppIDLength, "AppID")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ClientID, "ClientID")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.ClientID, MaxGitHubClientIDLength, "ClientID")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ClientSecret, "ClientSecret")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.ClientSecret, MaxGitHubClientSecretLength, "ClientSecret")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Key, MaxAlmKeyLength, "Key")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.PrivateKey, "PrivateKey")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.PrivateKey, MaxGitHubPrivateKeyLength, "PrivateKey")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.URL, "URL")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.URL, MaxAlmURLLength, "URL")
	if err != nil {
		return err
	}

	if opt.WebhookSecret != "" {
		err = ValidateMaxLength(opt.WebhookSecret, MaxGitHubWebhookSecretLength, "WebhookSecret")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateCreateGitlabOpt validates the options for the CreateGitlab method.
func (s *AlmSettingsService) ValidateCreateGitlabOpt(opt *AlmSettingsCreateGitlabOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Key, MaxAlmKeyLength, "Key")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.PersonalAccessToken, "PersonalAccessToken")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.PersonalAccessToken, MaxPersonalAccessTokenLength, "PersonalAccessToken")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.URL, "URL")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.URL, MaxAlmURLLength, "URL")
	if err != nil {
		return err
	}

	return nil
}

// ValidateDeleteOpt validates the options for the Delete method.
func (s *AlmSettingsService) ValidateDeleteOpt(opt *AlmSettingsDeleteOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	return nil
}

// ValidateGetBindingOpt validates the options for the GetBinding method.
func (s *AlmSettingsService) ValidateGetBindingOpt(opt *AlmSettingsGetBindingOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	return nil
}

// ValidateListOpt validates the options for the List method.
func (s *AlmSettingsService) ValidateListOpt(opt *AlmSettingsListOption) error {
	// Options are optional; nothing to validate.
	return nil
}

// ValidateUpdateAzureOpt validates the options for the UpdateAzure method.
func (s *AlmSettingsService) ValidateUpdateAzureOpt(opt *AlmSettingsUpdateAzureOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Key, MaxAlmKeyLength, "Key")
	if err != nil {
		return err
	}

	if opt.NewKey != "" {
		err = ValidateMaxLength(opt.NewKey, MaxAlmKeyLength, "NewKey")
		if err != nil {
			return err
		}
	}

	if opt.PersonalAccessToken != "" {
		err = ValidateMaxLength(opt.PersonalAccessToken, MaxPersonalAccessTokenLength, "PersonalAccessToken")
		if err != nil {
			return err
		}
	}

	err = ValidateRequired(opt.URL, "URL")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.URL, MaxAlmURLLength, "URL")
	if err != nil {
		return err
	}

	return nil
}

// ValidateUpdateBitbucketOpt validates the options for the UpdateBitbucket method.
func (s *AlmSettingsService) ValidateUpdateBitbucketOpt(opt *AlmSettingsUpdateBitbucketOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Key, MaxAlmKeyLength, "Key")
	if err != nil {
		return err
	}

	if opt.NewKey != "" {
		err = ValidateMaxLength(opt.NewKey, MaxAlmKeyLength, "NewKey")
		if err != nil {
			return err
		}
	}

	if opt.PersonalAccessToken != "" {
		err = ValidateMaxLength(opt.PersonalAccessToken, MaxPersonalAccessTokenLength, "PersonalAccessToken")
		if err != nil {
			return err
		}
	}

	err = ValidateRequired(opt.URL, "URL")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.URL, MaxAlmURLLength, "URL")
	if err != nil {
		return err
	}

	return nil
}

// ValidateUpdateBitbucketCloudOpt validates the options for the UpdateBitbucketCloud method.
//
//nolint:cyclop // Validation functions are naturally complex due to multiple checks
func (s *AlmSettingsService) ValidateUpdateBitbucketCloudOpt(opt *AlmSettingsUpdateBitbucketCloudOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ClientID, "ClientID")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.ClientID, MaxBitbucketCloudClientIDUpdateLength, "ClientID")
	if err != nil {
		return err
	}

	if opt.ClientSecret != "" {
		err = ValidateMaxLength(opt.ClientSecret, MaxBitbucketCloudClientSecretUpdateLength, "ClientSecret")
		if err != nil {
			return err
		}
	}

	err = ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Key, MaxAlmKeyLength, "Key")
	if err != nil {
		return err
	}

	if opt.NewKey != "" {
		err = ValidateMaxLength(opt.NewKey, MaxAlmKeyLength, "NewKey")
		if err != nil {
			return err
		}
	}

	err = ValidateRequired(opt.Workspace, "Workspace")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Workspace, MaxBitbucketCloudWorkspaceUpdateLength, "Workspace")
	if err != nil {
		return err
	}

	return nil
}

// ValidateUpdateGithubOpt validates the options for the UpdateGithub method.
//
//nolint:cyclop,funlen // Validation functions are naturally complex due to multiple checks
func (s *AlmSettingsService) ValidateUpdateGithubOpt(opt *AlmSettingsUpdateGithubOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AppID, "AppID")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.AppID, MaxGitHubAppIDLength, "AppID")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ClientID, "ClientID")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.ClientID, MaxGitHubClientIDLength, "ClientID")
	if err != nil {
		return err
	}

	if opt.ClientSecret != "" {
		err = ValidateMaxLength(opt.ClientSecret, MaxGitHubClientSecretLength, "ClientSecret")
		if err != nil {
			return err
		}
	}

	err = ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Key, MaxAlmKeyLength, "Key")
	if err != nil {
		return err
	}

	if opt.NewKey != "" {
		err = ValidateMaxLength(opt.NewKey, MaxAlmKeyLength, "NewKey")
		if err != nil {
			return err
		}
	}

	if opt.PrivateKey != "" {
		err = ValidateMaxLength(opt.PrivateKey, MaxGitHubPrivateKeyLength, "PrivateKey")
		if err != nil {
			return err
		}
	}

	err = ValidateRequired(opt.URL, "URL")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.URL, MaxAlmURLLength, "URL")
	if err != nil {
		return err
	}

	if opt.WebhookSecret != "" {
		err = ValidateMaxLength(opt.WebhookSecret, MaxGitHubWebhookSecretLength, "WebhookSecret")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateUpdateGitlabOpt validates the options for the UpdateGitlab method.
func (s *AlmSettingsService) ValidateUpdateGitlabOpt(opt *AlmSettingsUpdateGitlabOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Key, MaxAlmKeyLength, "Key")
	if err != nil {
		return err
	}

	if opt.NewKey != "" {
		err = ValidateMaxLength(opt.NewKey, MaxAlmKeyLength, "NewKey")
		if err != nil {
			return err
		}
	}

	if opt.PersonalAccessToken != "" {
		err = ValidateMaxLength(opt.PersonalAccessToken, MaxPersonalAccessTokenLength, "PersonalAccessToken")
		if err != nil {
			return err
		}
	}

	err = ValidateRequired(opt.URL, "URL")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.URL, MaxAlmURLLength, "URL")
	if err != nil {
		return err
	}

	return nil
}

// ValidateValidateOpt validates the options for the Validate method.
func (s *AlmSettingsService) ValidateValidateOpt(opt *AlmSettingsValidateOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Key, MaxAlmKeyLength, "Key")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// CountBinding counts the number of projects bound to a DevOps Platform setting.
// Requires the 'Administer System' permission.
//
// API endpoint: GET /api/alm_settings/count_binding.
// Since: 8.1.
func (s *AlmSettingsService) CountBinding(opt *AlmSettingsCountBindingOption) (*AlmSettingsCountBinding, *http.Response, error) {
	err := s.ValidateCountBindingOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "alm_settings/count_binding", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(AlmSettingsCountBinding)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateAzure creates an Azure DevOps instance setting.
// Requires the 'Administer System' permission.
//
// API endpoint: POST /api/alm_settings/create_azure.
// Since: 8.1.
func (s *AlmSettingsService) CreateAzure(opt *AlmSettingsCreateAzureOption) (*http.Response, error) {
	err := s.ValidateCreateAzureOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "alm_settings/create_azure", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// CreateBitbucket creates a Bitbucket Server instance setting.
// Requires the 'Administer System' permission.
//
// API endpoint: POST /api/alm_settings/create_bitbucket.
// Since: 8.1.
func (s *AlmSettingsService) CreateBitbucket(opt *AlmSettingsCreateBitbucketOption) (*http.Response, error) {
	err := s.ValidateCreateBitbucketOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "alm_settings/create_bitbucket", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// CreateBitbucketCloud configures a new instance of Bitbucket Cloud.
// Requires the 'Administer System' permission.
//
// API endpoint: POST /api/alm_settings/create_bitbucketcloud.
// Since: 8.7.
func (s *AlmSettingsService) CreateBitbucketCloud(opt *AlmSettingsCreateBitbucketCloudOption) (*http.Response, error) {
	err := s.ValidateCreateBitbucketCloudOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "alm_settings/create_bitbucketcloud", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// CreateGithub creates a GitHub instance setting.
// Requires the 'Administer System' permission.
//
// API endpoint: POST /api/alm_settings/create_github.
// Since: 8.1.
func (s *AlmSettingsService) CreateGithub(opt *AlmSettingsCreateGithubOption) (*http.Response, error) {
	err := s.ValidateCreateGithubOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "alm_settings/create_github", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// CreateGitlab creates a GitLab instance setting.
// Requires the 'Administer System' permission.
//
// API endpoint: POST /api/alm_settings/create_gitlab.
// Since: 8.1.
func (s *AlmSettingsService) CreateGitlab(opt *AlmSettingsCreateGitlabOption) (*http.Response, error) {
	err := s.ValidateCreateGitlabOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "alm_settings/create_gitlab", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Delete deletes a DevOps Platform setting.
// Requires the 'Administer System' permission.
//
// API endpoint: POST /api/alm_settings/delete.
// Since: 8.1.
func (s *AlmSettingsService) Delete(opt *AlmSettingsDeleteOption) (*http.Response, error) {
	err := s.ValidateDeleteOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "alm_settings/delete", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetBinding gets the DevOps Platform binding of a given project.
// Requires the 'Browse' permission on the project.
//
// API endpoint: GET /api/alm_settings/get_binding.
// Since: 8.1.
func (s *AlmSettingsService) GetBinding(opt *AlmSettingsGetBindingOption) (*AlmSettingsGetBinding, *http.Response, error) {
	err := s.ValidateGetBindingOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "alm_settings/get_binding", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(AlmSettingsGetBinding)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// List lists DevOps Platform settings available for a given project, sorted by DevOps Platform key.
// Requires the 'Administer project' permission if the 'project' parameter is provided,
// requires the 'Create Projects' permission otherwise.
//
// API endpoint: GET /api/alm_settings/list.
// Since: 8.1.
func (s *AlmSettingsService) List(opt *AlmSettingsListOption) (*AlmSettingsList, *http.Response, error) {
	err := s.ValidateListOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "alm_settings/list", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(AlmSettingsList)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ListDefinitions lists DevOps Platform settings, sorted by created date.
// Requires the 'Administer System' permission.
//
// API endpoint: GET /api/alm_settings/list_definitions.
// Since: 8.1.
func (s *AlmSettingsService) ListDefinitions() (*AlmSettingsListDefinitions, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "alm_settings/list_definitions", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(AlmSettingsListDefinitions)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateAzure updates an Azure DevOps instance setting.
// Requires the 'Administer System' permission.
//
// API endpoint: POST /api/alm_settings/update_azure.
// Since: 8.1.
func (s *AlmSettingsService) UpdateAzure(opt *AlmSettingsUpdateAzureOption) (*http.Response, error) {
	err := s.ValidateUpdateAzureOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "alm_settings/update_azure", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UpdateBitbucket updates a Bitbucket Server instance setting.
// Requires the 'Administer System' permission.
//
// API endpoint: POST /api/alm_settings/update_bitbucket.
// Since: 8.1.
func (s *AlmSettingsService) UpdateBitbucket(opt *AlmSettingsUpdateBitbucketOption) (*http.Response, error) {
	err := s.ValidateUpdateBitbucketOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "alm_settings/update_bitbucket", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UpdateBitbucketCloud updates a Bitbucket Cloud setting.
// Requires the 'Administer System' permission.
//
// API endpoint: POST /api/alm_settings/update_bitbucketcloud.
// Since: 8.7.
func (s *AlmSettingsService) UpdateBitbucketCloud(opt *AlmSettingsUpdateBitbucketCloudOption) (*http.Response, error) {
	err := s.ValidateUpdateBitbucketCloudOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "alm_settings/update_bitbucketcloud", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UpdateGithub updates a GitHub instance setting.
// Requires the 'Administer System' permission.
//
// API endpoint: POST /api/alm_settings/update_github.
// Since: 8.1.
func (s *AlmSettingsService) UpdateGithub(opt *AlmSettingsUpdateGithubOption) (*http.Response, error) {
	err := s.ValidateUpdateGithubOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "alm_settings/update_github", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// UpdateGitlab updates a GitLab instance setting.
// Requires the 'Administer System' permission.
//
// API endpoint: POST /api/alm_settings/update_gitlab.
// Since: 8.1.
func (s *AlmSettingsService) UpdateGitlab(opt *AlmSettingsUpdateGitlabOption) (*http.Response, error) {
	err := s.ValidateUpdateGitlabOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "alm_settings/update_gitlab", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Validate validates a DevOps Platform setting by checking connectivity and permissions.
// Requires the 'Administer System' permission.
//
// API endpoint: GET /api/alm_settings/validate.
// Since: 8.6.
func (s *AlmSettingsService) Validate(opt *AlmSettingsValidateOption) (*AlmSettingsValidation, *http.Response, error) {
	err := s.ValidateValidateOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "alm_settings/validate", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(AlmSettingsValidation)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
