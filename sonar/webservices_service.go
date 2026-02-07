package sonar

import "net/http"

// WebservicesService handles communication with the webservices related methods
// of the SonarQube API.
// This service provides access to API documentation and response examples.
//
// Since: 4.2.
type WebservicesService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// WebservicesList represents the response from listing webservices.
type WebservicesList struct {
	// Webservices is the list of available webservices.
	Webservices []Webservice `json:"webServices,omitempty"`
}

// Webservice represents a webservice controller.
//
//nolint:govet // fieldalignment - structure kept for readability
type Webservice struct {
	// Actions is the list of actions available in the controller.
	Actions []WebserviceAction `json:"actions,omitempty"`
	// Description is the controller description.
	Description string `json:"description,omitempty"`
	// Path is the controller path.
	Path string `json:"path,omitempty"`
	// Since indicates when the controller was introduced.
	Since string `json:"since,omitempty"`
}

// WebserviceAction represents an action within a webservice.
//
//nolint:govet // fieldalignment - structure kept for readability
type WebserviceAction struct {
	// Changelog is the list of changes.
	Changelog []WebserviceChangelog `json:"changelog,omitempty"`
	// DeprecatedSince indicates when the action was deprecated.
	DeprecatedSince string `json:"deprecatedSince,omitempty"`
	// Description is the action description.
	Description string `json:"description,omitempty"`
	// HasResponseExample indicates if a response example exists.
	HasResponseExample bool `json:"hasResponseExample,omitempty"`
	// Internal indicates if the action is internal.
	Internal bool `json:"internal,omitempty"`
	// Key is the action key.
	Key string `json:"key,omitempty"`
	// Params is the list of parameters.
	Params []WebserviceParam `json:"params,omitempty"`
	// Post indicates if the action uses POST method.
	Post bool `json:"post,omitempty"`
	// Since indicates when the action was introduced.
	Since string `json:"since,omitempty"`
}

// WebserviceChangelog represents a changelog entry for an action.
type WebserviceChangelog struct {
	// Description is the changelog description.
	Description string `json:"description,omitempty"`
	// Version is the version when the change occurred.
	Version string `json:"version,omitempty"`
}

// WebserviceParam represents a parameter for an action.
//
//nolint:govet // fieldalignment - structure kept for readability
type WebserviceParam struct {
	// DefaultValue is the parameter default value.
	DefaultValue string `json:"defaultValue,omitempty"`
	// DeprecatedKey is the deprecated key for the parameter.
	DeprecatedKey string `json:"deprecatedKey,omitempty"`
	// DeprecatedKeySince indicates when the key was deprecated.
	DeprecatedKeySince string `json:"deprecatedKeySince,omitempty"`
	// DeprecatedSince indicates when the parameter was deprecated.
	DeprecatedSince string `json:"deprecatedSince,omitempty"`
	// Description is the parameter description.
	Description string `json:"description,omitempty"`
	// ExampleValue is an example value for the parameter.
	ExampleValue string `json:"exampleValue,omitempty"`
	// Internal indicates if the parameter is internal.
	Internal bool `json:"internal,omitempty"`
	// Key is the parameter key.
	Key string `json:"key,omitempty"`
	// MaxValuesAllowed is the maximum number of values allowed.
	MaxValuesAllowed int `json:"maxValuesAllowed,omitempty"`
	// MaximumLength is the maximum length of the value.
	MaximumLength int `json:"maximumLength,omitempty"`
	// MaximumValue is the maximum value allowed.
	MaximumValue int `json:"maximumValue,omitempty"`
	// MinimumLength is the minimum length of the value.
	MinimumLength int `json:"minimumLength,omitempty"`
	// MinimumValue is the minimum value allowed.
	MinimumValue int `json:"minimumValue,omitempty"`
	// PossibleValues is the list of possible values.
	PossibleValues []string `json:"possibleValues,omitempty"`
	// Required indicates if the parameter is required.
	Required bool `json:"required,omitempty"`
	// Since indicates when the parameter was introduced.
	Since string `json:"since,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// WebservicesListOption represents options for listing webservices.
type WebservicesListOption struct {
	// IncludeInternals includes internal actions and parameters.
	// Default: false.
	IncludeInternals bool `url:"include_internals,omitempty"`
}

// WebservicesResponseExampleOption represents options for getting a response example.
type WebservicesResponseExampleOption struct {
	// Action is the action key (required).
	Action string `url:"action,omitempty"`
	// Controller is the controller key (required).
	Controller string `url:"controller,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateListOpt validates the options for the List method.
func (s *WebservicesService) ValidateListOpt(_ *WebservicesListOption) error {
	// No required fields, all options are optional
	return nil
}

// ValidateResponseExampleOpt validates the options for the ResponseExample method.
func (s *WebservicesService) ValidateResponseExampleOpt(opt *WebservicesResponseExampleOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Action, "Action")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Controller, "Controller")
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// List lists available web services.
// Returns the list of webservices with their actions and parameters.
//
// API endpoint: GET /api/webservices/list.
// Since: 4.2.
func (s *WebservicesService) List(opt *WebservicesListOption) (*WebservicesList, *http.Response, error) {
	err := s.ValidateListOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "webservices/list", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(WebservicesList)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ResponseExample returns an example response for a web service action.
// Returns the response example as a string (the example format depends on the action).
//
// API endpoint: GET /api/webservices/response_example.
// Since: 4.4.
func (s *WebservicesService) ResponseExample(opt *WebservicesResponseExampleOption) (*string, *http.Response, error) {
	err := s.ValidateResponseExampleOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "webservices/response_example", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(string)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
