package sdk

import (
	"context"
	"fmt"
	"strconv"
)

// GetAll returns all resources of a given type.
func GetAll[T Resource](ctx context.Context, c Civoer) ([]T, error) {
	var item T

	if item.IsSinglePaged() {
		return getSinglePage[T](ctx, c, item.GetAPIEndpoint())
	}

	return getPaginated[T](ctx, c, item.GetAPIEndpoint())
}

// getPaginated is a helper function to get results off an API endpoint that
// supports pagination using "page" and "perPage" query parameters.
func getPaginated[T Resource](ctx context.Context, c Civoer, endpoint string) ([]T, error) {
	var totalItems []T

	for page := 1; ; page++ {
		params := map[string]string{
			"page":     strconv.Itoa(page),
			"per_page": "100",
			"region":   c.GetRegion(),
		}

		var resp struct {
			Page    int `json:"page"`
			PerPage int `json:"per_page"`
			Pages   int `json:"pages"`
			Items   []T `json:"items"`
		}

		err := c.Do(ctx, endpoint, "GET", &resp, params)
		if err != nil {
			return nil, fmt.Errorf("unable to get items: %w", err)
		}

		totalItems = append(totalItems, resp.Items...)

		if resp.Page >= resp.Pages {
			break
		}
	}

	return totalItems, nil
}

// getSinglePage is a helper function to get results off an API endpoint that
// does not support pagination.
func getSinglePage[T Resource](ctx context.Context, c Civoer, endpoint string) ([]T, error) {
	var resp []T

	params := map[string]string{"region": c.GetRegion()}

	err := c.Do(ctx, endpoint, "GET", &resp, params)
	if err != nil {
		return nil, fmt.Errorf("unable to get items: %w", err)
	}

	return resp, nil
}
