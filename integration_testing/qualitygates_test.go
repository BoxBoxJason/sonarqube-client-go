package integration_testing_test

import (
"net/http"

. "github.com/onsi/ginkgo/v2"
. "github.com/onsi/gomega"

sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Qualitygates Service", Ordered, func() {
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
	// List
	// =========================================================================
	Describe("List", func() {
		It("should list all quality gates", func() {
			result, resp, err := client.Qualitygates.List()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			// SonarQube has a built-in default quality gate
			Expect(result.Qualitygates).NotTo(BeEmpty())
		})

		It("should include built-in quality gate", func() {
			result, resp, err := client.Qualitygates.List()
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			foundBuiltIn := false
			for _, gate := range result.Qualitygates {
				if gate.IsBuiltIn {
					foundBuiltIn = true
					break
				}
			}
			Expect(foundBuiltIn).To(BeTrue())
		})
	})

	// =========================================================================
	// Create
	// =========================================================================
	Describe("Create", func() {
		It("should create a new quality gate", func() {
			gateName := helpers.UniqueResourceName("qg")

			result, resp, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
				Name: gateName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Name).To(Equal(gateName))

			cleanup.RegisterCleanup("qualitygate", gateName, func() error {
				_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
					Name: gateName,
				})
				return err
			})
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Qualitygates.Create(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with empty name", func() {
				result, resp, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Context("error cases", func() {
			It("should fail with duplicate name", func() {
				gateName := helpers.UniqueResourceName("qg-dup")

				_, _, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
					Name: gateName,
				})
				Expect(err).NotTo(HaveOccurred())

				cleanup.RegisterCleanup("qualitygate", gateName, func() error {
					_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
						Name: gateName,
					})
					return err
				})

				// Try to create with same name
				result, resp, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
					Name: gateName,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Show
	// =========================================================================
	Describe("Show", func() {
		It("should show quality gate details", func() {
			gateName := helpers.UniqueResourceName("qg-show")

			_, _, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
				Name: gateName,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualitygate", gateName, func() error {
				_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
					Name: gateName,
				})
				return err
			})

			result, resp, err := client.Qualitygates.Show(&sonargo.QualitygatesShowOption{
				Name: gateName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			Expect(result.Name).To(Equal(gateName))
			Expect(result.IsBuiltIn).To(BeFalse())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Qualitygates.Show(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})

			It("should fail with empty name", func() {
				result, resp, err := client.Qualitygates.Show(&sonargo.QualitygatesShowOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
				Expect(result).To(BeNil())
			})
		})

		Context("error cases", func() {
			It("should fail for non-existent quality gate", func() {
				result, resp, err := client.Qualitygates.Show(&sonargo.QualitygatesShowOption{
					Name: "non-existent-gate-xyz",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
				Expect(result).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Rename
	// =========================================================================
	Describe("Rename", func() {
		It("should rename a quality gate", func() {
			originalName := helpers.UniqueResourceName("qg-rename")
			newName := helpers.UniqueResourceName("qg-renamed")

			_, _, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
				Name: originalName,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualitygate", newName, func() error {
				_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
					Name: newName,
				})
				return err
			})

			resp, err := client.Qualitygates.Rename(&sonargo.QualitygatesRenameOption{
				CurrentName: originalName,
				Name:        newName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			// Verify rename succeeded
			result, _, err := client.Qualitygates.Show(&sonargo.QualitygatesShowOption{
				Name: newName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Name).To(Equal(newName))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Qualitygates.Rename(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with empty current name", func() {
				resp, err := client.Qualitygates.Rename(&sonargo.QualitygatesRenameOption{
					Name: "new-name",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with empty new name", func() {
				resp, err := client.Qualitygates.Rename(&sonargo.QualitygatesRenameOption{
					CurrentName: "current-name",
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
		It("should copy a quality gate", func() {
			sourceName := helpers.UniqueResourceName("qg-source")
			copyName := helpers.UniqueResourceName("qg-copy")

			_, _, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
				Name: sourceName,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualitygate", sourceName, func() error {
				_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
					Name: sourceName,
				})
				return err
			})

			resp, err := client.Qualitygates.Copy(&sonargo.QualitygatesCopyOption{
				SourceName: sourceName,
				Name:       copyName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			cleanup.RegisterCleanup("qualitygate", copyName, func() error {
				_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
					Name: copyName,
				})
				return err
			})

			// Verify copy exists
			result, _, err := client.Qualitygates.Show(&sonargo.QualitygatesShowOption{
				Name: copyName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Name).To(Equal(copyName))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Qualitygates.Copy(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with empty source name", func() {
				resp, err := client.Qualitygates.Copy(&sonargo.QualitygatesCopyOption{
					Name: "new-copy",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with empty target name", func() {
				resp, err := client.Qualitygates.Copy(&sonargo.QualitygatesCopyOption{
					SourceName: "source",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Destroy
	// =========================================================================
	Describe("Destroy", func() {
		It("should delete a quality gate", func() {
			gateName := helpers.UniqueResourceName("qg-destroy")

			_, _, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
				Name: gateName,
			})
			Expect(err).NotTo(HaveOccurred())

			resp, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
				Name: gateName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify deletion
			_, _, err = client.Qualitygates.Show(&sonargo.QualitygatesShowOption{
				Name: gateName,
			})
			Expect(err).To(HaveOccurred())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Qualitygates.Destroy(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with empty name", func() {
				resp, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})

		Context("error cases", func() {
			It("should fail for non-existent quality gate", func() {
				resp, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
					Name: "non-existent-gate-xyz",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	// =========================================================================
	// SetAsDefault
	// =========================================================================
	Describe("SetAsDefault", func() {
		It("should set a quality gate as default", func() {
			gateName := helpers.UniqueResourceName("qg-default")

			_, _, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
				Name: gateName,
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("qualitygate", gateName, func() error {
				_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
					Name: gateName,
				})
				return err
			})

			resp, err := client.Qualitygates.SetAsDefault(&sonargo.QualitygatesSetAsDefaultOption{
				Name: gateName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify it's now default
result, _, err := client.Qualitygates.Show(&sonargo.QualitygatesShowOption{
Name: gateName,
})
Expect(err).NotTo(HaveOccurred())
Expect(result.IsDefault).To(BeTrue())

// Restore Sonar way as default
listResult, _, _ := client.Qualitygates.List()
for _, gate := range listResult.Qualitygates {
if gate.IsBuiltIn {
_, _ = client.Qualitygates.SetAsDefault(&sonargo.QualitygatesSetAsDefaultOption{
Name: gate.Name,
})
break
}
}
})

Context("parameter validation", func() {
It("should fail with nil options", func() {
resp, err := client.Qualitygates.SetAsDefault(nil)
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with empty name", func() {
resp, err := client.Qualitygates.SetAsDefault(&sonargo.QualitygatesSetAsDefaultOption{})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})
})
})

// =========================================================================
// CreateCondition
// =========================================================================
Describe("CreateCondition", func() {
It("should create a condition for a quality gate", func() {
gateName := helpers.UniqueResourceName("qg-cond")

_, _, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
Name: gateName,
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualitygate", gateName, func() error {
_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
Name: gateName,
})
return err
})

result, resp, err := client.Qualitygates.CreateCondition(&sonargo.QualitygatesCreateConditionOption{
GateName: gateName,
Metric:   "coverage",
Op:       "LT",
Error:    "80",
})
Expect(err).NotTo(HaveOccurred())
Expect(resp.StatusCode).To(Equal(http.StatusOK))
Expect(result).NotTo(BeNil())
Expect(result.Metric).To(Equal("coverage"))
Expect(result.Error).To(Equal("80"))
Expect(result.ID).NotTo(BeEmpty())
})

It("should create condition with GT operator", func() {
gateName := helpers.UniqueResourceName("qg-condgt")

_, _, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
Name: gateName,
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualitygate", gateName, func() error {
_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
Name: gateName,
})
return err
})

result, resp, err := client.Qualitygates.CreateCondition(&sonargo.QualitygatesCreateConditionOption{
GateName: gateName,
Metric:   "duplicated_lines_density",
Op:       "GT",
Error:    "5",
})
Expect(err).NotTo(HaveOccurred())
Expect(resp.StatusCode).To(Equal(http.StatusOK))
Expect(result.Op).To(Equal("GT"))
})

Context("parameter validation", func() {
It("should fail with nil options", func() {
result, resp, err := client.Qualitygates.CreateCondition(nil)
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
Expect(result).To(BeNil())
})

It("should fail with missing gate name", func() {
result, resp, err := client.Qualitygates.CreateCondition(&sonargo.QualitygatesCreateConditionOption{
Metric: "coverage",
Error:  "80",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
Expect(result).To(BeNil())
})

It("should fail with missing metric", func() {
result, resp, err := client.Qualitygates.CreateCondition(&sonargo.QualitygatesCreateConditionOption{
GateName: "some-gate",
Error:    "80",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
Expect(result).To(BeNil())
})

It("should fail with missing error threshold", func() {
result, resp, err := client.Qualitygates.CreateCondition(&sonargo.QualitygatesCreateConditionOption{
GateName: "some-gate",
Metric:   "coverage",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
Expect(result).To(BeNil())
})
})
})

// =========================================================================
// UpdateCondition
// =========================================================================
Describe("UpdateCondition", func() {
It("should update a condition", func() {
gateName := helpers.UniqueResourceName("qg-upd")

_, _, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
Name: gateName,
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualitygate", gateName, func() error {
_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
Name: gateName,
})
return err
})

// Create a condition first
condResult, _, err := client.Qualitygates.CreateCondition(&sonargo.QualitygatesCreateConditionOption{
GateName: gateName,
Metric:   "coverage",
Op:       "LT",
Error:    "80",
})
Expect(err).NotTo(HaveOccurred())

// Update it
resp, err := client.Qualitygates.UpdateCondition(&sonargo.QualitygatesUpdateConditionOption{
ID:     condResult.ID,
Metric: "coverage",
Op:     "LT",
Error:  "90",
})
Expect(err).NotTo(HaveOccurred())
Expect(resp.StatusCode).To(Equal(http.StatusOK))

// Verify update - find the condition with our ID
result, _, err := client.Qualitygates.Show(&sonargo.QualitygatesShowOption{
Name: gateName,
})
Expect(err).NotTo(HaveOccurred())
Expect(result.Conditions).NotTo(BeEmpty())
found := false
for _, cond := range result.Conditions {
if cond.ID == condResult.ID {
Expect(cond.Error).To(Equal("90"))
found = true
break
}
}
Expect(found).To(BeTrue(), "Updated condition not found")
})

Context("parameter validation", func() {
It("should fail with nil options", func() {
resp, err := client.Qualitygates.UpdateCondition(nil)
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with missing ID", func() {
resp, err := client.Qualitygates.UpdateCondition(&sonargo.QualitygatesUpdateConditionOption{
Metric: "coverage",
Error:  "80",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})
})
})

// =========================================================================
// DeleteCondition
// =========================================================================
Describe("DeleteCondition", func() {
It("should delete a condition", func() {
gateName := helpers.UniqueResourceName("qg-delcond")

_, _, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
Name: gateName,
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualitygate", gateName, func() error {
_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
Name: gateName,
})
return err
})

// Create a condition first
condResult, _, err := client.Qualitygates.CreateCondition(&sonargo.QualitygatesCreateConditionOption{
GateName: gateName,
Metric:   "coverage",
Op:       "LT",
Error:    "80",
})
Expect(err).NotTo(HaveOccurred())

// Delete it
resp, err := client.Qualitygates.DeleteCondition(&sonargo.QualitygatesDeleteConditionOption{
ID: condResult.ID,
})
Expect(err).NotTo(HaveOccurred())
Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

// Verify deletion - the specific condition should no longer exist
result, _, err := client.Qualitygates.Show(&sonargo.QualitygatesShowOption{
Name: gateName,
})
Expect(err).NotTo(HaveOccurred())
for _, cond := range result.Conditions {
Expect(cond.ID).NotTo(Equal(condResult.ID), "Deleted condition still exists")
}
})

Context("parameter validation", func() {
It("should fail with nil options", func() {
resp, err := client.Qualitygates.DeleteCondition(nil)
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with missing ID", func() {
resp, err := client.Qualitygates.DeleteCondition(&sonargo.QualitygatesDeleteConditionOption{})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})
})
})

// =========================================================================
// Select and Deselect (Project association)
// =========================================================================
Describe("Select", func() {
It("should associate a project with a quality gate", func() {
gateName := helpers.UniqueResourceName("qg-sel")
projectKey := helpers.UniqueResourceName("proj-qg")

_, _, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
Name: gateName,
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualitygate", gateName, func() error {
_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
Name: gateName,
})
return err
})

_, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
Name:    "QualityGate Select Test Project",
Project: projectKey,
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("project", projectKey, func() error {
_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
Project: projectKey,
})
return err
})

resp, err := client.Qualitygates.Select(&sonargo.QualitygatesSelectOption{
GateName:   gateName,
ProjectKey: projectKey,
})
Expect(err).NotTo(HaveOccurred())
Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

// Verify association
result, _, err := client.Qualitygates.GetByProject(&sonargo.QualitygatesGetByProjectOption{
Project: projectKey,
})
Expect(err).NotTo(HaveOccurred())
Expect(result.QualityGate.Name).To(Equal(gateName))
})

Context("parameter validation", func() {
It("should fail with nil options", func() {
resp, err := client.Qualitygates.Select(nil)
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with missing gate name", func() {
resp, err := client.Qualitygates.Select(&sonargo.QualitygatesSelectOption{
ProjectKey: "some-project",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with missing project key", func() {
resp, err := client.Qualitygates.Select(&sonargo.QualitygatesSelectOption{
GateName: "some-gate",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})
})
})

Describe("Deselect", func() {
It("should remove project association from quality gate", func() {
gateName := helpers.UniqueResourceName("qg-desel")
projectKey := helpers.UniqueResourceName("proj-qgde")

_, _, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
Name: gateName,
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualitygate", gateName, func() error {
_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
Name: gateName,
})
return err
})

_, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
Name:    "QualityGate Deselect Test Project",
Project: projectKey,
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("project", projectKey, func() error {
_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
Project: projectKey,
})
return err
})

// Associate first
_, err = client.Qualitygates.Select(&sonargo.QualitygatesSelectOption{
GateName:   gateName,
ProjectKey: projectKey,
})
Expect(err).NotTo(HaveOccurred())

// Deselect
resp, err := client.Qualitygates.Deselect(&sonargo.QualitygatesDeselectOption{
ProjectKey: projectKey,
})
Expect(err).NotTo(HaveOccurred())
Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

// Verify - project should now use default gate
result, _, err := client.Qualitygates.GetByProject(&sonargo.QualitygatesGetByProjectOption{
Project: projectKey,
})
Expect(err).NotTo(HaveOccurred())
Expect(result.QualityGate.Default).To(BeTrue())
})

Context("parameter validation", func() {
It("should fail with nil options", func() {
resp, err := client.Qualitygates.Deselect(nil)
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with missing project key", func() {
resp, err := client.Qualitygates.Deselect(&sonargo.QualitygatesDeselectOption{})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})
})
})

// =========================================================================
// GetByProject
// =========================================================================
Describe("GetByProject", func() {
It("should get quality gate for a project", func() {
projectKey := helpers.UniqueResourceName("proj-getqg")

_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
Name:    "GetByProject Test Project",
Project: projectKey,
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("project", projectKey, func() error {
_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
Project: projectKey,
})
return err
})

result, resp, err := client.Qualitygates.GetByProject(&sonargo.QualitygatesGetByProjectOption{
Project: projectKey,
})
Expect(err).NotTo(HaveOccurred())
Expect(resp.StatusCode).To(Equal(http.StatusOK))
Expect(result).NotTo(BeNil())
// New projects use the default quality gate
Expect(result.QualityGate.Name).NotTo(BeEmpty())
})

Context("parameter validation", func() {
It("should fail with nil options", func() {
result, resp, err := client.Qualitygates.GetByProject(nil)
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
Expect(result).To(BeNil())
})

It("should fail with empty project", func() {
result, resp, err := client.Qualitygates.GetByProject(&sonargo.QualitygatesGetByProjectOption{})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
Expect(result).To(BeNil())
})
})

Context("error cases", func() {
It("should fail for non-existent project", func() {
result, resp, err := client.Qualitygates.GetByProject(&sonargo.QualitygatesGetByProjectOption{
Project: "non-existent-project-xyz",
})
Expect(err).To(HaveOccurred())
Expect(resp).NotTo(BeNil())
Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
Expect(result).To(BeNil())
})
})
})

// =========================================================================
// Search (Projects)
// =========================================================================
Describe("Search", func() {
It("should search projects for a quality gate", func() {
gateName := helpers.UniqueResourceName("qg-search")

_, _, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
Name: gateName,
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualitygate", gateName, func() error {
_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
Name: gateName,
})
return err
})

result, resp, err := client.Qualitygates.Search(&sonargo.QualitygatesSearchOption{
GateName: gateName,
Selected: "all",
})
Expect(err).NotTo(HaveOccurred())
Expect(resp.StatusCode).To(Equal(http.StatusOK))
Expect(result).NotTo(BeNil())
})

It("should search with query filter", func() {
gateName := helpers.UniqueResourceName("qg-searchq")
projectKey := helpers.UniqueResourceName("proj-searchq")

_, _, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
Name: gateName,
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualitygate", gateName, func() error {
_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
Name: gateName,
})
return err
})

_, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
Name:    "Search Query Test Project",
Project: projectKey,
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("project", projectKey, func() error {
_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
Project: projectKey,
})
return err
})

// Associate project with gate
_, err = client.Qualitygates.Select(&sonargo.QualitygatesSelectOption{
GateName:   gateName,
ProjectKey: projectKey,
})
Expect(err).NotTo(HaveOccurred())

// Search
result, resp, err := client.Qualitygates.Search(&sonargo.QualitygatesSearchOption{
GateName: gateName,
Query:    projectKey,
})
Expect(err).NotTo(HaveOccurred())
Expect(resp.StatusCode).To(Equal(http.StatusOK))
Expect(result.Results).To(HaveLen(1))
Expect(result.Results[0].Key).To(Equal(projectKey))
})

Context("parameter validation", func() {
It("should fail with nil options", func() {
result, resp, err := client.Qualitygates.Search(nil)
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
Expect(result).To(BeNil())
})

It("should fail with missing gate name", func() {
result, resp, err := client.Qualitygates.Search(&sonargo.QualitygatesSearchOption{})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
Expect(result).To(BeNil())
})
})
})

// =========================================================================
// AddGroup, SearchGroups, RemoveGroup
// =========================================================================
Describe("AddGroup", func() {
It("should add group permission to quality gate", func() {
gateName := helpers.UniqueResourceName("qg-grp")

_, _, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
Name: gateName,
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualitygate", gateName, func() error {
_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
Name: gateName,
})
return err
})

resp, err := client.Qualitygates.AddGroup(&sonargo.QualitygatesAddGroupOption{
GateName:  gateName,
GroupName: "sonar-users",
})
Expect(err).NotTo(HaveOccurred())
Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
})

Context("parameter validation", func() {
It("should fail with nil options", func() {
resp, err := client.Qualitygates.AddGroup(nil)
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with missing gate name", func() {
resp, err := client.Qualitygates.AddGroup(&sonargo.QualitygatesAddGroupOption{
GroupName: "sonar-users",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with missing group name", func() {
resp, err := client.Qualitygates.AddGroup(&sonargo.QualitygatesAddGroupOption{
GateName: "some-gate",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})
})
})

Describe("SearchGroups", func() {
It("should search groups for a quality gate", func() {
gateName := helpers.UniqueResourceName("qg-sgrp")

_, _, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
Name: gateName,
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualitygate", gateName, func() error {
_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
Name: gateName,
})
return err
})

result, resp, err := client.Qualitygates.SearchGroups(&sonargo.QualitygatesSearchGroupsOption{
GateName: gateName,
Selected: "all",
})
Expect(err).NotTo(HaveOccurred())
Expect(resp.StatusCode).To(Equal(http.StatusOK))
Expect(result).NotTo(BeNil())
})

Context("parameter validation", func() {
It("should fail with nil options", func() {
result, resp, err := client.Qualitygates.SearchGroups(nil)
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
Expect(result).To(BeNil())
})

It("should fail with missing gate name", func() {
result, resp, err := client.Qualitygates.SearchGroups(&sonargo.QualitygatesSearchGroupsOption{})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
Expect(result).To(BeNil())
})
})
})

Describe("RemoveGroup", func() {
It("should remove group permission from quality gate", func() {
gateName := helpers.UniqueResourceName("qg-rmgrp")

_, _, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
Name: gateName,
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualitygate", gateName, func() error {
_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
Name: gateName,
})
return err
})

// Add group first
_, err = client.Qualitygates.AddGroup(&sonargo.QualitygatesAddGroupOption{
GateName:  gateName,
GroupName: "sonar-users",
})
Expect(err).NotTo(HaveOccurred())

// Remove group
resp, err := client.Qualitygates.RemoveGroup(&sonargo.QualitygatesRemoveGroupOption{
GateName:  gateName,
GroupName: "sonar-users",
})
Expect(err).NotTo(HaveOccurred())
Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
})

Context("parameter validation", func() {
It("should fail with nil options", func() {
resp, err := client.Qualitygates.RemoveGroup(nil)
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with missing gate name", func() {
resp, err := client.Qualitygates.RemoveGroup(&sonargo.QualitygatesRemoveGroupOption{
GroupName: "sonar-users",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with missing group name", func() {
resp, err := client.Qualitygates.RemoveGroup(&sonargo.QualitygatesRemoveGroupOption{
GateName: "some-gate",
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
It("should add user permission to quality gate", func() {
gateName := helpers.UniqueResourceName("qg-usr")

_, _, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
Name: gateName,
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualitygate", gateName, func() error {
_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
Name: gateName,
})
return err
})

resp, err := client.Qualitygates.AddUser(&sonargo.QualitygatesAddUserOption{
GateName: gateName,
Login:    "admin",
})
Expect(err).NotTo(HaveOccurred())
Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
})

Context("parameter validation", func() {
It("should fail with nil options", func() {
resp, err := client.Qualitygates.AddUser(nil)
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with missing gate name", func() {
resp, err := client.Qualitygates.AddUser(&sonargo.QualitygatesAddUserOption{
Login: "admin",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with missing login", func() {
resp, err := client.Qualitygates.AddUser(&sonargo.QualitygatesAddUserOption{
GateName: "some-gate",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})
})
})

Describe("SearchUsers", func() {
It("should search users for a quality gate", func() {
gateName := helpers.UniqueResourceName("qg-susr")

_, _, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
Name: gateName,
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualitygate", gateName, func() error {
_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
Name: gateName,
})
return err
})

result, resp, err := client.Qualitygates.SearchUsers(&sonargo.QualitygatesSearchUsersOption{
GateName: gateName,
Selected: "all",
})
Expect(err).NotTo(HaveOccurred())
Expect(resp.StatusCode).To(Equal(http.StatusOK))
Expect(result).NotTo(BeNil())
})

Context("parameter validation", func() {
It("should fail with nil options", func() {
result, resp, err := client.Qualitygates.SearchUsers(nil)
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
Expect(result).To(BeNil())
})

It("should fail with missing gate name", func() {
result, resp, err := client.Qualitygates.SearchUsers(&sonargo.QualitygatesSearchUsersOption{})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
Expect(result).To(BeNil())
})
})
})

Describe("RemoveUser", func() {
It("should remove user permission from quality gate", func() {
gateName := helpers.UniqueResourceName("qg-rmusr")

_, _, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
Name: gateName,
})
Expect(err).NotTo(HaveOccurred())

cleanup.RegisterCleanup("qualitygate", gateName, func() error {
_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
Name: gateName,
})
return err
})

// Add user first
_, err = client.Qualitygates.AddUser(&sonargo.QualitygatesAddUserOption{
GateName: gateName,
Login:    "admin",
})
Expect(err).NotTo(HaveOccurred())

// Remove user
resp, err := client.Qualitygates.RemoveUser(&sonargo.QualitygatesRemoveUserOption{
GateName: gateName,
Login:    "admin",
})
Expect(err).NotTo(HaveOccurred())
Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
})

Context("parameter validation", func() {
It("should fail with nil options", func() {
resp, err := client.Qualitygates.RemoveUser(nil)
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with missing gate name", func() {
resp, err := client.Qualitygates.RemoveUser(&sonargo.QualitygatesRemoveUserOption{
Login: "admin",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})

It("should fail with missing login", func() {
resp, err := client.Qualitygates.RemoveUser(&sonargo.QualitygatesRemoveUserOption{
GateName: "some-gate",
})
Expect(err).To(HaveOccurred())
Expect(resp).To(BeNil())
})
})
})

// =========================================================================
// Qualitygates Lifecycle
// =========================================================================
Describe("Qualitygates Lifecycle", func() {
It("should complete full quality gate lifecycle", func() {
gateName := helpers.UniqueResourceName("qg-lifecycle")
projectKey := helpers.UniqueResourceName("proj-lifecycle")

// Step 1: Create quality gate
createResult, _, err := client.Qualitygates.Create(&sonargo.QualitygatesCreateOption{
Name: gateName,
})
Expect(err).NotTo(HaveOccurred())
Expect(createResult.Name).To(Equal(gateName))

cleanup.RegisterCleanup("qualitygate", gateName, func() error {
_, err := client.Qualitygates.Destroy(&sonargo.QualitygatesDestroyOption{
Name: gateName,
})
return err
})

// Step 2: Show the quality gate
showResult, _, err := client.Qualitygates.Show(&sonargo.QualitygatesShowOption{
Name: gateName,
})
Expect(err).NotTo(HaveOccurred())
Expect(showResult).NotTo(BeNil())

// Step 3: Create condition
condResult, _, err := client.Qualitygates.CreateCondition(&sonargo.QualitygatesCreateConditionOption{
GateName: gateName,
Metric:   "coverage",
Op:       "LT",
Error:    "80",
})
Expect(err).NotTo(HaveOccurred())
conditionID := condResult.ID

// Step 4: Update condition
_, err = client.Qualitygates.UpdateCondition(&sonargo.QualitygatesUpdateConditionOption{
ID:     conditionID,
Metric: "coverage",
Op:     "LT",
Error:  "85",
})
Expect(err).NotTo(HaveOccurred())

// Step 5: Create project and associate
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

_, err = client.Qualitygates.Select(&sonargo.QualitygatesSelectOption{
GateName:   gateName,
ProjectKey: projectKey,
})
Expect(err).NotTo(HaveOccurred())

// Step 6: Add user permission
_, err = client.Qualitygates.AddUser(&sonargo.QualitygatesAddUserOption{
GateName: gateName,
Login:    "admin",
})
Expect(err).NotTo(HaveOccurred())

// Step 7: Search users
usersResult, _, err := client.Qualitygates.SearchUsers(&sonargo.QualitygatesSearchUsersOption{
GateName: gateName,
Selected: "selected",
})
Expect(err).NotTo(HaveOccurred())
Expect(usersResult.Users).To(HaveLen(1))

// Step 8: Remove user permission
_, err = client.Qualitygates.RemoveUser(&sonargo.QualitygatesRemoveUserOption{
GateName: gateName,
Login:    "admin",
})
Expect(err).NotTo(HaveOccurred())

// Step 9: Deselect project
_, err = client.Qualitygates.Deselect(&sonargo.QualitygatesDeselectOption{
ProjectKey: projectKey,
})
Expect(err).NotTo(HaveOccurred())

// Step 10: Delete condition
_, err = client.Qualitygates.DeleteCondition(&sonargo.QualitygatesDeleteConditionOption{
ID: conditionID,
})
Expect(err).NotTo(HaveOccurred())

// Step 11: Verify condition was deleted
showResult, _, err = client.Qualitygates.Show(&sonargo.QualitygatesShowOption{
Name: gateName,
})
Expect(err).NotTo(HaveOccurred())
for _, cond := range showResult.Conditions {
Expect(cond.ID).NotTo(Equal(conditionID), "Deleted condition still exists")
}
})
})
})
