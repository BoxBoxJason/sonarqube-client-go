package sonar

import (
	"context"
	"fmt"
	"net/http"
)

// BillingService handles communication with the billing, entitlement checks and consumption related methods of the
// SonarQube V2 API. This service is only available in Enterprise Edition.
type BillingService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

//nolint:gochecknoglobals // constant sets of allowed values
var (
	// allowedBillingGetEntitlementCheckResourceType is the set of allowed values for the corresponding option.
	allowedBillingGetEntitlementCheckResourceType = map[string]struct{}{
		"organization": {},
	}
	// allowedBillingCreateConsumptionRecordResourceType is the set of allowed values for the corresponding option.
	allowedBillingCreateConsumptionRecordResourceType = map[string]struct{}{
		"organization": {},
	}
)

// -----------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------

// BillingGetEntitlementCheckOptions contains parameters for the GetEntitlementCheck method.
type BillingGetEntitlementCheckOptions struct {
	// FeatureKey feature key to check entitlement for. This field is required.
	FeatureKey string `json:"featureKey"`
	// ResourceId organization ID, for compatibility with SQC.
	ResourceId string `json:"resourceId,omitempty"`
	// ResourceType resource type, only 'organization' is supported, for compatibility with SQC.
	// Allowed values: organization.
	// Default: organization.
	ResourceType string `json:"resourceType,omitempty"`
}

// BillingLimitRestResponse represents the LimitRestResponse object of the SonarQube V2 API.
type BillingLimitRestResponse struct {
	// Base consumption limit for the feature.
	Base int64 `json:"base,omitempty"`
	// Overage consumption limit for the feature.
	Overage int64 `json:"overage,omitempty"`
	// OverageEnabled whether overage is enabled for the feature.
	OverageEnabled bool `json:"overageEnabled,omitempty"`
}

// BillingPeriodRestResponse represents the PeriodRestResponse object of the SonarQube V2 API.
type BillingPeriodRestResponse struct {
	// Cadence metering cadence (DAILY, WEEKLY, MONTHLY, ANNUAL, PERPETUAL).
	Cadence string `json:"cadence,omitempty"`
	// Anchor optional period window anchor, e.g. 2026-07 for MONTHLY.
	Anchor string `json:"anchor,omitempty"`
}

// BillingMeteringRestResponse represents the MeteringRestResponse object of the SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type BillingMeteringRestResponse struct {
	// Used amount of credits already used for the feature.
	Used int64 `json:"used,omitempty"`
	// Period information for the metering.
	Period BillingPeriodRestResponse `json:"period,omitzero"`
}

// BillingConsumptionRestResponse represents the ConsumptionRestResponse object of the SonarQube V2
// API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type BillingConsumptionRestResponse struct {
	// Allowed whether consumption of the feature is allowed, also check for the entitled property.
	Allowed bool `json:"allowed,omitempty"`
	// Limit consumption limits for the feature.
	Limit BillingLimitRestResponse `json:"limit,omitzero"`
	// Metering consumption metering for the feature.
	Metering BillingMeteringRestResponse `json:"metering,omitzero"`
}

// BillingEntitlementRestResponse represents the EntitlementRestResponse object of the SonarQube V2
// API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type BillingEntitlementRestResponse struct {
	// FeatureKey internal feature key.
	FeatureKey string `json:"featureKey,omitempty"`
	// Entitled whether the feature is entitled for use, also check for the consumption.allowed
	// property if available.
	Entitled bool `json:"entitled,omitempty"`
	// Consumption information for a given feature, only available for CONSUMABLE features.
	Consumption BillingConsumptionRestResponse `json:"consumption,omitzero"`
	// Value for the feature.
	Value int64 `json:"value,omitempty"`
	// ExcludedValues excluded values for the feature.
	ExcludedValues []string `json:"excludedValues,omitempty"`
}

// BillingCreateConsumptionRecordOptions contains the request body for the CreateConsumptionRecord
// method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type BillingCreateConsumptionRecordOptions struct {
	// FeatureKey feature key for the consumed CONSUMABLE feature. This field is required.
	FeatureKey string `json:"featureKey"`
	// ResourceId organization ID, for compatibility with SQC.
	ResourceId string `json:"resourceId,omitempty"`
	// ResourceType resource type, only 'organization' is supported, for compatibility with SQC.
	// Allowed values: organization.
	ResourceType string `json:"resourceType,omitempty"`
	// Consumed units consumed in this event.
	Consumed *int64 `json:"consumed,omitempty"`
}

// BillingConsumptionRecordRestResponse represents the ConsumptionRecordRestResponse object of the
// SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type BillingConsumptionRecordRestResponse struct {
	// RecordId transient acceptance/tracing id. Not persisted and not a lookup key; retries
	// receive a new value.
	RecordId string `json:"recordId,omitempty"`
	// FeatureKey feature key that was recorded.
	FeatureKey string `json:"featureKey,omitempty"`
	// ResourceId organization ID the consumption was recorded against.
	ResourceId string `json:"resourceId,omitempty"`
	// ResourceType resource type the consumption was recorded against. Allowed values:
	// organization.
	ResourceType string `json:"resourceType,omitempty"`
	// Consumed units consumed in this event.
	Consumed int64 `json:"consumed,omitempty"`
	// ConsumedAt when the consumption was accepted.
	ConsumedAt string `json:"consumedAt,omitempty"`
}

// BillingSetOverageOptions contains the request body for the SetOverage method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type BillingSetOverageOptions struct {
	// FeatureKey feature key for the overage-eligible feature. This field is required.
	FeatureKey string `json:"featureKey"`
	// Enabled whether overage is enabled for this feature. This field is required.
	Enabled bool `json:"enabled"`
	// Limit maximum overage units, required when enabling.
	Limit *int64 `json:"limit,omitempty"`
	// TermsAccepted must be true when enabling, confirms acceptance of overage terms.
	TermsAccepted *bool `json:"termsAccepted,omitempty"`
}

// BillingLicenseFeatureRestResponse represents the LicenseFeatureRestResponse object of the
// SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type BillingLicenseFeatureRestResponse struct {
	// Name is the name.
	Name string `json:"name,omitempty"`
	// Parent is the parent.
	Parent string `json:"parent,omitempty"`
	// StartDate is the start date.
	StartDate string `json:"startDate,omitempty"`
	// EndDate is the end date.
	EndDate string `json:"endDate,omitempty"`
	// FeatureKey is the feature key.
	FeatureKey string `json:"featureKey,omitempty"`
	// MaxConsumption is the max consumption.
	MaxConsumption int64 `json:"maxConsumption,omitempty"`
	// MaxOverage is the max overage.
	MaxOverage int64 `json:"maxOverage,omitempty"`
	// OverageState is the overage state. Allowed values: NOT_ELIGIBLE, ELIGIBLE, NOT_ENABLED,
	// ENABLED, OFFLINE_BLOCKED.
	OverageState string `json:"overageState,omitempty"`
	// OveragesStep is the overages step.
	OveragesStep int64 `json:"overagesStep,omitempty"`
	// OverageGraceDays is the overage grace days.
	OverageGraceDays int32 `json:"overageGraceDays,omitempty"`
	// OverageUnitPrice is the overage unit price.
	OverageUnitPrice float64 `json:"overageUnitPrice,omitempty"`
	// OverageUnitCurrency is the overage unit currency.
	OverageUnitCurrency string `json:"overageUnitCurrency,omitempty"`
	// OverageLimit is the overage limit.
	OverageLimit int64 `json:"overageLimit,omitempty"`
	// MaxConsumptionUnit is the max consumption unit.
	MaxConsumptionUnit string `json:"maxConsumptionUnit,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateGetEntitlementCheckOpt validates the options for the GetEntitlementCheck method.
func (s *BillingService) ValidateGetEntitlementCheckOpt(opt *BillingGetEntitlementCheckOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.FeatureKey, "FeatureKey")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.ResourceType, allowedBillingGetEntitlementCheckResourceType, "ResourceType")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCreateConsumptionRecordOpt validates the options for the CreateConsumptionRecord method.
func (s *BillingService) ValidateCreateConsumptionRecordOpt(opt *BillingCreateConsumptionRecordOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.FeatureKey, "FeatureKey")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.ResourceType, allowedBillingCreateConsumptionRecordResourceType, "ResourceType")
	if err != nil {
		return err
	}

	return nil
}

// ValidateSetOverageOpt validates the options for the SetOverage method.
func (s *BillingService) ValidateSetOverageOpt(opt *BillingSetOverageOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.FeatureKey, "FeatureKey")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// GetEntitlementCheck returns the entitlement status for a given feature. Returns live entitlement
// status for a single feature. Used to get entitlement status, including consumption and remaining
// quota for CONSUMABLE features. Requires authentication.
//
// API endpoint: GET /api/v2/billing/entitlement-checks.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *BillingService) GetEntitlementCheck(ctx context.Context, opt *BillingGetEntitlementCheckOptions) (*BillingEntitlementRestResponse, *http.Response, error) {
	err := s.ValidateGetEntitlementCheckOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "billing/entitlement-checks", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(BillingEntitlementRestResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateConsumptionRecord records a consumption event for a CONSUMABLE feature. Records usage after
// a billable action. Does not enforce caps — use entitlement-checks first. Returns 202 because
// drain and LicenseSpring reporting remain asynchronous. Requires authentication.
//
// API endpoint: POST /api/v2/billing/consumption-records.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *BillingService) CreateConsumptionRecord(ctx context.Context, idempotencyKey string, opt *BillingCreateConsumptionRecordOptions) (*BillingConsumptionRecordRestResponse, *http.Response, error) {
	err := s.ValidateCreateConsumptionRecordOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "billing/consumption-records", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	if idempotencyKey != "" {
		req.Header.Set("Idempotency-Key", idempotencyKey)
	}

	result := new(BillingConsumptionRecordRestResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SetOverage enables,s updates, or disables the admin's overage arrangement for a feature. Enables,
// updates, or disables the overage arrangement for a feature. A single POST covers all three
// actions. Requires system administrator privileges.
//
// API endpoint: POST /api/v2/billing/overage.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *BillingService) SetOverage(ctx context.Context, opt *BillingSetOverageOptions) (*BillingLicenseFeatureRestResponse, *http.Response, error) {
	err := s.ValidateSetOverageOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "billing/overage", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(BillingLicenseFeatureRestResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
