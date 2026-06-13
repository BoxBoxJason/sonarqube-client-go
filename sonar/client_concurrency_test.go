package sonar

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
)

// TestClient_ConcurrentReconfigurationIsRaceFree exercises the client from many
// goroutines while its auth and base URL are reconfigured at runtime. It is
// designed to surface data races when run under `go test -race`; without the
// race detector it simply asserts that concurrent use does not panic.
func TestClient_ConcurrentReconfigurationIsRaceFree(t *testing.T) {
	t.Parallel()

	server := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})

	client := newTestClient(t, server.url())
	baseURL := server.url()

	const (
		goroutines = 8
		iterations = 50
	)

	var wg sync.WaitGroup

	wg.Add(goroutines * 2)

	// Writers continuously rotate credentials and reset the base URL.
	for i := range goroutines {
		go func(worker int) {
			defer wg.Done()

			for j := range iterations {
				client.SetPrivateToken(fmt.Sprintf("token-%d-%d", worker, j))
				client.SetBasicAuth("user", "pass")
				_ = client.SetBaseURL(&baseURL)
			}
		}(i)
	}

	// Readers continuously build and send requests, reading the same fields.
	for range goroutines {
		go func() {
			defer wg.Done()

			for range iterations {
				req, err := client.NewSonarQubeV1APIRequest(context.Background(), http.MethodGet, "ping", nil)
				if err != nil {
					t.Errorf("failed to build request: %v", err)

					return
				}

				_, _ = client.Do(req, nil)
			}
		}()
	}

	wg.Wait()
}
