package sonar

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// buildMetricItems returns a JSON array string with n metric objects.
func buildMetricItems(start, count int) string {
	parts := make([]string, count)
	for i := range parts {
		parts[i] = fmt.Sprintf(`{"key":"m%d"}`, start+i)
	}

	return "[" + strings.Join(parts, ",") + "]"
}

// TestAllPages_QueryParameters verifies that allPages sends the correct query
// parameters: ps defaults to MaxPageSize (500) when not set, and p advances by
// one for each subsequent request.
func TestAllPages_QueryParameters(t *testing.T) {
	var (
		mu      sync.Mutex
		queries []url.Values
	)

	const total = 501 // forces two pages at pageSize=500

	callCount := 0
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		queries = append(queries, r.URL.Query())
		callCount++
		local := callCount
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")

		if local == 1 {
			// First page: 500 items, 1 more remains.
			items := buildMetricItems(1, 500)
			fmt.Fprintf(w, `{"paging":{"pageIndex":1,"pageSize":500,"total":%d},"metrics":%s}`, total, items)
		} else {
			// Second page: final 1 item.
			fmt.Fprintf(w, `{"paging":{"pageIndex":2,"pageSize":500,"total":%d},"metrics":[{"key":"m501"}]}`, total)
		}
	})

	client := newTestClient(t, server.URL)

	// No page size specified — allPages must default to MaxPageSize (500).
	result, _, err := client.Metrics.SearchAll(context.Background(), &MetricsSearchOptions{})
	require.NoError(t, err)
	assert.Len(t, result, total)
	require.Len(t, queries, 2)

	// First request: must use MaxPageSize and start at page 1.
	assert.Equal(t, "500", queries[0].Get("ps"), "first request must use maximum page size")
	assert.Equal(t, "1", queries[0].Get("p"), "first request must start at page 1")

	// Second request: page size unchanged, page number advances.
	assert.Equal(t, "500", queries[1].Get("ps"), "second request must keep the same page size")
	assert.Equal(t, "2", queries[1].Get("p"), "second request must advance to page 2")
}

// TestAllPages_ContextCancellation verifies that when the context is cancelled
// between pages, allPages returns the partial results collected so far together
// with a wrapped context.Canceled error, without making another HTTP request.
//
// The fetch function is mocked directly so we avoid any race between the HTTP
// transport and context propagation.
func TestAllPages_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	callCount := 0
	page := int64(0)
	pageSize := int64(0) // will be defaulted to MaxPageSize by allPages

	fakeResp := &http.Response{StatusCode: http.StatusOK}

	fetch := func(_ context.Context) ([]string, int64, *http.Response, error) {
		callCount++
		// After delivering the first page, cancel so the next iteration sees a
		// cancelled context before issuing another request.
		cancel()

		return []string{"item1"}, 3, fakeResp, nil // total=3 requires more pages
	}

	result, resp, err := allPages(ctx, &page, &pageSize, fetch)

	// Partial results from the completed page must be present.
	require.Len(t, result, 1)
	assert.Equal(t, "item1", result[0])

	// The error must wrap context.Canceled.
	assert.ErrorIs(t, err, context.Canceled)

	// The response from the last successful fetch must be returned.
	assert.Equal(t, fakeResp, resp)

	// No second fetch must have been made.
	assert.Equal(t, 1, callCount, "must not issue a second request after cancellation")
}
