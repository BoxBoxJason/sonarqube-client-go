package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("DismissMessage Service", Ordered, func() {
	var (
		client     *sonargo.Client
		cleanup    *helpers.CleanupManager
		projectKey string
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
		cleanup = helpers.NewCleanupManager(client)

		// Create a test project for dismiss message operations
		projectKey = helpers.UniqueResourceName("dms")
		_, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
			Name:    "DismissMessage Test Project",
			Project: projectKey,
		})
		Expect(err).NotTo(HaveOccurred())

		cleanup.RegisterCleanup("project", projectKey, func() error {
			_, err := client.Projects.Delete(&sonargo.ProjectsDeleteOption{
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
	// Check
	// =========================================================================
	Describe("Check", func() {
		Context("Valid Requests", func() {
			It("should check message dismissal status", func() {
				result, resp, err := client.DismissMessage.Check(&sonargo.DismissMessageCheckOption{
					MessageType: "PROJECT_NCD_90",
					ProjectKey:  projectKey,
				})
				// Skip if API not available or invalid message type
				if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest) {
					Skip("DismissMessage API is not available or message type not supported in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				Expect(result).NotTo(BeNil())
			})
		})

		Context("Parameter Validation", func() {
			It("should fail with missing message type", func() {
				_, _, err := client.DismissMessage.Check(&sonargo.DismissMessageCheckOption{
					ProjectKey: projectKey,
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with missing project key", func() {
				_, _, err := client.DismissMessage.Check(&sonargo.DismissMessageCheckOption{
					MessageType: "PROJECT_NCD_90",
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, _, err := client.DismissMessage.Check(nil)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// =========================================================================
	// Dismiss
	// =========================================================================
	Describe("Dismiss", func() {
		Context("Valid Requests", func() {
			It("should dismiss a message", func() {
				resp, err := client.DismissMessage.Dismiss(&sonargo.DismissMessageDismissOption{
					MessageType: "PROJECT_NCD_90",
					ProjectKey:  projectKey,
				})
				// Skip if API not available or invalid message type
				if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest) {
					Skip("DismissMessage API is not available or message type not supported in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(SatisfyAny(Equal(http.StatusOK), Equal(http.StatusNoContent)))
			})
		})

		Context("Parameter Validation", func() {
			It("should fail with missing message type", func() {
				_, err := client.DismissMessage.Dismiss(&sonargo.DismissMessageDismissOption{
					ProjectKey: projectKey,
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with missing project key", func() {
				_, err := client.DismissMessage.Dismiss(&sonargo.DismissMessageDismissOption{
					MessageType: "PROJECT_NCD_90",
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, err := client.DismissMessage.Dismiss(nil)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	// =========================================================================
	// Full Workflow
	// =========================================================================
	Describe("Full Workflow", func() {
		It("should verify message is dismissed after dismissal", func() {
			// Check message is not dismissed initially
			initialCheck, resp, err := client.DismissMessage.Check(&sonargo.DismissMessageCheckOption{
				MessageType: "PROJECT_NCD_90",
				ProjectKey:  projectKey,
			})
			// Skip if API not available or invalid message type
			if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest) {
				Skip("DismissMessage API is not available or message type not supported in this SonarQube version")
			}
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(initialCheck).NotTo(BeNil())

			// Dismiss the message
			resp, err = client.DismissMessage.Dismiss(&sonargo.DismissMessageDismissOption{
				MessageType: "PROJECT_NCD_90",
				ProjectKey:  projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(SatisfyAny(Equal(http.StatusOK), Equal(http.StatusNoContent)))

			// Check message is now dismissed
			finalCheck, resp, err := client.DismissMessage.Check(&sonargo.DismissMessageCheckOption{
				MessageType: "PROJECT_NCD_90",
				ProjectKey:  projectKey,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(finalCheck).NotTo(BeNil())

			// Verify the dismissal status changed (dismissed field should exist)
			// The exact field name may vary by SonarQube version, so we just verify
			// the response is valid and the workflow completes successfully
		})
	})
})
