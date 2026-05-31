package sonar

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQualityprofiles_ActivateRule(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/activate_rule", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesActivateRuleOptions{
		Key:  "AU-TpxcA-iU5OvuD2FL0",
		Rule: "squid:AvoidCycles",
	}

	resp, err := client.Qualityprofiles.ActivateRule(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_ActivateRule_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.ActivateRule(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Key
	_, err = client.Qualityprofiles.ActivateRule(context.Background(), &QualityprofilesActivateRuleOptions{
		Rule: "squid:AvoidCycles",
	})
	assert.Error(t, err)

	// Test missing Rule
	_, err = client.Qualityprofiles.ActivateRule(context.Background(), &QualityprofilesActivateRuleOptions{
		Key: "AU-TpxcA-iU5OvuD2FL0",
	})
	assert.Error(t, err)

	// Test both Impacts and Severity set
	_, err = client.Qualityprofiles.ActivateRule(context.Background(), &QualityprofilesActivateRuleOptions{
		Key:      "AU-TpxcA-iU5OvuD2FL0",
		Rule:     "squid:AvoidCycles",
		Impacts:  map[string]string{SoftwareQualityMaintainability: RuleImpactSeverityHigh},
		Severity: RuleSeverityMajor,
	})
	assert.Error(t, err)

	// Test invalid Severity
	_, err = client.Qualityprofiles.ActivateRule(context.Background(), &QualityprofilesActivateRuleOptions{
		Key:      "AU-TpxcA-iU5OvuD2FL0",
		Rule:     "squid:AvoidCycles",
		Severity: "INVALID",
	})
	assert.Error(t, err)

	// Test invalid Impacts map key (software quality)
	_, err = client.Qualityprofiles.ActivateRule(context.Background(), &QualityprofilesActivateRuleOptions{
		Key:     "AU-TpxcA-iU5OvuD2FL0",
		Rule:    "squid:AvoidCycles",
		Impacts: map[string]string{"INVALID_QUALITY": RuleImpactSeverityHigh},
	})
	assert.Error(t, err)

	// Test invalid Impacts map value (severity)
	_, err = client.Qualityprofiles.ActivateRule(context.Background(), &QualityprofilesActivateRuleOptions{
		Key:     "AU-TpxcA-iU5OvuD2FL0",
		Rule:    "squid:AvoidCycles",
		Impacts: map[string]string{SoftwareQualityMaintainability: "INVALID_SEVERITY"},
	})
	assert.Error(t, err)
}

func TestQualityprofiles_ActivateRules(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/activate_rules", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesActivateRulesOptions{
		TargetKey: "AU-TpxcA-iU5OvuD2FL0",
		Languages: []string{"java"},
	}

	resp, err := client.Qualityprofiles.ActivateRules(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_ActivateRules_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.ActivateRules(context.Background(), nil)
	assert.Error(t, err)

	// Test missing TargetKey
	_, err = client.Qualityprofiles.ActivateRules(context.Background(), &QualityprofilesActivateRulesOptions{})
	assert.Error(t, err)

	// Test invalid language
	_, err = client.Qualityprofiles.ActivateRules(context.Background(), &QualityprofilesActivateRulesOptions{
		TargetKey: "AU-TpxcA-iU5OvuD2FL0",
		Languages: []string{"invalid_language"},
	})
	assert.Error(t, err)

	// Test invalid severity
	_, err = client.Qualityprofiles.ActivateRules(context.Background(), &QualityprofilesActivateRulesOptions{
		TargetKey:  "AU-TpxcA-iU5OvuD2FL0",
		Severities: []string{"INVALID_SEVERITY"},
	})
	assert.Error(t, err)

	// Test invalid impact severity
	_, err = client.Qualityprofiles.ActivateRules(context.Background(), &QualityprofilesActivateRulesOptions{
		TargetKey:        "AU-TpxcA-iU5OvuD2FL0",
		ImpactSeverities: []string{"INVALID"},
	})
	assert.Error(t, err)

	// Test invalid software quality
	_, err = client.Qualityprofiles.ActivateRules(context.Background(), &QualityprofilesActivateRulesOptions{
		TargetKey:               "AU-TpxcA-iU5OvuD2FL0",
		ImpactSoftwareQualities: []string{"INVALID"},
	})
	assert.Error(t, err)

	// Test invalid sort field
	_, err = client.Qualityprofiles.ActivateRules(context.Background(), &QualityprofilesActivateRulesOptions{
		TargetKey: "AU-TpxcA-iU5OvuD2FL0",
		Sort:      "invalid_sort",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_AddGroup(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/add_group", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesAddGroupOptions{
		Group:          "sonar-administrators",
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.AddGroup(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_AddGroup_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.AddGroup(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Group
	_, err = client.Qualityprofiles.AddGroup(context.Background(), &QualityprofilesAddGroupOptions{
		Language:       "java",
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing Language
	_, err = client.Qualityprofiles.AddGroup(context.Background(), &QualityprofilesAddGroupOptions{
		Group:          "sonar-administrators",
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test invalid Language
	_, err = client.Qualityprofiles.AddGroup(context.Background(), &QualityprofilesAddGroupOptions{
		Group:          "sonar-administrators",
		Language:       "invalid_lang",
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, err = client.Qualityprofiles.AddGroup(context.Background(), &QualityprofilesAddGroupOptions{
		Group:    "sonar-administrators",
		Language: "java",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_AddProject(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/add_project", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesAddProjectOptions{
		Language:       "java",
		Project:        "my_project",
		QualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.AddProject(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_AddProject_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.AddProject(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Language
	_, err = client.Qualityprofiles.AddProject(context.Background(), &QualityprofilesAddProjectOptions{
		Project:        "my_project",
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing Project
	_, err = client.Qualityprofiles.AddProject(context.Background(), &QualityprofilesAddProjectOptions{
		Language:       "java",
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, err = client.Qualityprofiles.AddProject(context.Background(), &QualityprofilesAddProjectOptions{
		Language: "java",
		Project:  "my_project",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_AddUser(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/add_user", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesAddUserOptions{
		Language:       "java",
		Login:          "john.doe",
		QualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.AddUser(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_AddUser_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.AddUser(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Language
	_, err = client.Qualityprofiles.AddUser(context.Background(), &QualityprofilesAddUserOptions{
		Login:          "john.doe",
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing Login
	_, err = client.Qualityprofiles.AddUser(context.Background(), &QualityprofilesAddUserOptions{
		Language:       "java",
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, err = client.Qualityprofiles.AddUser(context.Background(), &QualityprofilesAddUserOptions{
		Language: "java",
		Login:    "john.doe",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_Backup(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/qualityprofiles/backup")
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`<?xml version='1.0'?><profile><name>Sonar way</name></profile>`))
	})
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesBackupOptions{
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	result, resp, err := client.Qualityprofiles.Backup(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Contains(t, *result, "Sonar way")
}

func TestQualityprofiles_Backup_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.Backup(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Language
	_, _, err = client.Qualityprofiles.Backup(context.Background(), &QualityprofilesBackupOptions{
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, _, err = client.Qualityprofiles.Backup(context.Background(), &QualityprofilesBackupOptions{
		Language: "java",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_ChangeParent(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/change_parent", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesChangeParentOptions{
		Language:             "java",
		QualityProfile:       "My Profile",
		ParentQualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.ChangeParent(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_ChangeParent_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.ChangeParent(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Language
	_, err = client.Qualityprofiles.ChangeParent(context.Background(), &QualityprofilesChangeParentOptions{
		QualityProfile: "My Profile",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, err = client.Qualityprofiles.ChangeParent(context.Background(), &QualityprofilesChangeParentOptions{
		Language: "java",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_Changelog(t *testing.T) {
	response := &QualityprofilesChangelog{
		Paging: Paging{PageIndex: 1, PageSize: 25, Total: 1},
		Events: []ChangelogEvent{
			{
				Action:     "ACTIVATED",
				AuthorName: "John Doe",
				RuleKey:    "squid:S1234",
				RuleName:   "Some Rule",
			},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualityprofiles/changelog", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesChangelogOptions{
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	result, resp, err := client.Qualityprofiles.Changelog(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Events, 1)
	assert.Equal(t, "ACTIVATED", result.Events[0].Action)
}

func TestQualityprofiles_Changelog_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.Changelog(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Language
	_, _, err = client.Qualityprofiles.Changelog(context.Background(), &QualityprofilesChangelogOptions{
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, _, err = client.Qualityprofiles.Changelog(context.Background(), &QualityprofilesChangelogOptions{
		Language: "java",
	})
	assert.Error(t, err)

	// Test invalid FilterMode
	_, _, err = client.Qualityprofiles.Changelog(context.Background(), &QualityprofilesChangelogOptions{
		Language:       "java",
		QualityProfile: "Sonar way",
		FilterMode:     "INVALID",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_Compare(t *testing.T) {
	response := &QualityprofilesCompare{
		Left:  QualityprofilesCompareProfile{Key: "profile1", Name: "Profile 1"},
		Right: QualityprofilesCompareProfile{Key: "profile2", Name: "Profile 2"},
		InLeft: []QualityprofilesCompareRule{
			{Key: "squid:S1234", Name: "Rule in left"},
		},
		InRight: []QualityprofilesCompareRule{
			{Key: "squid:S5678", Name: "Rule in right"},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualityprofiles/compare", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesCompareOptions{
		LeftKey:  "profile1",
		RightKey: "profile2",
	}

	result, resp, err := client.Qualityprofiles.Compare(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "Profile 1", result.Left.Name)
	assert.Len(t, result.InLeft, 1)
}

func TestQualityprofiles_Compare_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.Compare(context.Background(), nil)
	assert.Error(t, err)

	// Test missing LeftKey
	_, _, err = client.Qualityprofiles.Compare(context.Background(), &QualityprofilesCompareOptions{
		RightKey: "profile2",
	})
	assert.Error(t, err)

	// Test missing RightKey
	_, _, err = client.Qualityprofiles.Compare(context.Background(), &QualityprofilesCompareOptions{
		LeftKey: "profile1",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_Copy(t *testing.T) {
	response := &QualityprofilesCopy{
		Key:          "new-profile-key",
		Name:         "My Profile Copy",
		Language:     "java",
		LanguageName: "Java",
		IsDefault:    false,
		IsInherited:  false,
	}

	server := newTestServer(t, mockHandler(t, http.MethodPost, "/qualityprofiles/copy", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesCopyOptions{
		FromKey: "source-profile-key",
		ToName:  "My Profile Copy",
	}

	result, resp, err := client.Qualityprofiles.Copy(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "My Profile Copy", result.Name)
}

func TestQualityprofiles_Copy_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.Copy(context.Background(), nil)
	assert.Error(t, err)

	// Test missing FromKey
	_, _, err = client.Qualityprofiles.Copy(context.Background(), &QualityprofilesCopyOptions{
		ToName: "New Profile",
	})
	assert.Error(t, err)

	// Test missing ToName
	_, _, err = client.Qualityprofiles.Copy(context.Background(), &QualityprofilesCopyOptions{
		FromKey: "source-key",
	})
	assert.Error(t, err)

	// Test ToName too long
	_, _, err = client.Qualityprofiles.Copy(context.Background(), &QualityprofilesCopyOptions{
		FromKey: "source-key",
		ToName:  strings.Repeat("a", MaxQualityProfileNameLength+1),
	})
	assert.Error(t, err)
}

func TestQualityprofiles_Create(t *testing.T) {
	response := &QualityprofilesCreate{
		Profile: QualityprofilesCreatedProfile{
			Key:          "new-profile-key",
			Name:         "My New Profile",
			Language:     "java",
			LanguageName: "Java",
			IsDefault:    false,
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodPost, "/qualityprofiles/create", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesCreateOptions{
		Language: "java",
		Name:     "My New Profile",
	}

	result, resp, err := client.Qualityprofiles.Create(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "My New Profile", result.Profile.Name)
}

func TestQualityprofiles_Create_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.Create(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Language
	_, _, err = client.Qualityprofiles.Create(context.Background(), &QualityprofilesCreateOptions{
		Name: "My Profile",
	})
	assert.Error(t, err)

	// Test missing Name
	_, _, err = client.Qualityprofiles.Create(context.Background(), &QualityprofilesCreateOptions{
		Language: "java",
	})
	assert.Error(t, err)

	// Test Name too long
	_, _, err = client.Qualityprofiles.Create(context.Background(), &QualityprofilesCreateOptions{
		Language: "java",
		Name:     strings.Repeat("a", MaxQualityProfileNameLength+1),
	})
	assert.Error(t, err)
}

func TestQualityprofiles_DeactivateRule(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/deactivate_rule", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesDeactivateRuleOptions{
		Key:  "AU-TpxcA-iU5OvuD2FL0",
		Rule: "squid:AvoidCycles",
	}

	resp, err := client.Qualityprofiles.DeactivateRule(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_DeactivateRule_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.DeactivateRule(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Key
	_, err = client.Qualityprofiles.DeactivateRule(context.Background(), &QualityprofilesDeactivateRuleOptions{
		Rule: "squid:AvoidCycles",
	})
	assert.Error(t, err)

	// Test missing Rule
	_, err = client.Qualityprofiles.DeactivateRule(context.Background(), &QualityprofilesDeactivateRuleOptions{
		Key: "AU-TpxcA-iU5OvuD2FL0",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_DeactivateRules(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/deactivate_rules", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesDeactivateRulesOptions{
		TargetKey: "AU-TpxcA-iU5OvuD2FL0",
		Languages: []string{"java"},
	}

	resp, err := client.Qualityprofiles.DeactivateRules(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_DeactivateRules_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.DeactivateRules(context.Background(), nil)
	assert.Error(t, err)

	// Test missing TargetKey
	_, err = client.Qualityprofiles.DeactivateRules(context.Background(), &QualityprofilesDeactivateRulesOptions{})
	assert.Error(t, err)
}

func TestQualityprofiles_Delete(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/delete", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesDeleteOptions{
		Language:       "java",
		QualityProfile: "My Profile",
	}

	resp, err := client.Qualityprofiles.Delete(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_Delete_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.Delete(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Language
	_, err = client.Qualityprofiles.Delete(context.Background(), &QualityprofilesDeleteOptions{
		QualityProfile: "My Profile",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, err = client.Qualityprofiles.Delete(context.Background(), &QualityprofilesDeleteOptions{
		Language: "java",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_Export(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/qualityprofiles/export")
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`<?xml version='1.0'?><profile><name>Sonar way</name></profile>`))
	})
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesExportOptions{
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	result, resp, err := client.Qualityprofiles.Export(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Contains(t, *result, "Sonar way")
}

func TestQualityprofiles_Export_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.Export(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Language
	_, _, err = client.Qualityprofiles.Export(context.Background(), &QualityprofilesExportOptions{})
	assert.Error(t, err)
}

func TestQualityprofiles_Exporters(t *testing.T) {
	response := &QualityprofilesExporters{
		Exporters: []QualityprofilesExporter{
			{Key: "checkstyle", Name: "Checkstyle", Languages: []string{"java"}},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualityprofiles/exporters", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Qualityprofiles.Exporters(context.Background(), )
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Exporters, 1)
	assert.Equal(t, "checkstyle", result.Exporters[0].Key)
}

func TestQualityprofiles_Importers(t *testing.T) {
	response := &QualityprofilesImporters{
		Importers: []QualityprofilesImporter{
			{Key: "pmd", Name: "PMD", Languages: []string{"java"}},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualityprofiles/importers", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Qualityprofiles.Importers(context.Background(), )
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Importers, 1)
	assert.Equal(t, "pmd", result.Importers[0].Key)
}

func TestQualityprofiles_Inheritance(t *testing.T) {
	response := &QualityprofilesInheritance{
		Profile: QualityprofilesInheritanceProfile{
			Key:             "my-profile-key",
			Name:            "My Profile",
			ActiveRuleCount: 150,
		},
		Ancestors: []QualityprofilesInheritanceProfile{
			{Key: "parent-key", Name: "Sonar way", ActiveRuleCount: 200},
		},
		Children: []QualityprofilesInheritanceProfile{},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualityprofiles/inheritance", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesInheritanceOptions{
		Language:       "java",
		QualityProfile: "My Profile",
	}

	result, resp, err := client.Qualityprofiles.Inheritance(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "My Profile", result.Profile.Name)
	assert.Len(t, result.Ancestors, 1)
	assert.Equal(t, "Sonar way", result.Ancestors[0].Name)
}

func TestQualityprofiles_Inheritance_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.Inheritance(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Language
	_, _, err = client.Qualityprofiles.Inheritance(context.Background(), &QualityprofilesInheritanceOptions{
		QualityProfile: "My Profile",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, _, err = client.Qualityprofiles.Inheritance(context.Background(), &QualityprofilesInheritanceOptions{
		Language: "java",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_Projects(t *testing.T) {
	response := &QualityprofilesProjects{
		Paging: Paging{PageIndex: 1, PageSize: 25, Total: 2},
		Results: []QualityprofilesProfileProject{
			{Key: "project1", Name: "Project 1", Selected: true},
			{Key: "project2", Name: "Project 2", Selected: false},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualityprofiles/projects", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesProjectsOptions{
		Key: "profile-key",
	}

	result, resp, err := client.Qualityprofiles.Projects(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Results, 2)
	assert.Equal(t, "Project 1", result.Results[0].Name)
}

func TestQualityprofiles_Projects_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.Projects(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Key
	_, _, err = client.Qualityprofiles.Projects(context.Background(), &QualityprofilesProjectsOptions{})
	assert.Error(t, err)

	// Test invalid Selected
	_, _, err = client.Qualityprofiles.Projects(context.Background(), &QualityprofilesProjectsOptions{
		Key:      "profile-key",
		Selected: "invalid",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_RemoveGroup(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/remove_group", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesRemoveGroupOptions{
		Group:          "sonar-administrators",
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.RemoveGroup(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_RemoveGroup_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.RemoveGroup(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Group
	_, err = client.Qualityprofiles.RemoveGroup(context.Background(), &QualityprofilesRemoveGroupOptions{
		Language:       "java",
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_RemoveProject(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/remove_project", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesRemoveProjectOptions{
		Language:       "java",
		Project:        "my_project",
		QualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.RemoveProject(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_RemoveProject_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.RemoveProject(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Language
	_, err = client.Qualityprofiles.RemoveProject(context.Background(), &QualityprofilesRemoveProjectOptions{
		Project:        "my_project",
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_RemoveUser(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/remove_user", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesRemoveUserOptions{
		Language:       "java",
		Login:          "john.doe",
		QualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.RemoveUser(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_RemoveUser_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.RemoveUser(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Language
	_, err = client.Qualityprofiles.RemoveUser(context.Background(), &QualityprofilesRemoveUserOptions{
		Login:          "john.doe",
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_Rename(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/rename", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesRenameOptions{
		Key:  "profile-key",
		Name: "New Profile Name",
	}

	resp, err := client.Qualityprofiles.Rename(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_Rename_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.Rename(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Key
	_, err = client.Qualityprofiles.Rename(context.Background(), &QualityprofilesRenameOptions{
		Name: "New Name",
	})
	assert.Error(t, err)

	// Test missing Name
	_, err = client.Qualityprofiles.Rename(context.Background(), &QualityprofilesRenameOptions{
		Key: "profile-key",
	})
	assert.Error(t, err)

	// Test Name too long
	_, err = client.Qualityprofiles.Rename(context.Background(), &QualityprofilesRenameOptions{
		Key:  "profile-key",
		Name: strings.Repeat("a", MaxQualityProfileNameLength+1),
	})
	assert.Error(t, err)
}

func TestQualityprofiles_Restore(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/restore", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesRestoreOptions{
		Backup: `<?xml version='1.0'?><profile><name>My Profile</name></profile>`,
	}

	resp, err := client.Qualityprofiles.Restore(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_Restore_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.Restore(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Backup
	_, err = client.Qualityprofiles.Restore(context.Background(), &QualityprofilesRestoreOptions{})
	assert.Error(t, err)
}

func TestQualityprofiles_Search(t *testing.T) {
	response := &QualityprofilesSearch{
		Actions: QualityprofilesActions{Create: true},
		Profiles: []QualityProfile{
			{
				Key:             "sonar-way-java",
				Name:            "Sonar way",
				Language:        "java",
				LanguageName:    "Java",
				IsDefault:       true,
				IsBuiltIn:       true,
				ActiveRuleCount: 200,
			},
			{
				Key:             "my-profile-java",
				Name:            "My Profile",
				Language:        "java",
				LanguageName:    "Java",
				IsDefault:       false,
				IsBuiltIn:       false,
				ActiveRuleCount: 150,
			},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualityprofiles/search", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesSearchOptions{
		Language: "java",
	}

	result, resp, err := client.Qualityprofiles.Search(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Profiles, 2)
	assert.Equal(t, "Sonar way", result.Profiles[0].Name)
	assert.True(t, result.Actions.Create)
}

func TestQualityprofiles_Search_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.Search(context.Background(), nil)
	assert.Error(t, err)
}

func TestQualityprofiles_SearchGroups(t *testing.T) {
	response := &QualityprofilesSearchGroups{
		Paging: Paging{PageIndex: 1, PageSize: 25, Total: 1},
		Groups: []QualityprofilesProfileGroup{
			{Name: "sonar-administrators", Description: "Admin group", Selected: true},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualityprofiles/search_groups", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesSearchGroupsOptions{
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	result, resp, err := client.Qualityprofiles.SearchGroups(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Groups, 1)
	assert.Equal(t, "sonar-administrators", result.Groups[0].Name)
}

func TestQualityprofiles_SearchGroups_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.SearchGroups(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Language
	_, _, err = client.Qualityprofiles.SearchGroups(context.Background(), &QualityprofilesSearchGroupsOptions{
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, _, err = client.Qualityprofiles.SearchGroups(context.Background(), &QualityprofilesSearchGroupsOptions{
		Language: "java",
	})
	assert.Error(t, err)

	// Test invalid Selected
	_, _, err = client.Qualityprofiles.SearchGroups(context.Background(), &QualityprofilesSearchGroupsOptions{
		Language:       "java",
		QualityProfile: "Sonar way",
		Selected:       "invalid",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_SearchUsers(t *testing.T) {
	response := &QualityprofilesSearchUsers{
		Paging: Paging{PageIndex: 1, PageSize: 25, Total: 1},
		Users: []QualityprofilesProfileUser{
			{Login: "john.doe", Name: "John Doe", Selected: true},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualityprofiles/search_users", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesSearchUsersOptions{
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	result, resp, err := client.Qualityprofiles.SearchUsers(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Users, 1)
	assert.Equal(t, "john.doe", result.Users[0].Login)
}

func TestQualityprofiles_SearchUsers_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.SearchUsers(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Language
	_, _, err = client.Qualityprofiles.SearchUsers(context.Background(), &QualityprofilesSearchUsersOptions{
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, _, err = client.Qualityprofiles.SearchUsers(context.Background(), &QualityprofilesSearchUsersOptions{
		Language: "java",
	})
	assert.Error(t, err)

	// Test invalid Selected
	_, _, err = client.Qualityprofiles.SearchUsers(context.Background(), &QualityprofilesSearchUsersOptions{
		Language:       "java",
		QualityProfile: "Sonar way",
		Selected:       "invalid",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_SetDefault(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/set_default", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesSetDefaultOptions{
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.SetDefault(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_SetDefault_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.SetDefault(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Language
	_, err = client.Qualityprofiles.SetDefault(context.Background(), &QualityprofilesSetDefaultOptions{
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, err = client.Qualityprofiles.SetDefault(context.Background(), &QualityprofilesSetDefaultOptions{
		Language: "java",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_Show(t *testing.T) {
	response := &QualityprofilesShow{
		Profile: QualityprofilesShownProfile{
			Key:             "sonar-way-java",
			Name:            "Sonar way",
			Language:        "java",
			LanguageName:    "Java",
			IsDefault:       true,
			IsBuiltIn:       true,
			ActiveRuleCount: 200,
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualityprofiles/show", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesShowOptions{
		Key: "sonar-way-java",
	}

	result, resp, err := client.Qualityprofiles.Show(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "Sonar way", result.Profile.Name)
	assert.Equal(t, int64(200), result.Profile.ActiveRuleCount)
}

func TestQualityprofiles_Show_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.Show(context.Background(), nil)
	assert.Error(t, err)

	// Test missing Key
	_, _, err = client.Qualityprofiles.Show(context.Background(), &QualityprofilesShowOptions{})
	assert.Error(t, err)
}

func TestQualityprofiles_ConvertActivateRuleOptForURL(t *testing.T) {
	client := newLocalhostClient(t)

	opt := &QualityprofilesActivateRuleOptions{
		Key:             "AU-TpxcA-iU5OvuD2FL0",
		Rule:            "squid:AvoidCycles",
		Impacts:         map[string]string{SoftwareQualityMaintainability: RuleImpactSeverityHigh, SoftwareQualitySecurity: RuleImpactSeverityMedium},
		Params:          map[string]string{"max": "10", "threshold": "5"},
		PrioritizedRule: true,
		Reset:           false,
		Severity:        "",
	}

	urlOpt := client.Qualityprofiles.convertActivateRuleOptForURL(opt)

	// Verify basic fields
	assert.Equal(t, opt.Key, urlOpt.Key)
	assert.Equal(t, opt.Rule, urlOpt.Rule)
	assert.Equal(t, opt.PrioritizedRule, urlOpt.PrioritizedRule)

	// Verify map conversions
	assert.NotEmpty(t, urlOpt.Impacts)

	// Impacts string should contain both entries
	assert.Contains(t, urlOpt.Impacts, "MAINTAINABILITY=HIGH")
	assert.Contains(t, urlOpt.Impacts, "SECURITY=MEDIUM")

	assert.NotEmpty(t, urlOpt.Params)

	// Params string should contain both entries
	assert.Contains(t, urlOpt.Params, "max=10")
	assert.Contains(t, urlOpt.Params, "threshold=5")
}

func TestQualityprofilesService_ChangelogAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		callCount := 0
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			callCount++
			w.Header().Set("Content-Type", "application/json")
			if callCount == 1 {
				_, _ = w.Write([]byte(`{"paging":{"pageIndex":1,"pageSize":500,"total":2},"events":[{"action":"ACTIVATED"}]}`))
			} else {
				_, _ = w.Write([]byte(`{"paging":{"pageIndex":2,"pageSize":500,"total":2},"events":[{"action":"DEACTIVATED"}]}`))
			}
		})

		client := newTestClient(t, server.URL)
		opt := &QualityprofilesChangelogOptions{Language: "java", QualityProfile: "Sonar way"}
		result, _, err := client.Qualityprofiles.ChangelogAll(context.Background(), opt)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, 2, callCount)
	})

	t.Run("nil option", func(t *testing.T) {
		client := newLocalhostClient(t)
		_, _, err := client.Qualityprofiles.ChangelogAll(context.Background(), nil)
		assert.Error(t, err)
	})
}

func TestQualityprofilesService_ProjectsAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		callCount := 0
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			callCount++
			w.Header().Set("Content-Type", "application/json")
			if callCount == 1 {
				_, _ = w.Write([]byte(`{"paging":{"pageIndex":1,"pageSize":500,"total":2},"results":[{"key":"p1"}]}`))
			} else {
				_, _ = w.Write([]byte(`{"paging":{"pageIndex":2,"pageSize":500,"total":2},"results":[{"key":"p2"}]}`))
			}
		})

		client := newTestClient(t, server.URL)
		opt := &QualityprofilesProjectsOptions{Key: "myprofile"}
		result, _, err := client.Qualityprofiles.ProjectsAll(context.Background(), opt)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, 2, callCount)
	})

	t.Run("nil option", func(t *testing.T) {
		client := newLocalhostClient(t)
		_, _, err := client.Qualityprofiles.ProjectsAll(context.Background(), nil)
		assert.Error(t, err)
	})
}

func TestQualityprofilesService_SearchGroupsAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		callCount := 0
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			callCount++
			w.Header().Set("Content-Type", "application/json")
			if callCount == 1 {
				_, _ = w.Write([]byte(`{"paging":{"pageIndex":1,"pageSize":500,"total":2},"groups":[{"name":"g1"}]}`))
			} else {
				_, _ = w.Write([]byte(`{"paging":{"pageIndex":2,"pageSize":500,"total":2},"groups":[{"name":"g2"}]}`))
			}
		})

		client := newTestClient(t, server.URL)
		opt := &QualityprofilesSearchGroupsOptions{Language: "java", QualityProfile: "Sonar way"}
		result, _, err := client.Qualityprofiles.SearchGroupsAll(context.Background(), opt)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, 2, callCount)
	})

	t.Run("nil option", func(t *testing.T) {
		client := newLocalhostClient(t)
		_, _, err := client.Qualityprofiles.SearchGroupsAll(context.Background(), nil)
		assert.Error(t, err)
	})
}

func TestQualityprofilesService_SearchUsersAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		callCount := 0
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			callCount++
			w.Header().Set("Content-Type", "application/json")
			if callCount == 1 {
				_, _ = w.Write([]byte(`{"paging":{"pageIndex":1,"pageSize":500,"total":2},"users":[{"login":"u1"}]}`))
			} else {
				_, _ = w.Write([]byte(`{"paging":{"pageIndex":2,"pageSize":500,"total":2},"users":[{"login":"u2"}]}`))
			}
		})

		client := newTestClient(t, server.URL)
		opt := &QualityprofilesSearchUsersOptions{Language: "java", QualityProfile: "Sonar way"}
		result, _, err := client.Qualityprofiles.SearchUsersAll(context.Background(), opt)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, 2, callCount)
	})

	t.Run("nil option", func(t *testing.T) {
		client := newLocalhostClient(t)
		_, _, err := client.Qualityprofiles.SearchUsersAll(context.Background(), nil)
		assert.Error(t, err)
	})
}
