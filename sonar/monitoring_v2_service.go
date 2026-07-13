package sonar

import (
	"context"
	"fmt"
	"net/http"
)

// MonitoringServiceV2 handles communication with the monitoring V2 API endpoints.
// This service is only available in Enterprise Edition. The underlying endpoint is
// marked internal by SonarQube (x-sonar-internal) and its request/response contract
// may change without notice between SonarQube versions.
type MonitoringServiceV2 struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// MonitoringAlert represents a single admin alert whose condition is currently active.
type MonitoringAlert struct {
	// Key is the unique identifier of the alert.
	Key string `json:"key,omitempty"`
	// Message is a human-readable message describing the alert condition shown to administrators.
	Message string `json:"message,omitempty"`
	// ActiveSince is the timestamp the alert was first triggered, in ISO-8601 format.
	ActiveSince string `json:"activeSince,omitempty"`
}

// MonitoringActiveAlerts represents the response from listing active admin alerts.
type MonitoringActiveAlerts struct {
	// Alerts is the list of currently active admin alerts.
	Alerts []MonitoringAlert `json:"alerts,omitempty"`
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// GetActiveAlerts returns the list of admin alerts whose condition is currently
// active. Each entry contains the alert identifier and a message describing the
// condition.
// Requires 'Administer System' permission.
//
// API endpoint: GET /api/v2/monitoring/alerts.
// Enterprise Edition only. Marked internal by SonarQube and subject to change
// without notice.
func (s *MonitoringServiceV2) GetActiveAlerts(ctx context.Context) (*MonitoringActiveAlerts, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(ctx, http.MethodGet, "monitoring/alerts", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(MonitoringActiveAlerts)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
