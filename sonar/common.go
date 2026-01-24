package sonargo

const (
	// MaxPageSize is the maximum allowed page size for pagination.
	MaxPageSize = 500
	// MinPageSize is the minimum allowed page size for pagination.
	MinPageSize = 1
)

// PaginationArgs contains common pagination parameters for API requests.
type PaginationArgs struct {
	// Page is the response page number. Must be strictly greater than 0.
	Page int64 `url:"p,omitempty"`
	// PageSize is the response page size. Must be greater than 0 and less than or equal to 500.
	PageSize int64 `url:"ps,omitempty"`
}

// Validate validates the pagination arguments.
func (p *PaginationArgs) Validate() error {
	return ValidatePagination(p.Page, p.PageSize)
}
