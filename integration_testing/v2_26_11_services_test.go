package integration_testing_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

// Specs for the V2 endpoints added with SonarQube 26.11 (components projects lookup).
var _ = Describe("SonarQube 26.11 V2 services", Ordered, func() {
	var (
		client *sonar.Client
		ctx    = context.Background()
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Components projects", func() {
		const projectKey = "e2e-components-v2-project"

		BeforeAll(func() {
			_, _, err := client.Projects.Create(ctx, &sonar.ProjectsCreateOptions{Name: "E2E Components V2", Project: projectKey})
			if err != nil {
				Skip("components/projects endpoint prerequisites unavailable: " + err.Error())
			}

			DeferCleanup(func() {
				_, _ = client.Projects.Delete(ctx, &sonar.ProjectsDeleteOptions{Project: projectKey})
			})

			_, resp, err := client.V2.Components.GetProjects(ctx, &sonar.ComponentsGetProjectsOptionsV2{Keys: []string{projectKey}})
			if err != nil && resp != nil && resp.StatusCode == http.StatusNotFound {
				Skip("components/projects is not available on this SonarQube version")
			}
		})

		It("should find a project by key and then by id", func() {
			byKey, _, err := client.V2.Components.GetProjects(ctx, &sonar.ComponentsGetProjectsOptionsV2{Keys: []string{projectKey}})
			Expect(err).NotTo(HaveOccurred())
			Expect(byKey.Projects).To(HaveLen(1))
			Expect(byKey.Projects[0].Key).To(Equal(projectKey))
			Expect(byKey.Page.Total).To(Equal(int32(1)))

			byID, _, err := client.V2.Components.GetProjects(ctx, &sonar.ComponentsGetProjectsOptionsV2{Ids: []string{byKey.Projects[0].Id}})
			Expect(err).NotTo(HaveOccurred())
			Expect(byID.Projects).To(HaveLen(1))
			Expect(byID.Projects[0].Key).To(Equal(projectKey))
		})

		It("should filter by tags", func() {
			res, _, err := client.V2.Components.GetProjects(ctx, &sonar.ComponentsGetProjectsOptionsV2{
				Keys: []string{projectKey}, Tags: []string{"no-such-tag"},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(res.Projects).To(BeEmpty())
		})

		It("should require ids or keys", func() {
			_, _, err := client.V2.Components.GetProjects(ctx, &sonar.ComponentsGetProjectsOptionsV2{})
			Expect(err).To(HaveOccurred())
		})
	})
})
