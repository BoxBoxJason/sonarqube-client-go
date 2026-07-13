package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScaService_ListClis(t *testing.T) {
	// The real API returns a bare JSON array (verified against a live SonarQube 2025.2
	// Enterprise instance), not an object wrapping a "clis" field.
	response := []ScaCliInfo{{Id: "cli-1", Filename: "sca-cli-linux", Os: "linux", Arch: "amd64"}}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/sca/clis", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.ListClis(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result, 1)
}

func TestScaService_GetCli(t *testing.T) {
	response := ScaCliInfo{Id: "cli-1", Filename: "sca-cli-linux", Os: "linux", Arch: "amd64"}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/sca/clis/cli-1", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.GetCli(context.Background(), &ScaCliGetOptions{Id: "cli-1"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "cli-1", result.Id)
}

func TestScaService_GetCli_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.GetCli(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.GetCli(context.Background(), &ScaCliGetOptions{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_GetEnablement(t *testing.T) {
	response := ScaFeatureEnablement{Enablement: true}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/sca/feature-enablements", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.GetEnablement(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, result.Enablement)
}

func TestScaService_SetEnablement(t *testing.T) {
	response := ScaFeatureEnablement{Enablement: false}
	server := newTestServer(t, mockPatchHandler(t, "/v2/sca/feature-enablements", http.StatusOK,
		map[string]any{"enablement": false}, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.SetEnablement(context.Background(), &ScaSetEnablementOptions{Enablement: false})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.False(t, result.Enablement)
}

func TestScaService_SetEnablement_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.SetEnablement(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_SearchDependencyRisks(t *testing.T) {
	response := ScaDependencyRisksSearch{
		IssuesReleases: []ScaDependencyRisk{{Key: "risk-1", Severity: "HIGH"}},
		Page:           PageResponseV2{PageIndex: 1, PageSize: 20, Total: 1},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/sca/issues-releases", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.SearchDependencyRisks(context.Background(), &ScaDependencyRisksSearchOptions{
		ProjectKey: "my-project",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.IssuesReleases, 1)
}

func TestScaService_SearchDependencyRisks_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.SearchDependencyRisks(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.SearchDependencyRisks(context.Background(), &ScaDependencyRisksSearchOptions{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_GetDependencyRisk(t *testing.T) {
	response := ScaDependencyRisk{Key: "risk-1", Severity: "HIGH", Type: "VULNERABILITY"}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/sca/issues-releases/risk-1", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.GetDependencyRisk(context.Background(), &ScaDependencyRiskGetOptions{Key: "risk-1"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "HIGH", result.Severity)
}

func TestScaService_GetDependencyRisk_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.GetDependencyRisk(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.GetDependencyRisk(context.Background(), &ScaDependencyRiskGetOptions{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_SearchReleases(t *testing.T) {
	// packageManagerCounts is confirmed present on the real API's response envelope (verified
	// live against a SonarQube 2025.2 Enterprise instance), alongside releases and page.
	response := ScaReleasesSearch{
		Releases:             []ScaReleaseSearchResource{{Key: "rel-1", PackageName: "lodash"}},
		PackageManagerCounts: []ScaReleasePackageManagerCount{{PackageManager: "npm", ReleaseCount: 1}},
		Page:                 PageResponseV2{Total: 1},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/sca/releases", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.SearchReleases(context.Background(), &ScaReleasesSearchOptions{
		ProjectKey: "my-project",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Releases, 1)
	assert.Len(t, result.PackageManagerCounts, 1)
	assert.Equal(t, "npm", result.PackageManagerCounts[0].PackageManager)
}

func TestScaService_SearchReleases_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.SearchReleases(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.SearchReleases(context.Background(), &ScaReleasesSearchOptions{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_GetRelease(t *testing.T) {
	response := ScaReleaseDetail{Key: "rel-1", PackageName: "lodash", Version: "4.17.21"}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/sca/releases/rel-1", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.GetRelease(context.Background(), &ScaReleaseGetOptions{Key: "rel-1"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "lodash", result.PackageName)
}

func TestScaService_GetRelease_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.GetRelease(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.GetRelease(context.Background(), &ScaReleaseGetOptions{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_GetSbomReport(t *testing.T) {
	// Verified live against a SonarQube 2025.2 Enterprise instance: the "type" query parameter
	// only accepts "cyclonedx" or "spdx_23", and the actual serialization (JSON vs XML) is
	// selected by the request's Accept header, not by "type" alone. A request without a
	// matching Accept header (e.g. the previous "CYCLONEDX_1_4_JSON" example value with no
	// Accept override) is rejected by the real server with 400 "No report found for type...".
	data := []byte(`{"bomFormat":"CycloneDX","specVersion":"1.4"}`)
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/vnd.cyclonedx+json", r.Header.Get("Accept"))
		mockBinaryHandler(t, http.MethodGet, "/v2/sca/sbom-reports", http.StatusOK, "application/json", data)(w, r)
	})
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.GetSbomReport(context.Background(), &ScaSbomReportOptions{
		Component: "my-project",
		Type:      ScaSbomReportTypeCycloneDX,
		Format:    ScaSbomReportFormatJSON,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, data, result)
}

func TestScaService_GetSbomReport_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.GetSbomReport(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.GetSbomReport(context.Background(), &ScaSbomReportOptions{
		Type:   ScaSbomReportTypeCycloneDX,
		Format: ScaSbomReportFormatJSON,
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.GetSbomReport(context.Background(), &ScaSbomReportOptions{
		Component: "my-project",
		Format:    ScaSbomReportFormatJSON,
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.GetSbomReport(context.Background(), &ScaSbomReportOptions{
		Component: "my-project",
		Type:      ScaSbomReportTypeCycloneDX,
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_ListLicenseProfiles(t *testing.T) {
	response := ScaLicenseProfileIndex{
		LicenseProfiles: []ScaLicenseProfile{{Id: "lp-1", Name: "Default"}},
		Actions:         ScaLicenseProfileCollectionActions{Create: true},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/sca/license-profiles", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.ListLicenseProfiles(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, result.LicenseProfiles, 1)
	assert.Equal(t, "lp-1", result.LicenseProfiles[0].Id)
	assert.True(t, result.Actions.Create)

	result, resp, err = client.V2.Sca.ListLicenseProfiles(context.Background(), &ScaLicenseProfileListOptions{ProjectKey: "my-project"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, result.LicenseProfiles, 1)
}

func TestScaService_CreateLicenseProfile(t *testing.T) {
	body := &ScaLicenseProfileCreateRequest{Name: "Default", Organization: "my-org"}
	response := ScaLicenseProfile{Id: "lp-1", Name: "Default"}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/sca/license-profiles", http.StatusCreated, body, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.CreateLicenseProfile(context.Background(), body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "lp-1", result.Id)
}

func TestScaService_CreateLicenseProfile_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.CreateLicenseProfile(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.CreateLicenseProfile(context.Background(), &ScaLicenseProfileCreateRequest{Organization: "my-org"})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.CreateLicenseProfile(context.Background(), &ScaLicenseProfileCreateRequest{Name: "Default"})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_ListLicenseProfileAssignableProjects(t *testing.T) {
	response := ScaLicenseProfileAssignableProjectsResponse{
		AssignableProjects: []ScaLicenseProfileAssignableProject{{ProjectKey: "my-project"}},
		Page:               PageResponseV2{Total: 1},
	}
	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/sca/license-profiles/assignable-projects", http.StatusOK,
		map[string]string{"licenseProfileId": "lp-1", "q": "my"}, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.ListLicenseProfileAssignableProjects(context.Background(), &ScaLicenseProfileAssignableProjectsOptions{
		LicenseProfileId: "lp-1",
		Q:                "my",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, result.AssignableProjects, 1)
	assert.Equal(t, "my-project", result.AssignableProjects[0].ProjectKey)
}

func TestScaService_ListLicenseProfileAssignableProjects_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.ListLicenseProfileAssignableProjects(context.Background(), &ScaLicenseProfileAssignableProjectsOptions{
		PaginationParamsV2: PaginationParamsV2{PageIndex: -1},
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_AssignLicenseProfileProject(t *testing.T) {
	body := &ScaLicenseProfileAssignmentRequest{ProjectKey: "my-project", LicenseProfileId: "lp-1"}
	response := ScaLicenseProfile{Id: "lp-1", Name: "Default"}
	server := newTestServer(t, mockPatchHandler(t, "/v2/sca/license-profiles/assigned-projects", http.StatusOK, body, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.AssignLicenseProfileProject(context.Background(), body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "lp-1", result.Id)
}

func TestScaService_AssignLicenseProfileProject_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.AssignLicenseProfileProject(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.AssignLicenseProfileProject(context.Background(), &ScaLicenseProfileAssignmentRequest{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_DeleteLicenseProfileAssignment(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodDelete, "/v2/sca/license-profiles/assigned-projects/my-project", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.V2.Sca.DeleteLicenseProfileAssignment(context.Background(), &ScaLicenseProfileAssignmentDeleteOptions{ProjectKey: "my-project"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestScaService_DeleteLicenseProfileAssignment_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.V2.Sca.DeleteLicenseProfileAssignment(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.V2.Sca.DeleteLicenseProfileAssignment(context.Background(), &ScaLicenseProfileAssignmentDeleteOptions{})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestScaService_GetLicenseProfile(t *testing.T) {
	response := ScaLicenseProfileDetails{
		Profile:    ScaLicenseProfile{Id: "lp-1", Name: "Default"},
		Categories: []ScaLicenseProfileCategory{{Id: "cat-1", Key: ScaLicenseCategoryCopyleftStrong, Policy: ScaLicensePolicyDeny}},
		Licenses:   []ScaLicensePolicyLicense{{Id: "lic-1", SpdxLicenseId: "MIT", Policy: ScaLicensePolicyAllow}},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/sca/license-profiles/lp-1", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.GetLicenseProfile(context.Background(), &ScaLicenseProfileGetOptions{Key: "lp-1"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "lp-1", result.Profile.Id)
	require.Len(t, result.Categories, 1)
	assert.Equal(t, ScaLicenseCategoryCopyleftStrong, result.Categories[0].Key)
	require.Len(t, result.Licenses, 1)
	assert.Equal(t, "MIT", result.Licenses[0].SpdxLicenseId)
}

func TestScaService_GetLicenseProfile_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.GetLicenseProfile(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.GetLicenseProfile(context.Background(), &ScaLicenseProfileGetOptions{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_DeleteLicenseProfile(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodDelete, "/v2/sca/license-profiles/lp-1", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.V2.Sca.DeleteLicenseProfile(context.Background(), &ScaLicenseProfileGetOptions{Key: "lp-1"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestScaService_DeleteLicenseProfile_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.V2.Sca.DeleteLicenseProfile(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.V2.Sca.DeleteLicenseProfile(context.Background(), &ScaLicenseProfileGetOptions{})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestScaService_UpdateLicenseProfile(t *testing.T) {
	body := &ScaLicenseProfileUpdateRequest{Name: "Renamed"}
	response := ScaLicenseProfile{Id: "lp-1", Name: "Renamed"}
	server := newTestServer(t, mockPatchHandler(t, "/v2/sca/license-profiles/lp-1", http.StatusOK, body, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.UpdateLicenseProfile(context.Background(), "lp-1", body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Renamed", result.Name)
}

func TestScaService_UpdateLicenseProfile_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.UpdateLicenseProfile(context.Background(), "", &ScaLicenseProfileUpdateRequest{Name: "Renamed"})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.UpdateLicenseProfile(context.Background(), "lp-1", nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_UpdateLicenseProfileCategory(t *testing.T) {
	body := &ScaLicenseProfileCategoryUpdateRequest{Policy: ScaLicensePolicyDeny}
	response := ScaLicenseProfileCategory{Id: "cat-1", Key: ScaLicenseCategoryCopyleftStrong, Policy: ScaLicensePolicyDeny}
	server := newTestServer(t, mockPatchHandler(t, "/v2/sca/license-profiles/lp-1/categories/COPYLEFT_STRONG", http.StatusOK, body, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.UpdateLicenseProfileCategory(context.Background(), &ScaLicenseProfileCategoryOptions{
		LicenseProfileKey: "lp-1",
		CategoryKey:       ScaLicenseCategoryCopyleftStrong,
	}, body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, ScaLicensePolicyDeny, result.Policy)
}

func TestScaService_UpdateLicenseProfileCategory_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.UpdateLicenseProfileCategory(context.Background(), nil, &ScaLicenseProfileCategoryUpdateRequest{Policy: ScaLicensePolicyDeny})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.UpdateLicenseProfileCategory(context.Background(), &ScaLicenseProfileCategoryOptions{
		LicenseProfileKey: "lp-1",
	}, &ScaLicenseProfileCategoryUpdateRequest{Policy: ScaLicensePolicyDeny})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.UpdateLicenseProfileCategory(context.Background(), &ScaLicenseProfileCategoryOptions{
		LicenseProfileKey: "lp-1",
		CategoryKey:       "NOT_A_REAL_CATEGORY",
	}, &ScaLicenseProfileCategoryUpdateRequest{Policy: ScaLicensePolicyDeny})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.UpdateLicenseProfileCategory(context.Background(), &ScaLicenseProfileCategoryOptions{
		LicenseProfileKey: "lp-1",
		CategoryKey:       ScaLicenseCategoryCopyleftStrong,
	}, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.UpdateLicenseProfileCategory(context.Background(), &ScaLicenseProfileCategoryOptions{
		LicenseProfileKey: "lp-1",
		CategoryKey:       ScaLicenseCategoryCopyleftStrong,
	}, &ScaLicenseProfileCategoryUpdateRequest{Policy: "NOT_A_REAL_POLICY"})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_UpdateLicenseProfileLicense(t *testing.T) {
	body := &ScaLicensePolicyLicenseUpdateRequest{Policy: ScaLicensePolicyAllow}
	response := ScaLicensePolicyLicense{Id: "lic-1", Policy: ScaLicensePolicyAllow}
	server := newTestServer(t, mockPatchHandler(t, "/v2/sca/license-profiles/lp-1/licenses/lic-1", http.StatusOK, body, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.UpdateLicenseProfileLicense(context.Background(), &ScaLicenseProfileLicenseOptions{
		LicenseProfileKey: "lp-1",
		LicensePolicyId:   "lic-1",
	}, body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, ScaLicensePolicyAllow, result.Policy)
}

func TestScaService_UpdateLicenseProfileLicense_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.UpdateLicenseProfileLicense(context.Background(), nil, &ScaLicensePolicyLicenseUpdateRequest{Policy: ScaLicensePolicyAllow})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.UpdateLicenseProfileLicense(context.Background(), &ScaLicenseProfileLicenseOptions{
		LicenseProfileKey: "lp-1",
	}, &ScaLicensePolicyLicenseUpdateRequest{Policy: ScaLicensePolicyAllow})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.UpdateLicenseProfileLicense(context.Background(), &ScaLicenseProfileLicenseOptions{
		LicenseProfileKey: "lp-1",
		LicensePolicyId:   "lic-1",
	}, nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.UpdateLicenseProfileLicense(context.Background(), &ScaLicenseProfileLicenseOptions{
		LicenseProfileKey: "lp-1",
		LicensePolicyId:   "lic-1",
	}, &ScaLicensePolicyLicenseUpdateRequest{Policy: "NOT_A_REAL_POLICY"})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_GetAnalysis(t *testing.T) {
	response := ScaAnalysisResource{Status: ScaAnalysisStatusCompleted}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/sca/analyses", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.GetAnalysis(context.Background(), &ScaAnalysisGetOptions{ProjectKey: "my-project"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, ScaAnalysisStatusCompleted, result.Status)
}

func TestScaService_GetAnalysis_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.GetAnalysis(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.GetAnalysis(context.Background(), &ScaAnalysisGetOptions{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_IsScaEnabled(t *testing.T) {
	response := ScaFeatureEnabledResult{Enabled: true}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/sca/enabled", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.IsScaEnabled(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, result.Enabled)
}

func TestScaService_GetFeatureEnabled(t *testing.T) {
	response := ScaFeatureEnabledResult{Enabled: true}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/sca/feature-enabled", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.GetFeatureEnabled(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, result.Enabled)
}

func TestScaService_SelfTest(t *testing.T) {
	response := ScaSelfTestResult{FeatureEnabled: true, SelfTestPassed: true}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/sca/self-test", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.SelfTest(context.Background())
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, result.SelfTestPassed)
}

func TestScaService_GetAllAssignees(t *testing.T) {
	response := []ScaUserResource{{Login: "jdoe", Name: "Jane Doe"}}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/sca/issues-releases/all-assignees", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.GetAllAssignees(context.Background(), &ScaAllAssigneesOptions{ProjectKey: "my-project"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result, 1)
	assert.Equal(t, "jdoe", result[0].Login)
}

func TestScaService_GetAllAssignees_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.GetAllAssignees(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.GetAllAssignees(context.Background(), &ScaAllAssigneesOptions{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_UpdateAssignee(t *testing.T) {
	body := &ScaUpdateAssigneeRequest{IssueReleaseKey: "risk-1", AssigneeLogin: "jdoe"}
	response := ScaIssueReleaseDetails{Key: "risk-1", Assignee: ScaUserResource{Login: "jdoe"}}
	server := newTestServer(t, mockPatchHandler(t, "/v2/sca/issues-releases/update-assignee", http.StatusOK, body, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.UpdateAssignee(context.Background(), body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "jdoe", result.Assignee.Login)
}

func TestScaService_UpdateAssignee_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.UpdateAssignee(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.UpdateAssignee(context.Background(), &ScaUpdateAssigneeRequest{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_ChangeStatus(t *testing.T) {
	body := &ScaChangeStatusRequest{IssueReleaseKey: "risk-1", TransitionKey: ScaTransitionConfirm}
	response := ScaIssueReleaseDetails{Key: "risk-1", Status: "CONFIRMED"}
	server := newTestServer(t, mockPatchHandler(t, "/v2/sca/issues-releases/change-status", http.StatusOK, body, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.ChangeStatus(context.Background(), body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "CONFIRMED", result.Status)
}

func TestScaService_ChangeStatus_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.ChangeStatus(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.ChangeStatus(context.Background(), &ScaChangeStatusRequest{IssueReleaseKey: "risk-1"})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.ChangeStatus(context.Background(), &ScaChangeStatusRequest{
		IssueReleaseKey: "risk-1",
		TransitionKey:   "NOT_A_REAL_TRANSITION",
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_SetSeverity(t *testing.T) {
	body := &ScaSetSeverityRequest{IssueReleaseKey: "risk-1", Severity: RuleImpactSeverityHigh}
	response := ScaIssueReleaseDetails{Key: "risk-1", ManualSeverity: RuleImpactSeverityHigh}
	server := newTestServer(t, mockHandler(t, http.MethodPost, "/v2/sca/issues-releases/set-severity", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.SetSeverity(context.Background(), body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, RuleImpactSeverityHigh, result.ManualSeverity)
}

func TestScaService_SetSeverity_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.SetSeverity(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.SetSeverity(context.Background(), &ScaSetSeverityRequest{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.SetSeverity(context.Background(), &ScaSetSeverityRequest{
		IssueReleaseKey: "risk-1",
		Severity:        "NOT_A_REAL_SEVERITY",
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_UpdateIssueReleaseSeverity(t *testing.T) {
	body := &ScaSetSeverityRequest{Severity: RuleImpactSeverityHigh}
	response := ScaIssueReleaseDetails{Key: "risk-1", ManualSeverity: RuleImpactSeverityHigh}
	server := newTestServer(t, mockPatchHandler(t, "/v2/sca/issues-releases/risk-1", http.StatusOK, body, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.UpdateIssueReleaseSeverity(context.Background(), "risk-1", body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, RuleImpactSeverityHigh, result.ManualSeverity)
}

func TestScaService_UpdateIssueReleaseSeverity_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.UpdateIssueReleaseSeverity(context.Background(), "", &ScaSetSeverityRequest{Severity: RuleImpactSeverityHigh})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.UpdateIssueReleaseSeverity(context.Background(), "risk-1", nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.UpdateIssueReleaseSeverity(context.Background(), "risk-1", &ScaSetSeverityRequest{Severity: "NOT_A_REAL_SEVERITY"})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_AddComment(t *testing.T) {
	body := &ScaAddCommentRequest{IssueReleaseKey: "risk-1", Comment: "hello"}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/sca/issues-releases/comments", http.StatusNoContent, body, nil))
	client := newTestClient(t, server.URL)

	resp, err := client.V2.Sca.AddComment(context.Background(), body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestScaService_AddComment_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.V2.Sca.AddComment(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.V2.Sca.AddComment(context.Background(), &ScaAddCommentRequest{})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestScaService_GetChangelog(t *testing.T) {
	response := ScaIssueReleaseChangelog{Changelog: []ScaIssueReleaseChange{{Key: "chg-1"}}}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/sca/issues-releases/risk-1/changelogs", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.GetChangelog(context.Background(), &ScaChangelogGetOptions{Key: "risk-1"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Changelog, 1)
}

func TestScaService_GetChangelog_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.GetChangelog(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.GetChangelog(context.Background(), &ScaChangelogGetOptions{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_DeleteChangelogEntry(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodDelete, "/v2/sca/issues-releases/risk-1/changelog", http.StatusNoContent))
	client := newTestClient(t, server.URL)

	resp, err := client.V2.Sca.DeleteChangelogEntry(context.Background(), &ScaChangelogDeleteOptions{
		Key:                   "risk-1",
		IssueReleaseChangeKey: "chg-1",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestScaService_DeleteChangelogEntry_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.V2.Sca.DeleteChangelogEntry(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.V2.Sca.DeleteChangelogEntry(context.Background(), &ScaChangelogDeleteOptions{Key: "risk-1"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestScaService_UpdateChangelogComment(t *testing.T) {
	body := &ScaChangelogUpdateRequest{IssueReleaseChangeKey: "chg-1", Comment: "updated"}
	server := newTestServer(t, mockPatchHandler(t, "/v2/sca/issues-releases/risk-1/changelog", http.StatusNoContent, body, nil))
	client := newTestClient(t, server.URL)

	resp, err := client.V2.Sca.UpdateChangelogComment(context.Background(), &ScaChangelogUpdateOptions{Key: "risk-1"}, body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestScaService_UpdateChangelogComment_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.V2.Sca.UpdateChangelogComment(context.Background(), nil, &ScaChangelogUpdateRequest{IssueReleaseChangeKey: "chg-1", Comment: "x"})
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.V2.Sca.UpdateChangelogComment(context.Background(), &ScaChangelogUpdateOptions{Key: "risk-1"}, nil)
	assert.Error(t, err)
	assert.Nil(t, resp)

	resp, err = client.V2.Sca.UpdateChangelogComment(context.Background(), &ScaChangelogUpdateOptions{Key: "risk-1"}, &ScaChangelogUpdateRequest{IssueReleaseChangeKey: "chg-1"})
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestScaService_ResearchReleases(t *testing.T) {
	body := &ScaReleaseResearchOptions{ProjectKey: "my-project", Purls: []string{"pkg:npm/lodash@4.17.21"}}
	response := ScaReleaseResearch{Releases: []ScaReleaseResearchEntry{{Key: "rel-1", PackageName: "lodash"}}}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/sca/release-research/releases", http.StatusOK, body, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.ResearchReleases(context.Background(), body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Releases, 1)
}

func TestScaService_ResearchReleases_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.ResearchReleases(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.ResearchReleases(context.Background(), &ScaReleaseResearchOptions{ProjectKey: "my-project"})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_TriggerReanalysis(t *testing.T) {
	body := &ScaReanalysisOptions{ProjectKey: "my-project"}
	response := ScaReanalysisResult{BranchesQueued: 2}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/sca/reanalysis", http.StatusOK, body, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.TriggerReanalysis(context.Background(), body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int32(2), result.BranchesQueued)
}

func TestScaService_TriggerReanalysis_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.TriggerReanalysis(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.TriggerReanalysis(context.Background(), &ScaReanalysisOptions{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_GetPackageInfo(t *testing.T) {
	body := &ScaPackageInfoOptions{ProjectKey: "my-project", Purls: []string{"pkg:npm/lodash@4.17.21"}}
	response := ScaPackageInfo{Packages: []ScaPackageInfoEntry{{PackageUrl: "pkg:npm/lodash@4.17.21"}}}
	server := newTestServer(t, mockJSONBodyHandler(t, http.MethodPost, "/v2/sca/package-info", http.StatusOK, body, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.GetPackageInfo(context.Background(), body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Packages, 1)
}

func TestScaService_GetPackageInfo_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.GetPackageInfo(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.GetPackageInfo(context.Background(), &ScaPackageInfoOptions{ProjectKey: "my-project"})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_SearchReleasesByPurl(t *testing.T) {
	response := ScaReleaseSearchByPurl{
		Branches: []ScaReleaseSearchBranch{{Key: "main", PackageUrl: "pkg:npm/lodash@4.17.21"}},
		Page:     PageResponseV2{Total: 1},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/sca/releases/search", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.SearchReleasesByPurl(context.Background(), &ScaReleaseSearchByPurlOptions{Purl: "pkg:npm/lodash@4.17.21"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Branches, 1)
}

func TestScaService_SearchReleasesByPurl_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.SearchReleasesByPurl(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.SearchReleasesByPurl(context.Background(), &ScaReleaseSearchByPurlOptions{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}

func TestScaService_GetRiskReport(t *testing.T) {
	response := []ScaRiskReportItem{{ProjectKey: "my-project", RiskType: "VULNERABILITY", RiskSeverity: RuleImpactSeverityHigh}}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/sca/risk-reports", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.GetRiskReport(context.Background(), &ScaRiskReportOptions{Component: "my-project"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result, 1)
}

func TestScaService_GetRiskReport_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.GetRiskReport(context.Background(), nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.GetRiskReport(context.Background(), &ScaRiskReportOptions{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}
