// Package response provides the client and response types for the SonarQube API.
package response

import (
	"fmt"
	"net/http"
	"net/url"
)

// Client communicates with the SonarQube API.
type Client struct {
	baseURL     *url.URL
	httpClient  *http.Client
	Webservices *WebservicesService
	username    string
	password    string
	token       string
	authType    authType
}

// NewClient returns a new SonarQube API client.
func NewClient(endpoint, username, password string) (*Client, error) {
	client := &Client{
		username:    username,
		password:    password,
		authType:    basicAuth,
		httpClient:  http.DefaultClient,
		baseURL:     nil,
		token:       "",
		Webservices: nil,
	}

	if endpoint == "" {
		err := client.SetBaseURL(defaultBaseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to set default base URL: %w", err)
		}
	} else {
		err := client.SetBaseURL(endpoint)
		if err != nil {
			return nil, err
		}
	}

	client.Webservices = &WebservicesService{client: client}

	return client, nil
}
