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
		ProjectKey: "my-project",
		Type:       ScaSbomReportTypeCycloneDX,
		Format:     ScaSbomReportFormatJSON,
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
		ProjectKey: "my-project",
		Format:     ScaSbomReportFormatJSON,
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.GetSbomReport(context.Background(), &ScaSbomReportOptions{
		ProjectKey: "my-project",
		Type:       ScaSbomReportTypeCycloneDX,
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}
