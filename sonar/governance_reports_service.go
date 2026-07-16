package sonar

import (
	"bytes"
	"context"
	"net/http"
)

// GovernanceReportsService handles communication with the governance reports
// related methods of the SonarQube API.
// This service is only available in Enterprise Edition.
type GovernanceReportsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// GovernanceReportsStatus represents the metadata of a governance report.
//
// Verified against a live SonarQube 2025.2 Enterprise Edition instance:
// componentKey/hasFile do not exist on the real response; the actual shape
// is canDownload/canSubscribe/canAdmin/globalFrequency/subscribed, plus
// componentRecipients/globalRecipients when recipients have been configured.
//
//nolint:govet // fieldalignment: keeping fields grouped in the order the API returns them
type GovernanceReportsStatus struct {
	// ComponentRecipients is the list of recipient emails configured on the
	// component itself. Only present when component-level recipients exist.
	ComponentRecipients []string `json:"componentRecipients,omitempty"`
	// GlobalRecipients is the list of recipient emails configured at the
	// global level. Only present when global recipients exist.
	GlobalRecipients []string `json:"globalRecipients,omitempty"`
	// GlobalFrequency is the report frequency configured at the global level
	// (e.g. "monthly").
	GlobalFrequency string `json:"globalFrequency,omitempty"`
	// CanDownload indicates whether the current user can download the report.
	CanDownload bool `json:"canDownload"`
	// CanSubscribe indicates whether the current user can subscribe to the report.
	CanSubscribe bool `json:"canSubscribe"`
	// CanAdmin indicates whether the current user can administer the report
	// (update its frequency and recipients).
	CanAdmin bool `json:"canAdmin"`
	// Subscribed indicates whether the current user is subscribed to reports.
	Subscribed bool `json:"subscribed"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// GovernanceReportsDownloadOptions contains parameters for the Download method.
type GovernanceReportsDownloadOptions struct {
	// BranchKey is the branch key. Optional.
	BranchKey string `url:"branchKey,omitempty"`
	// ComponentID is the component id. Optional.
	ComponentID string `url:"componentId,omitempty"`
	// ComponentKey is the component key. Optional.
	ComponentKey string `url:"componentKey,omitempty"`
}

// GovernanceReportsStatusOptions contains parameters for the Status method.
type GovernanceReportsStatusOptions struct {
	// BranchKey is the branch key. Optional.
	BranchKey string `url:"branchKey,omitempty"`
	// ComponentID is the component id. Optional.
	ComponentID string `url:"componentId,omitempty"`
	// ComponentKey is the component key. Optional.
	ComponentKey string `url:"componentKey,omitempty"`
}

// GovernanceReportsSubscribeOptions contains parameters for the Subscribe method.
type GovernanceReportsSubscribeOptions struct {
	// BranchKey is the branch key. Optional.
	BranchKey string `url:"branchKey,omitempty"`
	// ComponentID is the component id. Optional.
	ComponentID string `url:"componentId,omitempty"`
	// ComponentKey is the component key. Optional.
	ComponentKey string `url:"componentKey,omitempty"`
}

// GovernanceReportsUnsubscribeOptions contains parameters for the Unsubscribe method.
type GovernanceReportsUnsubscribeOptions struct {
	// BranchKey is the branch key. Optional.
	BranchKey string `url:"branchKey,omitempty"`
	// ComponentID is the component id. Optional.
	ComponentID string `url:"componentId,omitempty"`
	// ComponentKey is the component key. Optional.
	ComponentKey string `url:"componentKey,omitempty"`
}

// GovernanceReportsUpdateFrequencyOptions contains parameters for the UpdateFrequency method.
type GovernanceReportsUpdateFrequencyOptions struct {
	// BranchKey is the branch key. Optional.
	BranchKey string `url:"branchKey,omitempty"`
	// ComponentID is the component id. Optional.
	ComponentID string `url:"componentId,omitempty"`
	// ComponentKey is the component key. Optional.
	ComponentKey string `url:"componentKey,omitempty"`
	// Frequency is the report frequency. Optional. Valid values: daily, weekly, monthly.
	Frequency string `url:"frequency,omitempty"`
}

// GovernanceReportsUpdateRecipientsOptions contains parameters for the UpdateRecipients method.
type GovernanceReportsUpdateRecipientsOptions struct {
	// BranchKey is the branch key. Optional.
	BranchKey string `url:"branchKey,omitempty"`
	// ComponentID is the component id. Optional.
	ComponentID string `url:"componentId,omitempty"`
	// ComponentKey is the component key. Optional.
	ComponentKey string `url:"componentKey,omitempty"`
	// Recipients is the list of recipient emails. This field is required.
	Recipients string `url:"recipients"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateUpdateRecipientsOpt validates the options for the UpdateRecipients method.
func (s *GovernanceReportsService) ValidateUpdateRecipientsOpt(opt *GovernanceReportsUpdateRecipientsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Recipients, "Recipients")
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Download downloads the PDF report of a portfolio, sub-portfolio, project or application.
// Requires 'Browse' permission on the component.
//
// API endpoint: GET /api/governance_reports/download.
// Since: 1.0.
// Enterprise Edition only.
func (s *GovernanceReportsService) Download(ctx context.Context, opt *GovernanceReportsDownloadOptions) ([]byte, *http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "governance_reports/download", opt)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", "application/pdf")

	var buf bytes.Buffer

	resp, err := s.client.Do(req, &buf)
	if err != nil {
		return nil, resp, err
	}

	return buf.Bytes(), resp, nil
}

// Status returns PDF report metadata (action rights, report availability etc.).
// Requires 'Browse' permission on the component.
//
// API endpoint: GET /api/governance_reports/status.
// Since: 1.0.
// Enterprise Edition only.
func (s *GovernanceReportsService) Status(ctx context.Context, opt *GovernanceReportsStatusOptions) (*GovernanceReportsStatus, *http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "governance_reports/status", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(GovernanceReportsStatus)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// Subscribe subscribes the current user to reports for a component.
// Requires authentication and 'Browse' permission on the component.
//
// API endpoint: POST /api/governance_reports/subscribe.
// Since: 1.0.
// Enterprise Edition only.
func (s *GovernanceReportsService) Subscribe(ctx context.Context, opt *GovernanceReportsSubscribeOptions) (*http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "governance_reports/subscribe", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Unsubscribe unsubscribes the current user from reports for a component.
// Requires authentication and 'Browse' permission on the component.
//
// API endpoint: POST /api/governance_reports/unsubscribe.
// Since: 1.0.
// Enterprise Edition only.
func (s *GovernanceReportsService) Unsubscribe(ctx context.Context, opt *GovernanceReportsUnsubscribeOptions) (*http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "governance_reports/unsubscribe", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// UpdateFrequency updates the frequency at which a report is sent for a component.
// Requires 'Administer' permission.
//
// API endpoint: POST /api/governance_reports/update_frequency.
// Since: 1.0.
// Enterprise Edition only.
func (s *GovernanceReportsService) UpdateFrequency(ctx context.Context, opt *GovernanceReportsUpdateFrequencyOptions) (*http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "governance_reports/update_frequency", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// UpdateRecipients updates the list of users who will receive reports for a component.
// Requires 'Administer' permission.
//
// API endpoint: POST /api/governance_reports/update_recipients.
// Since: 1.0.
// Enterprise Edition only.
func (s *GovernanceReportsService) UpdateRecipients(ctx context.Context, opt *GovernanceReportsUpdateRecipientsOptions) (*http.Response, error) {
	err := s.ValidateUpdateRecipientsOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "governance_reports/update_recipients", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
