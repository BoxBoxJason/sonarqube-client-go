package sonargo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDismissMessage_Check(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}

		if r.URL.Path != "/api/dismiss_message/check" {
			t.Errorf("expected path /api/dismiss_message/check, got %s", r.URL.Path)
		}

		messageType := r.URL.Query().Get("messageType")
		if messageType != "INFO" {
			t.Errorf("expected messageType 'INFO', got %s", messageType)
		}

		projectKey := r.URL.Query().Get("projectKey")
		if projectKey != "my-project" {
			t.Errorf("expected projectKey 'my-project', got %s", projectKey)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"dismissed": true}`))
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &DismissMessageCheckOption{
		MessageType: "INFO",
		ProjectKey:  "my-project",
	}

	result, resp, err := client.DismissMessage.Check(opt)
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if !result.Dismissed {
		t.Error("expected dismissed to be true")
	}
}

func TestDismissMessage_Check_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, _, err := client.DismissMessage.Check(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing MessageType should fail validation.
	_, _, err = client.DismissMessage.Check(&DismissMessageCheckOption{
		ProjectKey: "my-project",
	})
	if err == nil {
		t.Error("expected error for missing MessageType")
	}

	// Missing ProjectKey should fail validation.
	_, _, err = client.DismissMessage.Check(&DismissMessageCheckOption{
		MessageType: "INFO",
	})
	if err == nil {
		t.Error("expected error for missing ProjectKey")
	}
}

func TestDismissMessage_Dismiss(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/dismiss_message/dismiss" {
			t.Errorf("expected path /api/dismiss_message/dismiss, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &DismissMessageDismissOption{
		MessageType: "WARNING",
		ProjectKey:  "my-project",
	}

	resp, err := client.DismissMessage.Dismiss(opt)
	if err != nil {
		t.Fatalf("Dismiss failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestDismissMessage_Dismiss_ValidationError(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Nil option should fail validation.
	_, err := client.DismissMessage.Dismiss(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}

	// Missing MessageType should fail validation.
	_, err = client.DismissMessage.Dismiss(&DismissMessageDismissOption{
		ProjectKey: "my-project",
	})
	if err == nil {
		t.Error("expected error for missing MessageType")
	}

	// Missing ProjectKey should fail validation.
	_, err = client.DismissMessage.Dismiss(&DismissMessageDismissOption{
		MessageType: "INFO",
	})
	if err == nil {
		t.Error("expected error for missing ProjectKey")
	}
}

func TestDismissMessage_ValidateCheckOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Valid option should pass.
	err := client.DismissMessage.ValidateCheckOpt(&DismissMessageCheckOption{
		MessageType: "INFO",
		ProjectKey:  "my-project",
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// Nil option should fail.
	err = client.DismissMessage.ValidateCheckOpt(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}
}

func TestDismissMessage_ValidateDismissOpt(t *testing.T) {
	client, _ := NewClient("http://localhost/api/", "user", "pass")

	// Valid option should pass.
	err := client.DismissMessage.ValidateDismissOpt(&DismissMessageDismissOption{
		MessageType: "INFO",
		ProjectKey:  "my-project",
	})
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// Nil option should fail.
	err = client.DismissMessage.ValidateDismissOpt(nil)
	if err == nil {
		t.Error("expected error for nil option")
	}
}
