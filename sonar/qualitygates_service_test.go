package sonargo

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQualityGates_AddGroup(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualitygates/add_group", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &QualitygatesAddGroupOption{
		GateName:  "SonarSource Way",
		GroupName: "sonar-administrators",
	}

	resp, err := client.Qualitygates.AddGroup(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestQualityGates_AddGroup_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualitygates.AddGroup(nil)
	assert.Error(t, err)

	// Test missing GateName
	_, err = client.Qualitygates.AddGroup(&QualitygatesAddGroupOption{
		GroupName: "group",
	})
	assert.Error(t, err)

	// Test missing GroupName
	_, err = client.Qualitygates.AddGroup(&QualitygatesAddGroupOption{
		GateName: "gate",
	})
	assert.Error(t, err)

	// Test GateName too long
	_, err = client.Qualitygates.AddGroup(&QualitygatesAddGroupOption{
		GateName:  strings.Repeat("a", MaxQualityGateNameLength+1),
		GroupName: "group",
	})
	assert.Error(t, err)
}

func TestQualityGates_AddUser(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualitygates/add_user", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &QualitygatesAddUserOption{
		GateName: "SonarSource Way",
		Login:    "john.doe",
	}

	resp, err := client.Qualitygates.AddUser(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestQualityGates_Copy(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualitygates/copy", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &QualitygatesCopyOption{
		Name:       "My New Quality Gate",
		SourceName: "SonarSource Way",
	}

	resp, err := client.Qualitygates.Copy(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestQualityGates_Create(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/qualitygates/create", http.StatusOK, `{"name":"My Quality Gate"}`))
	client := newTestClient(t, server.URL)

	opt := &QualitygatesCreateOption{
		Name: "My Quality Gate",
	}

	result, resp, err := client.Qualitygates.Create(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "My Quality Gate", result.Name)
}

func TestQualityGates_CreateCondition(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/qualitygates/create_condition", http.StatusOK, `{"id":"1","metric":"coverage","op":"LT","error":"80"}`))
	client := newTestClient(t, server.URL)

	opt := &QualitygatesCreateConditionOption{
		Error:    "80",
		GateName: "My Quality Gate",
		Metric:   "coverage",
		Op:       "LT",
	}

	result, resp, err := client.Qualitygates.CreateCondition(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "coverage", result.Metric)
}

func TestQualityGates_CreateCondition_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test invalid Op value
	_, _, err := client.Qualitygates.CreateCondition(&QualitygatesCreateConditionOption{
		Error:    "80",
		GateName: "gate",
		Metric:   "coverage",
		Op:       "INVALID",
	})
	assert.Error(t, err)

	// Test missing required fields
	_, _, err = client.Qualitygates.CreateCondition(&QualitygatesCreateConditionOption{
		GateName: "gate",
		Metric:   "coverage",
	})
	assert.Error(t, err)
}

func TestQualityGates_DeleteCondition(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualitygates/delete_condition", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &QualitygatesDeleteConditionOption{
		ID: "1",
	}

	resp, err := client.Qualitygates.DeleteCondition(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestQualityGates_Deselect(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualitygates/deselect", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &QualitygatesDeselectOption{
		ProjectKey: "my_project",
	}

	resp, err := client.Qualitygates.Deselect(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestQualityGates_Destroy(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualitygates/destroy", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &QualitygatesDestroyOption{
		Name: "My Quality Gate",
	}

	resp, err := client.Qualitygates.Destroy(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestQualityGates_GetByProject(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualitygates/get_by_project", http.StatusOK, `{"qualityGate":{"name":"SonarSource Way","default":true}}`))
	client := newTestClient(t, server.URL)

	opt := &QualitygatesGetByProjectOption{
		Project: "my_project",
	}

	result, resp, err := client.Qualitygates.GetByProject(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "SonarSource Way", result.QualityGate.Name)
	assert.True(t, result.QualityGate.Default)
}

func TestQualityGates_List(t *testing.T) {
	response := `{
		"actions":{"create":true},
		"qualitygates":[
			{"name":"SonarSource Way","isDefault":true,"isBuiltIn":true},
			{"name":"Custom Gate","isDefault":false,"isBuiltIn":false}
		]
	}`
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualitygates/list", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Qualitygates.List()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Qualitygates, 2)
	assert.True(t, result.Actions.Create)
	assert.Equal(t, "SonarSource Way", result.Qualitygates[0].Name)
}

func TestQualityGates_ProjectStatus(t *testing.T) {
	response := `{
		"projectStatus":{
			"status":"OK",
			"caycStatus":"compliant",
			"conditions":[
				{"status":"OK","metricKey":"coverage","comparator":"LT","errorThreshold":"80","actualValue":"85"}
			],
			"ignoredConditions":false
		}
	}`
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualitygates/project_status", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &QualitygatesProjectStatusOption{
		ProjectKey: "my_project",
	}

	result, resp, err := client.Qualitygates.ProjectStatus(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "OK", result.ProjectStatus.Status)
	assert.Len(t, result.ProjectStatus.Conditions, 1)
	assert.Equal(t, "coverage", result.ProjectStatus.Conditions[0].MetricKey)
}

func TestQualityGates_ProjectStatus_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualitygates.ProjectStatus(nil)
	assert.Error(t, err)

	// Test missing all required fields
	_, _, err = client.Qualitygates.ProjectStatus(&QualitygatesProjectStatusOption{})
	assert.Error(t, err)
}

func TestQualityGates_RemoveGroup(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualitygates/remove_group", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &QualitygatesRemoveGroupOption{
		GateName:  "SonarSource Way",
		GroupName: "sonar-administrators",
	}

	resp, err := client.Qualitygates.RemoveGroup(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestQualityGates_RemoveUser(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualitygates/remove_user", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &QualitygatesRemoveUserOption{
		GateName: "SonarSource Way",
		Login:    "john.doe",
	}

	resp, err := client.Qualitygates.RemoveUser(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestQualityGates_Rename(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualitygates/rename", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &QualitygatesRenameOption{
		CurrentName: "Old Name",
		Name:        "New Name",
	}

	resp, err := client.Qualitygates.Rename(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestQualityGates_Search(t *testing.T) {
	response := `{
		"paging":{"pageIndex":1,"pageSize":100,"total":2},
		"results":[
			{"key":"project1","name":"Project 1","selected":true},
			{"key":"project2","name":"Project 2","selected":false}
		]
	}`
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualitygates/search", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &QualitygatesSearchOption{
		GateName: "SonarSource Way",
	}

	result, resp, err := client.Qualitygates.Search(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Results, 2)
	assert.Equal(t, "project1", result.Results[0].Key)
}

func TestQualityGates_Search_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test invalid Selected value
	_, _, err := client.Qualitygates.Search(&QualitygatesSearchOption{
		GateName: "gate",
		Selected: "invalid",
	})
	assert.Error(t, err)
}

func TestQualityGates_SearchGroups(t *testing.T) {
	response := `{
		"paging":{"pageIndex":1,"pageSize":25,"total":1},
		"groups":[
			{"name":"sonar-administrators","description":"Administrators","selected":true}
		]
	}`
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualitygates/search_groups", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &QualitygatesSearchGroupsOption{
		GateName: "SonarSource Way",
	}

	result, resp, err := client.Qualitygates.SearchGroups(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Groups, 1)
	assert.Equal(t, "sonar-administrators", result.Groups[0].Name)
}

func TestQualityGates_SearchUsers(t *testing.T) {
	response := `{
		"paging":{"pageIndex":1,"pageSize":25,"total":1},
		"users":[
			{"login":"john.doe","name":"John Doe","selected":true}
		]
	}`
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualitygates/search_users", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &QualitygatesSearchUsersOption{
		GateName: "SonarSource Way",
	}

	result, resp, err := client.Qualitygates.SearchUsers(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Users, 1)
	assert.Equal(t, "john.doe", result.Users[0].Login)
}

func TestQualityGates_Select(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualitygates/select", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &QualitygatesSelectOption{
		GateName:   "SonarSource Way",
		ProjectKey: "my_project",
	}

	resp, err := client.Qualitygates.Select(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestQualityGates_SetAsDefault(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualitygates/set_as_default", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &QualitygatesSetAsDefaultOption{
		Name: "SonarSource Way",
	}

	resp, err := client.Qualitygates.SetAsDefault(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestQualityGates_Show(t *testing.T) {
	response := `{
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
	}`
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualitygates/show", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &QualitygatesShowOption{
		Name: "SonarSource Way",
	}

	result, resp, err := client.Qualitygates.Show(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "SonarSource Way", result.Name)
	assert.True(t, result.IsDefault)
	assert.True(t, result.IsBuiltIn)
	assert.Len(t, result.Conditions, 1)
	assert.Equal(t, "coverage", result.Conditions[0].Metric)
	assert.True(t, result.Actions.ManageConditions)
}

func TestQualityGates_UpdateCondition(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualitygates/update_condition", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	opt := &QualitygatesUpdateConditionOption{
		Error:  "85",
		ID:     "1",
		Metric: "coverage",
		Op:     "LT",
	}

	resp, err := client.Qualitygates.UpdateCondition(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestQualityGates_UpdateCondition_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualitygates.UpdateCondition(nil)
	assert.Error(t, err)

	// Test invalid Op value
	_, err = client.Qualitygates.UpdateCondition(&QualitygatesUpdateConditionOption{
		Error:  "80",
		ID:     "1",
		Metric: "coverage",
		Op:     "INVALID",
	})
	assert.Error(t, err)

	// Test missing required fields
	_, err = client.Qualitygates.UpdateCondition(&QualitygatesUpdateConditionOption{
		ID:     "1",
		Metric: "coverage",
	})
	assert.Error(t, err)
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
	require.NoError(t, err)
	assert.True(t, response.Actions.Create)
	require.Len(t, response.Qualitygates, 1)

	gate := response.Qualitygates[0]
	assert.Equal(t, "Sonar way", gate.Name)
	assert.True(t, gate.IsDefault)
	assert.True(t, gate.IsBuiltIn)
	assert.Equal(t, "compliant", gate.CaycStatus)
	assert.False(t, gate.HasMQRConditions)
	assert.True(t, gate.HasStandardConditions)
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
	require.NoError(t, err)
	assert.Equal(t, "ERROR", response.ProjectStatus.Status)
	assert.Equal(t, "non-compliant", response.ProjectStatus.CaycStatus)
	require.Len(t, response.ProjectStatus.Conditions, 2)

	cond := response.ProjectStatus.Conditions[0]
	assert.Equal(t, "new_coverage", cond.MetricKey)
	assert.Equal(t, "70.5", cond.ActualValue)
	assert.Equal(t, "PREVIOUS_VERSION", response.ProjectStatus.Period.Mode)
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
	require.NoError(t, err)
	assert.Equal(t, "My Quality Gate", response.Name)
	assert.False(t, response.IsDefault)
	assert.Equal(t, "over-compliant", response.CaycStatus)
	assert.True(t, response.IsAiCodeSupported)
	require.Len(t, response.Conditions, 2)
	assert.Equal(t, "AWGzl6C3r_TYPqzFCFqN", response.Conditions[0].ID)
	assert.True(t, response.Actions.Rename)
	assert.True(t, response.Actions.ManageAiCodeAssurance)
}

// TestValidation_EdgeCases tests edge cases for validation functions.
func TestValidation_EdgeCases(t *testing.T) {
	client := newLocalhostClient(t)

	// Test Copy validation - missing name
	err := client.Qualitygates.ValidateCopyOpt(&QualitygatesCopyOption{
		SourceName: "source",
	})
	assert.Error(t, err)

	// Test Delete condition - nil option
	err = client.Qualitygates.ValidateDeleteConditionOpt(nil)
	assert.Error(t, err)

	// Test Rename - both names too long
	err = client.Qualitygates.ValidateRenameOpt(&QualitygatesRenameOption{
		CurrentName: strings.Repeat("a", MaxQualityGateNameLength+1),
		Name:        "new",
	})
	assert.Error(t, err)

	// Test GetByProject - nil option
	err = client.Qualitygates.ValidateGetByProjectOpt(nil)
	assert.Error(t, err)

	// Test Search groups - invalid selected
	err = client.Qualitygates.ValidateSearchGroupsOpt(&QualitygatesSearchGroupsOption{
		GateName: "gate",
		Selected: "invalid",
	})
	assert.Error(t, err)

	// Test Search users - invalid selected
	err = client.Qualitygates.ValidateSearchUsersOpt(&QualitygatesSearchUsersOption{
		GateName: "gate",
		Selected: "invalid",
	})
	assert.Error(t, err)

	// Test AddUser - nil option
	err = client.Qualitygates.ValidateAddUserOpt(nil)
	assert.Error(t, err)

	// Test RemoveGroup - nil option
	err = client.Qualitygates.ValidateRemoveGroupOpt(nil)
	assert.Error(t, err)

	// Test RemoveUser - nil option
	err = client.Qualitygates.ValidateRemoveUserOpt(nil)
	assert.Error(t, err)

	// Test Deselect - nil option
	err = client.Qualitygates.ValidateDeselectOpt(nil)
	assert.Error(t, err)

	// Test Select - nil option
	err = client.Qualitygates.ValidateSelectOpt(nil)
	assert.Error(t, err)

	// Test SetAsDefault - nil option
	err = client.Qualitygates.ValidateSetAsDefaultOpt(nil)
	assert.Error(t, err)

	// Test Show - nil option
	err = client.Qualitygates.ValidateShowOpt(nil)
	assert.Error(t, err)

	// Test Create - nil option
	err = client.Qualitygates.ValidateCreateOpt(nil)
	assert.Error(t, err)

	// Test CreateCondition - Error too long
	err = client.Qualitygates.ValidateCreateConditionOpt(&QualitygatesCreateConditionOption{
		Error:    strings.Repeat("a", MaxConditionErrorLength+1),
		GateName: "gate",
		Metric:   "coverage",
	})
	assert.Error(t, err)
}
