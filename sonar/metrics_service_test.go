package sonar

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetrics_Search(t *testing.T) {
	response := `{
		"metrics": [
			{
				"id": "1",
				"key": "coverage",
				"name": "Coverage",
				"description": "Code coverage percentage",
				"domain": "Coverage",
				"type": "PERCENT",
				"direction": 1,
				"qualitative": true,
				"hidden": false,
				"custom": false
			},
			{
				"id": "2",
				"key": "bugs",
				"name": "Bugs",
				"description": "Number of bugs",
				"domain": "Reliability",
				"type": "INT",
				"direction": -1,
				"qualitative": true,
				"hidden": false,
				"custom": false
			}
		],
		"paging": {
			"pageIndex": 1,
			"pageSize": 100,
			"total": 2
		}
	}`
	handler := mockHandler(t, http.MethodGet, "/metrics/search", http.StatusOK, response)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, resp, err := client.Metrics.Search(nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Metrics, 2)
	assert.Equal(t, "coverage", result.Metrics[0].Key)
	assert.Equal(t, "Coverage", result.Metrics[0].Name)
	assert.Equal(t, int64(1), result.Metrics[0].Direction)
	assert.Equal(t, int64(2), result.Paging.Total)
}

func TestMetrics_Search_WithPagination(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/metrics/search", http.StatusOK, `{"metrics": [], "paging": {"pageIndex": 2, "pageSize": 50, "total": 0}}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &MetricsSearchOption{
		PaginationArgs: PaginationArgs{
			Page:     2,
			PageSize: 50,
		},
	}

	_, resp, err := client.Metrics.Search(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestMetrics_Types(t *testing.T) {
	response := `{"types": ["INT", "FLOAT", "PERCENT", "BOOL", "STRING", "MILLISEC", "RATING", "DATA", "DISTRIB", "LEVEL", "WORK_DUR"]}`
	handler := mockHandler(t, http.MethodGet, "/metrics/types", http.StatusOK, response)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, resp, err := client.Metrics.Types()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Types, 11)
	assert.Equal(t, "INT", result.Types[0])
}

func TestMetrics_ValidateSearchOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *MetricsSearchOption
		wantErr bool
	}{
		{"nil option", nil, false},
		{"empty option", &MetricsSearchOption{}, false},
		{"valid pagination", &MetricsSearchOption{PaginationArgs: PaginationArgs{Page: 1, PageSize: 100}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Metrics.ValidateSearchOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
