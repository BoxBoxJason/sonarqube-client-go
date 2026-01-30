package sonargo

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestDoJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte("{\"foo\":\"bar\"}"))
	}))
	defer ts.Close()
	baseURL, _ := url.Parse(ts.URL + "/")
	req, err := NewRequest(http.MethodGet, "test", baseURL, "u", "p", nil)
	if err != nil {
		t.Fatalf("NewRequest failed: %v", err)
	}
	var v map[string]interface{}
	resp, err := Do(http.DefaultClient, req, &v)
	if err != nil {
		t.Fatalf("Do failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}
	if v["foo"] != "bar" {
		t.Fatalf("unexpected body: %v", v)
	}
}
