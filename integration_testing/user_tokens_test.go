package integration_testing_test

import (
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("UserTokens Service", Ordered, func() {
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

	Describe("Generate", func() {
		Context("with USER_TOKEN type (default)", func() {
			It("should generate a user token", func() {
				tokenName := helpers.UniqueResourceName("token")

				result, resp, err := client.UserTokens.Generate(&sonar.UserTokensGenerateOption{
					Name: tokenName,
				})

				Expect(err).NotTo(HaveOccurred())

				// Register cleanup
				cleanup.RegisterCleanup("token", tokenName, func() error {
					_, err := client.UserTokens.Revoke(&sonar.UserTokensRevokeOption{
						Name: tokenName,
					})
					return err
				})

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Name).To(Equal(tokenName))
				Expect(result.Token).NotTo(BeEmpty())
				// Default type should be USER_TOKEN
				Expect(result.Type).To(Equal("USER_TOKEN"))
			})

			It("should generate a token with explicit USER_TOKEN type", func() {
				tokenName := helpers.UniqueResourceName("token-user")

				result, resp, err := client.UserTokens.Generate(&sonar.UserTokensGenerateOption{
					Name: tokenName,
					Type: "USER_TOKEN",
				})

				Expect(err).NotTo(HaveOccurred())

				// Register cleanup
				cleanup.RegisterCleanup("token", tokenName, func() error {
					_, err := client.UserTokens.Revoke(&sonar.UserTokensRevokeOption{
						Name: tokenName,
					})
					return err
				})

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Name).To(Equal(tokenName))
				Expect(result.Type).To(Equal("USER_TOKEN"))
				Expect(result.Token).NotTo(BeEmpty())
			})
		})

		Context("with GLOBAL_ANALYSIS_TOKEN type", func() {
			It("should generate a global analysis token", func() {
				tokenName := helpers.UniqueResourceName("token-global")

				result, resp, err := client.UserTokens.Generate(&sonar.UserTokensGenerateOption{
					Name: tokenName,
					Type: "GLOBAL_ANALYSIS_TOKEN",
				})

				Expect(err).NotTo(HaveOccurred())

				// Register cleanup
				cleanup.RegisterCleanup("token", tokenName, func() error {
					_, err := client.UserTokens.Revoke(&sonar.UserTokensRevokeOption{
						Name: tokenName,
					})
					return err
				})

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Name).To(Equal(tokenName))
				Expect(result.Type).To(Equal("GLOBAL_ANALYSIS_TOKEN"))
				Expect(result.Token).NotTo(BeEmpty())
			})
		})

		Context("with PROJECT_ANALYSIS_TOKEN type", func() {
			var testProjectKey string

			BeforeEach(func() {
				// Create a project for the token
				testProjectKey = helpers.UniqueResourceName("proj-token")

				_, _, err := client.Projects.Create(&sonar.ProjectsCreateOption{
					Name:    "Token Test Project",
					Project: testProjectKey,
				})
				Expect(err).NotTo(HaveOccurred())

				cleanup.RegisterCleanup("project", testProjectKey, func() error {
					_, err := client.Projects.Delete(&sonar.ProjectsDeleteOption{
						Project: testProjectKey,
					})
					return err
				})
			})

			It("should generate a project analysis token", func() {
				tokenName := helpers.UniqueResourceName("token-proj")

				result, resp, err := client.UserTokens.Generate(&sonar.UserTokensGenerateOption{
					Name:       tokenName,
					Type:       "PROJECT_ANALYSIS_TOKEN",
					ProjectKey: testProjectKey,
				})

				Expect(err).NotTo(HaveOccurred())

				// Register cleanup
				cleanup.RegisterCleanup("token", tokenName, func() error {
					_, err := client.UserTokens.Revoke(&sonar.UserTokensRevokeOption{
						Name: tokenName,
					})
					return err
				})

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Name).To(Equal(tokenName))
				Expect(result.Type).To(Equal("PROJECT_ANALYSIS_TOKEN"))
				Expect(result.Token).NotTo(BeEmpty())
			})
		})

		Context("with expiration date", func() {
			It("should generate a token with expiration date", func() {
				tokenName := helpers.UniqueResourceName("token-exp")
				// Set expiration date to 30 days from now
				expirationDate := "2027-01-01"

				result, resp, err := client.UserTokens.Generate(&sonar.UserTokensGenerateOption{
					Name:           tokenName,
					ExpirationDate: expirationDate,
				})

				Expect(err).NotTo(HaveOccurred())

				// Register cleanup
				cleanup.RegisterCleanup("token", tokenName, func() error {
					_, err := client.UserTokens.Revoke(&sonar.UserTokensRevokeOption{
						Name: tokenName,
					})
					return err
				})

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Name).To(Equal(tokenName))
				Expect(result.ExpirationDate).NotTo(BeEmpty())
			})
		})

		Context("for another user", func() {
			var testUserLogin string

			BeforeEach(func() {
				testUserLogin = helpers.UniqueResourceName("user-token")

				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, _, err := client.Users.Create(&sonar.UsersCreateOption{
					Login:    testUserLogin,
					Name:     "Token Test User",
					Password: "SecurePassword123!",
					Local:    true,
				})
				Expect(err).NotTo(HaveOccurred())

				cleanup.RegisterCleanup("user", testUserLogin, func() error {
					//nolint:staticcheck // Using deprecated API until v2 API is implemented
					_, _, err := client.Users.Deactivate(&sonar.UsersDeactivateOption{
						Login:     testUserLogin,
						Anonymize: true,
					})
					return err
				})
			})

			It("should generate a token for another user", func() {
				tokenName := helpers.UniqueResourceName("token-other")

				result, resp, err := client.UserTokens.Generate(&sonar.UserTokensGenerateOption{
					Name:  tokenName,
					Login: testUserLogin,
				})

				Expect(err).NotTo(HaveOccurred())

				// Register cleanup
				cleanup.RegisterCleanup("token", tokenName, func() error {
					_, err := client.UserTokens.Revoke(&sonar.UserTokensRevokeOption{
						Name:  tokenName,
						Login: testUserLogin,
					})
					return err
				})

				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Name).To(Equal(tokenName))
				Expect(result.Login).To(Equal(testUserLogin))
			})
		})

		Context("duplicate token", func() {
			It("should fail when generating token with existing name", func() {
				tokenName := helpers.UniqueResourceName("token-dup")

				// Generate first token
				_, _, err := client.UserTokens.Generate(&sonar.UserTokensGenerateOption{
					Name: tokenName,
				})
				Expect(err).NotTo(HaveOccurred())

				// Register cleanup
				cleanup.RegisterCleanup("token", tokenName, func() error {
					_, err := client.UserTokens.Revoke(&sonar.UserTokensRevokeOption{
						Name: tokenName,
					})
					return err
				})

				// Try to generate duplicate
				_, resp, err := client.UserTokens.Generate(&sonar.UserTokensGenerateOption{
					Name: tokenName,
				})
				Expect(err).To(HaveOccurred())
				if resp != nil {
					Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
				}
			})
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.UserTokens.Generate(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing name", func() {
				_, resp, err := client.UserTokens.Generate(&sonar.UserTokensGenerateOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with invalid type", func() {
				_, resp, err := client.UserTokens.Generate(&sonar.UserTokensGenerateOption{
					Name: "test-token",
					Type: "INVALID_TYPE",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with PROJECT_ANALYSIS_TOKEN without project key", func() {
				_, resp, err := client.UserTokens.Generate(&sonar.UserTokensGenerateOption{
					Name: "test-token",
					Type: "PROJECT_ANALYSIS_TOKEN",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with name too long", func() {
				longName := strings.Repeat("a", sonar.MaxTokenNameLength+1)
				_, resp, err := client.UserTokens.Generate(&sonar.UserTokensGenerateOption{
					Name: longName,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("Search", func() {
		var tokenName string

		BeforeEach(func() {
			tokenName = helpers.UniqueResourceName("token-search")

			_, _, err := client.UserTokens.Generate(&sonar.UserTokensGenerateOption{
				Name: tokenName,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("token", tokenName, func() error {
				_, err := client.UserTokens.Revoke(&sonar.UserTokensRevokeOption{
					Name: tokenName,
				})
				return err
			})
		})

		It("should search tokens for current user", func() {
			result, resp, err := client.UserTokens.Search(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Login).NotTo(BeEmpty())
			Expect(result.UserTokens).NotTo(BeEmpty())

			// Find our test token
			found := false
			for _, t := range result.UserTokens {
				if t.Name == tokenName {
					found = true
					Expect(t.Type).To(Equal("USER_TOKEN"))
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		It("should search tokens with empty options", func() {
			result, resp, err := client.UserTokens.Search(&sonar.UserTokensSearchOption{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.UserTokens).NotTo(BeEmpty())
		})

		Context("for another user", func() {
			var testUserLogin string
			var userTokenName string

			BeforeEach(func() {
				testUserLogin = helpers.UniqueResourceName("user-search")
				userTokenName = helpers.UniqueResourceName("token-usrsrch")

				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, _, err := client.Users.Create(&sonar.UsersCreateOption{
					Login:    testUserLogin,
					Name:     "Search Test User",
					Password: "SecurePassword123!",
					Local:    true,
				})
				Expect(err).NotTo(HaveOccurred())

				cleanup.RegisterCleanup("user", testUserLogin, func() error {
					//nolint:staticcheck // Using deprecated API until v2 API is implemented
					_, _, err := client.Users.Deactivate(&sonar.UsersDeactivateOption{
						Login:     testUserLogin,
						Anonymize: true,
					})
					return err
				})

				// Generate a token for the test user
				_, _, err = client.UserTokens.Generate(&sonar.UserTokensGenerateOption{
					Name:  userTokenName,
					Login: testUserLogin,
				})
				Expect(err).NotTo(HaveOccurred())

				cleanup.RegisterCleanup("token", userTokenName, func() error {
					_, err := client.UserTokens.Revoke(&sonar.UserTokensRevokeOption{
						Name:  userTokenName,
						Login: testUserLogin,
					})
					return err
				})
			})

			It("should search tokens for specific user", func() {
				result, resp, err := client.UserTokens.Search(&sonar.UserTokensSearchOption{
					Login: testUserLogin,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Login).To(Equal(testUserLogin))
				Expect(result.UserTokens).NotTo(BeEmpty())

				// Find the user's token
				found := false
				for _, t := range result.UserTokens {
					if t.Name == userTokenName {
						found = true
						break
					}
				}
				Expect(found).To(BeTrue())
			})
		})

		It("should return token details", func() {
			result, resp, err := client.UserTokens.Search(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())

			// Check token details
			for _, t := range result.UserTokens {
				if t.Name == tokenName {
					Expect(t.Name).NotTo(BeEmpty())
					Expect(t.Type).NotTo(BeEmpty())
					Expect(t.CreatedAt).NotTo(BeEmpty())
					break
				}
			}
		})
	})

	Describe("Revoke", func() {
		It("should revoke a token", func() {
			tokenName := helpers.UniqueResourceName("token-revoke")

			// Generate a token first
			_, _, err := client.UserTokens.Generate(&sonar.UserTokensGenerateOption{
				Name: tokenName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Revoke the token
			resp, err := client.UserTokens.Revoke(&sonar.UserTokensRevokeOption{
				Name: tokenName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify token is revoked by searching
			result, _, err := client.UserTokens.Search(nil)
			Expect(err).NotTo(HaveOccurred())
			for _, t := range result.UserTokens {
				Expect(t.Name).NotTo(Equal(tokenName))
			}
		})

		Context("for another user", func() {
			var testUserLogin string
			var userTokenName string

			BeforeEach(func() {
				testUserLogin = helpers.UniqueResourceName("user-revoke")
				userTokenName = helpers.UniqueResourceName("token-usrrev")

				//nolint:staticcheck // Using deprecated API until v2 API is implemented
				_, _, err := client.Users.Create(&sonar.UsersCreateOption{
					Login:    testUserLogin,
					Name:     "Revoke Test User",
					Password: "SecurePassword123!",
					Local:    true,
				})
				Expect(err).NotTo(HaveOccurred())

				cleanup.RegisterCleanup("user", testUserLogin, func() error {
					//nolint:staticcheck // Using deprecated API until v2 API is implemented
					_, _, err := client.Users.Deactivate(&sonar.UsersDeactivateOption{
						Login:     testUserLogin,
						Anonymize: true,
					})
					return err
				})

				// Generate a token for the test user
				_, _, err = client.UserTokens.Generate(&sonar.UserTokensGenerateOption{
					Name:  userTokenName,
					Login: testUserLogin,
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("should revoke a token for another user", func() {
				resp, err := client.UserTokens.Revoke(&sonar.UserTokensRevokeOption{
					Name:  userTokenName,
					Login: testUserLogin,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// Verify token is revoked
				result, _, err := client.UserTokens.Search(&sonar.UserTokensSearchOption{
					Login: testUserLogin,
				})
				Expect(err).NotTo(HaveOccurred())
				for _, t := range result.UserTokens {
					Expect(t.Name).NotTo(Equal(userTokenName))
				}
			})
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.UserTokens.Revoke(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing name", func() {
				resp, err := client.UserTokens.Revoke(&sonar.UserTokensRevokeOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should handle revoking non-existent token gracefully", func() {
				// SonarQube API is idempotent - revoking a non-existent token succeeds silently
				resp, err := client.UserTokens.Revoke(&sonar.UserTokensRevokeOption{
					Name: "nonexistent-token-xyz12345",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})
		})
	})

	Describe("Token Lifecycle", func() {
		It("should complete full generate/search/revoke cycle", func() {
			tokenName := helpers.UniqueResourceName("token-lifecycle")

			// Step 1: Generate token
			generateResult, _, err := client.UserTokens.Generate(&sonar.UserTokensGenerateOption{
				Name: tokenName,
				Type: "USER_TOKEN",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(generateResult.Name).To(Equal(tokenName))
			Expect(generateResult.Token).NotTo(BeEmpty())
			generatedToken := generateResult.Token

			// Step 2: Search and verify token exists
			searchResult, _, err := client.UserTokens.Search(nil)
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, t := range searchResult.UserTokens {
				if t.Name == tokenName {
					found = true
					Expect(t.Type).To(Equal("USER_TOKEN"))
					// Token value should NOT be returned in search
					break
				}
			}
			Expect(found).To(BeTrue())

			// Step 3: Verify token can be used for authentication
			cfg := helpers.LoadConfig()
			tokenClient, err := sonar.NewClient(nil,
				sonar.WithBaseURL(helpers.NormalizeBaseURL(cfg.BaseURL)),
				sonar.WithToken(generatedToken),
			)
			Expect(err).NotTo(HaveOccurred())

			// Use the token to make an authenticated request
			authResult, _, err := tokenClient.Authentication.Validate()
			Expect(err).NotTo(HaveOccurred())
			Expect(authResult.Valid).To(BeTrue())

			// Step 4: Revoke token
			_, err = client.UserTokens.Revoke(&sonar.UserTokensRevokeOption{
				Name: tokenName,
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 5: Verify token no longer exists
			searchResult, _, err = client.UserTokens.Search(nil)
			Expect(err).NotTo(HaveOccurred())
			for _, t := range searchResult.UserTokens {
				Expect(t.Name).NotTo(Equal(tokenName))
			}
		})
	})

	Describe("Multiple Tokens", func() {
		It("should manage multiple tokens for the same user", func() {
			token1Name := helpers.UniqueResourceName("token-multi1")
			token2Name := helpers.UniqueResourceName("token-multi2")
			token3Name := helpers.UniqueResourceName("token-multi3")

			// Generate multiple tokens
			_, _, err := client.UserTokens.Generate(&sonar.UserTokensGenerateOption{
				Name: token1Name,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("token", token1Name, func() error {
				_, err := client.UserTokens.Revoke(&sonar.UserTokensRevokeOption{
					Name: token1Name,
				})
				return err
			})

			_, _, err = client.UserTokens.Generate(&sonar.UserTokensGenerateOption{
				Name: token2Name,
				Type: "GLOBAL_ANALYSIS_TOKEN",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("token", token2Name, func() error {
				_, err := client.UserTokens.Revoke(&sonar.UserTokensRevokeOption{
					Name: token2Name,
				})
				return err
			})

			_, _, err = client.UserTokens.Generate(&sonar.UserTokensGenerateOption{
				Name: token3Name,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("token", token3Name, func() error {
				_, err := client.UserTokens.Revoke(&sonar.UserTokensRevokeOption{
					Name: token3Name,
				})
				return err
			})

			// Search and verify all tokens exist
			result, _, err := client.UserTokens.Search(nil)
			Expect(err).NotTo(HaveOccurred())

			foundTokens := make(map[string]bool)
			for _, t := range result.UserTokens {
				if t.Name == token1Name || t.Name == token2Name || t.Name == token3Name {
					foundTokens[t.Name] = true
				}
			}
			Expect(foundTokens).To(HaveKey(token1Name))
			Expect(foundTokens).To(HaveKey(token2Name))
			Expect(foundTokens).To(HaveKey(token3Name))
		})
	})
})
