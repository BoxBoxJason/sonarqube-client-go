package sonar

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthentication_Login(t *testing.T) {
	handler := mockEmptyHandler(t, http.MethodPost, "/authentication/login", http.StatusNoContent)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	opt := &AuthenticationLoginOptions{
		Login:    "admin",
		Password: "secret",
	}

	resp, err := client.Authentication.Login(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAuthentication_Login_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name string
		opt  *AuthenticationLoginOptions
	}{
		{"nil option", nil},
		{"missing Login", &AuthenticationLoginOptions{Password: "secret"}},
		{"missing Password", &AuthenticationLoginOptions{Login: "admin"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Authentication.Login(tt.opt)
			assert.Error(t, err)
		})
	}
}

func TestAuthentication_Logout(t *testing.T) {
	handler := mockEmptyHandler(t, http.MethodPost, "/authentication/logout", http.StatusNoContent)
	server := newTestServer(t, handler)
	client := newTestClient(t, server.url())

	resp, err := client.Authentication.Logout()
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAuthentication_Validate(t *testing.T) {
	tests := []struct {
		name     string
		response string
		expected bool
	}{
		{"valid", `{"valid": true}`, true},
		{"invalid", `{"valid": false}`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := mockHandler(t, http.MethodGet, "/authentication/validate", http.StatusOK, tt.response)
			server := newTestServer(t, handler)
			client := newTestClient(t, server.url())

			result, resp, err := client.Authentication.Validate()
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			require.NotNil(t, result)
			assert.Equal(t, tt.expected, result.Valid)
		})
	}
}

func TestAuthentication_ValidateLoginOpt(t *testing.T) {
	client := newLocalhostClient(t)

	tests := []struct {
		name    string
		opt     *AuthenticationLoginOptions
		wantErr bool
	}{
		{"valid option", &AuthenticationLoginOptions{Login: "admin", Password: "secret"}, false},
		{"nil option", nil, true},
		{"missing Login", &AuthenticationLoginOptions{Password: "secret"}, true},
		{"missing Password", &AuthenticationLoginOptions{Login: "admin"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := client.Authentication.ValidateLoginOpt(tt.opt)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
