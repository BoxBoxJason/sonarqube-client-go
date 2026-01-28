package sonargo

import (
	"net/http"
)

const (
	// MaxWebhookNameLength is the maximum length for a webhook name.
	MaxWebhookNameLength = 100
	// MaxWebhookProjectLength is the maximum length for a webhook project key.
	MaxWebhookProjectLength = 400
	// MaxWebhookSecretLength is the maximum length for a webhook secret.
	MaxWebhookSecretLength = 200
	// MinWebhookSecretLength is the minimum length for a webhook secret.
	// Since: 10.6.
	MinWebhookSecretLength = 16
	// MaxWebhookURLLength is the maximum length for a webhook URL.
	MaxWebhookURLLength = 512
	// MaxWebhookKeyLength is the maximum length for a webhook key.
	MaxWebhookKeyLength = 40
)

// WebhooksService handles communication with the Webhooks related methods of the SonarQube API.
// Webhooks allow to notify external services when a project analysis is done.
//
// Since: 6.2.
type WebhooksService struct {
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// Webhook represents a webhook configuration.
type Webhook struct {
	// Key is the unique identifier of the webhook.
	Key string `json:"key,omitempty"`
	// Name is the display name of the webhook.
	Name string `json:"name,omitempty"`
	// URL is the target endpoint for the webhook.
	URL string `json:"url,omitempty"`
	// HasSecret indicates if a secret is configured for HMAC signature.
	HasSecret bool `json:"hasSecret,omitempty"`
}

// WebhooksCreate represents the response from creating a webhook.
type WebhooksCreate struct {
	// Webhook contains the created webhook details.
	Webhook Webhook `json:"webhook,omitzero"`
}

// WebhooksDeliveries represents the response from listing webhook deliveries.
type WebhooksDeliveries struct {
	// Deliveries is the list of webhook deliveries.
	Deliveries []WebhookDelivery `json:"deliveries,omitempty"`
	// Paging contains pagination information.
	Paging Paging `json:"paging,omitzero"`
}

// WebhookDelivery represents a webhook delivery attempt.
type WebhookDelivery struct {
	// At is the timestamp of the delivery attempt.
	At string `json:"at,omitempty"`
	// CeTaskID is the Compute Engine task ID.
	//
	// Deprecated: Since 10.7.
	CeTaskID string `json:"ceTaskId,omitempty"`
	// ComponentKey is the project key.
	ComponentKey string `json:"componentKey,omitempty"`
	// ID is the unique identifier of the delivery.
	ID string `json:"id,omitempty"`
	// Name is the name of the webhook.
	Name string `json:"name,omitempty"`
	// Payload is the JSON payload sent (only in single delivery response).
	Payload string `json:"payload,omitempty"`
	// URL is the target URL.
	URL string `json:"url,omitempty"`
	// DurationMs is the duration of the request in milliseconds.
	DurationMs int64 `json:"durationMs,omitempty"`
	// HTTPStatus is the HTTP response status code.
	HTTPStatus int64 `json:"httpStatus,omitempty"`
	// Success indicates if the delivery was successful.
	Success bool `json:"success,omitempty"`
}

// WebhooksDelivery represents the response from getting a single delivery.
type WebhooksDelivery struct {
	// Delivery contains the detailed delivery information.
	Delivery WebhookDelivery `json:"delivery,omitzero"`
}

// WebhooksList represents the response from listing webhooks.
type WebhooksList struct {
	// Webhooks is the list of configured webhooks.
	Webhooks []Webhook `json:"webhooks,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// WebhooksCreateOption represents options for creating a webhook.
type WebhooksCreateOption struct {
	// Name is the display name of the webhook (required).
	// Maximum length: 100 characters.
	Name string `url:"name,omitempty"`
	// Project is the key of the project that will own the webhook (optional).
	// If not provided, the webhook will be global.
	// Maximum length: 400 characters.
	Project string `url:"project,omitempty"`
	// Secret is the HMAC secret for signing payloads (optional).
	// If provided, the secret will be used to generate the HMAC hex digest
	// in the 'X-Sonar-Webhook-HMAC-SHA256' header.
	// Minimum length: 16 characters (since 10.6).
	// Maximum length: 200 characters.
	Secret string `url:"secret,omitempty"`
	// URL is the target endpoint for the webhook (required).
	// Maximum length: 512 characters.
	URL string `url:"url,omitempty"`
}

// WebhooksDeleteOption represents options for deleting a webhook.
type WebhooksDeleteOption struct {
	// Webhook is the key of the webhook to delete (required).
	// Maximum length: 40 characters.
	Webhook string `url:"webhook,omitempty"`
}

// WebhooksDeliveriesOption represents options for listing webhook deliveries.
//
//nolint:govet // Embedded PaginationArgs makes optimal alignment impractical
type WebhooksDeliveriesOption struct {
	PaginationArgs

	// CeTaskID filters deliveries by Compute Engine task ID.
	//
	// Deprecated: Since 10.7.
	CeTaskID string `url:"ceTaskId,omitempty"`
	// ComponentKey filters deliveries by project key.
	//
	// Deprecated: Since 10.7.
	ComponentKey string `url:"componentKey,omitempty"`
	// Webhook filters deliveries by webhook key.
	Webhook string `url:"webhook,omitempty"`
}

// WebhooksDeliveryOption represents options for getting a single delivery.
type WebhooksDeliveryOption struct {
	// DeliveryID is the unique identifier of the delivery (required).
	DeliveryID string `url:"deliveryId,omitempty"`
}

// WebhooksListOption represents options for listing webhooks.
type WebhooksListOption struct {
	// Project filters webhooks by project key (optional).
	// If not provided, returns global webhooks.
	Project string `url:"project,omitempty"`
}

// WebhooksUpdateOption represents options for updating a webhook.
type WebhooksUpdateOption struct {
	// Name is the new name for the webhook (required).
	// Maximum length: 100 characters.
	Name string `url:"name,omitempty"`
	// Secret is the new HMAC secret (optional).
	// If blank, any existing secret will be removed.
	// If not set, the secret will remain unchanged.
	// Maximum length: 200 characters.
	Secret string `url:"secret,omitempty"`
	// URL is the new target endpoint (required).
	// Maximum length: 512 characters.
	URL string `url:"url,omitempty"`
	// Webhook is the key of the webhook to update (required).
	// Maximum length: 40 characters.
	Webhook string `url:"webhook,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Methods
// -----------------------------------------------------------------------------

// validateWebhookSecret validates the secret if provided.
func validateWebhookSecret(secret string) error {
	if secret == "" {
		return nil
	}

	err := ValidateMinLength(secret, MinWebhookSecretLength, "Secret")
	if err != nil {
		return err
	}

	return ValidateMaxLength(secret, MaxWebhookSecretLength, "Secret")
}

// ValidateCreateOpt validates the options for Create.
func (s *WebhooksService) ValidateCreateOpt(opt *WebhooksCreateOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Name, MaxWebhookNameLength, "Name")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.URL, "URL")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.URL, MaxWebhookURLLength, "URL")
	if err != nil {
		return err
	}

	if opt.Project != "" {
		err = ValidateMaxLength(opt.Project, MaxWebhookProjectLength, "Project")
		if err != nil {
			return err
		}
	}

	return validateWebhookSecret(opt.Secret)
}

// ValidateDeleteOpt validates the options for Delete.
func (s *WebhooksService) ValidateDeleteOpt(opt *WebhooksDeleteOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Webhook, "Webhook")
	if err != nil {
		return err
	}

	return ValidateMaxLength(opt.Webhook, MaxWebhookKeyLength, "Webhook")
}

// ValidateDeliveriesOpt validates the options for Deliveries.
func (s *WebhooksService) ValidateDeliveriesOpt(opt *WebhooksDeliveriesOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return opt.Validate()
}

// ValidateDeliveryOpt validates the options for Delivery.
func (s *WebhooksService) ValidateDeliveryOpt(opt *WebhooksDeliveryOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.DeliveryID, "DeliveryID")
}

// ValidateListOpt validates the options for List.
func (s *WebhooksService) ValidateListOpt(opt *WebhooksListOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return nil
}

// ValidateUpdateOpt validates the options for Update.
func (s *WebhooksService) ValidateUpdateOpt(opt *WebhooksUpdateOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Name, "Name")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Name, MaxWebhookNameLength, "Name")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.URL, "URL")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.URL, MaxWebhookURLLength, "URL")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Webhook, "Webhook")
	if err != nil {
		return err
	}

	err = ValidateMaxLength(opt.Webhook, MaxWebhookKeyLength, "Webhook")
	if err != nil {
		return err
	}

	if opt.Secret != "" {
		err = ValidateMaxLength(opt.Secret, MaxWebhookSecretLength, "Secret")
		if err != nil {
			return err
		}
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Create creates a webhook.
// Requires 'Administer' permission on the specified project, or global 'Administer' permission.
//
// Since: 7.1.
func (s *WebhooksService) Create(opt *WebhooksCreateOption) (*WebhooksCreate, *http.Response, error) {
	err := s.ValidateCreateOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "webhooks/create", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(WebhooksCreate)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Delete deletes a webhook.
// Requires 'Administer' permission on the specified project, or global 'Administer' permission.
//
// Since: 7.1.
func (s *WebhooksService) Delete(opt *WebhooksDeleteOption) (*http.Response, error) {
	err := s.ValidateDeleteOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "webhooks/delete", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Deliveries gets the recent deliveries for a specified project or Compute Engine task.
// Requires 'Administer' permission on the related project.
// Note that additional information is returned by api/webhooks/delivery.
//
// Since: 6.2.
func (s *WebhooksService) Deliveries(opt *WebhooksDeliveriesOption) (*WebhooksDeliveries, *http.Response, error) {
	err := s.ValidateDeliveriesOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "webhooks/deliveries", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(WebhooksDeliveries)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Delivery gets a webhook delivery by its ID.
// Requires 'Administer System' permission.
//
// Since: 6.2.
func (s *WebhooksService) Delivery(opt *WebhooksDeliveryOption) (*WebhooksDelivery, *http.Response, error) {
	err := s.ValidateDeliveryOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "webhooks/delivery", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(WebhooksDelivery)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// List searches for global webhooks or project webhooks.
// Webhooks are ordered by name.
// Requires 'Administer' permission on the specified project, or global 'Administer' permission.
//
// Since: 7.1.
func (s *WebhooksService) List(opt *WebhooksListOption) (*WebhooksList, *http.Response, error) {
	err := s.ValidateListOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "webhooks/list", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(WebhooksList)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Update updates a webhook.
// Requires 'Administer' permission on the specified project, or global 'Administer' permission.
//
// Since: 7.1.
func (s *WebhooksService) Update(opt *WebhooksUpdateOption) (*http.Response, error) {
	err := s.ValidateUpdateOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "webhooks/update", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
