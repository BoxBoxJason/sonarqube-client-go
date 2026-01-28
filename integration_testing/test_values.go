package sonargo_test

import (
	"fmt"

	sonar "github.com/boxboxjason/sonarqube-client-go/sonar"
)

var (
	testQualityGate string = ""
	testWebhook     string = ""
)

func SetupTestResources(c *sonar.Client) {
	// Create Test Project
	_, _, err := c.Projects.Create(&sonar.ProjectsCreateOption{
		Name:       "Test Project",
		Project:    "test-project",
		Visibility: "public",
	})
	if err != nil {
		fmt.Printf("Setup: Project creation failed (might exist): %v\n", err)
	} else {
		fmt.Println("Setup: Created test-project")
	}

	// Create Test User
	_, _, err = c.Users.Create(&sonar.UsersCreateOption{
		Login:    "test-user",
		Name:     "Test User",
		Password: "test-password",
		Local:    true,
	})
	if err != nil {
		fmt.Printf("Setup: User creation failed (might exist): %v\n", err)
	} else {
		fmt.Println("Setup: Created test-user")
	}

	// Create Quality Profile
	_, _, err = c.Qualityprofiles.Create(&sonar.QualityprofilesCreateOption{
		Language: "java",
		Name:     "test-profile",
	})
	if err != nil {
		fmt.Printf("Setup: Quality Profile creation failed (might exist): %v\n", err)
	} else {
		fmt.Println("Setup: Created test-profile")
	}

	// Create User Group
	_, _, err = c.UserGroups.Create(&sonar.UserGroupsCreateOption{
		Name: "test-group",
	})
	if err != nil {
		fmt.Printf("Setup: User Group creation failed (might exist): %v\n", err)
	} else {
		fmt.Println("Setup: Created test-group")
	}

	// Create Quality Gate
	qg, _, err := c.Qualitygates.Create(&sonar.QualitygatesCreateOption{
		Name: "test-quality-gate",
	})
	if err != nil {
		fmt.Printf("Setup: Quality Gate creation failed (might exist): %v\n", err)
		// Try to use default quality gate or hardcode "1"
		testQualityGate = "test-quality-gate"
	} else {
		testQualityGate = qg.ID
		fmt.Printf("Setup: Created quality gate: %s\n", testQualityGate)
	}

	// Create Webhook
	wh, _, err := c.Webhooks.Create(&sonar.WebhooksCreateOption{
		Name:    "test-webhook",
		URL:     "https://example.com/webhook",
		Project: "test-project",
	})
	if err != nil {
		fmt.Printf("Setup: Webhook creation failed (might exist): %v\n", err)
	} else {
		testWebhook = wh.Webhook.Key
		fmt.Printf("Setup: Created webhook: %s\n", testWebhook)
	}

	// Create Permission Template
	_, _, err = c.Permissions.CreateTemplate(&sonar.PermissionsCreateTemplateOption{
		Name: "test-permission-template",
	})
	if err != nil {
		fmt.Printf("Setup: Permission Template creation failed (might exist): %v\n", err)
	} else {
		fmt.Println("Setup: Created test-permission-template")
	}

	// Add user to group
	_, err = c.UserGroups.AddUser(&sonar.UserGroupsAddUserOption{
		Name:  "test-group",
		Login: "test-user",
	})
	if err != nil {
		fmt.Printf("Setup: Add user to group failed: %v\n", err)
	} else {
		fmt.Println("Setup: Added test-user to test-group")
	}

	// Add project to favorites
	_, err = c.Favorites.Add(&sonar.FavoritesAddOption{
		Component: "test-project",
	})
	if err != nil {
		fmt.Printf("Setup: Add favorite failed: %v\n", err)
	} else {
		fmt.Println("Setup: Added test-project to favorites")
	}
}

// GetTestValue returns test values for specific service/action/param combinations.
func GetTestValue(service, action, param string) string {
	// Quality Gate values
	if param == "gateName" || param == "gateId" {
		if testQualityGate != "" {
			return testQualityGate
		}

		return "test-quality-gate"
	}

	// Webhook values
	if param == "webhook" {
		if testWebhook != "" {
			return testWebhook
		}

		return "test-webhook"
	}

	// Common values
	if param == "project" || param == "projectKey" || param == "component" {
		return "test-project"
	}

	if param == "name" {
		return "Test Project"
	}

	if param == "login" {
		return "test-user"
	}

	if param == "password" {
		return "test-password"
	}

	if param == "key" {
		return "test-project"
	}

	if param == "profileName" || param == "qualityProfile" {
		return "test-profile"
	}

	if param == "language" {
		return "java"
	}

	if param == "templateName" {
		return "test-permission-template"
	}

	if param == "groupName" || param == "group" {
		return "test-group"
	}

	if param == "url" {
		return "https://example.com/webhook"
	}

	if param == "rule" || param == "ruleKey" {
		return "java:S100"
	}

	if param == "metric" || param == "metricKey" {
		return "coverage"
	}

	if param == "op" {
		return "LT"
	}

	if param == "error" {
		return "80"
	}

	if param == "value" || param == "values" {
		return "80"
	}

	return ""
}
