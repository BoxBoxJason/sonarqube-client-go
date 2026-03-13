package sonar

// AuthorizationsService handles communication with the Authorizations related
// methods of the SonarQube V2 API.
type AuthorizationsService struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}
