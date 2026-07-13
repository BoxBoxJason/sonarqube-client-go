package sonar

import (
	"context"
	"net/http"
)

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
	// MaxGitHubRepositoryLength is the maximum length for a GitHub repository binding.
	MaxGitHubRepositoryLength = 256
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

// AlmSettingsGithubManifest represents the response from initiating the GitHub App Manifest flow.
type AlmSettingsGithubManifest struct {
	// Manifest is the GitHub App Manifest JSON to be POSTed to GitHub by the browser.
	Manifest string `json:"manifest,omitempty"`
	// State is the single-use state token to be POSTed to GitHub alongside the manifest.
	State string `json:"state,omitempty"`
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

// AlmSettingsCountBindingOptions contains parameters for the CountBinding method.
type AlmSettingsCountBindingOptions struct {
	// AlmSetting is the DevOps Platform setting key.
	// This field is required.
	AlmSetting string `url:"almSetting"`
}

// AlmSettingsCreateAzureOptions contains parameters for the CreateAzure method.
type AlmSettingsCreateAzureOptions struct {
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

// AlmSettingsCreateBitbucketOptions contains parameters for the CreateBitbucket method.
type AlmSettingsCreateBitbucketOptions struct {
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

// AlmSettingsCreateBitbucketCloudOptions contains parameters for the CreateBitbucketCloud method.
type AlmSettingsCreateBitbucketCloudOptions struct {
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

// AlmSettingsCreateGithubOptions contains parameters for the CreateGithub method.
type AlmSettingsCreateGithubOptions struct {
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

// AlmSettingsCreateGitlabOptions contains parameters for the CreateGitlab method.
type AlmSettingsCreateGitlabOptions struct {
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

// AlmSettingsCreateGithubFromManifestOptions contains parameters for the CreateGithubFromManifest method.
type AlmSettingsCreateGithubFromManifestOptions struct {
	// Auth indicates whether to also set up GitHub authentication (sign-in) using this App.
	// Allowed values: true, false, yes, no
	// This field is optional. Default: false.
	Auth string `url:"auth,omitempty"`
	// Devops indicates whether to create the DevOps Platform integration (project import / PR analysis) for this App.
	// Allowed values: true, false, yes, no
	// This field is optional. Default: true.
	Devops string `url:"devops,omitempty"`
	// Key is the unique key of the GitHub instance setting that will be created. Required when Devops is true.
	// This field is optional. Maximum length: 200 characters.
	Key string `url:"key,omitempty"`
	// Name is the suggested name for the GitHub App (the user can change it on GitHub).
	// Defaults to 'SonarQube - <add_unique_name>'.
	// This field is optional. Maximum length: 200 characters.
	Name string `url:"name,omitempty"`
	// Organization is the GitHub organization the App should be created under.
	// Leave empty to create it under the user's personal account.
	// This field is optional. Maximum length: 200 characters.
	Organization string `url:"organization,omitempty"`
}

// AlmSettingsDeleteOptions contains parameters for the Delete method.
type AlmSettingsDeleteOptions struct {
	// Key is the DevOps Platform Setting key.
	// This field is required.
	Key string `url:"key"`
}

// AlmSettingsDeleteBindingOptions contains parameters for the DeleteBinding method.
type AlmSettingsDeleteBindingOptions struct {
	// Project is the project key.
	// This field is required.
	Project string `url:"project"`
}

// AlmSettingsGetBindingOptions contains parameters for the GetBinding method.
type AlmSettingsGetBindingOptions struct {
	// Project is the project key.
	// This field is required.
	Project string `url:"project"`
}

// AlmSettingsSetAzureBindingOptions contains parameters for the SetAzureBinding method.
type AlmSettingsSetAzureBindingOptions struct {
	// InlineAnnotationsEnabled enables inline annotations during Pull Request decoration.
	// This field is optional (since 2025.1). Default: true.
	InlineAnnotationsEnabled *bool `url:"inlineAnnotationsEnabled,omitempty"`
	// AlmSetting is the Azure DevOps setting key.
	// This field is required. Maximum length: 200 characters.
	AlmSetting string `url:"almSetting"`
	// Project is the SonarQube project key.
	// This field is required.
	Project string `url:"project"`
	// ProjectName is the Azure DevOps project name.
	// This field is required (since 8.6).
	ProjectName string `url:"projectName"`
	// RepositoryName is the Azure DevOps repository name.
	// This field is required (since 8.6).
	RepositoryName string `url:"repositoryName"`
	// Monorepo indicates if this project is part of a monorepo.
	// This field is required (since 8.7).
	Monorepo bool `url:"monorepo"`
}

// AlmSettingsSetBitbucketBindingOptions contains parameters for the SetBitbucketBinding method.
type AlmSettingsSetBitbucketBindingOptions struct {
	// AlmSetting is the Bitbucket Server setting key.
	// This field is required. Maximum length: 200 characters.
	AlmSetting string `url:"almSetting"`
	// Project is the project key.
	// This field is required.
	Project string `url:"project"`
	// Repository is the Bitbucket Server repository key.
	// This field is required.
	Repository string `url:"repository"`
	// Slug is the Bitbucket Server repository slug.
	// This field is required.
	Slug string `url:"slug"`
	// Monorepo indicates if this project is part of a monorepo.
	// This field is required (since 8.7).
	Monorepo bool `url:"monorepo"`
}

// AlmSettingsSetBitbucketCloudBindingOptions contains parameters for the SetBitbucketCloudBinding method.
type AlmSettingsSetBitbucketCloudBindingOptions struct {
	// AlmSetting is the Bitbucket Cloud setting key.
	// This field is required. Maximum length: 200 characters.
	AlmSetting string `url:"almSetting"`
	// Project is the project key.
	// This field is required.
	Project string `url:"project"`
	// Repository is the Bitbucket Cloud repository key.
	// This field is required.
	Repository string `url:"repository"`
	// Monorepo indicates if this project is part of a monorepo.
	// This field is required (since 8.8).
	Monorepo bool `url:"monorepo"`
}

// AlmSettingsSetGithubBindingOptions contains parameters for the SetGithubBinding method.
type AlmSettingsSetGithubBindingOptions struct {
	// SummaryCommentEnabled enables/disables the analysis summary in the PR discussion tab.
	// This field is optional (since 8.3). Default: true.
	SummaryCommentEnabled *bool `url:"summaryCommentEnabled,omitempty"`
	// AlmSetting is the GitHub setting key.
	// This field is required. Maximum length: 200 characters.
	AlmSetting string `url:"almSetting"`
	// Project is the project key.
	// This field is required.
	Project string `url:"project"`
	// Repository is the GitHub repository.
	// This field is required. Maximum length: 256 characters.
	Repository string `url:"repository"`
	// Monorepo indicates if this project is part of a monorepo.
	// This field is required (since 8.7).
	Monorepo bool `url:"monorepo"`
}

// AlmSettingsSetGitlabBindingOptions contains parameters for the SetGitlabBinding method.
type AlmSettingsSetGitlabBindingOptions struct {
	// AlmSetting is the GitLab setting key.
	// This field is required. Maximum length: 200 characters.
	AlmSetting string `url:"almSetting"`
	// Project is the project key.
	// This field is required.
	Project string `url:"project"`
	// Repository is the GitLab project ID.
	// This field is required.
	Repository string `url:"repository"`
	// Monorepo indicates if this project is part of a monorepo.
	// This field is required (since 8.7).
	Monorepo bool `url:"monorepo"`
}

// AlmSettingsListOptions contains parameters for the List method.
type AlmSettingsListOptions struct {
	// Project is the project key.
	// This field is optional.
	Project string `url:"project,omitempty"`
}

// AlmSettingsUpdateAzureOptions contains parameters for the UpdateAzure method.
type AlmSettingsUpdateAzureOptions struct {
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

// AlmSettingsUpdateBitbucketOptions contains parameters for the UpdateBitbucket method.
type AlmSettingsUpdateBitbucketOptions struct {
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

// AlmSettingsUpdateBitbucketCloudOptions contains parameters for the UpdateBitbucketCloud method.
type AlmSettingsUpdateBitbucketCloudOptions struct {
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

// AlmSettingsUpdateGithubOptions contains parameters for the UpdateGithub method.
type AlmSettingsUpdateGithubOptions struct {
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

// AlmSettingsUpdateGitlabOptions contains parameters for the UpdateGitlab method.
type AlmSettingsUpdateGitlabOptions struct {
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

// AlmSettingsValidateOptions contains parameters for the Validate method.
type AlmSettingsValidateOptions struct {
	// Key is the unique key of the DevOps Platform settings.
	// This field is required. Maximum length: 200 characters.
	Key string `url:"key"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateCountBindingOpt validates the options for the CountBinding method.
func (s *AlmSettingsService) ValidateCountBindingOpt(opt *AlmSettingsCountBindingOptions) error {
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
func (s *AlmSettingsService) ValidateCreateAzureOpt(opt *AlmSettingsCreateAzureOptions) error {
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
func (s *AlmSettingsService) ValidateCreateBitbucketOpt(opt *AlmSettingsCreateBitbucketOptions) error {
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
func (s *AlmSettingsService) ValidateCreateBitbucketCloudOpt(opt *AlmSettingsCreateBitbucketCloudOptions) error {
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
func (s *AlmSettingsService) ValidateCreateGithubOpt(opt *AlmSettingsCreateGithubOptions) error {
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
func (s *AlmSettingsService) ValidateCreateGitlabOpt(opt *AlmSettingsCreateGitlabOptions) error {
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

// ValidateCreateGithubFromManifestOpt validates the options for the CreateGithubFromManifest method.
func (s *AlmSettingsService) ValidateCreateGithubFromManifestOpt(opt *AlmSettingsCreateGithubFromManifestOptions) error {
	if opt == nil {
		// Options are optional; nothing to validate.
		return nil
	}

	if opt.Key != "" {
		err := ValidateMaxLength(opt.Key, MaxAlmKeyLength, "Key")
		if err != nil {
			return err
		}
	}

	if opt.Name != "" {
		err := ValidateMaxLength(opt.Name, MaxAlmKeyLength, "Name")
		if err != nil {
			return err
		}
	}

	if opt.Organization != "" {
		err := ValidateMaxLength(opt.Organization, MaxAlmKeyLength, "Organization")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateDeleteOpt validates the options for the Delete method.
func (s *AlmSettingsService) ValidateDeleteOpt(opt *AlmSettingsDeleteOptions) error {
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
func (s *AlmSettingsService) ValidateGetBindingOpt(opt *AlmSettingsGetBindingOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	return nil
}

// ValidateDeleteBindingOpt validates the options for the DeleteBinding method.
func (s *AlmSettingsService) ValidateDeleteBindingOpt(opt *AlmSettingsDeleteBindingOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSetAzureBindingOpt validates the options for the SetAzureBinding method.
func (s *AlmSettingsService) ValidateSetAzureBindingOpt(opt *AlmSettingsSetAzureBindingOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AlmSetting, "AlmSetting")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.AlmSetting, MaxAlmKeyLength, "AlmSetting")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ProjectName, "ProjectName")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.RepositoryName, "RepositoryName")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSetBitbucketBindingOpt validates the options for the SetBitbucketBinding method.
func (s *AlmSettingsService) ValidateSetBitbucketBindingOpt(opt *AlmSettingsSetBitbucketBindingOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AlmSetting, "AlmSetting")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.AlmSetting, MaxAlmKeyLength, "AlmSetting")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Repository, "Repository")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Slug, "Slug")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSetBitbucketCloudBindingOpt validates the options for the SetBitbucketCloudBinding method.
func (s *AlmSettingsService) ValidateSetBitbucketCloudBindingOpt(opt *AlmSettingsSetBitbucketCloudBindingOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AlmSetting, "AlmSetting")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.AlmSetting, MaxAlmKeyLength, "AlmSetting")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Repository, "Repository")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSetGithubBindingOpt validates the options for the SetGithubBinding method.
func (s *AlmSettingsService) ValidateSetGithubBindingOpt(opt *AlmSettingsSetGithubBindingOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AlmSetting, "AlmSetting")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.AlmSetting, MaxAlmKeyLength, "AlmSetting")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Repository, "Repository")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Repository, MaxGitHubRepositoryLength, "Repository")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSetGitlabBindingOpt validates the options for the SetGitlabBinding method.
func (s *AlmSettingsService) ValidateSetGitlabBindingOpt(opt *AlmSettingsSetGitlabBindingOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AlmSetting, "AlmSetting")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.AlmSetting, MaxAlmKeyLength, "AlmSetting")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Repository, "Repository")
	if err != nil {
		return err
	}

	return nil
}

// ValidateListOpt validates the options for the List method.
func (s *AlmSettingsService) ValidateListOpt(opt *AlmSettingsListOptions) error {
	// Options are optional; nothing to validate.
	return nil
}

// ValidateUpdateAzureOpt validates the options for the UpdateAzure method.
func (s *AlmSettingsService) ValidateUpdateAzureOpt(opt *AlmSettingsUpdateAzureOptions) error {
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
func (s *AlmSettingsService) ValidateUpdateBitbucketOpt(opt *AlmSettingsUpdateBitbucketOptions) error {
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
func (s *AlmSettingsService) ValidateUpdateBitbucketCloudOpt(opt *AlmSettingsUpdateBitbucketCloudOptions) error {
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
func (s *AlmSettingsService) ValidateUpdateGithubOpt(opt *AlmSettingsUpdateGithubOptions) error {
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
func (s *AlmSettingsService) ValidateUpdateGitlabOpt(opt *AlmSettingsUpdateGitlabOptions) error {
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
func (s *AlmSettingsService) ValidateValidateOpt(opt *AlmSettingsValidateOptions) error {
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
func (s *AlmSettingsService) CountBinding(ctx context.Context, opt *AlmSettingsCountBindingOptions) (*AlmSettingsCountBinding, *http.Response, error) {
	err := s.ValidateCountBindingOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "alm_settings/count_binding", opt)
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
func (s *AlmSettingsService) CreateAzure(ctx context.Context, opt *AlmSettingsCreateAzureOptions) (*http.Response, error) {
	err := s.ValidateCreateAzureOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "alm_settings/create_azure", opt)
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
func (s *AlmSettingsService) CreateBitbucket(ctx context.Context, opt *AlmSettingsCreateBitbucketOptions) (*http.Response, error) {
	err := s.ValidateCreateBitbucketOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "alm_settings/create_bitbucket", opt)
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
func (s *AlmSettingsService) CreateBitbucketCloud(ctx context.Context, opt *AlmSettingsCreateBitbucketCloudOptions) (*http.Response, error) {
	err := s.ValidateCreateBitbucketCloudOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "alm_settings/create_bitbucketcloud", opt)
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
func (s *AlmSettingsService) CreateGithub(ctx context.Context, opt *AlmSettingsCreateGithubOptions) (*http.Response, error) {
	err := s.ValidateCreateGithubOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "alm_settings/create_github", opt)
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
func (s *AlmSettingsService) CreateGitlab(ctx context.Context, opt *AlmSettingsCreateGitlabOptions) (*http.Response, error) {
	err := s.ValidateCreateGitlabOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "alm_settings/create_gitlab", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// CreateGithubFromManifest creates a GitHub App configuration from a manifest.
// Requires the 'Administer System' permission.
//
// API endpoint: POST /api/alm_settings/create_github_from_manifest.
// Since: 2026.4.
func (s *AlmSettingsService) CreateGithubFromManifest(ctx context.Context, opt *AlmSettingsCreateGithubFromManifestOptions) (*AlmSettingsGithubManifest, *http.Response, error) {
	err := s.ValidateCreateGithubFromManifestOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "alm_settings/create_github_from_manifest", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(AlmSettingsGithubManifest)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Delete deletes a DevOps Platform setting.
// Requires the 'Administer System' permission.
//
// API endpoint: POST /api/alm_settings/delete.
// Since: 8.1.
func (s *AlmSettingsService) Delete(ctx context.Context, opt *AlmSettingsDeleteOptions) (*http.Response, error) {
	err := s.ValidateDeleteOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "alm_settings/delete", opt)
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
func (s *AlmSettingsService) GetBinding(ctx context.Context, opt *AlmSettingsGetBindingOptions) (*AlmSettingsGetBinding, *http.Response, error) {
	err := s.ValidateGetBindingOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "alm_settings/get_binding", opt)
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

// DeleteBinding deletes the DevOps Platform binding of a project.
// Requires the 'Administer' permission on the project.
//
// API endpoint: POST /api/alm_settings/delete_binding.
// Since: 8.1.
func (s *AlmSettingsService) DeleteBinding(ctx context.Context, opt *AlmSettingsDeleteBindingOptions) (*http.Response, error) {
	err := s.ValidateDeleteBindingOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "alm_settings/delete_binding", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// SetAzureBinding binds an Azure DevOps instance to a project.
// If the project is already bound to a previous Azure DevOps instance, the binding will be updated to the new one.
// Requires the 'Administer' permission on the project.
//
// API endpoint: POST /api/alm_settings/set_azure_binding.
// Since: 8.1.
func (s *AlmSettingsService) SetAzureBinding(ctx context.Context, opt *AlmSettingsSetAzureBindingOptions) (*http.Response, error) {
	err := s.ValidateSetAzureBindingOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "alm_settings/set_azure_binding", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// SetBitbucketBinding binds a Bitbucket Server instance to a project.
// If the project is already bound to a previous Bitbucket instance, the binding will be updated to the new one.
// Requires the 'Administer' permission on the project.
//
// API endpoint: POST /api/alm_settings/set_bitbucket_binding.
// Since: 8.1.
func (s *AlmSettingsService) SetBitbucketBinding(ctx context.Context, opt *AlmSettingsSetBitbucketBindingOptions) (*http.Response, error) {
	err := s.ValidateSetBitbucketBindingOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "alm_settings/set_bitbucket_binding", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// SetBitbucketCloudBinding binds a Bitbucket Cloud setting to a project.
// If the project is already bound to a different Bitbucket Cloud setting, the binding will be updated to the new one.
// Requires the 'Administer' permission on the project.
//
// API endpoint: POST /api/alm_settings/set_bitbucketcloud_binding.
// Since: 8.7.
func (s *AlmSettingsService) SetBitbucketCloudBinding(ctx context.Context, opt *AlmSettingsSetBitbucketCloudBindingOptions) (*http.Response, error) {
	err := s.ValidateSetBitbucketCloudBindingOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "alm_settings/set_bitbucketcloud_binding", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// SetGithubBinding binds a GitHub instance to a project.
// If the project is already bound to a previous GitHub instance, the binding will be updated to the new one.
// Requires the 'Administer' permission on the project.
//
// API endpoint: POST /api/alm_settings/set_github_binding.
// Since: 8.1.
func (s *AlmSettingsService) SetGithubBinding(ctx context.Context, opt *AlmSettingsSetGithubBindingOptions) (*http.Response, error) {
	err := s.ValidateSetGithubBindingOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "alm_settings/set_github_binding", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// SetGitlabBinding binds a GitLab instance to a project.
// If the project is already bound to a previous GitLab instance, the binding will be updated to the new one.
// Requires the 'Administer' permission on the project.
//
// API endpoint: POST /api/alm_settings/set_gitlab_binding.
// Since: 8.1.
func (s *AlmSettingsService) SetGitlabBinding(ctx context.Context, opt *AlmSettingsSetGitlabBindingOptions) (*http.Response, error) {
	err := s.ValidateSetGitlabBindingOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "alm_settings/set_gitlab_binding", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// List lists DevOps Platform settings available for a given project, sorted by DevOps Platform key.
// Requires the 'Administer project' permission if the 'project' parameter is provided,
// requires the 'Create Projects' permission otherwise.
//
// API endpoint: GET /api/alm_settings/list.
// Since: 8.1.
func (s *AlmSettingsService) List(ctx context.Context, opt *AlmSettingsListOptions) (*AlmSettingsList, *http.Response, error) {
	err := s.ValidateListOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "alm_settings/list", opt)
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
func (s *AlmSettingsService) ListDefinitions(ctx context.Context) (*AlmSettingsListDefinitions, *http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "alm_settings/list_definitions", nil)
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
func (s *AlmSettingsService) UpdateAzure(ctx context.Context, opt *AlmSettingsUpdateAzureOptions) (*http.Response, error) {
	err := s.ValidateUpdateAzureOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "alm_settings/update_azure", opt)
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
func (s *AlmSettingsService) UpdateBitbucket(ctx context.Context, opt *AlmSettingsUpdateBitbucketOptions) (*http.Response, error) {
	err := s.ValidateUpdateBitbucketOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "alm_settings/update_bitbucket", opt)
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
func (s *AlmSettingsService) UpdateBitbucketCloud(ctx context.Context, opt *AlmSettingsUpdateBitbucketCloudOptions) (*http.Response, error) {
	err := s.ValidateUpdateBitbucketCloudOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "alm_settings/update_bitbucketcloud", opt)
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
func (s *AlmSettingsService) UpdateGithub(ctx context.Context, opt *AlmSettingsUpdateGithubOptions) (*http.Response, error) {
	err := s.ValidateUpdateGithubOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "alm_settings/update_github", opt)
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
func (s *AlmSettingsService) UpdateGitlab(ctx context.Context, opt *AlmSettingsUpdateGitlabOptions) (*http.Response, error) {
	err := s.ValidateUpdateGitlabOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "alm_settings/update_gitlab", opt)
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
func (s *AlmSettingsService) Validate(ctx context.Context, opt *AlmSettingsValidateOptions) (*AlmSettingsValidation, *http.Response, error) {
	err := s.ValidateValidateOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "alm_settings/validate", opt)
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
