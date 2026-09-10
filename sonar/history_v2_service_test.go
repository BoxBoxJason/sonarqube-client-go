package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -----------------------------------------------------------------------------
// GetIssueCountHistory
// -----------------------------------------------------------------------------

func TestHistoryService_GetIssueCountHistory(t *testing.T) {
	response := HistoryIssueCountHistoryResponse{
		IssueCountHistory: []HistoryIssueCount{
			{
				Date: "2026-01-01T12:00:00Z",
				Distribution: []HistoryIssueCountDistribution{
					{Key: "BUG", Value: 3},
					{Key: "CODE_SMELL", Value: 7},
				},
			},
		},
	}
	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/history/issue-count-history", http.StatusOK,
		map[string]string{
			"entityId":   "00000000-0000-0000-0000-000000000000",
			"entityType": "PROJECT_BRANCH",
			"startDate":  "2026-01-01T00:00:00Z",
			"sliceBy":    "issueType",
		}, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.History.GetIssueCountHistory(context.Background(), &HistoryIssueCountHistoryOptions{
		EntityID:   "00000000-0000-0000-0000-000000000000",
		EntityType: "PROJECT_BRANCH",
		StartDate:  "2026-01-01T00:00:00Z",
		SliceBy:    "issueType",
		IssueTypes: []string{"BUG", "CODE_SMELL"},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, result.IssueCountHistory, 1)
	assert.Equal(t, "2026-01-01T12:00:00Z", result.IssueCountHistory[0].Date)
	require.Len(t, result.IssueCountHistory[0].Distribution, 2)
	assert.Equal(t, int32(3), result.IssueCountHistory[0].Distribution[0].Value)
}

func TestHistoryService_GetIssueCountHistory_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &HistoryService{client: client}

	tests := []struct {
		opt  *HistoryIssueCountHistoryOptions
		name string
	}{
		{nil, "nil opt"},
		{&HistoryIssueCountHistoryOptions{EntityType: "PROJECT_BRANCH", StartDate: "2026-01-01T00:00:00Z"}, "missing EntityID"},
		{&HistoryIssueCountHistoryOptions{EntityID: "id", StartDate: "2026-01-01T00:00:00Z"}, "missing EntityType"},
		{&HistoryIssueCountHistoryOptions{EntityID: "id", EntityType: "PROJECT_BRANCH"}, "missing StartDate"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, resp, err := svc.GetIssueCountHistory(context.Background(), tt.opt)
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Nil(t, resp)
		})
	}
}

// -----------------------------------------------------------------------------
// GetMeasuresHistory
// -----------------------------------------------------------------------------

func TestHistoryService_GetMeasuresHistory(t *testing.T) {
	response := HistoryMeasuresHistoryResponse{
		MeasuresHistory: []HistoryMeasureItem{
			{
				Date: "2026-01-01T12:00:00Z",
				Measures: []HistoryMeasureEntry{
					{Metric: "coverage", Type: "PERCENT", Value: "81.2"},
				},
			},
		},
	}
	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/history/measures-history", http.StatusOK,
		map[string]string{
			"entityType": HistoryEntityTypeProjectBranch,
			"entityId":   "00000000-0000-0000-0000-000000000000",
			"startDate":  "2026-01-01T00:00:00Z",
		}, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.History.GetMeasuresHistory(context.Background(), &HistoryMeasuresHistoryOptions{
		EntityType: HistoryEntityTypeProjectBranch,
		EntityID:   "00000000-0000-0000-0000-000000000000",
		MetricKeys: []string{"coverage", "bugs"},
		StartDate:  "2026-01-01T00:00:00Z",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, result.MeasuresHistory, 1)
	require.Len(t, result.MeasuresHistory[0].Measures, 1)
	assert.Equal(t, "coverage", result.MeasuresHistory[0].Measures[0].Metric)
	assert.Equal(t, "81.2", result.MeasuresHistory[0].Measures[0].Value)
}

func TestHistoryService_GetMeasuresHistory_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &HistoryService{client: client}

	tests := []struct {
		opt  *HistoryMeasuresHistoryOptions
		name string
	}{
		{nil, "nil opt"},
		{&HistoryMeasuresHistoryOptions{EntityID: "id", MetricKeys: []string{"coverage"}, StartDate: "d"}, "missing EntityType"},
		{&HistoryMeasuresHistoryOptions{EntityType: "BRANCH", EntityID: "id", MetricKeys: []string{"coverage"}, StartDate: "d"}, "invalid EntityType"},
		{&HistoryMeasuresHistoryOptions{EntityType: HistoryEntityTypePortfolio, MetricKeys: []string{"coverage"}, StartDate: "d"}, "missing EntityID"},
		{&HistoryMeasuresHistoryOptions{EntityType: HistoryEntityTypePortfolio, EntityID: "id", StartDate: "d"}, "missing MetricKeys"},
		{&HistoryMeasuresHistoryOptions{EntityType: HistoryEntityTypePortfolio, EntityID: "id", MetricKeys: []string{"coverage"}}, "missing StartDate"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, resp, err := svc.GetMeasuresHistory(context.Background(), tt.opt)
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Nil(t, resp)
		})
	}
}

// -----------------------------------------------------------------------------
// GetProjectIssueCounts
// -----------------------------------------------------------------------------

func TestHistoryService_GetProjectIssueCounts(t *testing.T) {
	refCount := int64(4)
	response := HistoryProjectIssueCountsResponse{
		HiddenProjectCount: 2,
		ProjectIssueCounts: []HistoryProjectIssueCount{
			{
				BranchID:            "branch-1",
				ProjectName:         "My Project",
				ProjectKey:          "my-project",
				BranchName:          "main",
				IssueCount:          9,
				ReferenceIssueCount: &refCount,
			},
		},
		Page: PageResponseV2{PageIndex: 1, PageSize: 50, Total: 1},
	}
	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/history/project-issue-counts", http.StatusOK,
		map[string]string{
			"portfolioId":   "portfolio-1",
			"requireIssues": "true",
		}, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.History.GetProjectIssueCounts(context.Background(), &HistoryProjectIssueCountsOptions{
		PortfolioID:   "portfolio-1",
		RequireIssues: true,
		Sort:          []string{"-issueCount"},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, int32(2), result.HiddenProjectCount)
	require.Len(t, result.ProjectIssueCounts, 1)
	assert.Equal(t, "my-project", result.ProjectIssueCounts[0].ProjectKey)
	assert.Equal(t, int64(9), result.ProjectIssueCounts[0].IssueCount)
	require.NotNil(t, result.ProjectIssueCounts[0].ReferenceIssueCount)
	assert.Equal(t, int64(4), *result.ProjectIssueCounts[0].ReferenceIssueCount)
	assert.Equal(t, int32(1), result.Page.Total)
}

func TestHistoryService_GetProjectIssueCounts_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &HistoryService{client: client}

	tests := []struct {
		opt  *HistoryProjectIssueCountsOptions
		name string
	}{
		{nil, "nil opt"},
		{&HistoryProjectIssueCountsOptions{}, "no selector"},
		{&HistoryProjectIssueCountsOptions{PortfolioID: "p", EntityType: "PORTFOLIO", EntityID: "e"}, "both selector forms"},
		{&HistoryProjectIssueCountsOptions{EntityType: "PORTFOLIO"}, "entityType without entityId"},
		{&HistoryProjectIssueCountsOptions{EntityID: "e"}, "entityId without entityType"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, resp, err := svc.GetProjectIssueCounts(context.Background(), tt.opt)
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Nil(t, resp)
		})
	}
}

// -----------------------------------------------------------------------------
// GetProjectMeasures
// -----------------------------------------------------------------------------

func TestHistoryService_GetProjectMeasures(t *testing.T) {
	response := HistoryProjectMeasuresResponse{
		HiddenProjectCount: 0,
		ProjectMeasures: []HistoryProjectMeasure{
			{
				BranchID:    "branch-1",
				BranchName:  "main",
				ProjectKey:  "my-project",
				ProjectName: "My Project",
				Measure: HistoryProjectMeasureMetric{
					CurrentValue:   "81.2",
					Metric:         "coverage",
					ReferenceValue: "79.0",
					Type:           "PERCENT",
				},
			},
		},
		Page: PageResponseV2{PageIndex: 1, PageSize: 50, Total: 1},
	}
	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/history/project-measures", http.StatusOK,
		map[string]string{
			"metricKey":  "coverage",
			"entityType": "APPLICATION",
			"entityId":   "app-branch-1",
		}, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.History.GetProjectMeasures(context.Background(), &HistoryProjectMeasuresOptions{
		MetricKey:  "coverage",
		EntityType: "APPLICATION",
		EntityID:   "app-branch-1",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, result.ProjectMeasures, 1)
	assert.Equal(t, "my-project", result.ProjectMeasures[0].ProjectKey)
	assert.Equal(t, "81.2", result.ProjectMeasures[0].Measure.CurrentValue)
	assert.Equal(t, "PERCENT", result.ProjectMeasures[0].Measure.Type)
}

func TestHistoryService_GetProjectMeasures_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)
	svc := &HistoryService{client: client}

	tests := []struct {
		opt  *HistoryProjectMeasuresOptions
		name string
	}{
		{nil, "nil opt"},
		{&HistoryProjectMeasuresOptions{EntityType: "PORTFOLIO", EntityID: "e"}, "missing MetricKey"},
		{&HistoryProjectMeasuresOptions{MetricKey: "coverage"}, "no selector"},
		{&HistoryProjectMeasuresOptions{MetricKey: "coverage", PortfolioID: "p", EntityID: "e"}, "both selector forms"},
		{&HistoryProjectMeasuresOptions{MetricKey: "coverage", EntityType: "PORTFOLIO"}, "entityType without entityId"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, resp, err := svc.GetProjectMeasures(context.Background(), tt.opt)
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Nil(t, resp)
		})
	}
}
