package sonar

import (
	"context"
	"net/http"
)

// EditionsService handles communication with the editions related methods of the
// SonarQube API. This service is only available in Enterprise Edition.
type EditionsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// License represents a SonarQube Enterprise Edition license.
//
//nolint:govet // Field alignment is less important than logical grouping
type License struct {
	// ContactEmail is the contact email associated with the license.
	ContactEmail string `json:"contactEmail,omitempty"`
	// ExpiresAt is the license expiration date (ISO 8601 format).
	ExpiresAt string `json:"expiresAt,omitempty"`
	// IsExpired indicates whether the license has expired.
	IsExpired bool `json:"isExpired,omitempty"`
	// IsOfficialDistribution indicates whether the SonarQube instance is an
	// official distribution.
	IsOfficialDistribution bool `json:"isOfficialDistribution,omitempty"`
	// IsSupported indicates whether the license is currently valid and supported.
	IsSupported bool `json:"isSupported,omitempty"`
	// IsValidEdition indicates whether the license edition matches the running
	// SonarQube edition.
	IsValidEdition bool `json:"isValidEdition,omitempty"`
	// IsValidServerId indicates whether the license server ID matches the current
	// instance.
	IsValidServerId bool `json:"isValidServerId,omitempty"`
	// LOCsMax is the maximum number of lines of code allowed by the license.
	LOCsMax int64 `json:"locsMax,omitempty"`
	// LOCsRemaining is the number of lines of code remaining before the limit
	// is reached.
	LOCsRemaining int64 `json:"locsRemaining,omitempty"`
	// Organization is the name of the organization the license was issued to.
	Organization string `json:"organization,omitempty"`
	// ProductEdition is the SonarQube edition the license is for.
	ProductEdition string `json:"productEdition,omitempty"`
	// ServerId is the server ID the license is bound to.
	ServerId string `json:"serverId,omitempty"`
	// Type is the license type (e.g. PRODUCTION, EVALUATION).
	Type string `json:"type,omitempty"`
}

// LicenseGet wraps the License response from the show_license endpoint.
type LicenseGet struct {
	// License contains the license details.
	License License `json:"license,omitzero"`
}

// LicenseIsValid wraps the response from the is_valid_license endpoint.
type LicenseIsValid struct {
	// IsValid indicates whether the current license is valid.
	IsValid bool `json:"isValid,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// LicenseSetOptions contains parameters for the Set method.
type LicenseSetOptions struct {
	// License is the license key to activate. This field is required.
	License string `url:"license"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateSetOpt validates the options for the Set method.
func (s *EditionsService) ValidateSetOpt(opt *LicenseSetOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.License, "License")
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// ActivateGracePeriod enables a 7-day license grace period when the Server ID
// is invalid. Requires 'Administer System' permission.
//
// API endpoint: POST /api/editions/activate_grace_period.
// Since: 10.3.
// Enterprise Edition only.
func (s *EditionsService) ActivateGracePeriod(ctx context.Context) (*http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "editions/activate_grace_period", nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Get returns the details of the current SonarQube Enterprise Edition license.
// Requires 'Administer System' permission.
//
// API endpoint: GET /api/editions/show_license.
// Since: 7.2.
// Enterprise Edition only.
func (s *EditionsService) Get(ctx context.Context) (*LicenseGet, *http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "editions/show_license", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(LicenseGet)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// IsValidLicense returns the validity of the current license.
//
// API endpoint: GET /api/editions/is_valid_license.
// Since: 7.3.
// Enterprise Edition only.
func (s *EditionsService) IsValidLicense(ctx context.Context) (*LicenseIsValid, *http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "editions/is_valid_license", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(LicenseIsValid)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Set activates a new SonarQube Enterprise Edition license.
// Requires 'Administer System' permission.
//
// API endpoint: POST /api/editions/set_license.
// Since: 7.2.
// Enterprise Edition only.
func (s *EditionsService) Set(ctx context.Context, opt *LicenseSetOptions) (*http.Response, error) {
	err := s.ValidateSetOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "editions/set_license", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// UnsetLicense removes the current license.
// Requires 'Administer System' permission.
//
// API endpoint: POST /api/editions/unset_license.
// Since: 7.2.
// Enterprise Edition only.
func (s *EditionsService) UnsetLicense(ctx context.Context) (*http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "editions/unset_license", nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
