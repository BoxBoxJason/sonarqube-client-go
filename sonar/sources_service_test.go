package sonar

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// SourcesService Test Suite
// -----------------------------------------------------------------------------

// TestSourcesService_Index tests the Index method.
func TestSourcesService_Index(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/sources/index", http.StatusOK, &SourcesIndex{
		Sources: map[string]string{
			"1": "package main",
			"2": "",
			"3": "import \"fmt\"",
		},
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &SourcesIndexOption{
		Resource: "my-project:src/main.go",
		From:     1,
		To:       10,
	}

	result, resp, err := client.Sources.Index(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "package main", result.Sources["1"])
}

// TestSourcesService_Index_ValidationError tests validation for Index.
func TestSourcesService_Index_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test missing Resource
	opt := &SourcesIndexOption{}
	_, _, err := client.Sources.Index(opt)
	assert.Error(t, err)
}

// TestSourcesService_IssueSnippets tests the IssueSnippets method.
func TestSourcesService_IssueSnippets(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/sources/issue_snippets", http.StatusOK, &SourcesIssueSnippets{
		"my-project:src/main.go": {
			Component: SourcesComponent{
				Key:       "my-project:src/main.go",
				Name:      "main.go",
				Qualifier: "FIL",
			},
			Sources: []SourcesLine{
				{Line: 10, Code: "func main() {"},
			},
		},
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &SourcesIssueSnippetsOption{
		IssueKey: "AX1234567890",
	}

	result, resp, err := client.Sources.IssueSnippets(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	snippet, ok := (*result)["my-project:src/main.go"]
	assert.True(t, ok)
	assert.Equal(t, "my-project:src/main.go", snippet.Component.Key)
}

// TestSourcesService_IssueSnippets_ValidationError tests validation for IssueSnippets.
func TestSourcesService_IssueSnippets_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test missing IssueKey
	opt := &SourcesIssueSnippetsOption{}
	_, _, err := client.Sources.IssueSnippets(opt)
	assert.Error(t, err)
}

// TestSourcesService_Lines tests the Lines method.
func TestSourcesService_Lines(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/sources/lines", http.StatusOK, &SourcesLines{
		Sources: []SourcesLine{
			{
				Line:        1,
				Code:        "<span class=\"k\">package</span> main",
				SCMAuthor:   "john.doe@example.com",
				SCMDate:     "2024-01-15T10:30:00+0000",
				SCMRevision: "abc123",
				Duplicated:  false,
			},
		},
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &SourcesLinesOption{
		Key:    "my-project:src/main.go",
		Branch: "main",
		From:   1,
		To:     100,
	}

	result, resp, err := client.Sources.Lines(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Sources, 1)
	assert.Equal(t, int64(1), result.Sources[0].Line)
	assert.Equal(t, "john.doe@example.com", result.Sources[0].SCMAuthor)
}

// TestSourcesService_Lines_ValidationError tests validation for Lines.
func TestSourcesService_Lines_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test missing Key
	opt := &SourcesLinesOption{
		Branch: "main",
	}
	_, _, err := client.Sources.Lines(opt)
	assert.Error(t, err)
}

// TestSourcesService_Raw tests the Raw method.
func TestSourcesService_Raw(t *testing.T) {
	expectedContent := "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}\n"
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/sources/raw")
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(expectedContent))
	})
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &SourcesRawOption{
		Key: "my-project:src/main.go",
	}

	result, resp, err := client.Sources.Raw(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedContent, result)
}

// TestSourcesService_Raw_ValidationError tests validation for Raw.
func TestSourcesService_Raw_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test missing Key
	opt := &SourcesRawOption{}
	_, _, err := client.Sources.Raw(opt)
	assert.Error(t, err)
}

// TestSourcesService_Scm tests the Scm method.
func TestSourcesService_Scm(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/sources/scm", http.StatusOK, &SourcesScm{
		Scm: [][]any{
			{1, "john.doe@example.com", "2024-01-15T10:30:00+0000", "abc123"},
			{2, "jane.smith@example.com", "2024-01-14T09:00:00+0000", "def456"},
		},
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &SourcesScmOption{
		Key:           "my-project:src/main.go",
		CommitsByLine: true,
	}

	result, resp, err := client.Sources.Scm(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Scm, 2)
}

// TestSourcesService_Scm_ValidationError tests validation for Scm.
func TestSourcesService_Scm_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test missing Key
	opt := &SourcesScmOption{
		CommitsByLine: true,
	}
	_, _, err := client.Sources.Scm(opt)
	assert.Error(t, err)
}

// TestSourcesService_Show tests the Show method.
func TestSourcesService_Show(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/sources/show", http.StatusOK, &SourcesShow{
		Sources: [][]any{
			{1, "package main"},
			{2, ""},
			{3, "import \"fmt\""},
		},
	}))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &SourcesShowOption{
		Key:  "my-project:src/main.go",
		From: 1,
		To:   10,
	}

	result, resp, err := client.Sources.Show(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Len(t, result.Sources, 3)
}

// TestSourcesService_Show_ValidationError tests validation for Show.
func TestSourcesService_Show_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Test missing Key
	opt := &SourcesShowOption{}
	_, _, err := client.Sources.Show(opt)
	assert.Error(t, err)
}
