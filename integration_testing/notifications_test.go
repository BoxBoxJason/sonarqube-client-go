package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("Notifications Service", Ordered, func() {
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

		// Create a test project for project-scoped notifications
		projectKey = helpers.UniqueResourceName("notif")
		_, _, err = client.Projects.Create(&sonar.ProjectsCreateOption{
			Name:    "Notifications Test Project",
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
	// List
	// =========================================================================
	Describe("List", func() {
		Context("Valid Requests", func() {
			It("should list notifications with nil options", func() {
				result, resp, err := client.Notifications.List(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Channels).NotTo(BeNil())
				Expect(result.GlobalTypes).NotTo(BeNil())
			})

			It("should list notifications with empty options", func() {
				result, resp, err := client.Notifications.List(&sonar.NotificationsListOption{})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})
	})

	// =========================================================================
	// Add
	// =========================================================================
	Describe("Add", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Notifications.Add(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required type", func() {
				resp, err := client.Notifications.Add(&sonar.NotificationsAddOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Type"))
				Expect(resp).To(BeNil())
			})
		})

		Context("Valid Requests", func() {
			It("should add a global notification", func() {
				// First get the list to know valid types
				listResult, _, err := client.Notifications.List(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(listResult).NotTo(BeNil())
				Expect(listResult.GlobalTypes).NotTo(BeEmpty())

				notificationType := listResult.GlobalTypes[0]

				resp, err := client.Notifications.Add(&sonar.NotificationsAddOption{
					Type: notificationType,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// Clean up
				_, err = client.Notifications.Remove(&sonar.NotificationsRemoveOption{
					Type: notificationType,
				})
				if err != nil {
					GinkgoWriter.Printf("Cleanup error: failed to remove notification: %v\n", err)
				}
			})

			It("should add a project-scoped notification", func() {
				// First get the list to know valid types
				listResult, _, err := client.Notifications.List(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(listResult).NotTo(BeNil())
				Expect(listResult.PerProjectTypes).NotTo(BeEmpty())

				notificationType := listResult.PerProjectTypes[0]

				resp, err := client.Notifications.Add(&sonar.NotificationsAddOption{
					Type:    notificationType,
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// Clean up
				_, err = client.Notifications.Remove(&sonar.NotificationsRemoveOption{
					Type:    notificationType,
					Project: projectKey,
				})
				if err != nil {
					GinkgoWriter.Printf("Cleanup error: failed to remove notification: %v\n", err)
				}
			})
		})

		Context("Invalid Type", func() {
			It("should fail for invalid notification type", func() {
				resp, err := client.Notifications.Add(&sonar.NotificationsAddOption{
					Type: "invalid-type",
				})
				Expect(err).To(HaveOccurred())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// Remove
	// =========================================================================
	Describe("Remove", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Notifications.Remove(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required type", func() {
				resp, err := client.Notifications.Remove(&sonar.NotificationsRemoveOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Type"))
				Expect(resp).To(BeNil())
			})
		})

		Context("Valid Requests", func() {
			It("should remove a notification", func() {
				// First get valid types
				listResult, _, err := client.Notifications.List(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(listResult).NotTo(BeNil())
				Expect(listResult.GlobalTypes).NotTo(BeEmpty())

				notificationType := listResult.GlobalTypes[0]

				// Add the notification first
				_, err = client.Notifications.Add(&sonar.NotificationsAddOption{
					Type: notificationType,
				})
				Expect(err).NotTo(HaveOccurred())

				// Now remove it
				resp, err := client.Notifications.Remove(&sonar.NotificationsRemoveOption{
					Type: notificationType,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})
		})

		Context("Non-Existent Notification", func() {
			It("should fail for notification that doesn't exist", func() {
				// Use a valid type but one that's not subscribed
				listResult, _, err := client.Notifications.List(nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(listResult).NotTo(BeNil())
				Expect(listResult.GlobalTypes).NotTo(BeEmpty())

				// Pick a type that's likely not subscribed
				var unsubscribedType string
				for _, t := range listResult.GlobalTypes {
					found := false
					for _, n := range listResult.Notifications {
						if n.Type == t && n.Project == "" {
							found = true
							break
						}
					}
					if !found {
						unsubscribedType = t
						break
					}
				}

				if unsubscribedType != "" {
					resp, err := client.Notifications.Remove(&sonar.NotificationsRemoveOption{
						Type: unsubscribedType,
					})
					Expect(err).To(HaveOccurred())
					if resp != nil {
						Expect(resp.StatusCode).To(BeNumerically(">=", 400))
					}
				}
			})
		})
	})

	// =========================================================================
	// Full Workflow
	// =========================================================================
	Describe("Full Workflow", func() {
		It("should add, verify, and remove a notification", func() {
			// Get valid types
			listResult, _, err := client.Notifications.List(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(listResult).NotTo(BeNil())
			Expect(listResult.GlobalTypes).NotTo(BeEmpty())

			notificationType := listResult.GlobalTypes[0]

			// Remove if already exists (clean state)
			_, err = client.Notifications.Remove(&sonar.NotificationsRemoveOption{
				Type: notificationType,
			})
			if err != nil {
				GinkgoWriter.Printf("Cleanup error: failed to remove notification for clean state: %v\n", err)
			}

			// Add the notification
			resp, err := client.Notifications.Add(&sonar.NotificationsAddOption{
				Type: notificationType,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// List and verify it's there
			listResult, resp, err = client.Notifications.List(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			found := false
			for _, n := range listResult.Notifications {
				if n.Type == notificationType && n.Project == "" {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue(), "Notification should be in the list")

			// Remove the notification
			resp, err = client.Notifications.Remove(&sonar.NotificationsRemoveOption{
				Type: notificationType,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify it's gone
			listResult, _, err = client.Notifications.List(nil)
			Expect(err).NotTo(HaveOccurred())

			found = false
			for _, n := range listResult.Notifications {
				if n.Type == notificationType && n.Project == "" {
					found = true
					break
				}
			}
			Expect(found).To(BeFalse(), "Notification should not be in the list")
		})
	})
})
