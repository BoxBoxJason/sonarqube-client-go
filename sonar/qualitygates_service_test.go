package sonargo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestQualityGates_AddGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "qualitygates/add_group") {
			t.Errorf("expected path to contain qualitygates/add_group, got %s", r.URL.Path)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualitygatesAddGroupOption{
		GateName:  "SonarSource Way",
		GroupName: "sonar-administrators",
	}

	resp, err := client.Qualitygates.AddGroup(opt)
	if err != nil {
		t.Fatalf("AddGroup failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityGates_AddGroup_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.Qualitygates.AddGroup(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing GateName
	_, err = client.Qualitygates.AddGroup(&QualitygatesAddGroupOption{
		GroupName: "group",
	})
	if err == nil {
		t.Error("expected error for missing GateName")
	}

	// Test missing GroupName
	_, err = client.Qualitygates.AddGroup(&QualitygatesAddGroupOption{
		GateName: "gate",
	})
	if err == nil {
		t.Error("expected error for missing GroupName")
	}

	// Test GateName too long
	_, err = client.Qualitygates.AddGroup(&QualitygatesAddGroupOption{
		GateName:  strings.Repeat("a", MaxQualityGateNameLength+1),
		GroupName: "group",
	})
	if err == nil {
		t.Error("expected error for GateName exceeding max length")
	}
}

func TestQualityGates_AddUser(t *testing.T) {
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

	opt := &QualitygatesAddUserOption{
		GateName: "SonarSource Way",
		Login:    "john.doe",
	}

	resp, err := client.Qualitygates.AddUser(opt)
	if err != nil {
		t.Fatalf("AddUser failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityGates_Copy(t *testing.T) {
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

	opt := &QualitygatesCopyOption{
		Name:       "My New Quality Gate",
		SourceName: "SonarSource Way",
	}

	resp, err := client.Qualitygates.Copy(opt)
	if err != nil {
		t.Fatalf("Copy failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityGates_Create(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"name":"My Quality Gate"}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualitygatesCreateOption{
		Name: "My Quality Gate",
	}

	result, resp, err := client.Qualitygates.Create(opt)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || result.Name != "My Quality Gate" {
		t.Error("expected name 'My Quality Gate'")
	}
}

func TestQualityGates_CreateCondition(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"id":"1","metric":"coverage","op":"LT","error":"80"}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualitygatesCreateConditionOption{
		Error:    "80",
		GateName: "My Quality Gate",
		Metric:   "coverage",
		Op:       "LT",
	}

	result, resp, err := client.Qualitygates.CreateCondition(opt)
	if err != nil {
		t.Fatalf("CreateCondition failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || result.Metric != "coverage" {
		t.Error("unexpected result")
	}
}

func TestQualityGates_CreateCondition_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test invalid Op value
	_, _, err = client.Qualitygates.CreateCondition(&QualitygatesCreateConditionOption{
		Error:    "80",
		GateName: "gate",
		Metric:   "coverage",
		Op:       "INVALID",
	})
	if err == nil {
		t.Error("expected error for invalid Op")
	}

	// Test missing required fields
	_, _, err = client.Qualitygates.CreateCondition(&QualitygatesCreateConditionOption{
		GateName: "gate",
		Metric:   "coverage",
	})
	if err == nil {
		t.Error("expected error for missing Error")
	}
}

func TestQualityGates_DeleteCondition(t *testing.T) {
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

	opt := &QualitygatesDeleteConditionOption{
		ID: "1",
	}

	resp, err := client.Qualitygates.DeleteCondition(opt)
	if err != nil {
		t.Fatalf("DeleteCondition failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityGates_Deselect(t *testing.T) {
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

	opt := &QualitygatesDeselectOption{
		ProjectKey: "my_project",
	}

	resp, err := client.Qualitygates.Deselect(opt)
	if err != nil {
		t.Fatalf("Deselect failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityGates_Destroy(t *testing.T) {
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

	opt := &QualitygatesDestroyOption{
		Name: "My Quality Gate",
	}

	resp, err := client.Qualitygates.Destroy(opt)
	if err != nil {
		t.Fatalf("Destroy failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityGates_GetByProject(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{"qualityGate":{"name":"SonarSource Way","default":true}}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualitygatesGetByProjectOption{
		Project: "my_project",
	}

	result, resp, err := client.Qualitygates.GetByProject(opt)
	if err != nil {
		t.Fatalf("GetByProject failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || result.QualityGate.Name != "SonarSource Way" {
		t.Error("unexpected quality gate name")
	}

	if !result.QualityGate.Default {
		t.Error("expected quality gate to be default")
	}
}

func TestQualityGates_List(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{
			"actions":{"create":true},
			"qualitygates":[
				{"name":"SonarSource Way","isDefault":true,"isBuiltIn":true},
				{"name":"Custom Gate","isDefault":false,"isBuiltIn":false}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Qualitygates.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.Qualitygates) != 2 {
		t.Fatal("expected 2 quality gates")
	}

	if !result.Actions.Create {
		t.Error("expected create action to be true")
	}

	if result.Qualitygates[0].Name != "SonarSource Way" {
		t.Errorf("expected first gate name 'SonarSource Way', got '%s'", result.Qualitygates[0].Name)
	}
}

func TestQualityGates_ProjectStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{
			"projectStatus":{
				"status":"OK",
				"caycStatus":"compliant",
				"conditions":[
					{"status":"OK","metricKey":"coverage","comparator":"LT","errorThreshold":"80","actualValue":"85"}
				],
				"ignoredConditions":false
			}
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualitygatesProjectStatusOption{
		ProjectKey: "my_project",
	}

	result, resp, err := client.Qualitygates.ProjectStatus(opt)
	if err != nil {
		t.Fatalf("ProjectStatus failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || result.ProjectStatus.Status != "OK" {
		t.Error("expected status 'OK'")
	}

	if len(result.ProjectStatus.Conditions) != 1 {
		t.Error("expected 1 condition")
	}

	if result.ProjectStatus.Conditions[0].MetricKey != "coverage" {
		t.Error("expected metric key 'coverage'")
	}
}

func TestQualityGates_ProjectStatus_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.Qualitygates.ProjectStatus(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing all required fields
	_, _, err = client.Qualitygates.ProjectStatus(&QualitygatesProjectStatusOption{})
	if err == nil {
		t.Error("expected error when no identifier is provided")
	}
}

func TestQualityGates_RemoveGroup(t *testing.T) {
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

	opt := &QualitygatesRemoveGroupOption{
		GateName:  "SonarSource Way",
		GroupName: "sonar-administrators",
	}

	resp, err := client.Qualitygates.RemoveGroup(opt)
	if err != nil {
		t.Fatalf("RemoveGroup failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityGates_RemoveUser(t *testing.T) {
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

	opt := &QualitygatesRemoveUserOption{
		GateName: "SonarSource Way",
		Login:    "john.doe",
	}

	resp, err := client.Qualitygates.RemoveUser(opt)
	if err != nil {
		t.Fatalf("RemoveUser failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityGates_Rename(t *testing.T) {
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

	opt := &QualitygatesRenameOption{
		CurrentName: "Old Name",
		Name:        "New Name",
	}

	resp, err := client.Qualitygates.Rename(opt)
	if err != nil {
		t.Fatalf("Rename failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityGates_Search(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{
			"paging":{"pageIndex":1,"pageSize":100,"total":2},
			"results":[
				{"key":"project1","name":"Project 1","selected":true},
				{"key":"project2","name":"Project 2","selected":false}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualitygatesSearchOption{
		GateName: "SonarSource Way",
	}

	result, resp, err := client.Qualitygates.Search(opt)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.Results) != 2 {
		t.Fatal("expected 2 results")
	}

	if result.Results[0].Key != "project1" {
		t.Errorf("expected first project key 'project1', got '%s'", result.Results[0].Key)
	}
}

func TestQualityGates_Search_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test invalid Selected value
	_, _, err = client.Qualitygates.Search(&QualitygatesSearchOption{
		GateName: "gate",
		Selected: "invalid",
	})
	if err == nil {
		t.Error("expected error for invalid Selected value")
	}
}

func TestQualityGates_SearchGroups(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{
			"paging":{"pageIndex":1,"pageSize":25,"total":1},
			"groups":[
				{"name":"sonar-administrators","description":"Administrators","selected":true}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualitygatesSearchGroupsOption{
		GateName: "SonarSource Way",
	}

	result, resp, err := client.Qualitygates.SearchGroups(opt)
	if err != nil {
		t.Fatalf("SearchGroups failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.Groups) != 1 {
		t.Fatal("expected 1 group")
	}

	if result.Groups[0].Name != "sonar-administrators" {
		t.Errorf("expected group name 'sonar-administrators', got '%s'", result.Groups[0].Name)
	}
}

func TestQualityGates_SearchUsers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{
			"paging":{"pageIndex":1,"pageSize":25,"total":1},
			"users":[
				{"login":"john.doe","name":"John Doe","selected":true}
			]
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualitygatesSearchUsersOption{
		GateName: "SonarSource Way",
	}

	result, resp, err := client.Qualitygates.SearchUsers(opt)
	if err != nil {
		t.Fatalf("SearchUsers failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.Users) != 1 {
		t.Fatal("expected 1 user")
	}

	if result.Users[0].Login != "john.doe" {
		t.Errorf("expected login 'john.doe', got '%s'", result.Users[0].Login)
	}
}

func TestQualityGates_Select(t *testing.T) {
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

	opt := &QualitygatesSelectOption{
		GateName:   "SonarSource Way",
		ProjectKey: "my_project",
	}

	resp, err := client.Qualitygates.Select(opt)
	if err != nil {
		t.Fatalf("Select failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityGates_SetAsDefault(t *testing.T) {
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

	opt := &QualitygatesSetAsDefaultOption{
		Name: "SonarSource Way",
	}

	resp, err := client.Qualitygates.SetAsDefault(opt)
	if err != nil {
		t.Fatalf("SetAsDefault failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityGates_Show(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{
			"name":"SonarSource Way",
			"isDefault":true,
			"isBuiltIn":true,
			"caycStatus":"compliant",
			"conditions":[
				{"id":"1","metric":"coverage","op":"LT","error":"80"}
			],
			"actions":{
				"rename":true,
				"delete":false,
				"manageConditions":true
			}
		}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualitygatesShowOption{
		Name: "SonarSource Way",
	}

	result, resp, err := client.Qualitygates.Show(opt)
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || result.Name != "SonarSource Way" {
		t.Error("expected name 'SonarSource Way'")
	}

	if !result.IsDefault {
		t.Error("expected IsDefault to be true")
	}

	if !result.IsBuiltIn {
		t.Error("expected IsBuiltIn to be true")
	}

	if len(result.Conditions) != 1 {
		t.Fatal("expected 1 condition")
	}

	if result.Conditions[0].Metric != "coverage" {
		t.Errorf("expected metric 'coverage', got '%s'", result.Conditions[0].Metric)
	}

	if !result.Actions.ManageConditions {
		t.Error("expected ManageConditions action to be true")
	}
}

func TestQualityGates_UpdateCondition(t *testing.T) {
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

	opt := &QualitygatesUpdateConditionOption{
		Error:  "85",
		ID:     "1",
		Metric: "coverage",
		Op:     "LT",
	}

	resp, err := client.Qualitygates.UpdateCondition(opt)
	if err != nil {
		t.Fatalf("UpdateCondition failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityGates_UpdateCondition_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.Qualitygates.UpdateCondition(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test invalid Op value
	_, err = client.Qualitygates.UpdateCondition(&QualitygatesUpdateConditionOption{
		Error:  "80",
		ID:     "1",
		Metric: "coverage",
		Op:     "INVALID",
	})
	if err == nil {
		t.Error("expected error for invalid Op")
	}

	// Test missing required fields
	_, err = client.Qualitygates.UpdateCondition(&QualitygatesUpdateConditionOption{
		ID:     "1",
		Metric: "coverage",
	})
	if err == nil {
		t.Error("expected error for missing Error")
	}
}

// TestQualitygatesList_JSONUnmarshal verifies the JSON unmarshal for QualitygatesList.
func TestQualitygatesList_JSONUnmarshal(t *testing.T) {
	jsonData := `{
		"actions":{"create":true},
		"qualitygates":[
			{
				"name":"Sonar way",
				"isDefault":true,
				"isBuiltIn":true,
				"caycStatus":"compliant",
				"hasMQRConditions":false,
				"hasStandardConditions":true,
				"isAiCodeSupported":false,
				"actions":{
					"rename":false,
					"delete":false,
					"manageConditions":true,
					"copy":true,
					"setAsDefault":true,
					"associateProjects":true,
					"delegate":false,
					"manageAiCodeAssurance":false
				}
			}
		]
	}`

	var response QualitygatesList

	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if !response.Actions.Create {
		t.Error("expected Create action to be true")
	}

	if len(response.Qualitygates) != 1 {
		t.Fatalf("expected 1 quality gate, got %d", len(response.Qualitygates))
	}

	gate := response.Qualitygates[0]

	if gate.Name != "Sonar way" {
		t.Errorf("expected name 'Sonar way', got '%s'", gate.Name)
	}

	if !gate.IsDefault {
		t.Error("expected IsDefault to be true")
	}

	if !gate.IsBuiltIn {
		t.Error("expected IsBuiltIn to be true")
	}

	if gate.CaycStatus != "compliant" {
		t.Errorf("expected caycStatus 'compliant', got '%s'", gate.CaycStatus)
	}

	if gate.HasMQRConditions {
		t.Error("expected HasMQRConditions to be false")
	}

	if !gate.HasStandardConditions {
		t.Error("expected HasStandardConditions to be true")
	}
}

// TestQualitygatesProjectStatus_JSONUnmarshal verifies the JSON unmarshal for project status.
func TestQualitygatesProjectStatus_JSONUnmarshal(t *testing.T) {
	jsonData := `{
		"projectStatus":{
			"status":"ERROR",
			"caycStatus":"non-compliant",
			"conditions":[
				{
					"status":"ERROR",
					"metricKey":"new_coverage",
					"comparator":"LT",
					"errorThreshold":"80",
					"actualValue":"70.5"
				},
				{
					"status":"OK",
					"metricKey":"new_duplicated_lines_density",
					"comparator":"GT",
					"errorThreshold":"3",
					"actualValue":"1.2"
				}
			],
			"period":{
				"mode":"PREVIOUS_VERSION",
				"date":"2024-01-15T10:30:00+0000"
			},
			"ignoredConditions":false
		}
	}`

	var response QualitygatesProjectStatus

	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if response.ProjectStatus.Status != "ERROR" {
		t.Errorf("expected status 'ERROR', got '%s'", response.ProjectStatus.Status)
	}

	if response.ProjectStatus.CaycStatus != "non-compliant" {
		t.Errorf("expected caycStatus 'non-compliant', got '%s'", response.ProjectStatus.CaycStatus)
	}

	if len(response.ProjectStatus.Conditions) != 2 {
		t.Fatalf("expected 2 conditions, got %d", len(response.ProjectStatus.Conditions))
	}

	cond := response.ProjectStatus.Conditions[0]

	if cond.MetricKey != "new_coverage" {
		t.Errorf("expected metricKey 'new_coverage', got '%s'", cond.MetricKey)
	}

	if cond.ActualValue != "70.5" {
		t.Errorf("expected actualValue '70.5', got '%s'", cond.ActualValue)
	}

	if response.ProjectStatus.Period.Mode != "PREVIOUS_VERSION" {
		t.Errorf("expected period mode 'PREVIOUS_VERSION', got '%s'", response.ProjectStatus.Period.Mode)
	}
}

// TestQualitygatesShow_JSONUnmarshal verifies the JSON unmarshal for show response.
func TestQualitygatesShow_JSONUnmarshal(t *testing.T) {
	jsonData := `{
		"name":"My Quality Gate",
		"isDefault":false,
		"isBuiltIn":false,
		"caycStatus":"over-compliant",
		"isAiCodeSupported":true,
		"conditions":[
			{"id":"AWGzl6C3r_TYPqzFCFqN","metric":"new_coverage","op":"LT","error":"80"},
			{"id":"AWGzl6C3r_TYPqzFCFqO","metric":"new_duplicated_lines_density","op":"GT","error":"3"}
		],
		"actions":{
			"rename":true,
			"setAsDefault":true,
			"copy":true,
			"associateProjects":true,
			"delete":true,
			"manageConditions":true,
			"delegate":true,
			"manageAiCodeAssurance":true
		}
	}`

	var response QualitygatesShow

	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if response.Name != "My Quality Gate" {
		t.Errorf("expected name 'My Quality Gate', got '%s'", response.Name)
	}

	if response.IsDefault {
		t.Error("expected IsDefault to be false")
	}

	if response.CaycStatus != "over-compliant" {
		t.Errorf("expected caycStatus 'over-compliant', got '%s'", response.CaycStatus)
	}

	if !response.IsAiCodeSupported {
		t.Error("expected IsAiCodeSupported to be true")
	}

	if len(response.Conditions) != 2 {
		t.Fatalf("expected 2 conditions, got %d", len(response.Conditions))
	}

	if response.Conditions[0].ID != "AWGzl6C3r_TYPqzFCFqN" {
		t.Errorf("expected first condition id 'AWGzl6C3r_TYPqzFCFqN', got '%s'", response.Conditions[0].ID)
	}

	if !response.Actions.Rename {
		t.Error("expected Rename action to be true")
	}

	if !response.Actions.ManageAiCodeAssurance {
		t.Error("expected ManageAiCodeAssurance action to be true")
	}
}

// TestValidation_EdgeCases tests edge cases for validation functions.
func TestValidation_EdgeCases(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test Copy validation - missing name
	err = client.Qualitygates.ValidateCopyOpt(&QualitygatesCopyOption{
		SourceName: "source",
	})
	if err == nil {
		t.Error("expected error for missing Name in Copy")
	}

	// Test Delete condition - nil option
	err = client.Qualitygates.ValidateDeleteConditionOpt(nil)
	if err == nil {
		t.Error("expected error for nil DeleteCondition option")
	}

	// Test Rename - both names too long
	err = client.Qualitygates.ValidateRenameOpt(&QualitygatesRenameOption{
		CurrentName: strings.Repeat("a", MaxQualityGateNameLength+1),
		Name:        "new",
	})
	if err == nil {
		t.Error("expected error for CurrentName exceeding max length")
	}

	// Test GetByProject - nil option
	err = client.Qualitygates.ValidateGetByProjectOpt(nil)
	if err == nil {
		t.Error("expected error for nil GetByProject option")
	}

	// Test Search groups - invalid selected
	err = client.Qualitygates.ValidateSearchGroupsOpt(&QualitygatesSearchGroupsOption{
		GateName: "gate",
		Selected: "invalid",
	})
	if err == nil {
		t.Error("expected error for invalid Selected in SearchGroups")
	}

	// Test Search users - invalid selected
	err = client.Qualitygates.ValidateSearchUsersOpt(&QualitygatesSearchUsersOption{
		GateName: "gate",
		Selected: "invalid",
	})
	if err == nil {
		t.Error("expected error for invalid Selected in SearchUsers")
	}

	// Test AddUser - nil option
	err = client.Qualitygates.ValidateAddUserOpt(nil)
	if err == nil {
		t.Error("expected error for nil AddUser option")
	}

	// Test RemoveGroup - nil option
	err = client.Qualitygates.ValidateRemoveGroupOpt(nil)
	if err == nil {
		t.Error("expected error for nil RemoveGroup option")
	}

	// Test RemoveUser - nil option
	err = client.Qualitygates.ValidateRemoveUserOpt(nil)
	if err == nil {
		t.Error("expected error for nil RemoveUser option")
	}

	// Test Deselect - nil option
	err = client.Qualitygates.ValidateDeselectOpt(nil)
	if err == nil {
		t.Error("expected error for nil Deselect option")
	}

	// Test Select - nil option
	err = client.Qualitygates.ValidateSelectOpt(nil)
	if err == nil {
		t.Error("expected error for nil Select option")
	}

	// Test SetAsDefault - nil option
	err = client.Qualitygates.ValidateSetAsDefaultOpt(nil)
	if err == nil {
		t.Error("expected error for nil SetAsDefault option")
	}

	// Test Show - nil option
	err = client.Qualitygates.ValidateShowOpt(nil)
	if err == nil {
		t.Error("expected error for nil Show option")
	}

	// Test Create - nil option
	err = client.Qualitygates.ValidateCreateOpt(nil)
	if err == nil {
		t.Error("expected error for nil Create option")
	}

	// Test CreateCondition - Error too long
	err = client.Qualitygates.ValidateCreateConditionOpt(&QualitygatesCreateConditionOption{
		Error:    strings.Repeat("a", MaxConditionErrorLength+1),
		GateName: "gate",
		Metric:   "coverage",
	})
	if err == nil {
		t.Error("expected error for Error exceeding max length")
	}
}
