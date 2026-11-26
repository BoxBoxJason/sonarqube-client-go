package sonar

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApplicationsService_SetTags_Empty(t *testing.T) {
	server := newTestServer(t, wireHandler(t, "/applications/set_tags", "application=app&tags="))
	client := newTestClient(t, server.URL)

	_, err := client.Applications.SetTags(context.Background(), &ApplicationsSetTagsOptions{
		Application: "app",
		Tags:        []string{},
	})
	require.NoError(t, err)
}

func TestGovernanceReportsService_UpdateRecipients_Empty(t *testing.T) {
	server := newTestServer(t, wireHandler(t, "/governance_reports/update_recipients", "componentKey=portfolio&recipients="))
	client := newTestClient(t, server.URL)

	_, err := client.GovernanceReports.UpdateRecipients(context.Background(), &GovernanceReportsUpdateRecipientsOptions{
		ComponentKey: "portfolio",
	})
	require.NoError(t, err)
}
