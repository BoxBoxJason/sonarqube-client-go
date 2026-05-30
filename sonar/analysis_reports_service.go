package sonar

import (
	"context"
	"net/http"
)

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

// QueueStatus checks if the queue of Compute Engine is empty.
// Returns an AnalysisReportsQueueStatus indicating whether the queue is empty.
//
// API endpoint: GET /api/analysis_reports/is_queue_empty.
// WARNING: this is an internal API and may change without notice.
func (s *AnalysisReportsService) QueueStatus(ctx context.Context) (*AnalysisReportsQueueStatus, *http.Response, error) {
	req, err := s.client.NewSonarQubeV1APIRequest(ctx, http.MethodGet, "analysis_reports/is_queue_empty", nil)
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
