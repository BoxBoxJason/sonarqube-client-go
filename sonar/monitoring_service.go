package sonargo

import "net/http"

// MonitoringService handles communication with the monitoring related methods
// of the SonarQube API.
// This service provides monitoring metrics in Prometheus format.
type MonitoringService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// MonitoringMetrics represents the response from the monitoring metrics endpoint.
// The content is a string containing metrics in Prometheus format.
type MonitoringMetrics string

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// Metrics returns monitoring metrics in Prometheus format.
// Supports content type 'text/plain' (default) and 'application/openmetrics-text'.
// This endpoint can be accessed using a Bearer token, which needs to be defined
// in sonar.properties with the 'sonar.web.systemPasscode' key.
//
// API endpoint: GET /api/monitoring/metrics.
// Since: 9.3.
func (s *MonitoringService) Metrics() (*MonitoringMetrics, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "monitoring/metrics", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(MonitoringMetrics)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
