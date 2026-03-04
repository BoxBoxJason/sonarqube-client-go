package sonar

// UsersManagementServiceV2 handles communication with the Users Management
// related methods of the SonarQube V2 API.
type UsersManagementServiceV2 struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}
