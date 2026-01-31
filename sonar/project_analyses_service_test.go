package sonargo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectAnalysesService_CreateEvent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			require.NoError(t, r.ParseForm())
			assert.Equal(t, "AU-TpxcA-iU5OvuD2FL0", r.Form.Get("analysis"))
			assert.Equal(t, "1.0", r.Form.Get("name"))
			assert.Equal(t, "VERSION", r.Form.Get("category"))
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{
				"event": {
					"analysis": "AU-TpxcA-iU5OvuD2FL0",
					"key": "AU-TpxcA-iU5OvuD2FL1",
					"category": "VERSION",
					"name": "1.0"
				}
			}`))
			require.NoError(t, err)
		})

		client := newTestClient(t, server.URL)

		result, resp, err := client.ProjectAnalyses.CreateEvent(&ProjectAnalysesCreateEventOption{
			Analysis: "AU-TpxcA-iU5OvuD2FL0",
			Category: "VERSION",
			Name:     "1.0",
		})
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		require.NotNil(t, result)
	})

	t.Run("nil option", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, _, err := client.ProjectAnalyses.CreateEvent(nil)
		assert.Error(t, err)
	})

	t.Run("missing analysis", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, _, err := client.ProjectAnalyses.CreateEvent(&ProjectAnalysesCreateEventOption{
			Name: "1.0",
		})
		assert.Error(t, err)
	})

	t.Run("missing name", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, _, err := client.ProjectAnalyses.CreateEvent(&ProjectAnalysesCreateEventOption{
			Analysis: "AU-TpxcA-iU5OvuD2FL0",
		})
		assert.Error(t, err)
	})

	t.Run("invalid category", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, _, err := client.ProjectAnalyses.CreateEvent(&ProjectAnalysesCreateEventOption{
			Analysis: "AU-TpxcA-iU5OvuD2FL0",
			Category: "INVALID",
			Name:     "1.0",
		})
		assert.Error(t, err)
	})
}

func TestProjectAnalysesService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			require.NoError(t, r.ParseForm())
			assert.Equal(t, "AU-TpxcA-iU5OvuD2FL0", r.Form.Get("analysis"))
			w.WriteHeader(http.StatusNoContent)
		})

		client := newTestClient(t, server.URL)

		resp, err := client.ProjectAnalyses.Delete(&ProjectAnalysesDeleteOption{
			Analysis: "AU-TpxcA-iU5OvuD2FL0",
		})
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("nil option", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, err := client.ProjectAnalyses.Delete(nil)
		assert.Error(t, err)
	})

	t.Run("missing analysis", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, err := client.ProjectAnalyses.Delete(&ProjectAnalysesDeleteOption{})
		assert.Error(t, err)
	})
}

func TestProjectAnalysesService_DeleteEvent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			require.NoError(t, r.ParseForm())
			assert.Equal(t, "AU-TpxcA-iU5OvuD2FL1", r.Form.Get("event"))
			w.WriteHeader(http.StatusNoContent)
		})

		client := newTestClient(t, server.URL)

		resp, err := client.ProjectAnalyses.DeleteEvent(&ProjectAnalysesDeleteEventOption{
			Event: "AU-TpxcA-iU5OvuD2FL1",
		})
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("nil option", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, err := client.ProjectAnalyses.DeleteEvent(nil)
		assert.Error(t, err)
	})

	t.Run("missing event", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, err := client.ProjectAnalyses.DeleteEvent(&ProjectAnalysesDeleteEventOption{})
		assert.Error(t, err)
	})
}

func TestProjectAnalysesService_Search(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "my-project", r.URL.Query().Get("project"))
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{
				"paging": {
					"pageIndex": 1,
					"pageSize": 100,
					"total": 1
				},
				"analyses": [
					{
						"key": "AU-TpxcA-iU5OvuD2FL0",
						"date": "2022-01-15T10:00:00+0000",
						"projectVersion": "1.0",
						"revision": "abc123",
						"events": [
							{
								"key": "AU-TpxcA-iU5OvuD2FL1",
								"category": "VERSION",
								"name": "1.0"
							}
						]
					}
				]
			}`))
			require.NoError(t, err)
		})

		client := newTestClient(t, server.URL)

		result, resp, err := client.ProjectAnalyses.Search(&ProjectAnalysesSearchOption{
			Project: "my-project",
		})
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		require.NotNil(t, result)
		assert.Len(t, result.Analyses, 1)
		assert.Equal(t, "AU-TpxcA-iU5OvuD2FL0", result.Analyses[0].Key)
		assert.Len(t, result.Analyses[0].Events, 1)
	})

	t.Run("with filters", func(t *testing.T) {
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "my-project", r.URL.Query().Get("project"))
			assert.Equal(t, "main", r.URL.Query().Get("branch"))
			assert.Equal(t, "VERSION", r.URL.Query().Get("category"))
			assert.Equal(t, "2022-01-01", r.URL.Query().Get("from"))
			assert.Equal(t, "2022-12-31", r.URL.Query().Get("to"))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"paging": {"pageIndex": 1, "pageSize": 100, "total": 0}, "analyses": []}`))
		})

		client := newTestClient(t, server.URL)

		_, _, err := client.ProjectAnalyses.Search(&ProjectAnalysesSearchOption{
			Project:  "my-project",
			Branch:   "main",
			Category: "VERSION",
			From:     "2022-01-01",
			To:       "2022-12-31",
		})
		require.NoError(t, err)
	})

	t.Run("with datetime format", func(t *testing.T) {
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"paging": {"pageIndex": 1, "pageSize": 100, "total": 0}, "analyses": []}`))
		})

		client := newTestClient(t, server.URL)

		_, _, err := client.ProjectAnalyses.Search(&ProjectAnalysesSearchOption{
			Project: "my-project",
			From:    "2022-01-01T00:00:00Z",
			To:      "2022-12-31T23:59:59Z",
		})
		require.NoError(t, err)
	})

	t.Run("nil option", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, _, err := client.ProjectAnalyses.Search(nil)
		assert.Error(t, err)
	})

	t.Run("missing project", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, _, err := client.ProjectAnalyses.Search(&ProjectAnalysesSearchOption{})
		assert.Error(t, err)
	})

	t.Run("invalid category", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, _, err := client.ProjectAnalyses.Search(&ProjectAnalysesSearchOption{
			Project:  "my-project",
			Category: "INVALID",
		})
		assert.Error(t, err)
	})

	t.Run("invalid from date", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, _, err := client.ProjectAnalyses.Search(&ProjectAnalysesSearchOption{
			Project: "my-project",
			From:    "invalid-date",
		})
		assert.Error(t, err)
	})

	t.Run("invalid to date", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, _, err := client.ProjectAnalyses.Search(&ProjectAnalysesSearchOption{
			Project: "my-project",
			To:      "invalid-date",
		})
		assert.Error(t, err)
	})
}

func TestProjectAnalysesService_SearchAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		callCount := 0
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			callCount++
			w.Header().Set("Content-Type", "application/json")
			if callCount == 1 {
				_, _ = w.Write([]byte(`{
					"paging": {"pageIndex": 1, "pageSize": 1, "total": 2},
					"analyses": [{"key": "analysis1"}]
				}`))
			} else {
				_, _ = w.Write([]byte(`{
					"paging": {"pageIndex": 2, "pageSize": 1, "total": 2},
					"analyses": [{"key": "analysis2"}]
				}`))
			}
		})

		client := newTestClient(t, server.URL)

		opt := &ProjectAnalysesSearchOption{
			Project: "my-project",
		}
		opt.PageSize = 1

		result, _, err := client.ProjectAnalyses.SearchAll(opt)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, 2, callCount)
	})

	t.Run("nil option", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, _, err := client.ProjectAnalyses.SearchAll(nil)
		assert.Error(t, err)
	})
}

func TestProjectAnalysesService_UpdateEvent(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			require.NoError(t, r.ParseForm())
			assert.Equal(t, "AU-TpxcA-iU5OvuD2FL1", r.Form.Get("event"))
			assert.Equal(t, "2.0", r.Form.Get("name"))
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(`{
				"event": {
					"key": "AU-TpxcA-iU5OvuD2FL1",
					"category": "VERSION",
					"name": "2.0"
				}
			}`))
			require.NoError(t, err)
		})

		client := newTestClient(t, server.URL)

		result, resp, err := client.ProjectAnalyses.UpdateEvent(&ProjectAnalysesUpdateEventOption{
			Event: "AU-TpxcA-iU5OvuD2FL1",
			Name:  "2.0",
		})
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		require.NotNil(t, result)
	})

	t.Run("nil option", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, _, err := client.ProjectAnalyses.UpdateEvent(nil)
		assert.Error(t, err)
	})

	t.Run("missing event", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, _, err := client.ProjectAnalyses.UpdateEvent(&ProjectAnalysesUpdateEventOption{
			Name: "2.0",
		})
		assert.Error(t, err)
	})

	t.Run("missing name", func(t *testing.T) {
		client := newLocalhostClient(t)

		_, _, err := client.ProjectAnalyses.UpdateEvent(&ProjectAnalysesUpdateEventOption{
			Event: "AU-TpxcA-iU5OvuD2FL1",
		})
		assert.Error(t, err)
	})
}

func TestProjectAnalysesService_ValidateCreateEventOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *ProjectAnalysesCreateEventOption
		wantErr bool
	}{
		{"valid minimal", &ProjectAnalysesCreateEventOption{Analysis: "a1", Name: "1.0"}, false},
		{"valid with VERSION", &ProjectAnalysesCreateEventOption{Analysis: "a1", Name: "1.0", Category: "VERSION"}, false},
		{"valid with OTHER", &ProjectAnalysesCreateEventOption{Analysis: "a1", Name: "1.0", Category: "OTHER"}, false},
		{"nil option", nil, true},
		{"missing analysis", &ProjectAnalysesCreateEventOption{Name: "1.0"}, true},
		{"missing name", &ProjectAnalysesCreateEventOption{Analysis: "a1"}, true},
		{"invalid category", &ProjectAnalysesCreateEventOption{Analysis: "a1", Name: "1.0", Category: "INVALID"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.ProjectAnalyses.ValidateCreateEventOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProjectAnalysesService_ValidateDeleteOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *ProjectAnalysesDeleteOption
		wantErr bool
	}{
		{"valid", &ProjectAnalysesDeleteOption{Analysis: "a1"}, false},
		{"nil option", nil, true},
		{"empty analysis", &ProjectAnalysesDeleteOption{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.ProjectAnalyses.ValidateDeleteOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProjectAnalysesService_ValidateDeleteEventOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *ProjectAnalysesDeleteEventOption
		wantErr bool
	}{
		{"valid", &ProjectAnalysesDeleteEventOption{Event: "e1"}, false},
		{"nil option", nil, true},
		{"empty event", &ProjectAnalysesDeleteEventOption{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.ProjectAnalyses.ValidateDeleteEventOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProjectAnalysesService_ValidateSearchOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *ProjectAnalysesSearchOption
		wantErr bool
	}{
		{"valid minimal", &ProjectAnalysesSearchOption{Project: "p1"}, false},
		{"valid with VERSION", &ProjectAnalysesSearchOption{Project: "p1", Category: "VERSION"}, false},
		{"valid with QUALITY_GATE", &ProjectAnalysesSearchOption{Project: "p1", Category: "QUALITY_GATE"}, false},
		{"valid with date", &ProjectAnalysesSearchOption{Project: "p1", From: "2022-01-01"}, false},
		{"valid with datetime", &ProjectAnalysesSearchOption{Project: "p1", From: "2022-01-01T00:00:00Z"}, false},
		{"nil option", nil, true},
		{"missing project", &ProjectAnalysesSearchOption{}, true},
		{"invalid category", &ProjectAnalysesSearchOption{Project: "p1", Category: "INVALID"}, true},
		{"invalid from", &ProjectAnalysesSearchOption{Project: "p1", From: "bad"}, true},
		{"invalid to", &ProjectAnalysesSearchOption{Project: "p1", To: "bad"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.ProjectAnalyses.ValidateSearchOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProjectAnalysesService_ValidateUpdateEventOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *ProjectAnalysesUpdateEventOption
		wantErr bool
	}{
		{"valid", &ProjectAnalysesUpdateEventOption{Event: "e1", Name: "2.0"}, false},
		{"nil option", nil, true},
		{"missing event", &ProjectAnalysesUpdateEventOption{Name: "2.0"}, true},
		{"missing name", &ProjectAnalysesUpdateEventOption{Event: "e1"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.ProjectAnalyses.ValidateUpdateEventOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsValidDate(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"2022-01-15", true},
		{"2022-12-31", true},
		{"22-01-15", false},
		{"2022/01/15", false},
		{"2022-1-15", false},
		{"2022-01-1", false},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := isValidDate(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsValidDateTime(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"2022-01-15T10:30:00Z", true},
		{"2022-12-31T23:59:59+0000", true},
		{"2022-01-15", false},
		{"invalid", false},
		{"", false},
		{"2022-01-15 10:30:00", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := isValidDateTime(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
