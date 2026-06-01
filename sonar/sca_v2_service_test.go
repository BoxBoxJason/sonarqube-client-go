package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScaService_ListClis(t *testing.T) {
	response := ScaClis{Clis: []ScaCliInfo{{Id: "cli-1", Filename: "sca-cli-linux", Os: "linux", Arch: "amd64"}}}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/sca/clis", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.ListClis(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Clis, 1)
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
		Page:           ScaPageResponse{PageIndex: 1, PageSize: 20, Total: 1},
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
}

func TestScaService_SearchReleases(t *testing.T) {
	response := ScaReleasesSearch{
		Releases: []ScaReleaseSearchResource{{Key: "rel-1", PackageName: "lodash"}},
		Page:     ScaPageResponse{Total: 1},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/sca/releases", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.SearchReleases(context.Background(), &ScaReleasesSearchOptions{
		ProjectKey: "my-project",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Releases, 1)
}

func TestScaService_SearchReleases_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	result, resp, err := client.V2.Sca.SearchReleases(context.Background(), nil)
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
}

func TestScaService_GetSbomReport(t *testing.T) {
	data := []byte(`{"bomFormat":"CycloneDX","specVersion":"1.4"}`)
	server := newTestServer(t, mockBinaryHandler(t, http.MethodGet, "/v2/sca/sbom-reports", http.StatusOK, "application/json", data))
	client := newTestClient(t, server.URL)

	result, resp, err := client.V2.Sca.GetSbomReport(context.Background(), &ScaSbomReportOptions{
		ProjectKey: "my-project",
		Type:       "CYCLONEDX_1_4_JSON",
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

	result, resp, err = client.V2.Sca.GetSbomReport(context.Background(), &ScaSbomReportOptions{Type: "CYCLONEDX_1_4_JSON"})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)

	result, resp, err = client.V2.Sca.GetSbomReport(context.Background(), &ScaSbomReportOptions{ProjectKey: "my-project"})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, resp)
}
