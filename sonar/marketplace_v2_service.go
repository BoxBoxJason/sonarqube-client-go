package sonar

import (
	"fmt"
	"net/http"
)

// MarketplaceServiceV2 handles communication with the Marketplace related
// methods of the SonarQube V2 API.
type MarketplaceServiceV2 struct {
	// client is used to communicate with the SonarQube API.
	client *Client
}

// -----------------------------------------------------------------------------
// Response Types
// -----------------------------------------------------------------------------

// MarketplaceAzureBillingV2 represents the response from billing an Azure account.
type MarketplaceAzureBillingV2 struct {
	// Message contains an informational message about the billing operation.
	Message string `json:"message,omitempty"`
	// Success indicates whether the billing operation was successful.
	Success bool `json:"success,omitempty"`
}

// -----------------------------------------------------------------------------
// Service Methods
// -----------------------------------------------------------------------------

// BillAzureAccount bills the user's Azure account with the cost of the
// SonarQube Server license. Used by admins.
func (s *MarketplaceServiceV2) BillAzureAccount() (*MarketplaceAzureBillingV2, *http.Response, error) {
	req, err := s.client.NewSonarQubeV2APIRequest(http.MethodPost, "marketplace/azure/billing", nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(MarketplaceAzureBillingV2)

	resp, err := s.client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}
