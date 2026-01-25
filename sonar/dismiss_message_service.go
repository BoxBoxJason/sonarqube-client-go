package sonargo

import "net/http"

// DismissMessageService handles communication with the message dismissal related methods
// of the SonarQube API.
// This service manages message dismissal for users.
type DismissMessageService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// DismissMessageCheck represents the response from checking if a message has been dismissed.
type DismissMessageCheck struct {
	// Dismissed indicates whether the message has been dismissed.
	Dismissed bool `json:"dismissed,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// DismissMessageCheckOption contains parameters for the Check method.
type DismissMessageCheckOption struct {
	// MessageType is the type of the message to check.
	// This field is required.
	MessageType string `url:"messageType"`
	// ProjectKey is the project key.
	// This field is required.
	ProjectKey string `url:"projectKey"`
}

// DismissMessageDismissOption contains parameters for the Dismiss method.
type DismissMessageDismissOption struct {
	// MessageType is the type of the message to dismiss.
	// This field is required.
	MessageType string `url:"messageType"`
	// ProjectKey is the project key.
	// This field is required.
	ProjectKey string `url:"projectKey"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateCheckOpt validates the options for the Check method.
func (s *DismissMessageService) ValidateCheckOpt(opt *DismissMessageCheckOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.MessageType, "MessageType")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	return nil
}

// ValidateDismissOpt validates the options for the Dismiss method.
func (s *DismissMessageService) ValidateDismissOpt(opt *DismissMessageDismissOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.MessageType, "MessageType")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Check checks if a message has been dismissed.
//
// API endpoint: GET /api/dismiss_message/check.
// Since: 10.2.
func (s *DismissMessageService) Check(opt *DismissMessageCheckOption) (*DismissMessageCheck, *http.Response, error) {
	err := s.ValidateCheckOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "dismiss_message/check", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(DismissMessageCheck)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Dismiss dismisses a message for the current user.
//
// API endpoint: POST /api/dismiss_message/dismiss.
// Since: 10.2.
func (s *DismissMessageService) Dismiss(opt *DismissMessageDismissOption) (*http.Response, error) {
	err := s.ValidateDismissOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "dismiss_message/dismiss", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
