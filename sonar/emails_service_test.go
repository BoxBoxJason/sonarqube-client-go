package sonar

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmails_Send(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/emails/send", http.StatusNoContent))
	client := newTestClient(t, server.url())

	opt := &EmailsSendOption{
		Message: "Test message content",
		Subject: "Test Subject",
		To:      "test@example.com",
	}

	resp, err := client.Emails.Send(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestEmails_Send_ValidationErrors(t *testing.T) {
	tests := []struct {
		name      string
		opt       *EmailsSendOption
		wantField string
	}{
		{
			name:      "nil option",
			opt:       nil,
			wantField: "opt",
		},
		{
			name: "missing message",
			opt: &EmailsSendOption{
				To: "test@example.com",
			},
			wantField: "Message",
		},
		{
			name: "missing to",
			opt: &EmailsSendOption{
				Message: "Test message",
			},
			wantField: "To",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newLocalhostClient(t)

			_, err := client.Emails.Send(tt.opt)
			require.Error(t, err)

			var validationErr *ValidationError
			require.True(t, errors.As(err, &validationErr), "expected ValidationError, got %T", err)
			assert.Equal(t, tt.wantField, validationErr.Field)
		})
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

	client := newLocalhostClient(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Emails.ValidateSendOpt(tt.opt)

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantField != "" {
					var validationErr *ValidationError
					if errors.As(err, &validationErr) {
						assert.Equal(t, tt.wantField, validationErr.Field)
					}
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
