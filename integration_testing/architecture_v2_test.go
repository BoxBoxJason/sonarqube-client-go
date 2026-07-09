package integration_testing_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/sonar"
)

var _ = Describe("Architecture V2 Service", Ordered, func() {
	var client *sonar.Client

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	// =========================================================================
	// FileGraph
	// =========================================================================
	Describe("FileGraph", func() {
		Context("Parameter Validation", func() {
			It("should fail with nil options", func() {
				result, resp, err := client.V2.Architecture.FileGraph(context.Background(), nil)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("required"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})

			It("should fail without required project key", func() {
				result, resp, err := client.V2.Architecture.FileGraph(context.Background(), &sonar.ArchitectureFileGraphOptions{
					BranchKey: "main",
					Source:    "java",
				})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ProjectKey"))
				Expect(result).To(BeNil())
				Expect(resp).To(BeNil())
			})
		})

		Context("Functional Tests", func() {
			It("should return file graph or an enterprise-only error", func() {
				result, resp, err := client.V2.Architecture.FileGraph(context.Background(), &sonar.ArchitectureFileGraphOptions{
					ProjectKey: "nonexistent-project",
					BranchKey:  "main",
					Source:     "java",
				})
				if err != nil {
					Expect(resp).NotTo(BeNil())
					// Enterprise-only / not-found gate: accept the expected error codes only,
					// so unrelated failures (network errors, 5xx, decoding issues) are not
					// silently swallowed by this test.
					Expect(resp.StatusCode).To(BeElementOf(http.StatusNotFound, http.StatusForbidden, http.StatusPaymentRequired))
				} else {
					Expect(resp.StatusCode).To(BeNumerically("<", 400))
					Expect(result).NotTo(BeNil())
				}
			})
		})
	})
})
