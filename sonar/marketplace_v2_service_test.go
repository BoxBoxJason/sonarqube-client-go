package sonar

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// BillAzureAccount
// =============================================================================

func TestMarketplaceV2_BillAzureAccount(t *testing.T) {
	response := MarketplaceAzureBillingV2{
		Success: true,
		Message: "Billing successful",
	}
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/marketplace/azure/billing", http.StatusOK, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Marketplace.BillAzureAccount()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, result.Success)
	assert.Equal(t, "Billing successful", result.Message)
}
