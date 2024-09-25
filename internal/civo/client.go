package civo

import (
	"context"
	"errors"
	"fmt"

	"github.com/konstructio/dropkick/internal/civov2"
	"github.com/konstructio/dropkick/internal/logger"
)

const civoAPIURL = "https://api.civo.com"

// Civo is a client for the Civo API.
type Civo struct {
	client     *civov2.Client  // The underlying Civo API client.
	context    context.Context // The context for API requests.
	nuke       bool            // Whether to nuke resources.
	region     string          // The region for API requests.
	nameFilter string          // If set, only resources with a name containing this string will be deleted.
	token      string          // The API token.
	logger     customLogger    // The logger instance.
	apiURL     string          // The URL for the Civo API.
}

// Option is a function that configures a Civo.
type Option func(*Civo) error

// WithLogger sets the logger for a Civo.
func WithLogger(logger *logger.Logger) Option {
	return func(c *Civo) error {
		c.logger = logger
		return nil
	}
}

// WithToken sets the API token for a Civo.
func WithToken(token string) Option {
	return func(c *Civo) error {
		c.token = token
		return nil
	}
}

// WithRegion sets the region for a Civo.
func WithRegion(region string) Option {
	return func(c *Civo) error {
		c.region = region
		return nil
	}
}

// WithNuke sets whether to nuke resources for a Civo.
func WithNuke(nuke bool) Option {
	return func(c *Civo) error {
		c.nuke = nuke
		return nil
	}
}

// WithContext sets the context for a Civo.
func WithContext(ctx context.Context) Option {
	return func(c *Civo) error {
		c.context = ctx
		return nil
	}
}

// WithNameFilter sets the name filter for a Civo.
func WithNameFilter(nameFilter string) Option {
	return func(c *Civo) error {
		c.nameFilter = nameFilter
		return nil
	}
}

// WithAPIURL sets the API URL for a Civo.
func WithAPIURL(apiURL string) Option {
	return func(c *Civo) error {
		c.apiURL = apiURL
		return nil
	}
}

type customLogger interface {
	Errorf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
}

var _ customLogger = &logger.Logger{}

// New creates a new Civo with the given options.
// It returns an error if the token or region is not set, or if it fails to create the underlying Civo API client.
func New(opts ...Option) (*Civo, error) {
	c := &Civo{}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("unable to apply option: %w", err)
		}
	}

	if c.token == "" {
		return nil, errors.New("required token not found")
	}

	if c.region == "" {
		return nil, errors.New("required region not set")
	}

	if c.apiURL == "" {
		c.apiURL = civoAPIURL
	}

	if c.context == nil {
		c.context = context.Background()
	}

	if c.logger == nil {
		c.logger = logger.None
	}

	client, err := civov2.New(
		civov2.WithLogger(c.logger),
		civov2.WithRegion(c.region),
		civov2.WithJSONClient(nil, c.apiURL, c.token),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create Civo client: %w", err)
	}

	c.client = client

	return c, nil
}
