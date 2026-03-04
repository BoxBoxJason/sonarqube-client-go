package sonar

// AnalysisServiceV2 handles communication with the Analysis related methods of
// the SonarQube V2 API.
type AnalysisServiceV2 struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}
