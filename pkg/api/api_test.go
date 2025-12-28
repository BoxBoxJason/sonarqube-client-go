package api

import (
	"encoding/json"
	"testing"
)

func TestAPIStruct(t *testing.T) {
	jsonData := `{
		"webServices": [
			{
				"path": "api/projects",
				"description": "Manage project",
				"since": "1.0",
				"actions": [
					{
						"key": "search",
						"description": "Search for projects",
						"since": "2.0",
						"internal": false,
						"post": false,
						"hasResponseExample": true,
						"changelog": [],
						"params": []
					}
				]
			}
		]
	}`

	var api API
	err := json.Unmarshal([]byte(jsonData), &api)
	if err != nil {
		t.Fatalf("Failed to unmarshal API: %v", err)
	}

	if len(api.WebServices) != 1 {
		t.Errorf("Expected 1 web service, got %d", len(api.WebServices))
	}

	ws := api.WebServices[0]
	if ws.Path != "api/projects" {
		t.Errorf("Expected path 'api/projects', got '%s'", ws.Path)
	}
	if ws.Description != "Manage project" {
		t.Errorf("Expected description 'Manage project', got '%s'", ws.Description)
	}
	if ws.Since != "1.0" {
		t.Errorf("Expected since '1.0', got '%s'", ws.Since)
	}
}

func TestWebServiceStruct(t *testing.T) {
	jsonData := `{
		"path": "api/issues",
		"description": "Manage issues",
		"since": "3.0",
		"actions": [
			{
				"key": "search",
				"description": "Search issues",
				"since": "3.0",
				"internal": false,
				"post": false,
				"hasResponseExample": true,
				"changelog": [
					{"description": "Added field", "version": "3.1"}
				],
				"params": [
					{
						"key": "projectKey",
						"description": "Project key",
						"required": true,
						"internal": false
					}
				]
			}
		]
	}`

	var ws WebService
	err := json.Unmarshal([]byte(jsonData), &ws)
	if err != nil {
		t.Fatalf("Failed to unmarshal WebService: %v", err)
	}

	if ws.Path != "api/issues" {
		t.Errorf("Expected path 'api/issues', got '%s'", ws.Path)
	}
	if len(ws.Actions) != 1 {
		t.Errorf("Expected 1 action, got %d", len(ws.Actions))
	}
}

func TestActionStruct(t *testing.T) {
	jsonData := `{
		"key": "create",
		"description": "Create a project",
		"since": "2.0",
		"internal": true,
		"post": true,
		"hasResponseExample": false,
		"deprecatedSince": "5.0",
		"responseType": "json",
		"changelog": [
			{"description": "Initial version", "version": "2.0"},
			{"description": "Added features", "version": "3.0"}
		],
		"params": [
			{
				"key": "name",
				"description": "Project name",
				"required": true,
				"internal": false,
				"exampleValue": "My Project",
				"defaultValue": "",
				"possibleValues": ["val1", "val2"]
			}
		]
	}`

	var action Action
	err := json.Unmarshal([]byte(jsonData), &action)
	if err != nil {
		t.Fatalf("Failed to unmarshal Action: %v", err)
	}

	if action.Key != "create" {
		t.Errorf("Expected key 'create', got '%s'", action.Key)
	}
	if !action.Internal {
		t.Error("Expected internal to be true")
	}
	if !action.Post {
		t.Error("Expected post to be true")
	}
	if action.HasResponseExample {
		t.Error("Expected hasResponseExample to be false")
	}
	if action.DeprecatedSince != "5.0" {
		t.Errorf("Expected deprecatedSince '5.0', got '%s'", action.DeprecatedSince)
	}
	if action.ResponseType != "json" {
		t.Errorf("Expected responseType 'json', got '%s'", action.ResponseType)
	}
	if len(action.Changelog) != 2 {
		t.Errorf("Expected 2 changelog entries, got %d", len(action.Changelog))
	}
	if len(action.Params) != 1 {
		t.Errorf("Expected 1 param, got %d", len(action.Params))
	}
}

func TestParamStruct(t *testing.T) {
	jsonData := `{
		"key": "pageSize",
		"description": "Page size",
		"required": false,
		"internal": false,
		"exampleValue": "100",
		"deprecatedSince": "6.0",
		"defaultValue": "50",
		"possibleValues": ["10", "50", "100"],
		"deprecatedKey": "size",
		"deprecatedKeySince": "5.0",
		"maximumValue": 500,
		"since": "2.0",
		"minimumLength": 1,
		"maxValuesAllowed": 20,
		"maximumLength": 100
	}`

	var param Param
	err := json.Unmarshal([]byte(jsonData), &param)
	if err != nil {
		t.Fatalf("Failed to unmarshal Param: %v", err)
	}

	if param.Key != "pageSize" {
		t.Errorf("Expected key 'pageSize', got '%s'", param.Key)
	}
	if param.Required {
		t.Error("Expected required to be false")
	}
	if param.MaximumValue != 500 {
		t.Errorf("Expected maximumValue 500, got %d", param.MaximumValue)
	}
	if param.MinimumLength != 1 {
		t.Errorf("Expected minimumLength 1, got %d", param.MinimumLength)
	}
	if len(param.PossibleValues) != 3 {
		t.Errorf("Expected 3 possible values, got %d", len(param.PossibleValues))
	}
}

func TestChangelogStruct(t *testing.T) {
	jsonData := `{
		"description": "Added new parameter",
		"version": "4.5"
	}`

	var changelog Changelog
	err := json.Unmarshal([]byte(jsonData), &changelog)
	if err != nil {
		t.Fatalf("Failed to unmarshal Changelog: %v", err)
	}

	if changelog.Description != "Added new parameter" {
		t.Errorf("Expected description 'Added new parameter', got '%s'", changelog.Description)
	}
	if changelog.Version != "4.5" {
		t.Errorf("Expected version '4.5', got '%s'", changelog.Version)
	}
}

func TestAPIMarshalRoundTrip(t *testing.T) {
	api := API{
		WebServices: []WebService{
			{
				Path:        "api/test",
				Description: "Test service",
				Since:       "1.0",
				Actions: []Action{
					{
						Key:         "action1",
						Description: "Test action",
						Since:       "1.0",
						Post:        true,
					},
				},
			},
		},
	}

	data, err := json.Marshal(api)
	if err != nil {
		t.Fatalf("Failed to marshal API: %v", err)
	}

	var api2 API
	err = json.Unmarshal(data, &api2)
	if err != nil {
		t.Fatalf("Failed to unmarshal API: %v", err)
	}

	if len(api2.WebServices) != 1 {
		t.Errorf("Expected 1 web service, got %d", len(api2.WebServices))
	}
	if api2.WebServices[0].Path != "api/test" {
		t.Errorf("Expected path 'api/test', got '%s'", api2.WebServices[0].Path)
	}
}
