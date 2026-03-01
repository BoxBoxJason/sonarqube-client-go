package sonar

// AuthorizationsServiceV2 handles communication with the Authorizations related
// methods of the SonarQube V2 API.
type AuthorizationsServiceV2 struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}
