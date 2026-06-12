package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalysisReports_QueueStatus(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		expectedEmpty  bool
	}{
		{"queue is empty", "true", true},
		{"queue is not empty", "false", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := mockBinaryHandler(t, http.MethodGet, "/analysis_reports/is_queue_empty", http.StatusOK, "text/plain", []byte(tt.serverResponse))
			server := newTestServer(t, handler)
			client := newTestClient(t, server.url())

			result, resp, err := client.AnalysisReports.QueueStatus(context.Background())
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			require.NotNil(t, result)
			assert.Equal(t, tt.expectedEmpty, result.IsEmpty)
		})
	}
}
