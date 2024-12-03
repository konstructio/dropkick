package cloudflare

import (
	"errors"
	"fmt"
	"os"

	cloudflarego "github.com/cloudflare/cloudflare-go"
	"github.com/konstructio/dropkick/internal/logger"
)

// Cloudflare is a client for the Cloudflare API.
type Cloudflare struct {
	client    *cloudflarego.API // The underlying Cloudflare API client.
	nuke      bool              // Whether to nuke resources.
	token     string            // The API token.
	subdomain string            // The subdomain to delete records from.
	zoneID    string            // The zone ID for the domain.
	zoneName  string            // The domain to clean in cloudflare.
	logger    *logger.Logger    // The logger instance.
}

// Option is a function that configures a Cloudflare.
type Option func(*Cloudflare) error

// WithLogger sets the logger for a Cloudflare.
func WithLogger(logger *logger.Logger) Option {
	return func(c *Cloudflare) error {
		c.logger = logger
		return nil
	}
}

// WithZoneName sets the API token for a Cloudflare.
func WithZoneName(zoneName string) Option {
	return func(c *Cloudflare) error {
		c.zoneName = zoneName
		return nil
	}
}

// WithSubdomain sets the subdomain to filter on for aCloudflare.
func WithSubdomain(subdomain string) Option {
	return func(c *Cloudflare) error {
		c.subdomain = subdomain
		return nil
	}
}

// WithToken sets the API token for a Cloudflare.
func WithToken(token string) Option {
	return func(c *Cloudflare) error {
		c.token = token
		return nil
	}
}

// WithNuke sets whether to nuke resources for a Cloudflare.
func WithNuke(nuke bool) Option {
	return func(c *Cloudflare) error {
		c.nuke = nuke
		return nil
	}
}

// New creates a new Cloudflare with the given options.
// It returns an error if the token is not set, or if it fails to
// create the underlying Cloudflare API client.
func New(opts ...Option) (*Cloudflare, error) {
	c := &Cloudflare{}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("unable to apply option: %w", err)
		}
	}

	if c.token == "" {
		return nil, errors.New("required token not found for cloudflare client")
	}

	cloudflareAPI, err := cloudflarego.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
	if err != nil {
		return nil, fmt.Errorf("unable to authenticate Cloudflare client: %w", err)
	}

	zoneID, err := cloudflareAPI.ZoneIDByName(c.zoneName)
	if err != nil {
		return nil, fmt.Errorf("unable to get zone ID by name %q: %w", c.zoneName, err)
	}
	c.zoneID = zoneID

	c.client = cloudflareAPI

	if c.logger == nil {
		c.logger = logger.None
	}

	return c, nil
}
