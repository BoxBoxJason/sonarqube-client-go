package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectBranches_Delete(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/project_branches/delete", http.StatusNoContent))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectBranchesDeleteOptions{
		Branch:  "feature-1",
		Project: "my-project",
	}

	resp, err := client.ProjectBranches.Delete(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestProjectBranches_Delete_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.ProjectBranches.Delete(context.Background(), nil)
	assert.Error(t, err)

	// Missing Branch should fail validation.
	_, err = client.ProjectBranches.Delete(context.Background(), &ProjectBranchesDeleteOptions{
		Project: "my-project",
	})
	assert.Error(t, err)

	// Missing Project should fail validation.
	_, err = client.ProjectBranches.Delete(context.Background(), &ProjectBranchesDeleteOptions{
		Branch: "feature-1",
	})
	assert.Error(t, err)
}

func TestProjectBranches_List(t *testing.T) {
	response := &ProjectBranchesList{
		Branches: []ProjectBranch{
			{
				Name:              "main",
				IsMain:            true,
				Type:              "LONG",
				Status:            ProjectBranchStatus{QualityGateStatus: "OK"},
				AnalysisDate:      "2024-01-01T00:00:00+0000",
				ExcludedFromPurge: true,
			},
			{
				Name:              "feature-1",
				IsMain:            false,
				Type:              "BRANCH",
				Status:            ProjectBranchStatus{QualityGateStatus: "ERROR"},
				ExcludedFromPurge: false,
			},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/project_branches/list", http.StatusOK, response))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectBranchesListOptions{
		Project: "my-project",
	}

	result, resp, err := client.ProjectBranches.List(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, result)
	assert.Len(t, result.Branches, 2)
	assert.Equal(t, "main", result.Branches[0].Name)
	assert.True(t, result.Branches[0].IsMain)
	assert.Equal(t, "OK", result.Branches[0].Status.QualityGateStatus)
	assert.False(t, result.Branches[1].ExcludedFromPurge)
}

func TestProjectBranches_List_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, _, err := client.ProjectBranches.List(context.Background(), nil)
	assert.Error(t, err)

	// Missing Project should fail validation.
	_, _, err = client.ProjectBranches.List(context.Background(), &ProjectBranchesListOptions{})
	assert.Error(t, err)
}

func TestProjectBranches_Rename(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/project_branches/rename", http.StatusNoContent))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectBranchesRenameOptions{
		Name:    "main",
		Project: "my-project",
	}

	resp, err := client.ProjectBranches.Rename(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestProjectBranches_Rename_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.ProjectBranches.Rename(context.Background(), nil)
	assert.Error(t, err)

	// Missing Name should fail validation.
	_, err = client.ProjectBranches.Rename(context.Background(), &ProjectBranchesRenameOptions{
		Project: "my-project",
	})
	assert.Error(t, err)

	// Missing Project should fail validation.
	_, err = client.ProjectBranches.Rename(context.Background(), &ProjectBranchesRenameOptions{
		Name: "main",
	})
	assert.Error(t, err)

	// Name exceeding max length should fail.
	longName := ""
	for i := 0; i < MaxBranchNameLength+1; i++ {
		longName += "a"
	}
	_, err = client.ProjectBranches.Rename(context.Background(), &ProjectBranchesRenameOptions{
		Name:    longName,
		Project: "my-project",
	})
	assert.Error(t, err)
}

func TestProjectBranches_SetAutomaticDeletionProtection(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/project_branches/set_automatic_deletion_protection", http.StatusNoContent))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectBranchesSetAutomaticDeletionProtectionOptions{
		Branch:  "feature-1",
		Project: "my-project",
		Value:   true,
	}

	resp, err := client.ProjectBranches.SetAutomaticDeletionProtection(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestProjectBranches_SetAutomaticDeletionProtection_False(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/project_branches/set_automatic_deletion_protection", http.StatusNoContent))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectBranchesSetAutomaticDeletionProtectionOptions{
		Branch:  "feature-1",
		Project: "my-project",
		Value:   false,
	}

	resp, err := client.ProjectBranches.SetAutomaticDeletionProtection(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestProjectBranches_SetAutomaticDeletionProtection_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.ProjectBranches.SetAutomaticDeletionProtection(context.Background(), nil)
	assert.Error(t, err)

	// Missing Branch should fail validation.
	_, err = client.ProjectBranches.SetAutomaticDeletionProtection(context.Background(), &ProjectBranchesSetAutomaticDeletionProtectionOptions{
		Project: "my-project",
		Value:   true,
	})
	assert.Error(t, err)

	// Missing Project should fail validation.
	_, err = client.ProjectBranches.SetAutomaticDeletionProtection(context.Background(), &ProjectBranchesSetAutomaticDeletionProtectionOptions{
		Branch: "feature-1",
		Value:  true,
	})
	assert.Error(t, err)
}

func TestProjectBranches_SetMain(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/project_branches/set_main", http.StatusNoContent))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectBranchesSetMainOptions{
		Branch:  "main",
		Project: "my-project",
	}

	resp, err := client.ProjectBranches.SetMain(context.Background(), opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestProjectBranches_SetMain_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.ProjectBranches.SetMain(context.Background(), nil)
	assert.Error(t, err)

	// Missing Branch should fail validation.
	_, err = client.ProjectBranches.SetMain(context.Background(), &ProjectBranchesSetMainOptions{
		Project: "my-project",
	})
	assert.Error(t, err)

	// Missing Project should fail validation.
	_, err = client.ProjectBranches.SetMain(context.Background(), &ProjectBranchesSetMainOptions{
		Branch: "main",
	})
	assert.Error(t, err)
}

func TestProjectBranches_ValidateDeleteOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Valid option should pass.
	err := client.ProjectBranches.ValidateDeleteOpt(&ProjectBranchesDeleteOptions{
		Branch:  "feature-1",
		Project: "my-project",
	})
	assert.NoError(t, err)
}

func TestProjectBranches_ValidateListOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Valid option should pass.
	err := client.ProjectBranches.ValidateListOpt(&ProjectBranchesListOptions{
		Project: "my-project",
	})
	assert.NoError(t, err)
}

func TestProjectBranches_ValidateRenameOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Valid option should pass.
	err := client.ProjectBranches.ValidateRenameOpt(&ProjectBranchesRenameOptions{
		Name:    "main",
		Project: "my-project",
	})
	assert.NoError(t, err)
}

func TestProjectBranches_ValidateSetAutomaticDeletionProtectionOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Valid option with true should pass.
	err := client.ProjectBranches.ValidateSetAutomaticDeletionProtectionOpt(&ProjectBranchesSetAutomaticDeletionProtectionOptions{
		Branch:  "feature-1",
		Project: "my-project",
		Value:   true,
	})
	assert.NoError(t, err)

	// Valid option with false should pass.
	err = client.ProjectBranches.ValidateSetAutomaticDeletionProtectionOpt(&ProjectBranchesSetAutomaticDeletionProtectionOptions{
		Branch:  "feature-1",
		Project: "my-project",
		Value:   false,
	})
	assert.NoError(t, err)
}

func TestProjectBranches_ValidateSetMainOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Valid option should pass.
	err := client.ProjectBranches.ValidateSetMainOpt(&ProjectBranchesSetMainOptions{
		Branch:  "main",
		Project: "my-project",
	})
	assert.NoError(t, err)
}
