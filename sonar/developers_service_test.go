package sonar

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDevelopers_SearchEvents(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/developers/search_events", http.StatusOK, &DevelopersSearchEvents{
		Events: []DeveloperEvent{
			{
				Category: "NEW_ISSUES",
				Link:     "https://sonar.example.com/project?id=my-project",
				Message:  "10 new issues",
				Project:  "my-project",
			},
		},
	}))

	client := newTestClient(t, server.url())

	opt := &DevelopersSearchEventsOption{
		From:     []string{"2017-10-19T13:00:00+0200"},
		Projects: []string{"my-project"},
	}

	result, resp, err := client.Developers.SearchEvents(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Events, 1)

	event := result.Events[0]
	assert.Equal(t, "NEW_ISSUES", event.Category)
	assert.Equal(t, "my-project", event.Project)
}

func TestDevelopers_SearchEvents_ValidationErrors(t *testing.T) {
	tests := []struct {
		name      string
		opt       *DevelopersSearchEventsOption
		wantField string
	}{
		{
			name:      "nil option",
			opt:       nil,
			wantField: "opt",
		},
		{
			name: "missing from",
			opt: &DevelopersSearchEventsOption{
				Projects: []string{"my-project"},
			},
			wantField: "From",
		},
		{
			name: "missing projects",
			opt: &DevelopersSearchEventsOption{
				From: []string{"2017-10-19T13:00:00+0200"},
			},
			wantField: "Projects",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newLocalhostClient(t)

			_, _, err := client.Developers.SearchEvents(tt.opt)
			require.Error(t, err)

			var validationErr *ValidationError
			require.True(t, errors.As(err, &validationErr), "expected ValidationError, got %T", err)
			assert.Equal(t, tt.wantField, validationErr.Field)
		})
	}
}

func TestDevelopers_ValidateSearchEventsOpt(t *testing.T) {
	tests := []struct {
		name      string
		opt       *DevelopersSearchEventsOption
		wantErr   bool
		wantField string
	}{
		{
			name:      "nil option",
			opt:       nil,
			wantErr:   true,
			wantField: "opt",
		},
		{
			name: "missing from",
			opt: &DevelopersSearchEventsOption{
				Projects: []string{"my-project"},
			},
			wantErr:   true,
			wantField: "From",
		},
		{
			name: "missing projects",
			opt: &DevelopersSearchEventsOption{
				From: []string{"2017-10-19T13:00:00+0200"},
			},
			wantErr:   true,
			wantField: "Projects",
		},
		{
			name: "valid option",
			opt: &DevelopersSearchEventsOption{
				From:     []string{"2017-10-19T13:00:00+0200"},
				Projects: []string{"my-project"},
			},
			wantErr: false,
		},
		{
			name: "valid with multiple values",
			opt: &DevelopersSearchEventsOption{
				From:     []string{"2017-10-19T13:00:00+0200", "2017-10-20T13:00:00+0200"},
				Projects: []string{"my-project", "other-project"},
			},
			wantErr: false,
		},
	}

	client := newLocalhostClient(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Developers.ValidateSearchEventsOpt(tt.opt)

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantField != "" {
					var validationErr *ValidationError
					if errors.As(err, &validationErr) {
						assert.Equal(t, tt.wantField, validationErr.Field)
					}
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
