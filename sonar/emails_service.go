package sonargo

import "net/http"

// EmailsService handles communication with the email related methods
// of the SonarQube API.
// This service manages email operations, including sending test emails.
type EmailsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// EmailsSendOption contains parameters for the Send method.
type EmailsSendOption struct {
	// Message is the content of the email.
	// This field is required.
	Message string `url:"message"`
	// Subject is the subject of the email.
	// Optional, defaults to empty.
	Subject string `url:"subject,omitempty"`
	// To is the recipient email address.
	// This field is required.
	To string `url:"to"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateSendOpt validates the options for the Send method.
func (s *EmailsService) ValidateSendOpt(opt *EmailsSendOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Message, "Message")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.To, "To")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Send tests email configuration by sending an email.
// Requires 'Administer System' permission.
//
// API endpoint: POST /api/emails/send.
// Since: 6.1 (internal).
func (s *EmailsService) Send(opt *EmailsSendOption) (*http.Response, error) {
	err := s.ValidateSendOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "emails/send", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
