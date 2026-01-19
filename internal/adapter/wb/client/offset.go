package client

import "context"

type OffsetPage[T any] struct {
	Items   []T
	HasMore bool
	// или Total int, если WB отдаёт total
}

func fetchAllByOffset[T any](
	ctx context.Context,
	limit int,
	fetch func(ctx context.Context, offset, limit int) (OffsetPage[T], error),
) ([]T, error) {
	if limit <= 0 {
		limit = 100
	}

	var (
		out    []T
		offset = 0
	)

	for {
		page, err := fetch(ctx, offset, limit)
		if err != nil {
			return nil, err
		}

		out = append(out, page.Items...)

		// стоп условия:
		if !page.HasMore || len(page.Items) == 0 {
			break
		}

		offset += len(page.Items)
	}

	return out, nil
}
