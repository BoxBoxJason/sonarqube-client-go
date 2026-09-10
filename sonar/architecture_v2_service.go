package sonar

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

// ArchitectureService handles communication with the architecture V2 API endpoints.
// This service is only available in Enterprise Edition. The underlying endpoint is
// marked internal by SonarQube (x-sonar-internal) and its request/response contract
// may change without notice between SonarQube versions.
type ArchitectureService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// ArchitectureFileGraphOptions contains parameters for the FileGraph method.
type ArchitectureFileGraphOptions struct {
	// ProjectKey is the project key. This field is required.
	ProjectKey string `json:"projectKey"`
	// BranchKey is the branch key. This field is required.
	BranchKey string `json:"branchKey"`
	// Source is the language/analyzer that produced this graph, e.g. "java", "python", "js".
	// This field is required.
	Source string `json:"source"`
}

// ArchitectureSearchGraphsOptions contains parameters for the SearchGraphs method.
type ArchitectureSearchGraphsOptions struct {
	// ProjectKey is the key of the project. This field is required.
	ProjectKey string `json:"projectKey"`
	// BranchKey is the key of the branch. This field is required.
	BranchKey string `json:"branchKey"`
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// ArchitectureGraphMetadata represents the metadata of a single graph as returned
// by SearchGraphs. The SonarQube V2 API spec declares this as a generic/loosely
// typed schema, so it is decoded into a string-keyed map rather than a fixed struct.
type ArchitectureGraphMetadata map[string]any

// architectureSearchGraphsResponse is the wrapper object returned by the search
// graphs endpoint; only its "graphs" array is exposed to callers of SearchGraphs.
type architectureSearchGraphsResponse struct {
	// Graphs is the list of graph metadata entries.
	Graphs []ArchitectureGraphMetadata `json:"graphs,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateFileGraphOpt validates the options for the FileGraph method.
func (s *ArchitectureService) ValidateFileGraphOpt(opt *ArchitectureFileGraphOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.BranchKey, "BranchKey")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.Source, "Source")
}

// ValidateSearchGraphsOpt validates the options for the SearchGraphs method.
func (s *ArchitectureService) ValidateSearchGraphsOpt(opt *ArchitectureSearchGraphsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ProjectKey, "ProjectKey")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.BranchKey, "BranchKey")
}

// ValidateGetGraphOpt validates the graphID parameter for the GetGraph method.
func (s *ArchitectureService) ValidateGetGraphOpt(graphID string) error {
	return ValidateRequired(graphID, "graphID")
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// architectureFileGraphResponse is a defined string type (rather than a bare
// string) used solely to decode the FileGraph response body. client.Do treats
// a destination of type *string as an opaque text/plain payload and forces an
// "Accept: text/plain" request header for it. Unlike other opaque *string
// endpoints in this SDK (e.g. AnalysisService.GetVersion), the API spec
// declares this endpoint's 200 response as "application/json" with a string
// schema, not "text/plain" — and live verification against a SonarQube
// 2025.2 Enterprise instance confirmed that V2 endpoints strictly enforce
// their declared content type: requesting "Accept: text/plain" against a
// JSON-only V2 endpoint returns 406 Not Acceptable rather than the payload.
// Using a distinct named type keeps client.Do on its default JSON-decode
// path (default "Accept: application/json"), which both matches the
// endpoint's contract and correctly unescapes the JSON string payload.
type architectureFileGraphResponse string

// FileGraph returns the file dependency graph for a project branch, for the given
// source language. Requires 'Browse' permission on the project.
//
// The SonarQube API documents this endpoint's response as an opaque JSON string; its
// exact payload format (e.g. serialized graph nodes/edges, DOT graph text) is not
// published, so the decoded string is returned as-is for callers to parse as needed.
//
// API endpoint: GET /api/v2/architecture/file-graph.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *ArchitectureService) FileGraph(ctx context.Context, opt *ArchitectureFileGraphOptions) (*string, *http.Response, error) {
	err := s.ValidateFileGraphOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/file-graph", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result architectureFileGraphResponse

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	strResult := string(result)

	return &strResult, resp, nil
}

// SearchGraphs returns the metadata of all graphs currently available for a project
// branch. This endpoint does not include the graph data itself; use GetGraph to
// retrieve it. Requires 'Browse' permission on the project.
//
// API endpoint: GET /api/v2/architecture/graphs.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *ArchitectureService) SearchGraphs(ctx context.Context, opt *ArchitectureSearchGraphsOptions) ([]ArchitectureGraphMetadata, *http.Response, error) {
	err := s.ValidateSearchGraphsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/graphs", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(architectureSearchGraphsResponse)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result.Graphs, resp, nil
}

// architectureGraphResponse is a defined string type (rather than a bare string)
// used solely to decode the GetGraph response body. Like architectureFileGraphResponse
// above, this keeps client.Do on its default JSON-decode path (Accept:
// application/json) instead of the *string/text/plain special case, which matches
// this endpoint's declared "application/graph+json" content type carrying a JSON
// string payload.
type architectureGraphResponse string

// GetGraph returns the graph data for the given graph ID, as produced by analysis.
// Requires 'Browse' permission on the project.
//
// The SonarQube API documents this endpoint's response as an opaque string (content
// type "application/graph+json"); its exact payload format is not published, so the
// decoded string is returned as-is for callers to parse as needed.
//
// API endpoint: GET /api/v2/architecture/graphs/{id}.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *ArchitectureService) GetGraph(ctx context.Context, graphID string) (*string, *http.Response, error) {
	err := s.ValidateGetGraphOpt(graphID)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/graphs/"+url.PathEscape(graphID), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result architectureGraphResponse

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	strResult := string(result)

	return &strResult, resp, nil
}

// =============================================================================
// Architecture directives, models, snapshots and structure
//
// The endpoints below were added by SonarQube's architecture-as-code feature.
// They are Enterprise Edition only; several are marked internal by SonarQube
// (x-sonar-internal) and their request/response contracts may change without
// notice between SonarQube versions.
// =============================================================================

// -----------------------------------------------------------------------------
// Allowed Values
// -----------------------------------------------------------------------------

//nolint:gochecknoglobals // constant sets of allowed values
var (
	allowedArchitectureEcosystems = map[string]struct{}{
		"java": {},
		"js":   {},
		"ts":   {},
		"py":   {},
		"cs":   {},
	}
	allowedArchitectureGraphTypes = map[string]struct{}{
		"file_graph":      {},
		"namespace_graph": {},
	}
	allowedArchitectureDirectiveTypes = map[string]struct{}{
		"REMOVE_EDGE": {},
		"RENAME_NODE": {},
		"MOVE_NODE":   {},
	}
)

// -----------------------------------------------------------------------------
// Shared Types
// -----------------------------------------------------------------------------

// ArchitectureConstraint is a directed "from -> to" edge constraint used by
// intended-architecture groups and model perspectives.
type ArchitectureConstraint struct {
	// From is the key of the source node.
	From string `json:"from,omitempty"`
	// To is the key of the destination node.
	To string `json:"to,omitempty"`
}

// -----------------------------------------------------------------------------
// Directives
// -----------------------------------------------------------------------------

// ArchitectureDirectiveData describes the change a directive applies to an
// architecture graph.
type ArchitectureDirectiveData struct {
	// ToNodeKey is the key of the destination node. May be empty depending on Type.
	ToNodeKey string `json:"toNodeKey,omitempty"`
	// FromNodeKey is the key of the source node.
	FromNodeKey string `json:"fromNodeKey,omitempty"`
	// GraphPerspective is the perspective the directive applies to.
	GraphPerspective string `json:"graphPerspective,omitempty"`
	// GraphType is the graph type. One of "file_graph", "namespace_graph".
	GraphType string `json:"graphType,omitempty"`
	// GraphEcoSystem is the language ecosystem. One of "java", "js", "ts", "py", "cs".
	GraphEcoSystem string `json:"graphEcoSystem,omitempty"`
	// Type is the directive type. One of "REMOVE_EDGE", "RENAME_NODE", "MOVE_NODE".
	Type string `json:"type,omitempty"`
	// Description is an optional human-readable description.
	Description string `json:"description,omitempty"`
}

// ArchitectureDirective represents a stored architecture directive.
type ArchitectureDirective struct {
	// Id is the directive's unique identifier.
	Id string `json:"id,omitempty"`
	// ProjectId is the identifier of the project the directive belongs to.
	ProjectId string `json:"projectId,omitempty"`
	// Origin describes where the directive originated from.
	Origin string `json:"origin,omitempty"`
	// DirectiveData is the change the directive applies.
	DirectiveData ArchitectureDirectiveData `json:"directiveData,omitzero"`
}

// ArchitectureDirectivesList is the paginated response returned by ListDirectives.
type ArchitectureDirectivesList struct {
	// Directives is the current page of directives.
	Directives []ArchitectureDirective `json:"directives,omitempty"`
	// Page contains pagination information.
	Page PageResponseV2 `json:"page,omitzero"`
}

// ArchitectureListDirectivesOptions contains parameters for the ListDirectives method.
//
//nolint:govet // fieldalignment: embedded pagination is kept first for readability
type ArchitectureListDirectivesOptions struct {
	PaginationParamsV2

	// ProjectId is the project identifier. This field is required.
	ProjectId string `json:"projectId"`
}

// ArchitectureCreateDirectiveOptions contains the request body for the CreateDirective method.
type ArchitectureCreateDirectiveOptions struct {
	// ProjectId is the project identifier. This field is required.
	ProjectId string `json:"projectId"`
	// Origin describes where the directive originated from. Optional.
	Origin string `json:"origin,omitempty"`
	// DirectiveData is the change the directive applies. This field is required.
	DirectiveData ArchitectureDirectiveData `json:"directiveData"`
}

// ValidateListDirectivesOpt validates the options for the ListDirectives method.
func (s *ArchitectureService) ValidateListDirectivesOpt(opt *ArchitectureListDirectivesOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.ProjectId, "ProjectId")
}

// ValidateCreateDirectiveOpt validates the options for the CreateDirective method.
func (s *ArchitectureService) ValidateCreateDirectiveOpt(opt *ArchitectureCreateDirectiveOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.ProjectId, "ProjectId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.DirectiveData.FromNodeKey, "DirectiveData.FromNodeKey")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.DirectiveData.GraphType, allowedArchitectureGraphTypes, "DirectiveData.GraphType")
	if err != nil {
		return err
	}

	err = IsValueAuthorized(opt.DirectiveData.GraphEcoSystem, allowedArchitectureEcosystems, "DirectiveData.GraphEcoSystem")
	if err != nil {
		return err
	}

	return IsValueAuthorized(opt.DirectiveData.Type, allowedArchitectureDirectiveTypes, "DirectiveData.Type")
}

// ValidateDirectiveIDOpt validates a directive id path parameter.
func (s *ArchitectureService) ValidateDirectiveIDOpt(directiveID string) error {
	return ValidateRequired(directiveID, "directiveID")
}

// ListDirectives returns the paginated list of architecture directives stored for a project.
//
// API endpoint: GET /api/v2/architecture/directives.
// Enterprise Edition only.
func (s *ArchitectureService) ListDirectives(ctx context.Context, opt *ArchitectureListDirectivesOptions) (*ArchitectureDirectivesList, *http.Response, error) {
	err := s.ValidateListDirectivesOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/directives", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureDirectivesList)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateDirective stores a new architecture directive for a project.
//
// API endpoint: POST /api/v2/architecture/directives.
// Enterprise Edition only.
func (s *ArchitectureService) CreateDirective(ctx context.Context, opt *ArchitectureCreateDirectiveOptions) (*ArchitectureDirective, *http.Response, error) {
	err := s.ValidateCreateDirectiveOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "architecture/directives", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureDirective)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetDirective returns a single architecture directive by id.
//
// API endpoint: GET /api/v2/architecture/directives/{id}.
// Enterprise Edition only.
func (s *ArchitectureService) GetDirective(ctx context.Context, directiveID string) (*ArchitectureDirective, *http.Response, error) {
	err := s.ValidateDirectiveIDOpt(directiveID)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/directives/"+url.PathEscape(directiveID), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureDirective)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DeleteDirective removes an architecture directive by id.
//
// API endpoint: DELETE /api/v2/architecture/directives/{id}.
// Enterprise Edition only.
func (s *ArchitectureService) DeleteDirective(ctx context.Context, directiveID string) (*http.Response, error) {
	err := s.ValidateDirectiveIDOpt(directiveID)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodDelete, "architecture/directives/"+url.PathEscape(directiveID), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// -----------------------------------------------------------------------------
// Intended architecture
// -----------------------------------------------------------------------------

// ArchitectureIntendedGroup is a node grouping in an intended-architecture export.
type ArchitectureIntendedGroup struct {
	// Label is the group label.
	Label string `json:"label,omitempty"`
	// Interface indicates whether the group is an interface group.
	Interface *bool `json:"interface,omitempty"`
	// InterfaceAccess indicates whether interface access is allowed.
	InterfaceAccess *bool `json:"interfaceAccess,omitempty"`
	// Patterns is the list of glob patterns that assign nodes to the group.
	Patterns []string `json:"patterns,omitempty"`
	// Groups is the list of nested sub-groups.
	Groups []ArchitectureIntendedGroup `json:"groups,omitempty"`
}

// ArchitectureIntendedItem is one entry of an intended-architecture export.
type ArchitectureIntendedItem struct {
	// Id is the intended-architecture item identifier.
	Id string `json:"id,omitempty"`
	// Language is the language ecosystem. One of "java", "js", "ts", "py", "cs".
	Language string `json:"language,omitempty"`
	// Type is the item type. One of "file", "namespace".
	Type string `json:"type,omitempty"`
	// Constraints is the list of "from -> to" edge constraints.
	Constraints []ArchitectureConstraint `json:"constraints,omitempty"`
	// Groups is the list of node groups.
	Groups []ArchitectureIntendedGroup `json:"groups,omitempty"`
}

// ArchitectureIntendedOptions contains parameters for the GetIntendedArchitecture method.
type ArchitectureIntendedOptions struct {
	// ProjectId is the project identifier. This field is required.
	ProjectId string `json:"projectId"`
}

// ValidateIntendedOpt validates the options for the GetIntendedArchitecture method.
func (s *ArchitectureService) ValidateIntendedOpt(opt *ArchitectureIntendedOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.ProjectId, "ProjectId")
}

// GetIntendedArchitecture exports the intended architecture defined for a project.
//
// API endpoint: GET /api/v2/architecture/intended.
// Enterprise Edition only.
func (s *ArchitectureService) GetIntendedArchitecture(ctx context.Context, opt *ArchitectureIntendedOptions) ([]ArchitectureIntendedItem, *http.Response, error) {
	err := s.ValidateIntendedOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/intended", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result []ArchitectureIntendedItem

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// -----------------------------------------------------------------------------
// Models
// -----------------------------------------------------------------------------

// ArchitectureModelGroup is a node grouping within a model perspective.
type ArchitectureModelGroup struct {
	// Label is the group label.
	Label string `json:"label,omitempty"`
	// Description is an optional description.
	Description string `json:"description,omitempty"`
	// PatternId identifies the pattern set the group is derived from.
	PatternId string `json:"patternId,omitempty"`
	// Interface indicates whether the group is an interface group.
	Interface *bool `json:"interface,omitempty"`
	// InterfaceAccess indicates whether interface access is allowed.
	InterfaceAccess *bool `json:"interfaceAccess,omitempty"`
	// Patterns is the list of glob patterns that assign nodes to the group.
	Patterns []string `json:"patterns,omitempty"`
	// Groups is the list of nested sub-groups.
	Groups []ArchitectureModelGroup `json:"groups,omitempty"`
}

// ArchitectureModelPerspective is one perspective (a language/qualifier view) of a model.
type ArchitectureModelPerspective struct {
	// Qualifiers is the node qualifier the perspective applies to. One of "file", "namespace".
	Qualifiers string `json:"qualifiers,omitempty"`
	// Language is the language ecosystem. One of "java", "js", "ts", "py", "cs".
	Language string `json:"language,omitempty"`
	// Label is an optional perspective label.
	Label string `json:"label,omitempty"`
	// Description is an optional perspective description.
	Description string `json:"description,omitempty"`
	// Groups is the list of node groups.
	Groups []ArchitectureModelGroup `json:"groups,omitempty"`
	// Constraints is the list of "from -> to" edge constraints.
	Constraints []ArchitectureConstraint `json:"constraints,omitempty"`
}

// ArchitectureModelData is the perspective set that makes up a model.
type ArchitectureModelData struct {
	// Perspectives is the list of model perspectives.
	Perspectives []ArchitectureModelPerspective `json:"perspectives,omitempty"`
}

// ArchitectureModel represents a stored architecture model.
type ArchitectureModel struct {
	// Id is the model's unique identifier.
	Id string `json:"id,omitempty"`
	// ProjectId is the identifier of the project the model belongs to.
	ProjectId string `json:"projectId,omitempty"`
	// OrganizationId is the identifier of the owning organization.
	OrganizationId string `json:"organizationId,omitempty"`
	// Data is an opaque, implementation-defined payload.
	Data string `json:"data,omitempty"`
	// Model is the perspective set that makes up the model.
	Model ArchitectureModelData `json:"model,omitzero"`
}

// ArchitectureModelsListOptions contains parameters for the ListModels method.
type ArchitectureModelsListOptions struct {
	// ProjectId is the project identifier. This field is required.
	ProjectId string `json:"projectId"`
}

// ArchitectureCreateModelOptions contains the request body for the CreateModel method.
type ArchitectureCreateModelOptions struct {
	// Model is the perspective set that makes up the model. Optional.
	Model *ArchitectureModelData `json:"model,omitempty"`
	// ProjectId is the project identifier. This field is required.
	ProjectId string `json:"projectId"`
	// OrganizationId is the identifier of the owning organization. Optional.
	OrganizationId string `json:"organizationId,omitempty"`
	// Data is an opaque, implementation-defined payload. Optional.
	Data string `json:"data,omitempty"`
}

// ArchitectureUpdateModelOptions contains the merge-patch request body for the UpdateModel method.
type ArchitectureUpdateModelOptions struct {
	// Model is the perspective set that makes up the model. Optional.
	Model *ArchitectureModelData `json:"model,omitempty"`
	// Data is an opaque, implementation-defined payload. Optional.
	Data string `json:"data,omitempty"`
}

// ValidateListModelsOpt validates the options for the ListModels method.
func (s *ArchitectureService) ValidateListModelsOpt(opt *ArchitectureModelsListOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.ProjectId, "ProjectId")
}

// ValidateCreateModelOpt validates the options for the CreateModel method.
func (s *ArchitectureService) ValidateCreateModelOpt(opt *ArchitectureCreateModelOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.ProjectId, "ProjectId")
}

// ValidateModelIDOpt validates a model id path parameter.
func (s *ArchitectureService) ValidateModelIDOpt(modelID string) error {
	return ValidateRequired(modelID, "modelID")
}

// ListModels returns the architecture models stored for a project.
//
// API endpoint: GET /api/v2/architecture/models.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *ArchitectureService) ListModels(ctx context.Context, opt *ArchitectureModelsListOptions) ([]ArchitectureModel, *http.Response, error) {
	err := s.ValidateListModelsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/models", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result []ArchitectureModel

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// CreateModel stores a new architecture model for a project.
//
// API endpoint: POST /api/v2/architecture/models.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *ArchitectureService) CreateModel(ctx context.Context, opt *ArchitectureCreateModelOptions) (*ArchitectureModel, *http.Response, error) {
	err := s.ValidateCreateModelOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "architecture/models", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureModel)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// GetModel returns a single architecture model by id.
//
// API endpoint: GET /api/v2/architecture/models/{id}.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *ArchitectureService) GetModel(ctx context.Context, modelID string) (*ArchitectureModel, *http.Response, error) {
	err := s.ValidateModelIDOpt(modelID)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/models/"+url.PathEscape(modelID), nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureModel)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// UpdateModel applies a merge-patch update to an architecture model.
//
// API endpoint: PATCH /api/v2/architecture/models/{id}.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *ArchitectureService) UpdateModel(ctx context.Context, modelID string, opt *ArchitectureUpdateModelOptions) (*ArchitectureModel, *http.Response, error) {
	err := s.ValidateModelIDOpt(modelID)
	if err != nil {
		return nil, nil, err
	}

	if opt == nil {
		return nil, nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPatch, "architecture/models/"+url.PathEscape(modelID), nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureModel)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// DeleteModel removes an architecture model by id.
//
// API endpoint: DELETE /api/v2/architecture/models/{id}.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *ArchitectureService) DeleteModel(ctx context.Context, modelID string) (*http.Response, error) {
	err := s.ValidateModelIDOpt(modelID)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodDelete, "architecture/models/"+url.PathEscape(modelID), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// -----------------------------------------------------------------------------
// Project configurations
// -----------------------------------------------------------------------------

// ArchitectureProjectConfiguration describes the architecture configuration of a project.
type ArchitectureProjectConfiguration struct {
	// ProjectId is the project identifier.
	ProjectId string `json:"projectId,omitempty"`
	// IntendedArchitecture is the raw intended-architecture definition, if any.
	IntendedArchitecture string `json:"intendedArchitecture,omitempty"`
	// BoundaryDescriptors is the list of boundary descriptor keys.
	BoundaryDescriptors []string `json:"boundaryDescriptors,omitempty"`
	// Directives is the list of directive identifiers applied to the project.
	Directives []string `json:"directives,omitempty"`
}

// ArchitectureProjectConfigurationsOptions contains parameters for the
// GetProjectConfigurations and GetPrivateProjectConfigurations methods. At least
// one of ProjectId or ProjectKey should be provided.
type ArchitectureProjectConfigurationsOptions struct {
	// ProjectId filters by project identifier. Optional.
	ProjectId string `json:"projectId,omitempty"`
	// ProjectKey filters by project key. Optional.
	ProjectKey string `json:"projectKey,omitempty"`
}

// GetProjectConfigurations returns architecture configurations for one or more projects.
//
// API endpoint: GET /api/v2/architecture/project-configurations.
// Enterprise Edition only.
func (s *ArchitectureService) GetProjectConfigurations(ctx context.Context, opt *ArchitectureProjectConfigurationsOptions) ([]ArchitectureProjectConfiguration, *http.Response, error) {
	return s.getProjectConfigurations(ctx, "architecture/project-configurations", opt)
}

// GetPrivateProjectConfigurations returns architecture configurations for one or
// more projects via the internal variant of the endpoint.
//
// API endpoint: GET /api/v2/architecture/private/architecture/project-configurations.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *ArchitectureService) GetPrivateProjectConfigurations(ctx context.Context, opt *ArchitectureProjectConfigurationsOptions) ([]ArchitectureProjectConfiguration, *http.Response, error) {
	return s.getProjectConfigurations(ctx, "architecture/private/architecture/project-configurations", opt)
}

// -----------------------------------------------------------------------------
// Snapshots
// -----------------------------------------------------------------------------

// ArchitectureSnapshotMeasures holds the aggregated architecture measures of a snapshot.
type ArchitectureSnapshotMeasures struct {
	// ArchitectureCoverage is the fraction of the graph covered by the intended architecture.
	ArchitectureCoverage string `json:"architectureCoverage,omitempty"`
	// ArchitectureCoverageSummary is a human-readable coverage summary.
	ArchitectureCoverageSummary string `json:"architectureCoverageSummary,omitempty"`
	// NumStructuralDirectives is the number of structural directives.
	NumStructuralDirectives string `json:"numStructuralDirectives,omitempty"`
	// NumRelationshipDirectives is the number of relationship directives.
	NumRelationshipDirectives string `json:"numRelationshipDirectives,omitempty"`
	// NumUnimplemented is the number of unimplemented intended elements.
	NumUnimplemented string `json:"numUnimplemented,omitempty"`
	// NumWeakTangledNodes is the number of weakly tangled nodes.
	NumWeakTangledNodes string `json:"numWeakTangledNodes,omitempty"`
	// NumMisplaced is the number of misplaced nodes.
	NumMisplaced string `json:"numMisplaced,omitempty"`
	// NumTangles is the number of tangles.
	NumTangles string `json:"numTangles,omitempty"`
	// NumSplitNodes is the number of split nodes.
	NumSplitNodes string `json:"numSplitNodes,omitempty"`
	// NumWeakTangles is the number of weak tangles.
	NumWeakTangles string `json:"numWeakTangles,omitempty"`
	// NumGraphNodes is the number of nodes in the graph.
	NumGraphNodes string `json:"numGraphNodes,omitempty"`
	// NumGraphEdges is the number of edges in the graph.
	NumGraphEdges string `json:"numGraphEdges,omitempty"`
	// NumDisallowed is the number of disallowed edges.
	NumDisallowed string `json:"numDisallowed,omitempty"`
	// NumInterfaceIssues is the number of interface issues.
	NumInterfaceIssues string `json:"numInterfaceIssues,omitempty"`
	// MaxWeakTangleSize is the size of the largest weak tangle.
	MaxWeakTangleSize string `json:"maxWeakTangleSize,omitempty"`
	// NumStructuralIssues is the number of structural issues.
	NumStructuralIssues string `json:"numStructuralIssues,omitempty"`
}

// ArchitectureSnapshotEcosystemMeasures pairs a language ecosystem with its measures.
type ArchitectureSnapshotEcosystemMeasures struct {
	// Ecosystem is the language ecosystem. One of "java", "js", "ts", "py", "cs".
	Ecosystem string `json:"ecosystem,omitempty"`
	// Values holds the measures for the ecosystem.
	Values ArchitectureSnapshotMeasures `json:"values,omitzero"`
}

// ArchitectureSnapshot represents a single architecture analysis snapshot.
type ArchitectureSnapshot struct {
	// Id is the snapshot identifier.
	Id string `json:"id,omitempty"`
	// ProjectId is the project identifier.
	ProjectId string `json:"projectId,omitempty"`
	// OrganizationId is the owning organization identifier.
	OrganizationId string `json:"organizationId,omitempty"`
	// BranchId is the branch identifier.
	BranchId string `json:"branchId,omitempty"`
	// AnalysisId is the analysis identifier, if any.
	AnalysisId string `json:"analysisId,omitempty"`
	// AnalysisDate is the analysis date, if any.
	AnalysisDate string `json:"analysisDate,omitempty"`
	// AggregatedMeasures holds the measures aggregated across every ecosystem.
	AggregatedMeasures ArchitectureSnapshotMeasures `json:"aggregatedMeasures,omitzero"`
	// Measures holds the per-ecosystem measures.
	Measures []ArchitectureSnapshotEcosystemMeasures `json:"measures,omitempty"`
}

// ArchitectureSnapshotsList is the paginated response returned by ListSnapshots.
type ArchitectureSnapshotsList struct {
	// Snapshots is the current page of snapshots.
	Snapshots []ArchitectureSnapshot `json:"snapshots,omitempty"`
	// Page contains pagination information.
	Page PageResponseV2 `json:"page,omitzero"`
}

// ArchitectureListSnapshotsOptions contains parameters for the ListSnapshots method.
//
//nolint:govet // fieldalignment: embedded pagination is kept first for readability
type ArchitectureListSnapshotsOptions struct {
	PaginationParamsV2

	// OrganizationId is the owning organization identifier. This field is required.
	OrganizationId string `json:"organizationId"`
	// BranchId is the branch identifier. This field is required.
	BranchId string `json:"branchId"`
}

// ValidateListSnapshotsOpt validates the options for the ListSnapshots method.
func (s *ArchitectureService) ValidateListSnapshotsOpt(opt *ArchitectureListSnapshotsOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.OrganizationId, "OrganizationId")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.BranchId, "BranchId")
}

// ListSnapshots returns the paginated list of architecture snapshots for a branch.
//
// API endpoint: GET /api/v2/architecture/snapshots.
// Enterprise Edition only.
func (s *ArchitectureService) ListSnapshots(ctx context.Context, opt *ArchitectureListSnapshotsOptions) (*ArchitectureSnapshotsList, *http.Response, error) {
	err := s.ValidateListSnapshotsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/snapshots", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureSnapshotsList)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// -----------------------------------------------------------------------------
// Structure
// -----------------------------------------------------------------------------

// ArchitectureStructureEdge is a directed, weighted edge between two structure nodes.
type ArchitectureStructureEdge struct {
	// ToId is the identifier of the destination node.
	ToId int `json:"toId,omitempty"`
	// Weight is the edge weight.
	Weight int `json:"weight,omitempty"`
}

// ArchitectureStructureNode is a node of a structure graph. Nodes nest recursively.
type ArchitectureStructureNode struct {
	// Name is the node name.
	Name string `json:"name,omitempty"`
	// Kind is the node kind.
	Kind string `json:"kind,omitempty"`
	// MappedKind is the node kind mapped to a canonical taxonomy, if any.
	MappedKind string `json:"mappedKind,omitempty"`
	// Edges is the list of outgoing edges.
	Edges []ArchitectureStructureEdge `json:"edges,omitempty"`
	// Nodes is the list of nested child nodes.
	Nodes []ArchitectureStructureNode `json:"nodes,omitempty"`
	// Id is the node identifier, unique within the graph.
	Id int `json:"id,omitempty"`
}

// ArchitectureStructureGraph is a structure graph for one language ecosystem.
type ArchitectureStructureGraph struct {
	// Id is the graph identifier.
	Id string `json:"id,omitempty"`
	// Nodes is the list of top-level nodes.
	Nodes []ArchitectureStructureNode `json:"nodes,omitempty"`
	// NodeCount is the total number of nodes in the graph.
	NodeCount int `json:"nodeCount,omitempty"`
	// EdgeCount is the total number of edges in the graph.
	EdgeCount int `json:"edgeCount,omitempty"`
}

// ArchitectureStructureOptions contains parameters for the GetStructure method.
type ArchitectureStructureOptions struct {
	// OrganizationId is the owning organization identifier. This field is required.
	OrganizationId string `json:"organizationId"`
	// BranchId is the branch identifier. This field is required.
	BranchId string `json:"branchId"`
	// Ecosystem is the language ecosystem. This field is required.
	// One of "java", "js", "ts", "py", "cs".
	Ecosystem string `json:"ecosystem"`
}

// ValidateStructureOpt validates the options for the GetStructure method.
func (s *ArchitectureService) ValidateStructureOpt(opt *ArchitectureStructureOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.OrganizationId, "OrganizationId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.BranchId, "BranchId")
	if err != nil {
		return err
	}

	err = ValidateRequired(opt.Ecosystem, "Ecosystem")
	if err != nil {
		return err
	}

	return IsValueAuthorized(opt.Ecosystem, allowedArchitectureEcosystems, "Ecosystem")
}

// GetStructure returns the structure graph(s) for a branch and language ecosystem.
//
// API endpoint: GET /api/v2/architecture/structure.
// Enterprise Edition only.
func (s *ArchitectureService) GetStructure(ctx context.Context, opt *ArchitectureStructureOptions) ([]ArchitectureStructureGraph, *http.Response, error) {
	err := s.ValidateStructureOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/structure", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result []ArchitectureStructureGraph

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// -----------------------------------------------------------------------------
// Analysis and scanner data
// -----------------------------------------------------------------------------

// ArchitectureStoreAnalysisOptions contains the request body for the StoreAnalysis method.
type ArchitectureStoreAnalysisOptions struct {
	// AnalysisId is the analysis identifier. Optional.
	AnalysisId string `json:"analysisId,omitempty"`
	// CeTaskId is the Compute Engine task identifier. Optional.
	CeTaskId string `json:"ceTaskId,omitempty"`
}

// architectureStoreAnalysisResponse is a defined string type (rather than a bare
// string) used solely to decode the StoreAnalysis response body. As with
// architectureFileGraphResponse above, this keeps client.Do on its default
// JSON-decode path (Accept: application/json) instead of the *string/text/plain
// special case, matching the endpoint's declared "application/json" string payload.
type architectureStoreAnalysisResponse string

// StoreAnalysis stores architecture analysis data and returns the created
// analysis identifier as reported by the server.
//
// API endpoint: POST /api/v2/architecture/analysis.
// Enterprise Edition only.
func (s *ArchitectureService) StoreAnalysis(ctx context.Context, opt *ArchitectureStoreAnalysisOptions) (*string, *http.Response, error) {
	if opt == nil {
		return nil, nil, NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "architecture/analysis", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result architectureStoreAnalysisResponse

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	strResult := string(result)

	return &strResult, resp, nil
}

// ArchitectureScannerDataUpload is returned by InitiateScannerDataUpload and
// carries the URL the scanner should PUT its data to.
type ArchitectureScannerDataUpload struct {
	// Id is the identifier of the pending upload.
	Id string `json:"id,omitempty"`
	// UploadUrl is the URL to upload the scanner data to.
	UploadUrl string `json:"uploadUrl,omitempty"`
}

// ArchitectureScannerDataOptions contains the request body for the InitiateScannerDataUpload method.
type ArchitectureScannerDataOptions struct {
	// AnalysisId is the analysis identifier. This field is required.
	AnalysisId string `json:"analysisId"`
}

// ValidateScannerDataOpt validates the options for the InitiateScannerDataUpload method.
func (s *ArchitectureService) ValidateScannerDataOpt(opt *ArchitectureScannerDataOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.AnalysisId, "AnalysisId")
}

// InitiateScannerDataUpload registers a scanner-data upload for an analysis and
// returns the URL the scanner should upload its report to.
//
// API endpoint: POST /api/v2/architecture/scanner-data.
// Enterprise Edition only.
func (s *ArchitectureService) InitiateScannerDataUpload(ctx context.Context, opt *ArchitectureScannerDataOptions) (*ArchitectureScannerDataUpload, *http.Response, error) {
	err := s.ValidateScannerDataOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodPost, "architecture/scanner-data", nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureScannerDataUpload)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// ArchitectureUploadScannerDataOptions contains parameters for the UploadScannerData method.
type ArchitectureUploadScannerDataOptions struct {
	// CeTaskId is the Compute Engine task identifier. This field is required.
	CeTaskId string
	// Payload is the scanner data report bytes. This field is required.
	Payload io.Reader
	// PayloadFilename is the file name to attach to the multipart part. Optional;
	// defaults to "payload".
	PayloadFilename string
}

// ValidateUploadScannerDataOpt validates the options for the UploadScannerData method.
func (s *ArchitectureService) ValidateUploadScannerDataOpt(opt *ArchitectureUploadScannerDataOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.CeTaskId, "CeTaskId")
	if err != nil {
		return err
	}

	if opt.Payload == nil {
		return NewValidationError("Payload", "is required", ErrMissingRequired)
	}

	return nil
}

// UploadScannerData uploads a scanner data report directly as multipart/form-data.
//
// API endpoint: POST /api/v2/architecture/scanner-data-upload.
// Enterprise Edition only.
func (s *ArchitectureService) UploadScannerData(ctx context.Context, opt *ArchitectureUploadScannerDataOptions) (*http.Response, error) {
	err := s.ValidateUploadScannerDataOpt(opt)
	if err != nil {
		return nil, err
	}

	filename := opt.PayloadFilename
	if filename == "" {
		filename = "payload"
	}

	body, contentType, err := buildMultipartFileBody("payload", filename, opt.Payload)
	if err != nil {
		return nil, err
	}

	//nolint:exhaustruct // only the fields relevant to a multipart upload are set
	req, err := s.client.NewSonarQubeAPIRequest(ctx, SonarAPIRequestParameters{
		Method:   http.MethodPost,
		Path:     v2BasePath + "architecture/scanner-data-upload",
		RawQuery: url.Values{"ceTaskId": {opt.CeTaskId}},
		Headers:  map[string]string{headerContentType: contentType},
		Body:     body,
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

// getProjectConfigurations is the shared implementation behind the public and
// internal project-configuration endpoints.
func (s *ArchitectureService) getProjectConfigurations(ctx context.Context, path string, opt *ArchitectureProjectConfigurationsOptions) ([]ArchitectureProjectConfiguration, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, path, opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	var result []ArchitectureProjectConfiguration

	resp, err := s.client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// buildMultipartFileBody encodes one or more file parts sharing a single form
// field name into a multipart/form-data body. It returns the encoded payload
// together with the Content-Type header (which carries the generated boundary)
// that must accompany it.
func buildMultipartFileBody(field string, filename string, contents ...io.Reader) (io.Reader, string, error) {
	var buf bytes.Buffer

	writer := multipart.NewWriter(&buf)

	for _, content := range contents {
		part, err := writer.CreateFormFile(field, filename)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create multipart part: %w", err)
		}

		_, err = io.Copy(part, content)
		if err != nil {
			return nil, "", fmt.Errorf("failed to write multipart part: %w", err)
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, "", fmt.Errorf("failed to finalize multipart body: %w", err)
	}

	return &buf, writer.FormDataContentType(), nil
}
