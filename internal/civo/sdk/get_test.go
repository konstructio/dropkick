package sdk

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/konstructio/dropkick/internal/civo/sdk/json"
	"github.com/konstructio/dropkick/internal/civo/sdk/testutils"
)

func Test_getByID(t *testing.T) {
	t.Run("successful fetch of an instance by ID", func(t *testing.T) {
		ctx := context.TODO()

		expectedID := "123"
		expectedLocation := "/v2/instances/" + expectedID
		expectedMethod := http.MethodGet
		expectedInstanceName := "test-instance"
		expectedInstanceStatus := "ACTIVE"

		c := &testutils.MockCivo{
			FnDo: func(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
				// location should be /v2/instances/123 since we're getting ID 123
				testutils.AssertEqual(t, location, expectedLocation)

				// method should be GET since we're fetching an instance
				testutils.AssertEqual(t, method, expectedMethod)

				// params should contain the region
				testutils.AssertEqual(t, params["region"], "lon1")

				instance, ok := output.(*Instance)
				if !ok {
					t.Fatalf("expected output to be of type *Instance, got %T", output)
				}

				// instance should be populated with the response from the API
				testutils.AssertEqual(t, instance.ID, "123")

				// attach the custom values to it
				instance.Name = expectedInstanceName
				instance.Status = expectedInstanceStatus

				return nil
			},

			FnGetRegion: func() string {
				return "lon1"
			},
		}

		instance := Instance{ID: expectedID}
		err := getByID(ctx, c, &instance)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// instance should be populated with the custom values
		testutils.AssertEqual(t, instance.Name, expectedInstanceName)
		testutils.AssertEqual(t, instance.Status, expectedInstanceStatus)
	})

	t.Run("instance by ID not found", func(t *testing.T) {
		ctx := context.TODO()

		expectedID := "123"
		expectedLocation := "/v2/instances/" + expectedID
		expectedMethod := http.MethodGet

		c := &testutils.MockCivo{
			FnDo: func(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
				// location should be /v2/instances/123 since we're getting ID 123
				testutils.AssertEqual(t, location, expectedLocation)

				// method should be GET since we're fetching an instance
				testutils.AssertEqual(t, method, expectedMethod)

				// params should contain the region
				testutils.AssertEqual(t, params["region"], "lon1")

				return &json.HTTPError{Code: http.StatusNotFound}
			},

			FnGetRegion: func() string {
				return "lon1"
			},
		}

		instance := Instance{ID: expectedID}
		err := getByID(ctx, c, &instance)
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("expected error to be ErrNotFound, got %v", err)
		}
	})

	t.Run("empty ID field", func(t *testing.T) {
		ctx := context.TODO()

		c := &testutils.MockCivo{
			// FnDo doesn't get called because we fail before we make
			// the call to it.
			// If we were to pass, the mock returns an error too, which
			// will be caught by the error checking type below.

			FnGetRegion: func() string {
				return "lon1"
			},
		}

		instance := Instance{}
		err := getByID(ctx, c, &instance)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		if !errors.Is(err, &EmptyIDError{}) {
			t.Fatalf("expected error to be of type EmptyIDError, got %v", err)
		}
	})

	t.Run("unexpected error", func(t *testing.T) {
		ctx := context.TODO()
		errUnexpected := errors.New("unexpected error")

		c := &testutils.MockCivo{
			FnDo: func(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
				return errUnexpected
			},

			FnGetRegion: func() string {
				return "lon1"
			},
		}

		instance := Instance{ID: "123"}
		err := getByID(ctx, c, &instance)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		if !errors.Is(err, errUnexpected) {
			t.Fatalf("expected error to be %v, got %v", errUnexpected, err)
		}
	})
}
