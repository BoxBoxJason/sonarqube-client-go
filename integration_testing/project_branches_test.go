package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("ProjectBranches Service", Ordered, func() {
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
	// List
	// =========================================================================
	Describe("List", func() {
		var testProjectKey string

		BeforeEach(func() {
			testProjectKey = helpers.UniqueResourceName("proj-branch")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:    "Branch Test Project",
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

		It("should list branches of a project", func() {
			result, resp, err := client.ProjectBranches.List(&sonargo.ProjectBranchesListOption{
				Project: testProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result).NotTo(BeNil())
			// A newly created project has at least the main branch
			Expect(result.Branches).NotTo(BeEmpty())

			// Find and verify the main branch
			var mainBranch *sonargo.Branch
			for i := range result.Branches {
				if result.Branches[i].IsMain {
					mainBranch = &result.Branches[i]
					break
				}
			}
			Expect(mainBranch).NotTo(BeNil())
			Expect(mainBranch.IsMain).To(BeTrue())
			Expect(mainBranch.ExcludedFromPurge).To(BeTrue()) // Main branch is always protected
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				_, resp, err := client.ProjectBranches.List(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing project key", func() {
				_, resp, err := client.ProjectBranches.List(&sonargo.ProjectBranchesListOption{})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})

		Context("error cases", func() {
			It("should fail for non-existent project", func() {
				_, resp, err := client.ProjectBranches.List(&sonargo.ProjectBranchesListOption{
					Project: "non-existent-project-12345",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	// =========================================================================
	// Rename
	// =========================================================================
	Describe("Rename", func() {
		var testProjectKey string

		BeforeEach(func() {
			testProjectKey = helpers.UniqueResourceName("proj-rename")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "Rename Branch Test Project",
				Project:    testProjectKey,
				MainBranch: "main",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", testProjectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: testProjectKey,
				})
				return err
			})
		})

		It("should rename the main branch", func() {
			resp, err := client.ProjectBranches.Rename(&sonargo.ProjectBranchesRenameOption{
				Project: testProjectKey,
				Name:    "master",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify the branch was renamed
			result, _, err := client.ProjectBranches.List(&sonargo.ProjectBranchesListOption{
				Project: testProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			var mainBranch *sonargo.Branch
			for i := range result.Branches {
				if result.Branches[i].IsMain {
					mainBranch = &result.Branches[i]
					break
				}
			}
			Expect(mainBranch).NotTo(BeNil())
			Expect(mainBranch.Name).To(Equal("master"))
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.ProjectBranches.Rename(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing project key", func() {
				resp, err := client.ProjectBranches.Rename(&sonargo.ProjectBranchesRenameOption{
					Name: "new-main",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing name", func() {
				resp, err := client.ProjectBranches.Rename(&sonargo.ProjectBranchesRenameOption{
					Project: helpers.UniqueResourceName("proj"),
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})
	})

	// =========================================================================
	// Delete
	// =========================================================================
	Describe("Delete", func() {
		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.ProjectBranches.Delete(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing project key", func() {
				resp, err := client.ProjectBranches.Delete(&sonargo.ProjectBranchesDeleteOption{
					Branch: "feature-branch",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing branch name", func() {
				resp, err := client.ProjectBranches.Delete(&sonargo.ProjectBranchesDeleteOption{
					Project: helpers.UniqueResourceName("proj"),
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})

		Context("error cases", func() {
			It("should fail for non-existent project", func() {
				resp, err := client.ProjectBranches.Delete(&sonargo.ProjectBranchesDeleteOption{
					Project: "non-existent-project-12345",
					Branch:  "feature-branch",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})

			It("should fail when trying to delete main branch", func() {
				projectKey := helpers.UniqueResourceName("proj-del-main")

				_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
					Name:       "Delete Main Branch Test",
					Project:    projectKey,
					MainBranch: "main",
				})
				Expect(err).NotTo(HaveOccurred())

				cleanup.RegisterCleanup("project", projectKey, func() error {
					_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
						Project: projectKey,
					})
					return err
				})

				resp, err := client.ProjectBranches.Delete(&sonargo.ProjectBranchesDeleteOption{
					Project: projectKey,
					Branch:  "main",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})

			It("should fail for non-existent branch", func() {
				projectKey := helpers.UniqueResourceName("proj-del-noexist")

				_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
					Name:    "Delete Non-existent Branch Test",
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())

				cleanup.RegisterCleanup("project", projectKey, func() error {
					_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
						Project: projectKey,
					})
					return err
				})

				resp, err := client.ProjectBranches.Delete(&sonargo.ProjectBranchesDeleteOption{
					Project: projectKey,
					Branch:  "non-existent-branch",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	// =========================================================================
	// SetAutomaticDeletionProtection
	// =========================================================================
	Describe("SetAutomaticDeletionProtection", func() {
		var testProjectKey string

		BeforeEach(func() {
			testProjectKey = helpers.UniqueResourceName("proj-protect")

			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "Protection Test Project",
				Project:    testProjectKey,
				MainBranch: "main",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", testProjectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: testProjectKey,
				})
				return err
			})
		})

		It("should set automatic deletion protection on main branch", func() {
			// Main branch should already be protected, but we can still call the API
			resp, err := client.ProjectBranches.SetAutomaticDeletionProtection(&sonargo.ProjectBranchesSetAutomaticDeletionProtectionOption{
				Project: testProjectKey,
				Branch:  "main",
				Value:   true,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify protection is set
			result, _, err := client.ProjectBranches.List(&sonargo.ProjectBranchesListOption{
				Project: testProjectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			var mainBranch *sonargo.Branch
			for i := range result.Branches {
				if result.Branches[i].IsMain {
					mainBranch = &result.Branches[i]
					break
				}
			}
			Expect(mainBranch).NotTo(BeNil())
			Expect(mainBranch.ExcludedFromPurge).To(BeTrue())
		})

		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.ProjectBranches.SetAutomaticDeletionProtection(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing project key", func() {
				resp, err := client.ProjectBranches.SetAutomaticDeletionProtection(&sonargo.ProjectBranchesSetAutomaticDeletionProtectionOption{
					Branch: "main",
					Value:  true,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing branch name", func() {
				resp, err := client.ProjectBranches.SetAutomaticDeletionProtection(&sonargo.ProjectBranchesSetAutomaticDeletionProtectionOption{
					Project: helpers.UniqueResourceName("proj"),
					Value:   true,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})

		Context("error cases", func() {
			It("should fail for non-existent project", func() {
				resp, err := client.ProjectBranches.SetAutomaticDeletionProtection(&sonargo.ProjectBranchesSetAutomaticDeletionProtectionOption{
					Project: "non-existent-project-12345",
					Branch:  "main",
					Value:   true,
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	// =========================================================================
	// SetMain
	// =========================================================================
	Describe("SetMain", func() {
		Context("parameter validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.ProjectBranches.SetMain(nil)
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing project key", func() {
				resp, err := client.ProjectBranches.SetMain(&sonargo.ProjectBranchesSetMainOption{
					Branch: "new-main",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})

			It("should fail with missing branch name", func() {
				resp, err := client.ProjectBranches.SetMain(&sonargo.ProjectBranchesSetMainOption{
					Project: helpers.UniqueResourceName("proj"),
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).To(BeNil())
			})
		})

		Context("error cases", func() {
			It("should fail for non-existent project", func() {
				resp, err := client.ProjectBranches.SetMain(&sonargo.ProjectBranchesSetMainOption{
					Project: "non-existent-project-12345",
					Branch:  "develop",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})

			It("should fail for non-existent branch", func() {
				projectKey := helpers.UniqueResourceName("proj-setmain")

				_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
					Name:    "Set Main Test Project",
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())

				cleanup.RegisterCleanup("project", projectKey, func() error {
					_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
						Project: projectKey,
					})
					return err
				})

				resp, err := client.ProjectBranches.SetMain(&sonargo.ProjectBranchesSetMainOption{
					Project: projectKey,
					Branch:  "non-existent-branch",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	// =========================================================================
	// Branch Lifecycle
	// =========================================================================
	Describe("Branch Lifecycle", func() {
		It("should manage main branch name through rename", func() {
			projectKey := helpers.UniqueResourceName("proj-lifecycle")

			// Step 1: Create project with specific main branch name
			_, _, err := client.Projects.Create(&sonargo.ProjectsCreateOption{
				Name:       "Lifecycle Test Project",
				Project:    projectKey,
				MainBranch: "main",
			})
			Expect(err).NotTo(HaveOccurred())

			cleanup.RegisterCleanup("project", projectKey, func() error {
				_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
					Project: projectKey,
				})
				return err
			})

			// Step 2: List branches
			result, _, err := client.ProjectBranches.List(&sonargo.ProjectBranchesListOption{
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Branches).NotTo(BeEmpty())

			var mainBranch *sonargo.Branch
			for i := range result.Branches {
				if result.Branches[i].IsMain {
					mainBranch = &result.Branches[i]
					break
				}
			}
			Expect(mainBranch).NotTo(BeNil())
			Expect(mainBranch.Name).To(Equal("main"))

			// Step 3: Rename main branch
			_, err = client.ProjectBranches.Rename(&sonargo.ProjectBranchesRenameOption{
				Project: projectKey,
				Name:    "master",
			})
			Expect(err).NotTo(HaveOccurred())

			// Step 4: Verify rename
			result, _, err = client.ProjectBranches.List(&sonargo.ProjectBranchesListOption{
				Project: projectKey,
			})
			Expect(err).NotTo(HaveOccurred())

			mainBranch = nil
			for i := range result.Branches {
				if result.Branches[i].IsMain {
					mainBranch = &result.Branches[i]
					break
				}
			}
			Expect(mainBranch).NotTo(BeNil())
			Expect(mainBranch.Name).To(Equal("master"))

			// Step 5: Verify protection on main branch
			Expect(mainBranch.ExcludedFromPurge).To(BeTrue())
		})
	})
})
