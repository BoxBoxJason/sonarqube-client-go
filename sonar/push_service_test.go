package sonar

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPush_SonarlintEvents(t *testing.T) {
	handler := mockHandler(t, http.MethodGet, "/push/sonarlint_events", http.StatusOK, `{}`)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &PushSonarlintEventsOptions{
		Languages:   []string{"java", "go"},
		ProjectKeys: []string{"my-project"},
	}

	resp, err := client.Push.SonarlintEvents(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPush_SonarlintEvents_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name      string
		opt       *PushSonarlintEventsOptions
		wantField string
	}{
		{"nil option", nil, "opt"},
		{"missing languages", &PushSonarlintEventsOptions{ProjectKeys: []string{"my-project"}}, "Languages"},
		{"missing project keys", &PushSonarlintEventsOptions{Languages: []string{"java"}}, "ProjectKeys"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Push.SonarlintEvents(tt.opt)
			require.Error(t, err)

			var validationErr *ValidationError
			require.True(t, errors.As(err, &validationErr))
			assert.Equal(t, tt.wantField, validationErr.Field)
		})
	}
}

func TestPush_ValidateSonarlintEventsOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name      string
		opt       *PushSonarlintEventsOptions
		wantErr   bool
		wantField string
	}{
		{"nil option", nil, true, "opt"},
		{"missing languages", &PushSonarlintEventsOptions{ProjectKeys: []string{"my-project"}}, true, "Languages"},
		{"missing project keys", &PushSonarlintEventsOptions{Languages: []string{"java"}}, true, "ProjectKeys"},
		{"valid option", &PushSonarlintEventsOptions{Languages: []string{"java", "go"}, ProjectKeys: []string{"my-project"}}, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Push.ValidateSonarlintEventsOpt(tt.opt)
			if tt.wantErr {
				require.Error(t, err)
				var validationErr *ValidationError
				if errors.As(err, &validationErr) {
					assert.Equal(t, tt.wantField, validationErr.Field)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
