package civo

import (
	"context"
	"errors"
	"testing"

	"github.com/konstructio/dropkick/internal/civo/sdk"
	"github.com/konstructio/dropkick/internal/civo/sdk/testutils"
	"github.com/konstructio/dropkick/internal/logger"
)

func TestIterator(t *testing.T) {
	t.Run("successfully delete resource", func(t *testing.T) {
		deleteCalled := false

		mock := &mockClient{
			fnDelete: func(ctx context.Context, resource sdk.APIResource) error {
				deleteCalled = true
				return nil
			},
		}

		c := &Civo{
			client: mock,
			logger: logger.None,
			nuke:   true, // to truly call delete, nuke must be set to true
		}

		iterFunc := c.deleteIterator(context.Background())

		instance := sdk.Instance{
			ID:   "123",
			Name: "test-instance",
		}

		err := iterFunc(instance)
		testutils.AssertNoErrorf(t, err, "expected no error when calling iterator, got %v", err)
		testutils.AssertEqualf(t, true, deleteCalled, "expected delete to be called, got %v", deleteCalled)
	})

	t.Run("delete should not be called if nuke isn't set", func(t *testing.T) {
		deleteCalled := false

		mock := &mockClient{
			fnDelete: func(ctx context.Context, resource sdk.APIResource) error {
				deleteCalled = true
				return nil
			},
		}

		c := &Civo{
			client: mock,
			logger: logger.None,
			nuke:   false,
		}

		iterFunc := c.deleteIterator(context.Background())

		instance := sdk.Instance{
			ID:   "123",
			Name: "test-instance",
		}

		err := iterFunc(instance)
		testutils.AssertNoErrorf(t, err, "expected no error when calling iterator, got %v", err)
		testutils.AssertEqualf(t, false, deleteCalled, "expected delete to not be called, got %v", deleteCalled)
	})

	t.Run("delete should not be called if name filter doesn't match", func(t *testing.T) {
		deleteCalled := false

		mock := &mockClient{
			fnDelete: func(ctx context.Context, resource sdk.APIResource) error {
				deleteCalled = true
				return nil
			},
		}

		c := &Civo{
			client:     mock,
			logger:     logger.None,
			nuke:       true,
			nameFilter: "test-instance",
		}

		iterFunc := c.deleteIterator(context.Background())

		instance := sdk.Instance{
			ID:   "123",
			Name: "other-instance",
		}

		err := iterFunc(instance)
		testutils.AssertNoErrorf(t, err, "expected no error when calling iterator, got %v", err)
		testutils.AssertEqualf(t, false, deleteCalled, "expected delete to not be called, got %v", deleteCalled)
	})

	t.Run("delete should return error if delete fails", func(t *testing.T) {
		madeUpError := errors.New("made up!")

		mock := &mockClient{
			fnDelete: func(ctx context.Context, resource sdk.APIResource) error {
				return madeUpError
			},
		}

		c := &Civo{
			client: mock,
			logger: logger.None,
			nuke:   true,
		}

		iterFunc := c.deleteIterator(context.Background())

		instance := sdk.Instance{
			ID:   "123",
			Name: "test-instance",
		}

		err := iterFunc(instance)
		testutils.AssertErrorEqual(t, madeUpError, err)
	})
}
