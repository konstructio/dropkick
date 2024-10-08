package sdk

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/konstructio/dropkick/internal/civo/sdk/json"
)

// Civoer is the interface that represents a high-level Civo client. It differs
// from the JSONClient interface in that this is a higher-level abstraction that
// is used to interact with the Civo API.
type Civoer interface {
	Do(ctx context.Context, location, method string, output interface{}, params map[string]string) error
	GetRegion() string
}

// JSONClient is an interface that allows us to make requests to the Civo API.
// It's used to convey the low-level SDK client that can use generics to make
// requests to the Civo API.
type JSONClient interface {
	Do(ctx context.Context, location, method string, output interface{}, params map[string]string) error
	GetClient() *http.Client
	GetEndpoint() string
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
	// Defining a baseline timeout for the client.
	timeoutLimit := 30 * time.Second

	c := &Client{
		// Setting an opinionated client with a 30-second timeout.
		client: &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   timeoutLimit, // how long to wait for the connection to be established
					KeepAlive: timeoutLimit, // how long to keep the connections open
				}).DialContext,
				TLSHandshakeTimeout: 10 * time.Second, // how long to wait for the TLS handshake to complete
				IdleConnTimeout:     timeoutLimit,     // how long to keep idle connections open
			},
			Timeout: timeoutLimit * 2, // how long to wait for the request to complete
		},
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
