package sonargo

import "net/http"

// AuthenticationService handles communication with the authentication related methods
// of the SonarQube API.
// This service handles user authentication.
type AuthenticationService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// AuthenticationValidation represents the response from validating credentials.
type AuthenticationValidation struct {
	// Valid indicates whether the credentials are valid.
	Valid bool `json:"valid,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// AuthenticationLoginOption contains parameters for the Login method.
type AuthenticationLoginOption struct {
	// Login is the login/username of the user.
	// This field is required.
	Login string `url:"login"`
	// Password is the password of the user.
	// This field is required.
	Password string `url:"password"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateLoginOpt validates the options for the Login method.
func (s *AuthenticationService) ValidateLoginOpt(opt *AuthenticationLoginOption) error {
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

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Login authenticates a user.
//
// API endpoint: POST /api/authentication/login.
func (s *AuthenticationService) Login(opt *AuthenticationLoginOption) (*http.Response, error) {
	err := s.ValidateLoginOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "authentication/login", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Logout logs out the current user.
//
// API endpoint: POST /api/authentication/logout.
func (s *AuthenticationService) Logout() (*http.Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "authentication/logout", nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Validate checks if the current credentials are valid.
//
// API endpoint: GET /api/authentication/validate.
func (s *AuthenticationService) Validate() (*AuthenticationValidation, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "authentication/validate", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(AuthenticationValidation)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
