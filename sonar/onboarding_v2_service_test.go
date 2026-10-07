package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOnboardingService_GetOverview(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/onboarding/overview", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Onboarding.GetOverview(context.Background(), &OnboardingGetOverviewOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestOnboardingService_GetProjects(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/onboarding/projects", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Onboarding.GetProjects(context.Background(), &OnboardingGetProjectsOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}

func TestOnboardingService_GetProjects_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Onboarding.GetProjects(context.Background(), nil)
	require.Error(t, err)
}

func TestOnboardingService_GetStatistics(t *testing.T) {
	server := newTestServer(t, mockHandler(t, http.MethodGet, "/v2/onboarding/statistics", http.StatusOK, map[string]any{}))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Onboarding.GetStatistics(context.Background(), &OnboardingGetStatisticsOptions{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotNil(t, result)
}
