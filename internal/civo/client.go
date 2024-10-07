package civo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/konstructio/dropkick/internal/civo/sdk"
	"github.com/konstructio/dropkick/internal/logger"
)

const civoAPIURL = "https://api.civo.com"

// Client is the interface that wraps the basic Civo API client methods.
type Client interface {
	GetInstances(ctx context.Context) ([]sdk.Instance, error)
	GetFirewalls(ctx context.Context) ([]sdk.Firewall, error)
	GetVolumes(ctx context.Context) ([]sdk.Volume, error)
	GetKubernetesClusters(ctx context.Context) ([]sdk.KubernetesCluster, error)
	GetNetworks(ctx context.Context) ([]sdk.Network, error)
	GetObjectStores(ctx context.Context) ([]sdk.ObjectStore, error)
	GetObjectStoreCredentials(ctx context.Context) ([]sdk.ObjectStoreCredential, error)
	GetLoadBalancers(ctx context.Context) ([]sdk.LoadBalancer, error)
	GetSSHKeys(ctx context.Context) ([]sdk.SSHKey, error)
	Delete(ctx context.Context, resource sdk.APIResource) error
	Each(ctx context.Context, v sdk.APIResource, iterator func(sdk.APIResource) error) error
}

// Civo is a client for the Civo API.
type Civo struct {
	client     Client       // The underlying Civo API client.
	nuke       bool         // Whether to nuke resources.
	region     string       // The region for API requests.
	nameFilter string       // If set, only resources with a name containing this string will be deleted.
	token      string       // The API token.
	logger     customLogger // The logger instance.
	apiURL     string       // The URL for the Civo API.
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

// WithNameFilter sets the name filter for a Civo.
func WithNameFilter(nameFilter string) Option {
	return func(c *Civo) error {
		c.nameFilter = nameFilter
		return nil
	}
}

// customLogger is a custom logger interface.
type customLogger interface {
	Errorf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
}

// _ is a compile-time check to ensure that Civo implements
// the customLogger interface.
var _ customLogger = &logger.Logger{}

// debuggableHTTPClient is an HTTP client that logs requests if
// the HTTP_DEBUG environment variable is set.
var debuggableHTTPClient = &http.Client{
	Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if os.Getenv("HTTP_DEBUG") == "" {
			return http.DefaultTransport.RoundTrip(req)
		}

		log.Printf("Request: %s %s", req.Method, req.URL.String())
		return http.DefaultTransport.RoundTrip(req)
	}),
}

// roundTripperFunc is a function that implements the http.RoundTripper interface.
type roundTripperFunc func(*http.Request) (*http.Response, error)

// RoundTrip implements the http.RoundTripper interface.
func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
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

	if c.apiURL == "" {
		c.apiURL = civoAPIURL
	}

	if c.logger == nil {
		c.logger = logger.None
	}

	client, err := sdk.New(
		sdk.WithRegion(c.region),
		sdk.WithJSONClient(debuggableHTTPClient, c.apiURL, c.token),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create Civo client: %w", err)
	}

	c.client = client

	return c, nil
}
