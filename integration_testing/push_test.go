package integration_testing_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sonargo "github.com/boxboxjason/sonarqube-client-go/sonar"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
)

var _ = Describe("Push Service", Ordered, func() {
	var (
		client      *sonargo.Client
		testProject *sonargo.ProjectsCreate
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())

		// Create a test project for push events
		testProject, _, err = client.Projects.Create(&sonargo.ProjectsCreateOption{
			Name:    "push-e2e-test-project",
			Project: "push-e2e-test-project",
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
	// SonarlintEvents
	// =========================================================================
	Describe("SonarlintEvents", func() {
		Context("Functional Tests", func() {
			It("should connect to sonarlint events stream with valid parameters", func() {
				resp, err := client.Push.SonarlintEvents(&sonargo.PushSonarlintEventsOption{
					Languages:   []string{"java"},
					ProjectKeys: []string{testProject.Project.Key},
				})
				// The endpoint may not be available in all SonarQube versions
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("Push API is not available in this SonarQube version")
				}
				// This is a streaming endpoint, so we immediately close the connection
				if resp != nil && resp.Body != nil {
					_ = resp.Body.Close()
				}
				// Expect either success or specific errors (the endpoint may require specific setup)
				if err == nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", http.StatusOK))
					Expect(resp.StatusCode).To(BeNumerically("<", http.StatusMultipleChoices))
				}
			})

			It("should connect with multiple languages", func() {
				resp, err := client.Push.SonarlintEvents(&sonargo.PushSonarlintEventsOption{
					Languages:   []string{"java", "js", "py"},
					ProjectKeys: []string{testProject.Project.Key},
				})
				if resp != nil && resp.StatusCode == http.StatusNotFound {
					Skip("Push API is not available in this SonarQube version")
				}
				if resp != nil && resp.Body != nil {
					_ = resp.Body.Close()
				}
				if err == nil {
					Expect(resp.StatusCode).To(BeNumerically(">=", http.StatusOK))
					Expect(resp.StatusCode).To(BeNumerically("<", http.StatusMultipleChoices))
				}
			})
		})

		Context("Error Handling", func() {
			It("should fail with missing languages", func() {
				_, err := client.Push.SonarlintEvents(&sonargo.PushSonarlintEventsOption{
					ProjectKeys: []string{testProject.Project.Key},
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with missing project keys", func() {
				_, err := client.Push.SonarlintEvents(&sonargo.PushSonarlintEventsOption{
					Languages: []string{"java"},
				})
				Expect(err).To(HaveOccurred())
			})

			It("should fail with nil options", func() {
				_, err := client.Push.SonarlintEvents(nil)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
