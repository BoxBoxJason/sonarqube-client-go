// Handle authentication.
package sonargo

import "net/http"

type AuthenticationService struct {
	client *Client
}

type AuthenticationValidateObject struct {
	Valid bool `json:"valid,omitempty"`
}

type AuthenticationLoginOption struct {
	Login    string `url:"login,omitempty"`    // Description:"Login of the user",ExampleValue:""
	Password string `url:"password,omitempty"` // Description:"Password of the user",ExampleValue:""
}

// Login Authenticate a user.
func (s *AuthenticationService) Login(opt *AuthenticationLoginOption) (resp *http.Response, err error) {
	err = s.ValidateLoginOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "authentication/login", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

// Logout Logout a user.
func (s *AuthenticationService) Logout() (resp *http.Response, err error) {
	req, err := s.client.NewRequest("POST", "authentication/logout", nil)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

// Validate Check credentials.
func (s *AuthenticationService) Validate() (v *AuthenticationValidateObject, resp *http.Response, err error) {
	req, err := s.client.NewRequest("GET", "authentication/validate", nil)
	if err != nil {
		return
	}
	v = new(AuthenticationValidateObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}
