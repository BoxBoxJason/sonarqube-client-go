package sonar

import "net/http"

// FeaturesService handles communication with the features related methods
// of the SonarQube API.
// This service provides information about features available in SonarQube.
type FeaturesService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// FeaturesList represents the list of supported features returned by the API.
// This is a simple slice of feature names.
type FeaturesList []string

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// List returns the list of supported features in the SonarQube instance.
//
// API endpoint: GET /api/features/list.
// WARNING: This is an internal API and may change without notice.
func (s *FeaturesService) List() (*FeaturesList, *http.Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "features/list", nil)
	if err != nil {
		return nil, nil, err
	}

	result := new(FeaturesList)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
