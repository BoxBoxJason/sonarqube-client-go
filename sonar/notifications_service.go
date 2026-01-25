package sonargo

import "net/http"

// NotificationsService handles communication with the notifications related methods
// of the SonarQube API.
// This service manages notifications for the authenticated user.
type NotificationsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// NotificationsList represents the response from listing notifications.
type NotificationsList struct {
	// Channels is the list of available notification channels.
	Channels []string `json:"channels,omitempty"`
	// GlobalTypes is the list of global notification types.
	GlobalTypes []string `json:"globalTypes,omitempty"`
	// Notifications is the list of configured notifications.
	Notifications []Notification `json:"notifications,omitempty"`
	// PerProjectTypes is the list of per-project notification types.
	PerProjectTypes []string `json:"perProjectTypes,omitempty"`
}

// Notification represents a configured notification.
type Notification struct {
	// Channel is the notification channel (e.g., email).
	Channel string `json:"channel,omitempty"`
	// Organization is the organization key (deprecated).
	Organization string `json:"organization,omitempty"`
	// Project is the project key.
	Project string `json:"project,omitempty"`
	// ProjectName is the project name.
	ProjectName string `json:"projectName,omitempty"`
	// Type is the notification type.
	Type string `json:"type,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// NotificationsAddOption contains parameters for the Add method.
type NotificationsAddOption struct {
	// Channel is the channel through which the notification is sent.
	// Default is email.
	Channel string `url:"channel,omitempty"`
	// Login is the user login. If not provided, the authenticated user is used.
	Login string `url:"login,omitempty"`
	// Project is the project key for per-project notifications.
	Project string `url:"project,omitempty"`
	// Type is the notification type.
	// This field is required.
	Type string `url:"type"`
}

// NotificationsListOption contains parameters for the List method.
type NotificationsListOption struct {
	// Login is the user login. If not provided, the authenticated user is used.
	Login string `url:"login,omitempty"`
}

// NotificationsRemoveOption contains parameters for the Remove method.
type NotificationsRemoveOption struct {
	// Channel is the channel through which the notification is sent.
	// Default is email.
	Channel string `url:"channel,omitempty"`
	// Login is the user login. If not provided, the authenticated user is used.
	Login string `url:"login,omitempty"`
	// Project is the project key for per-project notifications.
	Project string `url:"project,omitempty"`
	// Type is the notification type.
	// This field is required.
	Type string `url:"type"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateAddOpt validates the options for the Add method.
func (s *NotificationsService) ValidateAddOpt(opt *NotificationsAddOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Type, "Type")
	if err != nil {
		return err
	}

	return nil
}

// ValidateListOpt validates the options for the List method.
func (s *NotificationsService) ValidateListOpt(opt *NotificationsListOption) error {
	// No required fields
	return nil
}

// ValidateRemoveOpt validates the options for the Remove method.
func (s *NotificationsService) ValidateRemoveOpt(opt *NotificationsRemoveOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Type, "Type")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Add adds a notification for the authenticated user.
// Requires authentication if no login is provided.
// Requires system administration if a login is provided.
// If a project is provided, requires the 'Browse' permission on the specified project.
//
// API endpoint: POST /api/notifications/add.
// Since: 6.3.
func (s *NotificationsService) Add(opt *NotificationsAddOption) (*http.Response, error) {
	err := s.ValidateAddOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "notifications/add", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// List lists notifications of the authenticated user.
// Requires authentication if no login is provided.
// Requires system administration if a login is provided.
//
// API endpoint: GET /api/notifications/list.
// Since: 6.3.
func (s *NotificationsService) List(opt *NotificationsListOption) (*NotificationsList, *http.Response, error) {
	err := s.ValidateListOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "notifications/list", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(NotificationsList)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Remove removes a notification for the authenticated user.
// Requires authentication if no login is provided.
// Requires system administration if a login is provided.
//
// API endpoint: POST /api/notifications/remove.
// Since: 6.3.
func (s *NotificationsService) Remove(opt *NotificationsRemoveOption) (*http.Response, error) {
	err := s.ValidateRemoveOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "notifications/remove", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
