package sonargo

import (
	"net/http"
)

const (
	// MaxSettingValueLength is the maximum length for a setting value.
	MaxSettingValueLength = 4000
)

// SettingsService handles communication with the Settings related methods of the SonarQube API.
// Manage settings.
//
// Since: 6.1.
type SettingsService struct {
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// SettingsCheckSecretKey represents the response from checking if a secret key is available.
type SettingsCheckSecretKey struct {
	// SecretKeyAvailable indicates if a secret key is available.
	SecretKeyAvailable bool `json:"secretKeyAvailable,omitempty"`
}

// SettingsEncrypt represents the response from encrypting a value.
type SettingsEncrypt struct {
	// EncryptedValue is the encrypted setting value.
	EncryptedValue string `json:"encryptedValue,omitempty"`
}

// SettingsGenerateSecretKey represents the response from generating a secret key.
type SettingsGenerateSecretKey struct {
	// SecretKey is the generated secret key.
	SecretKey string `json:"secretKey,omitempty"`
}

// SettingsListDefinitions represents the response from listing setting definitions.
type SettingsListDefinitions struct {
	// Definitions is the list of setting definitions.
	Definitions []SettingDefinition `json:"definitions,omitempty"`
}

// SettingDefinition represents a setting definition.
type SettingDefinition struct {
	// Category is the category of the setting.
	Category string `json:"category,omitempty"`
	// DefaultValue is the default value of the setting.
	DefaultValue string `json:"defaultValue,omitempty"`
	// Description is the description of the setting.
	Description string `json:"description,omitempty"`
	// Key is the unique key of the setting.
	Key string `json:"key,omitempty"`
	// Name is the display name of the setting.
	Name string `json:"name,omitempty"`
	// SubCategory is the sub-category of the setting.
	SubCategory string `json:"subCategory,omitempty"`
	// Type is the type of the setting (STRING, TEXT, BOOLEAN, etc.).
	Type string `json:"type,omitempty"`
	// Options is the list of possible values for the setting.
	Options []string `json:"options,omitempty"`
	// Fields is the list of field definitions for property set types.
	Fields []SettingField `json:"fields,omitempty"`
	// MultiValues indicates if the setting accepts multiple values.
	MultiValues bool `json:"multiValues,omitempty"`
}

// SettingField represents a field within a property set setting.
type SettingField struct {
	// Description is the description of the field.
	Description string `json:"description,omitempty"`
	// Key is the unique key of the field.
	Key string `json:"key,omitempty"`
	// Name is the display name of the field.
	Name string `json:"name,omitempty"`
	// Type is the type of the field.
	Type string `json:"type,omitempty"`
	// Options is the list of possible values for the field.
	Options []string `json:"options,omitempty"`
}

// SettingsLoginMessage represents the response from getting the login message.
type SettingsLoginMessage struct {
	// Message is the formatted login message.
	Message string `json:"message,omitempty"`
}

// SettingsValues represents the response from listing setting values.
type SettingsValues struct {
	// Settings is the list of settings.
	Settings []SettingValue `json:"settings,omitempty"`
	// SetSecuredSettings is the list of secured settings that have a value set.
	SetSecuredSettings []string `json:"setSecuredSettings,omitempty"`
}

// SettingValue represents a setting value.
type SettingValue struct {
	// Key is the unique key of the setting.
	Key string `json:"key,omitempty"`
	// Value is the value of the setting (for single-value settings).
	Value string `json:"value,omitempty"`
	// Values is the values of the setting (for multi-value settings).
	Values []string `json:"values,omitempty"`
	// FieldValues is the values of the setting (for property set types).
	FieldValues []map[string]string `json:"fieldValues,omitempty"`
	// Inherited indicates if the value is inherited from a parent component.
	Inherited bool `json:"inherited,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// SettingsEncryptOption represents options for encrypting a value.
type SettingsEncryptOption struct {
	// Value is the setting value to encrypt (required).
	Value string `url:"value,omitempty"`
}

// SettingsListDefinitionsOption represents options for listing setting definitions.
type SettingsListDefinitionsOption struct {
	// Component is the component key to get definitions for (optional).
	Component string `url:"component,omitempty"`
}

// SettingsResetOption represents options for resetting settings.
type SettingsResetOption struct {
	// Component is the component key (optional).
	// Only keys for projects, applications, portfolios or subportfolios are accepted.
	Component string `url:"component,omitempty"`
	// Keys is the list of setting keys to reset (required).
	Keys []string `url:"keys,omitempty,comma"`
}

// SettingsSetOption represents options for setting a value.
type SettingsSetOption struct {
	// Component is the component key (optional).
	// Only keys for projects, applications, portfolios or subportfolios are accepted.
	Component string `url:"component,omitempty"`
	// Key is the setting key (required).
	Key string `url:"key,omitempty"`
	// Value is the setting value (optional).
	// To reset a value, please use the reset web service.
	// Maximum length: 4000 characters.
	Value string `url:"value,omitempty"`
	// Values is the setting multi-value (optional).
	// To set several values, the parameter must be called once for each value.
	Values []string `url:"values,omitempty"`
	// FieldValues is the setting field values for property set types (optional).
	FieldValues []string `url:"fieldValues,omitempty"`
}

// SettingsValuesOption represents options for listing setting values.
type SettingsValuesOption struct {
	// Component is the component key (optional).
	Component string `url:"component,omitempty"`
	// Keys is the list of setting keys (optional).
	Keys []string `url:"keys,omitempty,comma"`
}

// -----------------------------------------------------------------------------
// Validation Methods
// -----------------------------------------------------------------------------

// ValidateEncryptOpt validates the options for Encrypt.
func (s *SettingsService) ValidateEncryptOpt(opt *SettingsEncryptOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Value, "Value")
}

// ValidateListDefinitionsOpt validates the options for ListDefinitions.
func (s *SettingsService) ValidateListDefinitionsOpt(opt *SettingsListDefinitionsOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return nil
}

// ValidateResetOpt validates the options for Reset.
func (s *SettingsService) ValidateResetOpt(opt *SettingsResetOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	if len(opt.Keys) == 0 {
		return NewValidationError("Keys", "at least one key is required", ErrMissingRequired)
	}

	return nil
}

// ValidateSetOpt validates the options for Set.
func (s *SettingsService) ValidateSetOpt(opt *SettingsSetOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Key, "Key")
	if err != nil {
		return err
	}

	if opt.Value != "" {
		err = ValidateMaxLength(opt.Value, MaxSettingValueLength, "Value")
		if err != nil {
			return err
		}
	}

	return nil
}

// ValidateValuesOpt validates the options for Values.
func (s *SettingsService) ValidateValuesOpt(opt *SettingsValuesOption) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// CheckSecretKey checks if a secret key is available.
// Requires the 'Administer System' permission.
//
// Since: 6.1.
func (s *SettingsService) CheckSecretKey() (*SettingsCheckSecretKey, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "settings/check_secret_key", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(SettingsCheckSecretKey)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Encrypt encrypts a setting value.
// Requires 'Administer System' permission.
//
// Since: 6.1.
func (s *SettingsService) Encrypt(opt *SettingsEncryptOption) (*SettingsEncrypt, *http.Response, error) {
	err := s.ValidateEncryptOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "settings/encrypt", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(SettingsEncrypt)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GenerateSecretKey generates a secret key.
// Requires the 'Administer System' permission.
//
// Since: 6.1.
func (s *SettingsService) GenerateSecretKey() (*SettingsGenerateSecretKey, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "settings/generate_secret_key", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(SettingsGenerateSecretKey)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ListDefinitions lists settings definitions.
// Requires 'Browse' permission when a component is specified.
// To access licensed settings, authentication is required.
// To access secured settings, one of the following permissions is required:
//   - 'Execute Analysis'
//   - 'Administer System'
//   - 'Administer' rights on the specified component
//
// Since: 6.3.
func (s *SettingsService) ListDefinitions(opt *SettingsListDefinitionsOption) (*SettingsListDefinitions, *http.Response, error) {
	err := s.ValidateListDefinitionsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "settings/list_definitions", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(SettingsListDefinitions)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// LoginMessage returns the formatted login message, set to the 'sonar.login.message' property.
//
// Since: 9.8.
func (s *SettingsService) LoginMessage() (*SettingsLoginMessage, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "settings/login_message", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(SettingsLoginMessage)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Reset removes a setting value.
// The settings defined in conf/sonar.properties are read-only and can't be changed.
// Requires one of the following permissions:
//   - 'Administer System'
//   - 'Administer' rights on the specified component
//
// Since: 6.1.
func (s *SettingsService) Reset(opt *SettingsResetOption) (*http.Response, error) {
	err := s.ValidateResetOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "settings/reset", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Set updates a setting value.
// Either 'value' or 'values' must be provided.
// The settings defined in conf/sonar.properties are read-only and can't be changed.
// Requires one of the following permissions:
//   - 'Administer System'
//   - 'Administer' rights on the specified component
//
// Since: 6.1.
func (s *SettingsService) Set(opt *SettingsSetOption) (*http.Response, error) {
	err := s.ValidateSetOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "settings/set", opt)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Values lists settings values.
// If no value has been set for a setting, then the default value is returned.
// The settings from conf/sonar.properties are excluded from results.
// Requires 'Browse' or 'Execute Analysis' permission when a component is specified.
// Secured settings values are not returned by the endpoint.
//
// Since: 6.3.
func (s *SettingsService) Values(opt *SettingsValuesOption) (*SettingsValues, *http.Response, error) {
	err := s.ValidateValuesOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, "settings/values", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(SettingsValues)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
