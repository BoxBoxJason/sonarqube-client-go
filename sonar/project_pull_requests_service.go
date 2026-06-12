package sonar

import (
	"context"
	"net/http"
)

// ProjectPullRequestsService handles communication with the project pull requests
// related methods of the SonarQube API.
type ProjectPullRequestsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// ProjectPullRequest represents a pull request analysis in SonarQube.
type ProjectPullRequest struct {
	// Key is the pull request key/identifier.
	Key string `json:"key,omitempty"`
	// Title is the pull request title.
	Title string `json:"title,omitempty"`
	// Branch is the source branch of the pull request.
	Branch string `json:"branch,omitempty"`
	// Base is the target branch of the pull request.
	Base string `json:"base,omitempty"`
	// URL is the URL of the pull request in the source code manager.
	URL string `json:"url,omitempty"`
	// Status contains the quality gate status and issue counts for the pull request.
	Status ProjectPullRequestStatus `json:"status,omitzero"`
	// IsOrphan indicates whether the pull request is orphaned (target branch no longer exists).
	IsOrphan bool `json:"isOrphan,omitempty"`
	// AnalysisDate is the date and time when the pull request was last analyzed, in ISO 8601 format.
	AnalysisDate string `json:"analysisDate,omitempty"`
}

// ProjectPullRequestStatus represents the status of a pull request analysis, including quality gate status and issue counts.
type ProjectPullRequestStatus struct {
	// QualityGateStatus is the quality gate status (e.g. "OK", "ERROR").
	QualityGateStatus string `json:"qualityGateStatus,omitempty"`
	// Bugs is the number of bugs found.
	Bugs int `json:"bugs,omitempty"`
	// Vulnerabilities is the number of vulnerabilities found.
	Vulnerabilities int `json:"vulnerabilities,omitempty"`
	// CodeSmells is the number of code smells found.
	CodeSmells int `json:"codeSmells,omitempty"`
}

// ProjectPullRequestsList represents the response from the list endpoint.
type ProjectPullRequestsList struct {
	// PullRequests is the list of pull requests for the project.
	PullRequests []ProjectPullRequest `json:"pullRequests,omitempty"`
}

// -----------------------------------------------------------------------------
// Option Types
// -----------------------------------------------------------------------------

// ProjectPullRequestsDeleteOptions contains parameters for the Delete method.
type ProjectPullRequestsDeleteOptions struct {
	// Project is the project key. This field is required.
	Project string `url:"project"`
	// PullRequest is the pull request key. This field is required.
	PullRequest string `url:"pullRequest"`
}

// ProjectPullRequestsListOptions contains parameters for the List method.
type ProjectPullRequestsListOptions struct {
	// Project is the project key. This field is required.
	Project string `url:"project"`
}

// -----------------------------------------------------------------------------
// Validation Functions
// -----------------------------------------------------------------------------

// ValidateDeleteOpt validates the options for the Delete method.
func (s *ProjectPullRequestsService) ValidateDeleteOpt(opt *ProjectPullRequestsDeleteOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	err := ValidateRequired(opt.Project, "Project")
	if err != nil {
		return err
	}

	return ValidateRequired(opt.PullRequest, "PullRequest")
}

// ValidateListOpt validates the options for the List method.
func (s *ProjectPullRequestsService) ValidateListOpt(opt *ProjectPullRequestsListOptions) error {
	if opt == nil {
		return NewValidationError("opt", "option struct is required", ErrMissingRequired)
	}

	return ValidateRequired(opt.Project, "Project")
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Delete deletes a pull request analysis.
// Requires Administer permission on the project.
//
// API endpoint: POST /api/project_pull_requests/delete.
// Since: 7.1.
func (s *ProjectPullRequestsService) Delete(ctx context.Context, opt *ProjectPullRequestsDeleteOptions) (*http.Response, error) {
	err := s.ValidateDeleteOpt(opt)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, "project_pull_requests/delete", opt)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// List lists the pull request analyses of a project.
// Requires Browse or Execute Analysis permission on the project.
//
// API endpoint: GET /api/project_pull_requests/list.
// Since: 7.1.
func (s *ProjectPullRequestsService) List(ctx context.Context, opt *ProjectPullRequestsListOptions) (*ProjectPullRequestsList, *http.Response, error) {
	err := s.ValidateListOpt(opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "project_pull_requests/list", opt)
	if err != nil {
		return nil, nil, err
	}

	result := new(ProjectPullRequestsList)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
