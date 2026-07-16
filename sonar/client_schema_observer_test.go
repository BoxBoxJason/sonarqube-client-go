package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithSchemaObserver_ReportsUnknownField(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, mockHandler(t, http.MethodGet, "/projects/search", http.StatusOK,
		`{"components":[{"key":"k","ghostField":"boo"}],"paging":{"pageIndex":1,"pageSize":10,"total":1}}`))

	var (
		gotEndpoint   string
		gotMismatches []SchemaMismatch
	)

	client, err := NewClient(nil, WithBaseURL(ts.url()), WithSchemaObserver(func(endpoint string, mismatches []SchemaMismatch) {
		gotEndpoint = endpoint
		gotMismatches = mismatches
	}))
	require.NoError(t, err)

	result, _, err := client.Projects.Search(context.Background(), &ProjectsSearchOptions{})
	require.NoError(t, err)
	require.NotNil(t, result)

	require.Len(t, gotMismatches, 1)
	assert.Equal(t, "components[].ghostField", gotMismatches[0].Path)
	assert.Contains(t, gotEndpoint, "/projects/search")
}

func TestWithSchemaObserver_NoMismatchesWhenSchemaMatches(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, mockHandler(t, http.MethodGet, "/projects/search", http.StatusOK,
		`{"components":[{"key":"k"}],"paging":{"pageIndex":1,"pageSize":10,"total":1}}`))

	called := false

	var gotMismatches []SchemaMismatch

	client, err := NewClient(nil, WithBaseURL(ts.url()), WithSchemaObserver(func(_ string, mismatches []SchemaMismatch) {
		called = true
		gotMismatches = mismatches
	}))
	require.NoError(t, err)

	_, _, err = client.Projects.Search(context.Background(), &ProjectsSearchOptions{})
	require.NoError(t, err)

	assert.True(t, called, "schema observer should have been invoked")
	assert.Empty(t, gotMismatches)
}

func TestWithSchemaObserver_NilObserverRejected(t *testing.T) {
	t.Parallel()

	_, err := NewClient(nil, WithBaseURL("http://example.com/"), WithSchemaObserver(nil))
	require.Error(t, err)
}

// TestDo_WithoutSchemaObserverBehavesLikeBefore verifies that a client with no
// schema observer configured decodes normally, ignoring unknown fields, since
// strict schema validation must stay fully opt-in.
func TestDo_WithoutSchemaObserverBehavesLikeBefore(t *testing.T) {
	t.Parallel()

	ts := newTestServer(t, mockHandler(t, http.MethodGet, "/projects/search", http.StatusOK,
		`{"components":[{"key":"k","ghostField":"boo"}]}`))

	client, err := NewClient(nil, WithBaseURL(ts.url()))
	require.NoError(t, err)

	result, _, err := client.Projects.Search(context.Background(), &ProjectsSearchOptions{})
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "k", result.Components[0].Key)
}
