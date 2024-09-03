package civo

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/civo/civogo"
	"github.com/konstructio/dropkick/internal/logger"
)

// Civo is a client for the Civo API.
type Civo struct {
	client     *civogo.Client  // The underlying Civo API client.
	context    context.Context // The context for API requests.
	nuke       bool            // Whether to nuke resources.
	region     string          // The region for API requests.
	nameFilter *regexp.Regexp  // If set, only resources with a name matching the regexp will be deleted.
	token      string          // The API token.
	logger     *logger.Logger  // The logger instance.
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
		reFilter, err := regexp.Compile(nameFilter)
		if err != nil {
			return fmt.Errorf("unable to compile name filter regexp %q: %w", nameFilter, err)
		}

		c.nameFilter = reFilter
		return nil
	}
}

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

	civoClient, err := civogo.NewClient(c.token, c.region)
	if err != nil {
		return nil, fmt.Errorf("unable to create new client: %w", err)
	}

	c.client = civoClient

	if c.context == nil {
		c.context = context.Background()
	}

	if c.logger == nil {
		c.logger = logger.None
	}

	return c, nil
}
