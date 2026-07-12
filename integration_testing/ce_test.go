package integration_testing_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

var _ = Describe("Ce Service", Ordered, func() {
	var (
		client         *sonar.Client
		cleanupManager *helpers.CleanupManager
		testProject    *sonar.ProjectsCreate
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())

		cleanupManager = helpers.NewCleanupManager(client)

		// Create a test project for CE operations
		projectName := helpers.UniqueResourceName("ce-test-project")
		testProject, _, err = client.Projects.Create(context.Background(), &sonar.ProjectsCreateOptions{
			Name:    projectName,
			Project: projectName,
		})
		Expect(err).NotTo(HaveOccurred())
		cleanupManager.RegisterCleanup("project", testProject.Project.Key, func() error {
			_, err := client.Projects.Delete(context.Background(), &sonar.ProjectsDeleteOptions{
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
				result, resp, err := client.Ce.Activity(context.Background(), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should list CE tasks with pagination", func() {
				result, resp, err := client.Ce.Activity(context.Background(), &sonar.CeActivityOptions{
					CePaginationArgs: sonar.CePaginationArgs{
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
				result, resp, err := client.Ce.Activity(context.Background(), &sonar.CeActivityOptions{
					Component: testProject.Project.Key,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should list CE tasks filtered by status", func() {
				result, resp, err := client.Ce.Activity(context.Background(), &sonar.CeActivityOptions{
					Statuses: []string{sonar.TaskStatusSuccess},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should list CE tasks filtered by type", func() {
				result, resp, err := client.Ce.Activity(context.Background(), &sonar.CeActivityOptions{
					Type: sonar.TaskTypeReport,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should list CE tasks with onlyCurrents filter", func() {
				result, resp, err := client.Ce.Activity(context.Background(), &sonar.CeActivityOptions{
					OnlyCurrents: true,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should search CE tasks by query", func() {
				result, resp, err := client.Ce.Activity(context.Background(), &sonar.CeActivityOptions{
					Query: testProject.Project.Key,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("Parameter Validation", func() {
			It("should fail with invalid status", func() {
				_, _, err := client.Ce.Activity(context.Background(), &sonar.CeActivityOptions{
					Statuses: []string{"INVALID_STATUS"},
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with invalid type", func() {
				_, _, err := client.Ce.Activity(context.Background(), &sonar.CeActivityOptions{
					Type: "INVALID_TYPE",
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with invalid page size", func() {
				_, _, err := client.Ce.Activity(context.Background(), &sonar.CeActivityOptions{
					CePaginationArgs: sonar.CePaginationArgs{
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
				result, resp, err := client.Ce.ActivityStatus(context.Background(), nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should get activity status for specific component", func() {
				result, resp, err := client.Ce.ActivityStatus(context.Background(), &sonar.CeActivityStatusOptions{
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
				result, resp, err := client.Ce.AnalysisStatus(context.Background(), &sonar.CeAnalysisStatusOptions{
					Component: testProject.Project.Key,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("Parameter Validation", func() {
			It("should fail with missing component", func() {
				_, _, err := client.Ce.AnalysisStatus(context.Background(), &sonar.CeAnalysisStatusOptions{})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, _, err := client.Ce.AnalysisStatus(context.Background(), nil)
				Expect(err).To(HaveOccurred())
			})

			It("should fail with non-existent component", func() {
				_, resp, err := client.Ce.AnalysisStatus(context.Background(), &sonar.CeAnalysisStatusOptions{
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
				_, err := client.Ce.Cancel(context.Background(), &sonar.CeCancelOptions{})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, err := client.Ce.Cancel(context.Background(), nil)
				Expect(err).To(HaveOccurred())
			})

			It("should handle non-existent task ID gracefully", func() {
				// SonarQube returns 204 even for non-existent task IDs
				resp, err := client.Ce.Cancel(context.Background(), &sonar.CeCancelOptions{
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
				resp, err := client.Ce.CancelAll(context.Background(), )
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
				result, resp, err := client.Ce.Component(context.Background(), &sonar.CeComponentOptions{
					Component: testProject.Project.Key,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("Parameter Validation", func() {
			It("should fail with missing component", func() {
				_, _, err := client.Ce.Component(context.Background(), &sonar.CeComponentOptions{})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, _, err := client.Ce.Component(context.Background(), nil)
				Expect(err).To(HaveOccurred())
			})

			It("should fail with non-existent component", func() {
				_, resp, err := client.Ce.Component(context.Background(), &sonar.CeComponentOptions{
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
				_, err := client.Ce.DismissAnalysisWarning(context.Background(), &sonar.CeDismissAnalysisWarningOptions{
					Warning: "some-warning",
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with missing warning", func() {
				_, err := client.Ce.DismissAnalysisWarning(context.Background(), &sonar.CeDismissAnalysisWarningOptions{
					Component: testProject.Project.Key,
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, err := client.Ce.DismissAnalysisWarning(context.Background(), nil)
				Expect(err).To(HaveOccurred())
			})

			It("should fail with non-existent warning", func() {
				resp, err := client.Ce.DismissAnalysisWarning(context.Background(), &sonar.CeDismissAnalysisWarningOptions{
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
				result, resp, err := client.Ce.IndexationStatus(context.Background(), )
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
				result, resp, err := client.Ce.Info(context.Background(), )
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
				resp, err := client.Ce.Pause(context.Background(), )
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(SatisfyAny(Equal(http.StatusOK), Equal(http.StatusNoContent)))

				// Resume CE workers
				resp, err = client.Ce.Resume(context.Background(), )
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
				_, _, err := client.Ce.Submit(context.Background(), &sonar.CeSubmitOptions{
					Report: "dummy-report",
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with missing report", func() {
				_, _, err := client.Ce.Submit(context.Background(), &sonar.CeSubmitOptions{
					ProjectKey: testProject.Project.Key,
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, _, err := client.Ce.Submit(context.Background(), nil)
				Expect(err).To(HaveOccurred())
			})

			It("should fail with project key too long", func() {
				longKey := string(make([]byte, 500))
				_, _, err := client.Ce.Submit(context.Background(), &sonar.CeSubmitOptions{
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
				activity, _, err := client.Ce.Activity(context.Background(), &sonar.CeActivityOptions{
					CePaginationArgs: sonar.CePaginationArgs{
						PageSize: 1,
					},
				})
				Expect(err).NotTo(HaveOccurred())

				if len(activity.Tasks) > 0 {
					taskID := activity.Tasks[0].ID
					result, resp, err := client.Ce.Task(context.Background(), &sonar.CeTaskOptions{
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
				activity, _, err := client.Ce.Activity(context.Background(), &sonar.CeActivityOptions{
					CePaginationArgs: sonar.CePaginationArgs{
						PageSize: 1,
					},
				})
				Expect(err).NotTo(HaveOccurred())

				if len(activity.Tasks) > 0 {
					taskID := activity.Tasks[0].ID
					result, resp, err := client.Ce.Task(context.Background(), &sonar.CeTaskOptions{
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
				_, _, err := client.Ce.Task(context.Background(), &sonar.CeTaskOptions{})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, _, err := client.Ce.Task(context.Background(), nil)
				Expect(err).To(HaveOccurred())
			})

			It("should fail with invalid additional fields", func() {
				_, _, err := client.Ce.Task(context.Background(), &sonar.CeTaskOptions{
					ID:               "some-task-id",
					AdditionalFields: []string{"invalid_field"},
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with non-existent task ID", func() {
				_, resp, err := client.Ce.Task(context.Background(), &sonar.CeTaskOptions{
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
				result, resp, err := client.Ce.TaskTypes(context.Background(), )
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
				result, resp, err := client.Ce.WorkerCount(context.Background(), )
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Value).To(BeNumerically(">=", 1))
			})
		})
	})
})
