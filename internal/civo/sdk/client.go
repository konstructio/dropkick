package sdk

import (
	"context"
	"net/http"

	"github.com/konstructio/dropkick/internal/civo/sdk/json"
)

// JSONClient is an interface that allows us to make requests to the Civo API.
type JSONClient interface {
	Do(ctx context.Context, location, method string, output interface{}, params map[string]string) error
}

// Client is a Civo client.
type Client struct {
	region string

	client    *http.Client
	requester JSONClient
}

// Option is a functional option for the Client.
type Option func(*Client) error

// WithJSONClient is an option to set a custom JSON client.
// The client can be nil, in which case http.DefaultClient will be used.
// The endpoint is the base URL for the Civo API.
// The bearerToken is the token to authenticate with the Civo API.
func WithJSONClient(client *http.Client, endpoint, bearerToken string) Option {
	return func(c *Client) error {
		c.requester = json.New(client, endpoint, bearerToken)
		return nil
	}
}

// WithRegion is an option to set the region.
func WithRegion(region string) Option {
	return func(c *Client) error {
		c.region = region
		return nil
	}
}

// New creates a new Civo client.
func New(opts ...Option) (*Client, error) {
	c := &Client{
		client: http.DefaultClient,
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// GetRegion returns the region of the client.
func (c *Client) GetRegion() string {
	return c.region
}

// Do wraps the underlying JSONClient's Do method.
func (c *Client) Do(ctx context.Context, location, method string, output interface{}, params map[string]string) error {
	return c.requester.Do(ctx, location, method, output, params) //nolint:wrapcheck // we control the downstream error too
}
