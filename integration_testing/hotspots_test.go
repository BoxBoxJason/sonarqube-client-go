package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("Hotspots Service", Ordered, func() {
	var (
		client     *sonar.Client
		cleanup    *helpers.CleanupManager
		projectKey string
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
		cleanup = helpers.NewCleanupManager(client)

		// Create a test project for hotspots-related operations
		projectKey = helpers.UniqueResourceName("hot")
		_, _, err = client.Projects.Create(&sonar.ProjectsCreateOption{
			Name:    "Hotspots Test Project",
			Project: projectKey,
		})
		Expect(err).NotTo(HaveOccurred())

		cleanup.RegisterCleanup("project", projectKey, func() error {
			_, err := client.Projects.Delete(&sonar.ProjectsDeleteOption{
				Project: projectKey,
			})
			return err
		})
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
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Hotspots.Search(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without project or hotspots parameter", func() {
				result, resp, err := client.Hotspots.Search(&sonar.HotspotsSearchOption{})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Valid Requests", func() {
			It("should search hotspots for a project", func() {
				result, resp, err := client.Hotspots.Search(&sonar.HotspotsSearchOption{
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should search hotspots with pagination", func() {
				result, resp, err := client.Hotspots.Search(&sonar.HotspotsSearchOption{
					Project: projectKey,
					PaginationArgs: sonar.PaginationArgs{
						PageSize: 10,
						Page:     1,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should search hotspots with status filter", func() {
				result, resp, err := client.Hotspots.Search(&sonar.HotspotsSearchOption{
					Project: projectKey,
					Status:  "TO_REVIEW",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should search hotspots in new code period", func() {
				result, resp, err := client.Hotspots.Search(&sonar.HotspotsSearchOption{
					Project:         projectKey,
					InNewCodePeriod: true,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("Non-Existent Project", func() {
			It("should fail for non-existent project", func() {
				result, resp, err := client.Hotspots.Search(&sonar.HotspotsSearchOption{
					Project: "non-existent-project",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// List
	// =========================================================================
	Describe("List", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Hotspots.List(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required project", func() {
				result, resp, err := client.Hotspots.List(&sonar.HotspotsListOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Project"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Valid Requests", func() {
			It("should list hotspots for a project", func() {
				result, resp, err := client.Hotspots.List(&sonar.HotspotsListOption{
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should list hotspots with pagination", func() {
				result, resp, err := client.Hotspots.List(&sonar.HotspotsListOption{
					Project: projectKey,
					PaginationArgs: sonar.PaginationArgs{
						PageSize: 10,
						Page:     1,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should list hotspots with status filter", func() {
				result, resp, err := client.Hotspots.List(&sonar.HotspotsListOption{
					Project: projectKey,
					Status:  "TO_REVIEW",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("Non-Existent Project", func() {
			It("should fail for non-existent project", func() {
				result, resp, err := client.Hotspots.List(&sonar.HotspotsListOption{
					Project: "non-existent-project",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// Show
	// =========================================================================
	Describe("Show", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Hotspots.Show(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required hotspot key", func() {
				result, resp, err := client.Hotspots.Show(&sonar.HotspotsShowOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Hotspot"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent Hotspot", func() {
			It("should fail with non-existent hotspot key", func() {
				result, resp, err := client.Hotspots.Show(&sonar.HotspotsShowOption{
					Hotspot: "AXxxxxxxxxxxxxxxxxxx",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// Pull
	// =========================================================================
	Describe("Pull", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Hotspots.Pull(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required project key", func() {
				result, resp, err := client.Hotspots.Pull(&sonar.HotspotsPullOption{
					BranchName: "main",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ProjectKey"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required branch name", func() {
				result, resp, err := client.Hotspots.Pull(&sonar.HotspotsPullOption{
					ProjectKey: projectKey,
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("BranchName"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent Project", func() {
			It("should fail for non-existent project", func() {
				result, resp, err := client.Hotspots.Pull(&sonar.HotspotsPullOption{
					ProjectKey: "non-existent-project",
					BranchName: "main",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// AddComment (requires existing hotspot)
	// =========================================================================
	Describe("AddComment", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Hotspots.AddComment(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required hotspot key", func() {
				resp, err := client.Hotspots.AddComment(&sonar.HotspotsAddCommentOption{
					Comment: "Test comment",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Hotspot"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required comment", func() {
				resp, err := client.Hotspots.AddComment(&sonar.HotspotsAddCommentOption{
					Hotspot: "AXxxxxxxxxxxxxxxxxxx",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Comment"))
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent Hotspot", func() {
			It("should fail with non-existent hotspot key", func() {
				resp, err := client.Hotspots.AddComment(&sonar.HotspotsAddCommentOption{
					Hotspot: "AXxxxxxxxxxxxxxxxxxx",
					Comment: "Test comment",
				})
				Expect(err).To(HaveOccurred())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// Assign (requires existing hotspot)
	// =========================================================================
	Describe("Assign", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Hotspots.Assign(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required hotspot key", func() {
				resp, err := client.Hotspots.Assign(&sonar.HotspotsAssignOption{
					Assignee: "admin",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Hotspot"))
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent Hotspot", func() {
			It("should fail with non-existent hotspot key", func() {
				resp, err := client.Hotspots.Assign(&sonar.HotspotsAssignOption{
					Hotspot:  "AXxxxxxxxxxxxxxxxxxx",
					Assignee: "admin",
				})
				Expect(err).To(HaveOccurred())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// ChangeStatus (requires existing hotspot)
	// =========================================================================
	Describe("ChangeStatus", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Hotspots.ChangeStatus(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required hotspot key", func() {
				resp, err := client.Hotspots.ChangeStatus(&sonar.HotspotsChangeStatusOption{
					Status: "REVIEWED",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Hotspot"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required status", func() {
				resp, err := client.Hotspots.ChangeStatus(&sonar.HotspotsChangeStatusOption{
					Hotspot: "AXxxxxxxxxxxxxxxxxxx",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Status"))
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent Hotspot", func() {
			It("should fail with non-existent hotspot key", func() {
				resp, err := client.Hotspots.ChangeStatus(&sonar.HotspotsChangeStatusOption{
					Hotspot:    "AXxxxxxxxxxxxxxxxxxx",
					Status:     "REVIEWED",
					Resolution: "SAFE",
				})
				Expect(err).To(HaveOccurred())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// DeleteComment (requires existing comment)
	// =========================================================================
	Describe("DeleteComment", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Hotspots.DeleteComment(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required comment key", func() {
				resp, err := client.Hotspots.DeleteComment(&sonar.HotspotsDeleteCommentOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Comment"))
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent Comment", func() {
			It("should fail with non-existent comment key", func() {
				resp, err := client.Hotspots.DeleteComment(&sonar.HotspotsDeleteCommentOption{
					Comment: "AXxxxxxxxxxxxxxxxxxx",
				})
				Expect(err).To(HaveOccurred())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// EditComment (requires existing comment)
	// =========================================================================
	Describe("EditComment", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Hotspots.EditComment(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required comment key", func() {
				result, resp, err := client.Hotspots.EditComment(&sonar.HotspotsEditCommentOption{
					Text: "Updated comment",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Comment"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required text", func() {
				result, resp, err := client.Hotspots.EditComment(&sonar.HotspotsEditCommentOption{
					Comment: "AXxxxxxxxxxxxxxxxxxx",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Text"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent Comment", func() {
			It("should fail with non-existent comment key", func() {
				result, resp, err := client.Hotspots.EditComment(&sonar.HotspotsEditCommentOption{
					Comment: "AXxxxxxxxxxxxxxxxxxx",
					Text:    "Updated comment",
				})
				Expect(err).To(HaveOccurred())
				Expect(result).To(BeNil())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})
})
