package sonargo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestL10N_Index(t *testing.T) {
	response := `{
		"locale": "en",
		"messages": {
			"quality_gates.operator.LT": "is less than",
			"quality_gates.operator.GT": "is greater than",
			"projects.no_projects.title": "There are no projects yet",
			"projects.create_project": "Create Project"
		}
	}`
	handler := mockHandler(t, http.MethodGet, "/l10n/index", http.StatusOK, response)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, resp, err := client.L10N.Index(nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, "en", result.Locale)
	assert.Len(t, result.Messages, 4)
	assert.Equal(t, "is less than", result.Messages["quality_gates.operator.LT"])
	assert.Equal(t, "Create Project", result.Messages["projects.create_project"])
}

func TestL10N_Index_WithLocale(t *testing.T) {
	response := `{
		"locale": "fr",
		"messages": {
			"quality_gates.operator.LT": "est inférieur à",
			"quality_gates.operator.GT": "est supérieur à"
		}
	}`
	handler := mockHandler(t, http.MethodGet, "/l10n/index", http.StatusOK, response)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, resp, err := client.L10N.Index(&L10NIndexOption{Locale: "fr"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "fr", result.Locale)
	assert.Equal(t, "est inférieur à", result.Messages["quality_gates.operator.LT"])
}

func TestL10N_Index_WithTimestamp(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/l10n/index", http.StatusOK, `{"locale": "en", "messages": {}}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	_, _, err := client.L10N.Index(&L10NIndexOption{Timestamp: "2024-01-01T00:00:00+0000"})
	require.NoError(t, err)
}

func TestL10N_Index_EmptyMessages(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/l10n/index", http.StatusOK, `{"locale": "en", "messages": {}}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	result, _, err := client.L10N.Index(nil)
	require.NoError(t, err)
	require.NotNil(t, result.Messages)
	assert.Empty(t, result.Messages)
}

func TestL10N_ValidateIndexOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *L10NIndexOption
		wantErr bool
	}{
		{"nil option", nil, false},
		{"empty option", &L10NIndexOption{}, false},
		{"with Locale", &L10NIndexOption{Locale: "en"}, false},
		{"with Timestamp", &L10NIndexOption{Timestamp: "2024-01-01T00:00:00+0000"}, false},
		{"with both", &L10NIndexOption{Locale: "fr", Timestamp: "2024-01-01T00:00:00+0000"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.L10N.ValidateIndexOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
