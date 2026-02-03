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
		client      *sonargo.Client
		testProject *sonargo.ProjectsCreate
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())

		// Create a test project for dismiss message operations
		testProject, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
			Name:    "dismiss-message-e2e-test-project",
			Project: "dismiss-message-e2e-test-project",
		})
		Expect(err).NotTo(HaveOccurred())
	})

	AfterAll(func() {
		if testProject != nil {
			_, _ = client.Projects.Delete(&sonargo.ProjectsDeleteOption{
				Project: testProject.Project.Key,
			})
		}
	})

	// =========================================================================
	// Check
	// =========================================================================
	Describe("Check", func() {
		Context("Functional Tests", func() {
			It("should check message dismissal status", func() {
				result, resp, err := client.DismissMessage.Check(&sonargo.DismissMessageCheckOption{
					MessageType: "PROJECT_NCD_90",
					ProjectKey:  testProject.Project.Key,
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

		Context("Error Handling", func() {
			It("should fail with missing message type", func() {
				_, _, err := client.DismissMessage.Check(&sonargo.DismissMessageCheckOption{
					ProjectKey: testProject.Project.Key,
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
		Context("Functional Tests", func() {
			It("should dismiss a message", func() {
				resp, err := client.DismissMessage.Dismiss(&sonargo.DismissMessageDismissOption{
					MessageType: "PROJECT_NCD_90",
					ProjectKey:  testProject.Project.Key,
				})
				// Skip if API not available or invalid message type
				if resp != nil && (resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest) {
					Skip("DismissMessage API is not available or message type not supported in this SonarQube version")
				}
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(SatisfyAny(Equal(http.StatusOK), Equal(http.StatusNoContent)))
			})
		})

		Context("Error Handling", func() {
			It("should fail with missing message type", func() {
				_, err := client.DismissMessage.Dismiss(&sonargo.DismissMessageDismissOption{
					ProjectKey: testProject.Project.Key,
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
})
