package sdk

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/konstructio/dropkick/internal/civo/sdk/testutils"
)

func Test_New(t *testing.T) {
	t.Run("successfully create a new client", func(t *testing.T) {
		endpoint := "https://example.com"
		token := "token"
		region := "abc123"

		c, err := New(
			WithJSONClient(nil, endpoint, token),
			WithRegion(region),
		)

		testutils.AssertNoError(t, err)
		testutils.AssertEqual(t, c.GetRegion(), region)
		testutils.AssertEqual(t, c.client, http.DefaultClient)
		testutils.AssertEqual(t, c.requester.GetEndpoint(), endpoint)
	})

	t.Run("fail to create a new client", func(t *testing.T) {
		fakeErr := errors.New("fake error")

		fakeOpt := func(c *Client) error {
			return fakeErr
		}

		_, err := New(fakeOpt)
		testutils.AssertErrorEqual(t, fakeErr, err)
	})

	t.Run("test \"do\" method", func(t *testing.T) {
		fakeErr := errors.New("fake error")

		client := &Client{
			client: http.DefaultClient,
			requester: &testutils.MockCivo{
				FnDo: func(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
					return fakeErr
				},
			},
		}

		err := client.Do(context.Background(), "", "", nil, nil)
		testutils.AssertErrorEqual(t, fakeErr, err)
	})
}
