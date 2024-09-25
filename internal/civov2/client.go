package civov2

import "net/http"

// customLogger is an interface that allows us to log messages.
type customLogger interface {
	Errorf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
}

// Client is a Civo client.
type Client struct {
	region string
	logger customLogger

	client    *http.Client
	requester *civoJSONClient
}

// Option is a functional option for the Client.
type Option func(*Client) error

// WithLogger is an option to set a custom logger.
func WithLogger(logger customLogger) Option {
	return func(c *Client) error {
		c.logger = logger
		return nil
	}
}

// WithJSONClient is an option to set a custom JSON client.
// The client can be nil, in which case http.DefaultClient will be used.
// The endpoint is the base URL for the Civo API.
// The bearerToken is the token to authenticate with the Civo API.
func WithJSONClient(client *http.Client, endpoint, bearerToken string) Option {
	return func(c *Client) error {
		c.requester = newCivoJSONClient(client, endpoint, bearerToken)
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
