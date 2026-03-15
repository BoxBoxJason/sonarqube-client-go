package sonar

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCodePeriods_List(t *testing.T) {
	response := `{
		"newCodePeriods": [
			{
				"projectKey": "my-project",
				"branchKey": "main",
				"type": "PREVIOUS_VERSION",
				"inherited": false
			},
			{
				"projectKey": "my-project",
				"branchKey": "feature-1",
				"type": "NUMBER_OF_DAYS",
				"value": "30",
				"inherited": true
			}
		]
	}`

	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/new_code_periods/list", r.URL.Path)
		assert.Equal(t, "my-project", r.URL.Query().Get("project"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	})

	client := newTestClient(t, server.URL)

	opt := &NewCodePeriodsListOptions{
		Project: "my-project",
	}

	result, resp, err := client.NewCodePeriods.List(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.NewCodePeriods, 2)
	assert.Equal(t, NewCodePeriodTypePreviousVersion, result.NewCodePeriods[0].Type)
	assert.Equal(t, "30", result.NewCodePeriods[1].Value)
}

func TestNewCodePeriods_List_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.NewCodePeriods.List(nil)
	assert.Error(t, err)

	// Missing Project should fail validation.
	_, _, err = client.NewCodePeriods.List(&NewCodePeriodsListOptions{})
	assert.Error(t, err)
}

func TestNewCodePeriods_Set(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/new_code_periods/set", r.URL.Path)
		assert.Equal(t, NewCodePeriodTypeNumberOfDays, r.URL.Query().Get("type"))
		assert.Equal(t, "30", r.URL.Query().Get("value"))

		w.WriteHeader(http.StatusOK)
	})

	client := newTestClient(t, server.URL)

	opt := &NewCodePeriodsSetOptions{
		Type:  NewCodePeriodTypeNumberOfDays,
		Value: "30",
	}

	resp, err := client.NewCodePeriods.Set(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestNewCodePeriods_Set_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.NewCodePeriods.Set(nil)
	assert.Error(t, err)

	// Missing Type should fail validation.
	_, err = client.NewCodePeriods.Set(&NewCodePeriodsSetOptions{})
	assert.Error(t, err)

	// Invalid Type should fail validation.
	_, err = client.NewCodePeriods.Set(&NewCodePeriodsSetOptions{
		Type: "INVALID_TYPE",
	})
	assert.Error(t, err)

	// SPECIFIC_ANALYSIS without Branch should fail validation.
	_, err = client.NewCodePeriods.Set(&NewCodePeriodsSetOptions{
		Type: NewCodePeriodTypeSpecificAnalysis,
	})
	assert.Error(t, err)

	// REFERENCE_BRANCH without Project should fail validation.
	_, err = client.NewCodePeriods.Set(&NewCodePeriodsSetOptions{
		Type: NewCodePeriodTypeReferenceBranch,
	})
	assert.Error(t, err)
}

func TestNewCodePeriods_Show(t *testing.T) {
	response := `{
		"type": "NUMBER_OF_DAYS",
		"value": "30",
		"inherited": false
	}`

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/new_code_periods/show", http.StatusOK, response))
	client := newTestClient(t, server.URL)

	result, resp, err := client.NewCodePeriods.Show(nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Equal(t, NewCodePeriodTypeNumberOfDays, result.Type)
	assert.Equal(t, "30", result.Value)
}

func TestNewCodePeriods_Show_WithOptions(t *testing.T) {
	response := `{
		"projectKey": "my-project",
		"branchKey": "main",
		"type": "REFERENCE_BRANCH",
		"value": "main",
		"inherited": true
	}`

	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/new_code_periods/show", r.URL.Path)
		assert.Equal(t, "my-project", r.URL.Query().Get("project"))
		assert.Equal(t, "main", r.URL.Query().Get("branch"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	})

	client := newTestClient(t, server.URL)

	opt := &NewCodePeriodsShowOptions{
		Project: "my-project",
		Branch:  "main",
	}

	result, _, err := client.NewCodePeriods.Show(opt)
	require.NoError(t, err)
	assert.Equal(t, "my-project", result.ProjectKey)
	assert.Equal(t, NewCodePeriodTypeReferenceBranch, result.Type)
	assert.Equal(t, "main", result.Value)
}

func TestNewCodePeriods_Unset(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/new_code_periods/unset", http.StatusOK))
	client := newTestClient(t, server.URL)

	resp, err := client.NewCodePeriods.Unset(nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestNewCodePeriods_Unset_WithOptions(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/new_code_periods/unset", r.URL.Path)
		assert.Equal(t, "my-project", r.URL.Query().Get("project"))

		w.WriteHeader(http.StatusOK)
	})

	client := newTestClient(t, server.URL)

	opt := &NewCodePeriodsUnsetOptions{
		Project: "my-project",
	}

	resp, err := client.NewCodePeriods.Unset(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestNewCodePeriods_ValidateSetOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// All valid types without special requirements should pass.
	validTypes := []string{NewCodePeriodTypePreviousVersion}
	for _, periodType := range validTypes {
		err := client.NewCodePeriods.ValidateSetOpt(&NewCodePeriodsSetOptions{
			Type: periodType,
		})
		assert.NoError(t, err, "expected nil error for type '%s'", periodType)
	}

	// NUMBER_OF_DAYS with valid Value should pass.
	err := client.NewCodePeriods.ValidateSetOpt(&NewCodePeriodsSetOptions{
		Type:  NewCodePeriodTypeNumberOfDays,
		Value: "30",
	})
	assert.NoError(t, err)

	// SPECIFIC_ANALYSIS with Branch should pass.
	err = client.NewCodePeriods.ValidateSetOpt(&NewCodePeriodsSetOptions{
		Type:   NewCodePeriodTypeSpecificAnalysis,
		Branch: "main",
	})
	assert.NoError(t, err)

	// REFERENCE_BRANCH with Project should pass.
	err = client.NewCodePeriods.ValidateSetOpt(&NewCodePeriodsSetOptions{
		Type:    NewCodePeriodTypeReferenceBranch,
		Project: "my-project",
	})
	assert.NoError(t, err)
}
