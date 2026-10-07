package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBillingService_GetEntitlementCheck(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/billing/entitlement-checks", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Billing.GetEntitlementCheck(context.Background(), &BillingGetEntitlementCheckOptions{FeatureKey: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestBillingService_GetEntitlementCheck_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Billing.GetEntitlementCheck(context.Background(), nil)
	require.Error(t, err)
}

func TestBillingService_CreateConsumptionRecord(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/billing/consumption-records", http.StatusAccepted, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Billing.CreateConsumptionRecord(context.Background(), "test-value", &BillingCreateConsumptionRecordOptions{FeatureKey: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusAccepted, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestBillingService_CreateConsumptionRecord_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Billing.CreateConsumptionRecord(context.Background(), "", nil)
	require.Error(t, err)
}

func TestBillingService_SetOverage(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/billing/overage", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Billing.SetOverage(context.Background(), &BillingSetOverageOptions{FeatureKey: "test", Enabled: true})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestBillingService_SetOverage_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Billing.SetOverage(context.Background(), nil)
	require.Error(t, err)
}
