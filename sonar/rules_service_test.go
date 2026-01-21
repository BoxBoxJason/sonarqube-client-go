package sonargo

import (
"encoding/json"
"net/http"
"net/http/httptest"
"testing"
)

func TestRules_App(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
if r.Method != "GET" {
t.Errorf("expected method GET, got %s", r.Method)
}

w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"canWrite":true,"languages":{"java":"Java"}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Rules.App()
	if err != nil {
		t.Fatalf("App failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || !result.CanWrite {
		t.Error("expected CanWrite to be true")
	}
}

func TestRules_Create(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
if r.Method != "POST" {
t.Errorf("expected method POST, got %s", r.Method)
}

w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"rule":{"key":"java:MyRule","name":"My Rule"}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &RulesCreateOption{CustomKey: "MyRule", Name: "My Rule"}
	result, resp, err := client.Rules.Create(opt)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || result.Rule.Key != "java:MyRule" {
		t.Error("unexpected rule key")
	}
}

func TestRules_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
if r.Method != "POST" {
t.Errorf("expected method POST, got %s", r.Method)
}

w.WriteHeader(204)
}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &RulesDeleteOption{Key: "java:MyRule"}
	resp, err := client.Rules.Delete(opt)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestRules_Search(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
if r.Method != "GET" {
t.Errorf("expected method GET, got %s", r.Method)
}

w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"paging":{"total":0},"rules":[],"actives":{}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &RulesSearchOption{Languages: "java"}
	result, resp, err := client.Rules.Search(opt)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

// TestRulesSearchResponse_DynamicActives verifies that the Actives field
// correctly handles dynamic rule keys as a map[string][]RuleActivation.
func TestRulesSearchResponse_DynamicActives(t *testing.T) {
	jsonData := `{
		"actives": {
			"squid:S1067": [
				{"qProfile": "profile1", "severity": "MAJOR"}
			],
			"squid:ClassCyclomaticComplexity": [
				{"qProfile": "profile2", "severity": "CRITICAL"}
			],
			"custom:MyRule": [
				{"qProfile": "profile3", "severity": "INFO"}
			]
		},
		"paging": {"pageIndex": 1, "pageSize": 10, "total": 3},
		"rules": []
	}`

	var response RulesSearchResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(response.Actives) != 3 {
		t.Errorf("expected 3 active rule keys, got %d", len(response.Actives))
	}

	// Verify dynamic keys exist
	if _, exists := response.Actives["squid:S1067"]; !exists {
		t.Error("expected 'squid:S1067' in actives")
	}

	if _, exists := response.Actives["squid:ClassCyclomaticComplexity"]; !exists {
		t.Error("expected 'squid:ClassCyclomaticComplexity' in actives")
	}

	if _, exists := response.Actives["custom:MyRule"]; !exists {
		t.Error("expected 'custom:MyRule' in actives")
	}

	// Verify activation details
	if len(response.Actives["squid:S1067"]) != 1 {
		t.Error("expected 1 activation for squid:S1067")
	}

	if response.Actives["squid:S1067"][0].QProfile != "profile1" {
		t.Errorf("expected qProfile 'profile1', got '%s'",
response.Actives["squid:S1067"][0].QProfile)
	}
}

func TestRules_Show(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
if r.Method != "GET" {
t.Errorf("expected method GET, got %s", r.Method)
}

w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"rule":{"key":"java:S1067","name":"Test Rule"}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &RulesShowOption{Key: "java:S1067"}
	result, resp, err := client.Rules.Show(opt)
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || result.Rule.Key != "java:S1067" {
		t.Error("unexpected rule key")
	}
}

func TestRules_Tags(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
if r.Method != "GET" {
t.Errorf("expected method GET, got %s", r.Method)
}

w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"tags":["security","bug"]}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &RulesTagsOption{Ps: "100"}
	result, resp, err := client.Rules.Tags(opt)
	if err != nil {
		t.Fatalf("Tags failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.Tags) != 2 {
		t.Error("expected 2 tags")
	}
}

func TestRules_Update(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
if r.Method != "POST" {
t.Errorf("expected method POST, got %s", r.Method)
}

w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"rule":{"key":"java:MyRule","name":"Updated Rule"}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &RulesUpdateOption{Key: "java:MyRule", Name: "Updated Rule"}
	result, resp, err := client.Rules.Update(opt)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || result.Rule.Name != "Updated Rule" {
		t.Error("unexpected rule name")
	}
}

func TestRules_Repositories(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
if r.Method != "GET" {
t.Errorf("expected method GET, got %s", r.Method)
}

w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"repositories":[{"key":"java","language":"java","name":"SonarAnalyzer"}]}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &RulesRepositoriesOption{}
	result, resp, err := client.Rules.Repositories(opt)
	if err != nil {
		t.Fatalf("Repositories failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.Repositories) != 1 {
		t.Error("expected 1 repository")
	}
}

func TestRules_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
if r.Method != "GET" {
t.Errorf("expected method GET, got %s", r.Method)
}

w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"rules":[]}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &RulesListOption{Ps: "50"}
	_, resp, err := client.Rules.List(opt)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}
