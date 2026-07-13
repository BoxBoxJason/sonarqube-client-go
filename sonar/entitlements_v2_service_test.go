package sonar

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// ActivateOnline
// =============================================================================

func TestEntitlementsV2_ActivateOnline(t *testing.T) {
	request := &EntitlementsActivateOnlineOptions{LicenseKey: "third-party-license-key"}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/entitlements/online-activation", http.StatusOK, request, nil))
	client := newTestClient(t, server.url())

	resp, err := client.V2.Entitlements.ActivateOnline(context.Background(), request)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestEntitlementsV2_ActivateOnline_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		opt  *EntitlementsActivateOnlineOptions
		name string
	}{
		{nil, "nil opt"},
		{&EntitlementsActivateOnlineOptions{}, "missing license key"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.V2.Entitlements.ActivateOnline(context.Background(), tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// DeactivateOffline
// =============================================================================

func TestEntitlementsV2_DeactivateOffline(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v2/entitlements/offline-deactivation", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode("deactivation-file-content")
	})
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Entitlements.DeactivateOffline(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "deactivation-file-content", *result)
}

// =============================================================================
// GetOfflineActivationRequest
// =============================================================================

func TestEntitlementsV2_GetOfflineActivationRequest(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v2/entitlements/offline-activation", r.URL.Path)
		assert.Equal(t, "ABCD-EFGH-IJKL-MNOP", r.Header.Get("License-Key"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode("request-file-content")
	})
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Entitlements.GetOfflineActivationRequest(context.Background(), &EntitlementsGetOfflineActivationRequestOptions{
		LicenseKey: "ABCD-EFGH-IJKL-MNOP",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "request-file-content", *result)
}

func TestEntitlementsV2_GetOfflineActivationRequest_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		opt  *EntitlementsGetOfflineActivationRequestOptions
		name string
	}{
		{nil, "nil opt"},
		{&EntitlementsGetOfflineActivationRequestOptions{}, "missing license key"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.V2.Entitlements.GetOfflineActivationRequest(context.Background(), tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// ActivateOffline
// =============================================================================

func TestEntitlementsV2_ActivateOffline(t *testing.T) {
	request := &EntitlementsActivateOfflineOptions{
		License:    "contents of the .lic file",
		LicenseKey: "ABCD-EFGH-IJKL-MNOP",
	}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/entitlements/offline-activation", http.StatusOK, request, nil))
	client := newTestClient(t, server.url())

	resp, err := client.V2.Entitlements.ActivateOffline(context.Background(), request)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestEntitlementsV2_ActivateOffline_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		opt  *EntitlementsActivateOfflineOptions
		name string
	}{
		{nil, "nil opt"},
		{&EntitlementsActivateOfflineOptions{LicenseKey: "key"}, "missing license"},
		{&EntitlementsActivateOfflineOptions{License: "content"}, "missing license key"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.V2.Entitlements.ActivateOffline(context.Background(), tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// ActivateLegacy
// =============================================================================

func TestEntitlementsV2_ActivateLegacy(t *testing.T) {
	request := &EntitlementsActivateLegacyOptions{LicenseKey: "sonar-issued-license-key"}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/entitlements/legacy-activation", http.StatusOK, request, nil))
	client := newTestClient(t, server.url())

	resp, err := client.V2.Entitlements.ActivateLegacy(context.Background(), request)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestEntitlementsV2_ActivateLegacy_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		opt  *EntitlementsActivateLegacyOptions
		name string
	}{
		{nil, "nil opt"},
		{&EntitlementsActivateLegacyOptions{}, "missing license key"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.V2.Entitlements.ActivateLegacy(context.Background(), tt.opt)
			assert.Error(t, err)
		})
	}
}

// =============================================================================
// GetLicense
// =============================================================================

func TestEntitlementsV2_GetLicense(t *testing.T) {
	response := LicenseV2{
		Edition:      "enterprise",
		Type:         "PRODUCTION",
		ServerId:     "server-1",
		LicenseKey:   "ABCD-EFGH-IJKL-MNOP",
		Supported:    true,
		Expired:      false,
		ValidEdition: true,
		Features: []LicenseFeatureV2{
			{Name: "governance", StartDate: "2024-01-01"},
		},
		MaxLoc: 1000000,
		Loc:    50000,
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/entitlements/license", http.StatusOK, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Entitlements.GetLicense(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "enterprise", result.Edition)
	assert.Equal(t, "PRODUCTION", result.Type)
	assert.True(t, result.Supported)
	assert.False(t, result.Expired)
	assert.Len(t, result.Features, 1)
	assert.Equal(t, "governance", result.Features[0].Name)
}

// =============================================================================
// DeleteLicense
// =============================================================================

func TestEntitlementsV2_DeleteLicense(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodDelete, "/v2/entitlements/license", http.StatusOK))
	client := newTestClient(t, server.url())

	resp, err := client.V2.Entitlements.DeleteLicense(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// =============================================================================
// UpdateLicense
// =============================================================================

func TestEntitlementsV2_UpdateLicense(t *testing.T) {
	server := newTestServer(t, mockPatchHandler(t, "/v2/entitlements/license", http.StatusOK, nil, nil))
	client := newTestClient(t, server.url())

	resp, err := client.V2.Entitlements.UpdateLicense(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// =============================================================================
// GetPurchasableFeatures
// =============================================================================

func TestEntitlementsV2_GetPurchasableFeatures(t *testing.T) {
	response := []PurchasableFeatureV2{
		{FeatureKey: "governance", IsEnabled: true, IsAvailable: true, URL: "https://example.com/governance"},
		{FeatureKey: "security-reports", Parent: "governance", IsEnabled: false, IsAvailable: true},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/entitlements/purchasable-features", http.StatusOK, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Entitlements.GetPurchasableFeatures(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result, 2)
	assert.Equal(t, "governance", result[0].FeatureKey)
	assert.True(t, result[0].IsEnabled)
	assert.Equal(t, "governance", result[1].Parent)
	assert.False(t, result[1].IsEnabled)
}
