package sonargo

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectBranches_Delete(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/project_branches/delete", http.StatusNoContent))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectBranchesDeleteOption{
		Branch:  "feature-1",
		Project: "my-project",
	}

	resp, err := client.ProjectBranches.Delete(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestProjectBranches_Delete_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.ProjectBranches.Delete(nil)
	assert.Error(t, err)

	// Missing Branch should fail validation.
	_, err = client.ProjectBranches.Delete(&ProjectBranchesDeleteOption{
		Project: "my-project",
	})
	assert.Error(t, err)

	// Missing Project should fail validation.
	_, err = client.ProjectBranches.Delete(&ProjectBranchesDeleteOption{
		Branch: "feature-1",
	})
	assert.Error(t, err)
}

func TestProjectBranches_List(t *testing.T) {
	response := &ProjectBranchesList{
		Branches: []Branch{
			{
				Name:              "main",
				IsMain:            true,
				Type:              "LONG",
				Status:            BranchStatus{QualityGateStatus: "OK"},
				AnalysisDate:      "2024-01-01T00:00:00+0000",
				ExcludedFromPurge: true,
			},
			{
				Name:              "feature-1",
				IsMain:            false,
				Type:              "BRANCH",
				Status:            BranchStatus{QualityGateStatus: "ERROR"},
				ExcludedFromPurge: false,
			},
		},
	}

	server := newTestServer(t, mockHandler(t, http.MethodGet, "/project_branches/list", http.StatusOK, response))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectBranchesListOption{
		Project: "my-project",
	}

	result, resp, err := client.ProjectBranches.List(opt)
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
	_, _, err := client.ProjectBranches.List(nil)
	assert.Error(t, err)

	// Missing Project should fail validation.
	_, _, err = client.ProjectBranches.List(&ProjectBranchesListOption{})
	assert.Error(t, err)
}

func TestProjectBranches_Rename(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/project_branches/rename", http.StatusNoContent))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectBranchesRenameOption{
		Name:    "main",
		Project: "my-project",
	}

	resp, err := client.ProjectBranches.Rename(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestProjectBranches_Rename_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.ProjectBranches.Rename(nil)
	assert.Error(t, err)

	// Missing Name should fail validation.
	_, err = client.ProjectBranches.Rename(&ProjectBranchesRenameOption{
		Project: "my-project",
	})
	assert.Error(t, err)

	// Missing Project should fail validation.
	_, err = client.ProjectBranches.Rename(&ProjectBranchesRenameOption{
		Name: "main",
	})
	assert.Error(t, err)

	// Name exceeding max length should fail.
	longName := ""
	for i := 0; i < MaxBranchNameLength+1; i++ {
		longName += "a"
	}
	_, err = client.ProjectBranches.Rename(&ProjectBranchesRenameOption{
		Name:    longName,
		Project: "my-project",
	})
	assert.Error(t, err)
}

func TestProjectBranches_SetAutomaticDeletionProtection(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/project_branches/set_automatic_deletion_protection", http.StatusNoContent))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectBranchesSetAutomaticDeletionProtectionOption{
		Branch:  "feature-1",
		Project: "my-project",
		Value:   true,
	}

	resp, err := client.ProjectBranches.SetAutomaticDeletionProtection(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestProjectBranches_SetAutomaticDeletionProtection_False(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/project_branches/set_automatic_deletion_protection", http.StatusNoContent))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectBranchesSetAutomaticDeletionProtectionOption{
		Branch:  "feature-1",
		Project: "my-project",
		Value:   false,
	}

	resp, err := client.ProjectBranches.SetAutomaticDeletionProtection(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestProjectBranches_SetAutomaticDeletionProtection_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.ProjectBranches.SetAutomaticDeletionProtection(nil)
	assert.Error(t, err)

	// Missing Branch should fail validation.
	_, err = client.ProjectBranches.SetAutomaticDeletionProtection(&ProjectBranchesSetAutomaticDeletionProtectionOption{
		Project: "my-project",
		Value:   true,
	})
	assert.Error(t, err)

	// Missing Project should fail validation.
	_, err = client.ProjectBranches.SetAutomaticDeletionProtection(&ProjectBranchesSetAutomaticDeletionProtectionOption{
		Branch: "feature-1",
		Value:  true,
	})
	assert.Error(t, err)
}

func TestProjectBranches_SetMain(t *testing.T) {
	server := newTestServer(t, mockEmptyHandler(t, http.MethodPost, "/project_branches/set_main", http.StatusNoContent))
	defer server.Close()

	client := newTestClient(t, server.URL)

	opt := &ProjectBranchesSetMainOption{
		Branch:  "main",
		Project: "my-project",
	}

	resp, err := client.ProjectBranches.SetMain(opt)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestProjectBranches_SetMain_ValidationError(t *testing.T) {
	client := newLocalhostClient(t)

	// Nil option should fail validation.
	_, err := client.ProjectBranches.SetMain(nil)
	assert.Error(t, err)

	// Missing Branch should fail validation.
	_, err = client.ProjectBranches.SetMain(&ProjectBranchesSetMainOption{
		Project: "my-project",
	})
	assert.Error(t, err)

	// Missing Project should fail validation.
	_, err = client.ProjectBranches.SetMain(&ProjectBranchesSetMainOption{
		Branch: "main",
	})
	assert.Error(t, err)
}

func TestProjectBranches_ValidateDeleteOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Valid option should pass.
	err := client.ProjectBranches.ValidateDeleteOpt(&ProjectBranchesDeleteOption{
		Branch:  "feature-1",
		Project: "my-project",
	})
	assert.NoError(t, err)
}

func TestProjectBranches_ValidateListOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Valid option should pass.
	err := client.ProjectBranches.ValidateListOpt(&ProjectBranchesListOption{
		Project: "my-project",
	})
	assert.NoError(t, err)
}

func TestProjectBranches_ValidateRenameOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Valid option should pass.
	err := client.ProjectBranches.ValidateRenameOpt(&ProjectBranchesRenameOption{
		Name:    "main",
		Project: "my-project",
	})
	assert.NoError(t, err)
}

func TestProjectBranches_ValidateSetAutomaticDeletionProtectionOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Valid option with true should pass.
	err := client.ProjectBranches.ValidateSetAutomaticDeletionProtectionOpt(&ProjectBranchesSetAutomaticDeletionProtectionOption{
		Branch:  "feature-1",
		Project: "my-project",
		Value:   true,
	})
	assert.NoError(t, err)

	// Valid option with false should pass.
	err = client.ProjectBranches.ValidateSetAutomaticDeletionProtectionOpt(&ProjectBranchesSetAutomaticDeletionProtectionOption{
		Branch:  "feature-1",
		Project: "my-project",
		Value:   false,
	})
	assert.NoError(t, err)
}

func TestProjectBranches_ValidateSetMainOpt(t *testing.T) {
	client := newLocalhostClient(t)

	// Valid option should pass.
	err := client.ProjectBranches.ValidateSetMainOpt(&ProjectBranchesSetMainOption{
		Branch:  "main",
		Project: "my-project",
	})
	assert.NoError(t, err)
}
