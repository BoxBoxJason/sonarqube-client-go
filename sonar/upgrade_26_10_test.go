package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectAnalysesService_Search_RelativePathFromScmRoot(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"paging":{"pageIndex":1,"pageSize":100,"total":1},"analyses":[{"key":"k","revision":"abc","relativePathFromScmRoot":"sub"}]}`))
		require.NoError(t, err)
	})
	client := newTestClient(t, server.URL)

	result, _, err := client.ProjectAnalyses.Search(context.Background(), &ProjectAnalysesSearchOptions{Project: "p"})
	require.NoError(t, err)
	require.Len(t, result.Analyses, 1)
	assert.Equal(t, "sub", result.Analyses[0].RelativePathFromScmRoot)
}

func TestComponentsService_Tree_RefQualifier(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"paging":{"pageIndex":1,"pageSize":100,"total":1},"baseComponent":{"key":"pf"},"components":[{"key":"a","qualifier":"SVW","refQualifier":"APP"}]}`))
		require.NoError(t, err)
	})
	client := newTestClient(t, server.URL)

	result, _, err := client.Components.Tree(context.Background(), &ComponentsTreeOptions{Component: "pf"})
	require.NoError(t, err)
	require.Len(t, result.Components, 1)
	assert.Equal(t, "APP", result.Components[0].RefQualifier)
}

func TestHistoryService_ProjectPageSizeValidation(t *testing.T) {
	client := newLocalhostClient(t)

	err := client.V2.History.ValidateProjectIssueCountsOpt(&HistoryProjectIssueCountsOptions{PortfolioID: "p", PageSize: 5001})
	require.Error(t, err)
	require.NoError(t, client.V2.History.ValidateProjectIssueCountsOpt(&HistoryProjectIssueCountsOptions{PortfolioID: "p", PageSize: 5000}))

	err = client.V2.History.ValidateProjectMeasuresOpt(&HistoryProjectMeasuresOptions{MetricKey: "m", PortfolioID: "p", PageSize: 5001})
	require.Error(t, err)
	require.NoError(t, client.V2.History.ValidateProjectMeasuresOpt(&HistoryProjectMeasuresOptions{MetricKey: "m", PortfolioID: "p", PageSize: 5000}))
}
