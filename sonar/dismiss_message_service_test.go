package sonargo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDismissMessage_Check(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/dismiss_message/check", http.StatusOK, &DismissMessageCheck{
		Dismissed: true,
	}))
	client := newTestClient(t, server.url())

	opt := &DismissMessageCheckOption{
		MessageType: "INFO",
		ProjectKey:  "my-project",
	}

	result, resp, err := client.DismissMessage.Check(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.True(t, result.Dismissed)
}

func TestDismissMessage_Check_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *DismissMessageCheckOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing MessageType",
			opt: &DismissMessageCheckOption{
				ProjectKey: "my-project",
			},
		},
		{
			name: "missing ProjectKey",
			opt: &DismissMessageCheckOption{
				MessageType: "INFO",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := client.DismissMessage.Check(tt.opt)
			require.Error(t, err)
		})
	}
}

func TestDismissMessage_Dismiss(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/dismiss_message/dismiss", http.StatusNoContent))
	client := newTestClient(t, server.url())

	opt := &DismissMessageDismissOption{
		MessageType: "WARNING",
		ProjectKey:  "my-project",
	}

	resp, err := client.DismissMessage.Dismiss(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestDismissMessage_Dismiss_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *DismissMessageDismissOption
	}{
		{
			name: "nil option",
			opt:  nil,
		},
		{
			name: "missing MessageType",
			opt: &DismissMessageDismissOption{
				ProjectKey: "my-project",
			},
		},
		{
			name: "missing ProjectKey",
			opt: &DismissMessageDismissOption{
				MessageType: "INFO",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.DismissMessage.Dismiss(tt.opt)
			require.Error(t, err)
		})
	}
}

func TestDismissMessage_ValidateCheckOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *DismissMessageCheckOption
		wantErr bool
	}{
		{
			name: "valid option",
			opt: &DismissMessageCheckOption{
				MessageType: "INFO",
				ProjectKey:  "my-project",
			},
			wantErr: false,
		},
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.DismissMessage.ValidateCheckOpt(tt.opt)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDismissMessage_ValidateDismissOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *DismissMessageDismissOption
		wantErr bool
	}{
		{
			name: "valid option",
			opt: &DismissMessageDismissOption{
				MessageType: "INFO",
				ProjectKey:  "my-project",
			},
			wantErr: false,
		},
		{
			name:    "nil option",
			opt:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.DismissMessage.ValidateDismissOpt(tt.opt)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
