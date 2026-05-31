package sonar

import (
	"context"
	"fmt"
	"net/http"
)

// ArchitectureService handles communication with the architecture V2 API endpoints.
// This service is only available in Enterprise Edition.
type ArchitectureService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// ArchitectureFileGraph represents the file dependency graph for a project branch.
type ArchitectureFileGraph struct {
	// Nodes is the list of file nodes in the dependency graph.
	Nodes []ArchitectureFileNode `json:"nodes,omitempty"`
	// Edges is the list of dependency edges between files.
	Edges []ArchitectureFileEdge `json:"edges,omitempty"`
}

// ArchitectureFileNode represents a file node in the architecture graph.
type ArchitectureFileNode struct {
	// Id is the unique identifier of the file node.
	Id string `json:"id,omitempty"`
	// Name is the display name of the file.
	Name string `json:"name,omitempty"`
	// Path is the file path within the project.
	Path string `json:"path,omitempty"`
}

// ArchitectureFileEdge represents a dependency edge between two file nodes.
type ArchitectureFileEdge struct {
	// Source is the id of the source file node.
	Source string `json:"source,omitempty"`
	// Target is the id of the target file node.
	Target string `json:"target,omitempty"`
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
	// Source is the source file path to start the graph from. This field is required.
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

// FileGraph returns the file dependency graph for a project branch starting from a source file.
// Requires 'Browse' permission on the project.
//
// API endpoint: GET /api/v2/architecture/file-graph.
// Enterprise Edition only.
func (s *ArchitectureService) FileGraph(ctx context.Context, opt *ArchitectureFileGraphOptions) (*ArchitectureFileGraph, *http.Response, error) {
	err := s.ValidateFileGraphOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "architecture/file-graph", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ArchitectureFileGraph)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
