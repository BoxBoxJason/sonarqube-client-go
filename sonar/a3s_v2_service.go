package sonar

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// A3sService handles communication with the Advanced Security Analysis (A3S) analyses, contexts and entitlements related methods of the
// SonarQube V2 API. This service is only available in Enterprise Edition.
type A3sService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

const (
	// A3sValueStandard is the "STANDARD" enumeration value.
	A3sValueStandard = "STANDARD"
)

//nolint:gochecknoglobals // constant sets of allowed values
var (
	// allowedA3sGetAnalysisStatsPeriod is the set of allowed values for the corresponding option.
	allowedA3sGetAnalysisStatsPeriod = map[string]struct{}{
		"LAST_DAY":     {},
		"LAST_7_DAYS":  {},
		"LAST_30_DAYS": {},
	}
	// allowedA3sCreatePublicAnalysisFileScope is the set of allowed values for the corresponding option.
	allowedA3sCreatePublicAnalysisFileScope = map[string]struct{}{
		"MAIN": {},
		"TEST": {},
	}
	// allowedA3sCreatePublicAnalysisAnalysisDepth is the set of allowed values for the corresponding option.
	allowedA3sCreatePublicAnalysisAnalysisDepth = map[string]struct{}{
		A3sValueStandard: {},
		"DEEP":           {},
	}
	// allowedA3sCreatePrivateAnalysisFileScope is the set of allowed values for the corresponding option.
	allowedA3sCreatePrivateAnalysisFileScope = map[string]struct{}{
		"MAIN": {},
		"TEST": {},
	}
	// allowedA3sCreatePrivateAnalysisAnalysisDepth is the set of allowed values for the corresponding option.
	allowedA3sCreatePrivateAnalysisAnalysisDepth = map[string]struct{}{
		A3sValueStandard: {},
		"DEEP":           {},
	}
)

// -----------------------------------------------------------------------------
// Types
// -----------------------------------------------------------------------------

// A3sGetAnalysisStatsOptions contains parameters for the GetAnalysisStats method.
type A3sGetAnalysisStatsOptions struct {
	// Period for the statistics. This field is required. Allowed values: LAST_DAY, LAST_7_DAYS,
	// LAST_30_DAYS.
	Period string `json:"period"`
}

// A3sProjectAnalysisCount represents a SonarQube V2 API object.
//
// Analysis count for a project.
type A3sProjectAnalysisCount struct {
	// Id the project ID.
	Id string `json:"id,omitempty"`
	// Name the project name.
	Name string `json:"name,omitempty"`
	// AnalysisCount number of analyses performed for the project over the requested period.
	AnalysisCount int64 `json:"analysisCount,omitempty"`
}

// A3sAnalysisStats represents a SonarQube V2 API object.
//
// Analysis statistics for an organization.
type A3sAnalysisStats struct {
	// Id the organization ID these statistics are about.
	Id string `json:"id,omitempty"`
	// Projects all allowed projects with their count of analysis over the select period, ordered
	// by analysis count descending.
	Projects []A3sProjectAnalysisCount `json:"projects,omitempty"`
}

// A3sGetPublicCollectionConfigOptions contains parameters for the GetPublicCollectionConfig method.
type A3sGetPublicCollectionConfigOptions struct {
	// OrganizationKey is the organization key. This field is required.
	OrganizationKey string `json:"organizationKey"`
	// ProjectKey is the project key. This field is required.
	ProjectKey string `json:"projectKey"`
}

// A3sCollectionConfig represents the CollectionConfig object of the SonarQube V2 API.
type A3sCollectionConfig struct {
	// Enabled whether collection is enabled for this project.
	Enabled bool `json:"enabled,omitempty"`
}

// A3sOrgEntitlement represents a SonarQube V2 API object.
//
// A3S entitlement for an organization.
type A3sOrgEntitlement struct {
	// Id an organization id.
	Id string `json:"id,omitempty"`
	// Allowed whether the organization is currently allowed to use A3S.
	Allowed bool `json:"allowed,omitempty"`
	// HasEntitlement whether the organization has any valid access path to A3S. When allowed is
	// false, true means the organization is still entitled but cannot use it right now, typically
	// because it has reached its current limit; false means the organization is not entitled at
	// all.
	HasEntitlement bool `json:"hasEntitlement,omitempty"`
}

// A3sGetPrivateCollectionConfigOptions contains parameters for the GetPrivateCollectionConfig method.
type A3sGetPrivateCollectionConfigOptions struct {
	// OrganizationKey is the organization key. This field is required.
	OrganizationKey string `json:"organizationKey"`
	// ProjectKey is the project key. This field is required.
	ProjectKey string `json:"projectKey"`
}

// A3sSearchPrivateContextsOptions contains parameters for the SearchPrivateContexts method.
type A3sSearchPrivateContextsOptions struct {
	// BranchId the identifier of a branch on SonarQube (a UUID on SonarQube Cloud, possibly a
	// legacy identifier on SonarQube Server). Must specify either branchId or branchName, but not
	// both.
	BranchId string `json:"branchId,omitempty"`
	// BranchName the name of a branch. Must specify either branchId or branchName, but not both.
	// When branchName is specified, projectId is required.
	BranchName string `json:"branchName,omitempty"`
	// ProjectId the identifier of a project (a UUID on SonarQube Cloud, possibly a legacy
	// identifier on SonarQube Server). Required when branchName is specified to disambiguate
	// branches with the same name across different projects.
	ProjectId string `json:"projectId,omitempty"`
}

// A3sSummarizedContext represents a SonarQube V2 API object.
//
// Summarized structure of a Context.
type A3sSummarizedContext struct {
	// Id the unique identifier of a context.
	Id string `json:"id,omitempty"`
	// AnalysisId the identifier of an analysis on SonarQube (a UUID on SonarQube Cloud, possibly a
	// legacy identifier on SonarQube Server).
	AnalysisId string `json:"analysisId,omitempty"`
	// Kind context kind. Must be unique for a particular analysis.
	Kind string `json:"kind,omitempty"`
}

// A3sContexts represents the Contexts object of the SonarQube V2 API.
type A3sContexts struct {
	// Contexts is the contexts.
	Contexts []A3sSummarizedContext `json:"contexts,omitempty"`
	// Page is the page.
	Page PageResponseV2 `json:"page,omitzero"`
}

// A3sContextItem represents a SonarQube V2 API object.
//
// Item object stored in a context
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type A3sContextItem struct {
	// Id the unique identifier of the item in its context.
	Id string `json:"id,omitempty"`
	// Type is the type. Allowed values: FILE.
	Type string `json:"type,omitempty"`
	// Sha256 A sha256 of the item content to detect changes.
	Sha256 string `json:"sha256,omitempty"`
	// Missing whether or the service has the context item already saved.
	Missing bool `json:"missing,omitempty"`
	// DownloadUrl the URL to download the item.
	DownloadUrl string `json:"downloadUrl,omitempty"`
}

// A3sContext represents a SonarQube V2 API object.
//
// Base structure of a Context.
type A3sContext struct {
	// Id the unique identifier of a context.
	Id string `json:"id,omitempty"`
	// AnalysisId the identifier of an analysis on SonarQube (a UUID on SonarQube Cloud, possibly a
	// legacy identifier on SonarQube Server).
	AnalysisId string `json:"analysisId,omitempty"`
	// Kind context kind. Must be unique for a particular analysis.
	Kind string `json:"kind,omitempty"`
	// Metadata any JSON-formatted object.
	Metadata string `json:"metadata,omitempty"`
	// Items is the items.
	Items []A3sContextItem `json:"items,omitempty"`
}

// A3sPublicAnalysisFile represents the PublicAnalysisFile object of the SonarQube V2 API.
type A3sPublicAnalysisFile struct {
	// Path project-relative path of the file.
	Path string `json:"path,omitempty"`
	// Content the original content of the file.
	Content string `json:"content,omitempty"`
	// Scope is the scope. Allowed values: MAIN, TEST.
	Scope string `json:"scope,omitempty"`
}

// A3sCreatePublicAnalysisOptions contains the request body for the CreatePublicAnalysis method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type A3sCreatePublicAnalysisOptions struct {
	// OrganizationId the unique identifier of the organization.
	OrganizationId string `json:"organizationId,omitempty"`
	// OrganizationKey the key of the organization.
	OrganizationKey string `json:"organizationKey,omitempty"`
	// ProjectId the unique identifier of the project.
	ProjectId string `json:"projectId,omitempty"`
	// ProjectKey the key of the project.
	ProjectKey string `json:"projectKey,omitempty"`
	// BranchId the unique identifier of the branch to retrieve the latest analysis.
	BranchId string `json:"branchId,omitempty"`
	// BranchName the branch to retrieve the latest analysis.
	BranchName string `json:"branchName,omitempty"`
	// FilePath project-relative path of the file. Deprecated — use 'files' instead.
	FilePath string `json:"filePath,omitempty"`
	// FileContent the original content of the file. Deprecated — use 'files' instead.
	FileContent string `json:"fileContent,omitempty"`
	// FileScope is the file scope. Allowed values: MAIN, TEST.
	FileScope string `json:"fileScope,omitempty"`
	// Files list of files to analyze. When provided, takes precedence over the top-level
	// filePath/fileContent fields. No duplicate paths allowed.
	Files []A3sPublicAnalysisFile `json:"files,omitempty"`
	// AnalysisDepth is the analysis depth. Allowed values: STANDARD, DEEP.
	AnalysisDepth string `json:"analysisDepth,omitempty"`
}

// A3sTextRange represents a SonarQube V2 API object.
//
// Location of an issue in a file.
type A3sTextRange struct {
	// StartLine is the start line.
	StartLine int32 `json:"startLine,omitempty"`
	// EndLine is the end line.
	EndLine int32 `json:"endLine,omitempty"`
	// StartOffset is the start offset.
	StartOffset int32 `json:"startOffset,omitempty"`
	// EndOffset is the end offset.
	EndOffset int32 `json:"endOffset,omitempty"`
}

// A3sHubIssueFlowsInnerLocationsInner represents the HubIssueFlowsInnerLocationsInner object of the
// SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type A3sHubIssueFlowsInnerLocationsInner struct {
	// TextRange is the text range.
	TextRange A3sTextRange `json:"textRange,omitzero"`
	// Message is the message.
	Message string `json:"message,omitempty"`
	// FilePath is the file path.
	FilePath string `json:"filePath,omitempty"`
}

// A3sHubIssueFlowsInner represents the HubIssueFlowsInner object of the SonarQube V2 API.
type A3sHubIssueFlowsInner struct {
	// Type the type of flow. Can be "UNDEFINED", "DATA", or "EXECUTION". Allowed values:
	// UNDEFINED, DATA, EXECUTION.
	Type string `json:"type,omitempty"`
	// Description the description of the flow, if any.
	Description string `json:"description,omitempty"`
	// Locations is the locations.
	Locations []A3sHubIssueFlowsInnerLocationsInner `json:"locations,omitempty"`
}

// A3sHubIssue represents a SonarQube V2 API object.
//
// Structure of an Issue
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type A3sHubIssue struct {
	// Id the unique identifier of the issue.
	Id string `json:"id,omitempty"`
	// FilePath project-relative path of the file containing the issue.
	FilePath string `json:"filePath,omitempty"`
	// Message primary message of the issue.
	Message string `json:"message,omitempty"`
	// Rule the rule key (e.g., java:S1854).
	Rule string `json:"rule,omitempty"`
	// TextRange is the text range.
	TextRange A3sTextRange `json:"textRange,omitzero"`
	// Flows secondary locations and flows for the issue.
	Flows []A3sHubIssueFlowsInner `json:"flows,omitempty"`
}

// A3sPatchResult represents a SonarQube V2 API object.
//
// Result of applying a patch to analyze changes.
type A3sPatchResult struct {
	// NewIssues issues that appear only in the patched version.
	NewIssues []A3sHubIssue `json:"newIssues,omitempty"`
	// MatchedIssues issues that exist in both original and patched versions.
	MatchedIssues []A3sHubIssue `json:"matchedIssues,omitempty"`
	// ClosedIssues uuids of issues that existed in the original but are gone in the patch.
	ClosedIssues []string `json:"closedIssues,omitempty"`
}

// A3sAnalysisError represents a SonarQube V2 API object.
//
// Error that occurred during analysis.
type A3sAnalysisError struct {
	// Code error code indicating the type of failure. Allowed values: SERVICE_CALL_ERROR,
	// PARSE_ERROR, INVALID_CONTEXT.
	Code string `json:"code,omitempty"`
	// Message human-readable error message.
	Message string `json:"message,omitempty"`
}

// A3sAnalysis represents a SonarQube V2 API object.
//
// Structure of an Analysis.
type A3sAnalysis struct {
	// Id the unique identifier of an analysis.
	Id string `json:"id,omitempty"`
	// Issues is the issues.
	Issues []A3sHubIssue `json:"issues,omitempty"`
	// PatchResult is the patch result.
	PatchResult A3sPatchResult `json:"patchResult,omitzero"`
	// Errors that occurred during analysis from language analysis services.
	Errors []A3sAnalysisError `json:"errors,omitempty"`
}

// A3sCreatePublicContextOptions contains the request body for the CreatePublicContext method.
type A3sCreatePublicContextOptions struct {
	// AnalysisId the identifier of an analysis on SonarQube (a UUID on SonarQube Cloud, possibly a
	// legacy identifier on SonarQube Server). This field is required.
	AnalysisId string `json:"analysisId"`
	// BranchName optional branch name.
	BranchName string `json:"branchName,omitempty"`
	// Kind context kind. Must be unique for a particular analysis. This field is required.
	Kind string `json:"kind"`
	// Metadata any JSON-formatted object. This field is required.
	Metadata string `json:"metadata"`
	// Items is the items.
	Items []A3sContextItem `json:"items,omitempty"`
}

// A3sContextWithUploadUrl represents the ContextWithUploadUrl object of the SonarQube V2 API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type A3sContextWithUploadUrl struct {
	// Id the unique identifier of a context.
	Id string `json:"id,omitempty"`
	// AnalysisId the identifier of an analysis on SonarQube (a UUID on SonarQube Cloud, possibly a
	// legacy identifier on SonarQube Server).
	AnalysisId string `json:"analysisId,omitempty"`
	// Kind context kind. Must be unique for a particular analysis.
	Kind string `json:"kind,omitempty"`
	// Metadata any JSON-formatted object.
	Metadata string `json:"metadata,omitempty"`
	// Items is the items.
	Items []A3sContextItem `json:"items,omitempty"`
	// UploadUrl an URL to upload missing context items.
	UploadUrl string `json:"uploadUrl,omitempty"`
}

// A3sCreatePrivateContextOptions contains the request body for the CreatePrivateContext method.
type A3sCreatePrivateContextOptions struct {
	// AnalysisId the identifier of an analysis on SonarQube (a UUID on SonarQube Cloud, possibly a
	// legacy identifier on SonarQube Server). This field is required.
	AnalysisId string `json:"analysisId"`
	// BranchName optional branch name.
	BranchName string `json:"branchName,omitempty"`
	// Kind context kind. Must be unique for a particular analysis. This field is required.
	Kind string `json:"kind"`
	// Metadata any JSON-formatted object. This field is required.
	Metadata string `json:"metadata"`
	// Items is the items.
	Items []A3sContextItem `json:"items,omitempty"`
}

// A3sPrivateAnalysisFile represents the PrivateAnalysisFile object of the SonarQube V2 API.
type A3sPrivateAnalysisFile struct {
	// Path project-relative path of the file.
	Path string `json:"path,omitempty"`
	// Content the original content of the file.
	Content string `json:"content,omitempty"`
	// Scope is the scope. Allowed values: MAIN, TEST.
	Scope string `json:"scope,omitempty"`
	// PatchContent the patch content to apply (unified diff format). If provided for one file,
	// must be provided for all files in the request.
	PatchContent string `json:"patchContent,omitempty"`
}

// A3sCreatePrivateAnalysisOptions contains the request body for the CreatePrivateAnalysis method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type A3sCreatePrivateAnalysisOptions struct {
	// OrganizationId the unique identifier of the organization. This field is required.
	OrganizationId string `json:"organizationId"`
	// ProjectId the unique identifier of the project.
	ProjectId string `json:"projectId,omitempty"`
	// ProjectKey the key of the project (deprecated, use projectId instead).
	ProjectKey string `json:"projectKey,omitempty"`
	// BranchId the unique identifier of the branch to retrieve the latest analysis.
	BranchId string `json:"branchId,omitempty"`
	// BranchName the name of the branch to retrieve the latest analysis (used when branchId is not
	// provided).
	BranchName string `json:"branchName,omitempty"`
	// FilePath project-relative path of the file. Deprecated — use 'files' instead.
	FilePath string `json:"filePath,omitempty"`
	// FileContent the original content of the file. Deprecated — use 'files' instead.
	FileContent string `json:"fileContent,omitempty"`
	// PatchContent the patch content to apply (unified diff format). Deprecated — use
	// 'patchContent' inside 'files' instead.
	PatchContent string `json:"patchContent,omitempty"`
	// FileScope is the file scope. Allowed values: MAIN, TEST.
	FileScope string `json:"fileScope,omitempty"`
	// Files list of files to analyze. When provided, takes precedence over the top-level
	// filePath/fileContent fields. No duplicate paths allowed.
	Files []A3sPrivateAnalysisFile `json:"files,omitempty"`
	// AnalysisDepth is the analysis depth. Allowed values: STANDARD, DEEP.
	AnalysisDepth string `json:"analysisDepth,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateGetAnalysisStatsOpt validates the options for the GetAnalysisStats method.
func (s *A3sService) ValidateGetAnalysisStatsOpt(opt *A3sGetAnalysisStatsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Period, "Period")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.Period, allowedA3sGetAnalysisStatsPeriod, "Period")
	if err != nil {
		return err
	}

	return nil
}

// ValidateGetPublicCollectionConfigOpt validates the options for the GetPublicCollectionConfig method.
func (s *A3sService) ValidateGetPublicCollectionConfigOpt(opt *A3sGetPublicCollectionConfigOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.OrganizationKey, "OrganizationKey")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	return nil
}

// ValidateGetPrivateCollectionConfigOpt validates the options for the GetPrivateCollectionConfig method.
func (s *A3sService) ValidateGetPrivateCollectionConfigOpt(opt *A3sGetPrivateCollectionConfigOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.OrganizationKey, "OrganizationKey")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCreatePublicAnalysisOpt validates the options for the CreatePublicAnalysis method.
func (s *A3sService) ValidateCreatePublicAnalysisOpt(opt *A3sCreatePublicAnalysisOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := IsValueAuthorized(opt.FileScope, allowedA3sCreatePublicAnalysisFileScope, "FileScope")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.AnalysisDepth, allowedA3sCreatePublicAnalysisAnalysisDepth, "AnalysisDepth")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCreatePublicContextOpt validates the options for the CreatePublicContext method.
func (s *A3sService) ValidateCreatePublicContextOpt(opt *A3sCreatePublicContextOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AnalysisId, "AnalysisId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Kind, "Kind")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Metadata, "Metadata")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCreatePrivateContextOpt validates the options for the CreatePrivateContext method.
func (s *A3sService) ValidateCreatePrivateContextOpt(opt *A3sCreatePrivateContextOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.AnalysisId, "AnalysisId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Kind, "Kind")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Metadata, "Metadata")
	if err != nil {
		return err
	}

	return nil
}

// ValidateCreatePrivateAnalysisOpt validates the options for the CreatePrivateAnalysis method.
func (s *A3sService) ValidateCreatePrivateAnalysisOpt(opt *A3sCreatePrivateAnalysisOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.OrganizationId, "OrganizationId")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.FileScope, allowedA3sCreatePrivateAnalysisFileScope, "FileScope")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.AnalysisDepth, allowedA3sCreatePrivateAnalysisAnalysisDepth, "AnalysisDepth")
	if err != nil {
		return err
	}

	return nil
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// GetAnalysisStats gets analysis statistics for an organization.
//
// API endpoint: GET /api/v2/a3s/analysis-stats/{id}.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *A3sService) GetAnalysisStats(ctx context.Context, organizationID string, opt *A3sGetAnalysisStatsOptions) (*A3sAnalysisStats, *http.Response, error) {
	err := ValidateRequired(organizationID, "organizationID")
	if err != nil {
		return nil, nil, err
	}

	err = s.ValidateGetAnalysisStatsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "a3s/analysis-stats/"+url.PathEscape(organizationID), opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(A3sAnalysisStats)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetPublicCollectionConfig gets collection configuration for an organization/project.
//
// API endpoint: GET /api/v2/a3s/collection-config.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *A3sService) GetPublicCollectionConfig(ctx context.Context, opt *A3sGetPublicCollectionConfigOptions) (*A3sCollectionConfig, *http.Response, error) {
	err := s.ValidateGetPublicCollectionConfigOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "a3s/collection-config", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(A3sCollectionConfig)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetPublicOrgEntitlement gets the A3S entitlement for an organization.
//
// API endpoint: GET /api/v2/a3s/org-entitlement/{id}.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *A3sService) GetPublicOrgEntitlement(ctx context.Context, organizationID string) (*A3sOrgEntitlement, *http.Response, error) {
	err := ValidateRequired(organizationID, "organizationID")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "a3s/org-entitlement/"+url.PathEscape(organizationID), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(A3sOrgEntitlement)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetPrivateOrgEntitlement gets the A3S entitlement for an organization (internal).
//
// API endpoint: GET /api/v2/a3s/private/a3s-analysis/org-entitlement/{id}.
// Enterprise Edition only.
func (s *A3sService) GetPrivateOrgEntitlement(ctx context.Context, organizationID string) (*A3sOrgEntitlement, *http.Response, error) {
	err := ValidateRequired(organizationID, "organizationID")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "a3s/private/a3s-analysis/org-entitlement/"+url.PathEscape(organizationID), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(A3sOrgEntitlement)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetPrivateCollectionConfig gets collection configuration for an organization/project (private).
//
// API endpoint: GET /api/v2/a3s/private/a3s-frontend-context/collection-config.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *A3sService) GetPrivateCollectionConfig(ctx context.Context, opt *A3sGetPrivateCollectionConfigOptions) (*A3sCollectionConfig, *http.Response, error) {
	err := s.ValidateGetPrivateCollectionConfigOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "a3s/private/a3s-frontend-context/collection-config", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(A3sCollectionConfig)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// SearchPrivateContexts searches for contexts.
//
// API endpoint: GET /api/v2/a3s/private/a3s-frontend-context/contexts.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *A3sService) SearchPrivateContexts(ctx context.Context, opt *A3sSearchPrivateContextsOptions) (*A3sContexts, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "a3s/private/a3s-frontend-context/contexts", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(A3sContexts)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetPrivateContext gets a context.
//
// API endpoint: GET /api/v2/a3s/private/a3s-frontend-context/contexts/{id}.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *A3sService) GetPrivateContext(ctx context.Context, contextID string) (*A3sContext, *http.Response, error) {
	err := ValidateRequired(contextID, "contextID")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "a3s/private/a3s-frontend-context/contexts/"+url.PathEscape(contextID), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(A3sContext)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreatePublicAnalysis creates a new analysis.
//
// API endpoint: POST /api/v2/a3s/analyses.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *A3sService) CreatePublicAnalysis(ctx context.Context, invocationID string, opt *A3sCreatePublicAnalysisOptions) (*A3sAnalysis, *http.Response, error) {
	err := s.ValidateCreatePublicAnalysisOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "a3s/analyses", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	if invocationID != "" {
		req.Header.Set("X-Sonar-Invocation-Id", invocationID)
	}

	result := new(A3sAnalysis)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreatePublicContext creates a new context.
//
// API endpoint: POST /api/v2/a3s/contexts.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *A3sService) CreatePublicContext(ctx context.Context, opt *A3sCreatePublicContextOptions) (*A3sContextWithUploadUrl, *http.Response, error) {
	err := s.ValidateCreatePublicContextOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "a3s/contexts", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(A3sContextWithUploadUrl)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreatePrivateContext creates a new context.
//
// API endpoint: POST /api/v2/a3s/private/a3s-frontend-context/contexts.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *A3sService) CreatePrivateContext(ctx context.Context, opt *A3sCreatePrivateContextOptions) (*A3sContextWithUploadUrl, *http.Response, error) {
	err := s.ValidateCreatePrivateContextOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "a3s/private/a3s-frontend-context/contexts", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(A3sContextWithUploadUrl)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreatePrivateAnalysis creates a new analysis.
//
// API endpoint: POST /api/v2/a3s/private/analyses.
// Enterprise Edition only.
func (s *A3sService) CreatePrivateAnalysis(ctx context.Context, opt *A3sCreatePrivateAnalysisOptions) (*A3sAnalysis, *http.Response, error) {
	err := s.ValidateCreatePrivateAnalysisOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "a3s/private/analyses", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(A3sAnalysis)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UploadContext uploads the content of a context created with CreatePublicContext or
// CreatePrivateContext, as a zip archive (the endpoint only accepts application/zip).
//
// API endpoint: PUT /api/v2/a3s/context-uploads/{contextId}.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *A3sService) UploadContext(ctx context.Context, contextID string, content io.Reader) (*http.Response, error) {
	err := ValidateRequired(contextID, "contextID")
	if err != nil {
		return nil, err
	}

	if content == nil {
		return nil, NewValidationError("content", "is required", ErrMissingRequired)
	}

	//nolint:exhaustruct // only the fields relevant to a raw upload are set
	req, err := s.client.NewSonarQubeAPIRequest(ctx, SonarAPIRequestParameters{
		Method:  http.MethodPut,
		Path:    v2BasePath + "a3s/context-uploads/" + url.PathEscape(contextID),
		Headers: map[string]string{headerContentType: "application/zip"},
		Body:    content,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
