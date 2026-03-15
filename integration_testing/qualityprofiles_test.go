package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("Qualityprofiles Service", Ordered, func() {
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
	// Search
	// =========================================================================
	Describe("Search", func() {
		It("should list all quality profiles", func() {
			result, resp, err := client.Qualityprofiles.Search(&sonar.QualityprofilesSearchOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Profiles).NotTo(BeEmpty())
		})

		It("should filter by language", func() {
			result, resp, err := client.Qualityprofiles.Search(&sonar.QualityprofilesSearchOptions{
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			for _, profile := range result.Profiles {
				Expect(profile.Language).To(Equal("java"))
			}
		})

		It("should include defaults", func() {
			result, resp, err := client.Qualityprofiles.Search(&sonar.QualityprofilesSearchOptions{
				Defaults: true,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			for _, profile := range result.Profiles {
				Expect(profile.IsDefault).To(BeTrue())
			}
		})
	})

	// =========================================================================
	// Create
	// =========================================================================
	Describe("Create", func() {
		It("should create a new quality profile", func() {
			profileName := helpers.UniqueResourceName("qp")

			result, resp, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Profile.Name).To(Equal(profileName))

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName,
					Language:       "java",
				})
				return err
			})
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Qualityprofiles.Create(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing name", func() {
				result, resp, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
					Language: "java",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing language", func() {
				result, resp, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
					Name: "test-profile",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Show
	// =========================================================================
	Describe("Show", func() {
		It("should show quality profile details", func() {
			profileName := helpers.UniqueResourceName("qp-show")

			createResult, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName,
					Language:       "java",
				})
				return err
			})

			result, resp, err := client.Qualityprofiles.Show(&sonar.QualityprofilesShowOptions{
				Key: createResult.Profile.Key,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Profile.Name).To(Equal(profileName))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Qualityprofiles.Show(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with empty key", func() {
				result, resp, err := client.Qualityprofiles.Show(&sonar.QualityprofilesShowOptions{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Rename
	// =========================================================================
	Describe("Rename", func() {
		It("should rename a quality profile", func() {
			originalName := helpers.UniqueResourceName("qp-rename")
			newName := helpers.UniqueResourceName("qp-renamed")

			createResult, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     originalName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			// Register cleanup that tries both names in case rename fails
			cleanup.RegisterCleanup("qualityprofile", originalName+"-or-"+newName, func() error {
				// Try deleting by new name first (expected case)
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: newName,
					Language:       "java",
				})
				if err == nil {
					return nil
				}
				// If that fails, try deleting by original name (rename failed)
				_, err = client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: originalName,
					Language:       "java",
				})
				return err
			})

			resp, err := client.Qualityprofiles.Rename(&sonar.QualityprofilesRenameOptions{
				Key:  createResult.Profile.Key,
				Name: newName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify rename
			result, _, err := client.Qualityprofiles.Show(&sonar.QualityprofilesShowOptions{
				Key: createResult.Profile.Key,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Profile.Name).To(Equal(newName))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Qualityprofiles.Rename(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing key", func() {
				resp, err := client.Qualityprofiles.Rename(&sonar.QualityprofilesRenameOptions{
					Name: "new-name",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing name", func() {
				resp, err := client.Qualityprofiles.Rename(&sonar.QualityprofilesRenameOptions{
					Key: "some-key",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Copy
	// =========================================================================
	Describe("Copy", func() {
		It("should copy a quality profile", func() {
			sourceName := helpers.UniqueResourceName("qp-source")
			copyName := helpers.UniqueResourceName("qp-copy")

			createResult, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     sourceName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", sourceName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: sourceName,
					Language:       "java",
				})
				return err
			})

			result, resp, err := client.Qualityprofiles.Copy(&sonar.QualityprofilesCopyOptions{
				FromKey: createResult.Profile.Key,
				ToName:  copyName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Name).To(Equal(copyName))

			cleanup.RegisterCleanup("qualityprofile", copyName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: copyName,
					Language:       "java",
				})
				return err
			})
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Qualityprofiles.Copy(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing from key", func() {
				result, resp, err := client.Qualityprofiles.Copy(&sonar.QualityprofilesCopyOptions{
					ToName: "new-copy",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing to name", func() {
				result, resp, err := client.Qualityprofiles.Copy(&sonar.QualityprofilesCopyOptions{
					FromKey: "some-key",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Delete
	// =========================================================================
	Describe("Delete", func() {
		It("should delete a quality profile", func() {
			profileName := helpers.UniqueResourceName("qp-delete")

			_, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			resp, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
				QualityProfile: profileName,
				Language:       "java",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Qualityprofiles.Delete(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing profile name", func() {
				resp, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					Language: "java",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing language", func() {
				resp, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: "test",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// SetDefault
	// =========================================================================
	Describe("SetDefault", func() {
		It("should set a quality profile as default", func() {
			profileName := helpers.UniqueResourceName("qp-default")

			// Capture the current default profile for this language FIRST
			searchResult, _, err := client.Qualityprofiles.Search(&sonar.QualityprofilesSearchOptions{
				Language: "java",
				Defaults: true,
			})
			Expect(err).NotTo(HaveOccurred())
			var originalDefaultName string
			for _, p := range searchResult.Profiles {
				if p.IsDefault && p.Language == "java" {
					originalDefaultName = p.Name
					break
				}
			}
			Expect(originalDefaultName).NotTo(BeEmpty(), "Should find original default Java profile")

			createResult, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			// Use DeferCleanup to ensure restoration happens even if test fails
			DeferCleanup(func() {
				// Restore the original default profile
				_, restoreErr := client.Qualityprofiles.SetDefault(&sonar.QualityprofilesSetDefaultOptions{
					QualityProfile: originalDefaultName,
					Language:       "java",
				})
				if restoreErr != nil {
					GinkgoWriter.Printf("Failed to restore default profile %s: %v\n", originalDefaultName, restoreErr)
				}
				// Delete the test profile
				_, _ = client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName,
					Language:       "java",
				})
			})

			resp, err := client.Qualityprofiles.SetDefault(&sonar.QualityprofilesSetDefaultOptions{
				QualityProfile: profileName,
				Language:       "java",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify it's now default
			result, _, err := client.Qualityprofiles.Show(&sonar.QualityprofilesShowOptions{
				Key: createResult.Profile.Key,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Profile.IsDefault).To(BeTrue())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Qualityprofiles.SetDefault(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing profile name", func() {
				resp, err := client.Qualityprofiles.SetDefault(&sonar.QualityprofilesSetDefaultOptions{
					Language: "java",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing language", func() {
				resp, err := client.Qualityprofiles.SetDefault(&sonar.QualityprofilesSetDefaultOptions{
					QualityProfile: "test",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// AddProject and RemoveProject
	// =========================================================================
	Describe("AddProject", func() {
		It("should associate a project with a quality profile", func() {
			profileName := helpers.UniqueResourceName("qp-proj")
			projectKey := helpers.UniqueResourceName("proj-qp")

			_, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName,
					Language:       "java",
				})
				return err
			})

			_, _, err = client.Projects.Create(&sonar.ProjectsCreateOptions{
				Name:    "QualityProfile AddProject Test",
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonar.ProjectsDeleteOptions{
					Project: projectKey,
				})
				return err
			})

			resp, err := client.Qualityprofiles.AddProject(&sonar.QualityprofilesAddProjectOptions{
				QualityProfile: profileName,
				Language:       "java",
				Project:        projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Qualityprofiles.AddProject(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing profile name", func() {
				resp, err := client.Qualityprofiles.AddProject(&sonar.QualityprofilesAddProjectOptions{
					Language: "java",
					Project:  "some-project",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing project", func() {
				resp, err := client.Qualityprofiles.AddProject(&sonar.QualityprofilesAddProjectOptions{
					QualityProfile: "test",
					Language:       "java",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("RemoveProject", func() {
		It("should remove project association from quality profile", func() {
			profileName := helpers.UniqueResourceName("qp-rmproj")
			projectKey := helpers.UniqueResourceName("proj-rmqp")

			_, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName,
					Language:       "java",
				})
				return err
			})

			_, _, err = client.Projects.Create(&sonar.ProjectsCreateOptions{
				Name:    "QualityProfile RemoveProject Test",
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonar.ProjectsDeleteOptions{
					Project: projectKey,
				})
				return err
			})

			// Add first
			_, err = client.Qualityprofiles.AddProject(&sonar.QualityprofilesAddProjectOptions{
				QualityProfile: profileName,
				Language:       "java",
				Project:        projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			// Remove
			resp, err := client.Qualityprofiles.RemoveProject(&sonar.QualityprofilesRemoveProjectOptions{
				QualityProfile: profileName,
				Language:       "java",
				Project:        projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Qualityprofiles.RemoveProject(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing profile name", func() {
				resp, err := client.Qualityprofiles.RemoveProject(&sonar.QualityprofilesRemoveProjectOptions{
					Language: "java",
					Project:  "some-project",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Projects
	// =========================================================================
	Describe("Projects", func() {
		It("should list projects for a quality profile", func() {
			profileName := helpers.UniqueResourceName("qp-projects")

			createResult, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName,
					Language:       "java",
				})
				return err
			})

			result, resp, err := client.Qualityprofiles.Projects(&sonar.QualityprofilesProjectsOptions{
				Key: createResult.Profile.Key,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Qualityprofiles.Projects(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing key", func() {
				result, resp, err := client.Qualityprofiles.Projects(&sonar.QualityprofilesProjectsOptions{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Inheritance
	// =========================================================================
	Describe("Inheritance", func() {
		It("should show profile inheritance", func() {
			profileName := helpers.UniqueResourceName("qp-inh")

			_, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName,
					Language:       "java",
				})
				return err
			})

			result, resp, err := client.Qualityprofiles.Inheritance(&sonar.QualityprofilesInheritanceOptions{
				QualityProfile: profileName,
				Language:       "java",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Profile.Name).To(Equal(profileName))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Qualityprofiles.Inheritance(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing profile name", func() {
				result, resp, err := client.Qualityprofiles.Inheritance(&sonar.QualityprofilesInheritanceOptions{
					Language: "java",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing language", func() {
				result, resp, err := client.Qualityprofiles.Inheritance(&sonar.QualityprofilesInheritanceOptions{
					QualityProfile: "test",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// ChangeParent
	// =========================================================================
	Describe("ChangeParent", func() {
		It("should set parent profile", func() {
			parentName := helpers.UniqueResourceName("qp-parent")
			childName := helpers.UniqueResourceName("qp-child")

			_, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     parentName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", parentName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: parentName,
					Language:       "java",
				})
				return err
			})

			_, _, err = client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     childName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", childName, func() error {
				// Remove parent first
				_, _ = client.Qualityprofiles.ChangeParent(&sonar.QualityprofilesChangeParentOptions{
					QualityProfile: childName,
					Language:       "java",
				})
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: childName,
					Language:       "java",
				})
				return err
			})

			resp, err := client.Qualityprofiles.ChangeParent(&sonar.QualityprofilesChangeParentOptions{
				QualityProfile:       childName,
				Language:             "java",
				ParentQualityProfile: parentName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify parent was set by checking ancestors
			result, _, err := client.Qualityprofiles.Inheritance(&sonar.QualityprofilesInheritanceOptions{
				QualityProfile: childName,
				Language:       "java",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Ancestors).To(HaveLen(1))
			Expect(result.Ancestors[0].Name).To(Equal(parentName))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Qualityprofiles.ChangeParent(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing profile name", func() {
				resp, err := client.Qualityprofiles.ChangeParent(&sonar.QualityprofilesChangeParentOptions{
					Language:             "java",
					ParentQualityProfile: "some-parent",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing language", func() {
				resp, err := client.Qualityprofiles.ChangeParent(&sonar.QualityprofilesChangeParentOptions{
					QualityProfile: "test",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Compare
	// =========================================================================
	Describe("Compare", func() {
		It("should compare two quality profiles", func() {
			profileName1 := helpers.UniqueResourceName("qp-cmp1")
			profileName2 := helpers.UniqueResourceName("qp-cmp2")

			result1, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName1,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", profileName1, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName1,
					Language:       "java",
				})
				return err
			})

			result2, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName2,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", profileName2, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName2,
					Language:       "java",
				})
				return err
			})

			compareResult, resp, err := client.Qualityprofiles.Compare(&sonar.QualityprofilesCompareOptions{
				LeftKey:  result1.Profile.Key,
				RightKey: result2.Profile.Key,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(compareResult).NotTo(BeNil())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Qualityprofiles.Compare(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing left key", func() {
				result, resp, err := client.Qualityprofiles.Compare(&sonar.QualityprofilesCompareOptions{
					RightKey: "some-key",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing right key", func() {
				result, resp, err := client.Qualityprofiles.Compare(&sonar.QualityprofilesCompareOptions{
					LeftKey: "some-key",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Changelog
	// =========================================================================
	Describe("Changelog", func() {
		It("should show quality profile changelog", func() {
			profileName := helpers.UniqueResourceName("qp-chlog")

			_, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName,
					Language:       "java",
				})
				return err
			})

			result, resp, err := client.Qualityprofiles.Changelog(&sonar.QualityprofilesChangelogOptions{
				QualityProfile: profileName,
				Language:       "java",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Qualityprofiles.Changelog(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing profile name", func() {
				result, resp, err := client.Qualityprofiles.Changelog(&sonar.QualityprofilesChangelogOptions{
					Language: "java",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing language", func() {
				result, resp, err := client.Qualityprofiles.Changelog(&sonar.QualityprofilesChangelogOptions{
					QualityProfile: "test",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Backup
	// =========================================================================
	Describe("Backup", func() {
		It("should backup a quality profile", func() {
			profileName := helpers.UniqueResourceName("qp-backup")

			_, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName,
					Language:       "java",
				})
				return err
			})

			result, resp, err := client.Qualityprofiles.Backup(&sonar.QualityprofilesBackupOptions{
				QualityProfile: profileName,
				Language:       "java",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(*result).NotTo(BeEmpty())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Qualityprofiles.Backup(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing profile name", func() {
				result, resp, err := client.Qualityprofiles.Backup(&sonar.QualityprofilesBackupOptions{
					Language: "java",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing language", func() {
				result, resp, err := client.Qualityprofiles.Backup(&sonar.QualityprofilesBackupOptions{
					QualityProfile: "test",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Restore
	// =========================================================================
	Describe("Restore", func() {
		// Note: The Restore endpoint requires multipart file upload which is not
		// currently supported by the SDK. The endpoint expects a file field named
		// "backup" containing the XML backup content as a file upload.
		// This is a known SDK limitation that should be addressed in a future version.

		It("should restore a quality profile from backup", func() {
			Skip("Restore endpoint requires multipart file upload, not currently supported by SDK")
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Qualityprofiles.Restore(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with empty backup", func() {
				resp, err := client.Qualityprofiles.Restore(&sonar.QualityprofilesRestoreOptions{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// ActivateRule and DeactivateRule
	// =========================================================================
	Describe("ActivateRule", func() {
		It("should activate a rule on a quality profile", func() {
			profileName := helpers.UniqueResourceName("qp-activate")

			// Create a profile
			createResult, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())
			profileKey := createResult.Profile.Key

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName,
					Language:       "java",
				})
				return err
			})

			// Find an available Java rule to activate
			rulesResult, _, err := client.Rules.Search(&sonar.RulesSearchOptions{
				Languages: []string{"java"},
				PaginationArgs: sonar.PaginationArgs{
					PageSize: 1,
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(rulesResult.Rules).NotTo(BeEmpty())
			ruleKey := rulesResult.Rules[0].Key

			// Activate the rule
			resp, err := client.Qualityprofiles.ActivateRule(&sonar.QualityprofilesActivateRuleOptions{
				Key:  profileKey,
				Rule: ruleKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Qualityprofiles.ActivateRule(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing key", func() {
				resp, err := client.Qualityprofiles.ActivateRule(&sonar.QualityprofilesActivateRuleOptions{
					Rule: "java:S1234",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing rule", func() {
				resp, err := client.Qualityprofiles.ActivateRule(&sonar.QualityprofilesActivateRuleOptions{
					Key: "some-key",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("DeactivateRule", func() {
		It("should deactivate a rule from a quality profile", func() {
			profileName := helpers.UniqueResourceName("qp-deactivate")

			// Create a profile
			createResult, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())
			profileKey := createResult.Profile.Key

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName,
					Language:       "java",
				})
				return err
			})

			// Find an available Java rule to activate then deactivate
			rulesResult, _, err := client.Rules.Search(&sonar.RulesSearchOptions{
				Languages: []string{"java"},
				PaginationArgs: sonar.PaginationArgs{
					PageSize: 1,
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(rulesResult.Rules).NotTo(BeEmpty())
			ruleKey := rulesResult.Rules[0].Key

			// First activate the rule
			_, err = client.Qualityprofiles.ActivateRule(&sonar.QualityprofilesActivateRuleOptions{
				Key:  profileKey,
				Rule: ruleKey,
			})
			Expect(err).NotTo(HaveOccurred())

			// Now deactivate it
			resp, err := client.Qualityprofiles.DeactivateRule(&sonar.QualityprofilesDeactivateRuleOptions{
				Key:  profileKey,
				Rule: ruleKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Qualityprofiles.DeactivateRule(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing key", func() {
				resp, err := client.Qualityprofiles.DeactivateRule(&sonar.QualityprofilesDeactivateRuleOptions{
					Rule: "java:S1234",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing rule", func() {
				resp, err := client.Qualityprofiles.DeactivateRule(&sonar.QualityprofilesDeactivateRuleOptions{
					Key: "some-key",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// ActivateRules and DeactivateRules (Bulk operations)
	// =========================================================================
	Describe("ActivateRules", func() {
		It("should bulk activate rules on a quality profile", func() {
			profileName := helpers.UniqueResourceName("qp-bulk-activate")

			// Create a profile
			createResult, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())
			profileKey := createResult.Profile.Key

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName,
					Language:       "java",
				})
				return err
			})

			// Bulk activate rules by language
			resp, err := client.Qualityprofiles.ActivateRules(&sonar.QualityprofilesActivateRulesOptions{
				TargetKey: profileKey,
				Languages: []string{"java"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Qualityprofiles.ActivateRules(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing target key", func() {
				resp, err := client.Qualityprofiles.ActivateRules(&sonar.QualityprofilesActivateRulesOptions{
					Languages: []string{"java"},
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("DeactivateRules", func() {
		It("should bulk deactivate rules from a quality profile", func() {
			profileName := helpers.UniqueResourceName("qp-bulk-deactivate")

			// Create a profile
			createResult, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())
			profileKey := createResult.Profile.Key

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName,
					Language:       "java",
				})
				return err
			})

			// Find an available Java rule and activate it first
			rulesResult, _, err := client.Rules.Search(&sonar.RulesSearchOptions{
				Languages: []string{"java"},
				PaginationArgs: sonar.PaginationArgs{
					PageSize: 1,
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(rulesResult.Rules).NotTo(BeEmpty())
			ruleKey := rulesResult.Rules[0].Key

			// Activate a rule first
			_, err = client.Qualityprofiles.ActivateRule(&sonar.QualityprofilesActivateRuleOptions{
				Key:  profileKey,
				Rule: ruleKey,
			})
			Expect(err).NotTo(HaveOccurred())

			// Bulk deactivate rules
			resp, err := client.Qualityprofiles.DeactivateRules(&sonar.QualityprofilesDeactivateRulesOptions{
				TargetKey: profileKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Qualityprofiles.DeactivateRules(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing target key", func() {
				resp, err := client.Qualityprofiles.DeactivateRules(&sonar.QualityprofilesDeactivateRulesOptions{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Exporters and Importers
	// =========================================================================
	Describe("Exporters", func() {
		It("should list exporters", func() {
			result, resp, err := client.Qualityprofiles.Exporters()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})
	})

	Describe("Importers", func() {
		It("should list importers", func() {
			result, resp, err := client.Qualityprofiles.Importers()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})
	})

	// =========================================================================
	// AddGroup, SearchGroups, RemoveGroup
	// =========================================================================
	Describe("AddGroup", func() {
		It("should add group permission to quality profile", func() {
			profileName := helpers.UniqueResourceName("qp-grp")

			_, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName,
					Language:       "java",
				})
				return err
			})

			resp, err := client.Qualityprofiles.AddGroup(&sonar.QualityprofilesAddGroupOptions{
				QualityProfile: profileName,
				Language:       "java",
				Group:          "sonar-users",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Qualityprofiles.AddGroup(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing profile name", func() {
				resp, err := client.Qualityprofiles.AddGroup(&sonar.QualityprofilesAddGroupOptions{
					Language: "java",
					Group:    "sonar-users",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing group", func() {
				resp, err := client.Qualityprofiles.AddGroup(&sonar.QualityprofilesAddGroupOptions{
					QualityProfile: "test",
					Language:       "java",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("SearchGroups", func() {
		It("should search groups for a quality profile", func() {
			profileName := helpers.UniqueResourceName("qp-sgrp")

			_, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName,
					Language:       "java",
				})
				return err
			})

			result, resp, err := client.Qualityprofiles.SearchGroups(&sonar.QualityprofilesSearchGroupsOptions{
				QualityProfile: profileName,
				Language:       "java",
				Selected:       "all",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Qualityprofiles.SearchGroups(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing profile name", func() {
				result, resp, err := client.Qualityprofiles.SearchGroups(&sonar.QualityprofilesSearchGroupsOptions{
					Language: "java",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing language", func() {
				result, resp, err := client.Qualityprofiles.SearchGroups(&sonar.QualityprofilesSearchGroupsOptions{
					QualityProfile: "test",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("RemoveGroup", func() {
		It("should remove group permission from quality profile", func() {
			profileName := helpers.UniqueResourceName("qp-rmgrp")

			_, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName,
					Language:       "java",
				})
				return err
			})

			// Add group first
			_, err = client.Qualityprofiles.AddGroup(&sonar.QualityprofilesAddGroupOptions{
				QualityProfile: profileName,
				Language:       "java",
				Group:          "sonar-users",
			})
			Expect(err).NotTo(HaveOccurred())

			// Remove group
			resp, err := client.Qualityprofiles.RemoveGroup(&sonar.QualityprofilesRemoveGroupOptions{
				QualityProfile: profileName,
				Language:       "java",
				Group:          "sonar-users",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Qualityprofiles.RemoveGroup(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing profile name", func() {
				resp, err := client.Qualityprofiles.RemoveGroup(&sonar.QualityprofilesRemoveGroupOptions{
					Language: "java",
					Group:    "sonar-users",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing group", func() {
				resp, err := client.Qualityprofiles.RemoveGroup(&sonar.QualityprofilesRemoveGroupOptions{
					QualityProfile: "test",
					Language:       "java",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// AddUser, SearchUsers, RemoveUser
	// =========================================================================
	Describe("AddUser", func() {
		It("should add user permission to quality profile", func() {
			profileName := helpers.UniqueResourceName("qp-usr")

			_, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName,
					Language:       "java",
				})
				return err
			})

			resp, err := client.Qualityprofiles.AddUser(&sonar.QualityprofilesAddUserOptions{
				QualityProfile: profileName,
				Language:       "java",
				Login:          "admin",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Qualityprofiles.AddUser(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing profile name", func() {
				resp, err := client.Qualityprofiles.AddUser(&sonar.QualityprofilesAddUserOptions{
					Language: "java",
					Login:    "admin",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing login", func() {
				resp, err := client.Qualityprofiles.AddUser(&sonar.QualityprofilesAddUserOptions{
					QualityProfile: "test",
					Language:       "java",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	Describe("SearchUsers", func() {
		It("should search users for a quality profile", func() {
			profileName := helpers.UniqueResourceName("qp-susr")

			_, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName,
					Language:       "java",
				})
				return err
			})

			result, resp, err := client.Qualityprofiles.SearchUsers(&sonar.QualityprofilesSearchUsersOptions{
				QualityProfile: profileName,
				Language:       "java",
				Selected:       "all",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Qualityprofiles.SearchUsers(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing profile name", func() {
				result, resp, err := client.Qualityprofiles.SearchUsers(&sonar.QualityprofilesSearchUsersOptions{
					Language: "java",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing language", func() {
				result, resp, err := client.Qualityprofiles.SearchUsers(&sonar.QualityprofilesSearchUsersOptions{
					QualityProfile: "test",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("RemoveUser", func() {
		It("should remove user permission from quality profile", func() {
			profileName := helpers.UniqueResourceName("qp-rmusr")

			_, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName,
					Language:       "java",
				})
				return err
			})

			// Add user first
			_, err = client.Qualityprofiles.AddUser(&sonar.QualityprofilesAddUserOptions{
				QualityProfile: profileName,
				Language:       "java",
				Login:          "admin",
			})
			Expect(err).NotTo(HaveOccurred())

			// Remove user
			resp, err := client.Qualityprofiles.RemoveUser(&sonar.QualityprofilesRemoveUserOptions{
				QualityProfile: profileName,
				Language:       "java",
				Login:          "admin",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Qualityprofiles.RemoveUser(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing profile name", func() {
				resp, err := client.Qualityprofiles.RemoveUser(&sonar.QualityprofilesRemoveUserOptions{
					Language: "java",
					Login:    "admin",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing login", func() {
				resp, err := client.Qualityprofiles.RemoveUser(&sonar.QualityprofilesRemoveUserOptions{
					QualityProfile: "test",
					Language:       "java",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Qualityprofiles Lifecycle
	// =========================================================================
	Describe("Qualityprofiles Lifecycle", func() {
		It("should complete full quality profile lifecycle", func() {
			profileName := helpers.UniqueResourceName("qp-lifecycle")
			projectKey := helpers.UniqueResourceName("proj-qplc")

			// Step 1: Create quality profile
			createResult, _, err := client.Qualityprofiles.Create(&sonar.QualityprofilesCreateOptions{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())
			profileKey := createResult.Profile.Key

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonar.QualityprofilesDeleteOptions{
					QualityProfile: profileName,
					Language:       "java",
				})
				return err
			})

			// Step 2: Show the profile
			showResult, _, err := client.Qualityprofiles.Show(&sonar.QualityprofilesShowOptions{
				Key: profileKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(showResult.Profile.Name).To(Equal(profileName))

			// Step 3: View changelog
			_, _, err = client.Qualityprofiles.Changelog(&sonar.QualityprofilesChangelogOptions{
				QualityProfile: profileName,
				Language:       "java",
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 4: Create project and associate
			_, _, err = client.Projects.Create(&sonar.ProjectsCreateOptions{
				Name:    "Lifecycle Test Project",
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonar.ProjectsDeleteOptions{
					Project: projectKey,
				})
				return err
			})

			_, err = client.Qualityprofiles.AddProject(&sonar.QualityprofilesAddProjectOptions{
				QualityProfile: profileName,
				Language:       "java",
				Project:        projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 5: List projects and verify our project is associated
			projectsResult, _, err := client.Qualityprofiles.Projects(&sonar.QualityprofilesProjectsOptions{
				Key: profileKey,
			})
			Expect(err).NotTo(HaveOccurred())
			foundProject := false
			for _, p := range projectsResult.Results {
				if p.Key == projectKey {
					foundProject = true
					break
				}
			}
			Expect(foundProject).To(BeTrue(), "Project %s should be associated with the quality profile", projectKey)

			// Step 6: Add user permission
			_, err = client.Qualityprofiles.AddUser(&sonar.QualityprofilesAddUserOptions{
				QualityProfile: profileName,
				Language:       "java",
				Login:          "admin",
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 7: Search users and verify admin is selected
			usersResult, _, err := client.Qualityprofiles.SearchUsers(&sonar.QualityprofilesSearchUsersOptions{
				QualityProfile: profileName,
				Language:       "java",
				Selected:       "selected",
			})
			Expect(err).NotTo(HaveOccurred())
			foundAdmin := false
			for _, u := range usersResult.Users {
				if u.Login == "admin" {
					foundAdmin = true
					break
				}
			}
			Expect(foundAdmin).To(BeTrue(), "Admin user should have permission on quality profile")

			// Step 8: Remove user permission
			_, err = client.Qualityprofiles.RemoveUser(&sonar.QualityprofilesRemoveUserOptions{
				QualityProfile: profileName,
				Language:       "java",
				Login:          "admin",
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 9: Remove project association
			_, err = client.Qualityprofiles.RemoveProject(&sonar.QualityprofilesRemoveProjectOptions{
				QualityProfile: profileName,
				Language:       "java",
				Project:        projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 10: Backup profile
			backupResult, _, err := client.Qualityprofiles.Backup(&sonar.QualityprofilesBackupOptions{
				QualityProfile: profileName,
				Language:       "java",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(*backupResult).NotTo(BeEmpty())
		})
	})
})
