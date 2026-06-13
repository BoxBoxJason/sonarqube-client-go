package sonar

import (
	"net/http"
	"testing"
	"time"
)

func TestNewClient_DefaultTimeout(t *testing.T) {
	t.Parallel()

	client, err := NewClient(nil)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	if client.httpClient.Timeout != defaultHTTPTimeout {
		t.Errorf("expected default timeout %s, got %s", defaultHTTPTimeout, client.httpClient.Timeout)
	}
}

func TestWithTimeout(t *testing.T) {
	t.Parallel()

	const custom = 5 * time.Second

	client, err := NewClient(nil, WithTimeout(custom))
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	if client.httpClient.Timeout != custom {
		t.Errorf("expected timeout %s, got %s", custom, client.httpClient.Timeout)
	}
}

func TestWithTimeout_Negative(t *testing.T) {
	t.Parallel()

	_, err := NewClient(nil, WithTimeout(-time.Second))
	if err == nil {
		t.Fatal("expected error for negative timeout, got nil")
	}
}

func TestNewClient_TransportConfigKeepsDefaultTimeout(t *testing.T) {
	t.Parallel()

	client, err := NewClient(nil, WithTransportConfig(TransportConfig{MaxIdleConns: 10})) //nolint:exhaustruct
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	if client.httpClient.Timeout != defaultHTTPTimeout {
		t.Errorf("expected default timeout %s, got %s", defaultHTTPTimeout, client.httpClient.Timeout)
	}

	if client.httpClient.Transport == nil {
		t.Error("expected transport to be configured from TransportConfig")
	}
}

func TestClientCreateOptions_Timeout(t *testing.T) {
	t.Parallel()

	custom := 7 * time.Second

	client, err := NewClient(&ClientCreateOptions{Timeout: &custom}) //nolint:exhaustruct
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	if client.httpClient.Timeout != custom {
		t.Errorf("expected timeout %s, got %s", custom, client.httpClient.Timeout)
	}
}

func TestClientCreateOptions_Timeout_Negative(t *testing.T) {
	t.Parallel()

	neg := -time.Second

	_, err := NewClient(&ClientCreateOptions{Timeout: &neg}) //nolint:exhaustruct
	if err == nil {
		t.Fatal("expected error for negative timeout, got nil")
	}
}

func TestWithHTTPClient_TimeoutNotOverridden(t *testing.T) {
	t.Parallel()

	const custom = 2 * time.Second
	//nolint:exhaustruct // only Timeout is relevant for this test
	provided := &http.Client{Timeout: custom}

	client, err := NewClient(nil, WithHTTPClient(provided), WithTimeout(time.Minute))
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}

	if client.httpClient.Timeout != custom {
		t.Errorf("caller-provided client timeout should be preserved: expected %s, got %s", custom, client.httpClient.Timeout)
	}
}
