package sonar

import (
	"fmt"
	"io"
	"net/http"
)

// AnalysisService handles communication with the Analysis related methods of
// the SonarQube V2 API.
type AnalysisService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Shared Types
// -----------------------------------------------------------------------------

// AnalysisRuleKey represents a rule key composed of a repository and rule identifier.
type AnalysisRuleKey struct {
	// Repository is the rule repository key.
	Repository string `json:"repository,omitempty"`
	// Rule is the rule identifier within the repository.
	Rule string `json:"rule,omitempty"`
}

// AnalysisParam represents a rule parameter key-value pair.
type AnalysisParam struct {
	// Key is the parameter key.
	Key string `json:"key,omitempty"`
	// Value is the parameter value.
	Value string `json:"value,omitempty"`
}

// AnalysisJre represents metadata for a Java Runtime Environment.
type AnalysisJre struct {
	// Arch is the CPU architecture (e.g. x64, aarch64).
	Arch string `json:"arch,omitempty"`
	// Filename is the JRE archive filename.
	Filename string `json:"filename,omitempty"`
	// Id is the unique identifier of the JRE.
	Id string `json:"id,omitempty"`
	// JavaPath is the path to the java executable within the archive.
	JavaPath string `json:"javaPath,omitempty"`
	// Os is the operating system (e.g. windows, linux, macos, alpine).
	Os string `json:"os,omitempty"`
	// Sha256 is the SHA-256 checksum of the JRE archive.
	Sha256 string `json:"sha256,omitempty"`
}

// AnalysisEngineInfo represents metadata for the Scanner Engine.
type AnalysisEngineInfo struct {
	// Filename is the scanner engine filename.
	Filename string `json:"filename,omitempty"`
	// Sha256 is the SHA-256 checksum of the scanner engine.
	Sha256 string `json:"sha256,omitempty"`
}

// AnalysisActiveRule represents an active rule for a project.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type AnalysisActiveRule struct {
	// CreatedAt is the datetime when the rule was created.
	CreatedAt string `json:"createdAt,omitempty"`
	// DeprecatedKeys is the list of deprecated rule keys.
	DeprecatedKeys []AnalysisRuleKey `json:"deprecatedKeys,omitempty"`
	// Impacts is a map of software quality to impact severity.
	Impacts map[string]string `json:"impacts,omitempty"`
	// InternalKey is the internal key of the rule.
	InternalKey string `json:"internalKey,omitempty"`
	// Language is the programming language key.
	Language string `json:"language,omitempty"`
	// Name is the rule name.
	Name string `json:"name,omitempty"`
	// Params is the list of rule parameters.
	Params []AnalysisParam `json:"params,omitempty"`
	// QProfileKey is the quality profile key.
	QProfileKey string `json:"qProfileKey,omitempty"`
	// RuleKey is the rule key.
	RuleKey AnalysisRuleKey `json:"ruleKey,omitzero"`
	// Severity is the rule severity.
	Severity string `json:"severity,omitempty"`
	// TemplateRuleKey is the template rule key if this is a custom rule.
	TemplateRuleKey string `json:"templateRuleKey,omitempty"`
	// UpdatedAt is the datetime when the rule was last updated.
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types (Query Parameters)
// -----------------------------------------------------------------------------

// AnalysisJresOptions contains query parameters for the GetJresMetadata method.
type AnalysisJresOptions struct {
	// Arch filters JREs by CPU architecture (x64, aarch64).
	Arch string `json:"arch,omitempty"`
	// Os filters JREs by operating system (windows, linux, macos, alpine).
	Os string `json:"os,omitempty"`
}

// AnalysisActiveRuleV2sOptions contains query parameters for the GetActiveRules method.
type AnalysisActiveRuleV2sOptions struct {
	// ProjectKey is the project key. This field is required.
	ProjectKey string `json:"projectKey"`
}

// -----------------------------------------------------------------------------
// Validation
// -----------------------------------------------------------------------------

// ValidateActiveRulesOpt validates the AnalysisActiveRuleV2sOptions.
func (s *AnalysisService) ValidateActiveRulesOpt(opt *AnalysisActiveRuleV2sOptions) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	return ValidateRequired(opt.ProjectKey, "ProjectKey")
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// GetVersion returns the Scanner Engine version as a plain text string.
func (s *AnalysisService) GetVersion() (*string, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodGet, "analysis/version", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result string

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// GetJresMetadata returns metadata for all available JREs, optionally filtered
// by operating system and CPU architecture.
func (s *AnalysisService) GetJresMetadata(opt *AnalysisJresOptions) ([]AnalysisJre, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodGet, "analysis/jres", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result []AnalysisJre

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DownloadJre downloads a JRE binary by ID into the provided writer.
// Set the Accept header to "application/octet-stream" to receive the binary.
func (s *AnalysisService) DownloadJre(jreID string, writer io.Writer) (*http.Response, error) {
	err := ValidateRequired(jreID, "Id")
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodGet, "analysis/jres/"+jreID, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/octet-stream")

	resp, err := s.client.Do(req, writer)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetJreMetadata returns metadata for a single JRE by ID.
func (s *AnalysisService) GetJreMetadata(jreID string) (*AnalysisJre, *http.Response, error) {
	err := ValidateRequired(jreID, "Id")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodGet, "analysis/jres/"+jreID, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(AnalysisJre)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DownloadScannerEngine downloads the Scanner Engine binary into the provided writer.
func (s *AnalysisService) DownloadScannerEngine(writer io.Writer) (*http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodGet, "analysis/engine", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/octet-stream")

	resp, err := s.client.Do(req, writer)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// GetScannerEngineMetadata returns metadata for the Scanner Engine.
func (s *AnalysisService) GetScannerEngineMetadata() (*AnalysisEngineInfo, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodGet, "analysis/engine", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(AnalysisEngineInfo)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetActiveRules returns all active rules for a specific project.
// Used by the scanner-engine.
func (s *AnalysisService) GetActiveRules(opt *AnalysisActiveRuleV2sOptions) ([]AnalysisActiveRule, *http.Response, error) {
	err := s.ValidateActiveRulesOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodGet, "analysis/active_rules", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result []AnalysisActiveRule

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
