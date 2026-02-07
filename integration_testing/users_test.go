package integration_testing_test

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("Users Service", Ordered, func() {
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

	Describe("Current", func() {
		It("should return the current authenticated user", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			user, resp, err := client.Users.Current()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(user).NotTo(BeNil())
			Expect(user.Login).NotTo(BeEmpty())
			Expect(user.IsLoggedIn).To(BeTrue())
		})

		It("should return user permissions for admin", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			user, resp, err := client.Users.Current()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(user).NotTo(BeNil())
			// Admin user should have global permissions
			Expect(user.Permissions.Global).NotTo(BeEmpty())
		})

		It("should return groups for the current user", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			user, resp, err := client.Users.Current()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(user).NotTo(BeNil())
			// User should belong to at least one group
			Expect(user.Groups).NotTo(BeEmpty())
		})
	})

	Describe("Search", func() {
		Context("without options", func() {
			It("should return list of users", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				result, resp, err := client.Users.Search(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Users).NotTo(BeEmpty())
			})
		})

		Context("with query filter", func() {
			It("should filter users by query", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				result, resp, err := client.Users.Search(&sonar.UsersSearchOption{
					Q: "admin",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				// Should find the admin user
				found := false
				for _, u := range result.Users {
					if u.Login == "admin" {
						found = true
						break
					}
				}
				Expect(found).To(BeTrue())
			})

			It("should return empty list for non-matching query", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				result, resp, err := client.Users.Search(&sonar.UsersSearchOption{
					Q: "nonexistentuserxyz123",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Users).To(BeEmpty())
			})
		})

		Context("with pagination", func() {
			It("should support page size", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				result, resp, err := client.Users.Search(&sonar.UsersSearchOption{
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
			It("should reject invalid page size", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, _, err := client.Users.Search(&sonar.UsersSearchOption{
					PaginationArgs: sonar.PaginationArgs{
						PageSize: 1000, // Too large
					},
				})
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Create", func() {
		Context("with required fields only", func() {
			It("should create a new user", func() {
				login := helpers.UniqueResourceName("user")

				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				result, resp, err := client.Users.Create(&sonar.UsersCreateOption{
					Login:    login,
					Name:     "E2E Test User",
					Password: "SecurePassword123!",
					Local:    true,
				})

				Expect(err).NotTo(HaveOccurred())

				// Register cleanup
				cleanup.RegisterCleanup("user", login, func() error {
					//nolint:staticcheck // Using deprecated API until v2 API is implemented
					_, _, err := client.Users.Deactivate(&sonar.UsersDeactivateOption{
						Login:     login,
						Anonymize: true,
					})
					return err
				})
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.User.Login).To(Equal(login))
				Expect(result.User.Name).To(Equal("E2E Test User"))
				Expect(result.User.Active).To(BeTrue())
			})
		})

		Context("with all optional fields", func() {
			It("should create a user with email and SCM accounts", func() {
				login := helpers.UniqueResourceName("user-full")
				// Use unique SCM accounts to avoid conflicts with other tests
				scmAccount1 := fmt.Sprintf("github:%s", login)
				scmAccount2 := fmt.Sprintf("gitlab:%s", login)

				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				result, resp, err := client.Users.Create(&sonar.UsersCreateOption{
					Login:       login,
					Name:        "E2E Full User",
					Email:       "e2e-test@example.com",
					Password:    "SecurePassword123!",
					Local:       true,
					ScmAccounts: []string{scmAccount1, scmAccount2},
				})

				Expect(err).NotTo(HaveOccurred())

				// Register cleanup
				cleanup.RegisterCleanup("user", login, func() error {
					//nolint:staticcheck // Using deprecated API until v2 API is implemented
					_, _, err := client.Users.Deactivate(&sonar.UsersDeactivateOption{
						Login:     login,
						Anonymize: true,
					})
					return err
				})
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.User.Login).To(Equal(login))
				Expect(result.User.Email).To(Equal("e2e-test@example.com"))
				Expect(result.User.ScmAccounts).To(ContainElements(scmAccount1, scmAccount2))
			})
		})

		Context("duplicate user", func() {
			It("should fail when creating user with existing login", func() {
				login := helpers.UniqueResourceName("user-dup")

				// Create first user
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, _, err := client.Users.Create(&sonar.UsersCreateOption{
					Login:    login,
					Name:     "First User",
					Password: "SecurePassword123!",
					Local:    true,
				})
				Expect(err).NotTo(HaveOccurred())

				// Register cleanup
				cleanup.RegisterCleanup("user", login, func() error {
					//nolint:staticcheck // Using deprecated API until v2 API is implemented
					_, _, err := client.Users.Deactivate(&sonar.UsersDeactivateOption{
						Login:     login,
						Anonymize: true,
					})
					return err
				})

				// Try to create duplicate
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.Users.Create(&sonar.UsersCreateOption{
					Login:    login,
					Name:     "Duplicate User",
					Password: "AnotherPassword123!",
					Local:    true,
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
				_, resp, err := client.Users.Create(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing login", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.Users.Create(&sonar.UsersCreateOption{
					Name:     "Test User",
					Password: "SecurePassword123!",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing name", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.Users.Create(&sonar.UsersCreateOption{
					Login:    "testuser",
					Password: "SecurePassword123!",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with login too short", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.Users.Create(&sonar.UsersCreateOption{
					Login:    "x",
					Name:     "Test User",
					Password: "SecurePassword123!",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing password for local user", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.Users.Create(&sonar.UsersCreateOption{
					Login: "testuser",
					Name:  "Test User",
					Local: true,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with password too short", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.Users.Create(&sonar.UsersCreateOption{
					Login:    "testuser",
					Name:     "Test User",
					Password: "short",
					Local:    true,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("Update", func() {
		var testUserLogin string

		BeforeEach(func() {
			testUserLogin = helpers.UniqueResourceName("user-update")
			// Capture login value for cleanup closure
			loginToCleanup := testUserLogin

			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err := client.Users.Create(&sonar.UsersCreateOption{
				Login:    testUserLogin,
				Name:     "Original Name",
				Email:    "original@example.com",
				Password: "SecurePassword123!",
				Local:    true,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("user", loginToCleanup, func() error {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, _, err := client.Users.Deactivate(&sonar.UsersDeactivateOption{
					Login:     loginToCleanup,
					Anonymize: true,
				})
				return err
			})
		})

		It("should update user name", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			result, resp, err := client.Users.Update(&sonar.UsersUpdateOption{
				Login: testUserLogin,
				Name:  "Updated Name",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.User.Name).To(Equal("Updated Name"))
		})

		It("should update user email", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			result, resp, err := client.Users.Update(&sonar.UsersUpdateOption{
				Login: testUserLogin,
				Email: "updated@example.com",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.User.Email).To(Equal("updated@example.com"))
		})

		It("should update SCM accounts", func() {
			// Use unique SCM accounts to avoid conflicts with other tests
			scmAccount1 := fmt.Sprintf("github:%s", testUserLogin)
			scmAccount2 := fmt.Sprintf("bitbucket:%s", testUserLogin)

			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			result, resp, err := client.Users.Update(&sonar.UsersUpdateOption{
				Login:       testUserLogin,
				ScmAccounts: []string{scmAccount1, scmAccount2},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.User.ScmAccounts).To(ContainElements(scmAccount1, scmAccount2))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.Users.Update(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing login", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.Users.Update(&sonar.UsersUpdateOption{
					Name: "Some Name",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail for non-existent user", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.Users.Update(&sonar.UsersUpdateOption{
					Login: "nonexistentuser12345",
					Name:  "New Name",
				})
				Expect(err).To(HaveOccurred())
				if resp != nil {
					Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				}
			})
		})
	})

	Describe("ChangePassword", func() {
		var testUserLogin string

		BeforeEach(func() {
			testUserLogin = helpers.UniqueResourceName("user-pwd")
			// Capture login value for cleanup closure
			loginToCleanup := testUserLogin

			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err := client.Users.Create(&sonar.UsersCreateOption{
				Login:    testUserLogin,
				Name:     "Password Test User",
				Password: "OldPassword123!",
				Local:    true,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("user", loginToCleanup, func() error {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, _, err := client.Users.Deactivate(&sonar.UsersDeactivateOption{
					Login:     loginToCleanup,
					Anonymize: true,
				})
				return err
			})
		})

		It("should change user password", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			resp, err := client.Users.ChangePassword(&sonar.UsersChangePasswordOption{
				Login:    testUserLogin,
				Password: "NewPassword456!",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify new password works by creating a client with the new credentials
			cfg := helpers.LoadConfig()
			newClient, err := sonar.NewClient(nil,
				sonar.WithBaseURL(helpers.NormalizeBaseURL(cfg.BaseURL)),
				sonar.WithBasicAuth(testUserLogin, "NewPassword456!"),
			)
			Expect(err).NotTo(HaveOccurred())

			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			currentUser, _, err := newClient.Users.Current()
			Expect(err).NotTo(HaveOccurred())
			Expect(currentUser.Login).To(Equal(testUserLogin))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.Users.ChangePassword(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing login", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.Users.ChangePassword(&sonar.UsersChangePasswordOption{
					Password: "NewPassword123!",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing password", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.Users.ChangePassword(&sonar.UsersChangePasswordOption{
					Login: testUserLogin,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with password too short", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.Users.ChangePassword(&sonar.UsersChangePasswordOption{
					Login:    testUserLogin,
					Password: "short",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("UpdateLogin", func() {
		It("should update user login", func() {
			oldLogin := helpers.UniqueResourceName("user-old")
			newLogin := helpers.UniqueResourceName("user-new")

			// Create user
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err := client.Users.Create(&sonar.UsersCreateOption{
				Login:    oldLogin,
				Name:     "Login Update User",
				Password: "SecurePassword123!",
				Local:    true,
			})
			Expect(err).NotTo(HaveOccurred())

			// Update login
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			resp, err := client.Users.UpdateLogin(&sonar.UsersUpdateLoginOption{
				Login:    oldLogin,
				NewLogin: newLogin,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Register cleanup with new login
			cleanup.RegisterCleanup("user", newLogin, func() error {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, _, err := client.Users.Deactivate(&sonar.UsersDeactivateOption{
					Login:     newLogin,
					Anonymize: true,
				})
				return err
			})

			// Verify new login exists
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			result, _, err := client.Users.Search(&sonar.UsersSearchOption{
				Q: newLogin,
			})
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, u := range result.Users {
				if u.Login == newLogin {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.Users.UpdateLogin(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing login", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.Users.UpdateLogin(&sonar.UsersUpdateLoginOption{
					NewLogin: "newlogin",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing new login", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.Users.UpdateLogin(&sonar.UsersUpdateLoginOption{
					Login: "oldlogin",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with new login too short", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.Users.UpdateLogin(&sonar.UsersUpdateLoginOption{
					Login:    "someuser",
					NewLogin: "x",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("Groups", func() {
		It("should return groups for a user", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			result, resp, err := client.Users.Groups(&sonar.UsersGroupsOption{
				Login: "admin",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			// Admin should belong to at least one group
			Expect(result.Groups).NotTo(BeEmpty())
		})

		It("should filter groups with query", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			result, resp, err := client.Users.Groups(&sonar.UsersGroupsOption{
				Login: "admin",
				Q:     "sonar",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.Users.Groups(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing login", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.Users.Groups(&sonar.UsersGroupsOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("IdentityProviders", func() {
		It("should return list of identity providers", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			result, resp, err := client.Users.IdentityProviders()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			// IdentityProviders may be empty if no external providers are configured
			// The API returns an empty list in a default SonarQube installation
		})

		It("should successfully call the API", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			result, resp, err := client.Users.IdentityProviders()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			// Verify we can iterate over the list (even if empty)
			for _, p := range result.IdentityProviders {
				Expect(p.Key).NotTo(BeEmpty())
			}
		})
	})

	Describe("SetHomepage", func() {
		It("should set homepage to PROJECTS", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			resp, err := client.Users.SetHomepage(&sonar.UsersSetHomepageOption{
				Type: "PROJECTS",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		It("should set homepage to ISSUES", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			resp, err := client.Users.SetHomepage(&sonar.UsersSetHomepageOption{
				Type: "ISSUES",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.Users.SetHomepage(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing type", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.Users.SetHomepage(&sonar.UsersSetHomepageOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid type", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.Users.SetHomepage(&sonar.UsersSetHomepageOption{
					Type: "INVALID_TYPE",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("DismissNotice", func() {
		It("should dismiss a notice", func() {
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			resp, err := client.Users.DismissNotice(&sonar.UsersDismissNoticeOption{
				Notice: "educationPrinciples",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		It("should handle already dismissed notice gracefully", func() {
			// Dismiss the same notice twice - should succeed silently
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			resp, err := client.Users.DismissNotice(&sonar.UsersDismissNoticeOption{
				Notice: "sonarlintAd",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Dismiss again
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			resp, err = client.Users.DismissNotice(&sonar.UsersDismissNoticeOption{
				Notice: "sonarlintAd",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.Users.DismissNotice(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing notice", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.Users.DismissNotice(&sonar.UsersDismissNoticeOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid notice type", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.Users.DismissNotice(&sonar.UsersDismissNoticeOption{
					Notice: "invalidNoticeType",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("Deactivate", func() {
		It("should deactivate a user", func() {
			login := helpers.UniqueResourceName("user-deact")

			// Create user
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err := client.Users.Create(&sonar.UsersCreateOption{
				Login:    login,
				Name:     "Deactivate Test User",
				Password: "SecurePassword123!",
				Local:    true,
			})
			Expect(err).NotTo(HaveOccurred())

			// Deactivate user
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			result, resp, err := client.Users.Deactivate(&sonar.UsersDeactivateOption{
				Login: login,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.User.Active).To(BeFalse())
			Expect(result.User.Login).To(Equal(login))

			// Verify user is deactivated
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			searchResult, _, err := client.Users.Search(&sonar.UsersSearchOption{
				Q:           login,
				Deactivated: true,
			})
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, u := range searchResult.Users {
				if u.Login == login {
					found = true
					Expect(u.Active).To(BeFalse())
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		It("should deactivate and anonymize a user", func() {
			login := helpers.UniqueResourceName("user-anon")

			// Create user
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err := client.Users.Create(&sonar.UsersCreateOption{
				Login:    login,
				Name:     "Anonymize Test User",
				Email:    "anon@example.com",
				Password: "SecurePassword123!",
				Local:    true,
			})
			Expect(err).NotTo(HaveOccurred())

			// Deactivate with anonymize
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			result, resp, err := client.Users.Deactivate(&sonar.UsersDeactivateOption{
				Login:     login,
				Anonymize: true,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.User.Active).To(BeFalse())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.Users.Deactivate(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing login", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.Users.Deactivate(&sonar.UsersDeactivateOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail for non-existent user", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, resp, err := client.Users.Deactivate(&sonar.UsersDeactivateOption{
					Login: "nonexistentuser12345",
				})
				Expect(err).To(HaveOccurred())
				if resp != nil {
					Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				}
			})
		})
	})

	Describe("Anonymize", func() {
		It("should anonymize a deactivated user", func() {
			login := helpers.UniqueResourceName("user-anon2")

			// Create user
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err := client.Users.Create(&sonar.UsersCreateOption{
				Login:    login,
				Name:     "Anonymize Test User 2",
				Email:    "anon2@example.com",
				Password: "SecurePassword123!",
				Local:    true,
			})
			Expect(err).NotTo(HaveOccurred())

			// Deactivate first (required before anonymize)
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err = client.Users.Deactivate(&sonar.UsersDeactivateOption{
				Login: login,
			})
			Expect(err).NotTo(HaveOccurred())

			// Anonymize
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			resp, err := client.Users.Anonymize(&sonar.UsersAnonymizeOption{
				Login: login,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.Users.Anonymize(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing login", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.Users.Anonymize(&sonar.UsersAnonymizeOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("UpdateIdentityProvider", func() {
		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.Users.UpdateIdentityProvider(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing login", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.Users.UpdateIdentityProvider(&sonar.UsersUpdateIdentityProviderOption{
					NewExternalProvider: "sonarqube",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing new external provider", func() {
				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				resp, err := client.Users.UpdateIdentityProvider(&sonar.UsersUpdateIdentityProviderOption{
					Login: "someuser",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("User Lifecycle", func() {
		It("should complete full create/update/search/deactivate cycle", func() {
			login := helpers.UniqueResourceName("user-lifecycle")
			timestamp := fmt.Sprintf("%d", time.Now().UnixNano())

			// Step 1: Create user
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			createResult, _, err := client.Users.Create(&sonar.UsersCreateOption{
				Login:    login,
				Name:     "Lifecycle User " + timestamp,
				Email:    "lifecycle" + timestamp + "@example.com",
				Password: "SecurePassword123!",
				Local:    true,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(createResult.User.Login).To(Equal(login))
			Expect(createResult.User.Active).To(BeTrue())

			// Step 2: Search and verify
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			searchResult, _, err := client.Users.Search(&sonar.UsersSearchOption{
				Q: login,
			})
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, u := range searchResult.Users {
				if u.Login == login {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())

			// Step 3: Update user
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			updateResult, _, err := client.Users.Update(&sonar.UsersUpdateOption{
				Login: login,
				Name:  "Updated Lifecycle User",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(updateResult.User.Name).To(Equal("Updated Lifecycle User"))

			// Step 4: Get user groups
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			_, _, err = client.Users.Groups(&sonar.UsersGroupsOption{
				Login: login,
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 5: Deactivate and anonymize
			//nolint:staticcheck // Using deprecated API until v2 API is implemented
			deactivateResult, _, err := client.Users.Deactivate(&sonar.UsersDeactivateOption{
				Login:     login,
				Anonymize: true,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(deactivateResult.User.Active).To(BeFalse())
		})
	})
})
