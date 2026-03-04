package sonar

// DopTranslationServiceV2 handles communication with the DevOps Platform
// Translation related methods of the SonarQube V2 API.
type DopTranslationServiceV2 struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}
