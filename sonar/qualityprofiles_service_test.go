package sonargo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestQualityprofiles_ActivateRule(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if !strings.Contains(r.URL.Path, "qualityprofiles/activate_rule") {
			t.Errorf("expected path to contain qualityprofiles/activate_rule, got %s", r.URL.Path)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesActivateRuleOption{
		Key:  "AU-TpxcA-iU5OvuD2FL0",
		Rule: "squid:AvoidCycles",
	}

	resp, err := client.Qualityprofiles.ActivateRule(opt)
	if err != nil {
		t.Fatalf("ActivateRule failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityprofiles_ActivateRule_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.Qualityprofiles.ActivateRule(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Key
	_, err = client.Qualityprofiles.ActivateRule(&QualityprofilesActivateRuleOption{
		Rule: "squid:AvoidCycles",
	})
	if err == nil {
		t.Error("expected error for missing Key")
	}

	// Test missing Rule
	_, err = client.Qualityprofiles.ActivateRule(&QualityprofilesActivateRuleOption{
		Key: "AU-TpxcA-iU5OvuD2FL0",
	})
	if err == nil {
		t.Error("expected error for missing Rule")
	}

	// Test both Impacts and Severity set
	_, err = client.Qualityprofiles.ActivateRule(&QualityprofilesActivateRuleOption{
		Key:      "AU-TpxcA-iU5OvuD2FL0",
		Rule:     "squid:AvoidCycles",
		Impacts:  map[string]string{"MAINTAINABILITY": "HIGH"},
		Severity: "MAJOR",
	})
	if err == nil {
		t.Error("expected error when both Impacts and Severity are set")
	}

	// Test invalid Severity
	_, err = client.Qualityprofiles.ActivateRule(&QualityprofilesActivateRuleOption{
		Key:      "AU-TpxcA-iU5OvuD2FL0",
		Rule:     "squid:AvoidCycles",
		Severity: "INVALID",
	})
	if err == nil {
		t.Error("expected error for invalid Severity")
	}

	// Test invalid Impacts map key (software quality)
	_, err = client.Qualityprofiles.ActivateRule(&QualityprofilesActivateRuleOption{
		Key:     "AU-TpxcA-iU5OvuD2FL0",
		Rule:    "squid:AvoidCycles",
		Impacts: map[string]string{"INVALID_QUALITY": "HIGH"},
	})
	if err == nil {
		t.Error("expected error for invalid Impacts map key")
	}

	// Test invalid Impacts map value (severity)
	_, err = client.Qualityprofiles.ActivateRule(&QualityprofilesActivateRuleOption{
		Key:     "AU-TpxcA-iU5OvuD2FL0",
		Rule:    "squid:AvoidCycles",
		Impacts: map[string]string{"MAINTAINABILITY": "INVALID_SEVERITY"},
	})
	if err == nil {
		t.Error("expected error for invalid Impacts map value")
	}
}

func TestQualityprofiles_ActivateRules(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesActivateRulesOption{
		TargetKey: "AU-TpxcA-iU5OvuD2FL0",
		Languages: []string{"java"},
	}

	resp, err := client.Qualityprofiles.ActivateRules(opt)
	if err != nil {
		t.Fatalf("ActivateRules failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityprofiles_ActivateRules_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.Qualityprofiles.ActivateRules(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing TargetKey
	_, err = client.Qualityprofiles.ActivateRules(&QualityprofilesActivateRulesOption{})
	if err == nil {
		t.Error("expected error for missing TargetKey")
	}

	// Test invalid language
	_, err = client.Qualityprofiles.ActivateRules(&QualityprofilesActivateRulesOption{
		TargetKey: "AU-TpxcA-iU5OvuD2FL0",
		Languages: []string{"invalid_language"},
	})
	if err == nil {
		t.Error("expected error for invalid language")
	}

	// Test invalid severity
	_, err = client.Qualityprofiles.ActivateRules(&QualityprofilesActivateRulesOption{
		TargetKey:  "AU-TpxcA-iU5OvuD2FL0",
		Severities: []string{"INVALID_SEVERITY"},
	})
	if err == nil {
		t.Error("expected error for invalid severity")
	}

	// Test invalid impact severity
	_, err = client.Qualityprofiles.ActivateRules(&QualityprofilesActivateRulesOption{
		TargetKey:        "AU-TpxcA-iU5OvuD2FL0",
		ImpactSeverities: []string{"INVALID"},
	})
	if err == nil {
		t.Error("expected error for invalid impact severity")
	}

	// Test invalid software quality
	_, err = client.Qualityprofiles.ActivateRules(&QualityprofilesActivateRulesOption{
		TargetKey:               "AU-TpxcA-iU5OvuD2FL0",
		ImpactSoftwareQualities: []string{"INVALID"},
	})
	if err == nil {
		t.Error("expected error for invalid software quality")
	}

	// Test invalid sort field
	_, err = client.Qualityprofiles.ActivateRules(&QualityprofilesActivateRulesOption{
		TargetKey: "AU-TpxcA-iU5OvuD2FL0",
		Sort:      "invalid_sort",
	})
	if err == nil {
		t.Error("expected error for invalid sort field")
	}
}

func TestQualityprofiles_AddGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesAddGroupOption{
		Group:          "sonar-administrators",
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.AddGroup(opt)
	if err != nil {
		t.Fatalf("AddGroup failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityprofiles_AddGroup_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.Qualityprofiles.AddGroup(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Group
	_, err = client.Qualityprofiles.AddGroup(&QualityprofilesAddGroupOption{
		Language:       "java",
		QualityProfile: "Sonar way",
	})
	if err == nil {
		t.Error("expected error for missing Group")
	}

	// Test missing Language
	_, err = client.Qualityprofiles.AddGroup(&QualityprofilesAddGroupOption{
		Group:          "sonar-administrators",
		QualityProfile: "Sonar way",
	})
	if err == nil {
		t.Error("expected error for missing Language")
	}

	// Test invalid Language
	_, err = client.Qualityprofiles.AddGroup(&QualityprofilesAddGroupOption{
		Group:          "sonar-administrators",
		Language:       "invalid_lang",
		QualityProfile: "Sonar way",
	})
	if err == nil {
		t.Error("expected error for invalid Language")
	}

	// Test missing QualityProfile
	_, err = client.Qualityprofiles.AddGroup(&QualityprofilesAddGroupOption{
		Group:    "sonar-administrators",
		Language: "java",
	})
	if err == nil {
		t.Error("expected error for missing QualityProfile")
	}
}

func TestQualityprofiles_AddProject(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesAddProjectOption{
		Language:       "java",
		Project:        "my_project",
		QualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.AddProject(opt)
	if err != nil {
		t.Fatalf("AddProject failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityprofiles_AddProject_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.Qualityprofiles.AddProject(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Language
	_, err = client.Qualityprofiles.AddProject(&QualityprofilesAddProjectOption{
		Project:        "my_project",
		QualityProfile: "Sonar way",
	})
	if err == nil {
		t.Error("expected error for missing Language")
	}

	// Test missing Project
	_, err = client.Qualityprofiles.AddProject(&QualityprofilesAddProjectOption{
		Language:       "java",
		QualityProfile: "Sonar way",
	})
	if err == nil {
		t.Error("expected error for missing Project")
	}

	// Test missing QualityProfile
	_, err = client.Qualityprofiles.AddProject(&QualityprofilesAddProjectOption{
		Language: "java",
		Project:  "my_project",
	})
	if err == nil {
		t.Error("expected error for missing QualityProfile")
	}
}

func TestQualityprofiles_AddUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesAddUserOption{
		Language:       "java",
		Login:          "john.doe",
		QualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.AddUser(opt)
	if err != nil {
		t.Fatalf("AddUser failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityprofiles_AddUser_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.Qualityprofiles.AddUser(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Language
	_, err = client.Qualityprofiles.AddUser(&QualityprofilesAddUserOption{
		Login:          "john.doe",
		QualityProfile: "Sonar way",
	})
	if err == nil {
		t.Error("expected error for missing Language")
	}

	// Test missing Login
	_, err = client.Qualityprofiles.AddUser(&QualityprofilesAddUserOption{
		Language:       "java",
		QualityProfile: "Sonar way",
	})
	if err == nil {
		t.Error("expected error for missing Login")
	}

	// Test missing QualityProfile
	_, err = client.Qualityprofiles.AddUser(&QualityprofilesAddUserOption{
		Language: "java",
		Login:    "john.doe",
	})
	if err == nil {
		t.Error("expected error for missing QualityProfile")
	}
}

func TestQualityprofiles_Backup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(200)
		w.Write([]byte(`<?xml version='1.0'?><profile><name>Sonar way</name></profile>`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesBackupOption{
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	result, resp, err := client.Qualityprofiles.Backup(opt)
	if err != nil {
		t.Fatalf("Backup failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || !strings.Contains(*result, "Sonar way") {
		t.Error("expected backup XML to contain 'Sonar way'")
	}
}

func TestQualityprofiles_Backup_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.Qualityprofiles.Backup(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Language
	_, _, err = client.Qualityprofiles.Backup(&QualityprofilesBackupOption{
		QualityProfile: "Sonar way",
	})
	if err == nil {
		t.Error("expected error for missing Language")
	}

	// Test missing QualityProfile
	_, _, err = client.Qualityprofiles.Backup(&QualityprofilesBackupOption{
		Language: "java",
	})
	if err == nil {
		t.Error("expected error for missing QualityProfile")
	}
}

func TestQualityprofiles_ChangeParent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesChangeParentOption{
		Language:             "java",
		QualityProfile:       "My Profile",
		ParentQualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.ChangeParent(opt)
	if err != nil {
		t.Fatalf("ChangeParent failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityprofiles_ChangeParent_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.Qualityprofiles.ChangeParent(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Language
	_, err = client.Qualityprofiles.ChangeParent(&QualityprofilesChangeParentOption{
		QualityProfile: "My Profile",
	})
	if err == nil {
		t.Error("expected error for missing Language")
	}

	// Test missing QualityProfile
	_, err = client.Qualityprofiles.ChangeParent(&QualityprofilesChangeParentOption{
		Language: "java",
	})
	if err == nil {
		t.Error("expected error for missing QualityProfile")
	}
}

func TestQualityprofiles_Changelog(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

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

		data, _ := json.Marshal(response)
		w.Write(data)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesChangelogOption{
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	result, resp, err := client.Qualityprofiles.Changelog(opt)
	if err != nil {
		t.Fatalf("Changelog failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.Events) != 1 {
		t.Fatal("expected 1 event")
	}

	if result.Events[0].Action != "ACTIVATED" {
		t.Errorf("expected action 'ACTIVATED', got '%s'", result.Events[0].Action)
	}
}

func TestQualityprofiles_Changelog_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.Qualityprofiles.Changelog(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Language
	_, _, err = client.Qualityprofiles.Changelog(&QualityprofilesChangelogOption{
		QualityProfile: "Sonar way",
	})
	if err == nil {
		t.Error("expected error for missing Language")
	}

	// Test missing QualityProfile
	_, _, err = client.Qualityprofiles.Changelog(&QualityprofilesChangelogOption{
		Language: "java",
	})
	if err == nil {
		t.Error("expected error for missing QualityProfile")
	}

	// Test invalid FilterMode
	_, _, err = client.Qualityprofiles.Changelog(&QualityprofilesChangelogOption{
		Language:       "java",
		QualityProfile: "Sonar way",
		FilterMode:     "INVALID",
	})
	if err == nil {
		t.Error("expected error for invalid FilterMode")
	}
}

func TestQualityprofiles_Compare(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

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

		data, _ := json.Marshal(response)
		w.Write(data)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesCompareOption{
		LeftKey:  "profile1",
		RightKey: "profile2",
	}

	result, resp, err := client.Qualityprofiles.Compare(opt)
	if err != nil {
		t.Fatalf("Compare failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected result")
	}

	if result.Left.Name != "Profile 1" {
		t.Errorf("expected left profile name 'Profile 1', got '%s'", result.Left.Name)
	}

	if len(result.InLeft) != 1 {
		t.Error("expected 1 rule in left")
	}
}

func TestQualityprofiles_Compare_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.Qualityprofiles.Compare(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing LeftKey
	_, _, err = client.Qualityprofiles.Compare(&QualityprofilesCompareOption{
		RightKey: "profile2",
	})
	if err == nil {
		t.Error("expected error for missing LeftKey")
	}

	// Test missing RightKey
	_, _, err = client.Qualityprofiles.Compare(&QualityprofilesCompareOption{
		LeftKey: "profile1",
	})
	if err == nil {
		t.Error("expected error for missing RightKey")
	}
}

func TestQualityprofiles_Copy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		response := &QualityprofilesCopy{
			Key:          "new-profile-key",
			Name:         "My Profile Copy",
			Language:     "java",
			LanguageName: "Java",
			IsDefault:    false,
			IsInherited:  false,
		}

		data, _ := json.Marshal(response)
		w.Write(data)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesCopyOption{
		FromKey: "source-profile-key",
		ToName:  "My Profile Copy",
	}

	result, resp, err := client.Qualityprofiles.Copy(opt)
	if err != nil {
		t.Fatalf("Copy failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || result.Name != "My Profile Copy" {
		t.Error("expected profile name 'My Profile Copy'")
	}
}

func TestQualityprofiles_Copy_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.Qualityprofiles.Copy(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing FromKey
	_, _, err = client.Qualityprofiles.Copy(&QualityprofilesCopyOption{
		ToName: "New Profile",
	})
	if err == nil {
		t.Error("expected error for missing FromKey")
	}

	// Test missing ToName
	_, _, err = client.Qualityprofiles.Copy(&QualityprofilesCopyOption{
		FromKey: "source-key",
	})
	if err == nil {
		t.Error("expected error for missing ToName")
	}

	// Test ToName too long
	_, _, err = client.Qualityprofiles.Copy(&QualityprofilesCopyOption{
		FromKey: "source-key",
		ToName:  strings.Repeat("a", MaxQualityProfileNameLength+1),
	})
	if err == nil {
		t.Error("expected error for ToName exceeding max length")
	}
}

func TestQualityprofiles_Create(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		response := &QualityprofilesCreate{
			Profile: CreatedProfile{
				Key:          "new-profile-key",
				Name:         "My New Profile",
				Language:     "java",
				LanguageName: "Java",
				IsDefault:    false,
			},
		}

		data, _ := json.Marshal(response)
		w.Write(data)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesCreateOption{
		Language: "java",
		Name:     "My New Profile",
	}

	result, resp, err := client.Qualityprofiles.Create(opt)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || result.Profile.Name != "My New Profile" {
		t.Error("expected profile name 'My New Profile'")
	}
}

func TestQualityprofiles_Create_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.Qualityprofiles.Create(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Language
	_, _, err = client.Qualityprofiles.Create(&QualityprofilesCreateOption{
		Name: "My Profile",
	})
	if err == nil {
		t.Error("expected error for missing Language")
	}

	// Test missing Name
	_, _, err = client.Qualityprofiles.Create(&QualityprofilesCreateOption{
		Language: "java",
	})
	if err == nil {
		t.Error("expected error for missing Name")
	}

	// Test Name too long
	_, _, err = client.Qualityprofiles.Create(&QualityprofilesCreateOption{
		Language: "java",
		Name:     strings.Repeat("a", MaxQualityProfileNameLength+1),
	})
	if err == nil {
		t.Error("expected error for Name exceeding max length")
	}
}

func TestQualityprofiles_DeactivateRule(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesDeactivateRuleOption{
		Key:  "AU-TpxcA-iU5OvuD2FL0",
		Rule: "squid:AvoidCycles",
	}

	resp, err := client.Qualityprofiles.DeactivateRule(opt)
	if err != nil {
		t.Fatalf("DeactivateRule failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityprofiles_DeactivateRule_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.Qualityprofiles.DeactivateRule(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Key
	_, err = client.Qualityprofiles.DeactivateRule(&QualityprofilesDeactivateRuleOption{
		Rule: "squid:AvoidCycles",
	})
	if err == nil {
		t.Error("expected error for missing Key")
	}

	// Test missing Rule
	_, err = client.Qualityprofiles.DeactivateRule(&QualityprofilesDeactivateRuleOption{
		Key: "AU-TpxcA-iU5OvuD2FL0",
	})
	if err == nil {
		t.Error("expected error for missing Rule")
	}
}

func TestQualityprofiles_DeactivateRules(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesDeactivateRulesOption{
		TargetKey: "AU-TpxcA-iU5OvuD2FL0",
		Languages: []string{"java"},
	}

	resp, err := client.Qualityprofiles.DeactivateRules(opt)
	if err != nil {
		t.Fatalf("DeactivateRules failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityprofiles_DeactivateRules_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.Qualityprofiles.DeactivateRules(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing TargetKey
	_, err = client.Qualityprofiles.DeactivateRules(&QualityprofilesDeactivateRulesOption{})
	if err == nil {
		t.Error("expected error for missing TargetKey")
	}
}

func TestQualityprofiles_Delete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesDeleteOption{
		Language:       "java",
		QualityProfile: "My Profile",
	}

	resp, err := client.Qualityprofiles.Delete(opt)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityprofiles_Delete_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.Qualityprofiles.Delete(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Language
	_, err = client.Qualityprofiles.Delete(&QualityprofilesDeleteOption{
		QualityProfile: "My Profile",
	})
	if err == nil {
		t.Error("expected error for missing Language")
	}

	// Test missing QualityProfile
	_, err = client.Qualityprofiles.Delete(&QualityprofilesDeleteOption{
		Language: "java",
	})
	if err == nil {
		t.Error("expected error for missing QualityProfile")
	}
}

func TestQualityprofiles_Export(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(200)
		w.Write([]byte(`<?xml version='1.0'?><profile><name>Sonar way</name></profile>`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesExportOption{
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	result, resp, err := client.Qualityprofiles.Export(opt)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || !strings.Contains(*result, "Sonar way") {
		t.Error("expected export to contain 'Sonar way'")
	}
}

func TestQualityprofiles_Export_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.Qualityprofiles.Export(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Language
	_, _, err = client.Qualityprofiles.Export(&QualityprofilesExportOption{})
	if err == nil {
		t.Error("expected error for missing Language")
	}
}

func TestQualityprofiles_Exporters(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		response := &QualityprofilesExporters{
			Exporters: []ProfileExporter{
				{Key: "checkstyle", Name: "Checkstyle", Languages: []string{"java"}},
			},
		}

		data, _ := json.Marshal(response)
		w.Write(data)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Qualityprofiles.Exporters()
	if err != nil {
		t.Fatalf("Exporters failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.Exporters) != 1 {
		t.Fatal("expected 1 exporter")
	}

	if result.Exporters[0].Key != "checkstyle" {
		t.Errorf("expected exporter key 'checkstyle', got '%s'", result.Exporters[0].Key)
	}
}

func TestQualityprofiles_Importers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		response := &QualityprofilesImporters{
			Importers: []ProfileImporter{
				{Key: "pmd", Name: "PMD", Languages: []string{"java"}},
			},
		}

		data, _ := json.Marshal(response)
		w.Write(data)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	result, resp, err := client.Qualityprofiles.Importers()
	if err != nil {
		t.Fatalf("Importers failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.Importers) != 1 {
		t.Fatal("expected 1 importer")
	}

	if result.Importers[0].Key != "pmd" {
		t.Errorf("expected importer key 'pmd', got '%s'", result.Importers[0].Key)
	}
}

func TestQualityprofiles_Inheritance(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

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

		data, _ := json.Marshal(response)
		w.Write(data)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesInheritanceOption{
		Language:       "java",
		QualityProfile: "My Profile",
	}

	result, resp, err := client.Qualityprofiles.Inheritance(opt)
	if err != nil {
		t.Fatalf("Inheritance failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected result")
	}

	if result.Profile.Name != "My Profile" {
		t.Errorf("expected profile name 'My Profile', got '%s'", result.Profile.Name)
	}

	if len(result.Ancestors) != 1 {
		t.Fatal("expected 1 ancestor")
	}

	if result.Ancestors[0].Name != "Sonar way" {
		t.Errorf("expected ancestor name 'Sonar way', got '%s'", result.Ancestors[0].Name)
	}
}

func TestQualityprofiles_Inheritance_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.Qualityprofiles.Inheritance(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Language
	_, _, err = client.Qualityprofiles.Inheritance(&QualityprofilesInheritanceOption{
		QualityProfile: "My Profile",
	})
	if err == nil {
		t.Error("expected error for missing Language")
	}

	// Test missing QualityProfile
	_, _, err = client.Qualityprofiles.Inheritance(&QualityprofilesInheritanceOption{
		Language: "java",
	})
	if err == nil {
		t.Error("expected error for missing QualityProfile")
	}
}

func TestQualityprofiles_Projects(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		response := &QualityprofilesProjects{
			Paging: Paging{PageIndex: 1, PageSize: 25, Total: 2},
			Results: []ProfileProject{
				{Key: "project1", Name: "Project 1", Selected: true},
				{Key: "project2", Name: "Project 2", Selected: false},
			},
		}

		data, _ := json.Marshal(response)
		w.Write(data)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesProjectsOption{
		Key: "profile-key",
	}

	result, resp, err := client.Qualityprofiles.Projects(opt)
	if err != nil {
		t.Fatalf("Projects failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.Results) != 2 {
		t.Fatal("expected 2 projects")
	}

	if result.Results[0].Name != "Project 1" {
		t.Errorf("expected first project name 'Project 1', got '%s'", result.Results[0].Name)
	}
}

func TestQualityprofiles_Projects_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.Qualityprofiles.Projects(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Key
	_, _, err = client.Qualityprofiles.Projects(&QualityprofilesProjectsOption{})
	if err == nil {
		t.Error("expected error for missing Key")
	}

	// Test invalid Selected
	_, _, err = client.Qualityprofiles.Projects(&QualityprofilesProjectsOption{
		Key:      "profile-key",
		Selected: "invalid",
	})
	if err == nil {
		t.Error("expected error for invalid Selected")
	}
}

func TestQualityprofiles_RemoveGroup(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesRemoveGroupOption{
		Group:          "sonar-administrators",
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.RemoveGroup(opt)
	if err != nil {
		t.Fatalf("RemoveGroup failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityprofiles_RemoveGroup_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.Qualityprofiles.RemoveGroup(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Group
	_, err = client.Qualityprofiles.RemoveGroup(&QualityprofilesRemoveGroupOption{
		Language:       "java",
		QualityProfile: "Sonar way",
	})
	if err == nil {
		t.Error("expected error for missing Group")
	}
}

func TestQualityprofiles_RemoveProject(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesRemoveProjectOption{
		Language:       "java",
		Project:        "my_project",
		QualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.RemoveProject(opt)
	if err != nil {
		t.Fatalf("RemoveProject failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityprofiles_RemoveProject_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.Qualityprofiles.RemoveProject(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Language
	_, err = client.Qualityprofiles.RemoveProject(&QualityprofilesRemoveProjectOption{
		Project:        "my_project",
		QualityProfile: "Sonar way",
	})
	if err == nil {
		t.Error("expected error for missing Language")
	}
}

func TestQualityprofiles_RemoveUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesRemoveUserOption{
		Language:       "java",
		Login:          "john.doe",
		QualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.RemoveUser(opt)
	if err != nil {
		t.Fatalf("RemoveUser failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityprofiles_RemoveUser_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.Qualityprofiles.RemoveUser(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Language
	_, err = client.Qualityprofiles.RemoveUser(&QualityprofilesRemoveUserOption{
		Login:          "john.doe",
		QualityProfile: "Sonar way",
	})
	if err == nil {
		t.Error("expected error for missing Language")
	}
}

func TestQualityprofiles_Rename(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesRenameOption{
		Key:  "profile-key",
		Name: "New Profile Name",
	}

	resp, err := client.Qualityprofiles.Rename(opt)
	if err != nil {
		t.Fatalf("Rename failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityprofiles_Rename_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.Qualityprofiles.Rename(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Key
	_, err = client.Qualityprofiles.Rename(&QualityprofilesRenameOption{
		Name: "New Name",
	})
	if err == nil {
		t.Error("expected error for missing Key")
	}

	// Test missing Name
	_, err = client.Qualityprofiles.Rename(&QualityprofilesRenameOption{
		Key: "profile-key",
	})
	if err == nil {
		t.Error("expected error for missing Name")
	}

	// Test Name too long
	_, err = client.Qualityprofiles.Rename(&QualityprofilesRenameOption{
		Key:  "profile-key",
		Name: strings.Repeat("a", MaxQualityProfileNameLength+1),
	})
	if err == nil {
		t.Error("expected error for Name exceeding max length")
	}
}

func TestQualityprofiles_Restore(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesRestoreOption{
		Backup: `<?xml version='1.0'?><profile><name>My Profile</name></profile>`,
	}

	resp, err := client.Qualityprofiles.Restore(opt)
	if err != nil {
		t.Fatalf("Restore failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityprofiles_Restore_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.Qualityprofiles.Restore(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Backup
	_, err = client.Qualityprofiles.Restore(&QualityprofilesRestoreOption{})
	if err == nil {
		t.Error("expected error for missing Backup")
	}
}

func TestQualityprofiles_Search(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

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

		data, _ := json.Marshal(response)
		w.Write(data)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesSearchOption{
		Language: "java",
	}

	result, resp, err := client.Qualityprofiles.Search(opt)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || len(result.Profiles) != 2 {
		t.Fatal("expected 2 profiles")
	}

	if result.Profiles[0].Name != "Sonar way" {
		t.Errorf("expected first profile name 'Sonar way', got '%s'", result.Profiles[0].Name)
	}

	if !result.Actions.Create {
		t.Error("expected create action to be true")
	}
}

func TestQualityprofiles_Search_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.Qualityprofiles.Search(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}
}

func TestQualityprofiles_SearchGroups(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		response := &QualityprofilesSearchGroups{
			Paging: Paging{PageIndex: 1, PageSize: 25, Total: 1},
			Groups: []ProfileGroup{
				{Name: "sonar-administrators", Description: "Admin group", Selected: true},
			},
		}

		data, _ := json.Marshal(response)
		w.Write(data)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesSearchGroupsOption{
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	result, resp, err := client.Qualityprofiles.SearchGroups(opt)
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

func TestQualityprofiles_SearchGroups_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.Qualityprofiles.SearchGroups(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Language
	_, _, err = client.Qualityprofiles.SearchGroups(&QualityprofilesSearchGroupsOption{
		QualityProfile: "Sonar way",
	})
	if err == nil {
		t.Error("expected error for missing Language")
	}

	// Test missing QualityProfile
	_, _, err = client.Qualityprofiles.SearchGroups(&QualityprofilesSearchGroupsOption{
		Language: "java",
	})
	if err == nil {
		t.Error("expected error for missing QualityProfile")
	}

	// Test invalid Selected
	_, _, err = client.Qualityprofiles.SearchGroups(&QualityprofilesSearchGroupsOption{
		Language:       "java",
		QualityProfile: "Sonar way",
		Selected:       "invalid",
	})
	if err == nil {
		t.Error("expected error for invalid Selected")
	}
}

func TestQualityprofiles_SearchUsers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

		response := &QualityprofilesSearchUsers{
			Paging: Paging{PageIndex: 1, PageSize: 25, Total: 1},
			Users: []ProfileUser{
				{Login: "john.doe", Name: "John Doe", Selected: true},
			},
		}

		data, _ := json.Marshal(response)
		w.Write(data)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesSearchUsersOption{
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	result, resp, err := client.Qualityprofiles.SearchUsers(opt)
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
		t.Errorf("expected user login 'john.doe', got '%s'", result.Users[0].Login)
	}
}

func TestQualityprofiles_SearchUsers_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.Qualityprofiles.SearchUsers(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Language
	_, _, err = client.Qualityprofiles.SearchUsers(&QualityprofilesSearchUsersOption{
		QualityProfile: "Sonar way",
	})
	if err == nil {
		t.Error("expected error for missing Language")
	}

	// Test missing QualityProfile
	_, _, err = client.Qualityprofiles.SearchUsers(&QualityprofilesSearchUsersOption{
		Language: "java",
	})
	if err == nil {
		t.Error("expected error for missing QualityProfile")
	}

	// Test invalid Selected
	_, _, err = client.Qualityprofiles.SearchUsers(&QualityprofilesSearchUsersOption{
		Language:       "java",
		QualityProfile: "Sonar way",
		Selected:       "invalid",
	})
	if err == nil {
		t.Error("expected error for invalid Selected")
	}
}

func TestQualityprofiles_SetDefault(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		w.WriteHeader(204)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesSetDefaultOption{
		Language:       "java",
		QualityProfile: "Sonar way",
	}

	resp, err := client.Qualityprofiles.SetDefault(opt)
	if err != nil {
		t.Fatalf("SetDefault failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestQualityprofiles_SetDefault_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, err = client.Qualityprofiles.SetDefault(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Language
	_, err = client.Qualityprofiles.SetDefault(&QualityprofilesSetDefaultOption{
		QualityProfile: "Sonar way",
	})
	if err == nil {
		t.Error("expected error for missing Language")
	}

	// Test missing QualityProfile
	_, err = client.Qualityprofiles.SetDefault(&QualityprofilesSetDefaultOption{
		Language: "java",
	})
	if err == nil {
		t.Error("expected error for missing QualityProfile")
	}
}

func TestQualityprofiles_Show(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)

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

		data, _ := json.Marshal(response)
		w.Write(data)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &QualityprofilesShowOption{
		Key: "sonar-way-java",
	}

	result, resp, err := client.Qualityprofiles.Show(opt)
	if err != nil {
		t.Fatalf("Show failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil || result.Profile.Name != "Sonar way" {
		t.Error("expected profile name 'Sonar way'")
	}

	if result.Profile.ActiveRuleCount != 200 {
		t.Errorf("expected 200 active rules, got %d", result.Profile.ActiveRuleCount)
	}
}

func TestQualityprofiles_Show_ValidationError(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test nil option
	_, _, err = client.Qualityprofiles.Show(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Test missing Key
	_, _, err = client.Qualityprofiles.Show(&QualityprofilesShowOption{})
	if err == nil {
		t.Error("expected error for missing Key")
	}
}

func TestQualityprofiles_ConvertActivateRuleOptForURL(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

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
	if urlOpt.Key != opt.Key {
		t.Errorf("expected Key %s, got %s", opt.Key, urlOpt.Key)
	}

	if urlOpt.Rule != opt.Rule {
		t.Errorf("expected Rule %s, got %s", opt.Rule, urlOpt.Rule)
	}

	if urlOpt.PrioritizedRule != opt.PrioritizedRule {
		t.Errorf("expected PrioritizedRule %v, got %v", opt.PrioritizedRule, urlOpt.PrioritizedRule)
	}

	// Verify map conversions
	if urlOpt.Impacts == "" {
		t.Error("expected Impacts to be converted to string")
	}

	// Impacts string should contain both entries
	if !strings.Contains(urlOpt.Impacts, "MAINTAINABILITY=HIGH") {
		t.Error("expected Impacts to contain MAINTAINABILITY=HIGH")
	}

	if !strings.Contains(urlOpt.Impacts, "SECURITY=MEDIUM") {
		t.Error("expected Impacts to contain SECURITY=MEDIUM")
	}

	if urlOpt.Params == "" {
		t.Error("expected Params to be converted to string")
	}

	// Params string should contain both entries
	if !strings.Contains(urlOpt.Params, "max=10") {
		t.Error("expected Params to contain max=10")
	}

	if !strings.Contains(urlOpt.Params, "threshold=5") {
		t.Error("expected Params to contain threshold=5")
	}
}
