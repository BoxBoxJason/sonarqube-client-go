package sonar

// SystemServiceV2 handles communication with the System related methods of the
// SonarQube V2 API.
type SystemServiceV2 struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}
