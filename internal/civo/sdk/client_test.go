package sdk

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func assertEqual[T comparable](t *testing.T, v1, v2 T) {
	t.Helper()

	if v1 != v2 {
		t.Fatalf("expected \"%v\", got \"%v\"", v1, v2)
	}
}

func assertEqualf[T comparable](t *testing.T, v1, v2 T, format string, args ...interface{}) {
	t.Helper()

	if v1 != v2 {
		t.Fatalf(format, args...)
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func assertErrorEqual(t *testing.T, expected error, got error) {
	t.Helper()

	if !errors.Is(expected, got) {
		t.Fatalf("expected error to be \"%#v\", got \"%#v\"", expected, got)
	}
}

// MockCivoer implements the JSONClient and Civoer interfaces.
var (
	_ JSONClient = &MockCivo{}
	_ Civoer     = &MockCivo{}
)

// MockCivo is a mock implementation of the Civoer interface.
type MockCivo struct {
	FnDo          func(ctx context.Context, location, method string, output interface{}, params map[string]string) error
	FnGetEndpoint func() string
	FnGetRegion   func() string
	FnGetClient   func() *http.Client
}

// Do is a mock implementation of the Civoer interface.
func (m *MockCivo) Do(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
	if m.FnDo == nil {
		return fmt.Errorf("method \"do\" not implemented")
	}

	return m.FnDo(ctx, location, method, output, params)
}

// GetClient is a mock implementation of the Civoer interface.
func (m *MockCivo) GetClient() *http.Client {
	if m.FnGetClient == nil {
		return nil
	}

	return m.FnGetClient()
}

// GetRegion is a mock implementation of the Civoer interface.
func (m *MockCivo) GetRegion() string {
	if m.FnGetRegion == nil {
		return ""
	}

	return m.FnGetRegion()
}

// GetEndpoint is a mock implementation of the Civoer interface.
func (m *MockCivo) GetEndpoint() string {
	if m.FnGetEndpoint == nil {
		return ""
	}

	return m.FnGetEndpoint()
}

func Test_New(t *testing.T) {
	t.Run("successfully create a new client", func(t *testing.T) {
		endpoint := "https://example.com"
		token := "token"
		region := "abc123"

		c, err := New(
			WithJSONClient(nil, endpoint, token),
			WithRegion(region),
		)

		assertNoError(t, err)
		assertEqual(t, c.GetRegion(), region)
		assertEqual(t, c.client, http.DefaultClient)
		assertEqual(t, c.requester.GetEndpoint(), endpoint)
	})

	t.Run("fail to create a new client", func(t *testing.T) {
		fakeErr := errors.New("fake error")

		fakeOpt := func(c *Client) error {
			return fakeErr
		}

		_, err := New(fakeOpt)
		assertErrorEqual(t, fakeErr, err)
	})

	t.Run("test \"do\" method", func(t *testing.T) {
		fakeErr := errors.New("fake error")

		client := &Client{
			client: http.DefaultClient,
			requester: &MockCivo{
				FnDo: func(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
					return fakeErr
				},
			},
		}

		err := client.Do(context.Background(), "", "", nil, nil)
		assertErrorEqual(t, fakeErr, err)
	})
}
