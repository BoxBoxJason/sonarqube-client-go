package sonar

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// GetVersion
// =============================================================================

func TestAnalysisV2_GetVersion(t *testing.T) {
	server := newTestServer(t, mockTextHandler(t, http.MethodGet, "/v2/analysis/version", http.StatusOK, "10.5.0.12345"))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Analysis.GetVersion()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "10.5.0.12345", *result)
}

// =============================================================================
// GetJresMetadata
// =============================================================================

func TestAnalysisV2_GetJresMetadata(t *testing.T) {
	response := []AnalysisJre{
		{Id: "jre-1", Os: "linux", Arch: "x64", Filename: "jre-linux-x64.tar.gz", Sha256: "abc123", JavaPath: "bin/java"},
		{Id: "jre-2", Os: "macos", Arch: "aarch64", Filename: "jre-macos-aarch64.tar.gz", Sha256: "def456", JavaPath: "bin/java"},
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/analysis/jres", http.StatusOK, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Analysis.GetJresMetadata(nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result, 2)
	assert.Equal(t, "linux", result[0].Os)
	assert.Equal(t, "aarch64", result[1].Arch)
}

func TestAnalysisV2_GetJresMetadata_WithFilter(t *testing.T) {
	response := []AnalysisJre{
		{Id: "jre-1", Os: "linux", Arch: "x64"},
	}
	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/analysis/jres", http.StatusOK,
		map[string]string{"os": "linux", "arch": "x64"},
		response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Analysis.GetJresMetadata(&AnalysisJresOptions{
		Os:   "linux",
		Arch: "x64",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result, 1)
}

// =============================================================================
// DownloadJre
// =============================================================================

func TestAnalysisV2_DownloadJre(t *testing.T) {
	binaryData := []byte("fake-jre-binary-data")
	server := newTestServer(t, mockBinaryHandler(t, http.MethodGet, "/v2/analysis/jres/jre-1", http.StatusOK,
		"application/octet-stream", binaryData))
	client := newTestClient(t, server.url())

	var buf bytes.Buffer
	resp, err := client.V2.Analysis.DownloadJre("jre-1", &buf)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, binaryData, buf.Bytes())
}

func TestAnalysisV2_DownloadJre_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name   string
		jreID  string
		writer io.Writer
	}{
		{name: "missing id", jreID: "", writer: &bytes.Buffer{}},
		{name: "nil writer", jreID: "jre-1", writer: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := client.V2.Analysis.DownloadJre(tt.jreID, tt.writer)
			assert.Error(t, err)
			assert.Nil(t, resp)
		})
	}
}

// =============================================================================
// GetJreMetadata
// =============================================================================

func TestAnalysisV2_GetJreMetadata(t *testing.T) {
	response := AnalysisJre{
		Id:       "jre-1",
		Os:       "linux",
		Arch:     "x64",
		Filename: "jre-linux-x64.tar.gz",
		Sha256:   "abc123",
		JavaPath: "bin/java",
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/analysis/jres/jre-1", http.StatusOK, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Analysis.GetJreMetadata("jre-1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "jre-1", result.Id)
	assert.Equal(t, "linux", result.Os)
}

func TestAnalysisV2_GetJreMetadata_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Analysis.GetJreMetadata("")
	assert.Error(t, err)
}

// =============================================================================
// DownloadScannerEngine
// =============================================================================

func TestAnalysisV2_DownloadScannerEngine(t *testing.T) {
	binaryData := []byte("fake-scanner-engine-binary")
	server := newTestServer(t, mockBinaryHandler(t, http.MethodGet, "/v2/analysis/engine", http.StatusOK,
		"application/octet-stream", binaryData))
	client := newTestClient(t, server.url())

	var buf bytes.Buffer
	resp, err := client.V2.Analysis.DownloadScannerEngine(&buf)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, binaryData, buf.Bytes())
}

func TestAnalysisV2_DownloadScannerEngine_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	resp, err := client.V2.Analysis.DownloadScannerEngine(nil)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

// =============================================================================
// GetScannerEngineMetadata
// =============================================================================

func TestAnalysisV2_GetScannerEngineMetadata(t *testing.T) {
	response := AnalysisEngineInfo{
		Filename: "scanner-engine.jar",
		Sha256:   "engine-sha256",
	}
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/analysis/engine", http.StatusOK, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Analysis.GetScannerEngineMetadata()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "scanner-engine.jar", result.Filename)
	assert.Equal(t, "engine-sha256", result.Sha256)
}

// =============================================================================
// GetActiveRules
// =============================================================================

func TestAnalysisV2_GetActiveRules(t *testing.T) {
	response := []AnalysisActiveRule{
		{
			RuleKey:     AnalysisRuleKey{Repository: "java", Rule: "S1234"},
			Name:        "Some Rule",
			Severity:    "MAJOR",
			Language:    "java",
			QProfileKey: "qp-1",
			Params:      []AnalysisParam{{Key: "max", Value: "10"}},
			Impacts:     map[string]string{"MAINTAINABILITY": "HIGH"},
		},
	}
	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/analysis/active_rules", http.StatusOK,
		map[string]string{"projectKey": "my-project"},
		response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Analysis.GetActiveRules(&AnalysisActiveRuleOptions{
		ProjectKey: "my-project",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result, 1)
	assert.Equal(t, "java", result[0].RuleKey.Repository)
	assert.Equal(t, "S1234", result[0].RuleKey.Rule)
	assert.Equal(t, "MAJOR", result[0].Severity)
}

func TestAnalysisV2_GetActiveRules_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *AnalysisActiveRuleOptions
	}{
		{"nil opt", nil},
		{"missing project key", &AnalysisActiveRuleOptions{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.V2.Analysis.GetActiveRules(tt.opt)
			assert.Error(t, err)
		})
	}
}
