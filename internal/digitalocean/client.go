package digitalocean

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/konstructio/dropkick/internal/logger"
)

// DigitalOcean is a client for the DigitalOcean API.
type DigitalOcean struct {
	client  *godo.Client    // The underlying DigitalOcean API client.
	context context.Context // The context for API requests.
	nuke    bool            // Whether to nuke resources.
	token   string          // The API token.
	logger  *logger.Logger  // The logger instance.
}

// DigitalOceanOption is a function that configures a DigitalOcean.
type DigitalOceanOption func(*DigitalOcean) error

// WithLogger sets the logger for a DigitalOcean.
func WithLogger(logger *logger.Logger) DigitalOceanOption {
	return func(c *DigitalOcean) error {
		c.logger = logger
		return nil
	}
}

// WithToken sets the API token for a DigitalOcean.
func WithToken(token string) DigitalOceanOption {
	return func(c *DigitalOcean) error {
		c.token = token
		return nil
	}
}

// WithNuke sets whether to nuke resources for a DigitalOcean.
func WithNuke(nuke bool) DigitalOceanOption {
	return func(c *DigitalOcean) error {
		c.nuke = nuke
		return nil
	}
}

// WithContext sets the context for a DigitalOcean.
func WithContext(ctx context.Context) DigitalOceanOption {
	return func(c *DigitalOcean) error {
		c.context = ctx
		return nil
	}
}

// New creates a new DigitalOcean with the given options.
// It returns an error if the token or region is not set, or if it fails to
// create the underlying DigitalOcean API client.
func New(opts ...DigitalOceanOption) (*DigitalOcean, error) {
	c := &DigitalOcean{}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("unable to apply option: %w", err)
		}
	}

	if c.token == "" {
		return nil, fmt.Errorf("required token not found")
	}

	// The DigitalOcean API client does not authenticate until a request is made,
	// so we're requesting the account but not using it to verify the token.
	godoClient := godo.NewFromToken(c.token)
	_, _, err := godoClient.Account.Get(c.context)
	if err != nil {
		return nil, fmt.Errorf("unable to authenticate DigitalOcean client: %w", err)
	}

	c.client = godoClient

	if c.context == nil {
		c.context = context.Background()
	}

	if c.logger == nil {
		c.logger = logger.None
	}

	return c, nil
}
