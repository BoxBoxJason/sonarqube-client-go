package integration_testing_test

import (
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Permissions Service", Ordered, func() {
	var (
		client  *sonargo.Client
		cleanup *helpers.CleanupManager
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
		cleanup = helpers.NewCleanupManager(client)
	})

	AfterAll(func() {
		// Cleanup all test resources
		errors := cleanup.Cleanup()
		for _, err := range errors {
			GinkgoWriter.Printf("Cleanup error: %v\n", err)
		}
	})

	// =========================================================================
	// User Permissions (Project-level)
	// =========================================================================
	Describe("AddUser", func() {
		var testProjectKey string
		var testUserLogin string

		BeforeEach(func() {
			testProjectKey = helpers.UniqueResourceName("proj-perm")
			testUserLogin = helpers.UniqueResourceName("user-perm")

			// Create private project (codeviewer/user permissions only work on private projects)
			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "Permission Test Project",
				Project:    testProjectKey,
				Visibility: "private",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", testProjectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: testProjectKey,
				})
				return err
			})

			// Create user
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err = client.Users.Create(&sonargo.UsersCreateOption{
				Login:    testUserLogin,
				Name:     "Permission Test User",
				Password: "SecurePassword123!",
				Local:    true,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("user", testUserLogin, func() error {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, _, err := client.Users.Deactivate(&sonargo.UsersDeactivateOption{
					Login:     testUserLogin,
					Anonymize: true,
				})
				return err
			})
		})

		It("should add user permission to project", func() {
			resp, err := client.Permissions.AddUser(&sonargo.PermissionsAddUserOption{
				Login:      testUserLogin,
				Permission: "codeviewer",
				ProjectKey: testProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify user has permission
			result, _, err := client.Permissions.Users(&sonargo.PermissionsUsersOption{
				ProjectKey: testProjectKey,
				Permission: "codeviewer",
			})
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, u := range result.Users {
				if u.Login == testUserLogin {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue(), "Expected user %s to have codeviewer permission on project %s", testUserLogin, testProjectKey)
		})

		It("should add admin permission to project", func() {
			resp, err := client.Permissions.AddUser(&sonargo.PermissionsAddUserOption{
				Login:      testUserLogin,
				Permission: "admin",
				ProjectKey: testProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		It("should add issueadmin permission to project", func() {
			resp, err := client.Permissions.AddUser(&sonargo.PermissionsAddUserOption{
				Login:      testUserLogin,
				Permission: "issueadmin",
				ProjectKey: testProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Permissions.AddUser(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing login", func() {
				resp, err := client.Permissions.AddUser(&sonargo.PermissionsAddUserOption{
					Permission: "codeviewer",
					ProjectKey: testProjectKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing permission", func() {
				resp, err := client.Permissions.AddUser(&sonargo.PermissionsAddUserOption{
					Login:      testUserLogin,
					ProjectKey: testProjectKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid permission", func() {
				resp, err := client.Permissions.AddUser(&sonargo.PermissionsAddUserOption{
					Login:      testUserLogin,
					Permission: "invalid_permission",
					ProjectKey: testProjectKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("RemoveUser", func() {
		var testProjectKey string
		var testUserLogin string

		BeforeEach(func() {
			testProjectKey = helpers.UniqueResourceName("proj-rmuser")
			testUserLogin = helpers.UniqueResourceName("user-rmperm")

			// Create private project (codeviewer permission only works on private projects)
			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "Remove User Permission Test",
				Project:    testProjectKey,
				Visibility: "private",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", testProjectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: testProjectKey,
				})
				return err
			})

			// Create user
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err = client.Users.Create(&sonargo.UsersCreateOption{
				Login:    testUserLogin,
				Name:     "Remove Permission Test User",
				Password: "SecurePassword123!",
				Local:    true,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("user", testUserLogin, func() error {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, _, err := client.Users.Deactivate(&sonargo.UsersDeactivateOption{
					Login:     testUserLogin,
					Anonymize: true,
				})
				return err
			})

			// Add permission first
			_, err = client.Permissions.AddUser(&sonargo.PermissionsAddUserOption{
				Login:      testUserLogin,
				Permission: "codeviewer",
				ProjectKey: testProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should remove user permission from project", func() {
			resp, err := client.Permissions.RemoveUser(&sonargo.PermissionsRemoveUserOption{
				Login:      testUserLogin,
				Permission: "codeviewer",
				ProjectKey: testProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify user no longer has permission
			result, _, err := client.Permissions.Users(&sonargo.PermissionsUsersOption{
				ProjectKey: testProjectKey,
				Permission: "codeviewer",
			})
			Expect(err).NotTo(HaveOccurred())
			for _, u := range result.Users {
				Expect(u.Login).NotTo(Equal(testUserLogin))
			}
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Permissions.RemoveUser(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing login", func() {
				resp, err := client.Permissions.RemoveUser(&sonargo.PermissionsRemoveUserOption{
					Permission: "codeviewer",
					ProjectKey: testProjectKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing permission", func() {
				resp, err := client.Permissions.RemoveUser(&sonargo.PermissionsRemoveUserOption{
					Login:      testUserLogin,
					ProjectKey: testProjectKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Group Permissions (Project-level)
	// =========================================================================
	Describe("AddGroup", func() {
		var testProjectKey string
		var testGroupName string

		BeforeEach(func() {
			testProjectKey = helpers.UniqueResourceName("proj-grpperm")
			testGroupName = helpers.UniqueResourceName("grp-perm")

			// Create private project (codeviewer permission only works on private projects)
			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "Group Permission Test Project",
				Project:    testProjectKey,
				Visibility: "private",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", testProjectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: testProjectKey,
				})
				return err
			})

			// Create group
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err = client.UserGroups.Create(&sonargo.UserGroupsCreateOption{
				Name:        testGroupName,
				Description: "Permission Test Group",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("group", testGroupName, func() error {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, err := client.UserGroups.Delete(&sonargo.UserGroupsDeleteOption{
					Name: testGroupName,
				})
				return err
			})
		})

		It("should add group permission to project", func() {
			resp, err := client.Permissions.AddGroup(&sonargo.PermissionsAddGroupOption{
				GroupName:  testGroupName,
				Permission: "codeviewer",
				ProjectKey: testProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify group has permission
			result, _, err := client.Permissions.Groups(&sonargo.PermissionsGroupsOption{
				ProjectKey: testProjectKey,
				Permission: "codeviewer",
			})
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, g := range result.Groups {
				if g.Name == testGroupName {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		It("should add admin permission to group", func() {
			resp, err := client.Permissions.AddGroup(&sonargo.PermissionsAddGroupOption{
				GroupName:  testGroupName,
				Permission: "admin",
				ProjectKey: testProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify group has admin permission
			result, _, err := client.Permissions.Groups(&sonargo.PermissionsGroupsOption{
				ProjectKey: testProjectKey,
				Permission: "admin",
			})
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, g := range result.Groups {
				if g.Name == testGroupName {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Permissions.AddGroup(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing group name", func() {
				resp, err := client.Permissions.AddGroup(&sonargo.PermissionsAddGroupOption{
					Permission: "codeviewer",
					ProjectKey: testProjectKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing permission", func() {
				resp, err := client.Permissions.AddGroup(&sonargo.PermissionsAddGroupOption{
					GroupName:  testGroupName,
					ProjectKey: testProjectKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid permission", func() {
				resp, err := client.Permissions.AddGroup(&sonargo.PermissionsAddGroupOption{
					GroupName:  testGroupName,
					Permission: "invalid_permission",
					ProjectKey: testProjectKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("RemoveGroup", func() {
		var testProjectKey string
		var testGroupName string

		BeforeEach(func() {
			testProjectKey = helpers.UniqueResourceName("proj-rmgrp")
			testGroupName = helpers.UniqueResourceName("grp-rmperm")

			// Create private project (codeviewer permission only works on private projects)
			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "Remove Group Permission Test",
				Project:    testProjectKey,
				Visibility: "private",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", testProjectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: testProjectKey,
				})
				return err
			})

			// Create group
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err = client.UserGroups.Create(&sonargo.UserGroupsCreateOption{
				Name:        testGroupName,
				Description: "Remove Permission Test Group",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("group", testGroupName, func() error {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, err := client.UserGroups.Delete(&sonargo.UserGroupsDeleteOption{
					Name: testGroupName,
				})
				return err
			})

			// Add permission first
			_, err = client.Permissions.AddGroup(&sonargo.PermissionsAddGroupOption{
				GroupName:  testGroupName,
				Permission: "codeviewer",
				ProjectKey: testProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should remove group permission from project", func() {
			resp, err := client.Permissions.RemoveGroup(&sonargo.PermissionsRemoveGroupOption{
				GroupName:  testGroupName,
				Permission: "codeviewer",
				ProjectKey: testProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify group no longer has permission
			result, _, err := client.Permissions.Groups(&sonargo.PermissionsGroupsOption{
				ProjectKey: testProjectKey,
				Permission: "codeviewer",
			})
			Expect(err).NotTo(HaveOccurred())
			for _, g := range result.Groups {
				Expect(g.Name).NotTo(Equal(testGroupName))
			}
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Permissions.RemoveGroup(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing group name", func() {
				resp, err := client.Permissions.RemoveGroup(&sonargo.PermissionsRemoveGroupOption{
					Permission: "codeviewer",
					ProjectKey: testProjectKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing permission", func() {
				resp, err := client.Permissions.RemoveGroup(&sonargo.PermissionsRemoveGroupOption{
					GroupName:  testGroupName,
					ProjectKey: testProjectKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Users and Groups listing
	// =========================================================================
	Describe("Users", func() {
		It("should list users globally", func() {
			result, resp, err := client.Permissions.Users(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			// Should return a list of users
			Expect(result.Users).NotTo(BeEmpty())
		})

		It("should list users with specific permission", func() {
			result, resp, err := client.Permissions.Users(&sonargo.PermissionsUsersOption{
				Permission: "admin",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		Context("with project scope", func() {
			var testProjectKey string

			BeforeEach(func() {
				testProjectKey = helpers.UniqueResourceName("proj-listuser")

				_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
					Name:    "List Users Test",
					Project: testProjectKey,
				})
				Expect(err).NotTo(HaveOccurred())

				cleanup.RegisterCleanup("project", testProjectKey, func() error {
					_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
						Project: testProjectKey,
					})
					return err
				})
			})

			It("should list users with project permissions", func() {
				result, resp, err := client.Permissions.Users(&sonargo.PermissionsUsersOption{
					ProjectKey: testProjectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("parameter validation", func() {
			It("should fail with query too short", func() {
				_, resp, err := client.Permissions.Users(&sonargo.PermissionsUsersOption{
					Query: "ab", // min 3 chars
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid permission", func() {
				_, resp, err := client.Permissions.Users(&sonargo.PermissionsUsersOption{
					Permission: "invalid_permission",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("Groups", func() {
		It("should list groups with global permissions", func() {
			result, resp, err := client.Permissions.Groups(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			// sonar-administrators should have global permissions
			found := false
			for _, g := range result.Groups {
				if g.Name == "sonar-administrators" {
					found = true
					Expect(g.Permissions).NotTo(BeEmpty())
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		It("should list groups with specific permission", func() {
			result, resp, err := client.Permissions.Groups(&sonargo.PermissionsGroupsOption{
				Permission: "admin",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		Context("with project scope", func() {
			var testProjectKey string

			BeforeEach(func() {
				testProjectKey = helpers.UniqueResourceName("proj-listgrp")

				_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
					Name:    "List Groups Test",
					Project: testProjectKey,
				})
				Expect(err).NotTo(HaveOccurred())

				cleanup.RegisterCleanup("project", testProjectKey, func() error {
					_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
						Project: testProjectKey,
					})
					return err
				})
			})

			It("should list groups with project permissions", func() {
				result, resp, err := client.Permissions.Groups(&sonargo.PermissionsGroupsOption{
					ProjectKey: testProjectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("parameter validation", func() {
			It("should fail with query too short", func() {
				_, resp, err := client.Permissions.Groups(&sonargo.PermissionsGroupsOption{
					Query: "ab", // min 3 chars
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid permission", func() {
				_, resp, err := client.Permissions.Groups(&sonargo.PermissionsGroupsOption{
					Permission: "invalid_permission",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Permission Templates
	// =========================================================================
	Describe("CreateTemplate", func() {
		It("should create a permission template", func() {
			templateName := helpers.UniqueResourceName("tpl")

			result, resp, err := client.Permissions.CreateTemplate(&sonargo.PermissionsCreateTemplateOption{
				Name:        templateName,
				Description: "E2E Test Template",
			})

			Expect(err).NotTo(HaveOccurred())

			// Register cleanup
			cleanup.RegisterCleanup("template", templateName, func() error {
				_, err := client.Permissions.DeleteTemplate(&sonargo.PermissionsDeleteTemplateOption{
					TemplateName: templateName,
				})
				return err
			})

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.PermissionTemplate.Name).To(Equal(templateName))
		})

		It("should create a template with project key pattern", func() {
			templateName := helpers.UniqueResourceName("tpl-pattern")

			result, resp, err := client.Permissions.CreateTemplate(&sonargo.PermissionsCreateTemplateOption{
				Name:              templateName,
				Description:       "Template with pattern",
				ProjectKeyPattern: "e2e-.*",
			})

			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("template", templateName, func() error {
				_, err := client.Permissions.DeleteTemplate(&sonargo.PermissionsDeleteTemplateOption{
					TemplateName: templateName,
				})
				return err
			})

			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.PermissionTemplate.ProjectKeyPattern).To(Equal("e2e-.*"))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.Permissions.CreateTemplate(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing name", func() {
				_, resp, err := client.Permissions.CreateTemplate(&sonargo.PermissionsCreateTemplateOption{
					Description: "Test",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("SearchTemplates", func() {
		var testTemplateName string

		BeforeEach(func() {
			testTemplateName = helpers.UniqueResourceName("tpl-search")

			_, _, err := client.Permissions.CreateTemplate(&sonargo.PermissionsCreateTemplateOption{
				Name:        testTemplateName,
				Description: "Search Test Template",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("template", testTemplateName, func() error {
				_, err := client.Permissions.DeleteTemplate(&sonargo.PermissionsDeleteTemplateOption{
					TemplateName: testTemplateName,
				})
				return err
			})
		})

		It("should search all templates", func() {
			result, resp, err := client.Permissions.SearchTemplates(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.PermissionTemplates).NotTo(BeEmpty())
		})

		It("should search templates by query", func() {
			result, resp, err := client.Permissions.SearchTemplates(&sonargo.PermissionsSearchTemplatesOption{
				Query: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			found := false
			for _, t := range result.PermissionTemplates {
				if t.Name == testTemplateName {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		It("should include default templates info", func() {
			result, resp, err := client.Permissions.SearchTemplates(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			// There should be at least one default template
			Expect(result.DefaultTemplates).NotTo(BeEmpty())
		})
	})

	Describe("DeleteTemplate", func() {
		It("should delete a template", func() {
			templateName := helpers.UniqueResourceName("tpl-del")

			_, _, err := client.Permissions.CreateTemplate(&sonargo.PermissionsCreateTemplateOption{
				Name: templateName,
			})
			Expect(err).NotTo(HaveOccurred())

			resp, err := client.Permissions.DeleteTemplate(&sonargo.PermissionsDeleteTemplateOption{
				TemplateName: templateName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify template is deleted
			result, _, err := client.Permissions.SearchTemplates(&sonargo.PermissionsSearchTemplatesOption{
				Query: templateName,
			})
			Expect(err).NotTo(HaveOccurred())
			for _, t := range result.PermissionTemplates {
				Expect(t.Name).NotTo(Equal(templateName))
			}
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Permissions.DeleteTemplate(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing template identifier", func() {
				resp, err := client.Permissions.DeleteTemplate(&sonargo.PermissionsDeleteTemplateOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("UpdateTemplate", func() {
		var testTemplateName string
		var templateID string

		BeforeEach(func() {
			testTemplateName = helpers.UniqueResourceName("tpl-update")

			result, _, err := client.Permissions.CreateTemplate(&sonargo.PermissionsCreateTemplateOption{
				Name:        testTemplateName,
				Description: "Original description",
			})
			Expect(err).NotTo(HaveOccurred())

			// Get template ID from search
			searchResult, _, err := client.Permissions.SearchTemplates(&sonargo.PermissionsSearchTemplatesOption{
				Query: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())
			for _, t := range searchResult.PermissionTemplates {
				if t.Name == testTemplateName {
					templateID = t.ID
					break
				}
			}
			Expect(templateID).NotTo(BeEmpty(), "Template ID not found for: %s, result: %+v", testTemplateName, result)

			cleanup.RegisterCleanup("template", testTemplateName, func() error {
				_, err := client.Permissions.DeleteTemplate(&sonargo.PermissionsDeleteTemplateOption{
					TemplateName: testTemplateName,
				})
				return err
			})
		})

		It("should update template description", func() {
			result, resp, err := client.Permissions.UpdateTemplate(&sonargo.PermissionsUpdateTemplateOption{
				ID:          templateID,
				Description: "Updated description",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.PermissionTemplate.Description).To(Equal("Updated description"))
		})

		It("should update template name", func() {
			newName := helpers.UniqueResourceName("tpl-renamed")

			result, resp, err := client.Permissions.UpdateTemplate(&sonargo.PermissionsUpdateTemplateOption{
				ID:   templateID,
				Name: newName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.PermissionTemplate.Name).To(Equal(newName))

			// Update cleanup to use new name
			cleanup.RegisterCleanup("template", newName, func() error {
				_, err := client.Permissions.DeleteTemplate(&sonargo.PermissionsDeleteTemplateOption{
					TemplateName: newName,
				})
				return err
			})
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.Permissions.UpdateTemplate(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing ID", func() {
				_, resp, err := client.Permissions.UpdateTemplate(&sonargo.PermissionsUpdateTemplateOption{
					Description: "test",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Template Users and Groups
	// =========================================================================
	Describe("AddUserToTemplate", func() {
		var testTemplateName string
		var testUserLogin string

		BeforeEach(func() {
			testTemplateName = helpers.UniqueResourceName("tpl-adduser")
			testUserLogin = helpers.UniqueResourceName("user-tpl")

			// Create template
			_, _, err := client.Permissions.CreateTemplate(&sonargo.PermissionsCreateTemplateOption{
				Name: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("template", testTemplateName, func() error {
				_, err := client.Permissions.DeleteTemplate(&sonargo.PermissionsDeleteTemplateOption{
					TemplateName: testTemplateName,
				})
				return err
			})

			// Create user
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err = client.Users.Create(&sonargo.UsersCreateOption{
				Login:    testUserLogin,
				Name:     "Template User Test",
				Password: "SecurePassword123!",
				Local:    true,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("user", testUserLogin, func() error {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, _, err := client.Users.Deactivate(&sonargo.UsersDeactivateOption{
					Login:     testUserLogin,
					Anonymize: true,
				})
				return err
			})
		})

		It("should add user to template", func() {
			resp, err := client.Permissions.AddUserToTemplate(&sonargo.PermissionsAddUserToTemplateOption{
				Login:        testUserLogin,
				Permission:   "codeviewer",
				TemplateName: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify user is in template
			result, _, err := client.Permissions.TemplateUsers(&sonargo.PermissionsTemplateUsersOption{
				TemplateName: testTemplateName,
				Permission:   "codeviewer",
			})
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, u := range result.Users {
				if u.Login == testUserLogin {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Permissions.AddUserToTemplate(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing login", func() {
				resp, err := client.Permissions.AddUserToTemplate(&sonargo.PermissionsAddUserToTemplateOption{
					Permission:   "codeviewer",
					TemplateName: testTemplateName,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing permission", func() {
				resp, err := client.Permissions.AddUserToTemplate(&sonargo.PermissionsAddUserToTemplateOption{
					Login:        testUserLogin,
					TemplateName: testTemplateName,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing template identifier", func() {
				resp, err := client.Permissions.AddUserToTemplate(&sonargo.PermissionsAddUserToTemplateOption{
					Login:      testUserLogin,
					Permission: "codeviewer",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("RemoveUserFromTemplate", func() {
		var testTemplateName string
		var testUserLogin string

		BeforeEach(func() {
			testTemplateName = helpers.UniqueResourceName("tpl-rmuser")
			testUserLogin = helpers.UniqueResourceName("user-tplrm")

			// Create template
			_, _, err := client.Permissions.CreateTemplate(&sonargo.PermissionsCreateTemplateOption{
				Name: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("template", testTemplateName, func() error {
				_, err := client.Permissions.DeleteTemplate(&sonargo.PermissionsDeleteTemplateOption{
					TemplateName: testTemplateName,
				})
				return err
			})

			// Create user
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err = client.Users.Create(&sonargo.UsersCreateOption{
				Login:    testUserLogin,
				Name:     "Remove Template User Test",
				Password: "SecurePassword123!",
				Local:    true,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("user", testUserLogin, func() error {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, _, err := client.Users.Deactivate(&sonargo.UsersDeactivateOption{
					Login:     testUserLogin,
					Anonymize: true,
				})
				return err
			})

			// Add user to template
			_, err = client.Permissions.AddUserToTemplate(&sonargo.PermissionsAddUserToTemplateOption{
				Login:        testUserLogin,
				Permission:   "codeviewer",
				TemplateName: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should remove user from template", func() {
			resp, err := client.Permissions.RemoveUserFromTemplate(&sonargo.PermissionsRemoveUserFromTemplateOption{
				Login:        testUserLogin,
				Permission:   "codeviewer",
				TemplateName: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify user is removed from template
			result, _, err := client.Permissions.TemplateUsers(&sonargo.PermissionsTemplateUsersOption{
				TemplateName: testTemplateName,
				Permission:   "codeviewer",
			})
			Expect(err).NotTo(HaveOccurred())
			for _, u := range result.Users {
				Expect(u.Login).NotTo(Equal(testUserLogin))
			}
		})
	})

	Describe("AddGroupToTemplate", func() {
		var testTemplateName string
		var testGroupName string

		BeforeEach(func() {
			testTemplateName = helpers.UniqueResourceName("tpl-addgrp")
			testGroupName = helpers.UniqueResourceName("grp-tpl")

			// Create template
			_, _, err := client.Permissions.CreateTemplate(&sonargo.PermissionsCreateTemplateOption{
				Name: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("template", testTemplateName, func() error {
				_, err := client.Permissions.DeleteTemplate(&sonargo.PermissionsDeleteTemplateOption{
					TemplateName: testTemplateName,
				})
				return err
			})

			// Create group
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err = client.UserGroups.Create(&sonargo.UserGroupsCreateOption{
				Name:        testGroupName,
				Description: "Template Group Test",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("group", testGroupName, func() error {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, err := client.UserGroups.Delete(&sonargo.UserGroupsDeleteOption{
					Name: testGroupName,
				})
				return err
			})
		})

		It("should add group to template", func() {
			resp, err := client.Permissions.AddGroupToTemplate(&sonargo.PermissionsAddGroupToTemplateOption{
				GroupName:    testGroupName,
				Permission:   "codeviewer",
				TemplateName: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify group is in template
			result, _, err := client.Permissions.TemplateGroups(&sonargo.PermissionsTemplateGroupsOption{
				TemplateName: testTemplateName,
				Permission:   "codeviewer",
			})
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, g := range result.Groups {
				if g.Name == testGroupName {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Permissions.AddGroupToTemplate(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing group name", func() {
				resp, err := client.Permissions.AddGroupToTemplate(&sonargo.PermissionsAddGroupToTemplateOption{
					Permission:   "codeviewer",
					TemplateName: testTemplateName,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing permission", func() {
				resp, err := client.Permissions.AddGroupToTemplate(&sonargo.PermissionsAddGroupToTemplateOption{
					GroupName:    testGroupName,
					TemplateName: testTemplateName,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing template identifier", func() {
				resp, err := client.Permissions.AddGroupToTemplate(&sonargo.PermissionsAddGroupToTemplateOption{
					GroupName:  testGroupName,
					Permission: "codeviewer",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("RemoveGroupFromTemplate", func() {
		var testTemplateName string
		var testGroupName string

		BeforeEach(func() {
			testTemplateName = helpers.UniqueResourceName("tpl-rmgrp")
			testGroupName = helpers.UniqueResourceName("grp-tplrm")

			// Create template
			_, _, err := client.Permissions.CreateTemplate(&sonargo.PermissionsCreateTemplateOption{
				Name: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("template", testTemplateName, func() error {
				_, err := client.Permissions.DeleteTemplate(&sonargo.PermissionsDeleteTemplateOption{
					TemplateName: testTemplateName,
				})
				return err
			})

			// Create group
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err = client.UserGroups.Create(&sonargo.UserGroupsCreateOption{
				Name:        testGroupName,
				Description: "Remove Template Group Test",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("group", testGroupName, func() error {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, err := client.UserGroups.Delete(&sonargo.UserGroupsDeleteOption{
					Name: testGroupName,
				})
				return err
			})

			// Add group to template
			_, err = client.Permissions.AddGroupToTemplate(&sonargo.PermissionsAddGroupToTemplateOption{
				GroupName:    testGroupName,
				Permission:   "codeviewer",
				TemplateName: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should remove group from template", func() {
			resp, err := client.Permissions.RemoveGroupFromTemplate(&sonargo.PermissionsRemoveGroupFromTemplateOption{
				GroupName:    testGroupName,
				Permission:   "codeviewer",
				TemplateName: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify group is removed from template
			result, _, err := client.Permissions.TemplateGroups(&sonargo.PermissionsTemplateGroupsOption{
				TemplateName: testTemplateName,
				Permission:   "codeviewer",
			})
			Expect(err).NotTo(HaveOccurred())
			for _, g := range result.Groups {
				Expect(g.Name).NotTo(Equal(testGroupName))
			}
		})
	})

	Describe("TemplateUsers", func() {
		var testTemplateName string

		BeforeEach(func() {
			testTemplateName = helpers.UniqueResourceName("tpl-usrlist")

			_, _, err := client.Permissions.CreateTemplate(&sonargo.PermissionsCreateTemplateOption{
				Name: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("template", testTemplateName, func() error {
				_, err := client.Permissions.DeleteTemplate(&sonargo.PermissionsDeleteTemplateOption{
					TemplateName: testTemplateName,
				})
				return err
			})
		})

		It("should list template users", func() {
			result, resp, err := client.Permissions.TemplateUsers(&sonargo.PermissionsTemplateUsersOption{
				TemplateName: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.Permissions.TemplateUsers(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing template identifier", func() {
				_, resp, err := client.Permissions.TemplateUsers(&sonargo.PermissionsTemplateUsersOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("TemplateGroups", func() {
		var testTemplateName string

		BeforeEach(func() {
			testTemplateName = helpers.UniqueResourceName("tpl-grplist")

			_, _, err := client.Permissions.CreateTemplate(&sonargo.PermissionsCreateTemplateOption{
				Name: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("template", testTemplateName, func() error {
				_, err := client.Permissions.DeleteTemplate(&sonargo.PermissionsDeleteTemplateOption{
					TemplateName: testTemplateName,
				})
				return err
			})
		})

		It("should list template groups", func() {
			result, resp, err := client.Permissions.TemplateGroups(&sonargo.PermissionsTemplateGroupsOption{
				TemplateName: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.Permissions.TemplateGroups(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing template identifier", func() {
				_, resp, err := client.Permissions.TemplateGroups(&sonargo.PermissionsTemplateGroupsOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Project Creator Template Permissions
	// =========================================================================
	Describe("AddProjectCreatorToTemplate", func() {
		var testTemplateName string

		BeforeEach(func() {
			testTemplateName = helpers.UniqueResourceName("tpl-creator")

			_, _, err := client.Permissions.CreateTemplate(&sonargo.PermissionsCreateTemplateOption{
				Name: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("template", testTemplateName, func() error {
				_, err := client.Permissions.DeleteTemplate(&sonargo.PermissionsDeleteTemplateOption{
					TemplateName: testTemplateName,
				})
				return err
			})
		})

		It("should add project creator to template", func() {
			resp, err := client.Permissions.AddProjectCreatorToTemplate(&sonargo.PermissionsAddProjectCreatorToTemplateOption{
				Permission:   "admin",
				TemplateName: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify by searching templates
			result, _, err := client.Permissions.SearchTemplates(&sonargo.PermissionsSearchTemplatesOption{
				Query: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, t := range result.PermissionTemplates {
				if t.Name == testTemplateName {
					for _, p := range t.Permissions {
						if p.Key == "admin" && p.WithProjectCreator {
							found = true
							break
						}
					}
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Permissions.AddProjectCreatorToTemplate(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing permission", func() {
				resp, err := client.Permissions.AddProjectCreatorToTemplate(&sonargo.PermissionsAddProjectCreatorToTemplateOption{
					TemplateName: testTemplateName,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing template identifier", func() {
				resp, err := client.Permissions.AddProjectCreatorToTemplate(&sonargo.PermissionsAddProjectCreatorToTemplateOption{
					Permission: "admin",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("RemoveProjectCreatorFromTemplate", func() {
		var testTemplateName string

		BeforeEach(func() {
			testTemplateName = helpers.UniqueResourceName("tpl-rmcreator")

			_, _, err := client.Permissions.CreateTemplate(&sonargo.PermissionsCreateTemplateOption{
				Name: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("template", testTemplateName, func() error {
				_, err := client.Permissions.DeleteTemplate(&sonargo.PermissionsDeleteTemplateOption{
					TemplateName: testTemplateName,
				})
				return err
			})

			// Add project creator first
			_, err = client.Permissions.AddProjectCreatorToTemplate(&sonargo.PermissionsAddProjectCreatorToTemplateOption{
				Permission:   "admin",
				TemplateName: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should remove project creator from template", func() {
			resp, err := client.Permissions.RemoveProjectCreatorFromTemplate(&sonargo.PermissionsRemoveProjectCreatorFromTemplateOption{
				Permission:   "admin",
				TemplateName: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify by searching templates
			result, _, err := client.Permissions.SearchTemplates(&sonargo.PermissionsSearchTemplatesOption{
				Query: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())
			for _, t := range result.PermissionTemplates {
				if t.Name == testTemplateName {
					for _, p := range t.Permissions {
						if p.Key == "admin" {
							Expect(p.WithProjectCreator).To(BeFalse())
						}
					}
					break
				}
			}
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Permissions.RemoveProjectCreatorFromTemplate(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing permission", func() {
				resp, err := client.Permissions.RemoveProjectCreatorFromTemplate(&sonargo.PermissionsRemoveProjectCreatorFromTemplateOption{
					TemplateName: testTemplateName,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing template identifier", func() {
				resp, err := client.Permissions.RemoveProjectCreatorFromTemplate(&sonargo.PermissionsRemoveProjectCreatorFromTemplateOption{
					Permission: "admin",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Apply Template
	// =========================================================================
	Describe("ApplyTemplate", func() {
		var testProjectKey string
		var testTemplateName string
		var testUserLogin string

		BeforeEach(func() {
			testProjectKey = helpers.UniqueResourceName("proj-apply")
			testTemplateName = helpers.UniqueResourceName("tpl-apply")
			testUserLogin = helpers.UniqueResourceName("user-apply")

			// Create private project (codeviewer permission only works on private projects)
			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "Apply Template Test",
				Project:    testProjectKey,
				Visibility: "private",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", testProjectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: testProjectKey,
				})
				return err
			})

			// Create template
			_, _, err = client.Permissions.CreateTemplate(&sonargo.PermissionsCreateTemplateOption{
				Name: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("template", testTemplateName, func() error {
				_, err := client.Permissions.DeleteTemplate(&sonargo.PermissionsDeleteTemplateOption{
					TemplateName: testTemplateName,
				})
				return err
			})

			// Create user and add to template
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err = client.Users.Create(&sonargo.UsersCreateOption{
				Login:    testUserLogin,
				Name:     "Apply Template User",
				Password: "SecurePassword123!",
				Local:    true,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("user", testUserLogin, func() error {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, _, err := client.Users.Deactivate(&sonargo.UsersDeactivateOption{
					Login:     testUserLogin,
					Anonymize: true,
				})
				return err
			})

			// Add user to template
			_, err = client.Permissions.AddUserToTemplate(&sonargo.PermissionsAddUserToTemplateOption{
				Login:        testUserLogin,
				Permission:   "codeviewer",
				TemplateName: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should apply template to project", func() {
			resp, err := client.Permissions.ApplyTemplate(&sonargo.PermissionsApplyTemplateOption{
				ProjectKey:   testProjectKey,
				TemplateName: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify user has permission on project
			result, _, err := client.Permissions.Users(&sonargo.PermissionsUsersOption{
				ProjectKey: testProjectKey,
				Permission: "codeviewer",
			})
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, u := range result.Users {
				if u.Login == testUserLogin {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Permissions.ApplyTemplate(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing project identifier", func() {
				resp, err := client.Permissions.ApplyTemplate(&sonargo.PermissionsApplyTemplateOption{
					TemplateName: testTemplateName,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing template identifier", func() {
				resp, err := client.Permissions.ApplyTemplate(&sonargo.PermissionsApplyTemplateOption{
					ProjectKey: testProjectKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("BulkApplyTemplate", func() {
		var testProjectKey1 string
		var testProjectKey2 string
		var testTemplateName string

		BeforeEach(func() {
			testProjectKey1 = helpers.UniqueResourceName("proj-bulk1")
			testProjectKey2 = helpers.UniqueResourceName("proj-bulk2")
			testTemplateName = helpers.UniqueResourceName("tpl-bulk")

			// Create projects
			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    "Bulk Apply Test 1",
				Project: testProjectKey1,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", testProjectKey1, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: testProjectKey1,
				})
				return err
			})

			_, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    "Bulk Apply Test 2",
				Project: testProjectKey2,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", testProjectKey2, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: testProjectKey2,
				})
				return err
			})

			// Create template
			_, _, err = client.Permissions.CreateTemplate(&sonargo.PermissionsCreateTemplateOption{
				Name: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("template", testTemplateName, func() error {
				_, err := client.Permissions.DeleteTemplate(&sonargo.PermissionsDeleteTemplateOption{
					TemplateName: testTemplateName,
				})
				return err
			})
		})

		It("should bulk apply template to projects", func() {
			resp, err := client.Permissions.BulkApplyTemplate(&sonargo.PermissionsBulkApplyTemplateOption{
				TemplateName: testTemplateName,
				Projects:     []string{testProjectKey1, testProjectKey2},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		It("should bulk apply template with query", func() {
			// Get prefix part for query
			queryPrefix := strings.TrimPrefix(testProjectKey1, helpers.E2EResourcePrefix)
			queryPrefix = helpers.E2EResourcePrefix + queryPrefix[:10] // Take first part

			resp, err := client.Permissions.BulkApplyTemplate(&sonargo.PermissionsBulkApplyTemplateOption{
				TemplateName: testTemplateName,
				Query:        queryPrefix,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Permissions.BulkApplyTemplate(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing template identifier", func() {
				resp, err := client.Permissions.BulkApplyTemplate(&sonargo.PermissionsBulkApplyTemplateOption{
					Projects: []string{testProjectKey1},
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid qualifier", func() {
				resp, err := client.Permissions.BulkApplyTemplate(&sonargo.PermissionsBulkApplyTemplateOption{
					TemplateName: testTemplateName,
					Qualifiers:   "INVALID",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("SetDefaultTemplate", func() {
		var testTemplateName string

		BeforeEach(func() {
			testTemplateName = helpers.UniqueResourceName("tpl-default")

			_, _, err := client.Permissions.CreateTemplate(&sonargo.PermissionsCreateTemplateOption{
				Name: testTemplateName,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("template", testTemplateName, func() error {
				_, err := client.Permissions.DeleteTemplate(&sonargo.PermissionsDeleteTemplateOption{
					TemplateName: testTemplateName,
				})
				return err
			})
		})

		It("should set default template for projects", func() {
			// First, get current default template to restore later
			searchResult, _, err := client.Permissions.SearchTemplates(nil)
			Expect(err).NotTo(HaveOccurred())

			var originalDefaultID string
			for _, dt := range searchResult.DefaultTemplates {
				if dt.Qualifier == "TRK" {
					originalDefaultID = dt.TemplateID
					break
				}
			}

			// Set our template as default
			resp, err := client.Permissions.SetDefaultTemplate(&sonargo.PermissionsSetDefaultTemplateOption{
				TemplateName: testTemplateName,
				Qualifier:    "TRK",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify it's now the default
			searchResult, _, err = client.Permissions.SearchTemplates(nil)
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, t := range searchResult.PermissionTemplates {
				if t.Name == testTemplateName {
					for _, dt := range searchResult.DefaultTemplates {
						if dt.TemplateID == t.ID && dt.Qualifier == "TRK" {
							found = true
							break
						}
					}
					break
				}
			}
			Expect(found).To(BeTrue())

			// Restore original default
			if originalDefaultID != "" {
				_, _ = client.Permissions.SetDefaultTemplate(&sonargo.PermissionsSetDefaultTemplateOption{
					TemplateID: originalDefaultID,
					Qualifier:  "TRK",
				})
			}
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Permissions.SetDefaultTemplate(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing template identifier", func() {
				resp, err := client.Permissions.SetDefaultTemplate(&sonargo.PermissionsSetDefaultTemplateOption{
					Qualifier: "TRK",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid qualifier", func() {
				resp, err := client.Permissions.SetDefaultTemplate(&sonargo.PermissionsSetDefaultTemplateOption{
					TemplateName: testTemplateName,
					Qualifier:    "INVALID",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Permission Lifecycle
	// =========================================================================
	Describe("Permission Lifecycle", func() {
		It("should complete full user permission lifecycle", func() {
			projectKey := helpers.UniqueResourceName("proj-lifecycle")
			userLogin := helpers.UniqueResourceName("user-lifecycle")

			// Step 1: Create private project (codeviewer permission only works on private projects)
			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "Lifecycle Test Project",
				Project:    projectKey,
				Visibility: "private",
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 2: Create user
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err = client.Users.Create(&sonargo.UsersCreateOption{
				Login:    userLogin,
				Name:     "Lifecycle Test User",
				Password: "SecurePassword123!",
				Local:    true,
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 3: Add permission
			_, err = client.Permissions.AddUser(&sonargo.PermissionsAddUserOption{
				Login:      userLogin,
				Permission: "codeviewer",
				ProjectKey: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 4: Verify permission
			result, _, err := client.Permissions.Users(&sonargo.PermissionsUsersOption{
				ProjectKey: projectKey,
				Permission: "codeviewer",
			})
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, u := range result.Users {
				if u.Login == userLogin {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())

			// Step 5: Remove permission
			_, err = client.Permissions.RemoveUser(&sonargo.PermissionsRemoveUserOption{
				Login:      userLogin,
				Permission: "codeviewer",
				ProjectKey: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 6: Verify permission removed
			result, _, err = client.Permissions.Users(&sonargo.PermissionsUsersOption{
				ProjectKey: projectKey,
				Permission: "codeviewer",
			})
			Expect(err).NotTo(HaveOccurred())
			for _, u := range result.Users {
				Expect(u.Login).NotTo(Equal(userLogin))
			}

			// Cleanup
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, _ = client.Users.Deactivate(&sonargo.UsersDeactivateOption{
				Login:     userLogin,
				Anonymize: true,
			})
			_, _ = client.Projects.Delete(&sonargo.ProjectsDeleteOption{
				Project: projectKey,
			})
		})

		It("should complete full template lifecycle", func() {
			templateName := helpers.UniqueResourceName("tpl-lifecycle")
			groupName := helpers.UniqueResourceName("grp-lifecycle")
			projectKey := helpers.UniqueResourceName("proj-tpllife")

			// Step 1: Create template
			_, _, err := client.Permissions.CreateTemplate(&sonargo.PermissionsCreateTemplateOption{
				Name:        templateName,
				Description: "Lifecycle test template",
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 2: Create group
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err = client.UserGroups.Create(&sonargo.UserGroupsCreateOption{
				Name: groupName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 3: Add group to template
			_, err = client.Permissions.AddGroupToTemplate(&sonargo.PermissionsAddGroupToTemplateOption{
				GroupName:    groupName,
				Permission:   "codeviewer",
				TemplateName: templateName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 4: Add project creator to template
			_, err = client.Permissions.AddProjectCreatorToTemplate(&sonargo.PermissionsAddProjectCreatorToTemplateOption{
				Permission:   "admin",
				TemplateName: templateName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 5: Create private project (codeviewer permission only works on private projects)
			_, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "Template Lifecycle Project",
				Project:    projectKey,
				Visibility: "private",
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 6: Apply template to project
			_, err = client.Permissions.ApplyTemplate(&sonargo.PermissionsApplyTemplateOption{
				ProjectKey:   projectKey,
				TemplateName: templateName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 7: Verify group has permission on project
			result, _, err := client.Permissions.Groups(&sonargo.PermissionsGroupsOption{
				ProjectKey: projectKey,
				Permission: "codeviewer",
			})
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, g := range result.Groups {
				if g.Name == groupName {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())

			// Cleanup
			_, _ = client.Projects.Delete(&sonargo.ProjectsDeleteOption{
				Project: projectKey,
			})
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _ = client.UserGroups.Delete(&sonargo.UserGroupsDeleteOption{
				Name: groupName,
			})
			_, _ = client.Permissions.DeleteTemplate(&sonargo.PermissionsDeleteTemplateOption{
				TemplateName: templateName,
			})
		})
	})
})
