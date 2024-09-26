package sdk

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"

	"github.com/konstructio/dropkick/internal/civo/sdk/json"
)

// GetByID gets a resource by ID as provided by the value in "resource".
// The ID value in resource must not be empty.
func GetByID[T Resource](ctx context.Context, c Civoer, resource *T) error {
	params := map[string]string{"region": c.GetRegion()}

	endpoint := (*resource).GetAPIEndpoint()
	name := (*resource).GetResourceType()
	id := (*resource).GetID()

	if id == "" {
		return fmt.Errorf("the ID field in the resource %s is empty", name)
	}

	fullpath := path.Join(endpoint, id)
	if err := c.Do(ctx, fullpath, http.MethodGet, resource, params); err != nil {
		if errors.Is(err, &json.HTTPError{Code: http.StatusNotFound}) {
			return ErrNotFound
		}

		return fmt.Errorf("unable to get item: %w", err)
	}

	return nil
}
