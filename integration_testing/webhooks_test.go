package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("Webhooks Service", Ordered, func() {
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

		// Create a test project for project-scoped webhooks
		projectKey = helpers.UniqueResourceName("webhook")
		_, _, err = client.Projects.Create(&sonar.ProjectsCreateOption{
			Name:    "Webhooks Test Project",
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
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Webhooks.List(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Valid Requests", func() {
			It("should list global webhooks", func() {
				result, resp, err := client.Webhooks.List(&sonar.WebhooksListOption{})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should list project webhooks", func() {
				result, resp, err := client.Webhooks.List(&sonar.WebhooksListOption{
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})
	})

	// =========================================================================
	// Create
	// =========================================================================
	Describe("Create", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Webhooks.Create(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required name", func() {
				result, resp, err := client.Webhooks.Create(&sonar.WebhooksCreateOption{
					URL: "https://example.com/webhook",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Name"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required URL", func() {
				result, resp, err := client.Webhooks.Create(&sonar.WebhooksCreateOption{
					Name: "Test Webhook",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("URL"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail with secret too short", func() {
				result, resp, err := client.Webhooks.Create(&sonar.WebhooksCreateOption{
					Name:   "Test Webhook",
					URL:    "https://example.com/webhook",
					Secret: "short",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Secret"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Valid Requests", func() {
			It("should create a global webhook", func() {
				webhookName := helpers.UniqueResourceName("wh")
				result, resp, err := client.Webhooks.Create(&sonar.WebhooksCreateOption{
					Name: webhookName,
					URL:  "https://example.com/webhook",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Webhook.Key).NotTo(BeEmpty())
				Expect(result.Webhook.Name).To(Equal(webhookName))

				// Clean up
				_, _ = client.Webhooks.Delete(&sonar.WebhooksDeleteOption{
					Webhook: result.Webhook.Key,
				})
			})

			It("should create a project webhook", func() {
				webhookName := helpers.UniqueResourceName("wh")
				result, resp, err := client.Webhooks.Create(&sonar.WebhooksCreateOption{
					Name:    webhookName,
					URL:     "https://example.com/webhook",
					Project: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Webhook.Key).NotTo(BeEmpty())

				// Clean up
				_, _ = client.Webhooks.Delete(&sonar.WebhooksDeleteOption{
					Webhook: result.Webhook.Key,
				})
			})

			It("should create a webhook with secret", func() {
				webhookName := helpers.UniqueResourceName("wh")
				result, resp, err := client.Webhooks.Create(&sonar.WebhooksCreateOption{
					Name:   webhookName,
					URL:    "https://example.com/webhook",
					Secret: "super-secret-key-16chars",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
				Expect(result.Webhook.HasSecret).To(BeTrue())

				// Clean up
				_, _ = client.Webhooks.Delete(&sonar.WebhooksDeleteOption{
					Webhook: result.Webhook.Key,
				})
			})
		})
	})

	// =========================================================================
	// Update
	// =========================================================================
	Describe("Update", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Webhooks.Update(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required webhook key", func() {
				resp, err := client.Webhooks.Update(&sonar.WebhooksUpdateOption{
					Name: "Updated Name",
					URL:  "https://example.com/updated",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Webhook"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required name", func() {
				resp, err := client.Webhooks.Update(&sonar.WebhooksUpdateOption{
					Webhook: "some-key",
					URL:     "https://example.com/updated",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Name"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required URL", func() {
				resp, err := client.Webhooks.Update(&sonar.WebhooksUpdateOption{
					Webhook: "some-key",
					Name:    "Updated Name",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("URL"))
				Expect(resp).To(BeNil())
			})
		})

		Context("Valid Requests", func() {
			It("should update a webhook", func() {
				// Create a webhook first
				webhookName := helpers.UniqueResourceName("wh")
				createResult, _, err := client.Webhooks.Create(&sonar.WebhooksCreateOption{
					Name: webhookName,
					URL:  "https://example.com/webhook",
				})
				Expect(err).NotTo(HaveOccurred())

				// Update it
				updatedName := webhookName + "-updated"
				resp, err := client.Webhooks.Update(&sonar.WebhooksUpdateOption{
					Webhook: createResult.Webhook.Key,
					Name:    updatedName,
					URL:     "https://example.com/updated",
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// Clean up
				_, _ = client.Webhooks.Delete(&sonar.WebhooksDeleteOption{
					Webhook: createResult.Webhook.Key,
				})
			})
		})

		Context("Non-Existent Webhook", func() {
			It("should fail for non-existent webhook", func() {
				resp, err := client.Webhooks.Update(&sonar.WebhooksUpdateOption{
					Webhook: "non-existent-webhook-key",
					Name:    "Updated Name",
					URL:     "https://example.com/updated",
				})
				Expect(err).To(HaveOccurred())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// Delete
	// =========================================================================
	Describe("Delete", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.Webhooks.Delete(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required webhook key", func() {
				resp, err := client.Webhooks.Delete(&sonar.WebhooksDeleteOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Webhook"))
				Expect(resp).To(BeNil())
			})
		})

		Context("Valid Requests", func() {
			It("should delete a webhook", func() {
				// Create a webhook first
				webhookName := helpers.UniqueResourceName("wh")
				createResult, _, err := client.Webhooks.Create(&sonar.WebhooksCreateOption{
					Name: webhookName,
					URL:  "https://example.com/webhook",
				})
				Expect(err).NotTo(HaveOccurred())

				// Delete it
				resp, err := client.Webhooks.Delete(&sonar.WebhooksDeleteOption{
					Webhook: createResult.Webhook.Key,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})
		})

		Context("Non-Existent Webhook", func() {
			It("should fail for non-existent webhook", func() {
				resp, err := client.Webhooks.Delete(&sonar.WebhooksDeleteOption{
					Webhook: "non-existent-webhook-key",
				})
				Expect(err).To(HaveOccurred())
				if resp != nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", 400))
				}
			})
		})
	})

	// =========================================================================
	// Deliveries
	// =========================================================================
	Describe("Deliveries", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Webhooks.Deliveries(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Valid Requests", func() {
			It("should list deliveries for a project", func() {
				result, resp, err := client.Webhooks.Deliveries(&sonar.WebhooksDeliveriesOption{
					ComponentKey: projectKey,
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})

			It("should list deliveries with pagination", func() {
				result, resp, err := client.Webhooks.Deliveries(&sonar.WebhooksDeliveriesOption{
					ComponentKey: projectKey,
					PaginationArgs: sonar.PaginationArgs{
						PageSize: 10,
						Page:     1,
					},
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})
	})

	// =========================================================================
	// Delivery
	// =========================================================================
	Describe("Delivery", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.Webhooks.Delivery(nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required DeliveryID", func() {
				result, resp, err := client.Webhooks.Delivery(&sonar.WebhooksDeliveryOption{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("DeliveryID"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Non-Existent Delivery", func() {
			It("should fail for non-existent delivery", func() {
				result, resp, err := client.Webhooks.Delivery(&sonar.WebhooksDeliveryOption{
					DeliveryID: "non-existent-delivery-id",
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
	// Full Workflow
	// =========================================================================
	Describe("Full Workflow", func() {
		It("should create, list, update, and delete a webhook", func() {
			// Create a webhook
			webhookName := helpers.UniqueResourceName("wh")
			createResult, resp, err := client.Webhooks.Create(&sonar.WebhooksCreateOption{
				Name: webhookName,
				URL:  "https://example.com/webhook",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(createResult).NotTo(BeNil())
			webhookKey := createResult.Webhook.Key

			// List and verify it's there
			listResult, resp, err := client.Webhooks.List(&sonar.WebhooksListOption{})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			found := false
			for _, w := range listResult.Webhooks {
				if w.Key == webhookKey {
					found = true
					Expect(w.Name).To(Equal(webhookName))
					break
				}
			}
			Expect(found).To(BeTrue(), "Webhook should be in the list")

			// Update the webhook
			updatedName := webhookName + "-updated"
			resp, err = client.Webhooks.Update(&sonar.WebhooksUpdateOption{
				Webhook: webhookKey,
				Name:    updatedName,
				URL:     "https://example.com/updated",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify update
			listResult, _, err = client.Webhooks.List(&sonar.WebhooksListOption{})
			Expect(err).NotTo(HaveOccurred())

			for _, w := range listResult.Webhooks {
				if w.Key == webhookKey {
					Expect(w.Name).To(Equal(updatedName))
					Expect(w.URL).To(Equal("https://example.com/updated"))
					break
				}
			}

			// Delete the webhook
			resp, err = client.Webhooks.Delete(&sonar.WebhooksDeleteOption{
				Webhook: webhookKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

			// Verify it's gone
			listResult, _, err = client.Webhooks.List(&sonar.WebhooksListOption{})
			Expect(err).NotTo(HaveOccurred())

			found = false
			for _, w := range listResult.Webhooks {
				if w.Key == webhookKey {
					found = true
					break
				}
			}
			Expect(found).To(BeFalse(), "Webhook should not be in the list")
		})
	})
})
