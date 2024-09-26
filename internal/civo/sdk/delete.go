package sdk

import (
	"context"
	"fmt"
	"net/http"
	"path"
)

// Delete removes a resource from the Civo API. The resource must have
// a non-empty ID.
func Delete[T Resource](ctx context.Context, c Civoer, resource T) error {
	if resource.GetID() == "" {
		return fmt.Errorf("the ID field in the resource %s is empty", resource.GetResourceType())
	}

	params := map[string]string{"region": c.GetRegion()}

	var output struct {
		ID     string `json:"id"`
		Result string `json:"result"`
	}

	fullpath := path.Join(resource.GetAPIEndpoint(), resource.GetID())
	if err := c.Do(ctx, fullpath, http.MethodDelete, &output, params); err != nil {
		return fmt.Errorf("unable to delete item: %w", err)
	}

	return nil
}
