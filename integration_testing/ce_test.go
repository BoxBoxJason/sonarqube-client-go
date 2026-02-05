package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Ce Service", Ordered, func() {
	var (
		client         *sonargo.Client
		cleanupManager *helpers.CleanupManager
		testProject    *sonargo.ProjectsCreate
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())

		cleanupManager = helpers.NewCleanupManager(client)

		// Create a test project for CE operations
		projectName := helpers.UniqueResourceName("ce-test-project")
		testProject, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
			Name:    projectName,
			Project: projectName,
		})
		Expect(err).NotTo(HaveOccurred())
		cleanupManager.RegisterCleanup("project", testProject.Project.Key, func() error {
			_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
				Project: testProject.Project.Key,
			})
			return err
		})
	})

	AfterAll(func() {
		errors := cleanupManager.Cleanup()
		for _, err := range errors {
			GinkgoWriter.Printf("Cleanup error: %v\n", err)
		}
	})

	// =========================================================================
	// Activity
	// =========================================================================
	Describe("Activity", func() {
		Context("Functional Tests", func() {
			It("should list CE tasks with no options", func() {
				result, resp, err := client.Ce.Activity(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should list CE tasks with pagination", func() {
				result, resp, err := client.Ce.Activity(&sonargo.CeActivityOption{
					CePaginationArgs: sonargo.CePaginationArgs{
						Page:     1,
						PageSize: 10,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Paging.PageSize).To(Equal(int64(10)))
			})

			It("should list CE tasks filtered by component", func() {
				result, resp, err := client.Ce.Activity(&sonargo.CeActivityOption{
					Component: testProject.Project.Key,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should list CE tasks filtered by status", func() {
				result, resp, err := client.Ce.Activity(&sonargo.CeActivityOption{
					Statuses: []string{"SUCCESS"},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should list CE tasks filtered by type", func() {
				result, resp, err := client.Ce.Activity(&sonargo.CeActivityOption{
					Type: "REPORT",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should list CE tasks with onlyCurrents filter", func() {
				result, resp, err := client.Ce.Activity(&sonargo.CeActivityOption{
					OnlyCurrents: true,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should search CE tasks by query", func() {
				result, resp, err := client.Ce.Activity(&sonargo.CeActivityOption{
					Q: testProject.Project.Key,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("Parameter Validation", func() {
			It("should fail with invalid status", func() {
				_, _, err := client.Ce.Activity(&sonargo.CeActivityOption{
					Statuses: []string{"INVALID_STATUS"},
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with invalid type", func() {
				_, _, err := client.Ce.Activity(&sonargo.CeActivityOption{
					Type: "INVALID_TYPE",
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with invalid page size", func() {
				_, _, err := client.Ce.Activity(&sonargo.CeActivityOption{
					CePaginationArgs: sonargo.CePaginationArgs{
						PageSize: 10000,
					},
				})
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// =========================================================================
	// ActivityStatus
	// =========================================================================
	Describe("ActivityStatus", func() {
		Context("Functional Tests", func() {
			It("should get activity status with no options", func() {
				result, resp, err := client.Ce.ActivityStatus(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should get activity status for specific component", func() {
				result, resp, err := client.Ce.ActivityStatus(&sonargo.CeActivityStatusOption{
					Component: testProject.Project.Key,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})
	})

	// =========================================================================
	// AnalysisStatus
	// =========================================================================
	Describe("AnalysisStatus", func() {
		Context("Functional Tests", func() {
			It("should get analysis status for component", func() {
				result, resp, err := client.Ce.AnalysisStatus(&sonargo.CeAnalysisStatusOption{
					Component: testProject.Project.Key,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("Parameter Validation", func() {
			It("should fail with missing component", func() {
				_, _, err := client.Ce.AnalysisStatus(&sonargo.CeAnalysisStatusOption{})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, _, err := client.Ce.AnalysisStatus(nil)
				Expect(err).To(HaveOccurred())
			})

			It("should fail with non-existent component", func() {
				_, resp, err := client.Ce.AnalysisStatus(&sonargo.CeAnalysisStatusOption{
					Component: "non-existent-project-12345",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	// =========================================================================
	// Cancel
	// =========================================================================
	Describe("Cancel", func() {
		Context("Parameter Validation", func() {
			It("should fail with missing task ID", func() {
				_, err := client.Ce.Cancel(&sonargo.CeCancelOption{})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, err := client.Ce.Cancel(nil)
				Expect(err).To(HaveOccurred())
			})

			It("should handle non-existent task ID gracefully", func() {
				// SonarQube returns 204 even for non-existent task IDs
				resp, err := client.Ce.Cancel(&sonargo.CeCancelOption{
					ID: "non-existent-task-id",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})
		})
	})

	// =========================================================================
	// CancelAll
	// =========================================================================
	Describe("CancelAll", func() {
		Context("Functional Tests", func() {
			It("should cancel all pending tasks", func() {
				resp, err := client.Ce.CancelAll()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(SatisfyAny(Equal(http.StatusOK), Equal(http.StatusNoContent)))
			})
		})
	})

	// =========================================================================
	// Component
	// =========================================================================
	Describe("Component", func() {
		Context("Functional Tests", func() {
			It("should get component CE status", func() {
				result, resp, err := client.Ce.Component(&sonargo.CeComponentOption{
					Component: testProject.Project.Key,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("Parameter Validation", func() {
			It("should fail with missing component", func() {
				_, _, err := client.Ce.Component(&sonargo.CeComponentOption{})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, _, err := client.Ce.Component(nil)
				Expect(err).To(HaveOccurred())
			})

			It("should fail with non-existent component", func() {
				_, resp, err := client.Ce.Component(&sonargo.CeComponentOption{
					Component: "non-existent-project-12345",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	// =========================================================================
	// DismissAnalysisWarning
	// =========================================================================
	Describe("DismissAnalysisWarning", func() {
		Context("Parameter Validation", func() {
			It("should fail with missing component", func() {
				_, err := client.Ce.DismissAnalysisWarning(&sonargo.CeDismissAnalysisWarningOption{
					Warning: "some-warning",
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with missing warning", func() {
				_, err := client.Ce.DismissAnalysisWarning(&sonargo.CeDismissAnalysisWarningOption{
					Component: testProject.Project.Key,
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, err := client.Ce.DismissAnalysisWarning(nil)
				Expect(err).To(HaveOccurred())
			})

			It("should fail with non-existent warning", func() {
				resp, err := client.Ce.DismissAnalysisWarning(&sonargo.CeDismissAnalysisWarningOption{
					Component: testProject.Project.Key,
					Warning:   "non-existent-warning",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
			})
		})
	})

	// =========================================================================
	// IndexationStatus
	// =========================================================================
	Describe("IndexationStatus", func() {
		Context("Functional Tests", func() {
			It("should get indexation status", func() {
				result, resp, err := client.Ce.IndexationStatus()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})
	})

	// =========================================================================
	// Info
	// =========================================================================
	Describe("Info", func() {
		Context("Functional Tests", func() {
			It("should get CE info", func() {
				result, resp, err := client.Ce.Info()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})
	})

	// =========================================================================
	// Pause and Resume
	// =========================================================================
	Describe("Pause and Resume", func() {
		Context("Functional Tests", func() {
			It("should pause and resume CE workers", func() {
				// Pause CE workers
				resp, err := client.Ce.Pause()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(SatisfyAny(Equal(http.StatusOK), Equal(http.StatusNoContent)))

				// Resume CE workers
				resp, err = client.Ce.Resume()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(SatisfyAny(Equal(http.StatusOK), Equal(http.StatusNoContent)))
			})
		})
	})

	// =========================================================================
	// Submit
	// =========================================================================
	Describe("Submit", func() {
		Context("Parameter Validation", func() {
			It("should fail with missing project key", func() {
				_, _, err := client.Ce.Submit(&sonargo.CeSubmitOption{
					Report: "dummy-report",
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with missing report", func() {
				_, _, err := client.Ce.Submit(&sonargo.CeSubmitOption{
					ProjectKey: testProject.Project.Key,
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, _, err := client.Ce.Submit(nil)
				Expect(err).To(HaveOccurred())
			})

			It("should fail with project key too long", func() {
				longKey := string(make([]byte, 500))
				_, _, err := client.Ce.Submit(&sonargo.CeSubmitOption{
					ProjectKey: longKey,
					Report:     "dummy-report",
				})
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// =========================================================================
	// Task
	// =========================================================================
	Describe("Task", func() {
		Context("Functional Tests", func() {
			It("should get task details when tasks exist", func() {
				// First get any existing task from activity
				activity, _, err := client.Ce.Activity(&sonargo.CeActivityOption{
					CePaginationArgs: sonargo.CePaginationArgs{
						PageSize: 1,
					},
				})
				Expect(err).NotTo(HaveOccurred())

				if len(activity.Tasks) > 0 {
					taskID := activity.Tasks[0].ID
					result, resp, err := client.Ce.Task(&sonargo.CeTaskOption{
						ID: taskID,
					})
					Expect(err).NotTo(HaveOccurred())
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					Expect(result).NotTo(BeNil())
					Expect(result.Task.ID).To(Equal(taskID))
				} else {
					Skip("No tasks available to test")
				}
			})

			It("should get task details with additional fields", func() {
				// First get any existing task from activity
				activity, _, err := client.Ce.Activity(&sonargo.CeActivityOption{
					CePaginationArgs: sonargo.CePaginationArgs{
						PageSize: 1,
					},
				})
				Expect(err).NotTo(HaveOccurred())

				if len(activity.Tasks) > 0 {
					taskID := activity.Tasks[0].ID
					result, resp, err := client.Ce.Task(&sonargo.CeTaskOption{
						ID:               taskID,
						AdditionalFields: []string{"warnings"},
					})
					Expect(err).NotTo(HaveOccurred())
					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					Expect(result).NotTo(BeNil())
				} else {
					Skip("No tasks available to test")
				}
			})
		})

		Context("Parameter Validation", func() {
			It("should fail with missing task ID", func() {
				_, _, err := client.Ce.Task(&sonargo.CeTaskOption{})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, _, err := client.Ce.Task(nil)
				Expect(err).To(HaveOccurred())
			})

			It("should fail with invalid additional fields", func() {
				_, _, err := client.Ce.Task(&sonargo.CeTaskOption{
					ID:               "some-task-id",
					AdditionalFields: []string{"invalid_field"},
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with non-existent task ID", func() {
				_, resp, err := client.Ce.Task(&sonargo.CeTaskOption{
					ID: "non-existent-task-id",
				})
				Expect(err).To(HaveOccurred())
				Expect(resp).NotTo(BeNil())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	})

	// =========================================================================
	// TaskTypes
	// =========================================================================
	Describe("TaskTypes", func() {
		Context("Functional Tests", func() {
			It("should list available task types", func() {
				result, resp, err := client.Ce.TaskTypes()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.TaskTypes).NotTo(BeEmpty())
			})
		})
	})

	// =========================================================================
	// WorkerCount
	// =========================================================================
	Describe("WorkerCount", func() {
		Context("Functional Tests", func() {
			It("should get CE worker count", func() {
				result, resp, err := client.Ce.WorkerCount()
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Value).To(BeNumerically(">=", 1))
			})
		})
	})
})
