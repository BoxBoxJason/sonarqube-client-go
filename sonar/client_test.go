package sonargo

import (
	"testing"
)

func TestNewClientWithToken_DefaultBaseURL(t *testing.T) {
	c, err := NewClientWithToken("", "token123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatalf("client is nil")
	}
	if c.BaseURL().String() != defaultBaseURL {
		t.Fatalf("expected base URL %s, got %s", defaultBaseURL, c.BaseURL().String())
	}
	if c.authType != privateToken {
		t.Fatalf("expected authType privateToken")
	}

	// ensure SetBaseURL accepts a custom URL
	if u, err := SetBaseURLUtil("http://example.com/api"); err != nil {
		t.Fatalf("SetBaseURLUtil failed: %v", err)
	} else if u.String() != "http://example.com/api/" {
		t.Fatalf("unexpected parsed url: %s", u.String())
	}
}
