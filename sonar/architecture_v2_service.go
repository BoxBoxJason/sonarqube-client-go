package sonar

import (
	"context"
	"fmt"
	"net/http"
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
