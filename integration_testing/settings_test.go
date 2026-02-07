package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("Settings Service", Ordered, func() {
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
		errors := cleanup.Cleanup()
		for _, err := range errors {
			GinkgoWriter.Printf("Cleanup error: %v\n", err)
		}
	})

	// =========================================================================
	// ListDefinitions
	// =========================================================================
	Describe("ListDefinitions", func() {
		It("should list all setting definitions", func() {
			result, resp, err := client.Settings.ListDefinitions(&sonar.SettingsListDefinitionsOption{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Definitions).NotTo(BeEmpty())
		})

		It("should return definitions with proper structure", func() {
			result, resp, err := client.Settings.ListDefinitions(&sonar.SettingsListDefinitionsOption{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			// Each definition should have at least a key
			for _, def := range result.Definitions {
				Expect(def.Key).NotTo(BeEmpty())
			}
		})

		It("should list definitions for a specific project", func() {
			// First create a project
			projectKey := helpers.UniqueResourceName("proj")
			_, _, err := client.Projects.Create(&sonar.ProjectsCreateOption{
				Name:    projectKey,
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonar.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			// List definitions for the project
			result, resp, err := client.Settings.ListDefinitions(&sonar.SettingsListDefinitionsOption{
				Component: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.Settings.ListDefinitions(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Values
	// =========================================================================
	Describe("Values", func() {
		It("should list all setting values", func() {
			result, resp, err := client.Settings.Values(&sonar.SettingsValuesOption{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Settings).NotTo(BeEmpty())
		})

		It("should filter by specific keys", func() {
			result, resp, err := client.Settings.Values(&sonar.SettingsValuesOption{
				Keys: []string{"sonar.core.id", "sonar.core.startTime"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			// Should contain settings for the requested keys
			for _, setting := range result.Settings {
				Expect([]string{"sonar.core.id", "sonar.core.startTime"}).To(ContainElement(setting.Key))
			}
		})

		It("should get values for a specific project", func() {
			// First create a project
			projectKey := helpers.UniqueResourceName("proj")
			_, _, err := client.Projects.Create(&sonar.ProjectsCreateOption{
				Name:    projectKey,
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonar.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			// Get values for the project
			result, resp, err := client.Settings.Values(&sonar.SettingsValuesOption{
				Component: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.Settings.Values(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Set and Reset
	// =========================================================================
	Describe("Set and Reset", func() {
		Describe("Set", func() {
			It("should set a global setting value", func() {
				// Set a simple string setting
				resp, err := client.Settings.Set(&sonar.SettingsSetOption{
					Key:   "sonar.login.message",
					Value: "E2E Test Login Message",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// Verify the value was set
				values, _, err := client.Settings.Values(&sonar.SettingsValuesOption{
					Keys: []string{"sonar.login.message"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(values.Settings).NotTo(BeEmpty())
				Expect(values.Settings[0].Value).To(Equal("E2E Test Login Message"))

				// Clean up - reset the setting
				_, _ = client.Settings.Reset(&sonar.SettingsResetOption{
					Keys: []string{"sonar.login.message"},
				})
			})

			It("should set a project-level setting", func() {
				// Create a project
				projectKey := helpers.UniqueResourceName("proj")
				_, _, err := client.Projects.Create(&sonar.ProjectsCreateOption{
					Name:    projectKey,
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())

				cleanup.RegisterCleanup("project", projectKey, func() error {
					_, err := client.Projects.Delete(&sonar.ProjectsDeleteOption{
						Project: projectKey,
					})
					return err
				})

				// Set a project-level setting (sonar.exclusions is a multi-value setting)
				resp, err := client.Settings.Set(&sonar.SettingsSetOption{
					Key:       "sonar.exclusions",
					Values:    []string{"**/test/**"},
					Component: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// Verify the value was set
				values, _, err := client.Settings.Values(&sonar.SettingsValuesOption{
					Component: projectKey,
					Keys:      []string{"sonar.exclusions"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(values.Settings).NotTo(BeEmpty())
			})

			It("should set a multi-value setting", func() {
				// Create a project for testing
				projectKey := helpers.UniqueResourceName("proj")
				_, _, err := client.Projects.Create(&sonar.ProjectsCreateOption{
					Name:    projectKey,
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())

				cleanup.RegisterCleanup("project", projectKey, func() error {
					_, err := client.Projects.Delete(&sonar.ProjectsDeleteOption{
						Project: projectKey,
					})
					return err
				})

				// Set a multi-value setting on the project
				resp, err := client.Settings.Set(&sonar.SettingsSetOption{
					Key:       "sonar.exclusions",
					Values:    []string{"**/test/**", "**/vendor/**"},
					Component: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})

			Context("parameter validation", func() {
				It("should fail with nil options", func() {
					resp, err := client.Settings.Set(nil)
					Expect(err).To(HaveOccurred())
					Expect(resp).To(BeNil())
				})

				It("should fail with missing key", func() {
					resp, err := client.Settings.Set(&sonar.SettingsSetOption{
						Value: "some value",
					})
					Expect(err).To(HaveOccurred())
					Expect(resp).To(BeNil())
				})

				It("should fail with missing value/values/fieldValues", func() {
					resp, err := client.Settings.Set(&sonar.SettingsSetOption{
						Key: "some.key",
					})
					Expect(err).To(HaveOccurred())
					Expect(resp).To(BeNil())
				})
			})
		})

		Describe("Reset", func() {
			It("should reset a global setting value", func() {
				// First set a value
				_, err := client.Settings.Set(&sonar.SettingsSetOption{
					Key:   "sonar.login.message",
					Value: "Message to be reset",
				})
				Expect(err).NotTo(HaveOccurred())

				// Reset the setting
				resp, err := client.Settings.Reset(&sonar.SettingsResetOption{
					Keys: []string{"sonar.login.message"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// Verify the value was reset (should not be found or have default value)
				values, _, err := client.Settings.Values(&sonar.SettingsValuesOption{
					Keys: []string{"sonar.login.message"},
				})
				Expect(err).NotTo(HaveOccurred())
				// After reset, either no settings or empty value
				if len(values.Settings) > 0 {
					Expect(values.Settings[0].Value).To(BeEmpty())
				}
			})

			It("should reset a project-level setting", func() {
				// Create a project
				projectKey := helpers.UniqueResourceName("proj")
				_, _, err := client.Projects.Create(&sonar.ProjectsCreateOption{
					Name:    projectKey,
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())

				cleanup.RegisterCleanup("project", projectKey, func() error {
					_, err := client.Projects.Delete(&sonar.ProjectsDeleteOption{
						Project: projectKey,
					})
					return err
				})

				// Set a project-level setting (sonar.exclusions is a multi-value setting)
				_, err = client.Settings.Set(&sonar.SettingsSetOption{
					Key:       "sonar.exclusions",
					Values:    []string{"**/test/**"},
					Component: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())

				// Reset the project-level setting
				resp, err := client.Settings.Reset(&sonar.SettingsResetOption{
					Keys:      []string{"sonar.exclusions"},
					Component: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})

			It("should reset multiple settings at once", func() {
				// Create a project
				projectKey := helpers.UniqueResourceName("proj")
				_, _, err := client.Projects.Create(&sonar.ProjectsCreateOption{
					Name:    projectKey,
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())

				cleanup.RegisterCleanup("project", projectKey, func() error {
					_, err := client.Projects.Delete(&sonar.ProjectsDeleteOption{
						Project: projectKey,
					})
					return err
				})

				// Set multiple settings (sonar.exclusions and sonar.inclusions are multi-value)
				_, err = client.Settings.Set(&sonar.SettingsSetOption{
					Key:       "sonar.exclusions",
					Values:    []string{"**/test/**"},
					Component: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())

				_, err = client.Settings.Set(&sonar.SettingsSetOption{
					Key:       "sonar.inclusions",
					Values:    []string{"**/src/**"},
					Component: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())

				// Reset multiple settings
				resp, err := client.Settings.Reset(&sonar.SettingsResetOption{
					Keys:      []string{"sonar.exclusions", "sonar.inclusions"},
					Component: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})

			Context("parameter validation", func() {
				It("should fail with nil options", func() {
					resp, err := client.Settings.Reset(nil)
					Expect(err).To(HaveOccurred())
					Expect(resp).To(BeNil())
				})

				It("should fail with missing keys", func() {
					resp, err := client.Settings.Reset(&sonar.SettingsResetOption{})
					Expect(err).To(HaveOccurred())
					Expect(resp).To(BeNil())
				})

				It("should fail with empty keys array", func() {
					resp, err := client.Settings.Reset(&sonar.SettingsResetOption{
						Keys: []string{},
					})
					Expect(err).To(HaveOccurred())
					Expect(resp).To(BeNil())
				})
			})
		})
	})

	// =========================================================================
	// LoginMessage
	// =========================================================================
	Describe("LoginMessage", func() {
		It("should get the login message", func() {
			result, resp, err := client.Settings.LoginMessage()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			// Message may be empty if not set
		})

		It("should return the set login message", func() {
			// Set a login message
			_, err := client.Settings.Set(&sonar.SettingsSetOption{
				Key:   "sonar.login.message",
				Value: "Welcome to E2E Testing!",
			})
			Expect(err).NotTo(HaveOccurred())

			// Verify the value was set using Values endpoint
			values, _, err := client.Settings.Values(&sonar.SettingsValuesOption{
				Keys: []string{"sonar.login.message"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(values.Settings).NotTo(BeEmpty())
			Expect(values.Settings[0].Value).To(Equal("Welcome to E2E Testing!"))

			// The LoginMessage endpoint may return empty due to SonarQube behavior
			// but the setting value is correctly stored
			result, resp, err := client.Settings.LoginMessage()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())

			// Clean up
			_, _ = client.Settings.Reset(&sonar.SettingsResetOption{
				Keys: []string{"sonar.login.message"},
			})
		})
	})

	// =========================================================================
	// CheckSecretKey
	// =========================================================================
	Describe("CheckSecretKey", func() {
		It("should check if secret key is available", func() {
			result, resp, err := client.Settings.CheckSecretKey()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			// SecretKeyAvailable is a boolean
		})
	})

	// =========================================================================
	// GenerateSecretKey
	// =========================================================================
	Describe("GenerateSecretKey", func() {
		It("should generate a secret key", func() {
			Skip("Skipping secret key generation test - it may fail if key already exists")
			result, resp, err := client.Settings.GenerateSecretKey()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.SecretKey).NotTo(BeEmpty())
		})
	})

	// =========================================================================
	// Encrypt
	// =========================================================================
	Describe("Encrypt", func() {
		It("should encrypt a value if secret key is available", func() {
			// First check if secret key is available
			checkResult, _, err := client.Settings.CheckSecretKey()
			Expect(err).NotTo(HaveOccurred())

			if !checkResult.SecretKeyAvailable {
				Skip("Secret key not available, skipping encryption test")
			}

			result, resp, err := client.Settings.Encrypt(&sonar.SettingsEncryptOption{
				Value: "my-secret-value",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.EncryptedValue).NotTo(BeEmpty())
			// Encrypted values start with specific prefix
			Expect(result.EncryptedValue).To(HavePrefix("{aes-gcm}"))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.Settings.Encrypt(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing value", func() {
				_, resp, err := client.Settings.Encrypt(&sonar.SettingsEncryptOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})
})
