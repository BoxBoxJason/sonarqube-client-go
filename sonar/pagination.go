package sonar

import (
	"context"
	"fmt"
	"net/http"
)

// allPages fetches every page from a V1 paginated endpoint, accumulating items
// into a single slice. page and pageSize are pointers into the caller's options
// struct so the helper can advance pagination without an extra allocation.
//
// fetch wraps the single-page call and extracts (items, total, response, err).
// If the context is cancelled between pages, allPages returns the items
// collected so far together with the context error.
func allPages[T any](
	ctx context.Context,
	page *int64,
	pageSize *int64,
	fetch func(context.Context) ([]T, int64, *http.Response, error),
) ([]T, *http.Response, error) {
	*page = 1

	if *pageSize == 0 {
		*pageSize = MaxPageSize
	}

	var all []T

	var resp *http.Response

	for {
		ctxErr := ctx.Err()
		if ctxErr != nil {
			return all, resp, fmt.Errorf("%w", ctxErr)
		}

		items, total, r, err := fetch(ctx)
		resp = r

		if err != nil {
			ctxErr = ctx.Err()
			if ctxErr != nil {
				return all, resp, fmt.Errorf("%w", ctxErr)
			}

			return nil, resp, err
		}

		all = append(all, items...)

		if int64(len(all)) >= total {
			return all, resp, nil
		}

		*page++
	}
}
