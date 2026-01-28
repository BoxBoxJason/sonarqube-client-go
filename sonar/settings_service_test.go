package sonargo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// -----------------------------------------------------------------------------
// SettingsService Test Suite
// -----------------------------------------------------------------------------

// TestSettingsService_CheckSecretKey tests the CheckSecretKey method.
func TestSettingsService_CheckSecretKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/settings/check_secret_key") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"secretKeyAvailable":true}`))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	result, resp, err := client.Settings.CheckSecretKey()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if !result.SecretKeyAvailable {
		t.Errorf("expected secretKeyAvailable to be true")
	}
}

// TestSettingsService_Encrypt tests the Encrypt method.
func TestSettingsService_Encrypt(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/settings/encrypt") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"encryptedValue":"{aes-gcm}encrypted-secret"}`))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &SettingsEncryptOption{
		Value: "my-secret-value",
	}

	result, resp, err := client.Settings.Encrypt(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if result.EncryptedValue != "{aes-gcm}encrypted-secret" {
		t.Errorf("unexpected encrypted value: %s", result.EncryptedValue)
	}
}

// TestSettingsService_Encrypt_ValidationError tests validation for Encrypt.
func TestSettingsService_Encrypt_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test missing Value
	opt := &SettingsEncryptOption{}
	_, _, err := client.Settings.Encrypt(opt)
	if err == nil {
		t.Error("expected validation error for missing Value")
	}
}

// TestSettingsService_GenerateSecretKey tests the GenerateSecretKey method.
func TestSettingsService_GenerateSecretKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/settings/generate_secret_key") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"secretKey":"AaBbCcDdEeFfGgHhIiJjKk=="}`))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	result, resp, err := client.Settings.GenerateSecretKey()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if result.SecretKey != "AaBbCcDdEeFfGgHhIiJjKk==" {
		t.Errorf("unexpected secret key: %s", result.SecretKey)
	}
}

// TestSettingsService_ListDefinitions tests the ListDefinitions method.
func TestSettingsService_ListDefinitions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/settings/list_definitions") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"definitions": [
				{
					"key": "sonar.coverage.jacoco.xmlReportPaths",
					"name": "XML report paths",
					"description": "Paths to JaCoCo XML coverage reports",
					"type": "STRING",
					"category": "Java",
					"subCategory": "Code Coverage",
					"multiValues": true
				}
			]
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &SettingsListDefinitionsOption{
		Component: "my-project",
	}

	result, resp, err := client.Settings.ListDefinitions(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if len(result.Definitions) != 1 {
		t.Errorf("expected 1 definition, got %d", len(result.Definitions))
	}
	if result.Definitions[0].Key != "sonar.coverage.jacoco.xmlReportPaths" {
		t.Errorf("unexpected definition key: %s", result.Definitions[0].Key)
	}
	if !result.Definitions[0].MultiValues {
		t.Errorf("expected multiValues to be true")
	}
}

// TestSettingsService_LoginMessage tests the LoginMessage method.
func TestSettingsService_LoginMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/settings/login_message") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"Welcome to SonarQube!"}`))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	result, resp, err := client.Settings.LoginMessage()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if result.Message != "Welcome to SonarQube!" {
		t.Errorf("unexpected message: %s", result.Message)
	}
}

// TestSettingsService_Reset tests the Reset method.
func TestSettingsService_Reset(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/settings/reset") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("failed to parse form: %v", err)
		}
		if r.FormValue("keys") != "sonar.test.key,sonar.other.key" {
			t.Errorf("unexpected keys: %s", r.FormValue("keys"))
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &SettingsResetOption{
		Keys: []string{"sonar.test.key", "sonar.other.key"},
	}

	resp, err := client.Settings.Reset(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

// TestSettingsService_Reset_ValidationError tests validation for Reset.
func TestSettingsService_Reset_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test missing Keys
	opt := &SettingsResetOption{}
	_, err := client.Settings.Reset(opt)
	if err == nil {
		t.Error("expected validation error for missing Keys")
	}
}

// TestSettingsService_Set tests the Set method.
func TestSettingsService_Set(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/settings/set") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("failed to parse form: %v", err)
		}
		if r.FormValue("key") != "sonar.test.key" {
			t.Errorf("unexpected key: %s", r.FormValue("key"))
		}
		if r.FormValue("value") != "test-value" {
			t.Errorf("unexpected value: %s", r.FormValue("value"))
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &SettingsSetOption{
		Key:   "sonar.test.key",
		Value: "test-value",
	}

	resp, err := client.Settings.Set(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

// TestSettingsService_Set_ValidationError tests validation for Set.
func TestSettingsService_Set_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Test missing Key
	opt := &SettingsSetOption{
		Value: "test-value",
	}
	_, err := client.Settings.Set(opt)
	if err == nil {
		t.Error("expected validation error for missing Key")
	}

	// Test Value too long
	opt = &SettingsSetOption{
		Key:   "sonar.test.key",
		Value: strings.Repeat("a", MaxSettingValueLength+1),
	}
	_, err = client.Settings.Set(opt)
	if err == nil {
		t.Error("expected validation error for Value too long")
	}

	// Test missing Value, Values, and FieldValues
	opt = &SettingsSetOption{
		Key: "sonar.test.key",
	}
	_, err = client.Settings.Set(opt)
	if err == nil {
		t.Error("expected validation error when all of Value, Values, and FieldValues are empty")
	}
}

// TestSettingsService_Set_WithMultiValues tests the Set method with multiple values.
func TestSettingsService_Set_WithMultiValues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/settings/set") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("failed to parse form: %v", err)
		}
		if r.FormValue("key") != "sonar.test.multikey" {
			t.Errorf("unexpected key: %s", r.FormValue("key"))
		}
		values := r.Form["values"]
		if len(values) != 3 {
			t.Errorf("expected 3 values, got %d", len(values))
		}
		if values[0] != "value1" || values[1] != "value2" || values[2] != "value3" {
			t.Errorf("unexpected values: %v", values)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &SettingsSetOption{
		Key:    "sonar.test.multikey",
		Values: []string{"value1", "value2", "value3"},
	}

	resp, err := client.Settings.Set(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

// TestSettingsService_Set_WithFieldValues tests the Set method with field values (property set).
func TestSettingsService_Set_WithFieldValues(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/settings/set") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("failed to parse form: %v", err)
		}
		if r.FormValue("key") != "sonar.test.fieldkey" {
			t.Errorf("unexpected key: %s", r.FormValue("key"))
		}

		fieldValuesStr := r.FormValue("fieldValues")
		if fieldValuesStr == "" {
			t.Error("expected fieldValues to be present")
		}

		// Verify it's valid JSON
		var decoded map[string]any
		if err := json.NewDecoder(strings.NewReader(fieldValuesStr)).Decode(&decoded); err != nil {
			t.Errorf("failed to decode fieldValues JSON: %v", err)
		}

		// Verify the content
		if decoded["field1"] != "value1" {
			t.Errorf("unexpected field1 value: %v", decoded["field1"])
		}
		if decoded["field2"] != "value2" {
			t.Errorf("unexpected field2 value: %v", decoded["field2"])
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &SettingsSetOption{
		Key: "sonar.test.fieldkey",
		FieldValues: map[string]any{
			"field1": "value1",
			"field2": "value2",
		},
	}

	resp, err := client.Settings.Set(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

// TestSettingsService_Values tests the Values method.
func TestSettingsService_Values(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		if !strings.HasSuffix(r.URL.Path, "/settings/values") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"settings": [
				{
					"key": "sonar.test.key",
					"value": "test-value",
					"inherited": false
				},
				{
					"key": "sonar.multi.key",
					"values": ["value1", "value2"],
					"inherited": true
				}
			],
			"setSecuredSettings": ["sonar.secured.key"]
		}`))
	}))
	defer server.Close()

	client, _ := NewClient(server.URL+"/api/", "user", "pass")

	opt := &SettingsValuesOption{
		Keys: []string{"sonar.test.key", "sonar.multi.key"},
	}

	result, resp, err := client.Settings.Values(opt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
	if len(result.Settings) != 2 {
		t.Errorf("expected 2 settings, got %d", len(result.Settings))
	}
	if result.Settings[0].Key != "sonar.test.key" {
		t.Errorf("unexpected setting key: %s", result.Settings[0].Key)
	}
	if result.Settings[0].Value != "test-value" {
		t.Errorf("unexpected setting value: %s", result.Settings[0].Value)
	}
	if len(result.Settings[1].Values) != 2 {
		t.Errorf("expected 2 values, got %d", len(result.Settings[1].Values))
	}
	if len(result.SetSecuredSettings) != 1 {
		t.Errorf("expected 1 secured setting, got %d", len(result.SetSecuredSettings))
	}
}
