package sonar_test

import (
	"log"
	"net/http"
	"time"

	"github.com/boxboxjason/sonarqube-client-go/v2/sonar"
)

// roundTripperFunc adapts a function to http.RoundTripper.
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// Attach a request-logging middleware to the client transport. Middleware sits
// inside the retry layer, so it observes every individual attempt. The same
// pattern wraps tracing/metrics transports such as otelhttp.NewTransport.
func ExampleWithMiddleware() {
	loggingMiddleware := func(next http.RoundTripper) http.RoundTripper {
		return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			start := time.Now()
			resp, err := next.RoundTrip(req)
			log.Printf("sonar %s %s (%s)", req.Method, req.URL.Path, time.Since(start))

			return resp, err
		})
	}

	client, err := sonar.NewClient(nil,
		sonar.WithBaseURL("https://sonar.example.com/api/"),
		sonar.WithToken("my-token"),
		sonar.WithMiddleware(loggingMiddleware),
	)
	if err != nil {
		panic(err)
	}

	_ = client
}
