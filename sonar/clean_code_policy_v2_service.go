package sonar

// CleanCodePolicyServiceV2 handles communication with the Clean Code Policy
// related methods of the SonarQube V2 API.
type CleanCodePolicyServiceV2 struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}
