package sonar

// MarketplaceServiceV2 handles communication with the Marketplace related
// methods of the SonarQube V2 API.
type MarketplaceServiceV2 struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}
