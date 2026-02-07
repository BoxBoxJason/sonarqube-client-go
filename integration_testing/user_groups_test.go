package integration_testing_test

import (
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("UserGroups Service", Ordered, func() {
	var (
		client  *sonar.Client
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

	Describe("Search", func() {
		Context("with default options", func() {
			It("should return list of groups", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				result, resp, err := client.UserGroups.Search(&sonar.UserGroupsSearchOption{})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				// Default SonarQube has at least sonar-users and sonar-administrators
				Expect(result.Groups).NotTo(BeEmpty())
			})
		})

		Context("with query filter", func() {
			It("should filter groups by query", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				result, resp, err := client.UserGroups.Search(&sonar.UserGroupsSearchOption{
					Query: "admin",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				// Should find sonar-administrators
				found := false
				for _, g := range result.Groups {
					if strings.Contains(g.Name, "admin") {
						found = true
						break
					}
				}
				Expect(found).To(BeTrue())
			})

			It("should return empty list for non-matching query", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				result, resp, err := client.UserGroups.Search(&sonar.UserGroupsSearchOption{
					Query: "nonexistentgroupxyz123",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Groups).To(BeEmpty())
			})
		})

		Context("with pagination", func() {
			It("should support page size", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				result, resp, err := client.UserGroups.Search(&sonar.UserGroupsSearchOption{
					PaginationArgs: sonar.PaginationArgs{
						PageSize: 1,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(len(result.Groups)).To(BeNumerically("<=", 1))
			})
		})

		Context("with fields selection", func() {
			It("should return specified fields", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				result, resp, err := client.UserGroups.Search(&sonar.UserGroupsSearchOption{
					Fields: []string{"name", "description", "membersCount"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Groups).NotTo(BeEmpty())
				// Verify at least name is populated
				for _, g := range result.Groups {
					Expect(g.Name).NotTo(BeEmpty())
				}
			})
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.UserGroups.Search(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid page size", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, _, err := client.UserGroups.Search(&sonar.UserGroupsSearchOption{
					PaginationArgs: sonar.PaginationArgs{
						PageSize: 1000, // Too large
					},
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with invalid field", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.UserGroups.Search(&sonar.UserGroupsSearchOption{
					Fields: []string{"invalid_field"},
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("Create", func() {
		Context("with required fields only", func() {
			It("should create a new group", func() {
				groupName := helpers.UniqueResourceName("group")

				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				result, resp, err := client.UserGroups.Create(&sonar.UserGroupsCreateOption{
					Name: groupName,
				})

				Expect(err).NotTo(HaveOccurred())

				// Register cleanup
				cleanup.RegisterCleanup("group", groupName, func() error {
					//nolint:staticcheck // Using deprecated API until v2 API is implemented
					_, err := client.UserGroups.Delete(&sonar.UserGroupsDeleteOption{
						Name: groupName,
					})
					return err
				})

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Group.Name).To(Equal(groupName))
				Expect(result.Group.ID).NotTo(BeEmpty())
			})
		})

		Context("with optional description", func() {
			It("should create a group with description", func() {
				groupName := helpers.UniqueResourceName("group-desc")

				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				result, resp, err := client.UserGroups.Create(&sonar.UserGroupsCreateOption{
					Name:        groupName,
					Description: "E2E test group with description",
				})

				Expect(err).NotTo(HaveOccurred())

				// Register cleanup
				cleanup.RegisterCleanup("group", groupName, func() error {
					//nolint:staticcheck // Using deprecated API until v2 API is implemented
					_, err := client.UserGroups.Delete(&sonar.UserGroupsDeleteOption{
						Name: groupName,
					})
					return err
				})

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Group.Name).To(Equal(groupName))
				Expect(result.Group.Description).To(Equal("E2E test group with description"))
			})
		})

		Context("duplicate group", func() {
			It("should fail when creating group with existing name", func() {
				groupName := helpers.UniqueResourceName("group-dup")

				// Create first group
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, _, err := client.UserGroups.Create(&sonar.UserGroupsCreateOption{
					Name: groupName,
				})
				Expect(err).NotTo(HaveOccurred())

				// Register cleanup
				cleanup.RegisterCleanup("group", groupName, func() error {
					//nolint:staticcheck // Using deprecated API until v2 API is implemented
					_, err := client.UserGroups.Delete(&sonar.UserGroupsDeleteOption{
						Name: groupName,
					})
					return err
				})

				// Try to create duplicate
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.UserGroups.Create(&sonar.UserGroupsCreateOption{
					Name: groupName,
				})
				Expect(err).To(HaveOccurred())
				if resp != nil {
					Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
				}
			})
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.UserGroups.Create(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing name", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.UserGroups.Create(&sonar.UserGroupsCreateOption{
					Description: "Test description",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with name too long", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.UserGroups.Create(&sonar.UserGroupsCreateOption{
					Name: strings.Repeat("a", sonar.MaxGroupNameLength+1),
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with description too long", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.UserGroups.Create(&sonar.UserGroupsCreateOption{
					Name:        "validname",
					Description: strings.Repeat("a", sonar.MaxGroupDescriptionLength+1),
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("Update", func() {
		var testGroupName string

		BeforeEach(func() {
			testGroupName = helpers.UniqueResourceName("group-update")
			// Capture name value for cleanup closure
			nameToCleanup := testGroupName

			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err := client.UserGroups.Create(&sonar.UserGroupsCreateOption{
				Name:        testGroupName,
				Description: "Original description",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("group", nameToCleanup, func() error {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, err := client.UserGroups.Delete(&sonar.UserGroupsDeleteOption{
					Name: nameToCleanup,
				})
				return err
			})
		})

		It("should update group description", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			resp, err := client.UserGroups.Update(&sonar.UserGroupsUpdateOption{
				CurrentName: testGroupName,
				Description: "Updated description",
			})
			Expect(err).NotTo(HaveOccurred())
			// API may return 200 OK or 204 NoContent depending on SonarQube version
			Expect(resp.StatusCode).To(BeNumerically(">=", http.StatusOK))
			Expect(resp.StatusCode).To(BeNumerically("<", 300))

			// Verify update
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			result, _, err := client.UserGroups.Search(&sonar.UserGroupsSearchOption{
				Query:  testGroupName,
				Fields: []string{"name", "description"},
			})
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, g := range result.Groups {
				if g.Name == testGroupName {
					found = true
					Expect(g.Description).To(Equal("Updated description"))
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		It("should update group name", func() {
			newName := helpers.UniqueResourceName("group-renamed")

			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			resp, err := client.UserGroups.Update(&sonar.UserGroupsUpdateOption{
				CurrentName: testGroupName,
				Name:        newName,
			})
			Expect(err).NotTo(HaveOccurred())
			// API may return 200 OK or 204 NoContent depending on SonarQube version
			Expect(resp.StatusCode).To(BeNumerically(">=", http.StatusOK))
			Expect(resp.StatusCode).To(BeNumerically("<", 300))

			// Ensure existing cleanup uses the new group name instead of the old one.
			testGroupName = newName

			// Verify update
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			result, _, err := client.UserGroups.Search(&sonar.UserGroupsSearchOption{
				Query: newName,
			})
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, g := range result.Groups {
				if g.Name == newName {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.UserGroups.Update(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing current name", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.UserGroups.Update(&sonar.UserGroupsUpdateOption{
					Name: "newname",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with name too long", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.UserGroups.Update(&sonar.UserGroupsUpdateOption{
					CurrentName: testGroupName,
					Name:        strings.Repeat("a", sonar.MaxGroupNameLength+1),
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with description too long", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.UserGroups.Update(&sonar.UserGroupsUpdateOption{
					CurrentName: testGroupName,
					Description: strings.Repeat("a", sonar.MaxGroupDescriptionLength+1),
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail for non-existent group", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, err := client.UserGroups.Update(&sonar.UserGroupsUpdateOption{
					CurrentName: "nonexistentgroup12345",
					Description: "New description",
				})
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Delete", func() {
		It("should delete a group", func() {
			groupName := helpers.UniqueResourceName("group-del")

			// Create group
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err := client.UserGroups.Create(&sonar.UserGroupsCreateOption{
				Name: groupName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Delete group
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			resp, err := client.UserGroups.Delete(&sonar.UserGroupsDeleteOption{
				Name: groupName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify deletion
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			result, _, err := client.UserGroups.Search(&sonar.UserGroupsSearchOption{
				Query: groupName,
			})
			Expect(err).NotTo(HaveOccurred())
			for _, g := range result.Groups {
				Expect(g.Name).NotTo(Equal(groupName))
			}
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.UserGroups.Delete(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing name", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.UserGroups.Delete(&sonar.UserGroupsDeleteOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail for non-existent group", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, err := client.UserGroups.Delete(&sonar.UserGroupsDeleteOption{
					Name: "nonexistentgroup12345",
				})
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("AddUser", func() {
		var testGroupName string
		var testUserLogin string

		BeforeEach(func() {
			testGroupName = helpers.UniqueResourceName("group-adduser")
			testUserLogin = helpers.UniqueResourceName("user-grp")

			// Capture values for cleanup closures
			groupToCleanup := testGroupName
			userToCleanup := testUserLogin

			// Create group first
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err := client.UserGroups.Create(&sonar.UserGroupsCreateOption{
				Name: testGroupName,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("group", groupToCleanup, func() error {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, err := client.UserGroups.Delete(&sonar.UserGroupsDeleteOption{
					Name: groupToCleanup,
				})
				return err
			})

			// Create user
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err = client.Users.Create(&sonar.UsersCreateOption{
				Login:    testUserLogin,
				Name:     "Test User for Group",
				Password: "SecurePassword123!",
				Local:    true,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("user", userToCleanup, func() error {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, _, err := client.Users.Deactivate(&sonar.UsersDeactivateOption{
					Login:     userToCleanup,
					Anonymize: true,
				})
				return err
			})
		})

		It("should add user to group", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			resp, err := client.UserGroups.AddUser(&sonar.UserGroupsAddUserOption{
				Name:  testGroupName,
				Login: testUserLogin,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify user is in group
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			result, _, err := client.UserGroups.Users(&sonar.UserGroupsUsersOption{
				Name:     testGroupName,
				Selected: "selected",
			})
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, u := range result.Users {
				if u.Login == testUserLogin {
					found = true
					Expect(u.Selected).To(BeTrue())
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		It("should add user to group by login only", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			resp, err := client.UserGroups.AddUser(&sonar.UserGroupsAddUserOption{
				Name:  testGroupName,
				Login: testUserLogin,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.UserGroups.AddUser(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing group name", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.UserGroups.AddUser(&sonar.UserGroupsAddUserOption{
					Login: testUserLogin,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail for non-existent group", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, err := client.UserGroups.AddUser(&sonar.UserGroupsAddUserOption{
					Name:  "nonexistentgroup12345",
					Login: testUserLogin,
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail for non-existent user", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, err := client.UserGroups.AddUser(&sonar.UserGroupsAddUserOption{
					Name:  testGroupName,
					Login: "nonexistentuser12345",
				})
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("RemoveUser", func() {
		var testGroupName string
		var testUserLogin string

		BeforeEach(func() {
			testGroupName = helpers.UniqueResourceName("group-rmuser")
			testUserLogin = helpers.UniqueResourceName("user-rmgrp")

			// Capture values for cleanup closures
			groupToCleanup := testGroupName
			userToCleanup := testUserLogin

			// Create group first
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err := client.UserGroups.Create(&sonar.UserGroupsCreateOption{
				Name: testGroupName,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("group", groupToCleanup, func() error {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, err := client.UserGroups.Delete(&sonar.UserGroupsDeleteOption{
					Name: groupToCleanup,
				})
				return err
			})

			// Create user
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err = client.Users.Create(&sonar.UsersCreateOption{
				Login:    testUserLogin,
				Name:     "Test User for Remove",
				Password: "SecurePassword123!",
				Local:    true,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("user", userToCleanup, func() error {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, _, err := client.Users.Deactivate(&sonar.UsersDeactivateOption{
					Login:     userToCleanup,
					Anonymize: true,
				})
				return err
			})

			// Add user to group
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, err = client.UserGroups.AddUser(&sonar.UserGroupsAddUserOption{
				Name:  testGroupName,
				Login: testUserLogin,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should remove user from group", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			resp, err := client.UserGroups.RemoveUser(&sonar.UserGroupsRemoveUserOption{
				Name:  testGroupName,
				Login: testUserLogin,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify user is not in group
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			result, _, err := client.UserGroups.Users(&sonar.UserGroupsUsersOption{
				Name:     testGroupName,
				Selected: "selected",
			})
			Expect(err).NotTo(HaveOccurred())
			for _, u := range result.Users {
				Expect(u.Login).NotTo(Equal(testUserLogin))
			}
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.UserGroups.RemoveUser(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing group name", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.UserGroups.RemoveUser(&sonar.UserGroupsRemoveUserOption{
					Login: testUserLogin,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail for non-existent group", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, err := client.UserGroups.RemoveUser(&sonar.UserGroupsRemoveUserOption{
					Name:  "nonexistentgroup12345",
					Login: testUserLogin,
				})
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Users", func() {
		var testGroupName string
		var testUserLogin string

		BeforeEach(func() {
			testGroupName = helpers.UniqueResourceName("group-users")
			testUserLogin = helpers.UniqueResourceName("user-ingrp")

			// Capture values for cleanup closures
			groupToCleanup := testGroupName
			userToCleanup := testUserLogin

			// Create group first
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err := client.UserGroups.Create(&sonar.UserGroupsCreateOption{
				Name: testGroupName,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("group", groupToCleanup, func() error {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, err := client.UserGroups.Delete(&sonar.UserGroupsDeleteOption{
					Name: groupToCleanup,
				})
				return err
			})

			// Create user
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err = client.Users.Create(&sonar.UsersCreateOption{
				Login:    testUserLogin,
				Name:     "Test User in Group",
				Password: "SecurePassword123!",
				Local:    true,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("user", userToCleanup, func() error {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, _, err := client.Users.Deactivate(&sonar.UsersDeactivateOption{
					Login:     userToCleanup,
					Anonymize: true,
				})
				return err
			})

			// Add user to group
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, err = client.UserGroups.AddUser(&sonar.UserGroupsAddUserOption{
				Name:  testGroupName,
				Login: testUserLogin,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should list selected users in group", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			result, resp, err := client.UserGroups.Users(&sonar.UserGroupsUsersOption{
				Name:     testGroupName,
				Selected: "selected",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Users).NotTo(BeEmpty())

			// Find our test user
			found := false
			for _, u := range result.Users {
				if u.Login == testUserLogin {
					found = true
					Expect(u.Selected).To(BeTrue())
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		It("should list all users with membership info", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			result, resp, err := client.UserGroups.Users(&sonar.UserGroupsUsersOption{
				Name:     testGroupName,
				Selected: "all",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			// Should have at least our test user
			Expect(result.Users).NotTo(BeEmpty())
		})

		It("should list deselected users", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			result, resp, err := client.UserGroups.Users(&sonar.UserGroupsUsersOption{
				Name:     testGroupName,
				Selected: "deselected",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			// admin user should be deselected (not in our test group)
			for _, u := range result.Users {
				Expect(u.Selected).To(BeFalse())
			}
		})

		It("should filter users by query", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			result, resp, err := client.UserGroups.Users(&sonar.UserGroupsUsersOption{
				Name:     testGroupName,
				Query:    testUserLogin,
				Selected: "all",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			// Should find our test user
			found := false
			for _, u := range result.Users {
				if strings.Contains(u.Login, testUserLogin) {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		Context("with pagination", func() {
			It("should support page size", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				result, resp, err := client.UserGroups.Users(&sonar.UserGroupsUsersOption{
					Name: testGroupName,
					PaginationArgs: sonar.PaginationArgs{
						PageSize: 1,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(len(result.Users)).To(BeNumerically("<=", 1))
			})
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.UserGroups.Users(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing name", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.UserGroups.Users(&sonar.UserGroupsUsersOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid selected value", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.UserGroups.Users(&sonar.UserGroupsUsersOption{
					Name:     testGroupName,
					Selected: "invalid",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail for non-existent group", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, _, err := client.UserGroups.Users(&sonar.UserGroupsUsersOption{
					Name: "nonexistentgroup12345",
				})
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Group Lifecycle", func() {
		It("should complete full create/update/add-user/remove-user/delete cycle", func() {
			groupName := helpers.UniqueResourceName("group-lifecycle")
			userLogin := helpers.UniqueResourceName("user-lc")

			// Step 1: Create group
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			createResult, _, err := client.UserGroups.Create(&sonar.UserGroupsCreateOption{
				Name:        groupName,
				Description: "Lifecycle test group",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(createResult.Group.Name).To(Equal(groupName))
			Expect(createResult.Group.Description).To(Equal("Lifecycle test group"))
			Expect(createResult.Group.MembersCount).To(Equal(int64(0)))

			// Step 2: Create user
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err = client.Users.Create(&sonar.UsersCreateOption{
				Login:    userLogin,
				Name:     "Lifecycle Test User",
				Password: "SecurePassword123!",
				Local:    true,
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 3: Add user to group
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, err = client.UserGroups.AddUser(&sonar.UserGroupsAddUserOption{
				Name:  groupName,
				Login: userLogin,
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 4: Verify user in group
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			usersResult, _, err := client.UserGroups.Users(&sonar.UserGroupsUsersOption{
				Name:     groupName,
				Selected: "selected",
			})
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, u := range usersResult.Users {
				if u.Login == userLogin {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())

			// Step 5: Update group
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, err = client.UserGroups.Update(&sonar.UserGroupsUpdateOption{
				CurrentName: groupName,
				Description: "Updated lifecycle description",
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 6: Remove user from group
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, err = client.UserGroups.RemoveUser(&sonar.UserGroupsRemoveUserOption{
				Name:  groupName,
				Login: userLogin,
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 7: Verify user removed
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			usersResult, _, err = client.UserGroups.Users(&sonar.UserGroupsUsersOption{
				Name:     groupName,
				Selected: "selected",
			})
			Expect(err).NotTo(HaveOccurred())
			for _, u := range usersResult.Users {
				Expect(u.Login).NotTo(Equal(userLogin))
			}

			// Step 8: Cleanup - deactivate user first
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err = client.Users.Deactivate(&sonar.UsersDeactivateOption{
				Login:     userLogin,
				Anonymize: true,
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 9: Delete group
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, err = client.UserGroups.Delete(&sonar.UserGroupsDeleteOption{
				Name: groupName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Verify group deleted
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			searchResult, _, err := client.UserGroups.Search(&sonar.UserGroupsSearchOption{
				Query: groupName,
			})
			Expect(err).NotTo(HaveOccurred())
			for _, g := range searchResult.Groups {
				Expect(g.Name).NotTo(Equal(groupName))
			}
		})
	})

	Describe("Default Groups", func() {
		It("should find sonar-users group", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			result, resp, err := client.UserGroups.Search(&sonar.UserGroupsSearchOption{
				Query: "sonar-users",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			found := false
			for _, g := range result.Groups {
				if g.Name == "sonar-users" {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		It("should find sonar-administrators group", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			result, resp, err := client.UserGroups.Search(&sonar.UserGroupsSearchOption{
				Query: "sonar-administrators",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			found := false
			for _, g := range result.Groups {
				if g.Name == "sonar-administrators" {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		It("should not delete default groups", func() {
			// Try to delete sonar-users (default group) - should fail
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, err := client.UserGroups.Delete(&sonar.UserGroupsDeleteOption{
				Name: "sonar-users",
			})
			// Default groups cannot be deleted
			Expect(err).To(HaveOccurred())
		})
	})
})
