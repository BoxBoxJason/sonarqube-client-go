package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEditionsService_ActivateGracePeriod(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/editions/activate_grace_period", http.StatusNoContent))

	client := newTestClient(t, server.URL)

	resp, err := client.Editions.ActivateGracePeriod(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestEditionsService_Get(t *testing.T) {
	response := LicenseGet{
		License: License{
			ContactEmail:           "admin@example.com",
			ExpiresAt:              "2030-01-01",
			IsExpired:              false,
			IsOfficialDistribution: true,
			IsSupported:            true,
			IsValidEdition:         true,
			IsValidServerId:        true,
			LOCsMax:                1000000,
			LOCsRemaining:          750000,
			Organization:           "Example Corp",
			ProductEdition:         "ENTERPRISE",
			ServerId:               "server-id-123",
			Type:                   "PRODUCTION",
		},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/editions/show_license", http.StatusOK, response))

	client := newTestClient(t, server.URL)

	result, resp, err := client.Editions.Get(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "admin@example.com", result.License.ContactEmail)
	assert.Equal(t, "Example Corp", result.License.Organization)
	assert.Equal(t, "ENTERPRISE", result.License.ProductEdition)
	assert.True(t, result.License.IsSupported)
	assert.Equal(t, int64(1000000), result.License.LOCsMax)
	assert.Equal(t, int64(750000), result.License.LOCsRemaining)
}

func TestEditionsService_Get_Empty(t *testing.T) {
	response := LicenseGet{}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/editions/show_license", http.StatusOK, response))

	client := newTestClient(t, server.URL)

	result, resp, err := client.Editions.Get(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
}

func TestEditionsService_IsValidLicense(t *testing.T) {
	response := LicenseIsValid{IsValid: true}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/editions/is_valid_license", http.StatusOK, response))

	client := newTestClient(t, server.URL)

	result, resp, err := client.Editions.IsValidLicense(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.True(t, result.IsValid)
}

func TestEditionsService_IsValidLicense_Invalid(t *testing.T) {
	response := LicenseIsValid{IsValid: false}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/editions/is_valid_license", http.StatusOK, response))

	client := newTestClient(t, server.URL)

	result, resp, err := client.Editions.IsValidLicense(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.False(t, result.IsValid)
}

func TestEditionsService_Set(t *testing.T) {
	server := newTestServer(t, mockEmptyHandlerWithParams(t, http.MethodPost, "/editions/set_license", http.StatusNoContent, map[string]string{"license": "my-license-key"}))

	client := newTestClient(t, server.URL)

	resp, err := client.Editions.Set(context.Background(), &LicenseSetOptions{
		License: "my-license-key",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestEditionsService_Set_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	resp, err := client.Editions.Set(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	// Missing License key should fail validation.
	resp, err = client.Editions.Set(context.Background(), &LicenseSetOptions{})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestEditionsService_UnsetLicense(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/editions/unset_license", http.StatusNoContent))

	client := newTestClient(t, server.URL)

	resp, err := client.Editions.UnsetLicense(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestEditionsService_ValidateSetOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Valid option should pass.
	err := client.Editions.ValidateSetOpt(&LicenseSetOptions{
		License: "my-license-key",
	})
	assert.NoError(t, err)

	// Nil option should fail.
	err = client.Editions.ValidateSetOpt(nil)
	assert.Error(t, err)

	// Missing License key should fail.
	err = client.Editions.ValidateSetOpt(&LicenseSetOptions{})
	assert.Error(t, err)
}
