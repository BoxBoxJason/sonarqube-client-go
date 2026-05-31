package integration_testing_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("Project Pull Requests Service", Ordered, func() {
	var client *sonar.Client

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	// =========================================================================
	// Delete
	// =========================================================================
	Describe("Delete", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				resp, err := client.ProjectPullRequests.Delete(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required project", func() {
				resp, err := client.ProjectPullRequests.Delete(context.Background(), &sonar.ProjectPullRequestsDeleteOptions{
					PullRequest: "123",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Project"))
				Expect(resp).To(BeNil())
			})

			It("should fail without required pull request", func() {
				resp, err := client.ProjectPullRequests.Delete(context.Background(), &sonar.ProjectPullRequestsDeleteOptions{
					Project: "my-project",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("PullRequest"))
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should succeed or return an expected error", func() {
				resp, err := client.ProjectPullRequests.Delete(context.Background(), &sonar.ProjectPullRequestsDeleteOptions{
					Project:     "nonexistent-project",
					PullRequest: "nonexistent-pr",
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
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
				result, resp, err := client.ProjectPullRequests.List(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required project", func() {
				result, resp, err := client.ProjectPullRequests.List(context.Background(), &sonar.ProjectPullRequestsListOptions{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Project"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should succeed or return an expected error", func() {
				result, resp, err := client.ProjectPullRequests.List(context.Background(), &sonar.ProjectPullRequestsListOptions{
					Project: "nonexistent-project",
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})
})
