// List, create, and delete a user's access tokens.
package sonargo

import "net/http"

type UserTokensService struct {
	client *Client
}

type UserTokensGenerateObject struct {
	CreatedAt      string `json:"createdAt,omitempty"`
	ExpirationDate string `json:"expirationDate,omitempty"`
	Login          string `json:"login,omitempty"`
	Name           string `json:"name,omitempty"`
	Token          string `json:"token,omitempty"`
	Type           string `json:"type,omitempty"`
}

type UserTokensSearchObject struct {
	Login      string                        `json:"login,omitempty"`
	UserTokens []UserTokensSearchObject_sub2 `json:"userTokens,omitempty"`
}

type UserTokensSearchObject_sub2 struct {
	CreatedAt      string                      `json:"createdAt,omitempty"`
	ExpirationDate string                      `json:"expirationDate,omitempty"`
	IsExpired      bool                        `json:"isExpired,omitempty"`
	Name           string                      `json:"name,omitempty"`
	Project        UserTokensSearchObject_sub1 `json:"project,omitempty"`
	Type           string                      `json:"type,omitempty"`
}

type UserTokensSearchObject_sub1 struct {
	Key  string `json:"key,omitempty"`
	Name string `json:"name,omitempty"`
}

type UserTokensGenerateOption struct {
	ExpirationDate string `url:"expirationDate,omitempty"` // Description:"The expiration date of the token being generated, in ISO 8601 format (YYYY-MM-DD). If not set, default to no expiration.",ExampleValue:""
	Login          string `url:"login,omitempty"`          // Description:"User login. If not set, the token is generated for the authenticated user.",ExampleValue:"g.hopper"
	Name           string `url:"name,omitempty"`           // Description:"Token name",ExampleValue:"Project scan on Travis"
	ProjectKey     string `url:"projectKey,omitempty"`     // Description:"The key of the only project that can be analyzed by the PROJECT_ANALYSIS_TOKEN being generated.",ExampleValue:""
	Type           string `url:"type,omitempty"`           // Description:"Token Type. If this parameters is set to PROJECT_ANALYSIS_TOKEN, it is necessary to provide the projectKey parameter too.",ExampleValue:""
}

// Generate Generate a user access token. <br />Please keep your tokens secret. They enable to authenticate and analyze projects.<br />It requires administration permissions to specify a 'login' and generate a token for another user. Otherwise, a token is generated for the current user.
func (s *UserTokensService) Generate(opt *UserTokensGenerateOption) (v *UserTokensGenerateObject, resp *http.Response, err error) {
	err = s.ValidateGenerateOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "user_tokens/generate", opt)
	if err != nil {
		return
	}
	v = new(UserTokensGenerateObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type UserTokensRevokeOption struct {
	Login string `url:"login,omitempty"` // Description:"User login",ExampleValue:"g.hopper"
	Name  string `url:"name,omitempty"`  // Description:"Token name",ExampleValue:"Project scan on Travis"
}

// Revoke Revoke a user access token. <br/>It requires administration permissions to specify a 'login' and revoke a token for another user. Otherwise, the token for the current user is revoked.
func (s *UserTokensService) Revoke(opt *UserTokensRevokeOption) (resp *http.Response, err error) {
	err = s.ValidateRevokeOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "user_tokens/revoke", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type UserTokensSearchOption struct {
	Login string `url:"login,omitempty"` // Description:"User login",ExampleValue:"g.hopper"
}

// Search List the access tokens of a user.<br>The login must exist and active.<br>Field 'lastConnectionDate' is only updated every hour, so it may not be accurate, for instance when a user is using a token many times in less than one hour.<br> It requires administration permissions to specify a 'login' and list the tokens of another user. Otherwise, tokens for the current user are listed. <br> Authentication is required for this API endpoint
func (s *UserTokensService) Search(opt *UserTokensSearchOption) (v *UserTokensSearchObject, resp *http.Response, err error) {
	err = s.ValidateSearchOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "user_tokens/search", opt)
	if err != nil {
		return
	}
	v = new(UserTokensSearchObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}
