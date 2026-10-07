package sonar

import (
	"context"
	"net/http"
)

//nolint:gochecknoglobals // Allowed-value sets
var (
	// allowedNotificationFilters is the set of allowed notification list filters.
	allowedNotificationFilters = map[string]struct{}{
		"user":              {},
		"groupSubscription": {},
		"all":               {},
	}
	// allowedNotificationGroupTypes is the set of allowed notification types for group subscriptions.
	allowedNotificationGroupTypes = map[string]struct{}{
		"security-alert-raised": {},
	}
)

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

// NotificationsGroupSubscription represents a group subscribed to a notification type.
type NotificationsGroupSubscription struct {
	// GroupUuid is the UUID of the subscribed group.
	GroupUuid string `json:"groupUuid,omitempty"`
	// GroupName is the name of the subscribed group.
	GroupName string `json:"groupName,omitempty"`
	// NotificationType is the notification type (e.g. security-alert-raised).
	NotificationType string `json:"notificationType,omitempty"`
	// ChannelKey is the key of the notification channel (e.g. EmailNotificationChannel).
	ChannelKey string `json:"channelKey,omitempty"`
}

// NotificationsListGroups represents the response from listing all group notification subscriptions.
type NotificationsListGroups struct {
	// Subscriptions is the list of group notification subscriptions.
	Subscriptions []NotificationsGroupSubscription `json:"subscriptions,omitempty"`
}

// NotificationsUserGroupSubscription represents a notification subscription of one of the current user's groups.
type NotificationsUserGroupSubscription struct {
	// GroupName is the name of the subscribed group.
	GroupName string `json:"groupName,omitempty"`
	// NotificationType is the notification type (e.g. security-alert-raised).
	NotificationType string `json:"notificationType,omitempty"`
}

// NotificationsListGroupSubscriptions represents the response from listing the current user's group subscriptions.
type NotificationsListGroupSubscriptions struct {
	// GroupSubscriptions is the list of notification subscriptions of the current user's groups.
	GroupSubscriptions []NotificationsUserGroupSubscription `json:"groupSubscriptions,omitempty"`
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

// NotificationsAddOptions contains parameters for the Add method.
type NotificationsAddOptions struct {
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

// NotificationsGroupOptions contains parameters for the AddGroup and RemoveGroup methods.
type NotificationsGroupOptions struct {
	// GroupUuid is the UUID of the group.
	// This field is required.
	GroupUuid string `url:"groupUuid"`
	// Type is the notification type.
	// This field is required.
	// Allowed values: security-alert-raised.
	Type string `url:"type"`
}

// NotificationsListOptions contains parameters for the List method.
type NotificationsListOptions struct {
	// Filter is the category of notification types to return.
	// Allowed values: user, groupSubscription, all. Default: all. Since 2026.5.
	Filter string `url:"filter,omitempty"`
	// Login is the user login. If not provided, the authenticated user is used.
	Login string `url:"login,omitempty"`
}

// NotificationsRemoveOptions contains parameters for the Remove method.
type NotificationsRemoveOptions struct {
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
func (s *NotificationsService) ValidateAddOpt(opt *NotificationsAddOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Type, "Type")
	if err != nil {
		return err
	}

	return nil
}

// ValidateGroupOpt validates the options for the AddGroup and RemoveGroup methods.
func (s *NotificationsService) ValidateGroupOpt(opt *NotificationsGroupOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.GroupUuid, "GroupUuid")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Type, "Type")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Type, allowedNotificationGroupTypes, "Type")
	if err != nil {
		return err
	}

	return nil
}

// ValidateListOpt validates the options for the List method.
func (s *NotificationsService) ValidateListOpt(opt *NotificationsListOptions) error {
	if opt == nil {
		return nil
	}

	return IsValueAuthorized(opt.Filter, allowedNotificationFilters, "Filter")
}

// ValidateRemoveOpt validates the options for the Remove method.
func (s *NotificationsService) ValidateRemoveOpt(opt *NotificationsRemoveOptions) error {
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
func (s *NotificationsService) Add(ctx context.Context, opt *NotificationsAddOptions) (*http.Response, error) {
	err := s.ValidateAddOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "notifications/add", opt)
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
func (s *NotificationsService) List(ctx context.Context, opt *NotificationsListOptions) (*NotificationsList, *http.Response, error) {
	err := s.ValidateListOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "notifications/list", opt)
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
func (s *NotificationsService) Remove(ctx context.Context, opt *NotificationsRemoveOptions) (*http.Response, error) {
	err := s.ValidateRemoveOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "notifications/remove", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// AddGroup subscribes a group to a notification type.
// Requires system administration permission.
//
// API endpoint: POST /api/notifications/add_group.
// Since: 2026.5.
// Internal: true.
func (s *NotificationsService) AddGroup(ctx context.Context, opt *NotificationsGroupOptions) (*http.Response, error) {
	err := s.ValidateGroupOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "notifications/add_group", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// ListGroupSubscriptions lists the notification group subscriptions for the current user's groups.
// Requires authentication.
//
// API endpoint: GET /api/notifications/list_group_subscriptions.
// Since: 2026.5.
// Internal: true.
func (s *NotificationsService) ListGroupSubscriptions(ctx context.Context) (*NotificationsListGroupSubscriptions, *http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "notifications/list_group_subscriptions", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(NotificationsListGroupSubscriptions)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ListGroups lists the group notification subscriptions.
// Requires system administration permission.
//
// API endpoint: GET /api/notifications/list_groups.
// Since: 2026.5.
// Internal: true.
func (s *NotificationsService) ListGroups(ctx context.Context) (*NotificationsListGroups, *http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "notifications/list_groups", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(NotificationsListGroups)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// RemoveGroup removes a group notification subscription.
// This is a no-op if the subscription is not found.
// Requires system administration permission.
//
// API endpoint: POST /api/notifications/remove_group.
// Since: 2026.5.
// Internal: true.
func (s *NotificationsService) RemoveGroup(ctx context.Context, opt *NotificationsGroupOptions) (*http.Response, error) {
	err := s.ValidateGroupOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "notifications/remove_group", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
