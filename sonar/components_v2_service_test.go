package sonar

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComponentsV2_GetProjects(t *testing.T) {
	response := ComponentsProjectsV2{
		Projects: []ComponentsProjectV2{{
			Id: "295f84b3-c0dd-427d-9c88-11afda4a6668", Key: "probe-proj", Name: "Probe", Visibility: "public",
			Tags:  []string{"a", "b"},
			Links: []ComponentsProjectLinkV2{{Name: "home", Type: "home", Url: "http://x.io"}},
		}},
		Page: PageResponseV2{PageIndex: 1, PageSize: 50, Total: 1},
	}
	server := newTestServer(t, mockHandlerWithParams(t, http.MethodGet, "/v2/components/projects", http.StatusOK,
		map[string]string{"keys": "probe-proj", "pageSize": "10"}, response))
	client := newTestClient(t, server.url())

	result, resp, err := client.V2.Components.GetProjects(context.Background(), &ComponentsGetProjectsOptionsV2{
		PaginationParamsV2: PaginationParamsV2{PageSize: 10},
		Keys:               []string{"probe-proj"},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.Len(t, result.Projects, 1)
	assert.Equal(t, "probe-proj", result.Projects[0].Key)
	assert.Equal(t, []string{"a", "b"}, result.Projects[0].Tags)
	assert.Equal(t, "home", result.Projects[0].Links[0].Type)
	assert.Equal(t, int32(1), result.Page.Total)
}

func TestComponentsV2_GetProjects_Validation(t *testing.T) {
	client := newLocalhostClient(t)

	_, _, err := client.V2.Components.GetProjects(context.Background(), nil)
	require.Error(t, err)

	_, _, err = client.V2.Components.GetProjects(context.Background(), &ComponentsGetProjectsOptionsV2{Tags: []string{"a"}})
	require.Error(t, err)

	_, _, err = client.V2.Components.GetProjects(context.Background(), &ComponentsGetProjectsOptionsV2{
		PaginationParamsV2: PaginationParamsV2{PageSize: 5001},
		Ids:                []string{"id"},
	})
	require.Error(t, err)
}
