package sonar

import (
	"context"
	"fmt"
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
