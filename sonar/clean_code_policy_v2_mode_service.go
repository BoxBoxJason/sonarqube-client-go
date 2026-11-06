package sonar

import (
	"context"
	"net/http"
)

const (
	// InstanceModeMQR is the Multi-Quality Rule mode.
	InstanceModeMQR = "MQR"
	// InstanceModeStandardExperience is the Standard Experience mode.
	InstanceModeStandardExperience = "STANDARD_EXPERIENCE"
)

//nolint:gochecknoglobals // allowed-value sets
var allowedInstanceModes = map[string]struct{}{InstanceModeMQR: {}, InstanceModeStandardExperience: {}}

// CleanCodePolicyMode represents the instance mode.
type CleanCodePolicyMode struct {
	// Mode is the instance mode. Allowed values: MQR, STANDARD_EXPERIENCE.
	Mode string `json:"mode,omitempty"`
	// Modified indicates whether the mode was modified from its default.
	Modified bool `json:"modified,omitempty"`
}

// CleanCodePolicyUpdateModeOptions contains parameters for the UpdateMode method.
type CleanCodePolicyUpdateModeOptions struct {
	// Mode is the new instance mode. Allowed values: MQR, STANDARD_EXPERIENCE. This field is required.
	Mode string `json:"mode"`
}

// ValidateUpdateModeOpt validates the options for the UpdateMode method.
func (s *CleanCodePolicyService) ValidateUpdateModeOpt(opt *CleanCodePolicyUpdateModeOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Mode, "Mode")
	if err != nil {
		return err
	}

	return IsValueAuthorized(opt.Mode, allowedInstanceModes, "Mode")
}

// GetMode retrieves the current instance mode.
//
// API endpoint: GET /api/v2/clean-code-policy/mode.
func (s *CleanCodePolicyService) GetMode(ctx context.Context) (*CleanCodePolicyMode, *http.Response, error) {
	return callV2[CleanCodePolicyMode](ctx, s.client, http.MethodGet, "clean-code-policy/mode", nil, nil)
}

// UpdateMode updates the current instance mode.
// Requires the 'Administer System' permission.
//
// API endpoint: PATCH /api/v2/clean-code-policy/mode.
func (s *CleanCodePolicyService) UpdateMode(ctx context.Context, opt *CleanCodePolicyUpdateModeOptions) (*CleanCodePolicyMode, *http.Response, error) {
	err := s.ValidateUpdateModeOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	return callV2[CleanCodePolicyMode](ctx, s.client, http.MethodPatch, "clean-code-policy/mode", nil, opt)
}
