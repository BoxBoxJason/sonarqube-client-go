package sonar

import (
	"context"
	"fmt"
	"net/http"
)

// EntitlementsService handles communication with the Entitlements related
// methods of the SonarQube V2 API. It covers license activation/deactivation
// (online, offline and legacy flows) and purchasable feature listing.
//
// This service is only available in Enterprise Edition. All underlying
// endpoints are marked internal by SonarQube (x-sonar-internal) and their
// request/response contract may change without notice between SonarQube
// versions.
//
// It supersedes the deprecated V1 license methods on EditionsService
// (Get, Set, UnsetLicense), which are deprecated since SonarQube 2025.6.
type EntitlementsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// LicenseFeatureV2 represents a feature covered by a SonarQube license.
type LicenseFeatureV2 struct {
	// Name is the name of the feature.
	Name string `json:"name,omitempty"`
	// Parent is the parent feature key, if any.
	Parent string `json:"parent,omitempty"`
	// StartDate is the date from which the feature is enabled.
	StartDate string `json:"startDate,omitempty"`
	// EndDate is the date until which the feature is enabled.
	EndDate string `json:"endDate,omitempty"`
}

// LicenseV2 represents detailed information about the current SonarQube license.
type LicenseV2 struct {
	LicenseKey             string             `json:"licenseKey,omitempty"`
	LastRefreshDate        string             `json:"lastRefreshDate,omitempty"`
	Edition                string             `json:"edition,omitempty"`
	StartDate              string             `json:"startDate,omitempty"`
	ExpirationDate         string             `json:"expirationDate,omitempty"`
	GracePeriodEndDate     string             `json:"gracePeriodEndDate,omitempty"`
	ServerId               string             `json:"serverId,omitempty"`
	Type                   string             `json:"type,omitempty"`
	ContactEmail           string             `json:"contactEmail,omitempty"`
	Features               []LicenseFeatureV2 `json:"features,omitempty"`
	RemainingLocThreshold  int64              `json:"remainingLocThreshold,omitempty"`
	Loc                    int64              `json:"loc,omitempty"`
	MaxLoc                 int64              `json:"maxLoc,omitempty"`
	ExtraDays              int32              `json:"extraDays,omitempty"`
	GracePeriodExpired     bool               `json:"gracePeriodExpired,omitempty"`
	ActivatedOnline        bool               `json:"activatedOnline,omitempty"`
	CanActivateGracePeriod bool               `json:"canActivateGracePeriod,omitempty"`
	Supported              bool               `json:"supported,omitempty"`
	Expired                bool               `json:"expired,omitempty"`
	Legacy                 bool               `json:"legacy,omitempty"`
	ValidEdition           bool               `json:"validEdition,omitempty"`
	ValidServerId          bool               `json:"validServerId,omitempty"`
	OfficialDistribution   bool               `json:"officialDistribution,omitempty"`
	Disabled               bool               `json:"disabled,omitempty"`
}

// PurchasableFeatureV2 represents a feature that can be purchased for the
// current SonarQube edition.
type PurchasableFeatureV2 struct {
	FeatureKey  string `json:"featureKey,omitempty"`
	Parent      string `json:"parent,omitempty"`
	URL         string `json:"url,omitempty"`
	IsEnabled   bool   `json:"isEnabled,omitempty"`
	IsAvailable bool   `json:"isAvailable,omitempty"`
}

// entitlementsOfflineFileResponse is a defined string type (rather than a bare
// string) used solely to decode opaque offline-activation-flow response
// bodies. client.Do treats a destination of type *string as an opaque
// text/plain payload and forces an "Accept: text/plain" request header for
// it. The API spec declares these endpoints' 200 response as
// "application/json" with a binary-formatted string schema, not
// "text/plain" — V2 endpoints strictly enforce their declared content type,
// so requesting "Accept: text/plain" against a JSON-only V2 endpoint returns
// 406 Not Acceptable rather than the payload. Using a distinct named type
// keeps client.Do on its default JSON-decode path (default "Accept:
// application/json"), which both matches the endpoint's contract and
// correctly unescapes the JSON string payload. See ArchitectureService.FileGraph
// for the same pattern.
type entitlementsOfflineFileResponse string

// -----------------------------------------------------------------------------
// Request Types
// -----------------------------------------------------------------------------

// EntitlementsActivateOnlineOptions contains parameters for activating a
// license online via a third-party license key.
type EntitlementsActivateOnlineOptions struct {
	// LicenseKey is the new license key to activate. This field is required.
	LicenseKey string `json:"licenseKey"`
}

// EntitlementsActivateLegacyOptions contains parameters for activating a
// license received from Sonar using the legacy flow.
type EntitlementsActivateLegacyOptions struct {
	// LicenseKey is the new license key to activate. This field is required.
	LicenseKey string `json:"licenseKey"`
}

// EntitlementsGetOfflineActivationRequestOptions contains parameters for
// retrieving a .req file for offline license activation.
type EntitlementsGetOfflineActivationRequestOptions struct {
	// LicenseKey is the unique license key associated with the license, in the
	// format 'ABCD-EFGH-IJKL-MNOP'. Sent as the 'License-Key' HTTP header.
	// This field is required.
	LicenseKey string `json:"-"`
}

// EntitlementsActivateOfflineOptions contains parameters for uploading a
// license file to complete offline activation.
type EntitlementsActivateOfflineOptions struct {
	// License is the content of a valid .lic file. This field is required.
	License string `json:"license"`
	// LicenseKey is the license key for the license. This field is required.
	LicenseKey string `json:"licenseKey"`
}

// -----------------------------------------------------------------------------
// Validation
// -----------------------------------------------------------------------------

// ValidateActivateOnlineOpt validates the options for the ActivateOnline method.
func (s *EntitlementsService) ValidateActivateOnlineOpt(opt *EntitlementsActivateOnlineOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.LicenseKey, "LicenseKey")
}

// ValidateActivateLegacyOpt validates the options for the ActivateLegacy method.
func (s *EntitlementsService) ValidateActivateLegacyOpt(opt *EntitlementsActivateLegacyOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.LicenseKey, "LicenseKey")
}

// ValidateGetOfflineActivationRequestOpt validates the options for the
// GetOfflineActivationRequest method.
func (s *EntitlementsService) ValidateGetOfflineActivationRequestOpt(opt *EntitlementsGetOfflineActivationRequestOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.LicenseKey, "LicenseKey")
}

// ValidateActivateOfflineOpt validates the options for the ActivateOffline method.
func (s *EntitlementsService) ValidateActivateOfflineOpt(opt *EntitlementsActivateOfflineOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.License, "License")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.LicenseKey, "LicenseKey")
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// ActivateOnline activates a new SonarQube Enterprise Edition license using a
// third-party license key. Only third party license keys are accepted.
// Requires 'Administer System' permission.
//
// API endpoint: POST /api/v2/entitlements/online-activation.
// Enterprise Edition only.
func (s *EntitlementsService) ActivateOnline(ctx context.Context, opt *EntitlementsActivateOnlineOptions) (*http.Response, error) {
	err := s.ValidateActivateOnlineOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "entitlements/online-activation", nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return s.client.Do(req, nil)
}

// DeactivateOffline fully deactivates a license that was activated using the
// offline method (e.g. in cases of migrating the installation to a new
// server). Requires 'Administer System' permission.
//
// API endpoint: POST /api/v2/entitlements/offline-deactivation.
// Enterprise Edition only.
func (s *EntitlementsService) DeactivateOffline(ctx context.Context) (*string, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "entitlements/offline-deactivation", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result entitlementsOfflineFileResponse

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	strResult := string(result)

	return &strResult, resp, nil
}

// GetOfflineActivationRequest retrieves a .req file for offline activation of
// a license. The license key must be passed via the License-Key HTTP header.
// Requires 'Administer System' permission.
//
// API endpoint: GET /api/v2/entitlements/offline-activation.
// Enterprise Edition only.
func (s *EntitlementsService) GetOfflineActivationRequest(ctx context.Context, opt *EntitlementsGetOfflineActivationRequestOptions) (*string, *http.Response, error) {
	err := s.ValidateGetOfflineActivationRequestOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "entitlements/offline-activation", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("License-Key", opt.LicenseKey)

	var result entitlementsOfflineFileResponse

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	strResult := string(result)

	return &strResult, resp, nil
}

// ActivateOffline uploads a license file to complete offline activation. The
// content should be the text of the license file, along with a valid license
// key. Requires 'Administer System' permission.
//
// API endpoint: POST /api/v2/entitlements/offline-activation.
// Enterprise Edition only.
func (s *EntitlementsService) ActivateOffline(ctx context.Context, opt *EntitlementsActivateOfflineOptions) (*http.Response, error) {
	err := s.ValidateActivateOfflineOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "entitlements/offline-activation", nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return s.client.Do(req, nil)
}

// ActivateLegacy activates a new SonarQube Enterprise Edition license using a
// license key received from Sonar. Only license keys received from Sonar are
// accepted. Requires 'Administer System' permission.
//
// API endpoint: POST /api/v2/entitlements/legacy-activation.
// Enterprise Edition only.
func (s *EntitlementsService) ActivateLegacy(ctx context.Context, opt *EntitlementsActivateLegacyOptions) (*http.Response, error) {
	err := s.ValidateActivateLegacyOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "entitlements/legacy-activation", nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return s.client.Do(req, nil)
}

// GetLicense returns information about the current license.
// Requires 'Administer System' permission.
//
// API endpoint: GET /api/v2/entitlements/license.
// Enterprise Edition only.
func (s *EntitlementsService) GetLicense(ctx context.Context) (*LicenseV2, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "entitlements/license", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(LicenseV2)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DeleteLicense deletes the current license. Both license keys received from
// Sonar and from a 3rd party system can be removed.
// Requires 'Administer System' permission.
//
// API endpoint: DELETE /api/v2/entitlements/license.
// Enterprise Edition only.
func (s *EntitlementsService) DeleteLicense(ctx context.Context) (*http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodDelete, "entitlements/license", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return s.client.Do(req, nil)
}

// UpdateLicense fetches the latest information about the license and updates
// it on the SonarQube Server instance.
// Requires 'Administer System' permission.
//
// API endpoint: PATCH /api/v2/entitlements/license.
// Enterprise Edition only.
func (s *EntitlementsService) UpdateLicense(ctx context.Context) (*http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "entitlements/license", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return s.client.Do(req, nil)
}

// GetPurchasableFeatures returns the list of all available purchasable
// features for this edition, including the ones already purchased.
//
// API endpoint: GET /api/v2/entitlements/purchasable-features.
// Enterprise Edition only.
func (s *EntitlementsService) GetPurchasableFeatures(ctx context.Context) ([]PurchasableFeatureV2, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "entitlements/purchasable-features", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result []PurchasableFeatureV2

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
