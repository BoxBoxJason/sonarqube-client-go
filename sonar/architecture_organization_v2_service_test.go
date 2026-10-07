package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArchitectureService_DeleteBoundaryDescriptor(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodDelete, "/v2/architecture/boundary-descriptors/test-id", http.StatusNoContent))
	client := newTestClient(t, server.url())

	resp, err := client.V2.Architecture.DeleteBoundaryDescriptor(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestArchitectureService_DeleteBoundaryDescriptor_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, err := client.V2.Architecture.DeleteBoundaryDescriptor(context.Background(), "")
	require.Error(t, err)
}

func TestArchitectureService_DeleteOrganizationPlaceholder(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodDelete, "/v2/architecture/organization-placeholders/test-id", http.StatusNoContent))
	client := newTestClient(t, server.url())

	resp, err := client.V2.Architecture.DeleteOrganizationPlaceholder(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestArchitectureService_DeleteOrganizationPlaceholder_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, err := client.V2.Architecture.DeleteOrganizationPlaceholder(context.Background(), "")
	require.Error(t, err)
}

func TestArchitectureService_DeletePattern(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodDelete, "/v2/architecture/patterns/test-id", http.StatusNoContent))
	client := newTestClient(t, server.url())

	resp, err := client.V2.Architecture.DeletePattern(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestArchitectureService_DeletePattern_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, err := client.V2.Architecture.DeletePattern(context.Background(), "")
	require.Error(t, err)
}

func TestArchitectureService_DeleteProjectRelationship(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodDelete, "/v2/architecture/project-relationships/test-id", http.StatusNoContent))
	client := newTestClient(t, server.url())

	resp, err := client.V2.Architecture.DeleteProjectRelationship(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestArchitectureService_DeleteProjectRelationship_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, err := client.V2.Architecture.DeleteProjectRelationship(context.Background(), "")
	require.Error(t, err)
}

func TestArchitectureService_DeleteSdk(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodDelete, "/v2/architecture/sdks/test-id", http.StatusNoContent))
	client := newTestClient(t, server.url())

	resp, err := client.V2.Architecture.DeleteSdk(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestArchitectureService_DeleteSdk_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, err := client.V2.Architecture.DeleteSdk(context.Background(), "")
	require.Error(t, err)
}

func TestArchitectureService_ListBoundaryDescriptors(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/architecture/boundary-descriptors", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.ListBoundaryDescriptors(context.Background(), &ArchitectureListBoundaryDescriptorsOptions{ProjectId: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_ListBoundaryDescriptors_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.ListBoundaryDescriptors(context.Background(), nil)
	require.Error(t, err)
}

func TestArchitectureService_GetBoundaryDescriptor(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/architecture/boundary-descriptors/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.GetBoundaryDescriptor(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_GetBoundaryDescriptor_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.GetBoundaryDescriptor(context.Background(), "")
	require.Error(t, err)
}

func TestArchitectureService_ListExternalInterfaces(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/architecture/external-interfaces", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.ListExternalInterfaces(context.Background(), &ArchitectureListExternalInterfacesOptions{OrganizationId: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_ListExternalInterfaces_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.ListExternalInterfaces(context.Background(), nil)
	require.Error(t, err)
}

func TestArchitectureService_GetOrganizationArchitecture(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/architecture/organization-architectures", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.GetOrganizationArchitecture(context.Background(), &ArchitectureGetOrganizationArchitectureOptions{OrganizationId: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_GetOrganizationArchitecture_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.GetOrganizationArchitecture(context.Background(), nil)
	require.Error(t, err)
}

func TestArchitectureService_ListOrganizationComponents(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/architecture/organization-components", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.ListOrganizationComponents(context.Background(), &ArchitectureListOrganizationComponentsOptions{OrganizationId: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_ListOrganizationComponents_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.ListOrganizationComponents(context.Background(), nil)
	require.Error(t, err)
}

func TestArchitectureService_GetOrganizationComponent(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/architecture/organization-components/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.GetOrganizationComponent(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_GetOrganizationComponent_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.GetOrganizationComponent(context.Background(), "")
	require.Error(t, err)
}

func TestArchitectureService_ListOrganizationPlaceholders(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/architecture/organization-placeholders", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.ListOrganizationPlaceholders(context.Background(), &ArchitectureListOrganizationPlaceholdersOptions{OrganizationId: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_ListOrganizationPlaceholders_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.ListOrganizationPlaceholders(context.Background(), nil)
	require.Error(t, err)
}

func TestArchitectureService_GetOrganizationPlaceholder(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/architecture/organization-placeholders/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.GetOrganizationPlaceholder(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_GetOrganizationPlaceholder_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.GetOrganizationPlaceholder(context.Background(), "")
	require.Error(t, err)
}

func TestArchitectureService_ListPatterns(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/architecture/patterns", http.StatusOK, []any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.ListPatterns(context.Background(), &ArchitectureListPatternsOptions{OrganizationId: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_ListPatterns_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.ListPatterns(context.Background(), nil)
	require.Error(t, err)
}

func TestArchitectureService_GetPattern(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/architecture/patterns/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.GetPattern(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_GetPattern_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.GetPattern(context.Background(), "")
	require.Error(t, err)
}

func TestArchitectureService_ListProjectRelationships(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/architecture/project-relationships", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.ListProjectRelationships(context.Background(), &ArchitectureListProjectRelationshipsOptions{ProjectId: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_ListProjectRelationships_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.ListProjectRelationships(context.Background(), nil)
	require.Error(t, err)
}

func TestArchitectureService_GetProjectRelationship(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/architecture/project-relationships/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.GetProjectRelationship(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_GetProjectRelationship_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.GetProjectRelationship(context.Background(), "")
	require.Error(t, err)
}

func TestArchitectureService_ListSdks(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/architecture/sdks", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.ListSdks(context.Background(), &ArchitectureListSdksOptions{OrganizationId: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_ListSdks_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.ListSdks(context.Background(), nil)
	require.Error(t, err)
}

func TestArchitectureService_GetSdk(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/architecture/sdks/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.GetSdk(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_GetSdk_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.GetSdk(context.Background(), "")
	require.Error(t, err)
}

func TestArchitectureService_UpdateBoundaryDescriptor(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPatch, "/v2/architecture/boundary-descriptors/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.UpdateBoundaryDescriptor(context.Background(), "test-id", &ArchitectureUpdateBoundaryDescriptorOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_UpdateBoundaryDescriptor_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.UpdateBoundaryDescriptor(context.Background(), "", nil)
	require.Error(t, err)
}

func TestArchitectureService_UpdateOrganizationPlaceholder(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPatch, "/v2/architecture/organization-placeholders/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.UpdateOrganizationPlaceholder(context.Background(), "test-id", &ArchitectureUpdateOrganizationPlaceholderOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_UpdateOrganizationPlaceholder_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.UpdateOrganizationPlaceholder(context.Background(), "", nil)
	require.Error(t, err)
}

func TestArchitectureService_UpdatePattern(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPatch, "/v2/architecture/patterns/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.UpdatePattern(context.Background(), "test-id", &ArchitectureUpdatePatternOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_UpdatePattern_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.UpdatePattern(context.Background(), "", nil)
	require.Error(t, err)
}

func TestArchitectureService_UpdateProjectRelationship(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPatch, "/v2/architecture/project-relationships/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.UpdateProjectRelationship(context.Background(), "test-id", &ArchitectureUpdateProjectRelationshipOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_UpdateProjectRelationship_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.UpdateProjectRelationship(context.Background(), "", nil)
	require.Error(t, err)
}

func TestArchitectureService_UpdateSdk(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPatch, "/v2/architecture/sdks/test-id", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.UpdateSdk(context.Background(), "test-id", &ArchitectureUpdateSdkOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_UpdateSdk_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.UpdateSdk(context.Background(), "", nil)
	require.Error(t, err)
}

func TestArchitectureService_CreateBoundaryDescriptor(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/architecture/boundary-descriptors", http.StatusCreated, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.CreateBoundaryDescriptor(context.Background(), &ArchitectureCreateBoundaryDescriptorOptions{Ecosystem: "java", KeyTemplate: "test", Query: "test", Direction: "EXIT_POINT", ProjectId: "test", Name: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_CreateBoundaryDescriptor_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.CreateBoundaryDescriptor(context.Background(), nil)
	require.Error(t, err)
}

func TestArchitectureService_CreateOrganizationPlaceholder(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/architecture/organization-placeholders", http.StatusCreated, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.CreateOrganizationPlaceholder(context.Background(), &ArchitectureCreateOrganizationPlaceholderOptions{DisplayName: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_CreateOrganizationPlaceholder_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.CreateOrganizationPlaceholder(context.Background(), nil)
	require.Error(t, err)
}

func TestArchitectureService_CreatePattern(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/architecture/patterns", http.StatusCreated, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.CreatePattern(context.Background(), &ArchitectureCreatePatternOptions{OrganizationId: "test", Pattern: ArchitecturePatternData{}, Name: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_CreatePattern_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.CreatePattern(context.Background(), nil)
	require.Error(t, err)
}

func TestArchitectureService_CreateProjectRelationship(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/architecture/project-relationships", http.StatusCreated, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.CreateProjectRelationship(context.Background(), &ArchitectureCreateProjectRelationshipOptions{BoundaryKey: "test", ProjectId: "test", TargetComponentId: "test", TargetEntryPointKey: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_CreateProjectRelationship_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.CreateProjectRelationship(context.Background(), nil)
	require.Error(t, err)
}

func TestArchitectureService_CreateSdk(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/architecture/sdks", http.StatusCreated, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Architecture.CreateSdk(context.Background(), &ArchitectureCreateSdkOptions{TargetComponentId: "test", SdkBoundaryDescriptors: []ArchitectureSdkBoundaryDescriptor{{}}, Name: "test", SdkKey: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestArchitectureService_CreateSdk_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Architecture.CreateSdk(context.Background(), nil)
	require.Error(t, err)
}
