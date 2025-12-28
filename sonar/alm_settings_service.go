// Manage DevOps Platform Settings
package sonargo

import "net/http"

type AlmSettingsService struct {
	client *Client
}

type AlmSettingsCountBindingObject struct {
	Key      string `json:"key,omitempty"`
	Projects int64  `json:"projects,omitempty"`
}

type AlmSettingsGetBindingObject struct {
	Alm                   string `json:"alm,omitempty"`
	Key                   string `json:"key,omitempty"`
	Monorepo              bool   `json:"monorepo,omitempty"`
	Repository            string `json:"repository,omitempty"`
	RepositoryURL         string `json:"repositoryUrl,omitempty"`
	SummaryCommentEnabled bool   `json:"summaryCommentEnabled,omitempty"`
	URL                   string `json:"url,omitempty"`
}

type AlmSettingsListObject struct {
	AlmSettings []AlmSettingsListObject_sub1 `json:"almSettings,omitempty"`
}

type AlmSettingsListObject_sub1 struct {
	Alm string `json:"alm,omitempty"`
	Key string `json:"key,omitempty"`
	URL string `json:"url,omitempty"`
}

type AlmSettingsListDefinitionsObject struct {
	Azure          []interface{}                           `json:"azure,omitempty"`
	Bitbucket      []interface{}                           `json:"bitbucket,omitempty"`
	Bitbucketcloud []interface{}                           `json:"bitbucketcloud,omitempty"`
	Github         []AlmSettingsListDefinitionsObject_sub1 `json:"github,omitempty"`
	Gitlab         []interface{}                           `json:"gitlab,omitempty"`
}

type AlmSettingsListDefinitionsObject_sub1 struct {
	AppID    string `json:"appId,omitempty"`
	ClientID string `json:"clientId,omitempty"`
	Key      string `json:"key,omitempty"`
	URL      string `json:"url,omitempty"`
}

type AlmSettingsValidateObject struct {
	Errors []AlmSettingsValidateObject_sub1 `json:"errors,omitempty"`
}

type AlmSettingsValidateObject_sub1 struct {
	Msg string `json:"msg,omitempty"`
}

type AlmSettingsCountBindingOption struct {
	AlmSetting string `url:"almSetting,omitempty"` // Description:"DevOps Platform setting key",ExampleValue:""
}

// CountBinding Count number of project bound to an DevOps Platform setting.<br/>Requires the 'Administer System' permission
func (s *AlmSettingsService) CountBinding(opt *AlmSettingsCountBindingOption) (v *AlmSettingsCountBindingObject, resp *http.Response, err error) {
	err = s.ValidateCountBindingOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "alm_settings/count_binding", opt)
	if err != nil {
		return
	}
	v = new(AlmSettingsCountBindingObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type AlmSettingsCreateAzureOption struct {
	Key                 string `url:"key,omitempty"`                 // Description:"Unique key of the Azure Devops instance setting",ExampleValue:""
	PersonalAccessToken string `url:"personalAccessToken,omitempty"` // Description:"Azure Devops personal access token",ExampleValue:""
	Url                 string `url:"url,omitempty"`                 // Description:"Azure API URL",ExampleValue:""
}

// CreateAzure Create Azure instance Setting. <br/>Requires the 'Administer System' permission
func (s *AlmSettingsService) CreateAzure(opt *AlmSettingsCreateAzureOption) (resp *http.Response, err error) {
	err = s.ValidateCreateAzureOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "alm_settings/create_azure", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type AlmSettingsCreateBitbucketOption struct {
	Key                 string `url:"key,omitempty"`                 // Description:"Unique key of the Bitbucket instance setting",ExampleValue:""
	PersonalAccessToken string `url:"personalAccessToken,omitempty"` // Description:"Bitbucket personal access token",ExampleValue:""
	Url                 string `url:"url,omitempty"`                 // Description:"BitBucket server API URL",ExampleValue:""
}

// CreateBitbucket Create Bitbucket instance Setting. <br/>Requires the 'Administer System' permission
func (s *AlmSettingsService) CreateBitbucket(opt *AlmSettingsCreateBitbucketOption) (resp *http.Response, err error) {
	err = s.ValidateCreateBitbucketOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "alm_settings/create_bitbucket", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type AlmSettingsCreateBitbucketcloudOption struct {
	ClientId     string `url:"clientId,omitempty"`     // Description:"Bitbucket Cloud Client ID",ExampleValue:""
	ClientSecret string `url:"clientSecret,omitempty"` // Description:"Bitbucket Cloud Client Secret",ExampleValue:""
	Key          string `url:"key,omitempty"`          // Description:"Unique key of the Bitbucket Cloud setting",ExampleValue:""
	Workspace    string `url:"workspace,omitempty"`    // Description:"Bitbucket Cloud workspace ID",ExampleValue:""
}

// CreateBitbucketcloud Configure a new instance of Bitbucket Cloud. <br/>Requires the 'Administer System' permission
func (s *AlmSettingsService) CreateBitbucketcloud(opt *AlmSettingsCreateBitbucketcloudOption) (resp *http.Response, err error) {
	err = s.ValidateCreateBitbucketcloudOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "alm_settings/create_bitbucketcloud", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type AlmSettingsCreateGithubOption struct {
	AppId         string `url:"appId,omitempty"`         // Description:"GitHub App ID",ExampleValue:""
	ClientId      string `url:"clientId,omitempty"`      // Description:"GitHub App Client ID",ExampleValue:""
	ClientSecret  string `url:"clientSecret,omitempty"`  // Description:"GitHub App Client Secret",ExampleValue:""
	Key           string `url:"key,omitempty"`           // Description:"Unique key of the GitHub instance setting",ExampleValue:""
	PrivateKey    string `url:"privateKey,omitempty"`    // Description:"GitHub App private key",ExampleValue:""
	Url           string `url:"url,omitempty"`           // Description:"GitHub API URL",ExampleValue:""
	WebhookSecret string `url:"webhookSecret,omitempty"` // Description:"GitHub App Webhook Secret",ExampleValue:""
}

// CreateGithub Create GitHub instance Setting. <br/>Requires the 'Administer System' permission
func (s *AlmSettingsService) CreateGithub(opt *AlmSettingsCreateGithubOption) (resp *http.Response, err error) {
	err = s.ValidateCreateGithubOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "alm_settings/create_github", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type AlmSettingsCreateGitlabOption struct {
	Key                 string `url:"key,omitempty"`                 // Description:"Unique key of the GitLab instance setting",ExampleValue:""
	PersonalAccessToken string `url:"personalAccessToken,omitempty"` // Description:"GitLab personal access token",ExampleValue:""
	Url                 string `url:"url,omitempty"`                 // Description:"GitLab API URL",ExampleValue:""
}

// CreateGitlab Create GitLab instance Setting. <br/>Requires the 'Administer System' permission
func (s *AlmSettingsService) CreateGitlab(opt *AlmSettingsCreateGitlabOption) (resp *http.Response, err error) {
	err = s.ValidateCreateGitlabOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "alm_settings/create_gitlab", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type AlmSettingsDeleteOption struct {
	Key string `url:"key,omitempty"` // Description:"DevOps Platform Setting key",ExampleValue:""
}

// Delete Delete an DevOps Platform Setting.<br/>Requires the 'Administer System' permission
func (s *AlmSettingsService) Delete(opt *AlmSettingsDeleteOption) (resp *http.Response, err error) {
	err = s.ValidateDeleteOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "alm_settings/delete", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type AlmSettingsGetBindingOption struct {
	Project string `url:"project,omitempty"` // Description:"Project key",ExampleValue:""
}

// GetBinding Get DevOps Platform binding of a given project.<br/>Requires the 'Browse' permission on the project
func (s *AlmSettingsService) GetBinding(opt *AlmSettingsGetBindingOption) (v *AlmSettingsGetBindingObject, resp *http.Response, err error) {
	err = s.ValidateGetBindingOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "alm_settings/get_binding", opt)
	if err != nil {
		return
	}
	v = new(AlmSettingsGetBindingObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type AlmSettingsListOption struct {
	Project string `url:"project,omitempty"` // Description:"Project key",ExampleValue:""
}

// List List DevOps Platform setting available for a given project, sorted by DevOps Platform key<br/>Requires the 'Administer project' permission if the 'project' parameter is provided, requires the 'Create Projects' permission otherwise.
func (s *AlmSettingsService) List(opt *AlmSettingsListOption) (v *AlmSettingsListObject, resp *http.Response, err error) {
	err = s.ValidateListOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "alm_settings/list", opt)
	if err != nil {
		return
	}
	v = new(AlmSettingsListObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

// ListDefinitions List DevOps Platform Settings, sorted by created date.<br/>Requires the 'Administer System' permission
func (s *AlmSettingsService) ListDefinitions() (v *AlmSettingsListDefinitionsObject, resp *http.Response, err error) {
	req, err := s.client.NewRequest("GET", "alm_settings/list_definitions", nil)
	if err != nil {
		return
	}
	v = new(AlmSettingsListDefinitionsObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type AlmSettingsUpdateAzureOption struct {
	Key                 string `url:"key,omitempty"`                 // Description:"Unique key of the Azure instance setting",ExampleValue:""
	NewKey              string `url:"newKey,omitempty"`              // Description:"Optional new value for an unique key of the Azure Devops instance setting",ExampleValue:""
	PersonalAccessToken string `url:"personalAccessToken,omitempty"` // Description:"Azure Devops personal access token",ExampleValue:""
	Url                 string `url:"url,omitempty"`                 // Description:"Azure API URL",ExampleValue:""
}

// UpdateAzure Update Azure instance Setting. <br/>Requires the 'Administer System' permission
func (s *AlmSettingsService) UpdateAzure(opt *AlmSettingsUpdateAzureOption) (resp *http.Response, err error) {
	err = s.ValidateUpdateAzureOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "alm_settings/update_azure", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type AlmSettingsUpdateBitbucketOption struct {
	Key                 string `url:"key,omitempty"`                 // Description:"Unique key of the Bitbucket instance setting",ExampleValue:""
	NewKey              string `url:"newKey,omitempty"`              // Description:"Optional new value for an unique key of the Bitbucket instance setting",ExampleValue:""
	PersonalAccessToken string `url:"personalAccessToken,omitempty"` // Description:"Bitbucket personal access token",ExampleValue:""
	Url                 string `url:"url,omitempty"`                 // Description:"Bitbucket API URL",ExampleValue:""
}

// UpdateBitbucket Update Bitbucket instance Setting. <br/>Requires the 'Administer System' permission
func (s *AlmSettingsService) UpdateBitbucket(opt *AlmSettingsUpdateBitbucketOption) (resp *http.Response, err error) {
	err = s.ValidateUpdateBitbucketOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "alm_settings/update_bitbucket", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type AlmSettingsUpdateBitbucketcloudOption struct {
	ClientId     string `url:"clientId,omitempty"`     // Description:"Bitbucket Cloud Client ID",ExampleValue:""
	ClientSecret string `url:"clientSecret,omitempty"` // Description:"Optional new value for the Bitbucket Cloud client secret",ExampleValue:""
	Key          string `url:"key,omitempty"`          // Description:"Unique key of the Bitbucket Cloud setting",ExampleValue:""
	NewKey       string `url:"newKey,omitempty"`       // Description:"Optional new value for an unique key of the Bitbucket Cloud setting",ExampleValue:""
	Workspace    string `url:"workspace,omitempty"`    // Description:"Bitbucket Cloud workspace ID",ExampleValue:""
}

// UpdateBitbucketcloud Update Bitbucket Cloud Setting. <br/>Requires the 'Administer System' permission
func (s *AlmSettingsService) UpdateBitbucketcloud(opt *AlmSettingsUpdateBitbucketcloudOption) (resp *http.Response, err error) {
	err = s.ValidateUpdateBitbucketcloudOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "alm_settings/update_bitbucketcloud", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type AlmSettingsUpdateGithubOption struct {
	AppId         string `url:"appId,omitempty"`         // Description:"GitHub API ID",ExampleValue:""
	ClientId      string `url:"clientId,omitempty"`      // Description:"GitHub App Client ID",ExampleValue:""
	ClientSecret  string `url:"clientSecret,omitempty"`  // Description:"GitHub App Client Secret",ExampleValue:""
	Key           string `url:"key,omitempty"`           // Description:"Unique key of the GitHub instance setting",ExampleValue:""
	NewKey        string `url:"newKey,omitempty"`        // Description:"Optional new value for an unique key of the GitHub instance setting",ExampleValue:""
	PrivateKey    string `url:"privateKey,omitempty"`    // Description:"GitHub App private key",ExampleValue:""
	Url           string `url:"url,omitempty"`           // Description:"GitHub API URL",ExampleValue:""
	WebhookSecret string `url:"webhookSecret,omitempty"` // Description:"GitHub App Webhook Secret",ExampleValue:""
}

// UpdateGithub Update GitHub instance Setting. <br/>Requires the 'Administer System' permission
func (s *AlmSettingsService) UpdateGithub(opt *AlmSettingsUpdateGithubOption) (resp *http.Response, err error) {
	err = s.ValidateUpdateGithubOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "alm_settings/update_github", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type AlmSettingsUpdateGitlabOption struct {
	Key                 string `url:"key,omitempty"`                 // Description:"Unique key of the GitLab instance setting",ExampleValue:""
	NewKey              string `url:"newKey,omitempty"`              // Description:"Optional new value for an unique key of the GitLab instance setting",ExampleValue:""
	PersonalAccessToken string `url:"personalAccessToken,omitempty"` // Description:"GitLab personal access token",ExampleValue:""
	Url                 string `url:"url,omitempty"`                 // Description:"GitLab API URL",ExampleValue:""
}

// UpdateGitlab Update GitLab instance Setting. <br/>Requires the 'Administer System' permission
func (s *AlmSettingsService) UpdateGitlab(opt *AlmSettingsUpdateGitlabOption) (resp *http.Response, err error) {
	err = s.ValidateUpdateGitlabOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "alm_settings/update_gitlab", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type AlmSettingsValidateOption struct {
	Key string `url:"key,omitempty"` // Description:"Unique key of the DevOps Platform settings",ExampleValue:""
}

// Validate Validate an DevOps Platform Setting by checking connectivity and permissions<br/>Requires the 'Administer System' permission
func (s *AlmSettingsService) Validate(opt *AlmSettingsValidateOption) (v *AlmSettingsValidateObject, resp *http.Response, err error) {
	err = s.ValidateValidateOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "alm_settings/validate", opt)
	if err != nil {
		return
	}
	v = new(AlmSettingsValidateObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}
