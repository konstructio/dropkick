package sdk

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path"

	"github.com/konstructio/dropkick/internal/civo/sdk/json"
)

// EmptyIDError is returned when the ID field in a resource is empty.
type EmptyIDError struct {
	ResourceType string
}

// Error returns the error message, by implementing the error interface.
func (e *EmptyIDError) Error() string {
	return fmt.Sprintf("the ID field in the resource %q is empty", e.ResourceType)
}

// Is checks if the target error is an EmptyIDError.
func (e *EmptyIDError) Is(target error) bool {
	_, ok := target.(*EmptyIDError)
	return ok
}

// GetByID gets a resource by ID as provided by the value in "resource".
// The ID value in resource must not be empty.
func GetByID[T Resource](ctx context.Context, c Civoer, resource *T) error {
	params := map[string]string{"region": c.GetRegion()}

	endpoint := (*resource).GetAPIEndpoint()
	restype := (*resource).GetResourceType()
	id := (*resource).GetID()

	if id == "" {
		return &EmptyIDError{ResourceType: restype}
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
