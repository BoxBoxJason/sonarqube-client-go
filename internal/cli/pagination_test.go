package cli

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// paginatedOption is a test option struct with pagination support.
type paginatedOption struct {
	PaginationArgs
	Query string `url:"q,omitempty"`
}

// PaginationArgs mimics sonar.PaginationArgs (exported name for reflection).
type PaginationArgs struct {
	Page     int64 `url:"p,omitempty"`
	PageSize int64 `url:"ps,omitempty"`
}

// paginatedResponse is a test response struct with paging info.
type paginatedResponse struct {
	Items  []fakeResponse
	Paging testPaging
}

// testPaging mimics sonar.Paging.
type testPaging struct {
	PageIndex int64
	PageSize  int64
	Total     int64
}

// nonPaginatedOption is a test option struct without pagination.
type nonPaginatedOption struct {
	Query string `url:"q,omitempty"`
}

// TestHasPagination tests detection of PaginationArgs embedding.
func TestHasPagination(t *testing.T) {
	tests := []struct {
		name string
		typ  reflect.Type
		want bool
	}{
		{
			name: "with pagination",
			typ:  reflect.TypeOf(paginatedOption{}),
			want: true,
		},
		{
			name: "without pagination",
			typ:  reflect.TypeOf(nonPaginatedOption{}),
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := hasPagination(tc.typ)
			assert.Equal(t, tc.want, got)
		})
	}
}

// TestResponseHasPaging tests detection of Paging field on response types.
func TestResponseHasPaging(t *testing.T) {
	tests := []struct {
		name string
		typ  reflect.Type
		want bool
	}{
		{
			name: "struct with paging",
			typ:  reflect.TypeOf(&paginatedResponse{}),
			want: true,
		},
		{
			name: "struct without paging",
			typ:  reflect.TypeOf(&fakeResponse{}),
			want: false,
		},
		{
			name: "slice type (non-struct)",
			typ:  reflect.TypeOf([]fakeResponse{}),
			want: false,
		},
		{
			name: "pointer to slice",
			typ:  reflect.TypeOf(&[]fakeResponse{}),
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := responseHasPaging(tc.typ)
			assert.Equal(t, tc.want, got)
		})
	}
}

// TestFindSliceField tests finding the first slice field in a response type.
func TestFindSliceField(t *testing.T) {
	tests := []struct {
		name     string
		typ      reflect.Type
		wantName string
		wantOK   bool
	}{
		{
			name:     "response with items",
			typ:      reflect.TypeOf(&paginatedResponse{}),
			wantName: "Items",
			wantOK:   true,
		},
		{
			name:   "response without slice",
			typ:    reflect.TypeOf(&fakeResponse{}),
			wantOK: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			name, ok := findSliceField(tc.typ)
			assert.Equal(t, tc.wantOK, ok)

			if ok {
				assert.Equal(t, tc.wantName, name)
			}
		})
	}
}

// paginatedService mocks a service with a paginated method.
type paginatedService struct {
	callCount int
}

// Search simulates a paginated method that returns pages of results.
func (s *paginatedService) Search(opt *paginatedOption) (*paginatedResponse, *http.Response, error) {
	s.callCount++

	// Return different data based on page.
	switch opt.Page {
	case 1:
		return &paginatedResponse{
			Items:  []fakeResponse{{Name: "a"}, {Name: "b"}},
			Paging: testPaging{PageIndex: 1, PageSize: 2, Total: 3},
		}, nil, nil
	case 2:
		return &paginatedResponse{
			Items:  []fakeResponse{{Name: "c"}},
			Paging: testPaging{PageIndex: 2, PageSize: 2, Total: 3},
		}, nil, nil
	default:
		return &paginatedResponse{
			Items:  nil,
			Paging: testPaging{PageIndex: opt.Page, PageSize: 2, Total: 3},
		}, nil, nil
	}
}

// TestPaginateAll tests multi-page result collection.
func TestPaginateAll(t *testing.T) {
	svc := &paginatedService{}
	opt := &paginatedOption{Query: "test"}
	optValue := reflect.New(reflect.TypeOf(*opt))
	optValue.Elem().Set(reflect.ValueOf(*opt))
	svcValue := reflect.ValueOf(svc)
	responseType := reflect.TypeOf(&paginatedResponse{})

	result, err := PaginateAll(svcValue, "Search", optValue, PatternResponseBody, responseType)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	resp, ok := result.(*paginatedResponse)
	assert.True(t, ok)
	assert.Len(t, resp.Items, 3)
	assert.Equal(t, "a", resp.Items[0].Name)
	assert.Equal(t, "b", resp.Items[1].Name)
	assert.Equal(t, "c", resp.Items[2].Name)

	// Should have been called twice (page 1 and page 2).
	assert.Equal(t, 2, svc.callCount)
}

// TestSetPageField tests setting pagination fields on option structs.
func TestSetPageField(t *testing.T) {
	opt := &paginatedOption{}
	optVal := reflect.ValueOf(opt).Elem()

	setPageField(optVal, "Page", 5)
	assert.Equal(t, int64(5), opt.Page)

	setPageField(optVal, "PageSize", 100)
	assert.Equal(t, int64(100), opt.PageSize)
}
