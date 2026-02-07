package sonar

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQualityprofiles_ActivateRule(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/activate_rule", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesActivateRuleOption{
		Key:  "AU-TpxcA-iU5OvuD2FL0",
		Rule: "squid:AvoidCycles",
	}

	resp, err := client.Qualityprofiles.ActivateRule(opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_ActivateRule_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.ActivateRule(nil)
	assert.Error(t, err)

	// Test missing Key
	_, err = client.Qualityprofiles.ActivateRule(&QualityprofilesActivateRuleOption{
		Rule: "squid:AvoidCycles",
	})
	assert.Error(t, err)

	// Test missing Rule
	_, err = client.Qualityprofiles.ActivateRule(&QualityprofilesActivateRuleOption{
		Key: "AU-TpxcA-iU5OvuD2FL0",
	})
	assert.Error(t, err)

	// Test both Impacts and Severity set
	_, err = client.Qualityprofiles.ActivateRule(&QualityprofilesActivateRuleOption{
		Key:      "AU-TpxcA-iU5OvuD2FL0",
		Rule:     "squid:AvoidCycles",
		Impacts:  map[string]string{"MAINTAINABILITY": "HIGH"},
		Severity: "MAJOR",
	})
	assert.Error(t, err)

	// Test invalid Severity
	_, err = client.Qualityprofiles.ActivateRule(&QualityprofilesActivateRuleOption{
		Key:      "AU-TpxcA-iU5OvuD2FL0",
		Rule:     "squid:AvoidCycles",
		Severity: "INVALID",
	})
	assert.Error(t, err)

	// Test invalid Impacts map key (software quality)
	_, err = client.Qualityprofiles.ActivateRule(&QualityprofilesActivateRuleOption{
		Key:     "AU-TpxcA-iU5OvuD2FL0",
		Rule:    "squid:AvoidCycles",
		Impacts: map[string]string{"INVALID_QUALITY": "HIGH"},
	})
	assert.Error(t, err)

	// Test invalid Impacts map value (severity)
	_, err = client.Qualityprofiles.ActivateRule(&QualityprofilesActivateRuleOption{
		Key:     "AU-TpxcA-iU5OvuD2FL0",
		Rule:    "squid:AvoidCycles",
		Impacts: map[string]string{"MAINTAINABILITY": "INVALID_SEVERITY"},
	})
	assert.Error(t, err)
}

func TestQualityprofiles_ActivateRules(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/activate_rules", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesActivateRulesOption{
		TargetKey: "AU-TpxcA-iU5OvuD2FL0",
		Languages: []string{"java"},
	}

	resp, err := client.Qualityprofiles.ActivateRules(opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_ActivateRules_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.ActivateRules(nil)
	assert.Error(t, err)

	// Test missing TargetKey
	_, err = client.Qualityprofiles.ActivateRules(&QualityprofilesActivateRulesOption{})
	assert.Error(t, err)

	// Test invalid language
	_, err = client.Qualityprofiles.ActivateRules(&QualityprofilesActivateRulesOption{
		TargetKey: "AU-TpxcA-iU5OvuD2FL0",
		Languages: []string{"invalid_language"},
	})
	assert.Error(t, err)

	// Test invalid severity
	_, err = client.Qualityprofiles.ActivateRules(&QualityprofilesActivateRulesOption{
		TargetKey:  "AU-TpxcA-iU5OvuD2FL0",
		Severities: []string{"INVALID_SEVERITY"},
	})
	assert.Error(t, err)

	// Test invalid impact severity
	_, err = client.Qualityprofiles.ActivateRules(&QualityprofilesActivateRulesOption{
		TargetKey:        "AU-TpxcA-iU5OvuD2FL0",
		ImpactSeverities: []string{"INVALID"},
	})
	assert.Error(t, err)

	// Test invalid software quality
	_, err = client.Qualityprofiles.ActivateRules(&QualityprofilesActivateRulesOption{
		TargetKey:               "AU-TpxcA-iU5OvuD2FL0",
		ImpactSoftwareQualities: []string{"INVALID"},
	})
	assert.Error(t, err)

	// Test invalid sort field
	_, err = client.Qualityprofiles.ActivateRules(&QualityprofilesActivateRulesOption{
		TargetKey: "AU-TpxcA-iU5OvuD2FL0",
		Sort:      "invalid_sort",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_AddGroup(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/add_group", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesAddGroupOption{
		Group:          "sonar-administrators",
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.AddGroup(opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_AddGroup_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.AddGroup(nil)
	assert.Error(t, err)

	// Test missing Group
	_, err = client.Qualityprofiles.AddGroup(&QualityprofilesAddGroupOption{
		Language:       "java",
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing Language
	_, err = client.Qualityprofiles.AddGroup(&QualityprofilesAddGroupOption{
		Group:          "sonar-administrators",
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test invalid Language
	_, err = client.Qualityprofiles.AddGroup(&QualityprofilesAddGroupOption{
		Group:          "sonar-administrators",
		Language:       "invalid_lang",
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, err = client.Qualityprofiles.AddGroup(&QualityprofilesAddGroupOption{
		Group:    "sonar-administrators",
		Language: "java",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_AddProject(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/add_project", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesAddProjectOption{
		Language:       "java",
		Project:        "my_project",
		QualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.AddProject(opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_AddProject_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.AddProject(nil)
	assert.Error(t, err)

	// Test missing Language
	_, err = client.Qualityprofiles.AddProject(&QualityprofilesAddProjectOption{
		Project:        "my_project",
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing Project
	_, err = client.Qualityprofiles.AddProject(&QualityprofilesAddProjectOption{
		Language:       "java",
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, err = client.Qualityprofiles.AddProject(&QualityprofilesAddProjectOption{
		Language: "java",
		Project:  "my_project",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_AddUser(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/add_user", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesAddUserOption{
		Language:       "java",
		Login:          "john.doe",
		QualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.AddUser(opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_AddUser_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.AddUser(nil)
	assert.Error(t, err)

	// Test missing Language
	_, err = client.Qualityprofiles.AddUser(&QualityprofilesAddUserOption{
		Login:          "john.doe",
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing Login
	_, err = client.Qualityprofiles.AddUser(&QualityprofilesAddUserOption{
		Language:       "java",
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, err = client.Qualityprofiles.AddUser(&QualityprofilesAddUserOption{
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

	opt := &QualityprofilesBackupOption{
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	result, resp, err := client.Qualityprofiles.Backup(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Contains(t, *result, "Sonar way")
}

func TestQualityprofiles_Backup_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.Backup(nil)
	assert.Error(t, err)

	// Test missing Language
	_, _, err = client.Qualityprofiles.Backup(&QualityprofilesBackupOption{
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, _, err = client.Qualityprofiles.Backup(&QualityprofilesBackupOption{
		Language: "java",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_ChangeParent(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/change_parent", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesChangeParentOption{
		Language:             "java",
		QualityProfile:       "My Profile",
		ParentQualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.ChangeParent(opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_ChangeParent_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.ChangeParent(nil)
	assert.Error(t, err)

	// Test missing Language
	_, err = client.Qualityprofiles.ChangeParent(&QualityprofilesChangeParentOption{
		QualityProfile: "My Profile",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, err = client.Qualityprofiles.ChangeParent(&QualityprofilesChangeParentOption{
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

	opt := &QualityprofilesChangelogOption{
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	result, resp, err := client.Qualityprofiles.Changelog(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Events, 1)
	assert.Equal(t, "ACTIVATED", result.Events[0].Action)
}

func TestQualityprofiles_Changelog_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.Changelog(nil)
	assert.Error(t, err)

	// Test missing Language
	_, _, err = client.Qualityprofiles.Changelog(&QualityprofilesChangelogOption{
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, _, err = client.Qualityprofiles.Changelog(&QualityprofilesChangelogOption{
		Language: "java",
	})
	assert.Error(t, err)

	// Test invalid FilterMode
	_, _, err = client.Qualityprofiles.Changelog(&QualityprofilesChangelogOption{
		Language:       "java",
		QualityProfile: "Sonar way",
		FilterMode:     "INVALID",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_Compare(t *testing.T) {
	response := &QualityprofilesCompare{
		Left:  CompareProfile{Key: "profile1", Name: "Profile 1"},
		Right: CompareProfile{Key: "profile2", Name: "Profile 2"},
		InLeft: []CompareRule{
			{Key: "squid:S1234", Name: "Rule in left"},
		},
		InRight: []CompareRule{
			{Key: "squid:S5678", Name: "Rule in right"},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualityprofiles/compare", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesCompareOption{
		LeftKey:  "profile1",
		RightKey: "profile2",
	}

	result, resp, err := client.Qualityprofiles.Compare(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "Profile 1", result.Left.Name)
	assert.Len(t, result.InLeft, 1)
}

func TestQualityprofiles_Compare_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.Compare(nil)
	assert.Error(t, err)

	// Test missing LeftKey
	_, _, err = client.Qualityprofiles.Compare(&QualityprofilesCompareOption{
		RightKey: "profile2",
	})
	assert.Error(t, err)

	// Test missing RightKey
	_, _, err = client.Qualityprofiles.Compare(&QualityprofilesCompareOption{
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

	opt := &QualityprofilesCopyOption{
		FromKey: "source-profile-key",
		ToName:  "My Profile Copy",
	}

	result, resp, err := client.Qualityprofiles.Copy(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "My Profile Copy", result.Name)
}

func TestQualityprofiles_Copy_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.Copy(nil)
	assert.Error(t, err)

	// Test missing FromKey
	_, _, err = client.Qualityprofiles.Copy(&QualityprofilesCopyOption{
		ToName: "New Profile",
	})
	assert.Error(t, err)

	// Test missing ToName
	_, _, err = client.Qualityprofiles.Copy(&QualityprofilesCopyOption{
		FromKey: "source-key",
	})
	assert.Error(t, err)

	// Test ToName too long
	_, _, err = client.Qualityprofiles.Copy(&QualityprofilesCopyOption{
		FromKey: "source-key",
		ToName:  strings.Repeat("a", MaxQualityProfileNameLength+1),
	})
	assert.Error(t, err)
}

func TestQualityprofiles_Create(t *testing.T) {
	response := &QualityprofilesCreate{
		Profile: CreatedProfile{
			Key:          "new-profile-key",
			Name:         "My New Profile",
			Language:     "java",
			LanguageName: "Java",
			IsDefault:    false,
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodPost, "/qualityprofiles/create", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesCreateOption{
		Language: "java",
		Name:     "My New Profile",
	}

	result, resp, err := client.Qualityprofiles.Create(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "My New Profile", result.Profile.Name)
}

func TestQualityprofiles_Create_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.Create(nil)
	assert.Error(t, err)

	// Test missing Language
	_, _, err = client.Qualityprofiles.Create(&QualityprofilesCreateOption{
		Name: "My Profile",
	})
	assert.Error(t, err)

	// Test missing Name
	_, _, err = client.Qualityprofiles.Create(&QualityprofilesCreateOption{
		Language: "java",
	})
	assert.Error(t, err)

	// Test Name too long
	_, _, err = client.Qualityprofiles.Create(&QualityprofilesCreateOption{
		Language: "java",
		Name:     strings.Repeat("a", MaxQualityProfileNameLength+1),
	})
	assert.Error(t, err)
}

func TestQualityprofiles_DeactivateRule(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/deactivate_rule", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesDeactivateRuleOption{
		Key:  "AU-TpxcA-iU5OvuD2FL0",
		Rule: "squid:AvoidCycles",
	}

	resp, err := client.Qualityprofiles.DeactivateRule(opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_DeactivateRule_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.DeactivateRule(nil)
	assert.Error(t, err)

	// Test missing Key
	_, err = client.Qualityprofiles.DeactivateRule(&QualityprofilesDeactivateRuleOption{
		Rule: "squid:AvoidCycles",
	})
	assert.Error(t, err)

	// Test missing Rule
	_, err = client.Qualityprofiles.DeactivateRule(&QualityprofilesDeactivateRuleOption{
		Key: "AU-TpxcA-iU5OvuD2FL0",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_DeactivateRules(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/deactivate_rules", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesDeactivateRulesOption{
		TargetKey: "AU-TpxcA-iU5OvuD2FL0",
		Languages: []string{"java"},
	}

	resp, err := client.Qualityprofiles.DeactivateRules(opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_DeactivateRules_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.DeactivateRules(nil)
	assert.Error(t, err)

	// Test missing TargetKey
	_, err = client.Qualityprofiles.DeactivateRules(&QualityprofilesDeactivateRulesOption{})
	assert.Error(t, err)
}

func TestQualityprofiles_Delete(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/delete", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesDeleteOption{
		Language:       "java",
		QualityProfile: "My Profile",
	}

	resp, err := client.Qualityprofiles.Delete(opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_Delete_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.Delete(nil)
	assert.Error(t, err)

	// Test missing Language
	_, err = client.Qualityprofiles.Delete(&QualityprofilesDeleteOption{
		QualityProfile: "My Profile",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, err = client.Qualityprofiles.Delete(&QualityprofilesDeleteOption{
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

	opt := &QualityprofilesExportOption{
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	result, resp, err := client.Qualityprofiles.Export(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Contains(t, *result, "Sonar way")
}

func TestQualityprofiles_Export_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.Export(nil)
	assert.Error(t, err)

	// Test missing Language
	_, _, err = client.Qualityprofiles.Export(&QualityprofilesExportOption{})
	assert.Error(t, err)
}

func TestQualityprofiles_Exporters(t *testing.T) {
	response := &QualityprofilesExporters{
		Exporters: []ProfileExporter{
			{Key: "checkstyle", Name: "Checkstyle", Languages: []string{"java"}},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualityprofiles/exporters", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Qualityprofiles.Exporters()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Exporters, 1)
	assert.Equal(t, "checkstyle", result.Exporters[0].Key)
}

func TestQualityprofiles_Importers(t *testing.T) {
	response := &QualityprofilesImporters{
		Importers: []ProfileImporter{
			{Key: "pmd", Name: "PMD", Languages: []string{"java"}},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualityprofiles/importers", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.Qualityprofiles.Importers()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Importers, 1)
	assert.Equal(t, "pmd", result.Importers[0].Key)
}

func TestQualityprofiles_Inheritance(t *testing.T) {
	response := &QualityprofilesInheritance{
		Profile: InheritanceProfile{
			Key:             "my-profile-key",
			Name:            "My Profile",
			ActiveRuleCount: 150,
		},
		Ancestors: []InheritanceProfile{
			{Key: "parent-key", Name: "Sonar way", ActiveRuleCount: 200},
		},
		Children: []InheritanceProfile{},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualityprofiles/inheritance", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesInheritanceOption{
		Language:       "java",
		QualityProfile: "My Profile",
	}

	result, resp, err := client.Qualityprofiles.Inheritance(opt)
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
	_, _, err := client.Qualityprofiles.Inheritance(nil)
	assert.Error(t, err)

	// Test missing Language
	_, _, err = client.Qualityprofiles.Inheritance(&QualityprofilesInheritanceOption{
		QualityProfile: "My Profile",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, _, err = client.Qualityprofiles.Inheritance(&QualityprofilesInheritanceOption{
		Language: "java",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_Projects(t *testing.T) {
	response := &QualityprofilesProjects{
		Paging: Paging{PageIndex: 1, PageSize: 25, Total: 2},
		Results: []ProfileProject{
			{Key: "project1", Name: "Project 1", Selected: true},
			{Key: "project2", Name: "Project 2", Selected: false},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualityprofiles/projects", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesProjectsOption{
		Key: "profile-key",
	}

	result, resp, err := client.Qualityprofiles.Projects(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Results, 2)
	assert.Equal(t, "Project 1", result.Results[0].Name)
}

func TestQualityprofiles_Projects_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.Projects(nil)
	assert.Error(t, err)

	// Test missing Key
	_, _, err = client.Qualityprofiles.Projects(&QualityprofilesProjectsOption{})
	assert.Error(t, err)

	// Test invalid Selected
	_, _, err = client.Qualityprofiles.Projects(&QualityprofilesProjectsOption{
		Key:      "profile-key",
		Selected: "invalid",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_RemoveGroup(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/remove_group", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesRemoveGroupOption{
		Group:          "sonar-administrators",
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.RemoveGroup(opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_RemoveGroup_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.RemoveGroup(nil)
	assert.Error(t, err)

	// Test missing Group
	_, err = client.Qualityprofiles.RemoveGroup(&QualityprofilesRemoveGroupOption{
		Language:       "java",
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_RemoveProject(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/remove_project", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesRemoveProjectOption{
		Language:       "java",
		Project:        "my_project",
		QualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.RemoveProject(opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_RemoveProject_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.RemoveProject(nil)
	assert.Error(t, err)

	// Test missing Language
	_, err = client.Qualityprofiles.RemoveProject(&QualityprofilesRemoveProjectOption{
		Project:        "my_project",
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_RemoveUser(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/remove_user", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesRemoveUserOption{
		Language:       "java",
		Login:          "john.doe",
		QualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.RemoveUser(opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_RemoveUser_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.RemoveUser(nil)
	assert.Error(t, err)

	// Test missing Language
	_, err = client.Qualityprofiles.RemoveUser(&QualityprofilesRemoveUserOption{
		Login:          "john.doe",
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_Rename(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/rename", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesRenameOption{
		Key:  "profile-key",
		Name: "New Profile Name",
	}

	resp, err := client.Qualityprofiles.Rename(opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_Rename_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.Rename(nil)
	assert.Error(t, err)

	// Test missing Key
	_, err = client.Qualityprofiles.Rename(&QualityprofilesRenameOption{
		Name: "New Name",
	})
	assert.Error(t, err)

	// Test missing Name
	_, err = client.Qualityprofiles.Rename(&QualityprofilesRenameOption{
		Key: "profile-key",
	})
	assert.Error(t, err)

	// Test Name too long
	_, err = client.Qualityprofiles.Rename(&QualityprofilesRenameOption{
		Key:  "profile-key",
		Name: strings.Repeat("a", MaxQualityProfileNameLength+1),
	})
	assert.Error(t, err)
}

func TestQualityprofiles_Restore(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/restore", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesRestoreOption{
		Backup: `<?xml version='1.0'?><profile><name>My Profile</name></profile>`,
	}

	resp, err := client.Qualityprofiles.Restore(opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_Restore_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.Restore(nil)
	assert.Error(t, err)

	// Test missing Backup
	_, err = client.Qualityprofiles.Restore(&QualityprofilesRestoreOption{})
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

	opt := &QualityprofilesSearchOption{
		Language: "java",
	}

	result, resp, err := client.Qualityprofiles.Search(opt)
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
	_, _, err := client.Qualityprofiles.Search(nil)
	assert.Error(t, err)
}

func TestQualityprofiles_SearchGroups(t *testing.T) {
	response := &QualityprofilesSearchGroups{
		Paging: Paging{PageIndex: 1, PageSize: 25, Total: 1},
		Groups: []ProfileGroup{
			{Name: "sonar-administrators", Description: "Admin group", Selected: true},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualityprofiles/search_groups", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesSearchGroupsOption{
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	result, resp, err := client.Qualityprofiles.SearchGroups(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Groups, 1)
	assert.Equal(t, "sonar-administrators", result.Groups[0].Name)
}

func TestQualityprofiles_SearchGroups_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.SearchGroups(nil)
	assert.Error(t, err)

	// Test missing Language
	_, _, err = client.Qualityprofiles.SearchGroups(&QualityprofilesSearchGroupsOption{
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, _, err = client.Qualityprofiles.SearchGroups(&QualityprofilesSearchGroupsOption{
		Language: "java",
	})
	assert.Error(t, err)

	// Test invalid Selected
	_, _, err = client.Qualityprofiles.SearchGroups(&QualityprofilesSearchGroupsOption{
		Language:       "java",
		QualityProfile: "Sonar way",
		Selected:       "invalid",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_SearchUsers(t *testing.T) {
	response := &QualityprofilesSearchUsers{
		Paging: Paging{PageIndex: 1, PageSize: 25, Total: 1},
		Users: []ProfileUser{
			{Login: "john.doe", Name: "John Doe", Selected: true},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/qualityprofiles/search_users", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesSearchUsersOption{
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	result, resp, err := client.Qualityprofiles.SearchUsers(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Users, 1)
	assert.Equal(t, "john.doe", result.Users[0].Login)
}

func TestQualityprofiles_SearchUsers_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.SearchUsers(nil)
	assert.Error(t, err)

	// Test missing Language
	_, _, err = client.Qualityprofiles.SearchUsers(&QualityprofilesSearchUsersOption{
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, _, err = client.Qualityprofiles.SearchUsers(&QualityprofilesSearchUsersOption{
		Language: "java",
	})
	assert.Error(t, err)

	// Test invalid Selected
	_, _, err = client.Qualityprofiles.SearchUsers(&QualityprofilesSearchUsersOption{
		Language:       "java",
		QualityProfile: "Sonar way",
		Selected:       "invalid",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_SetDefault(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/qualityprofiles/set_default", 204))
	client := newTestClient(t, server.URL)

	opt := &QualityprofilesSetDefaultOption{
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.SetDefault(opt)
	require.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}

func TestQualityprofiles_SetDefault_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, err := client.Qualityprofiles.SetDefault(nil)
	assert.Error(t, err)

	// Test missing Language
	_, err = client.Qualityprofiles.SetDefault(&QualityprofilesSetDefaultOption{
		QualityProfile: "Sonar way",
	})
	assert.Error(t, err)

	// Test missing QualityProfile
	_, err = client.Qualityprofiles.SetDefault(&QualityprofilesSetDefaultOption{
		Language: "java",
	})
	assert.Error(t, err)
}

func TestQualityprofiles_Show(t *testing.T) {
	response := &QualityprofilesShow{
		Profile: ShownProfile{
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

	opt := &QualityprofilesShowOption{
		Key: "sonar-way-java",
	}

	result, resp, err := client.Qualityprofiles.Show(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "Sonar way", result.Profile.Name)
	assert.Equal(t, int64(200), result.Profile.ActiveRuleCount)
}

func TestQualityprofiles_Show_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test nil option
	_, _, err := client.Qualityprofiles.Show(nil)
	assert.Error(t, err)

	// Test missing Key
	_, _, err = client.Qualityprofiles.Show(&QualityprofilesShowOption{})
	assert.Error(t, err)
}

func TestQualityprofiles_ConvertActivateRuleOptForURL(t *testing.T) {
	client := newLocalhostClient(t)

	opt := &QualityprofilesActivateRuleOption{
		Key:             "AU-TpxcA-iU5OvuD2FL0",
		Rule:            "squid:AvoidCycles",
		Impacts:         map[string]string{"MAINTAINABILITY": "HIGH", "SECURITY": "MEDIUM"},
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
