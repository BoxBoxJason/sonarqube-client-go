package sonar

import (
	"context"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Response payloads follow those of SonarQube Enterprise 2026.5.4.

func TestGovernanceReportsStatus_ComponentFrequency(t *testing.T) {
	var status GovernanceReportsStatus

	require.NoError(t, json.Unmarshal([]byte(`{"componentFrequency":"weekly","globalFrequency":"monthly"}`), &status))
	assert.Equal(t, "weekly", status.ComponentFrequency)
}

func TestViewDetails_NestedApplication(t *testing.T) {
	var details ViewDetails

	require.NoError(t, json.Unmarshal([]byte(`{"key":"pf1","subViews":[{"key":"pf1:app1","qualifier":"APP","selectedBranches":["main"],"originalKey":"app1"}]}`), &details))
	require.Len(t, details.SubViews, 1)
	assert.Equal(t, []string{"main"}, details.SubViews[0].SelectedBranches)
}

func TestIssuesSandboxSettings_DeferralDays(t *testing.T) {
	var settings IssuesSandboxSettings

	require.NoError(t, json.Unmarshal([]byte(`{"deferralDays":3}`), &settings))
	assert.Equal(t, int64(3), settings.DeferralDays)
}

func TestFixSuggestionsSetEnablementOptions_Provider(t *testing.T) {
	opt := &FixSuggestionsSetEnablementOptions{
		Enablement: "ENABLED_FOR_ALL_PROJECTS",
		Provider:   &FixSuggestionsProviderOptions{Type: "OPENAI", Model: "model"},
	}

	data, err := json.Marshal(opt)
	require.NoError(t, err)
	assert.JSONEq(t, `{"enablement":"ENABLED_FOR_ALL_PROJECTS","provider":{"type":"OPENAI","model":"model"}}`, string(data))

	client := newLocalhostClient(t)
	_, err = client.V2.FixSuggestions.SetEnablement(context.Background(), &FixSuggestionsSetEnablementOptions{
		Enablement: "ENABLED_FOR_ALL_PROJECTS",
		Provider:   &FixSuggestionsProviderOptions{},
	})
	assert.Error(t, err)
}

func TestFixSuggestionsSetEnablementOptions_ProjectKeys(t *testing.T) {
	tests := map[string]struct {
		keys []string
		body string
	}{
		"Nil":   {body: `{"enablement":"DISABLED"}`},
		"Empty": {keys: []string{}, body: `{"enablement":"DISABLED","enabledProjectKeys":[]}`},
		"Set":   {keys: []string{"a"}, body: `{"enablement":"DISABLED","enabledProjectKeys":["a"]}`},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			reader, err := marshalBody(&FixSuggestionsSetEnablementOptions{
				Enablement:         "DISABLED",
				EnabledProjectKeys: test.keys,
			})
			require.NoError(t, err)

			data, err := io.ReadAll(reader)
			require.NoError(t, err)
			assert.JSONEq(t, test.body, string(data))
		})
	}
}
