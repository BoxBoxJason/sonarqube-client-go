package sonargo

import "net/http"

// ServerService handles communication with the Server related methods of the SonarQube API.
type ServerService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// Version returns the SonarQube server version.
func (s *ServerService) Version() (v *string, resp *http.Response, err error) {
	req, err := s.client.NewRequest(http.MethodGet, "server/version", nil)
	if err != nil {
		return nil, nil, err
	}

	v = new(string)

	resp, err = s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, nil
}
