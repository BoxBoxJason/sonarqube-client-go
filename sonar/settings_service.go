// Manage settings.
package sonargo

import "net/http"

type SettingsService struct {
	client *Client
}

type SettingsCheckSecretKeyObject struct {
	SecretKeyAvailable bool `json:"secretKeyAvailable,omitempty"`
}

type SettingsEncryptObject struct {
	EncryptedValue string `json:"encryptedValue,omitempty"`
}

type SettingsGenerateSecretKeyObject struct {
	SecretKey string `json:"secretKey,omitempty"`
}

type SettingsListDefinitionsObject struct {
	Definitions []SettingsListDefinitionsObject_sub2 `json:"definitions,omitempty"`
}

type SettingsListDefinitionsObject_sub2 struct {
	Category     string                               `json:"category,omitempty"`
	DefaultValue string                               `json:"defaultValue,omitempty"`
	Description  string                               `json:"description,omitempty"`
	Fields       []SettingsListDefinitionsObject_sub1 `json:"fields,omitempty"`
	Key          string                               `json:"key,omitempty"`
	MultiValues  bool                                 `json:"multiValues,omitempty"`
	Name         string                               `json:"name,omitempty"`
	Options      []string                             `json:"options,omitempty"`
	SubCategory  string                               `json:"subCategory,omitempty"`
	Type         string                               `json:"type,omitempty"`
}

type SettingsListDefinitionsObject_sub1 struct {
	Description string   `json:"description,omitempty"`
	Key         string   `json:"key,omitempty"`
	Name        string   `json:"name,omitempty"`
	Options     []string `json:"options,omitempty"`
	Type        string   `json:"type,omitempty"`
}

type SettingsLoginMessageObject struct {
	Message string `json:"message,omitempty"`
}

type SettingsValuesObject struct {
	SetSecuredSettings []string                    `json:"setSecuredSettings,omitempty"`
	Settings           []SettingsValuesObject_sub2 `json:"settings,omitempty"`
}

type SettingsValuesObject_sub1 struct {
	Boolean string `json:"boolean,omitempty"`
	Text    string `json:"text,omitempty"`
}

type SettingsValuesObject_sub2 struct {
	FieldValues []SettingsValuesObject_sub1 `json:"fieldValues,omitempty"`
	Inherited   bool                        `json:"inherited,omitempty"`
	Key         string                      `json:"key,omitempty"`
	Value       string                      `json:"value,omitempty"`
	Values      []string                    `json:"values,omitempty"`
}

// CheckSecretKey Check if a secret key is available.<br>Requires the 'Administer System' permission.
func (s *SettingsService) CheckSecretKey() (v *SettingsCheckSecretKeyObject, resp *http.Response, err error) {
	req, err := s.client.NewRequest("GET", "settings/check_secret_key", nil)
	if err != nil {
		return
	}
	v = new(SettingsCheckSecretKeyObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type SettingsEncryptOption struct {
	Value string `url:"value,omitempty"` // Description:"Setting value to encrypt",ExampleValue:"my value"
}

// Encrypt Encrypt a setting value.<br>Requires 'Administer System' permission.
func (s *SettingsService) Encrypt(opt *SettingsEncryptOption) (v *SettingsEncryptObject, resp *http.Response, err error) {
	err = s.ValidateEncryptOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "settings/encrypt", opt)
	if err != nil {
		return
	}
	v = new(SettingsEncryptObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

// GenerateSecretKey Generate a secret key.<br>Requires the 'Administer System' permission
func (s *SettingsService) GenerateSecretKey() (v *SettingsGenerateSecretKeyObject, resp *http.Response, err error) {
	req, err := s.client.NewRequest("GET", "settings/generate_secret_key", nil)
	if err != nil {
		return
	}
	v = new(SettingsGenerateSecretKeyObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type SettingsListDefinitionsOption struct {
	Component string `url:"component,omitempty"` // Description:"Component key",ExampleValue:"my_project"
}

// ListDefinitions List settings definitions.<br>Requires 'Browse' permission when a component is specified<br/>To access licensed settings, authentication is required<br/>To access secured settings, one of the following permissions is required: <ul><li>'Execute Analysis'</li><li>'Administer System'</li><li>'Administer' rights on the specified component</li></ul>
func (s *SettingsService) ListDefinitions(opt *SettingsListDefinitionsOption) (v *SettingsListDefinitionsObject, resp *http.Response, err error) {
	err = s.ValidateListDefinitionsOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "settings/list_definitions", opt)
	if err != nil {
		return
	}
	v = new(SettingsListDefinitionsObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

// LoginMessage Returns the formatted login message, set to the 'sonar.login.message' property.
func (s *SettingsService) LoginMessage() (v *SettingsLoginMessageObject, resp *http.Response, err error) {
	req, err := s.client.NewRequest("GET", "settings/login_message", nil)
	if err != nil {
		return
	}
	v = new(SettingsLoginMessageObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}

type SettingsResetOption struct {
	Component string `url:"component,omitempty"` // Description:"Component key. Only keys for projects, applications, portfolios or subportfolios are accepted.",ExampleValue:"my_project"
	Keys      string `url:"keys,omitempty"`      // Description:"Comma-separated list of keys",ExampleValue:"sonar.links.scm,sonar.debt.hoursInDay"
}

// Reset Remove a setting value.<br>The settings defined in conf/sonar.properties are read-only and can't be changed.<br/>Requires one of the following permissions: <ul><li>'Administer System'</li><li>'Administer' rights on the specified component</li></ul>
func (s *SettingsService) Reset(opt *SettingsResetOption) (resp *http.Response, err error) {
	err = s.ValidateResetOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "settings/reset", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type SettingsSetOption struct {
	Component   string `url:"component,omitempty"`   // Description:"Component key. Only keys for projects, applications, portfolios or subportfolios are accepted.",ExampleValue:"my_project"
	FieldValues string `url:"fieldValues,omitempty"` // Description:"Setting field values. To set several values, the parameter must be called once for each value.",ExampleValue:"fieldValues={"firstField":"first value", "secondField":"second value", "thirdField":"third value"}"
	Key         string `url:"key,omitempty"`         // Description:"Setting key",ExampleValue:"sonar.core.serverBaseURL"
	Value       string `url:"value,omitempty"`       // Description:"Setting value. To reset a value, please use the reset web service.",ExampleValue:"http://my-sonarqube-instance.com"
	Values      string `url:"values,omitempty"`      // Description:"Setting multi value. To set several values, the parameter must be called once for each value.",ExampleValue:"values=firstValue&values=secondValue&values=thirdValue"
}

// Set Update a setting value.<br>Either 'value' or 'values' must be provided.<br> The settings defined in conf/sonar.properties are read-only and can't be changed.<br/>Requires one of the following permissions: <ul><li>'Administer System'</li><li>'Administer' rights on the specified component</li></ul>
func (s *SettingsService) Set(opt *SettingsSetOption) (resp *http.Response, err error) {
	err = s.ValidateSetOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("POST", "settings/set", opt)
	if err != nil {
		return
	}
	resp, err = s.client.Do(req, nil)
	if err != nil {
		return
	}
	return
}

type SettingsValuesOption struct {
	Component string `url:"component,omitempty"` // Description:"Component key",ExampleValue:"my_project"
	Keys      string `url:"keys,omitempty"`      // Description:"List of setting keys",ExampleValue:"sonar.test.inclusions,sonar.exclusions"
}

// Values List settings values.<br>If no value has been set for a setting, then the default value is returned.<br>The settings from conf/sonar.properties are excluded from results.<br>Requires 'Browse' or 'Execute Analysis' permission when a component is specified.<br/>Secured settings values are not returned by the endpoint.<br/>
func (s *SettingsService) Values(opt *SettingsValuesOption) (v *SettingsValuesObject, resp *http.Response, err error) {
	err = s.ValidateValuesOpt(opt)
	if err != nil {
		return
	}
	req, err := s.client.NewRequest("GET", "settings/values", opt)
	if err != nil {
		return
	}
	v = new(SettingsValuesObject)
	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}
	return
}
