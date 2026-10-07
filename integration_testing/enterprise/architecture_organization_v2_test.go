package enterprise_test

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/boxboxjason/sonarqube-client-go/v2/integration_testing/helpers"
	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

// defaultOrganizationID is the identifier SonarQube Server uses for its single,
// implicit organization in V2 endpoints that take an organizationId.
const defaultOrganizationID = "00000000-0000-4000-0000-000000000000"

// architectureEnterpriseSetting enables the organization-level architecture
// features (placeholders, patterns, SDKs, boundary descriptors, ...). Without
// it every organization architecture endpoint answers 403.
const architectureEnterpriseSetting = "sonar.architecture.enterprise.enabled"

// projectUUID returns the project UUID (not the main branch UUID) of a project,
// which is what V2 endpoints expect as "projectId".
func projectUUID(client *sonar.Client, projectKey string) string {
	GinkgoHelper()

	result, _, err := client.Projects.Search(context.Background(), &sonar.ProjectsSearchOptions{Projects: []string{projectKey}})
	Expect(err).NotTo(HaveOccurred())
	Expect(result.Components).To(HaveLen(1))

	return result.Components[0].ProjectUuid
}

var _ = Describe("Architecture V2 organization features", Ordered, func() {
	var (
		client     *sonar.Client
		cleanup    *helpers.CleanupManager
		projectKey string
		projectID  string
	)

	BeforeAll(func() {
		var err error
		client, err = helpers.NewDefaultClient()
		Expect(err).NotTo(HaveOccurred())
		cleanup = helpers.NewCleanupManager(client)

		_, err = client.Settings.Set(context.Background(), &sonar.SettingsSetOptions{Key: architectureEnterpriseSetting, Value: "true"})
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			_, _ = client.Settings.Reset(context.Background(), &sonar.SettingsResetOptions{Keys: []string{architectureEnterpriseSetting}})
		})

		projectKey = helpers.UniqueResourceName("arch-org")
		_, _, err = client.Projects.Create(context.Background(), &sonar.ProjectsCreateOptions{Name: projectKey, Project: projectKey})
		Expect(err).NotTo(HaveOccurred())
		cleanup.RegisterCleanup("project", projectKey, func() error {
			_, err := client.Projects.Delete(context.Background(), &sonar.ProjectsDeleteOptions{Project: projectKey})
			return err
		})
		projectID = projectUUID(client, projectKey)
	})

	AfterAll(func() {
		for _, err := range cleanup.Cleanup() {
			GinkgoWriter.Printf("Cleanup error: %v\n", err)
		}
	})

	Describe("Boundary descriptors", Ordered, func() {
		var descriptorID string

		It("should create a boundary descriptor", func() {
			result, resp, err := client.V2.Architecture.CreateBoundaryDescriptor(context.Background(), &sonar.ArchitectureCreateBoundaryDescriptorOptions{
				ProjectId:   projectID,
				Name:        "payments-client",
				Direction:   "EXIT_POINT",
				Ecosystem:   "java",
				KeyTemplate: "payments-{x}",
				Query:       "call:com.example.Payments#charge",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))
			Expect(result.Id).NotTo(BeEmpty())
			Expect(result.ProjectId).To(Equal(projectID))
			descriptorID = result.Id
		})

		It("should get and list the boundary descriptor", func() {
			result, _, err := client.V2.Architecture.GetBoundaryDescriptor(context.Background(), descriptorID)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Name).To(Equal("payments-client"))

			list, _, err := client.V2.Architecture.ListBoundaryDescriptors(context.Background(), &sonar.ArchitectureListBoundaryDescriptorsOptions{ProjectId: projectID})
			Expect(err).NotTo(HaveOccurred())
			Expect(list.BoundaryDescriptors).To(HaveLen(1))
			Expect(list.Page.Total).To(BeEquivalentTo(1))
		})

		It("should update the boundary descriptor", func() {
			result, resp, err := client.V2.Architecture.UpdateBoundaryDescriptor(context.Background(), descriptorID, &sonar.ArchitectureUpdateBoundaryDescriptorOptions{
				Name:        new("payments-client-v2"),
				Direction:   new("ENTRY_POINT"),
				Ecosystem:   new("java"),
				KeyTemplate: new("payments-{x}"),
				Query:       new("call:com.example.Payments#refund"),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(result.Name).To(Equal("payments-client-v2"))
			Expect(result.Direction).To(Equal("ENTRY_POINT"))
		})

		It("should delete the boundary descriptor", func() {
			resp, err := client.V2.Architecture.DeleteBoundaryDescriptor(context.Background(), descriptorID)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})
	})

	Describe("Placeholders, components, SDKs and external interfaces", Ordered, func() {
		var placeholderID, sdkID string

		AfterAll(func() {
			if sdkID != "" {
				_, _ = client.V2.Architecture.DeleteSdk(context.Background(), sdkID)
			}
			if placeholderID != "" {
				_, _ = client.V2.Architecture.DeleteOrganizationPlaceholder(context.Background(), placeholderID)
			}
		})

		It("should create a placeholder", func() {
			result, resp, err := client.V2.Architecture.CreateOrganizationPlaceholder(context.Background(), &sonar.ArchitectureCreateOrganizationPlaceholderOptions{
				OrganizationId: defaultOrganizationID,
				DisplayName:    "Payments DB",
				Type:           "DB",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))
			Expect(result.Type).To(Equal("DB"))
			Expect(result.OrganizationId).To(Equal(defaultOrganizationID))
			placeholderID = result.Id
		})

		It("should get, list and update the placeholder", func() {
			result, _, err := client.V2.Architecture.GetOrganizationPlaceholder(context.Background(), placeholderID)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.DisplayName).To(Equal("Payments DB"))

			list, _, err := client.V2.Architecture.ListOrganizationPlaceholders(context.Background(), &sonar.ArchitectureListOrganizationPlaceholdersOptions{OrganizationId: defaultOrganizationID})
			Expect(err).NotTo(HaveOccurred())
			Expect(list.Placeholders).NotTo(BeEmpty())

			updated, _, err := client.V2.Architecture.UpdateOrganizationPlaceholder(context.Background(), placeholderID, &sonar.ArchitectureUpdateOrganizationPlaceholderOptions{
				DisplayName: new("Payments queue"),
				Type:        new("QUEUE"),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.DisplayName).To(Equal("Payments queue"))
			Expect(updated.Type).To(Equal("QUEUE"))
		})

		It("should list the organization components, including the project and the placeholder", func() {
			list, resp, err := client.V2.Architecture.ListOrganizationComponents(context.Background(), &sonar.ArchitectureListOrganizationComponentsOptions{OrganizationId: defaultOrganizationID})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var projectComponentID string
			ids := make([]string, 0, len(list.Components))
			for _, component := range list.Components {
				ids = append(ids, component.Id)
				if component.ProjectId == projectID {
					Expect(component.Type).To(Equal("PROJECT"))
					projectComponentID = component.Id
				}
			}
			Expect(ids).To(ContainElement(placeholderID))
			Expect(projectComponentID).NotTo(BeEmpty())

			component, _, err := client.V2.Architecture.GetOrganizationComponent(context.Background(), projectComponentID)
			Expect(err).NotTo(HaveOccurred())
			Expect(component.ProjectId).To(Equal(projectID))
		})

		It("should create, get, list and update an SDK", func() {
			created, resp, err := client.V2.Architecture.CreateSdk(context.Background(), &sonar.ArchitectureCreateSdkOptions{
				OrganizationId:    defaultOrganizationID,
				Name:              "payments-sdk",
				SdkKey:            helpers.UniqueResourceName("payments-sdk"),
				TargetComponentId: placeholderID,
				SdkBoundaryDescriptors: []sonar.ArchitectureSdkBoundaryDescriptor{
					{Ecosystem: "java", Query: "call:com.example.Payments#charge"},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))
			Expect(created.SdkBoundaryDescriptors).To(HaveLen(1))
			sdkID = created.Id

			got, _, err := client.V2.Architecture.GetSdk(context.Background(), sdkID)
			Expect(err).NotTo(HaveOccurred())
			Expect(got.Name).To(Equal("payments-sdk"))

			list, _, err := client.V2.Architecture.ListSdks(context.Background(), &sonar.ArchitectureListSdksOptions{OrganizationId: defaultOrganizationID})
			Expect(err).NotTo(HaveOccurred())
			Expect(list.Sdks).NotTo(BeEmpty())

			updated, _, err := client.V2.Architecture.UpdateSdk(context.Background(), sdkID, &sonar.ArchitectureUpdateSdkOptions{Name: new("payments-sdk-v2")})
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.Name).To(Equal("payments-sdk-v2"))
		})

		It("should list external interfaces and the organization architecture", func() {
			interfaces, resp, err := client.V2.Architecture.ListExternalInterfaces(context.Background(), &sonar.ArchitectureListExternalInterfacesOptions{OrganizationId: defaultOrganizationID})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(interfaces).NotTo(BeNil())

			architecture, resp, err := client.V2.Architecture.GetOrganizationArchitecture(context.Background(), &sonar.ArchitectureGetOrganizationArchitectureOptions{OrganizationId: defaultOrganizationID})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(architecture.OrganizationArchitectures).NotTo(BeEmpty())
		})

		It("should delete the SDK and the placeholder", func() {
			resp, err := client.V2.Architecture.DeleteSdk(context.Background(), sdkID)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			sdkID = ""

			resp, err = client.V2.Architecture.DeleteOrganizationPlaceholder(context.Background(), placeholderID)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			placeholderID = ""
		})
	})

	Describe("Patterns", Ordered, func() {
		var patternID string

		It("should create, get, list, update and delete a pattern", func() {
			created, resp, err := client.V2.Architecture.CreatePattern(context.Background(), &sonar.ArchitectureCreatePatternOptions{
				OrganizationId: defaultOrganizationID,
				Name:           "layered",
				Pattern: sonar.ArchitecturePatternData{
					Label:       "Layered",
					Description: "API depends on core",
					Groups: []sonar.ArchitectureModelGroup{
						{Label: "api", Patterns: []string{"com.example.api.**"}},
						{Label: "core", Patterns: []string{"com.example.core.**"}},
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))
			Expect(created.Pattern.Groups).To(HaveLen(2))
			patternID = created.Id

			got, _, err := client.V2.Architecture.GetPattern(context.Background(), patternID)
			Expect(err).NotTo(HaveOccurred())
			Expect(got.Name).To(Equal("layered"))

			list, _, err := client.V2.Architecture.ListPatterns(context.Background(), &sonar.ArchitectureListPatternsOptions{OrganizationId: defaultOrganizationID})
			Expect(err).NotTo(HaveOccurred())
			Expect(list).NotTo(BeEmpty())

			updated, _, err := client.V2.Architecture.UpdatePattern(context.Background(), patternID, &sonar.ArchitectureUpdatePatternOptions{Name: new("layered-v2")})
			Expect(err).NotTo(HaveOccurred())
			Expect(updated.Name).To(Equal("layered-v2"))

			resp, err = client.V2.Architecture.DeletePattern(context.Background(), patternID)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
		})
	})

	Describe("Project relationships", func() {
		It("should list the (empty) relationships of a project", func() {
			list, resp, err := client.V2.Architecture.ListProjectRelationships(context.Background(), &sonar.ArchitectureListProjectRelationshipsOptions{ProjectId: projectID})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(list.ProjectRelationships).To(BeEmpty())
		})

		// A relationship can only be created from an exit boundary found by an
		// analysis; this suite has no scanner, so the server rejects it with 400.
		It("should reject a relationship without a matching analysed exit boundary", func() {
			_, resp, err := client.V2.Architecture.CreateProjectRelationship(context.Background(), &sonar.ArchitectureCreateProjectRelationshipOptions{
				ProjectId:           projectID,
				BoundaryKey:         "payments-1",
				TargetComponentId:   projectID,
				TargetEntryPointKey: "default",
			})
			Expect(err).To(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})
	})
})
