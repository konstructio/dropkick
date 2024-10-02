package sdk

import (
	"context"
	"testing"

	"github.com/konstructio/dropkick/internal/civo/sdk/json"
)

func Test_Delete(t *testing.T) {
	t.Run("successfully delete a resource", func(t *testing.T) {
		ctx := context.TODO()
		expectedInstanceID := "1"

		c := &MockCivo{
			FnDo: func(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
				// ensure context is being passed
				assertEqual(t, ctx, context.TODO())

				// ensure the appropriate endpoint is being called
				assertEqual(t, location, Instance{}.GetAPIEndpoint()+"/"+expectedInstanceID)

				// ensure the region is being passed
				assertEqual(t, params["region"], "lon1")

				return nil
			},
			FnGetRegion: func() string { return "lon1" },
		}

		// We create an instance with this ID, then request it to be
		// deleted via the API.
		instance := Instance{ID: expectedInstanceID}

		// Call the function we are testing.
		assertNoError(t, Delete(ctx, c, instance))
	})

	t.Run("fail to delete a resource with an empty ID", func(t *testing.T) {
		ctx := context.TODO()

		// We don't need any methods to be called on the Civoer interface.
		c := &MockCivo{}

		// We create an instance with an empty ID, then request it to be
		// deleted via the API.
		instance := Instance{}

		// Call the function we are testing.
		err := Delete(ctx, c, instance)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})

	t.Run("fail to delete a resource due to an API error", func(t *testing.T) {
		ctx := context.TODO()
		expectedInstanceID := "1"

		c := &MockCivo{
			FnDo: func(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
				return &json.HTTPError{Code: 500}
			},
			FnGetRegion: func() string { return "lon1" },
		}

		// We create an instance with this ID, then request it to be
		// deleted via the API.
		instance := Instance{ID: expectedInstanceID}

		// Call the function we are testing.
		err := Delete(ctx, c, instance)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})
}
