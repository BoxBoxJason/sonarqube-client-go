package sonar

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDuplications_Show(t *testing.T) {
	response := `{
		"duplications": [
			{
				"blocks": [
					{"_ref": "1", "from": 10, "size": 5},
					{"_ref": "2", "from": 20, "size": 5}
				]
			}
		],
		"files": {
			"1": {"key": "com.example:MyFile.java", "name": "MyFile.java", "projectName": "My Project"},
			"2": {"key": "com.example:OtherFile.java", "name": "OtherFile.java", "projectName": "My Project"}
		}
	}`
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/duplications/show", http.StatusOK, response))
	client := newTestClient(t, server.url())

	opt := &DuplicationsShowOption{
		Key: "com.example:MyFile.java",
	}

	result, resp, err := client.Duplications.Show(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Duplications, 1)
	assert.Len(t, result.Duplications[0].Blocks, 2)
	assert.Equal(t, "1", result.Duplications[0].Blocks[0].Ref)
	assert.Equal(t, int64(10), result.Duplications[0].Blocks[0].From)
	assert.Len(t, result.Files, 2)

	file, ok := result.Files["1"]
	require.True(t, ok, "expected file with key '1'")
	assert.Equal(t, "MyFile.java", file.Name)
}

func TestDuplications_Show_WithBranch(t *testing.T) {
	response := `{"duplications": [], "files": {}}`
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/duplications/show", http.StatusOK, response))
	client := newTestClient(t, server.url())

	opt := &DuplicationsShowOption{
		Key:    "com.example:MyFile.java",
		Branch: "feature",
	}

	_, resp, err := client.Duplications.Show(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDuplications_Show_WithPullRequest(t *testing.T) {
	response := `{"duplications": [], "files": {}}`
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/duplications/show", http.StatusOK, response))
	client := newTestClient(t, server.url())

	opt := &DuplicationsShowOption{
		Key:         "com.example:MyFile.java",
		PullRequest: "123",
	}

	_, resp, err := client.Duplications.Show(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDuplications_Show_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *DuplicationsShowOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing Key",
			opt:  &DuplicationsShowOption{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.Duplications.Show(tt.opt)
			require.Error(t, err)
		})
	}
}

func TestDuplications_ValidateShowOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *DuplicationsShowOption
		wantErr bool
	}{
		{
			name: "valid option",
			opt: &DuplicationsShowOption{
				Key: "com.example:MyFile.java",
			},
			wantErr: false,
		},
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
		},
		{
			name:    "missing Key",
			opt:     &DuplicationsShowOption{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Duplications.ValidateShowOpt(tt.opt)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
