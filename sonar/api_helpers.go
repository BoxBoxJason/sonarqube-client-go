package sonar

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// callV2 builds a V2 API request, executes it and decodes the JSON response into a new T.
// queryOpt is encoded as query parameters (json tags) and body is sent as the JSON request body.
func callV2[T any](ctx context.Context, client *Client, method, path string, queryOpt, body any) (*T, *http.Response, error) {
	req, err := client.NewSonarQubeV2APIRequest(ctx, method, path, queryOpt, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(T)

	resp, err := client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// deleteV2 issues a DELETE V2 API request, expecting no response body.
func deleteV2(ctx context.Context, client *Client, path string) (*http.Response, error) {
	req, err := client.NewSonarQubeV2APIRequest(ctx, http.MethodDelete, path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// v2PathWithID returns prefix + "/" + the escaped ident, after checking that ident is not empty.
func v2PathWithID(prefix, ident, fieldName string) (string, error) {
	err := ValidateRequired(ident, fieldName)
	if err != nil {
		return "", err
	}

	return prefix + "/" + url.PathEscape(ident), nil
}

// callV1 builds a V1 API request, executes it and decodes the JSON response into a new T.
// opt is encoded as query/form parameters using its url struct tags.
func callV1[T any](ctx context.Context, client *Client, method, path string, opt any) (*T, *http.Response, error) {
	req, err := client.NewSonarQubeV1APIRequest(ctx, method, path, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	result := new(T)

	resp, err := client.Do(req, result)
	if err != nil {
		return nil, resp, err
	}

	return result, resp, nil
}

// postV1NoContent issues a POST V1 API request, expecting no response body.
func postV1NoContent(ctx context.Context, client *Client, path string, opt any) (*http.Response, error) {
	req, err := client.NewSonarQubeV1APIRequest(ctx, http.MethodPost, path, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
