package sonargo

import "net/http"

// AnalysisReportsService handles communication with the analysis reports related methods
// of the SonarQube API.
// This service provides information about Compute Engine tasks.
type AnalysisReportsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// AnalysisReportsQueueStatus represents the response from checking if the Compute Engine queue is empty.
type AnalysisReportsQueueStatus struct {
	// IsEmpty indicates whether the Compute Engine queue is empty.
	IsEmpty bool
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// IsQueueEmpty checks if the queue of Compute Engine is empty.
// Returns true if the queue is empty, false otherwise.
//
// API endpoint: GET /api/analysis_reports/is_queue_empty.
// WARNING: this is an internal API and may change without notice.
func (s *AnalysisReportsService) IsQueueEmpty() (*AnalysisReportsQueueStatus, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "analysis_reports/is_queue_empty", nil)
	if err != nil {
		return nil, nil, err
	}

	// The API returns a plain text "true" or "false"
	var rawResponse string

	resp, err := s.client.Do(req, &rawResponse)
	if err != nil {
		return nil, resp, err
	}

	return &AnalysisReportsQueueStatus{
		IsEmpty: rawResponse == "true",
	}, resp, nil
}
