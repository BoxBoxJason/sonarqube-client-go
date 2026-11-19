package sonar

import (
	"context"
	"fmt"
	"net/http"
)

// ComponentsServiceV2 handles communication with the Components related
// methods of the SonarQube V2 API.
type ComponentsServiceV2 struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// ComponentsProjectLinkV2 represents a project link returned by the V2 components API.
type ComponentsProjectLinkV2 struct {
	// Name is the link name. It may be null on the server side.
	Name string `json:"name,omitempty"`
	// Type is the link type (e.g. home, ci, issue, scm).
	Type string `json:"type,omitempty"`
	// Url is the link URL.
	Url string `json:"url,omitempty"`
}

// ComponentsProjectV2 represents a project returned by the V2 components API.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type ComponentsProjectV2 struct {
	// CreatedAt is the creation date of the project (ISO 8601).
	CreatedAt string `json:"createdAt,omitempty"`
	// Description is the project description. It may be null on the server side.
	Description string `json:"description,omitempty"`
	// Id is the project UUID.
	Id string `json:"id,omitempty"`
	// Key is the project key.
	Key string `json:"key,omitempty"`
	// Links is the list of project links.
	Links []ComponentsProjectLinkV2 `json:"links,omitempty"`
	// Name is the project name.
	Name string `json:"name,omitempty"`
	// Tags is the list of project tags.
	Tags []string `json:"tags,omitempty"`
	// UpdatedAt is the last update date of the project (ISO 8601).
	UpdatedAt string `json:"updatedAt,omitempty"`
	// Visibility is the project visibility (public or private).
	Visibility string `json:"visibility,omitempty"`
}

// ComponentsProjectsV2 represents the response from the GetProjects method.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type ComponentsProjectsV2 struct {
	// Page contains pagination information.
	Page PageResponseV2 `json:"page,omitzero"`
	// Projects is the list of matching projects.
	Projects []ComponentsProjectV2 `json:"projects,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types (Query Parameters)
// -----------------------------------------------------------------------------

// ComponentsGetProjectsOptionsV2 contains query parameters for the GetProjects method.
// At least one of Ids or Keys is required.
//
//nolint:govet // Field alignment less important than maintaining consistent field order for readability
type ComponentsGetProjectsOptionsV2 struct {
	PaginationParamsV2

	// Ids filters projects by UUID.
	Ids []string `json:"ids,omitempty"`
	// Keys filters projects by key.
	Keys []string `json:"keys,omitempty"`
	// Tags matches projects having any of the supplied tags.
	Tags []string `json:"tags,omitempty"`
}

// -----------------------------------------------------------------------------
// Validation
// -----------------------------------------------------------------------------

// ValidateGetProjectsOpt validates the ComponentsGetProjectsOptionsV2.
func (s *ComponentsServiceV2) ValidateGetProjectsOpt(opt *ComponentsGetProjectsOptionsV2) error {
	if opt == nil {
		return NewValidationError("opt", "must not be nil", ErrMissingRequired)
	}

	if len(opt.Ids) == 0 && len(opt.Keys) == 0 {
		return NewValidationError("Ids", "at least one of Ids or Keys is required", ErrMissingRequired)
	}

	return opt.Validate()
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// GetProjects returns the projects matching the given ids, keys and tags.
// Only projects on which the caller has the Browse permission are returned.
// Results are ordered by creation time, then project ID.
func (s *ComponentsServiceV2) GetProjects(ctx context.Context, opt *ComponentsGetProjectsOptionsV2) (*ComponentsProjectsV2, *http.Response, error) {
	err := s.ValidateGetProjectsOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "components/projects", opt, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(ComponentsProjectsV2)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
