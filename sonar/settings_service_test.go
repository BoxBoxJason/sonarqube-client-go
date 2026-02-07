package sonar

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// SettingsService Test Suite
// -----------------------------------------------------------------------------

// TestSettingsService_CheckSecretKey tests the CheckSecretKey method.
func TestSettingsService_CheckSecretKey(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/settings/check_secret_key", http.StatusOK, &SettingsCheckSecretKey{
		SecretKeyAvailable: true,
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)

	result, resp, err := client.Settings.CheckSecretKey()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, result.SecretKeyAvailable)
}

// TestSettingsService_Encrypt tests the Encrypt method.
func TestSettingsService_Encrypt(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/settings/encrypt", http.StatusOK, &SettingsEncrypt{
		EncryptedValue: "{aes-gcm}encrypted-secret",
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &SettingsEncryptOption{
		Value: "my-secret-value",
	}

	result, resp, err := client.Settings.Encrypt(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "{aes-gcm}encrypted-secret", result.EncryptedValue)
}

// TestSettingsService_Encrypt_ValidationError tests validation for Encrypt.
func TestSettingsService_Encrypt_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test missing Value
	opt := &SettingsEncryptOption{}
	_, _, err := client.Settings.Encrypt(opt)
	assert.Error(t, err)
}

// TestSettingsService_GenerateSecretKey tests the GenerateSecretKey method.
func TestSettingsService_GenerateSecretKey(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/settings/generate_secret_key", http.StatusOK, &SettingsGenerateSecretKey{
		SecretKey: "AaBbCcDdEeFfGgHhIiJjKk==",
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)

	result, resp, err := client.Settings.GenerateSecretKey()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "AaBbCcDdEeFfGgHhIiJjKk==", result.SecretKey)
}

// TestSettingsService_ListDefinitions tests the ListDefinitions method.
func TestSettingsService_ListDefinitions(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/settings/list_definitions", http.StatusOK, &SettingsListDefinitions{
		Definitions: []SettingDefinition{
			{
				Key:         "sonar.coverage.jacoco.xmlReportPaths",
				Name:        "XML report paths",
				Description: "Paths to JaCoCo XML coverage reports",
				Type:        "STRING",
				Category:    "Java",
				SubCategory: "Code Coverage",
				MultiValues: true,
			},
		},
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &SettingsListDefinitionsOption{
		Component: "my-project",
	}

	result, resp, err := client.Settings.ListDefinitions(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Definitions, 1)
	assert.Equal(t, "sonar.coverage.jacoco.xmlReportPaths", result.Definitions[0].Key)
	assert.True(t, result.Definitions[0].MultiValues)
}

// TestSettingsService_LoginMessage tests the LoginMessage method.
func TestSettingsService_LoginMessage(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/settings/login_message", http.StatusOK, &SettingsLoginMessage{
		Message: "Welcome to SonarQube!",
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)

	result, resp, err := client.Settings.LoginMessage()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Welcome to SonarQube!", result.Message)
}

// TestSettingsService_Reset tests the Reset method.
func TestSettingsService_Reset(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/settings/reset"))

		err := r.ParseForm()
		require.NoError(t, err)
		assert.Equal(t, "sonar.test.key,sonar.other.key", r.FormValue("keys"))

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &SettingsResetOption{
		Keys: []string{"sonar.test.key", "sonar.other.key"},
	}

	resp, err := client.Settings.Reset(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// TestSettingsService_Reset_ValidationError tests validation for Reset.
func TestSettingsService_Reset_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test missing Keys
	opt := &SettingsResetOption{}
	_, err := client.Settings.Reset(opt)
	assert.Error(t, err)
}

// TestSettingsService_Set tests the Set method.
func TestSettingsService_Set(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/settings/set"))

		err := r.ParseForm()
		require.NoError(t, err)
		assert.Equal(t, "sonar.test.key", r.FormValue("key"))
		assert.Equal(t, "test-value", r.FormValue("value"))

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &SettingsSetOption{
		Key:   "sonar.test.key",
		Value: "test-value",
	}

	resp, err := client.Settings.Set(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// TestSettingsService_Set_ValidationError tests validation for Set.
func TestSettingsService_Set_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test missing Key
	opt := &SettingsSetOption{
		Value: "test-value",
	}
	_, err := client.Settings.Set(opt)
	assert.Error(t, err)

	// Test Value too long
	opt = &SettingsSetOption{
		Key:   "sonar.test.key",
		Value: strings.Repeat("a", MaxSettingValueLength+1),
	}
	_, err = client.Settings.Set(opt)
	assert.Error(t, err)

	// Test missing Value, Values, and FieldValues
	opt = &SettingsSetOption{
		Key: "sonar.test.key",
	}
	_, err = client.Settings.Set(opt)
	assert.Error(t, err)
}

// TestSettingsService_Set_WithMultiValues tests the Set method with multiple values.
func TestSettingsService_Set_WithMultiValues(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/settings/set"))

		err := r.ParseForm()
		require.NoError(t, err)
		assert.Equal(t, "sonar.test.multikey", r.FormValue("key"))
		values := r.Form["values"]
		assert.Len(t, values, 3)
		assert.Equal(t, []string{"value1", "value2", "value3"}, values)

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &SettingsSetOption{
		Key:    "sonar.test.multikey",
		Values: []string{"value1", "value2", "value3"},
	}

	resp, err := client.Settings.Set(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// TestSettingsService_Set_WithFieldValues tests the Set method with field values (property set).
func TestSettingsService_Set_WithFieldValues(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/settings/set"))

		err := r.ParseForm()
		require.NoError(t, err)
		assert.Equal(t, "sonar.test.fieldkey", r.FormValue("key"))

		fieldValuesStr := r.FormValue("fieldValues")
		assert.NotEmpty(t, fieldValuesStr)

		// Verify it's valid JSON
		var decoded map[string]any
		err = json.NewDecoder(strings.NewReader(fieldValuesStr)).Decode(&decoded)
		require.NoError(t, err)

		// Verify the content
		assert.Equal(t, "value1", decoded["field1"])
		assert.Equal(t, "value2", decoded["field2"])

		w.WriteHeader(http.StatusNoContent)
	})
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &SettingsSetOption{
		Key: "sonar.test.fieldkey",
		FieldValues: map[string]any{
			"field1": "value1",
			"field2": "value2",
		},
	}

	resp, err := client.Settings.Set(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

// TestSettingsService_Values tests the Values method.
func TestSettingsService_Values(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/settings/values", http.StatusOK, &SettingsValues{
		Settings: []SettingValue{
			{
				Key:       "sonar.test.key",
				Value:     "test-value",
				Inherited: false,
			},
			{
				Key:       "sonar.multi.key",
				Values:    []string{"value1", "value2"},
				Inherited: true,
			},
		},
		SetSecuredSettings: []string{"sonar.secured.key"},
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &SettingsValuesOption{
		Keys: []string{"sonar.test.key", "sonar.multi.key"},
	}

	result, resp, err := client.Settings.Values(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Settings, 2)
	assert.Equal(t, "sonar.test.key", result.Settings[0].Key)
	assert.Equal(t, "test-value", result.Settings[0].Value)
	assert.Len(t, result.Settings[1].Values, 2)
	assert.Len(t, result.SetSecuredSettings, 1)
}
