package sonargo

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEmails_Send(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		if r.URL.Path != "/api/emails/send" {
			t.Errorf("expected path /api/emails/send, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client, err := NewClient(ts.URL+"/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &EmailsSendOption{
		Message: "Test message content",
		Subject: "Test Subject",
		To:      "test@example.com",
	}

	resp, err := client.Emails.Send(opt)
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", resp.StatusCode)
	}
}

func TestEmails_Send_ValidationError_NilOption(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	_, err = client.Emails.Send(nil)
	if err == nil {
		t.Fatal("expected error for nil option")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
}

func TestEmails_Send_ValidationError_MissingMessage(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &EmailsSendOption{
		To: "test@example.com",
	}

	_, err = client.Emails.Send(opt)
	if err == nil {
		t.Fatal("expected error for missing Message")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	if validationErr.Field != "Message" {
		t.Errorf("expected field 'Message', got '%s'", validationErr.Field)
	}
}

func TestEmails_Send_ValidationError_MissingTo(t *testing.T) {
	client, err := NewClient("http://localhost/api/", "user", "pass")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	opt := &EmailsSendOption{
		Message: "Test message",
	}

	_, err = client.Emails.Send(opt)
	if err == nil {
		t.Fatal("expected error for missing To")
	}

	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("expected ValidationError, got %T", err)
	}

	if validationErr.Field != "To" {
		t.Errorf("expected field 'To', got '%s'", validationErr.Field)
	}
}

func TestEmails_ValidateSendOpt(t *testing.T) {
	tests := []struct {
		name      string
		opt       *EmailsSendOption
		wantErr   bool
		wantField string
	}{
		{
			name:      "nil option",
			opt:       nil,
			wantErr:   true,
			wantField: "opt",
		},
		{
			name: "missing message",
			opt: &EmailsSendOption{
				To: "test@example.com",
			},
			wantErr:   true,
			wantField: "Message",
		},
		{
			name: "missing to",
			opt: &EmailsSendOption{
				Message: "Test message",
			},
			wantErr:   true,
			wantField: "To",
		},
		{
			name: "valid option",
			opt: &EmailsSendOption{
				Message: "Test message",
				To:      "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "valid with subject",
			opt: &EmailsSendOption{
				Message: "Test message",
				Subject: "Subject",
				To:      "test@example.com",
			},
			wantErr: false,
		},
	}

	client, _ := NewClient("http://localhost/api/", "user", "pass")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Emails.ValidateSendOpt(tt.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSendOpt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.wantField != "" {
				var validationErr *ValidationError
				if errors.As(err, &validationErr) {
					if validationErr.Field != tt.wantField {
						t.Errorf("expected field '%s', got '%s'", tt.wantField, validationErr.Field)
					}
				}
			}
		})
	}
}
