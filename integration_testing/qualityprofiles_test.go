package integration_testing_test

import (
"net/http"

. "github.com/onsi/ginkgo/v2"
. "github.com/onsi/gomega"

sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Qualityprofiles Service", Ordered, func() {
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
			result, resp, err := client.Qualityprofiles.Search(&sonargo.QualityprofilesSearchOption{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Profiles).NotTo(BeEmpty())
		})

		It("should filter by language", func() {
			result, resp, err := client.Qualityprofiles.Search(&sonargo.QualityprofilesSearchOption{
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			for _, profile := range result.Profiles {
				Expect(profile.Language).To(Equal("java"))
			}
		})

		It("should include defaults", func() {
			result, resp, err := client.Qualityprofiles.Search(&sonargo.QualityprofilesSearchOption{
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

			result, resp, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Profile.Name).To(Equal(profileName))

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
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
				result, resp, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
					Language: "java",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing language", func() {
				result, resp, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
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

			createResult, _, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
					QualityProfile: profileName,
					Language:       "java",
				})
				return err
			})

			result, resp, err := client.Qualityprofiles.Show(&sonargo.QualityprofilesShowOption{
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
				result, resp, err := client.Qualityprofiles.Show(&sonargo.QualityprofilesShowOption{})
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

			createResult, _, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
				Name:     originalName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", newName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
					QualityProfile: newName,
					Language:       "java",
				})
				return err
			})

			resp, err := client.Qualityprofiles.Rename(&sonargo.QualityprofilesRenameOption{
				Key:  createResult.Profile.Key,
				Name: newName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify rename
			result, _, err := client.Qualityprofiles.Show(&sonargo.QualityprofilesShowOption{
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
				resp, err := client.Qualityprofiles.Rename(&sonargo.QualityprofilesRenameOption{
					Name: "new-name",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing name", func() {
				resp, err := client.Qualityprofiles.Rename(&sonargo.QualityprofilesRenameOption{
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

			createResult, _, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
				Name:     sourceName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", sourceName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
					QualityProfile: sourceName,
					Language:       "java",
				})
				return err
			})

			result, resp, err := client.Qualityprofiles.Copy(&sonargo.QualityprofilesCopyOption{
				FromKey: createResult.Profile.Key,
				ToName:  copyName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Name).To(Equal(copyName))

			cleanup.RegisterCleanup("qualityprofile", copyName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
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
				result, resp, err := client.Qualityprofiles.Copy(&sonargo.QualityprofilesCopyOption{
					ToName: "new-copy",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with missing to name", func() {
				result, resp, err := client.Qualityprofiles.Copy(&sonargo.QualityprofilesCopyOption{
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

			_, _, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			resp, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
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
				resp, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
					Language: "java",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing language", func() {
				resp, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
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

			createResult, _, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
				Name:     profileName,
				Language: "java",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
				_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
					QualityProfile: profileName,
					Language:       "java",
				})
				return err
			})

			resp, err := client.Qualityprofiles.SetDefault(&sonargo.QualityprofilesSetDefaultOption{
				QualityProfile: profileName,
				Language:       "java",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify it's now default
result, _, err := client.Qualityprofiles.Show(&sonargo.QualityprofilesShowOption{
Key: createResult.Profile.Key,
})
Expect(err).NotTo(HaveOccurred())
Expect(result.Profile.IsDefault).To(BeTrue())

// Restore the original default
searchResult, _, _ := client.Qualityprofiles.Search(&sonargo.QualityprofilesSearchOption{
Language: "java",
})
for _, p := range searchResult.Profiles {
if p.IsBuiltIn && p.Key != createResult.Profile.Key {
_, _ = client.Qualityprofiles.SetDefault(&sonargo.QualityprofilesSetDefaultOption{
QualityProfile: p.Name,
Language:       "java",
})
break
}
}
})

Context("parameter validation", func() {
It("should fail with nil options", func() {
resp, err := client.Qualityprofiles.SetDefault(nil)
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with missing profile name", func() {
resp, err := client.Qualityprofiles.SetDefault(&sonargo.QualityprofilesSetDefaultOption{
Language: "java",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with missing language", func() {
resp, err := client.Qualityprofiles.SetDefault(&sonargo.QualityprofilesSetDefaultOption{
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

_, _, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
Name:     profileName,
Language: "java",
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
QualityProfile: profileName,
Language:       "java",
})
return err
})

_, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
Name:    "QualityProfile AddProject Test",
Project: projectKey,
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("project", projectKey, func() error {
_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
Project: projectKey,
})
return err
})

resp, err := client.Qualityprofiles.AddProject(&sonargo.QualityprofilesAddProjectOption{
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
resp, err := client.Qualityprofiles.AddProject(&sonargo.QualityprofilesAddProjectOption{
Language: "java",
Project:  "some-project",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with missing project", func() {
resp, err := client.Qualityprofiles.AddProject(&sonargo.QualityprofilesAddProjectOption{
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

_, _, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
Name:     profileName,
Language: "java",
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
QualityProfile: profileName,
Language:       "java",
})
return err
})

_, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
Name:    "QualityProfile RemoveProject Test",
Project: projectKey,
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("project", projectKey, func() error {
_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
Project: projectKey,
})
return err
})

// Add first
_, err = client.Qualityprofiles.AddProject(&sonargo.QualityprofilesAddProjectOption{
QualityProfile: profileName,
Language:       "java",
Project:        projectKey,
})
Expect(err).NotTo(HaveOccurred())

// Remove
resp, err := client.Qualityprofiles.RemoveProject(&sonargo.QualityprofilesRemoveProjectOption{
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
resp, err := client.Qualityprofiles.RemoveProject(&sonargo.QualityprofilesRemoveProjectOption{
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

createResult, _, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
Name:     profileName,
Language: "java",
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
QualityProfile: profileName,
Language:       "java",
})
return err
})

result, resp, err := client.Qualityprofiles.Projects(&sonargo.QualityprofilesProjectsOption{
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
result, resp, err := client.Qualityprofiles.Projects(&sonargo.QualityprofilesProjectsOption{})
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

_, _, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
Name:     profileName,
Language: "java",
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
QualityProfile: profileName,
Language:       "java",
})
return err
})

result, resp, err := client.Qualityprofiles.Inheritance(&sonargo.QualityprofilesInheritanceOption{
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
result, resp, err := client.Qualityprofiles.Inheritance(&sonargo.QualityprofilesInheritanceOption{
Language: "java",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
Expect(result).To(BeNil())
})

It("should fail with missing language", func() {
result, resp, err := client.Qualityprofiles.Inheritance(&sonargo.QualityprofilesInheritanceOption{
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

_, _, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
Name:     parentName,
Language: "java",
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualityprofile", parentName, func() error {
_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
QualityProfile: parentName,
Language:       "java",
})
return err
})

_, _, err = client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
Name:     childName,
Language: "java",
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualityprofile", childName, func() error {
// Remove parent first
_, _ = client.Qualityprofiles.ChangeParent(&sonargo.QualityprofilesChangeParentOption{
QualityProfile: childName,
Language:       "java",
})
_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
QualityProfile: childName,
Language:       "java",
})
return err
})

resp, err := client.Qualityprofiles.ChangeParent(&sonargo.QualityprofilesChangeParentOption{
QualityProfile:       childName,
Language:             "java",
ParentQualityProfile: parentName,
})
Expect(err).NotTo(HaveOccurred())
Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

// Verify parent was set by checking ancestors
result, _, err := client.Qualityprofiles.Inheritance(&sonargo.QualityprofilesInheritanceOption{
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
resp, err := client.Qualityprofiles.ChangeParent(&sonargo.QualityprofilesChangeParentOption{
Language:             "java",
ParentQualityProfile: "some-parent",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with missing language", func() {
resp, err := client.Qualityprofiles.ChangeParent(&sonargo.QualityprofilesChangeParentOption{
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

result1, _, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
Name:     profileName1,
Language: "java",
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualityprofile", profileName1, func() error {
_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
QualityProfile: profileName1,
Language:       "java",
})
return err
})

result2, _, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
Name:     profileName2,
Language: "java",
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualityprofile", profileName2, func() error {
_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
QualityProfile: profileName2,
Language:       "java",
})
return err
})

compareResult, resp, err := client.Qualityprofiles.Compare(&sonargo.QualityprofilesCompareOption{
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
result, resp, err := client.Qualityprofiles.Compare(&sonargo.QualityprofilesCompareOption{
RightKey: "some-key",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
Expect(result).To(BeNil())
})

It("should fail with missing right key", func() {
result, resp, err := client.Qualityprofiles.Compare(&sonargo.QualityprofilesCompareOption{
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

_, _, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
Name:     profileName,
Language: "java",
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
QualityProfile: profileName,
Language:       "java",
})
return err
})

result, resp, err := client.Qualityprofiles.Changelog(&sonargo.QualityprofilesChangelogOption{
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
result, resp, err := client.Qualityprofiles.Changelog(&sonargo.QualityprofilesChangelogOption{
Language: "java",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
Expect(result).To(BeNil())
})

It("should fail with missing language", func() {
result, resp, err := client.Qualityprofiles.Changelog(&sonargo.QualityprofilesChangelogOption{
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

_, _, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
Name:     profileName,
Language: "java",
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
QualityProfile: profileName,
Language:       "java",
})
return err
})

result, resp, err := client.Qualityprofiles.Backup(&sonargo.QualityprofilesBackupOption{
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
result, resp, err := client.Qualityprofiles.Backup(&sonargo.QualityprofilesBackupOption{
Language: "java",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
Expect(result).To(BeNil())
})

It("should fail with missing language", func() {
result, resp, err := client.Qualityprofiles.Backup(&sonargo.QualityprofilesBackupOption{
QualityProfile: "test",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
Expect(result).To(BeNil())
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

_, _, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
Name:     profileName,
Language: "java",
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
QualityProfile: profileName,
Language:       "java",
})
return err
})

resp, err := client.Qualityprofiles.AddGroup(&sonargo.QualityprofilesAddGroupOption{
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
resp, err := client.Qualityprofiles.AddGroup(&sonargo.QualityprofilesAddGroupOption{
Language: "java",
Group:    "sonar-users",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with missing group", func() {
resp, err := client.Qualityprofiles.AddGroup(&sonargo.QualityprofilesAddGroupOption{
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

_, _, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
Name:     profileName,
Language: "java",
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
QualityProfile: profileName,
Language:       "java",
})
return err
})

result, resp, err := client.Qualityprofiles.SearchGroups(&sonargo.QualityprofilesSearchGroupsOption{
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
result, resp, err := client.Qualityprofiles.SearchGroups(&sonargo.QualityprofilesSearchGroupsOption{
Language: "java",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
Expect(result).To(BeNil())
})

It("should fail with missing language", func() {
result, resp, err := client.Qualityprofiles.SearchGroups(&sonargo.QualityprofilesSearchGroupsOption{
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

_, _, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
Name:     profileName,
Language: "java",
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
QualityProfile: profileName,
Language:       "java",
})
return err
})

// Add group first
_, err = client.Qualityprofiles.AddGroup(&sonargo.QualityprofilesAddGroupOption{
QualityProfile: profileName,
Language:       "java",
Group:          "sonar-users",
})
Expect(err).NotTo(HaveOccurred())

// Remove group
resp, err := client.Qualityprofiles.RemoveGroup(&sonargo.QualityprofilesRemoveGroupOption{
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
resp, err := client.Qualityprofiles.RemoveGroup(&sonargo.QualityprofilesRemoveGroupOption{
Language: "java",
Group:    "sonar-users",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with missing group", func() {
resp, err := client.Qualityprofiles.RemoveGroup(&sonargo.QualityprofilesRemoveGroupOption{
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

_, _, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
Name:     profileName,
Language: "java",
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
QualityProfile: profileName,
Language:       "java",
})
return err
})

resp, err := client.Qualityprofiles.AddUser(&sonargo.QualityprofilesAddUserOption{
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
resp, err := client.Qualityprofiles.AddUser(&sonargo.QualityprofilesAddUserOption{
Language: "java",
Login:    "admin",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with missing login", func() {
resp, err := client.Qualityprofiles.AddUser(&sonargo.QualityprofilesAddUserOption{
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

_, _, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
Name:     profileName,
Language: "java",
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
QualityProfile: profileName,
Language:       "java",
})
return err
})

result, resp, err := client.Qualityprofiles.SearchUsers(&sonargo.QualityprofilesSearchUsersOption{
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
result, resp, err := client.Qualityprofiles.SearchUsers(&sonargo.QualityprofilesSearchUsersOption{
Language: "java",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
Expect(result).To(BeNil())
})

It("should fail with missing language", func() {
result, resp, err := client.Qualityprofiles.SearchUsers(&sonargo.QualityprofilesSearchUsersOption{
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

_, _, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
Name:     profileName,
Language: "java",
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
QualityProfile: profileName,
Language:       "java",
})
return err
})

// Add user first
_, err = client.Qualityprofiles.AddUser(&sonargo.QualityprofilesAddUserOption{
QualityProfile: profileName,
Language:       "java",
Login:          "admin",
})
Expect(err).NotTo(HaveOccurred())

// Remove user
resp, err := client.Qualityprofiles.RemoveUser(&sonargo.QualityprofilesRemoveUserOption{
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
resp, err := client.Qualityprofiles.RemoveUser(&sonargo.QualityprofilesRemoveUserOption{
Language: "java",
Login:    "admin",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with missing login", func() {
resp, err := client.Qualityprofiles.RemoveUser(&sonargo.QualityprofilesRemoveUserOption{
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
createResult, _, err := client.Qualityprofiles.Create(&sonargo.QualityprofilesCreateOption{
Name:     profileName,
Language: "java",
})
Expect(err).NotTo(HaveOccurred())
profileKey := createResult.Profile.Key

cleanup.RegisterCleanup("qualityprofile", profileName, func() error {
_, err := client.Qualityprofiles.Delete(&sonargo.QualityprofilesDeleteOption{
QualityProfile: profileName,
Language:       "java",
})
return err
})

// Step 2: Show the profile
showResult, _, err := client.Qualityprofiles.Show(&sonargo.QualityprofilesShowOption{
Key: profileKey,
})
Expect(err).NotTo(HaveOccurred())
Expect(showResult.Profile.Name).To(Equal(profileName))

// Step 3: View changelog
_, _, err = client.Qualityprofiles.Changelog(&sonargo.QualityprofilesChangelogOption{
QualityProfile: profileName,
Language:       "java",
})
Expect(err).NotTo(HaveOccurred())

// Step 4: Create project and associate
_, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
Name:    "Lifecycle Test Project",
Project: projectKey,
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("project", projectKey, func() error {
_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
Project: projectKey,
})
return err
})

_, err = client.Qualityprofiles.AddProject(&sonargo.QualityprofilesAddProjectOption{
QualityProfile: profileName,
Language:       "java",
Project:        projectKey,
})
Expect(err).NotTo(HaveOccurred())

// Step 5: List projects
projectsResult, _, err := client.Qualityprofiles.Projects(&sonargo.QualityprofilesProjectsOption{
Key: profileKey,
})
Expect(err).NotTo(HaveOccurred())
Expect(projectsResult.Results).To(HaveLen(1))

// Step 6: Add user permission
_, err = client.Qualityprofiles.AddUser(&sonargo.QualityprofilesAddUserOption{
QualityProfile: profileName,
Language:       "java",
Login:          "admin",
})
Expect(err).NotTo(HaveOccurred())

// Step 7: Search users
usersResult, _, err := client.Qualityprofiles.SearchUsers(&sonargo.QualityprofilesSearchUsersOption{
QualityProfile: profileName,
Language:       "java",
Selected:       "selected",
})
Expect(err).NotTo(HaveOccurred())
Expect(usersResult.Users).To(HaveLen(1))

// Step 8: Remove user permission
_, err = client.Qualityprofiles.RemoveUser(&sonargo.QualityprofilesRemoveUserOption{
QualityProfile: profileName,
Language:       "java",
Login:          "admin",
})
Expect(err).NotTo(HaveOccurred())

// Step 9: Remove project association
_, err = client.Qualityprofiles.RemoveProject(&sonargo.QualityprofilesRemoveProjectOption{
QualityProfile: profileName,
Language:       "java",
Project:        projectKey,
})
Expect(err).NotTo(HaveOccurred())

// Step 10: Backup profile
backupResult, _, err := client.Qualityprofiles.Backup(&sonargo.QualityprofilesBackupOption{
QualityProfile: profileName,
Language:       "java",
})
Expect(err).NotTo(HaveOccurred())
Expect(*backupResult).NotTo(BeEmpty())
})
})
})
